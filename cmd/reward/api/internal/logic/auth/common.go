package auth

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/qnfnypen/gzocomm/merror"
	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/svc"
	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/types"
)

// Sctx 类型再定义
type sctx svc.ServiceContext

func (s *sctx) generateToken(ctx context.Context, uid int64, uuid string) (string, error) {
	key := fmt.Sprintf("%s_%s", types.TokenPre, uuid)

	td, err := time.ParseDuration(s.Config.Auth.AccessExpire)
	if err != nil {
		td = 24 * time.Hour
	}

	// 先从 redis 中获取，不存在则重新创建
	token, err := s.getToken(ctx, key, td)
	if err != nil {
		err = merror.NewError(fmt.Errorf("get token from redis error:%w", err))
		return "", err
	}
	if token != "" {
		return token, nil
	}
	// 不存在则生成 token，这里和go-zero对应使用jwt.MapClaims
	claims := jwt.MapClaims{
		"uid":  uid,
		"uuid": uuid,
		"iss":  "titan",
		"exp":  time.Now().Add(td).Unix(),
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenClaims.SignedString([]byte(s.Config.Auth.AccessSecret))
	if err != nil {
		err = merror.NewError(fmt.Errorf("signed jwt token error:%w", err))
		return "", err
	}
	// 将token存储到redis中去
	_, err = s.RedisCli.SetnxExCtx(ctx, key, token, int(td.Seconds()))
	if err != nil {
		err = merror.NewError(fmt.Errorf("redis set token error:%w", err))
		return "", err
	}

	return token, nil
}

// getToken 获取用户token
func (s *sctx) getToken(ctx context.Context, key string, td time.Duration) (string, error) {
	// 根据key获取value的ttl时间
	ttl, err := s.RedisCli.Ttl(key)
	if err != nil {
		return "", err
	}
	// 如果时间低于1/10，则清除redis
	if float64(ttl) < td.Seconds()*0.1 {
		s.RedisCli.DelCtx(ctx, key)
		return "", nil
	}
	// 根据key获取value
	token, err := s.RedisCli.GetCtx(ctx, key)
	if err != nil {
		return "", err
	}

	return token, nil
}

func checkUsername(un string) unKind {
	const (
		emailRegexPattern       = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		bitcoinRegexPattern     = `^(1|3)[a-km-zA-HJ-NP-Z1-9]{25,34}$`
		ethereumRegexPattern    = `^0x[a-fA-F0-9]{40}$`
		bitcoinCashRegexPattern = `^(bitcoincash:)?(q|p)[a-z0-9]{41}$`
	)

	emailRegex := regexp.MustCompile(emailRegexPattern)
	bitcoinRegex := regexp.MustCompile(bitcoinRegexPattern)
	ethereumRegex := regexp.MustCompile(ethereumRegexPattern)
	bitcoinCashRegex := regexp.MustCompile(bitcoinCashRegexPattern)

	switch {
	case emailRegex.MatchString(un):
		return emailKind
	case ethereumRegex.MatchString(un) || bitcoinRegex.MatchString(un) || bitcoinCashRegex.MatchString(un):
		return addrKind
	default:
		return errKind
	}
}
