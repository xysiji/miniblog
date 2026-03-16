package logic

import (
	"context"
	"fmt"

	"miniblog/app/interaction/model"
	"miniblog/app/interaction/rpc/interaction"
	"miniblog/app/interaction/rpc/internal/svc"

	"github.com/bwmarrin/snowflake"
	"github.com/zeromicro/go-zero/core/logx"
)

type CommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentLogic {
	return &CommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Comment 发表评论
func (l *CommentLogic) Comment(in *interaction.CommentRequest) (*interaction.CommentResponse, error) {
	// 1. 生成评论专属的雪花算法 ID
	node, err := snowflake.NewNode(3) // 节点 ID 设为 3，区分于 User(1) 和 Post(2)
	if err != nil {
		return nil, fmt.Errorf("生成评论ID失败: %v", err)
	}
	commentId := node.Generate().Int64()

	// 2. 组装数据，直接同步写入 MySQL
	newComment := &model.Comment{
		Id:      commentId,
		PostId:  in.PostId,
		UserId:  in.UserId,
		Content: in.Content,
	}

	// 使用 goctl 自动生成的 Insert 方法写入
	_, err = l.svcCtx.CommentModel.Insert(l.ctx, newComment)
	if err != nil {
		l.Logger.Errorf("评论写入 MySQL 失败: %v", err)
		return nil, fmt.Errorf("评论发布失败，请稍后重试")
	}

	l.Logger.Infof("用户 %d 对博文 %d 发表了评论", in.UserId, in.PostId)

	return &interaction.CommentResponse{
		CommentId: commentId,
	}, nil
}
