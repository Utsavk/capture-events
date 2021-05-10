package server

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type response struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

func sendResponse(ctx *fasthttp.RequestCtx, msg string, code int) {
	res := response{Msg: msg, Code: code}
	resBytes, _ := json.Marshal(res)
	ctx.SetBody(resBytes)
	ctx.SetStatusCode(code)
}
