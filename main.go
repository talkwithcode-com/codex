package main

import (
	"log"
	"net/http"
	"os"

	"github.com/talkwithcode-com/codex/engine"
	"github.com/talkwithcode-com/codex/server"
)

var port = os.Getenv("CODEX_PORT")

func main() {
	eg := engine.New(os.TempDir(), "main")
	cs := server.New(eg)
	log.Fatal(http.ListenAndServe(":"+port, cs))
}
