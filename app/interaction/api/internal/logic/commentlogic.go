package logic

import (
	"context"
	"encoding/json"

	"miniblog/app/interaction/api/internal/svc"
	"miniblog/app/interaction/api/internal/types"
	"miniblog/app/interaction/rpc/interaction"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentLogic {
	return &CommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentLogic) Comment(req *types.CommentReq) (resp *types.CommentResp, err error) {
	// 1. 获取当前登录用户的 ID
	userIdNumber, ok := l.ctx.Value("jwtUserId").(json.Number)
	var userId int64
	if ok {
		userId, _ = userIdNumber.Int64()
	} else {
		uid, _ := l.ctx.Value("userId").(json.Number)
		userId, _ = uid.Int64()
	}

	// 2. 调用后端 RPC 的评论方法
	rpcResp, err := l.svcCtx.InteractionRpc.Comment(l.ctx, &interaction.CommentRequest{
		UserId:  userId,
		PostId:  req.PostId,
		Content: req.Content,
	})

	if err != nil {
		return nil, err
	}

	// 3. 返回生成的评论 ID
	return &types.CommentResp{
		CommentId: rpcResp.CommentId,
	}, nil
}
