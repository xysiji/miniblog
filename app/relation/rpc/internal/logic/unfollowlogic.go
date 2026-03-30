package logic

import (
	"context"
	"fmt"

	"miniblog/app/relation/model"
	"miniblog/app/relation/rpc/internal/svc"
	"miniblog/app/relation/rpc/relation"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnfollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnfollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnfollowLogic {
	return &UnfollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnfollowLogic) Unfollow(in *relation.UnfollowReq) (*relation.UnfollowResp, error) {
	// 1. 查出关系记录的主键 ID
	rel, err := l.svcCtx.RelationModel.FindOneByFollowerIdFollowingId(l.ctx, in.FollowerId, in.FollowingId)
	if err != nil {
		if err == model.ErrNotFound {
			// 本来就没有关注，直接当做取消成功返回（幂等）
			return &relation.UnfollowResp{}, nil
		}
		l.Logger.Errorf("查询关注记录失败: %v", err)
		return nil, fmt.Errorf("操作失败，请稍后重试")
	}

	// 2. 根据主键 ID 执行物理删除
	err = l.svcCtx.RelationModel.Delete(l.ctx, rel.Id)
	if err != nil {
		l.Logger.Errorf("删除关注记录失败: %v", err)
		return nil, fmt.Errorf("取消关注失败，请稍后重试")
	}

	return &relation.UnfollowResp{}, nil
}
