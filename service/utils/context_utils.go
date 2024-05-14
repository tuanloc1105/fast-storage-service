package utils

import (
	"context"
	"fast-storage-go-service/constant"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PrepareContext(c *gin.Context, isBypassCurrentUserCheck ...bool) (context.Context, bool) {
	ctx := context.Background()

	currentUser, isCurrentUserExist := GetCurrentUsername(c)
	currentUserId, isCurrentUserIdExist := GetCurrentUserId(c)

	ctx = context.WithValue(ctx, constant.TraceIdLogKey, GetTraceId(c))

	if len(isBypassCurrentUserCheck) > 0 && isBypassCurrentUserCheck[0] {
		return ctx, true
	} else {
		if isCurrentUserExist != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				ReturnResponse(
					c,
					constant.Unauthorized,
					nil,
					isCurrentUserExist.Error(),
				),
			)
			return ctx, false
		}
		if isCurrentUserIdExist != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				ReturnResponse(
					c,
					constant.Unauthorized,
					nil,
					isCurrentUserIdExist.Error(),
				),
			)
			return ctx, false
		}
		ctx = context.WithValue(ctx, constant.UsernameLogKey, *currentUser)
		ctx = context.WithValue(ctx, constant.UserIdLogKey, *currentUserId)
		return ctx, true
	}
}

func GetCurrentUsernameFromContextForInsertOrUpdateDataInDb(ctx context.Context) string {
	var currentUsernameInsertOrUpdateData = ""
	var usernameFromContext = ctx.Value(constant.UsernameLogKey)
	if usernameFromContext != nil {
		currentUsernameInsertOrUpdateData = usernameFromContext.(string)
	}
	return currentUsernameInsertOrUpdateData
}
