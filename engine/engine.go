package engine

import (
	"context"

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

// Run code and return the outputs.
func (e *Engine) Run(ctx context.Context, code *lib.Code, input []byte) (lib.Output, error) {

	output := lib.Output{}
	return output, nil
}

// WriteCode into file system.
func (e *Engine) WriteCode(language lang.Language, sourceCode string) (*lib.Code, error) {
	return nil, nil
}

// DeleteCode from file system.
func (e *Engine) DeleteCode(code *lib.Code) {

}
