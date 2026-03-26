package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"miniblog/app/interaction/rpc/interaction"
	"miniblog/app/interaction/rpc/internal/svc"
	"miniblog/app/notice/rpc/noticerpc" // 【关键约束】：修正引入包，与评论保持完全一致

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

	// 2. O(1) 复杂度判断用户是否已经点赞
	isMember, err := l.svcCtx.BizRedis.SismemberCtx(l.ctx, likeUsersKey, in.UserId)
	if err != nil {
		l.Logger.Errorf("Redis Sismember 失败: %v", err)
		return nil, fmt.Errorf("系统繁忙，请稍后再试")
	}

	if isMember {
		l.Logger.Infof("=> [RPC逻辑层] 触发拦截器：已拒绝用户 %d 的重复点赞请求", in.UserId)
		return nil, fmt.Errorf("您已经对该博文点过赞了")
	}

	// 3. 执行点赞的缓存操作
	_, _ = l.svcCtx.BizRedis.SaddCtx(l.ctx, likeUsersKey, in.UserId)
	_, _ = l.svcCtx.BizRedis.IncrCtx(l.ctx, likeCountKey)

	// 4. 将落库任务投递到 MQ
	taskData := map[string]interface{}{
		"post_id": in.PostId,
		"user_id": in.UserId,
		"action":  1,
	}
	taskBytes, _ := json.Marshal(taskData)
	_, err = l.svcCtx.BizRedis.LpushCtx(l.ctx, mqKey, string(taskBytes))

	if err != nil {
		l.Logger.Errorf("【警告】点赞任务推入 MQ 失败: %v", err)
	} else {
		l.Logger.Infof("=> [RPC逻辑层] 点赞成功，底层任务已悄悄塞入队列！(UserId:%d)", in.UserId)
	}

	// ==========================================
	// 5. 微服务联动：基于防灾协程的异步通知
	// ==========================================
	threading.GoSafe(func() {
		// 【关键约束】：将查库和 RPC 调用统一塞入后台上下文，彻底解决提前 cancel 问题
		bgCtx := context.Background()

		postInfo, err := l.svcCtx.PostModel.FindOne(bgCtx, in.PostId)
		if err != nil || postInfo == nil {
			logx.Errorf("[异步通知异常] 未能查询到博文详情, postId: %d", in.PostId)
			return
		}

		// 防御：自己赞自己，不发通知
		if postInfo.UserId == in.UserId {
			return
		}

		_, noticeErr := l.svcCtx.NoticeRpc.CreateNotice(bgCtx, &noticerpc.CreateNoticeReq{
			UserId:  postInfo.UserId,
			ActorId: in.UserId,
			PostId:  in.PostId,
			Type:    1,
		})

		if noticeErr != nil {
			logx.Errorf("异步发送点赞通知失败: %v", noticeErr)
		} else {
			logx.Infof("[异步通知成功] 已向用户 %d 推送了来自 %d 的点赞通知", postInfo.UserId, in.UserId)
		}
	})

	return &interaction.LikeResponse{}, nil
}
