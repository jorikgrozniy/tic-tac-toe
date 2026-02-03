package contextkey

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
)

func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	value := ctx.Value(UserIDKey)
	if value == nil {
		return uuid.Nil, false
	}

	id, ok := value.(uuid.UUID)
	return id, ok
}
