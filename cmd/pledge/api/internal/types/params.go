package types

type (
	// HeaderKey 从请求头中获取的key
	HeaderKey string
)

const (
	// LangKey 上下文中语言的key
	LangKey HeaderKey = "language"

	// CodeRedisPre redis中随机码前缀
	CodeRedisPre string = "titan_pledge_code"
	// TokenPre token前缀
	TokenPre string = "titan_pledge_token"

	// RateKey 验证者节点预期年收益
	RateKey = "titan_pledge_rate"
	// DelegateCoinKey 质押的金额
	DelegateCoinKey = "titan_pledge_delegation"
	// LastDelegateReward 上一次统计的质押收益
	LastDelegateReward = "titan_pledge_last_reward"
)
