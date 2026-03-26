package logic

import (
	"context"
	"fmt"

	"miniblog/app/interaction/model" // 引入你的 model
	"miniblog/app/interaction/rpc/interaction"
	"miniblog/app/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
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
	// 1. 查询点赞记录是否存在
	record, err := l.svcCtx.LikeRecordModel.FindOneByPostIdUserId(l.ctx, in.PostId, in.UserId)
	if err != nil && err != model.ErrNotFound {
		l.Logger.Errorf("查询点赞记录失败: %v", err)
		return nil, fmt.Errorf("操作失败，请稍后重试")
	}

	// 2. 如果找到了记录，且当前状态是点赞(1)，则修改为取消点赞(0)
	if record != nil && record.Status == 1 {
		record.Status = 0
		err = l.svcCtx.LikeRecordModel.Update(l.ctx, record)
		if err != nil {
			l.Logger.Errorf("更新点赞状态失败: %v", err)
			return nil, fmt.Errorf("取消点赞失败")
		}

		l.Logger.Infof("用户 %d 成功取消了对博文 %d 的点赞", in.UserId, in.PostId)
		// 进阶扩展点：这里可以发送消息队列，通知 Post 服务将博文的 like_count 减 1
	}

	return &interaction.UnlikeResponse{}, nil
}
