package main

import (
	"net/http"

	"github.com/smackem/ylang/internal/goobar"
)

func registerHTTP() {
	http.Handle("/open", goobar.Get(func(x *goobar.Exchange) goobar.Responder {
		return goobar.Nop()
	}))
}
