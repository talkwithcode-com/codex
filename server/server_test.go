package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/talkwithcode-com/codex/engine"
	"github.com/talkwithcode-com/codex/lib"
	"github.com/talkwithcode-com/codex/lib/lang"
)

const baseURL = "localhost:3000/run"

type mockEngine struct {
}

func (me *mockEngine) Run(ctx context.Context, code *lib.Code, input []byte) (*lib.Output, error) {
	return &lib.Output{}, nil
}
func (me *mockEngine) WriteCode(language lang.Language, sourceCode string) (*lib.Code, error) {
	return &lib.Code{
		Path: "",
	}, errors.New("myerror")
}
func (me *mockEngine) DeleteCode(code *lib.Code) {

}

func TestCodexServer_handleRun(t *testing.T) {

	data := map[string]interface{}{
		"source_code": `
			console.log("hello")
			console.error("myerror")
		`,
		"language": "js",
		"inputs":   []string{""},
	}

	b, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Can'marshal submision")
	}

	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewReader(b))

	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rec := httptest.NewRecorder()

	e := engine.New(os.TempDir(), "main")
	cs := New(e)

	cs.handleRun(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}
	var r Result
	json.Unmarshal(b, &r)

	if r.Outputs[0].Stdout != "hello\n" {
		t.Errorf("Expect %s; got %s", r.Outputs[0].Stdout, "hello\n")
	}

	if r.Outputs[0].Stderr != "myerror\n" {
		t.Errorf("Expect %s; got %s", r.Outputs[0].Stderr, "myerror")
	}

}

func TestCodexServer_handleRun_BadRequest(t *testing.T) {

	data := map[string]interface{}{
		"source_code": `
			console.log("hello")
			console.error("myerror")
		`,
		"language": "js",
		"inputs":   "",
	}

	b, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Can'marshal submision")
	}

	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewReader(b))

	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rec := httptest.NewRecorder()

	e := engine.New(os.TempDir(), "main")
	cs := New(e)

	cs.handleRun(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Error("Expected bad request", res.StatusCode)
	}

}

func TestCodexServer_handleRun_InternalServerError(t *testing.T) {

	data := map[string]interface{}{
		"source_code": `
			console.log("hello")
			console.error("myerror")
		`,
		"language": "js",
		"inputs":   []string{""},
	}

	b, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Can'marshal submision")
	}

	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewReader(b))

	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rec := httptest.NewRecorder()

	e := new(mockEngine)
	cs := New(e)

	cs.handleRun(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusInternalServerError {
		t.Error("Expected bad request", res.StatusCode)
	}

}

func TestRouting_Ok(t *testing.T) {
	cs := New(engine.New(os.TempDir(), "main"))
	srv := httptest.NewServer(cs)
	defer srv.Close()

	data := map[string]interface{}{
		"source_code": `
			console.log("hello")
			console.error("myerror")
		`,
		"language": "js",
		"inputs":   []string{""},
	}

	b, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Can'marshal submision")
	}

	res, err := http.Post(fmt.Sprintf("%s/run", srv.URL), "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("could not send GET request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.Status)
	}

}

func TestRouting_NotFound(t *testing.T) {
	cs := New(engine.New(os.TempDir(), "main"))
	srv := httptest.NewServer(cs)
	defer srv.Close()

	data := map[string]interface{}{
		"source_code": `
			console.log("hello")
			console.error("myerror")
		`,
		"language": "js",
		"inputs":   []string{""},
	}

	b, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Can'marshal submision")
	}

	res, err := http.Post(fmt.Sprintf("%s/runx", srv.URL), "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("could not send GET request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected status OK; got %v", res.Status)
	}

}

func TestRouting_MethodNotAllowed(t *testing.T) {
	cs := New(engine.New(os.TempDir(), "main"))
	srv := httptest.NewServer(cs)
	defer srv.Close()

	res, err := http.Get(fmt.Sprintf("%s/run", srv.URL))
	if err != nil {
		t.Fatalf("could not send GET request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status OK; got %v", res.Status)
	}

}
