// Code generated by goctl. DO NOT EDIT.
package types

type CancelUnDelegateReq struct {
	Validator string  `json:"validator"` // 验证者
	Amount    float64 `json:"amount"`    // 数量
	Height    int64   `json:"height"`    // 高度
}

type CodeResp struct {
	Code string `json:"code"`
}

type DelegateReq struct {
	Validator string  `json:"validator"` // 验证者
	Amount    float64 `json:"amount"`    // 数量
}

type GetCodeReq struct {
	Wallet string `form:"wallet"`
}

type GetValidatorReq struct {
	Kind   int8   `form:"kind,options=0|1,optional"` // 获取验证者节点信息 0-所有 1-质押
	Key    string `form:"key,optional"`              // key
	Page   uint64 `form:"page,optional,default=1"`
	Size   uint64 `form:"size,optional,default=10"`
	SortBy int8   `form:"sortby,options=0|1,optional"` // 排序方式 0-质押总量 1-质押手续费
	Sort   int8   `form:"sort,options=0|1,optional"`   // 排序方式 0-倒序 1-正序
}

type LoginReq struct {
	Wallet    string `json:"wallet"`    // 钱包地址
	Sign      string `json:"sign"`      // 签名
	PublicKey string `json:"publicKey"` // 公钥
}

type LoginResp struct {
	Token string `json:"token"`
}

type ReDelegateReq struct {
	SrcValidator string  `json:"scrValidator"` // 原验证者
	DstValidator string  `json:"dstValidator"` // 目标验证者
	Amount       float64 `json:"amount"`       // 数量
}

type UnbindingDelegateInfo struct {
	ID              int64   `json:"id"`
	Image           string  `json:"image"`           // 验证者头像
	Name            string  `json:"name"`            // 验证者名称
	Validator       string  `json:"validator"`       // 验证者
	Tokens          float64 `json:"tokens"`          // 数量
	UnbindingPeriod float64 `json:"unbindingPeriod"` // 解绑期，最低的解绑到期时间戳
	Unit            string  `json:"unit"`            // 单位
	Height          int64   `json:"height"`          // 高度
	Status          bool    `json:"status"`          // true可用 false不可用
}

type UserInfo struct {
	TotalToken     float64  `json:"totalToken"`     // 总数
	AvailableToken float64  `json:"availableToken"` // 可用余额
	StakedToken    float64  `json:"stakedToken"`    // 质押数量
	Reward         float64  `json:"reward"`         // 质押收益
	UnstakedToken  float64  `json:"unstakedToken"`  // 锁仓质押
	ValidatorAddr  []string `json:"validatorAddr"`  // 质押验证者地址
}

type ValidatorInfo struct {
	ID              int64   `json:"id"`
	Image           string  `json:"image"`           // 验证者头像
	Name            string  `json:"name"`            // 验证者名称
	Validator       string  `json:"validator"`       // 验证者
	StakedTokens    float64 `json:"stakedTokens"`    // 总质押量
	Rate            float64 `json:"rate"`            // 预期年利率
	VotingPower     float64 `json:"votingPower"`     // 投票权
	UnbindingPeriod string  `json:"unbindingPeriod"` // 解绑期，最低的解绑到期时间戳
	HandlingFees    float64 `json:"handlingFees"`    // 质押手续费
	Status          bool    `json:"status"`          // true可用 false不可用
}

type Validators struct {
	Total int64           `json:"total"`
	List  []ValidatorInfo `json:"list"`
}
