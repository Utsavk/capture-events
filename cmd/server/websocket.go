package server

import (
	"net/http"

	ingestPkg "github.com/Utsavk/capture-events/ingest"
	"github.com/fasthttp/websocket"
	"github.com/kataras/golog"
	"github.com/valyala/fasthttp"
)

var upgrader = websocket.FastHTTPUpgrader{
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool { return true },
}

func websocketHandler(ws *websocket.Conn) {
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			golog.Error(err)
			return
		}
		ingestPkg.Ingest(message)
		ws.WriteMessage(mt, []byte("ok"))
	}
}

func onWSHandShake(ctx *fasthttp.RequestCtx) error {
	return nil
}

func onWSRequest(ctx *fasthttp.RequestCtx) (string, int) {
	err := onWSHandShake(ctx)
	if err != nil {
		ctx.Response.SetConnectionClose()
		return "handshake failure", http.StatusBadRequest
	}
	err = upgrader.Upgrade(ctx, websocketHandler)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); ok {
			golog.Error(err)
		}
		return "connection upgrade failure", http.StatusBadRequest
	}
	return "", http.StatusOK
}
