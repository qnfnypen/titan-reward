package myerror

// ErrCode 类型再定义
type ErrCode int64

const (
	// UsernameErrCode 邮箱或钱包地址错误
	UsernameErrCode ErrCode = 1000 + iota
	// GetVerifyCodeErrCode 获取验证码/随机码失败
	GetVerifyCodeErrCode
	// CaptchaErrCode 人机校验失败
	CaptchaErrCode
	// LoginCodeErrCode 用户登陆失败
	LoginCodeErrCode
	// EmailOrPasswordErrCode 邮箱或密码错误
	EmailOrPasswordErrCode
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
	// RelatedErrCode 绑定邮箱或钱包错误
	RelatedErrCode
	// MulRelatedErrCode 重复绑定绑定邮箱或钱包错误
	MulRelatedErrCode
	// BoundErrCode 邮箱或钱包已被绑定错误
	BoundErrCode
	// KeplrBoundErrCode keplr账户绑定后不能进行关联账户的绑定
	KeplrBoundErrCode

	// GetValitorErrCode 获取验证节点信息错误
	GetValitorErrCode
	// DelegateTokenErrCode 质押token错误
	DelegateTokenErrCode
	// UnDelegateTokenErrCode 解除质押错误
	UnDelegateTokenErrCode
	// ReDelegateTokenErrCode 质押转移错误
	ReDelegateTokenErrCode
	// CancelUnDelegateTokenErrCode 取消解除质押错误
	CancelUnDelegateTokenErrCode
	// WithdrawRewardsErrCode 提取收益错误
	WithdrawRewardsErrCode
)
