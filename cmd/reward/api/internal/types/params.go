package types

type (
	// HeaderKey 从请求头中获取的key
	HeaderKey string
)

const (
	// LangKey 上下文中语言的key
	LangKey HeaderKey = "language"

	// CodeRedisPre redis中随机码前缀
	CodeRedisPre string = "titan_reward_code"
	// TokenPre token前缀
	TokenPre string = "titan_reward_token"
)
