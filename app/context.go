package app

import (
	"context"

	"github.com/workdestiny/watgok_web/entity"
)

type contextKey int

const (
	ckMyID contextKey = iota
	ckUser
	ckRole
)

// WithMyID id put to Context
func WithMyID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ckMyID, id)
}

// GetMyID เอา ID user ออกจาก Context เพื่อเอาไปใช้
func GetMyID(ctx context.Context) string {
	return ctx.Value(ckMyID).(string)
}

// WithUserRole role put to Context
func WithUserRole(ctx context.Context, role entity.Role) context.Context {
	return context.WithValue(ctx, ckRole, role)
}

// GetUserRole เอา Role user ออกจาก Context เพื่อเอาไปใช้
func GetUserRole(ctx context.Context) entity.Role {
	return ctx.Value(ckRole).(entity.Role)
}

// WithUser put user to Context
func WithUser(ctx context.Context, u *entity.UserModel) context.Context {
	return context.WithValue(ctx, ckUser, u)
}

func getUser(ctx context.Context) *entity.UserModel {
	return ctx.Value(ckUser).(*entity.UserModel)
}
