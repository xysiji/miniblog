package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"miniblog/app/relation/api/internal/svc"
	"miniblog/app/relation/api/internal/types"
	"miniblog/app/relation/rpc/relationclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FollowLogic) Follow(req *types.FollowReq) (resp *types.FollowResp, err error) {
	// 1. 从 JWT 中安全提取当前登录的 userId (即粉丝 ID)
	var currentUserId int64
	if uidVal := l.ctx.Value("userId"); uidVal != nil {
		switch v := uidVal.(type) {
		case json.Number:
			currentUserId, _ = v.Int64()
		case float64:
			currentUserId = int64(v)
		case string:
			fmt.Sscanf(v, "%d", &currentUserId)
		}
	}

	if currentUserId == 0 {
		return nil, fmt.Errorf("未授权的访问，请重新登录")
	}

	// 2. 呼叫底层 RPC 执行关注
	_, err = l.svcCtx.RelationRpc.Follow(l.ctx, &relationclient.FollowReq{
		FollowerId:  currentUserId,
		FollowingId: req.TargetUserId, // 前端传来的目标用户 ID
	})

	if err != nil {
		return nil, err
	}

	return &types.FollowResp{}, nil
}
