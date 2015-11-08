package tkits

var (
	INVALID_URL = Error{
		"invalid_url",
		"url is invalid or method is not correct.",
	}

	INVALID_AUTH = Error{
		"invalid_token",
		"not found token in header or the value is invalid.",
	}

	INVALID_BODY = Error{
		"invalid_request",
		"request body is not correct for this url.",
	}

	DB_ERROR = Error{
		"system_error",
		"operate db failed.",
	}

	SYS_ERROR = Error{
		"system_error",
		"unkown error.",
	}
)

// Common Error Response
type Error struct {
	ErrorMsg string `json:"error"`
	Detail   string `json:"detail"`
}
