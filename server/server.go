package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/talkwithcode-com/codex/lib"
	"github.com/talkwithcode-com/codex/lib/lang"
)

// CodexServer ...
type CodexServer struct {
	engine lib.Codexer
}

var mapLang = map[string]lang.Language{
	"js": lang.JavaScript,
	"go": lang.Go,
}

// Submission ...
type Submission struct {
	SourceCode string   `json:"source_code"`
	Language   string   `json:"language"`
	Inputs     []string `json:"inputs"`
}

// Result ...
type Result struct {
	Outputs []lib.Output `json:"outputs"`
}

// New create CodexServer instance
func New(engine lib.Codexer) *CodexServer {
	return &CodexServer{
		engine: engine,
	}
}

func (cs *CodexServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if r.Method == http.MethodPost {
		switch path {
		case "/run":
			cs.handleRun(rw, r)
		default:
			rw.WriteHeader(http.StatusNotFound)
		}
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (cs *CodexServer) handleRun(rw http.ResponseWriter, r *http.Request) {
	// Read request body
	// Create code from request body
	// for each inputs Run code and store to ouputs
	// write output as

	var s Submission
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(rw, "Bad Request", http.StatusBadRequest)
		return
	}

	code, err := cs.engine.WriteCode(mapLang[s.Language], s.SourceCode)
	defer cs.engine.DeleteCode(code)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	result := new(Result)

	for _, input := range s.Inputs {
		out, err := cs.engine.Run(context.Background(), code, []byte(input))
		if err == nil {
			result.Outputs = append(result.Outputs, *out)
		}

	}

	b, err := json.Marshal(result)

	if err != nil {
		log.Fatal(err)
	}
	rw.Write(b)
}
