package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"miniblog/app/interaction/api/internal/svc"
	"miniblog/app/interaction/api/internal/types"
	"miniblog/app/interaction/rpc/interaction"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentDeleteLogic {
	return &CommentDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentDeleteLogic) CommentDelete(req *types.CommentDeleteReq) (resp *types.CommentDeleteResp, err error) {
	// 1. 安全提取 JWT 中的用户 ID
	var userId int64
	if uidVal := l.ctx.Value("userId"); uidVal != nil {
		switch v := uidVal.(type) {
		case json.Number:
			userId, _ = v.Int64()
		case float64:
			userId = int64(v)
		case string:
			fmt.Sscanf(v, "%d", &userId)
		}
	}

	if userId == 0 {
		return nil, fmt.Errorf("未授权的访问")
	}

	// 2. 呼叫底层 RPC 进行删除，强制传递 PostId 以便寻址
	_, err = l.svcCtx.InteractionRpc.CommentDelete(l.ctx, &interaction.CommentDeleteRequest{
		UserId:    userId,
		CommentId: req.CommentId,
		PostId:    req.PostId,
	})
	if err != nil {
		return nil, err
	}

	return &types.CommentDeleteResp{}, nil
}
