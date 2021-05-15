package config

import (
	"encoding/json"
	"fmt"
	"net/http"

	httputil "github.com/Utsavk/capture-events/utils/http"
)

type Conf struct {
	Hosts        []string
	Indexes      map[string]string
	SSL          bool
	ReqTimeoutMS int64
}

type InsertReq struct {
	Index string
	ID    string
	Doc   []byte
	ReqCommon
}

type ReqCommon struct {
	ReqTimeoutMS int64
	HttpReqObj   *httputil.ReqParams
}

type ErrorResponse struct {
	Error interface{} `json:"error"`
}

type Response struct {
	Body       interface{}
	Error      interface{}
	StatusCode int
}

func (e *Conf) ValidateConf() error {
	if len(e.Hosts) == 0 {
		return fmt.Errorf("invalid ES config")
	}
	return nil
}

func (e *Conf) GetHttpReqObject() *httputil.ReqParams {
	return &httputil.ReqParams{
		Hosts:     e.Hosts,
		SSL:       e.SSL,
		Headers:   map[string]string{"Content-type": "Application/json"},
		TimeoutMS: e.ReqTimeoutMS,
	}
}

func (e *Conf) Insert(r *InsertReq) (*Response, error) {
	if r.HttpReqObj == nil {
		r.HttpReqObj = e.GetHttpReqObject()
	}
	httpReq := r.HttpReqObj
	httpReq.Method = httputil.POST
	httpReq.Path = fmt.Sprintf("%s/_doc/%s", r.Index, r.ID)
	if r.ReqTimeoutMS > 0 {
		httpReq.TimeoutMS = r.ReqTimeoutMS
	}
	httpReq.RawBody = r.Doc
	rawRes, err := httputil.Send(httpReq)
	if err != nil {
		return nil, err
	}
	if rawRes.StatusCode != http.StatusCreated && rawRes.StatusCode != http.StatusOK {
		var errRes ErrorResponse
		errMarshal := json.Unmarshal(rawRes.Body, &errRes)
		if errMarshal != nil {
			errRes.Error = "error in unmarshalling ES insert error node"
		}
		return &Response{Error: errRes, StatusCode: rawRes.StatusCode}, nil
	}
	return &Response{Body: rawRes.Body, StatusCode: rawRes.StatusCode}, nil
}
