package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"sort"
	"strings"
	"unsafe"

	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/qnfnypen/gzocomm/merror"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/svc"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/types"
	"github.com/qnfnypen/titan-reward/common/myerror"
	"github.com/shopspring/decimal"

	"github.com/zeromicro/go-zero/core/logx"
)

type validatorDesc struct {
	Data struct {
		ValDesc []struct {
			AvatarURL interface{} `json:"avatar_url"`
			ValAddr   string      `json:"validator_address"`
		} `json:"validator_description"`
	} `json:"data"`
}

// ValidatorsLogic 获取验证者信息
type ValidatorsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewValidatorsLogic 新建 获取验证者信息
func NewValidatorsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidatorsLogic {
	return &ValidatorsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Validators 实现 获取验证者信息
func (l *ValidatorsLogic) Validators(req *types.GetValidatorReq) (resp *types.Validators, err error) {
	var (
		gzErr     merror.GzErr
		comctx    = (*sctx)(unsafe.Pointer(l.svcCtx))
		tokenMaps = make(map[string]*big.Int)
	)
	resp = new(types.Validators)
	resp.List = make([]types.ValidatorInfo, 0)

	lan := l.ctx.Value(types.LangKey).(string)
	wallet := l.ctx.Value("wallet").(string)
	gzErr.RespErr = myerror.GetMsg(myerror.GetValitorErrCode, lan)

	// 获取质押解绑期
	ut, err := l.svcCtx.TitanCli.GetValidatorUnbondingTime(l.ctx)
	if err != nil {
		gzErr.LogErr = merror.NewError(err).Error()
		return nil, gzErr
	}

	// 获取所有的验证者的token
	tokens, total, validators, err := l.getAllTokensPage(req.SortBy, req.Sort, req.Key)
	if err != nil {
		gzErr.LogErr = merror.NewError(err).Error()
		return nil, gzErr
	}
	// var validators []staking.Validator
	switch req.Kind {
	case 0:
		resp.Total = total
		validators = pageValidators(validators, req.Page, req.Size)
	case 1:
		resp.Total, validators, tokenMaps, err = l.getDelgatorVidatorNums(wallet, req.SortBy, req.Sort, req.Key)
		if err != nil {
			gzErr.LogErr = merror.NewError(err).Error()
			return nil, gzErr
		}
		validators = pageValidators(validators, req.Page, req.Size)
	}
	for i, v := range validators {
		var vpf float64
		token := v.Tokens.BigInt()
		info := types.ValidatorInfo{}
		info.Name = v.Description.Moniker
		info.Validator = v.OperatorAddress
		if req.Kind == 1 {
			vpf, _ = new(big.Float).Quo(new(big.Float).SetInt(tokenMaps[v.OperatorAddress]), new(big.Float).SetInt(tokens)).Float64()
		} else {
			vpf, _ = new(big.Float).Quo(new(big.Float).SetInt(v.Tokens.BigInt()), new(big.Float).SetInt(tokens)).Float64()
		}
		info.ID = int64(int(req.Page*req.Size) + i + 1)
		info.StakedTokens = getTTNT(token)
		info.Rate = comctx.getRate(l.ctx)
		info.VotingPower, _ = decimal.NewFromFloat(vpf).Round(4).Mul(decimal.NewFromInt(100)).Float64()
		info.UnbindingPeriod = converTimeDur(ut, lan)
		dc, _ := decimal.NewFromString(v.Commission.Rate.String())
		info.HandlingFees, _ = dc.Round(4).Mul(decimal.NewFromInt(100)).Float64()
		info.Status = v.IsBonded()
		// info.Image = aus[v.OperatorAddress]
		resp.List = append(resp.List, info)
	}

	return resp, nil
}

func (l *ValidatorsLogic) getAllTokens(key string) (*big.Int, int64, error) {
	var (
		totalTokens = new(big.Int)
		count       int64
	)

	validators, err := l.svcCtx.TitanCli.QueryValidators(l.ctx, 0, 0, key)
	if err != nil {
		return nil, 0, fmt.Errorf("get all tokens of validators error:%w", err)
	}
	for _, v := range validators {
		totalTokens = totalTokens.Add(totalTokens, v.Tokens.BigInt())
		count++
	}

	return totalTokens, count, nil
}

func (l *ValidatorsLogic) getDelgatorVidatorNums(addr string, orderBy int8, order int8, key string) (int64, []staking.Validator, map[string]*big.Int, error) {
	var pTokens = make(map[string]*big.Int)
	vs, err := l.svcCtx.TitanCli.QueryDelgatorVlidators(l.ctx, addr, 0, 0)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("get total of delgator vidators error:%w", err)
	}

	for i, v := range vs {
		pTokens[v.OperatorAddress] = v.Tokens.BigInt()
		del, err := l.svcCtx.TitanCli.QueryDelegation(l.ctx, addr, v.OperatorAddress)
		if err != nil {
			continue
		}
		vs[i].Tokens = del.DelegationResponse.Balance.Amount
	}

	vs, nvs := searchValidators(vs, key)
	vs = orderValidators(vs, orderBy, order)
	vs = append(vs, nvs...)

	return int64(len(vs)), vs, pTokens, nil
}

