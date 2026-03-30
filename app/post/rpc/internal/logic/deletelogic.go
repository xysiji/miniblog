package logic

import (
	"context"
	"fmt"
	"strconv"

	"miniblog/app/post/rpc/internal/svc"
	"miniblog/app/post/rpc/post"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteLogic) Delete(in *post.DeleteRequest) (*post.DeleteResponse, error) {
	// 1. 鉴权查询：由于是删除前校验，为了防止主从延迟，可以直接查主库
	postInfo, err := l.svcCtx.PostMasterModel.FindOne(l.ctx, in.PostId)
	if err != nil {
		return nil, fmt.Errorf("博文不存在")
	}

	// 2. 确保只能删除自己的博文
	if postInfo.UserId != in.UserId {
		return nil, fmt.Errorf("越权操作：只能删除自己的博文")
	}

	// 3. 【架构规范】：写操作路由到 Master 主库
	err = l.svcCtx.PostMasterModel.Delete(l.ctx, in.PostId)
	if err != nil {
		return nil, fmt.Errorf("删除博文失败")
	}

	// 4. 清理 Redis Timeline 缓存中的该条记录
	timelineKey := "biz:post:timeline:global"
	_, _ = l.svcCtx.BizRedis.ZremCtx(l.ctx, timelineKey, strconv.FormatInt(in.PostId, 10))

	return &post.DeleteResponse{}, nil
}
