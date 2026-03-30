package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"miniblog/app/interaction/api/internal/svc"
	"miniblog/app/interaction/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikedListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikedListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikedListLogic {
	return &LikedListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikedListLogic) LikedList(req *types.LikedListReq) (resp *types.LikedListResp, err error) {
	var userId int64
	if uidNum, ok := l.ctx.Value("userId").(json.Number); ok {
		userId, _ = uidNum.Int64()
	}

	if userId == 0 {
		l.Logger.Errorf("获取用户信息失败，Token异常")
		return nil, fmt.Errorf("未授权访问，请重新登录")
	}

	// 直接从底层拿到的也是 []int64
	postIds, err := l.svcCtx.LikeRecordModel.FindLikedPostIdsByUserId(l.ctx, userId)
	if err != nil {
		l.Logger.Errorf("查询用户点赞列表失败: userId=%d, err=%v", userId, err)
		return nil, fmt.Errorf("查询点赞列表失败")
	}

	// 如果查询结果为 nil，初始化为一个空切片，保证前端拿到 [] 而不是 null
	if postIds == nil {
		postIds = make([]int64, 0)
	}

	// 【修复点】：直接返回 []int64，与 types.go 保持完全一致
	return &types.LikedListResp{
		PostIds: postIds,
	}, nil
}
