package ingest

import (
	"github.com/Utsavk/capture-events/config"
	es "github.com/Utsavk/capture-events/db/es"
	httputil "github.com/Utsavk/capture-events/utils/http"
	"github.com/kataras/golog"
)

//create a conn struct which will contain session id

type Session struct {
	esHttp *httputil.ReqParams
}

func CreateSession() Session {
	return Session{esHttp: config.ESConf.GetHttpReqObject()}
}

func (s *Session) Ingest(rawData []byte) string {
	//Dump into ES
	esInsertReq := &es.InsertReq{
		Index: config.ESConf.Indexes["Logs"],
		ReqCommon: es.ReqCommon{
			HttpReqObj: s.esHttp,
		},
		Doc: []byte(`{"a":"b"}`),
	}
	res, err := config.ESConf.Insert(esInsertReq)
	if err != nil || res.Error != nil {
		golog.Info(res.Error)
		return "rejected"
	}
	return "ok"
}
