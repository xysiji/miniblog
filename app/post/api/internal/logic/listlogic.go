// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"encoding/json"
	"time"

	"miniblog/app/post/api/internal/svc"
	"miniblog/app/post/api/internal/types"
	"miniblog/app/post/rpc/postclient"

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
	// 1. 发起呼叫：调用底层 Post RPC 服务，索要分页数据
	rpcResp, err := l.svcCtx.PostRpc.List(l.ctx, &postclient.ListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	// 2. 数据清洗与转化 (BFF 模式核心表现)
	// 【安全防范】：使用 make 初始化切片，确保就算没数据给前端的也是 [] 而不是 null
	list := make([]types.PostItem, 0)

	if rpcResp.List != nil {
		for _, item := range rpcResp.List {
			// 【新增】：将底层的 JSON 字符串转化为前端需要的 []string 切片
			images := make([]string, 0)
			if item.Images != "" && item.Images != "null" {
				_ = json.Unmarshal([]byte(item.Images), &images)
			}

			list = append(list, types.PostItem{
				Id:           item.Id,
				UserId:       item.UserId,
				Content:      item.Content,
				Images:       images,            // 映射到切片
				LikeCount:    item.LikeCount,    // 【精准补齐】：映射点赞数
				CommentCount: item.CommentCount, // 【精准补齐】：映射评论数
				// 核心转化：利用 Go 独特的 2006-01-02 15:04:05 诞生时间来格式化
				CreateAt: time.Unix(item.CreateTime, 0).Format("2006-01-02 15:04:05"),
			})
		}
	}

	// 3. 组装最终给前端的 JSON 响应
	return &types.ListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}
