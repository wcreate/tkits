package tkits

var (
	INVALID_URL = Error{
		"invalid url!",
		"url is invalid or method is not correct.",
	}

	INVALID_AUTH = Error{
		"invalid request!",
		"not found Authorization in header or the value is invalid.",
	}

	INVALID_BODY = Error{
		"invalid request!",
		"request body is not correct for this url.",
	}

	DB_ERROR = Error{
		"system error!",
		"operate db failed.",
	}

	SYS_ERROR = Error{
		"system error!",
		"unkown error.",
	}
)

// Common Error Response
type Error struct {
	ErrorMsg string `json:"error"`
	Detail   string `json:"detail"`
}
