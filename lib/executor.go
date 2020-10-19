package lib

import "context"

// Executor ...
type Executor interface {
	CommandContext(ctx context.Context, name string, arg ...string) Commander
}
