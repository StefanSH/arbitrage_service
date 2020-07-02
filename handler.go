package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

//import "context"

// NewHandler - create a new Handler
func NewHandler(conf *Tuner) *Handler {
	h := &Handler{
		Conf: conf,
		Log:  logrus.New(),
	}
	return h
}

// Handler structure
type Handler struct {
	Conf *Tuner
	Log  *logrus.Logger
}

// HelloWorld - handler method for example
func (h *Handler) Default(w http.ResponseWriter, req *http.Request) {
	if req != nil {
		w.Header().Del("Content-Type")
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	}
	return
}
