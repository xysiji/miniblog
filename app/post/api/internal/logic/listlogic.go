// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"miniblog/app/interaction/rpc/interactionclient"
	"miniblog/app/post/api/internal/svc"
	"miniblog/app/post/api/internal/types"
	"miniblog/app/post/rpc/postclient"
	"miniblog/app/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req *types.ListReq) (resp *types.ListResp, err error) {
	// 1. 获取底层 Post RPC 分页数据
	rpcResp, err := l.svcCtx.PostRpc.List(l.ctx, &postclient.ListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	// 2. 从上下文中提取 userId（由 SoftAuthMiddleware 静默注入）
	var currentUserId int64
	if uidVal := l.ctx.Value("userId"); uidVal != nil {
		switch v := uidVal.(type) {
		case json.Number:
			currentUserId, _ = v.Int64()
		case float64:
			currentUserId = int64(v)
		case string:
			fmt.Sscanf(v, "%d", &currentUserId)
		}
	}

	// 3. 异步并发获取点赞状态，避免阻塞主流程
	likedMap := make(map[int64]bool)
	if currentUserId > 0 {
		likedResp, err := l.svcCtx.InteractionRpc.LikedList(l.ctx, &interactionclient.LikedListRequest{
			UserId: currentUserId,
		})
		if err == nil && likedResp != nil {
			for _, pid := range likedResp.PostIds {
				likedMap[pid] = true
			}
		} else {
			l.Logger.Errorf("获取用户点赞列表失败: %v", err)
		}
	}

	// 4. 提取本次查询的所有作者ID (去重)
	userIdSet := make(map[int64]struct{})
	uniqueUserIds := make([]int64, 0)
	if rpcResp.List != nil {
		for _, item := range rpcResp.List {
			if _, exists := userIdSet[item.UserId]; !exists {
				userIdSet[item.UserId] = struct{}{}
				uniqueUserIds = append(uniqueUserIds, item.UserId)
			}
		}
	}

	// 5. 呼叫 User RPC 的批量获取接口 (完美解决 N+1)
	userInfoMap := make(map[int64]*userclient.UserInfoResponse)
	if len(uniqueUserIds) > 0 {
		batchResp, err := l.svcCtx.UserRpc.BatchGetUserInfo(l.ctx, &userclient.BatchUserInfoReq{
			UserIds: uniqueUserIds,
		})
		if err == nil && batchResp != nil && batchResp.Users != nil {
			userInfoMap = batchResp.Users
		} else {
			l.Logger.Errorf("批量获取用户信息失败: %v", err)
		}
	}

	// 6. 最终数据拼装清洗
	list := make([]types.PostItem, 0)
	if rpcResp.List != nil {
		for _, item := range rpcResp.List {
			images := make([]string, 0)
			if item.Images != "" && item.Images != "null" {
				_ = json.Unmarshal([]byte(item.Images), &images)
			}

			// 获取该条记录对应的作者信息
			authorName := "匿名用户"
			authorAvatar := ""
			if uInfo, ok := userInfoMap[item.UserId]; ok {
				authorName = uInfo.Username
				authorAvatar = uInfo.Avatar
			}

			list = append(list, types.PostItem{
				Id:           item.Id,
				UserId:       item.UserId,
				AuthorName:   authorName,
				AuthorAvatar: authorAvatar,
				Content:      item.Content,
				Images:       images,
				LikeCount:    item.LikeCount,
				CommentCount: item.CommentCount,
				IsLiked:      likedMap[item.Id],
				CreateAt:     time.Unix(item.CreateTime, 0).Format("2006-01-02 15:04:05"),
			})
		}
	}

	return &types.ListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}
