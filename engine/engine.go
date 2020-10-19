package engine

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/talkwithcode-com/codex/lib"
	"github.com/talkwithcode-com/codex/lib/lang"
)

// Engine is a codexer engine which will run the code by accepting the source code,
// programming language, and execution time limit as input.
type Engine struct {
	tempDir  string
	exec     lib.Executor
	fm       lib.FileCreateRemover
	filename string
}

// New create new engine
func New(tempDir, filename string) *Engine {
	return &Engine{
		tempDir:  tempDir,
		filename: filename,
		fm:       new(FileManger),
	}
}

// Run code and return the outputs.
func (e *Engine) Run(ctx context.Context, code *lib.Code, input []byte) (lib.Output, error) {

	output := lib.Output{}
	return output, nil
}

// WriteCode into file system.
func (e *Engine) WriteCode(language lang.Language, sourceCode string) (*lib.Code, error) {
	config := lang.LanguageConfig[language]

	filePath := fmt.Sprintf("%s/%s.%s", e.tempDir, e.filename, config.Extension)

	file, err := e.fm.Create(filePath)

	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader([]byte(sourceCode))
	io.Copy(file, reader)

	code := &lib.Code{
		Language: language,
		Path:     filePath,
	}

	return code, nil

}

// DeleteCode from file system.
func (e *Engine) DeleteCode(code *lib.Code) {

}
