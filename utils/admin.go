package utils

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrGetUserID = errors.New("cant get userID from ctx")

func CheckPermission(ctx context.Context, permission []string) bool {
	roles, ok := ctx.Value("acl.permission").(map[string]bool)
	if !ok {
		return false
	}
	for _, p := range permission {
		if _, ok2 := roles[p]; ok2 {
			return true
		}
	}
	return false
}

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value("userID").(uuid.UUID)
	if !ok {
		return uuid.UUID{}, ErrGetUserID
	}
	return userID, nil
}
