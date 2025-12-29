package main

import (
	"bytes"
	"log"

	coap "github.com/plgd-dev/go-coap/v3"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/mux"
)

func loggingMiddleware(next mux.Handler) mux.Handler {
	return mux.HandlerFunc(
		func(w mux.ResponseWriter, resp *mux.Message) {
			log.Printf("ClientAddress %v, %v\n", w.Conn().RemoteAddr(), resp.String())
			next.ServeCOAP(w, resp)
		},
	)
}

func handleHello(writer mux.ResponseWriter, resp *mux.Message) {
	customResp := writer.Conn().AcquireMessage(resp.Context())
	defer writer.Conn().ReleaseMessage(customResp)
	customResp.SetCode(codes.Content)
	customResp.SetToken(resp.Token())
	customResp.SetContentFormat(message.TextPlain)
	customResp.SetBody(bytes.NewReader([]byte("Hello world")))
	err := writer.Conn().WriteMessage(customResp)
	if err != nil {
		log.Printf("cannot set response: %v", err)
	}
}

func main() {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.Handle("/api/hello", mux.HandlerFunc(handleHello))

	log.Fatal(coap.ListenAndServe("udp", ":5688", r))
}
