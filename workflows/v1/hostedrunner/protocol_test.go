package hostedrunner

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
	"google.golang.org/protobuf/proto"
)

type shortWriter struct{ bytes.Buffer }

func (w *shortWriter) Write(data []byte) (int, error) { return w.Buffer.Write(data[:1]) }

func TestFramesBoundedAndShortIO(t *testing.T) {
	identity := Identity{AssignmentID: "assignment", AssignmentEpoch: 7, GuestBootID: "boot"}
	envelope := identity.Envelope(1)
	envelope.SetHello(&workflowsv1.HostedRunnerHello{})
	var writer shortWriter
	require.NoError(t, WriteFrame(&writer, envelope))
	got, err := ReadFrame(&writer)
	require.NoError(t, err)
	assert.True(t, proto.Equal(envelope, got))
	require.NoError(t, identity.Check(got))
	for _, size := range []uint32{0, MaxFrameSize + 1, ^uint32(0)} {
		var header [4]byte
		binary.BigEndian.PutUint32(header[:], size)
		_, err := ReadFrame(bytes.NewReader(header[:]))
		require.ErrorContains(t, err, "frame size")
	}
	var full bytes.Buffer
	require.NoError(t, WriteFrame(&full, envelope))
	for size := range full.Len() {
		_, err := ReadFrame(bytes.NewReader(full.Bytes()[:size]))
		require.Error(t, err)
	}
	envelope.SetGuestBootId(strings.Repeat("x", MaxFrameSize))
	require.Error(t, WriteFrame(io.Discard, envelope))
}

func TestIdentityRejectsStaleAndIncompatibleMessages(t *testing.T) {
	identity := Identity{AssignmentID: "assignment", AssignmentEpoch: 7, GuestBootID: "boot"}
	for _, mutate := range []func(*workflowsv1.HostedRunnerEnvelope){
		func(e *workflowsv1.HostedRunnerEnvelope) { e.SetProtocolVersion(2) },
		func(e *workflowsv1.HostedRunnerEnvelope) { e.SetAssignmentEpoch(6) },
		func(e *workflowsv1.HostedRunnerEnvelope) { e.SetAssignmentId("other") },
		func(e *workflowsv1.HostedRunnerEnvelope) { e.SetGuestBootId("old") },
		func(e *workflowsv1.HostedRunnerEnvelope) { e.SetSequence(0) },
	} {
		envelope := identity.Envelope(1)
		mutate(envelope)
		require.Error(t, identity.Check(envelope))
	}
}
