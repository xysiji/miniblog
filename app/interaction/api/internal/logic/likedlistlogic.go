package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv" // 引入 strconv 用于数字转字符串

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

	postIds, err := l.svcCtx.LikeRecordModel.FindLikedPostIdsByUserId(l.ctx, userId)
	if err != nil {
		l.Logger.Errorf("查询用户点赞列表失败: userId=%d, err=%v", userId, err)
		return nil, fmt.Errorf("查询点赞列表失败")
	}

	// 将查询到的 []int64 转换为 []string，以避免前端 JS 精度丢失
	strPostIds := make([]string, 0, len(postIds))
	if postIds != nil {
		for _, id := range postIds {
			strPostIds = append(strPostIds, strconv.FormatInt(id, 10))
		}
	}

	return &types.LikedListResp{
		PostIds: strPostIds,
	}, nil
}
