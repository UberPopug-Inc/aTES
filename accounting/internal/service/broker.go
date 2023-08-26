package service

import (
	"context"

	"github.com/UberPopug-Inc/aTES/accounting/internal/events"
)

type TaskEventer interface {
	AssignedPull(ctx context.Context) (*events.TaskV1, error)
	DonePull(ctx context.Context) (*events.TaskV1, error)
}
