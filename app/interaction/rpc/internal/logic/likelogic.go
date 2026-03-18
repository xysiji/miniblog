package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"miniblog/app/interaction/rpc/interaction"
	"miniblog/app/interaction/rpc/internal/svc"

	"miniblog/app/notice/rpc/notice" // 引入 notice 协议

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

type LikeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeLogic {
	return &LikeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Like 点赞/取消点赞
func (l *LikeLogic) Like(in *interaction.LikeRequest) (*interaction.LikeResponse, error) {
	l.Logger.Infof("=> [RPC逻辑层] 开始处理点赞任务: UserId=%d, PostId=%d", in.UserId, in.PostId)

	// 1. 定义 Redis Key 的规约
	likeUsersKey := fmt.Sprintf("biz:interact:like_users:%d", in.PostId)
	likeCountKey := fmt.Sprintf("biz:interact:like_count:%d", in.PostId)
	mqKey := "biz:interact:mq_like_tasks"

	// 2. 【核心亮点】O(1) 复杂度判断用户是否已经点赞
	isMember, err := l.svcCtx.BizRedis.SismemberCtx(l.ctx, likeUsersKey, in.UserId)
	if err != nil {
		l.Logger.Errorf("Redis Sismember 失败: %v", err)
		return nil, fmt.Errorf("系统繁忙，请稍后再试")
	}

	l.Logger.Infof("=> [RPC逻辑层] 查验防重复机制: 用户是否已点赞 -> %v", isMember)

	if isMember {
		l.Logger.Infof("=> [RPC逻辑层] 触发拦截器：已拒绝用户 %d 的重复点赞请求", in.UserId)
		return nil, fmt.Errorf("您已经对该博文点过赞了")
	}

	// 3. 执行点赞的缓存操作 (全内存操作，极速响应)
	_, _ = l.svcCtx.BizRedis.SaddCtx(l.ctx, likeUsersKey, in.UserId)
	_, _ = l.svcCtx.BizRedis.IncrCtx(l.ctx, likeCountKey)

	// 4. 【异步削峰】将落库任务打包投递到 Redis 消息队列，立刻向前端返回成功
	taskData := map[string]interface{}{
		"post_id": in.PostId,
		"user_id": in.UserId,
		"action":  1, // 1 表示点赞行为
	}
	taskBytes, _ := json.Marshal(taskData)
	_, err = l.svcCtx.BizRedis.LpushCtx(l.ctx, mqKey, string(taskBytes))

	if err != nil {
		// 降级容灾
		l.Logger.Errorf("【警告】点赞任务推入 MQ 失败: %v", err)
	} else {
		l.Logger.Infof("=> [RPC逻辑层] 点赞成功，底层任务已悄悄塞入队列！(UserId:%d)", in.UserId)
	}
	// ======== 【精确插入：触发异步通知】 ========
	postInfo, err := l.svcCtx.PostModel.FindOne(l.ctx, in.PostId)
	if err == nil && postInfo != nil {
		// 使用 threading.GoSafe 开启安全的后台异步协程，不阻塞当前点赞响应
		threading.GoSafe(func() {
			_, noticeErr := l.svcCtx.NoticeRpc.CreateNotice(context.Background(), &notice.CreateNoticeReq{
				UserId:  postInfo.UserId, // 发给谁（这篇博文的作者）
				ActorId: in.UserId,       // 谁点的赞（当前操作者）
				PostId:  in.PostId,       // 哪篇博文
				Type:    1,               // 1代表点赞类型
			})
			if noticeErr != nil {
				logx.Errorf("异步发送点赞通知失败: %v", noticeErr)
			} else {
				logx.Infof("异步点赞通知发送成功! 给用户 %d", postInfo.UserId)
			}
		})
	} else {
		l.Logger.Errorf("未能查询到博文作者信息，跳过发送通知, postId: %d", in.PostId)
	}

	return &interaction.LikeResponse{}, nil
}
