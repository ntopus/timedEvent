package server

import "context"

type IServer interface {
	RunServer(ctx context.Context) error
}
