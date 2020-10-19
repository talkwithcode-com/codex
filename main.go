package main

import (
	"net/http"
	"os"

	"github.com/talkwithcode-com/codex/engine"
	"github.com/talkwithcode-com/codex/server"
)

var port = os.Getenv("CODEX_PORT")

func main() {
	eg := new(engine.Engine)
	cs := server.New(eg)
	http.ListenAndServe(":"+port, cs)
}
