package workflows

import (
	"context"
	"sync"

	"github.com/tilebox/tilebox-go/client"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
)

// PollingRunnerOption configures an explicitly opted-in polling runner.
type PollingRunnerOption func(*PollingTaskRunner)

// WithHostedLifecycle shares intake admission and aggregate observations across slots.
// Configure every slot before calling MarkReady. A lifecycle cannot be reused for a
// new guest boot; unresolved claims/reports deliberately have no time-based reset.
func WithHostedLifecycle(lifecycle *HostedLifecycle) PollingRunnerOption {
	return func(r *PollingTaskRunner) {
		if lifecycle == nil {
			return
		}
		lifecycle.mu.Lock()
		defer lifecycle.mu.Unlock()
		state := &hostedSlot{lifecycle: lifecycle}
		lifecycle.slots = append(lifecycle.slots, state)
		r.hosted = state
	}
}

// HostedLifecycle owns the StartIntake gate and irreversible normal drain. Its
// snapshot never waits for task execution, reporting, or a backend RPC.
type HostedLifecycle struct {
	mu                       sync.Mutex
	slots                    []*hostedSlot
	ready, started, draining bool
	wake                     chan struct{}
	activity                 uint64
}

type hostedSlot struct {
	lifecycle                                     *HostedLifecycle
	claim, active, pending                        bool
	claimUncertain, reportUncertain, leaseProblem bool
	abandoned                                     uint32
	backend                                       workflowsv1.HostedRunnerBackendState
}

// NewHostedLifecycle creates a closed intake gate.
func NewHostedLifecycle() *HostedLifecycle {
	return &HostedLifecycle{wake: make(chan struct{})}
}

// MarkReady records local initialization, not successful lazy runtime startup.
func (h *HostedLifecycle) MarkReady() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.ready = true
}

// StartIntake returns false before readiness or after any accepted drain.
func (h *HostedLifecycle) StartIntake() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	if !h.ready || h.draining {
		return false
	}
	if !h.started {
		h.started = true
		close(h.wake)
	}
	return true
}

// Drain atomically prevents new claims but does not cancel admitted RPCs or tasks.
func (h *HostedLifecycle) Drain() {
	h.mu.Lock()
	defer h.mu.Unlock()
	if !h.draining {
		h.draining = true
		if !h.started {
			close(h.wake)
		}
	}
}

// Snapshot returns an independent protobuf aggregate safe for concurrent marshaling.
// Empty polls do not increment ActivityGeneration. Received work and uncertain
// outcomes do, even if they are resolved between successive heartbeats.
func (h *HostedLifecycle) Snapshot() *workflowsv1.HostedRunnerStatus {
	h.mu.Lock()
	defer h.mu.Unlock()
	status := workflowsv1.HostedRunnerStatus_builder{
		Initialized: h.ready, IntakeOpen: h.started && !h.draining,
		Draining: h.draining, ActivityGeneration: h.activity,
	}.Build()
	for _, slot := range h.slots {
		if slot.claim {
			status.SetActiveClaims(status.GetActiveClaims() + 1)
		}
		if slot.active {
			status.SetActiveTasks(status.GetActiveTasks() + 1)
		}
		if slot.pending {
			status.SetPendingReports(status.GetPendingReports() + 1)
		}
		status.SetPendingReports(status.GetPendingReports() + slot.abandoned)
		status.SetClaimUncertain(status.GetClaimUncertain() || slot.claimUncertain)
		status.SetReportUncertain(status.GetReportUncertain() || slot.reportUncertain || slot.abandoned > 0)
		status.SetTaskLeaseProblem(status.GetTaskLeaseProblem() || slot.leaseProblem)
		if slot.backend > status.GetBackend() {
			status.SetBackend(slot.backend)
		}
	}
	settled := status.GetActiveClaims() == 0 && status.GetActiveTasks() == 0 && status.GetPendingReports() == 0 && !status.GetClaimUncertain() && !status.GetReportUncertain()
	status.SetIdle(h.ready && settled)
	status.SetSafeToStop(h.ready && h.draining && settled)
	return status
}

func (h *HostedLifecycle) wait(ctx context.Context) {
	select {
	case <-h.wake:
	case <-ctx.Done():
	}
}

func (h *HostedLifecycle) requesting() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.started && !h.draining
}

// nextTask serializes claim admission against Drain without holding a lock over I/O.
func (r *PollingTaskRunner) nextTask(ctx context.Context, computed *workflowsv1.ComputedTask, next *workflowsv1.NextTaskToRun) (*workflowsv1.NextTaskResponse, error) {
	slot := r.hosted
	if slot == nil {
		return r.service.NextTask(ctx, computed, next)
	}
	h := slot.lifecycle
	h.mu.Lock()
	if !h.started || h.draining {
		next = nil
	}
	slot.claim = next != nil
	h.mu.Unlock()
	if computed == nil && next == nil {
		return &workflowsv1.NextTaskResponse{}, nil
	}
	if next != nil {
		ctx = client.WithoutHTTPRetries(ctx)
	}
	response, err := r.service.NextTask(ctx, computed, next)
	h.mu.Lock()
	defer h.mu.Unlock()
	slot.claim = false
	slot.observeBackend(err)
	if err != nil {
		if next != nil && shouldRetryRPCAfterTimeout(err) {
			slot.claimUncertain = true
			h.activity++
		}
		if computed != nil {
			slot.reportUncertain = true
			h.activity++
		}
	} else {
		if computed != nil {
			slot.pending = false
			slot.reportUncertain = false
		}
		if response.GetNextTask() != nil {
			slot.active = true
			h.activity++
			if isEmpty(response.GetNextTask().GetId()) {
				slot.active = false
				slot.claimUncertain = true
			}
		}
	}
	return response, err
}

func (s *hostedSlot) observeBackend(err error) {
	s.backend = workflowsv1.HostedRunnerBackendState_HOSTED_RUNNER_BACKEND_STATE_REACHABLE
	if err != nil {
		s.backend = workflowsv1.HostedRunnerBackendState_HOSTED_RUNNER_BACKEND_STATE_UNREACHABLE
	}
}

func (s *hostedSlot) executed() {
	if s == nil {
		return
	}
	s.lifecycle.mu.Lock()
	defer s.lifecycle.mu.Unlock()
	s.active = false
	s.pending = true
}

func (s *hostedSlot) reported(err error) {
	if s == nil {
		return
	}
	s.lifecycle.mu.Lock()
	defer s.lifecycle.mu.Unlock()
	s.observeBackend(err)
	s.reportUncertain = err != nil
	if err == nil {
		s.pending = false
	}
}

func (s *hostedSlot) abandonReport() {
	if s == nil {
		return
	}
	s.lifecycle.mu.Lock()
	defer s.lifecycle.mu.Unlock()
	s.abandoned++
	s.pending = false
	s.reportUncertain = false
}

func (s *hostedSlot) leaseFailed() {
	if s == nil {
		return
	}
	s.lifecycle.mu.Lock()
	defer s.lifecycle.mu.Unlock()
	s.leaseProblem = true
}
