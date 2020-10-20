package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
	"github.com/talkwithcode-com/codex/engine"
	"github.com/talkwithcode-com/codex/server"
)

var port = os.Getenv("CODEX_PORT")

func main() {
	eg := engine.New(os.TempDir(), "main")
	cs := server.New(eg)
	handler := cors.Default().Handler(cs)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
