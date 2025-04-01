package grpc // import "github.com/tilebox/tilebox-go/grpc"

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
)

type replayRPC struct {
	Body       []byte
	StatusCode int
}

var _ http.RoundTripper = &RecordRoundTripper{}

type RecordRoundTripper struct {
	writer io.Writer
}

func NewRecordRoundTripper(writer io.Writer) *RecordRoundTripper {
	return &RecordRoundTripper{
		writer: writer,
	}
}

func (rt *RecordRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	response, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// reset the response body with the original body
	response.Body = io.NopCloser(bytes.NewReader(body))

	replay := &replayRPC{
		Body:       body,
		StatusCode: response.StatusCode,
	}

	err = rt.writeReplay(replay)
	if err != nil {
		return nil, fmt.Errorf("failed to write replay: %w", err)
	}

	return response, nil
}

func (rt *RecordRoundTripper) writeReplay(replay *replayRPC) error {
	err := binary.Write(rt.writer, binary.LittleEndian, int64(replay.StatusCode))
	if err != nil {
		return fmt.Errorf("failed to write status code: %w", err)
	}

	err = binary.Write(rt.writer, binary.LittleEndian, uint64(len(replay.Body)))
	if err != nil {
		return fmt.Errorf("failed to write body size: %w", err)
	}

	_, err = rt.writer.Write(replay.Body)
	if err != nil {
		return fmt.Errorf("failed to write body: %w", err)
	}

	return nil
}

var _ http.RoundTripper = &ReplayRoundTripper{}

type ReplayRoundTripper struct {
	reader io.Reader
}

func NewReplayRoundTripper(reader io.Reader) *ReplayRoundTripper {
	return &ReplayRoundTripper{
		reader: reader,
	}
}

func (rt *ReplayRoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	replay, err := rt.readReplay()
	if err != nil {
		return nil, fmt.Errorf("failed to read replay: %w", err)
	}

	return &http.Response{
		StatusCode:    replay.StatusCode,
		Body:          io.NopCloser(bytes.NewReader(replay.Body)),
		ContentLength: int64(-1),
		Header: map[string][]string{
			"Accept-Encoding":  {"gzip"},
			"Content-Encoding": {"gzip"},
			"Content-Type":     {"application/proto"},
		},
	}, nil
}

func (rt *ReplayRoundTripper) readReplay() (*replayRPC, error) {
	var statusCode int64
	err := binary.Read(rt.reader, binary.LittleEndian, &statusCode)
	if err != nil {
		return nil, fmt.Errorf("failed to read status code: %w", err)
	}

	var size uint64
	err = binary.Read(rt.reader, binary.LittleEndian, &size)
	if err != nil {
		return nil, fmt.Errorf("failed to read body size: %w", err)
	}

	body := make([]byte, size)
	_, err = io.ReadFull(rt.reader, body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	return &replayRPC{
		Body:       body,
		StatusCode: int(statusCode),
	}, nil
}
