package service

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
)

type UserEvents interface {
	Logged(ctx context.Context, user gocloak.UserInfo) error
	Created(ctx context.Context, userID string) error
	Updated(ctx context.Context, user gocloak.User) error
	Deleted(ctx context.Context, userID string) error
}
