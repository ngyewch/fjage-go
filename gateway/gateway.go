package gateway

import (
	"context"
	"io"
)

type Gateway interface {
	io.Closer

	Agents(ctx context.Context) error
}
