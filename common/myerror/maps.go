package myerror

var lanErrMaps = map[string]map[ErrCode]string{
	"en": enMaps,
	"cn": cnMaps,
}

var enMaps = map[ErrCode]string{
	UsernameErrCode:        "Please enter the correct email or wallet address",
	GetVerifyCodeErrCode:   "Failed to obtain verification code/random code",
	LoginCodeErrCode:       "Login failed",
	ParamErrCode:           "Param error",
	AddrSignOrCodeErrCode:  "Wallet address signature or email verification code error",
	LoginOutErrCode:        "Login out error",
	BindKeplrErrCode:       "Binding Keplr error",
	GetUserInfoErrCode:     "Get info of user error",
	GetRewardDetailErrCode: "Get detail of reward error",
}
var cnMaps = map[ErrCode]string{
	UsernameErrCode:        "请输入正确的邮箱或钱包地址",
	GetVerifyCodeErrCode:   "获取验证码/随机码失败",
	LoginCodeErrCode:       "登陆失败",
	ParamErrCode:           "参数错误",
	AddrSignOrCodeErrCode:  "钱包地址签名或邮箱验证码错误",
	LoginOutErrCode:        "退出登陆失败",
	BindKeplrErrCode:       "绑定keplr钱包失败",
	GetUserInfoErrCode:     "获取用户信息失败",
	GetRewardDetailErrCode: "获取奖励详情错误",
}
