package server

import (
	"net/http"

	"github.com/talkwithcode-com/codex/lib"
)

// CodexServer ...
type CodexServer struct {
	engine lib.Codexer
}

// New create CodexServer instance
func New(engine lib.Codexer) *CodexServer {
	return &CodexServer{
		engine: engine,
	}
}

func (cs *CodexServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

}
