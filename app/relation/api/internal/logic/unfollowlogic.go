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

type UnfollowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnfollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnfollowLogic {
	return &UnfollowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnfollowLogic) Unfollow(req *types.UnfollowReq) (resp *types.UnfollowResp, err error) {
	// 1. 解析 JWT Token
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

	// 2. 呼叫底层 RPC 取消关注
	_, err = l.svcCtx.RelationRpc.Unfollow(l.ctx, &relationclient.UnfollowReq{
		FollowerId:  currentUserId,
		FollowingId: req.TargetUserId,
	})

	if err != nil {
		return nil, err
	}

	return &types.UnfollowResp{}, nil
}
