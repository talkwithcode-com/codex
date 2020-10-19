package engine

import (
	"context"
	"os/exec"

	"github.com/talkwithcode-com/codex/lib"
)

// Exec ..
type Exec struct {
}

// CommandContext ...
func (s Exec) CommandContext(ctx context.Context, name string, arg ...string) lib.Commander {
	return exec.CommandContext(ctx, name, arg...)
}
