package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

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
	TestCases  []struct {
		ID     string `json:"id"`
		Input  string `json:"input"`
		Output string `json:"output"`
	} `json:"test_cases"`
}

// ResponseBody ...
type ResponseBody struct {
	Language string `json:"language"`
	Logs     []Log  `json:"logs"`
}

// Log ...
type Log struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Stderr string `json:"stderr"`
	Stdout string `json:"stdout"`
	Stdint string `json:"stdin"`
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
		http.Error(rw, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	code, err := cs.engine.WriteCode(mapLang[s.Language], s.SourceCode)
	defer cs.engine.DeleteCode(code)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var rb ResponseBody

	rb.Language = s.Language

	for _, test := range s.TestCases {
		out, err := cs.engine.Run(context.Background(), code, []byte(test.Input))
		l := Log{
			ID:     test.ID,
			Stdint: test.Input,
		}

		if err == nil {
			l.Stderr = strings.TrimSuffix(out.Stderr, "\n")
			l.Stdout = strings.TrimSuffix(out.Stdout, "\n")
		}

		l.Status = "fail"
		if l.Stdout == test.Output {
			l.Status = "success"
		}

		rb.Logs = append(rb.Logs, l)
	}

	b, err := json.Marshal(rb)

	if err != nil {
		log.Fatal(err)
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(b)
}
