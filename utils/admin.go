package utils

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrGetUserID = errors.New("cant get userID from ctx")

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value("user.id").(uuid.UUID)
	if !ok {
		return uuid.UUID{}, ErrGetUserID
	}
	return userID, nil
}
