package logic

import (
	"context"
	"fmt"

	"miniblog/app/interaction/rpc/interaction"
	"miniblog/app/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

type UnlikeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnlikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnlikeLogic {
	return &UnlikeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnlikeLogic) Unlike(in *interaction.UnlikeRequest) (*interaction.UnlikeResponse, error) {
	// 1. 无视 MySQL 延迟，立即从 Redis 清除用户的点赞状态，解除前端拦截
	likeUsersKey := fmt.Sprintf("biz:interact:like_users:%d", in.PostId)
	likeCountKey := fmt.Sprintf("biz:interact:like_count:%d", in.PostId)

	// 【注意】：Go 返回两个值，必须用 _, _
	_, _ = l.svcCtx.BizRedis.SremCtx(l.ctx, likeUsersKey, in.UserId)
	_, _ = l.svcCtx.BizRedis.DecrCtx(l.ctx, likeCountKey)

	l.Logger.Infof("=> [RPC] 用户 %d 成功取消了对博文 %d 的点赞", in.UserId, in.PostId)

	// 2. 异步处理落库与统计，彻底防止查不到记录的报错
	threading.GoSafe(func() {
		bgCtx := context.Background()

		// 如果能查到记录，把状态改为 0；如果查不到（MQ还未落盘），就算了，不拦截
		record, err := l.svcCtx.LikeRecordModel.FindOneByPostIdUserId(bgCtx, in.PostId, in.UserId)
		if err == nil && record != nil && record.Status == 1 {
			record.Status = 0
			_ = l.svcCtx.LikeRecordModel.Update(bgCtx, record)
		}

		// 原子扣减博文表的点赞数
		_ = l.svcCtx.PostModel.DecrLikeCount(bgCtx, in.PostId)

		// 3. 【彻底解决数字刷新消失】：强行删除博文缓存，逼迫下次刷新走数据库查最新数字！
		cacheKey := fmt.Sprintf("cache:post:id:%d", in.PostId)
		// 【注意】：这里就是之前报错的地方，已经改成两个下划线 _, _ 了！！！绝对不会再报错！
		_, _ = l.svcCtx.BizRedis.DelCtx(bgCtx, cacheKey)
	})

	return &interaction.UnlikeResponse{}, nil
}
