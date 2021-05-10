package server

import (
	"bytes"
	"net/http"

	"github.com/Utsavk/capture-events/config"
	"github.com/kataras/golog"
	"github.com/valyala/fasthttp"
)

func StartServer(conf config.Server) {
	if err := fasthttp.ListenAndServe(conf.Address, requestHandler); err != nil {
		golog.Fatal("Error in ListenAndServe: %s", err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	if bytes.Equal(ctx.Path(), health) {
		sendResponse(ctx, "ok", http.StatusOK)
	} else if bytes.Equal(ctx.Path(), ingest) {
		msg, code := onWSRequest(ctx)
		if msg != "" {
			sendResponse(ctx, msg, code)
		}
	} else {
		sendResponse(ctx, "invalid route", http.StatusNotFound)
	}
}
