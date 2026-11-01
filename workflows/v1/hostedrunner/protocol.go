// Package hostedrunner implements the bounded, versioned local guest lifecycle stream.
// Socket identity/authentication belongs to the host, not to these caller-supplied IDs.
package hostedrunner

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"google.golang.org/protobuf/proto"
)

const (
	ProtocolVersion = 1
	MaxFrameSize    = 64 * 1024
)

// Identity must be bound to the assigned guest's socket by the host daemon.
type Identity struct {
	AssignmentID    string
	AssignmentEpoch uint64
	GuestBootID     string
}

func (i Identity) Validate() error {
	if i.AssignmentID == "" || i.AssignmentEpoch == 0 || i.GuestBootID == "" {
		return errors.New("hosted lifecycle requires assignment ID, nonzero epoch, and guest boot ID")
	}
	return nil
}

// Envelope constructs the identity header; the caller sets exactly one payload.
func (i Identity) Envelope(sequence uint64) *workflowsv1.HostedRunnerEnvelope {
	return workflowsv1.HostedRunnerEnvelope_builder{
		ProtocolVersion: ProtocolVersion, AssignmentId: i.AssignmentID,
		AssignmentEpoch: i.AssignmentEpoch, GuestBootId: i.GuestBootID, Sequence: sequence,
	}.Build()
}

// Check validates version and identity before processing any payload or sequence.
func (i Identity) Check(envelope *workflowsv1.HostedRunnerEnvelope) error {
	if envelope.GetProtocolVersion() != ProtocolVersion {
		return errors.New("incompatible hosted lifecycle version")
	}
	if envelope.GetAssignmentId() != i.AssignmentID || envelope.GetAssignmentEpoch() != i.AssignmentEpoch || envelope.GetGuestBootId() != i.GuestBootID {
		return errors.New("stale or mismatched hosted lifecycle identity")
	}
	if envelope.GetSequence() == 0 {
		return errors.New("hosted lifecycle sequence must be nonzero")
	}
	return nil
}

// ReadFrame rejects oversized lengths before allocating the payload.
func ReadFrame(reader io.Reader) (*workflowsv1.HostedRunnerEnvelope, error) {
	var header [4]byte
	if _, err := io.ReadFull(reader, header[:]); err != nil {
		return nil, err
	}
	size := binary.BigEndian.Uint32(header[:])
	if size == 0 || size > MaxFrameSize {
		return nil, fmt.Errorf("invalid hosted lifecycle frame size %d", size)
	}
	data := make([]byte, size)
	if _, err := io.ReadFull(reader, data); err != nil {
		return nil, err
	}
	envelope := &workflowsv1.HostedRunnerEnvelope{}
	if err := proto.Unmarshal(data, envelope); err != nil {
		return nil, err
	}
	return envelope, nil
}

// WriteFrame must have a single writer per stream. Short writes are supported.
func WriteFrame(writer io.Writer, envelope *workflowsv1.HostedRunnerEnvelope) error {
	if size := proto.Size(envelope); size == 0 || size > MaxFrameSize {
		return errors.New("invalid hosted lifecycle message size")
	}
	data, err := proto.Marshal(envelope)
	if err != nil {
		return err
	}
	if len(data) > MaxFrameSize {
		return errors.New("invalid hosted lifecycle message size")
	}
	frame := make([]byte, 4, 4+len(data))
	binary.BigEndian.PutUint32(frame, uint32(len(data))) //nolint:gosec // Length is bounded to 64 KiB above.
	frame = append(frame, data...)
	for len(frame) > 0 {
		n, err := writer.Write(frame)
		if err != nil {
			return err
		}
		if n <= 0 || n > len(frame) {
			return io.ErrShortWrite
		}
		frame = frame[n:]
	}
	return nil
}
