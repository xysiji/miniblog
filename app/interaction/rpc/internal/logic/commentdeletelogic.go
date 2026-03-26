package logic

import (
	"context"
	"fmt"

	"miniblog/app/interaction/rpc/interaction"
	"miniblog/app/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommentDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentDeleteLogic {
	return &CommentDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CommentDeleteLogic) CommentDelete(in *interaction.CommentDeleteRequest) (*interaction.CommentDeleteResponse, error) {
	// 1. 根据 ID 查找评论是否存在
	comment, err := l.svcCtx.CommentModel.FindOne(l.ctx, in.CommentId)
	if err != nil {
		l.Logger.Errorf("查询评论失败: %v", err)
		return nil, fmt.Errorf("评论不存在或已被删除")
	}

	// 2. 核心安全校验：只能删除自己发表的评论
	if comment.UserId != in.UserId {
		return nil, fmt.Errorf("越权操作：无权删除他人的评论")
	}

	// 3. 执行软删除 (将状态置为 0)
	comment.Status = 0
	err = l.svcCtx.CommentModel.Update(l.ctx, comment)
	if err != nil {
		l.Logger.Errorf("软删除评论失败: %v", err)
		return nil, fmt.Errorf("删除失败，请稍后重试")
	}

	l.Logger.Infof("用户 %d 成功删除了评论 %d", in.UserId, in.CommentId)
	// 进阶扩展点：这里可以发送消息队列，通知 Post 服务将博文的 comment_count 减 1

	return &interaction.CommentDeleteResponse{}, nil
}
