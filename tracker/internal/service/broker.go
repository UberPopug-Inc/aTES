package service

import (
	"context"
)

type TaskEventer interface {
	Done(ctx context.Context, taskID string) error
	Created(ctx context.Context, taskID string) error
	Assigned(ctx context.Context, task Task) error
}
