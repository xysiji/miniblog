package logic

import (
	"context"
	"fmt"
	"time"

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

// Publish 发布微型博客
func (l *PublishLogic) Publish(in *post.PublishRequest) (*post.PublishResponse, error) {
	// 1. 生成雪花算法全局唯一博文 ID (为后续分库分表做准备，节点ID为2)
	node, err := snowflake.NewNode(2)
	if err != nil {
		return nil, fmt.Errorf("生成分布式ID失败: %v", err)
	}
	postId := node.Generate().Int64()

	// 2. 组装待插入的数据库模型数据

	newPost := &model.Post{
		Id:         postId,
		UserId:     in.UserId,
		Content:    in.Content,
		CreateTime: time.Now(), // ✅ 修正为 CreateTime
	}
	// 3. 【分布式架构核心】：写操作强制路由到 Master 主库 (PostMasterModel)
	_, err = l.svcCtx.PostMasterModel.Insert(l.ctx, newPost)
	if err != nil {
		l.Logger.Errorf("写入主库失败: %v", err)
		return nil, fmt.Errorf("发布失败，请稍后重试")
	}

	// 4. 将新发布的博文推送到全局 Timeline (Redis ZSet，Score 为时间戳)
	timelineKey := "biz:post:timeline:global"
	score := time.Now().Unix()
	_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, timelineKey, score, fmt.Sprintf("%d", postId))
	if err != nil {
		l.Logger.Errorf("同步到 Timeline 缓存失败: %v", err)
	} else {
		l.Logger.Infof("博文 %d 已成功推送到全局 Timeline, 路由: Master", postId)
	}

	return &post.PublishResponse{
		PostId: postId,
	}, nil
}
