package utils

import (
	"context"
	"github.com/AleksK1NG/api-mc/pkg/httpErrors"
	"github.com/AleksK1NG/api-mc/pkg/logger"
)

// Validate is user from owner of content
func ValidateIsOwner(ctx context.Context, creatorId string) error {
	user, err := GetUserFromCtx(ctx)

	if err != nil {
		return err
	}

	if user.UserID.String() != creatorId {
		logger.Errorf(
			"ValidateIsOwner, userID: %v, creatorID: %v",
			user.UserID.String(),
			creatorId,
		)
		return httpErrors.Forbidden
	}

	return nil
}
