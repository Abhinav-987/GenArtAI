package view

import (
	"context"
	"strconv"

	"github.com/Abhinav-987/GenArtAI/models"
)

func AuthenticatedUser(ctx context.Context) models.AuthenticatedUser {
	user, ok := ctx.Value(models.UserContextKey).(models.AuthenticatedUser)
	if !ok {
		return models.AuthenticatedUser{}
	}
	return user
}

func String(i int) string {
	return strconv.Itoa(i)
}
