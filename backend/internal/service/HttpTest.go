package service

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

type HttpTestReq struct {
	Method string            `json:"method"` // GET / POST / PUT / DELETE
	URL    string            `json:"url"`
	Header map[string]string `json:"header"`
	Body   string            `json:"body"`
}

type HttpTestResp struct {
	Status     string              `json:"status"`
	StatusCode int                 `json:"status_code"`
	Header     map[string][]string `json:"header"`
	Body       string              `json:"body"`
}

type HttpTestService struct {
	client *http.Client
}

func NewHttpTestService() *HttpTestService {
	return &HttpTestService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *HttpTestService) Do(req *HttpTestReq) (*HttpTestResp, error) {
	if strings.TrimSpace(req.URL) == "" {
		return nil, errors.New("url is empty")
	}

	method := strings.ToUpper(strings.TrimSpace(req.Method))
	if method == "" {
		method = http.MethodGet
	}

	var bodyReader io.Reader
	if req.Body != "" && method != http.MethodGet {
		bodyReader = strings.NewReader(req.Body)
	}

	httpReq, err := http.NewRequest(method, req.URL, bodyReader)
	if err != nil {
		return nil, err
	}

	for k, v := range req.Header {
		httpReq.Header.Set(k, v)
	}

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &HttpTestResp{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       string(respBodyBytes),
	}, nil
}
