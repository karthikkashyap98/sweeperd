package actions

import "context"

type Action interface {
	Plan(ctx context.Context) ([]string, error)
	Execute(ctx context.Context, files []string) error
}
