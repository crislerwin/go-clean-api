package auth

import "context"

const UserIDKey = "userID"

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return "", false
	}
	return userID, ok
}
