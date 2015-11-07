package tkits

import (
	"fmt"
	"strings"

	"gopkg.in/macaron.v1"
)

func CheckAuth(ctx *macaron.Context, uid int64) bool {
	auth := strings.TrimSpace(ctx.Header().Get("Authorization"))

	if !GetSimpleToken().Validate(auth, ctx.RemoteAddr(), fmt.Sprintf("%v", uid)) {
		ctx.JSON(404, INVALID_AUTH)
		return false
	}
	return true
}