// getAvatrURL 获取验证者节点头像
func getAvatrURL() (map[string]string, error) {
	var (
		rp   validatorDesc
		maps = make(map[string]string)
	)
	payload := []byte(`{"query":"query MyQuery {validator_description(distinct_on: validator_address) {avatar_url validator_address}}","variables":null,"operationName":"MyQuery"}`)

	resp, err := http.Post("http://8.217.10.76:8080/v1/graphql", "application/json", bytes.NewReader(payload))
	if err != nil {
		return maps, fmt.Errorf("get description of validator error:%w", err)
	}
	defer resp.Body.Close()

	payload, err = io.ReadAll(resp.Body)
	if err != nil {
		return maps, fmt.Errorf("read body of get validator's description:%w", err)
	}

	err = json.Unmarshal(payload, &rp)
	if err != nil {
		return maps, fmt.Errorf("json unmarshal of get validator's description error:%w", err)
	}

	for _, v := range rp.Data.ValDesc {
		if au, ok := v.AvatarURL.(string); ok && au != "" {
			maps[v.ValAddr] = au
		}
	}

	return maps, nil
}

// getAllTokensPage 分页查询
func (l *ValidatorsLogic) getAllTokensPage(orderBy, order int8, key string) (*big.Int, int64, []staking.Validator, error) {
	var (
		totalTokens = new(big.Int)
		count       int64
	)

	validators, err := l.svcCtx.TitanCli.QueryValidators(l.ctx, 0, 0, "")
	if err != nil {
		return nil, 0, nil, fmt.Errorf("get all tokens of validators error:%w", err)
	}
	for _, v := range validators {
		totalTokens = totalTokens.Add(totalTokens, v.Tokens.BigInt())
		count++
	}

	validators, nvs := searchValidators(validators, key)
	validators = orderValidators(validators, orderBy, order)
	validators = append(validators, nvs...)
	count = int64(len(validators))

	return totalTokens, count, validators, nil
}

func orderValidators(validators []staking.Validator, orderBy, order int8) []staking.Validator {
	sort.Slice(validators, func(i, j int) bool {
		switch orderBy {
		case 0:
			if order == 0 {
				return validators[i].Tokens.GT(validators[j].Tokens)
			}
			return validators[i].Tokens.LT(validators[j].Tokens)
		case 1:
			if order == 0 {
				return validators[i].Commission.Rate.GT(validators[j].Commission.Rate)
			}
			return validators[i].Commission.Rate.LT(validators[j].Commission.Rate)
		}

		return false
	})

	return validators
}

func searchValidators(validators []staking.Validator, key string) ([]staking.Validator, []staking.Validator) {
	var vs, nvs []staking.Validator

	key = strings.TrimSpace(key)

	for _, v := range validators {
		if key != "" {
			if strings.Contains(v.Description.Moniker, key) || strings.Contains(v.OperatorAddress, key) {
				if v.IsBonded() {
					vs = append(vs, v)
				} else {
					nvs = append(nvs, v)
				}
			}
		} else {
			if v.IsBonded() {
				vs = append(vs, v)
			} else {
				nvs = append(nvs, v)
			}
		}
	}

	return vs, nvs
}

func pageValidators(validators []staking.Validator, page, size uint64) []staking.Validator {
	var zv = make([]staking.Validator, 0)
	offset := (page - 1) * size

	if offset > uint64(len(validators)) {
		return zv
	}

	if uint64(len(validators))-offset <= size {
		return validators[offset:]
	}

	return validators[offset : offset+size]
}
