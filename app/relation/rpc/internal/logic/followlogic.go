package logic

import (
	"context"
	"fmt"
	"strings"

	"miniblog/app/relation/model"
	"miniblog/app/relation/rpc/internal/svc"
	"miniblog/app/relation/rpc/relation"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowLogic) Follow(in *relation.FollowReq) (*relation.FollowResp, error) {
	// 1. 基础校验：不能自己关注自己
	if in.FollowerId == in.FollowingId {
		return nil, fmt.Errorf("不能关注自己")
	}

	// 2. 幂等性校验：先查一次是否已经关注过
	// goctl 会根据我们的 UNIQUE KEY `idx_follower_following` 自动生成这个查询方法
	_, err := l.svcCtx.RelationModel.FindOneByFollowerIdFollowingId(l.ctx, in.FollowerId, in.FollowingId)
	if err == nil {
		// 记录已存在，说明已经关注过了，直接静默返回成功（幂等）
		return &relation.FollowResp{}, nil
	}
	if err != model.ErrNotFound {
		l.Logger.Errorf("查询关注记录失败: %v", err)
		return nil, fmt.Errorf("服务器开小差了，请稍后重试")
	}

	// 3. 执行关注插入
	_, err = l.svcCtx.RelationModel.Insert(l.ctx, &model.Relation{
		FollowerId:  in.FollowerId,
		FollowingId: in.FollowingId,
	})

	if err != nil {
		// 4. 极端并发兜底：两条请求同时突破了上面的查询校验，触发 MySQL 唯一索引报错
		if strings.Contains(err.Error(), "Duplicate entry") {
			return &relation.FollowResp{}, nil
		}
		l.Logger.Errorf("插入关注记录失败: %v", err)
		return nil, fmt.Errorf("关注失败，请稍后重试")
	}

	return &relation.FollowResp{}, nil
}
