package http

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Utsavk/capture-events/utils"
)

type Method string

const (
	POST   Method = "POST"
	GET    Method = "GET"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
	HEAD   Method = "HEAD"
)

type ReqParams struct {
	Hosts              []string
	SSL                bool
	Path               string
	Method             Method
	Headers            map[string]string
	Form               url.Values
	TimeoutMS          int64
	InsecureSkipVerify bool
	RawBody            []byte
	MaxRetry           int
	retry              int
}

const ReqDefaultMaxRetry = 3

type Response struct {
	Body       []byte
	StatusCode int
	Headers    map[string][]string
}

func getClient(timeoutMS int64, insecureSkipVerify bool) *http.Client {
	var tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
	}
	if timeoutMS > 0 {
		return &http.Client{Transport: tr, Timeout: time.Duration(timeoutMS) * time.Millisecond}
	}
	return &http.Client{Transport: tr, Timeout: 15 * time.Millisecond}
}

func Send(params *ReqParams) (*Response, error) {
	if params.MaxRetry > 0 {
		if params.retry > params.MaxRetry {
			return nil, fmt.Errorf("maximum retry %d exhausted", params.MaxRetry)
		}
	} else if params.retry > ReqDefaultMaxRetry {
		return nil, fmt.Errorf("maximum retry %d exhausted", ReqDefaultMaxRetry)
	}
	var bodyReader io.Reader
	var hostsCopy []string

	if params.RawBody != nil {
		bodyReader = bytes.NewReader(params.RawBody)
	} else if params.Form != nil {
		bodyReader = strings.NewReader(params.Form.Encode())
	}

	if len(params.Hosts) == 1 {
		params.Path = params.Hosts[0] + "/" + params.Path
	} else if len(params.Hosts) > 1 {
		hostsCopy = make([]string, len(params.Hosts))
		copy(hostsCopy, params.Hosts)
		host, index := utils.GetRandomSliceObject(hostsCopy)
		params.Path = host + "/" + params.Path
		hostsCopy[index] = hostsCopy[len(hostsCopy)-1]
		params.Hosts = hostsCopy
	}

	if params.SSL {
		params.Path = fmt.Sprintf("https://%s", params.Path)
	} else {
		params.Path = fmt.Sprintf("http://%s", params.Path)
	}

	req, err := http.NewRequest(string(params.Method), params.Path, bodyReader)
	if err != nil {
		return nil, err
	}

	if params.Headers != nil {
		for k, v := range params.Headers {
			req.Header.Set(k, v)
		}
	}

	client := getClient(params.TimeoutMS, params.InsecureSkipVerify)
	rawRes, err := client.Do(req)
	var response Response
	if err != nil {
		params.retry++
		return Send(params)
	}

	body, err := ioutil.ReadAll(rawRes.Body)
	rawRes.Body.Close()
	if err != nil {
		return nil, err
	}
	response.Body = body
	response.StatusCode = rawRes.StatusCode
	response.Headers = rawRes.Header
	return &response, nil
}
