package logic

import (
	"context"
	"fmt"

	"miniblog/app/post/model"
	"miniblog/app/post/rpc/internal/svc"
	"miniblog/app/post/rpc/post"

	"github.com/bwmarrin/snowflake"
	"github.com/zeromicro/go-zero/core/logx"
)

type PublishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishLogic) Publish(in *post.PublishRequest) (*post.PublishResponse, error) {
	// 1. 生成雪花算法全局唯一博文 ID (为后续分库分表做准备)
	node, err := snowflake.NewNode(2) // 节点ID设为2，区分于用户服务
	if err != nil {
		return nil, fmt.Errorf("生成博文ID失败: %v", err)
	}
	postId := node.Generate().Int64()

	// 2. 组装数据并写入 MySQL
	newPost := &model.Post{
		Id:      postId,
		UserId:  in.UserId,
		Content: in.Content,
	}

	_, err = l.svcCtx.PostModel.InsertWithId(l.ctx, newPost)
	if err != nil {
		return nil, fmt.Errorf("博文插入数据库失败: %v", err)
	}

	// 3. 返回生成的 PostId
	return &post.PublishResponse{
		PostId: postId,
	}, nil
}
