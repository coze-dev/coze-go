package coze

import (
	"bufio"
	"io"
	"net/http"
	"strings"

	"github.com/coze-dev/coze-go/log"
)

type streamable interface {
	ChatEvent | WorkflowEvent
}

type eventProcessor[T streamable] func(line []byte, reader *bufio.Reader) (*T, bool, error)

type streamReader[T streamable] struct {
	isFinished bool

	reader       *bufio.Reader
	response     *http.Response
	processor    eventProcessor[T]
	httpResponse *httpResponse
}

func (s *streamReader[T]) Recv() (response *T, err error) {
	return s.processLines()
}

//nolint:gocognit
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
		event, isDone, err := s.processor(line, s.reader)
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
			log.Warnf("Error reading response body: ", err)
			return err
		}
		return isResponseSuccess(&baseResponse{}, respStr, s.httpResponse)
	}
	return nil
}

func (s *streamReader[T]) Close() error {
	return s.response.Body.Close()
}

func (s *streamReader[T]) HTTPResponse() HTTPResponse {
	return s.httpResponse
}
