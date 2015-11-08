package tkits

import (
	"net/http"
	"strconv"
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
	log.Debugf("retrive uid = %s, token = % from cookie", suid, token)

	vtime := true // only check the time expire if token not exits in cookie
	iuid := uid

	// 1. login user
	if suid != "" && token != "" {
		vtime = false // existed cookie, not check expire
		if isuid, err := strconv.ParseInt(suid, 10, 0); err != nil {
			ctx.JSON(http.StatusUnauthorized, INVALID_AUTH)
			return false
		} else {
			if iuid != isuid {
				ctx.JSON(http.StatusUnauthorized, INVALID_AUTH)
				return false
			}
		}
	} else {
		token = strings.TrimSpace(ctx.Header().Get("Authorization"))
		log.Debugf("retrive token = % from header", token)
	}

	if uid == -1 {
		vtime = false // system user, not check expire
	}

	if !GetSimpleToken().Validate(token, ctx.RemoteAddr(), iuid, vtime) {
		ctx.JSON(http.StatusUnauthorized, INVALID_AUTH)
		return false
	}
	return true
}
