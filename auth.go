package tkits

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/macaron.v1"
)

// Support two authrization ways
// 1. authrizate by cookie
// 2  authrizate by http header 'Authorization'
func CheckAuth(ctx *macaron.Context, uid int64) bool {

	// firstly retrive the cookie
	suid := ctx.GetCookie("uid")
	token := ctx.GetCookie("token")
	vtime := true
	log.Debugf("retrive uid = %s, token = % from cookie", suid, token)

	ssuid := fmt.Sprintf("%v", uid)
	// only check the time expire if token not exits in cookie
	if suid != "" && token != "" {
		vtime = false

		if ssuid != suid {
			ctx.JSON(http.StatusUnauthorized, INVALID_AUTH)
			return false
		}
	} else {
		token = strings.TrimSpace(ctx.Header().Get("Authorization"))
		suid = ssuid
		log.Debugf("retrive token = % from header", token)
	}

	if !GetSimpleToken().Validate(token, ctx.RemoteAddr(), suid, vtime) {
		ctx.JSON(http.StatusUnauthorized, INVALID_AUTH)
		return false
	}
	return true
}
