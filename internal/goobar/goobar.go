package goobar

import (
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

// Handler stores a record of functions that handle different HTTP methods
// invoked for a common path. It stores Actions for GET, POST, PUT, DELETE
// and OPTIONS.
type Handler struct {
	Get, Post, Put, Delete, Options Action
}

// Action is a function that handles an HTTP request.
type Action func(x *Exchange) Responder

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer recoverFromActionPanic(w, r)

	u, _ := url.QueryUnescape(r.URL.String())
	log.Printf("%s %s", r.Method, u)

	action := h.getAction(r.Method)
	if action == nil {
		http.NotFound(w, r)
		return
	}

	x := makeExchange(w, r)
	responder := action(x)
	if ct := responder.ContentType(); len(strings.TrimSpace(ct)) > 0 {
		w.Header().Set("Content-Type", ct)
	}
	responder.Respond(w)
}

func (h *Handler) getAction(method string) Action {
	switch method {
	case "GET":
		return h.Get
	case "POST":
		return h.Post
	case "PUT":
		return h.Put
	case "DELETE":
		return h.Delete
	case "OPTIONS":
		return h.Options
	default:
		return nil
	}
}

// Get returns a Handler with only one action for the HTTP method GET.
func Get(action Action) *Handler {
	return &Handler{Get: action}
}

// Post returns a Handler with only one action for the HTTP method Post.
func Post(action Action) *Handler {
	return &Handler{Post: action}
}

// Put returns a Handler with only one action for the HTTP method PUT.
func Put(action Action) *Handler {
	return &Handler{Put: action}
}

// Delete returns a Handler with only one action for the HTTP method DELETE.
func Delete(action Action) *Handler {
	return &Handler{Delete: action}
}

// AnyMethod returns a Handler with only one action for all supported HTTP actions
// (GET, POST, PUT, DELETE, OPTIONS)
func AnyMethod(action Action) *Handler {
	return &Handler{
		Get:     action,
		Post:    action,
		Put:     action,
		Delete:  action,
		Options: action,
	}
}

// SetViewFolder sets the file path of the local directory that contains
// the view templates.
func SetViewFolder(folder string) {
	viewFolder = filepath.Clean(folder)
}

// ViewFolder returns the file path of the local directory that contains
// the view templates.
func ViewFolder() string {
	return viewFolder
}

var viewFolder = "resource/view"

func recoverFromActionPanic(w http.ResponseWriter, r *http.Request) {
	switch x := recover(); p := x.(type) {
	case nil: // nothing thrown, ignore
	case actionPanic:
		u, _ := url.QueryUnescape(r.URL.String())
		log.Printf("%s: %s", u, p.msg)
		http.NotFound(w, r)
	default:
		panic(x) // rethrow
	}
}
