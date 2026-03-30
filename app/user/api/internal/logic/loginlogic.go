package logic

import (
	"context"
	"time"

	"miniblog/app/user/api/internal/svc"
	"miniblog/app/user/api/internal/types"
	"miniblog/app/user/rpc/userclient"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 辅助函数：生成 JWT Token
func (l *LoginLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds // 过期时间
	claims["iat"] = iat           // 签发时间
	claims["userId"] = userId     // 载荷：把 userId 藏在 token 里

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 1. 调用 RPC 进行账号密码校验
	rpcRes, err := l.svcCtx.UserRpc.Login(l.ctx, &userclient.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	// 2. 密码校验通过，在 API 网关层独立生成 JWT Token
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	jwtToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, accessExpire, rpcRes.UserId)
	if err != nil {
		return nil, err
	}

	// 3. 【修复点】：严格按照 types.LoginResp 结构体返回给前端
	return &types.LoginResp{
		UserId: rpcRes.UserId,
		Token:  jwtToken,
	}, nil
}
