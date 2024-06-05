package auth

import (
	"context"
	"crypto/rand"
	"embed"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/qnfnypen/gzocomm/merror"
	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/svc"
	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/types"
	"github.com/qnfnypen/titan-reward/common/myerror"

	"github.com/zeromicro/go-zero/core/logx"
)

//go:embed emailtmps
var emailTmpl embed.FS

type unKind int

const (
	errKind unKind = iota - 1
	emailKind
	addrKind
)

// GetVerifyCodeLogic 获取邮箱验证码/钱包地址的随机码
type GetVerifyCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetVerifyCodeLogic 新建 获取邮箱验证码/钱包地址的随机码
func NewGetVerifyCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVerifyCodeLogic {
	return &GetVerifyCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetVerifyCode 实现 获取邮箱验证码/钱包地址的随机码
func (l *GetVerifyCodeLogic) GetVerifyCode(req *types.GetVerifyCodeReq) (resp *types.VerifyCodeResp, err error) {
	var (
		gzErr = merror.GzErr{}
	)
	resp = new(types.VerifyCodeResp)

	// 获取语言
	lan := l.ctx.Value(types.LangKey).(string)
	gzErr.RespErr = myerror.GetMsg(myerror.GetVerifyCodeErrCode, lan)

	// 生成随机码，并设置redis缓存
	nonce := generateNonce(6)
	err = l.svcCtx.RedisCli.SetexCtx(l.ctx, fmt.Sprintf("%s_%s", types.CodeRedisPre, req.Username), nonce, 5*60)
	if err != nil {
		return nil, gzErr
	}
	// 判断用户名是邮箱还是钱包地址
	switch checkUsername(req.Username) {
	case emailKind:
		go l.sendCodeEmail(nonce, lan, req.Username)
		// resp.VerifyCode = nonce
	case addrKind:
		resp.VerifyCode = fmt.Sprintf("TitanNetWork(%s)", nonce)
	default:
		gzErr.RespErr = myerror.GetMsg(myerror.UsernameErrCode, lan)
		return nil, gzErr
	}

	return resp, nil
}

func (l *GetVerifyCodeLogic) sendCodeEmail(code, lang string, email ...string) error {
	var (
		sendBody = new(strings.Builder)

		tmpl    *template.Template
		err     error
		subject string
		codes   []string
	)

	for _, v := range code {
		codes = append(codes, string(v))
	}
	switch lang {
	case "cn":
		subject = "[Titan Network] 您的验证码"
		tmpl, err = template.ParseFS(emailTmpl, "emailtmps/mail_cn.html")
	default:
		subject = "[Titan Network] Your verification code"
		tmpl, err = template.ParseFS(emailTmpl, "emailtmps/mail_en.html")
	}
	if err != nil {
		return fmt.Errorf("send verify_code of email error: %w", err)
	}

	err = tmpl.Execute(sendBody, map[string]interface{}{
		"Code": codes,
	})
	if err != nil {
		return fmt.Errorf("template execute verify_code of email error: %w", err)
	}

	from := fmt.Sprintf("Titan Network <%s>", l.svcCtx.Config.Email.Username)
	return l.svcCtx.EmailCli.SendEmail(subject, from, "", sendBody, "text/html", email...)
}

func generateNonce(len int) string {
	randBytes := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

	nonce := make([]byte, len)
	io.ReadFull(rand.Reader, nonce)

	for i, v := range nonce {
		nonce[i] = randBytes[v%10]
	}

	return string(nonce)
}
