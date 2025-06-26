package known

const (
	// XRequestIDKey 用来定义 Gin 上下文中的键，代表请求的 uuid.
	XRequestIDKey = "X-Request-ID"

	// XUserIDKey 用来定义 Gin 上下文的键，代表请求的所有者.
	XUserIDKey = "X-User-ID"

	// IpWhoKey 用来定义 Gin 上下文的键，代表请求的所有者.
	XIpWhoKey = "X-IpWho"
)
