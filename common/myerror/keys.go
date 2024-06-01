package myerror

// ErrCode 类型再定义
type ErrCode int64

const (
	// UsernameErrCode 邮箱或钱包地址错误
	UsernameErrCode ErrCode = 1000 + iota
	// GetVerifyCodeErrCode 获取验证码/随机码失败
	GetVerifyCodeErrCode
	// LoginCodeErrCode 用户登陆失败
	LoginCodeErrCode
	// ParamErrCode 参数错误
	ParamErrCode
	// AddrSignOrCodeErrCode 钱包地址签名或邮箱验证码错误
	AddrSignOrCodeErrCode
	// LoginOutErrCode 用户登出失败
	LoginOutErrCode
	// BindKeplrErrCode 绑定keplr失败
	BindKeplrErrCode
	// GetUserInfoErrCode 获取用户信息错误
	GetUserInfoErrCode
	// GetRewardDetailErrCode 获取奖励详情错误
	GetRewardDetailErrCode
)
