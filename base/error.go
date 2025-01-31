package base

var (
	InvalidParamErr   = NewException(400, "InvalidParam", "An invalid request payload is supplied.", "请求参数错误。")
	ResourceNotFound  = NewException(404, "ResourceNotFound", "The specified resource is not found.", "资源未找到。")
	InternalError     = NewException(500, "InternalError", "There is an internal error occurred.", "内部错误，请联系管理员。")
	ResourceInUse     = NewException(409, "ResourceInUse", "The specified resource already exists.", "资源已存在。")
	RateLimitExceeded = NewException(429, "RateLimitExceeded", "Request is due to rate limit.", "请求过于频繁，请稍后再试。")
	PermissionDenied  = NewException(403, "PermissionDenied", "You have no permission to do this operation.", "您没有权限进行此操作。")
	Unauthorized      = NewException(401, "Unauthorized", "Unauthorized identity.", "未授权的身份。")
)
