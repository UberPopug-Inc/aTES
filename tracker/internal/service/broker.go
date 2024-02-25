package service

import (
	"context"
)

type TaskEventer interface {
	Done(ctx context.Context, task Task) error
	Created(ctx context.Context, task Task) error
	Assigned(ctx context.Context, task Task) error
}
