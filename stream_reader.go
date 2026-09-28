package coze

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"strings"
)

type streamable interface {
	ChatEvent | WorkflowEvent | NopEvent
}

type NopEvent struct{}

type Stream[T streamable] interface {
	Responser
	Close() error
	// Recv returns the next event of the stream. After the terminal event has
	// been returned, every further call reports io.EOF. Callers are still
	// responsible for calling Close.
	Recv() (*T, error)
}

type eventProcessor[T streamable] func(ctx context.Context, core *core, line []byte, reader *bufio.Reader) (*T, bool, error)

type streamReader[T streamable] struct {
	// un-mutable
	ctx          context.Context
	core         *core
	response     *http.Response
	httpResponse *httpResponse
	processor    eventProcessor[T]

	isFinished bool
	reader     *bufio.Reader
}

func newStream[T streamable](ctx context.Context, core *core, resp *http.Response, processor eventProcessor[T]) Stream[T] {
	if resp == nil {
		return nil
	}
	return &streamReader[T]{
		ctx:          ctx,
		core:         core,
		response:     resp,
		httpResponse: newHTTPResponse(resp),
		processor:    processor,
		reader:       bufio.NewReader(resp.Body),
	}
}

// Recv returns the next event of the stream. Once the terminal event has been
// delivered, the stream is finished and Recv reports io.EOF without reading
// further from the HTTP response.
func (s *streamReader[T]) Recv() (response *T, err error) {
	if s.isFinished {
		return nil, io.EOF
	}
	return s.processLines()
}

func (s *streamReader[T]) processLines() (*T, error) {
	err := s.checkRespErr()
	if err != nil {
		return nil, err
	}
	for {
		line, _, readErr := s.reader.ReadLine()
		if readErr != nil {
			return nil, readErr
		}

		if line == nil {
			s.isFinished = true
			break
		}
		if len(line) == 0 {
			continue
		}
		event, isDone, err := s.processor(s.ctx, s.core, line, s.reader)
		if err != nil {
			return nil, err
		}
		s.isFinished = isDone
		if event == nil {
			continue
		}
		return event, nil
	}
	return nil, io.EOF
}

func (s *streamReader[T]) checkRespErr() error {
	contentType := s.response.Header.Get("Content-Type")
	if contentType != "" && strings.Contains(contentType, "application/json") {
		respStr, err := io.ReadAll(s.response.Body)
		if err != nil {
			logger.Warnf(s.ctx, "Error reading response body: ", err)
			return err
		}
		return isResponseSuccess(s.ctx, &baseResponse{}, respStr, s.httpResponse)
	}
	return nil
}

func (s *streamReader[T]) Close() error {
	return s.response.Body.Close()
}

func (s *streamReader[T]) Response() HTTPResponse {
	return s.httpResponse
}
