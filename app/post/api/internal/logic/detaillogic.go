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

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取单条博文详情
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.DetailReq) (resp *types.DetailResp, err error) {
	rpcResp, err := l.svcCtx.PostRpc.Detail(l.ctx, &postclient.DetailRequest{
		PostId: req.PostId,
	})
	if err != nil {
		return nil, err
	}

	// ✅ 修复：初始化为空切片而不是 nil，保证 JSON 输出为 []
	images := make([]string, 0)

	// ✅ 修复：防止底层返回的 Post 为空指针引发 Panic
	var postItem types.PostItem
	if rpcResp.Post != nil {
		if rpcResp.Post.Images != "" && rpcResp.Post.Images != "null" {
			_ = json.Unmarshal([]byte(rpcResp.Post.Images), &images)
		}

		postItem = types.PostItem{
			Id:       rpcResp.Post.Id,
			UserId:   rpcResp.Post.UserId,
			Content:  rpcResp.Post.Content,
			Images:   images,
			CreateAt: time.Unix(rpcResp.Post.CreateTime, 0).Format("2006-01-02 15:04:05"),
		}
	}

	return &types.DetailResp{
		Post: postItem,
	}, nil
}
