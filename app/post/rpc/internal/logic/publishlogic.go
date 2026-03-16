package logic

import (
	"context"
	"fmt"
	"time" // 新增：为了获取当前时间戳作为 Redis ZSet 的分数

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

	// ==========================================
	// 本次新增核心逻辑：3. 同步写入 Redis Timeline
	// ==========================================
	// 使用 Redis 的 ZSet (有序集合)，按时间戳排序
	timelineKey := "biz:post:global_timeline"
	// 将博文 ID 转为字符串存入 Redis，分数为当前 Unix 时间戳
	_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, timelineKey, time.Now().Unix(), fmt.Sprintf("%d", postId))
	if err != nil {
		// 缓存降级：就算 Redis 挂了，博文已经发到 MySQL 了，不能阻断用户，所以只打 Error 日志不 return err
		l.Logger.Errorf("【缓存警告】博文 %d 写入 Redis Timeline 失败: %v", postId, err)
	} else {
		l.Logger.Infof("博文 %d 已成功推送到全局 Timeline", postId)
	}

	// 4. 返回生成的 PostId
	return &post.PublishResponse{
		PostId: postId,
	}, nil
}
