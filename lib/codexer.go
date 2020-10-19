package lib

import (
	"context"

	"github.com/talkwithcode-com/codex/lib/lang"
)

// Code ...
type Code struct {
	Path     string
	Language lang.Language
}

// Output ...
type Output struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

// Codexer ...
type Codexer interface {
	Run(ctx context.Context, code *Code, input []byte) (Output, error)
	WriteCode(language lang.Language, sourceCode string) (*Code, error)
	DeleteCode(code *Code)
}
