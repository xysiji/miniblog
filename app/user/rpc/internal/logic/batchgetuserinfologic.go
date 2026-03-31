package logic

import (
	"context"

	"miniblog/app/user/rpc/internal/svc"
	"miniblog/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchGetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetUserInfoLogic {
	return &BatchGetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BatchGetUserInfo 批量获取用户信息
func (l *BatchGetUserInfoLogic) BatchGetUserInfo(in *user.BatchUserInfoReq) (*user.BatchUserInfoResp, error) {
	resp := &user.BatchUserInfoResp{
		Users: make(map[int64]*user.UserInfoResponse),
	}

	// 复用已有的单条查询逻辑（避免重新写数据库查询，非常安全）
	userInfoLogic := NewUserInfoLogic(l.ctx, l.svcCtx)

	for _, uid := range in.UserIds {
		// 去重判断
		if _, exists := resp.Users[uid]; !exists {
			info, err := userInfoLogic.UserInfo(&user.UserInfoRequest{UserId: uid})
			if err == nil && info != nil {
				resp.Users[uid] = info
			} else {
				l.Logger.Errorf("BatchGetUserInfo 获取 uid=%d 失败: %v", uid, err)
			}
		}
	}

	return resp, nil
}
