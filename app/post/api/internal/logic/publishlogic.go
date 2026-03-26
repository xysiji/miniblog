// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"miniblog/app/post/api/internal/svc"
	"miniblog/app/post/api/internal/types"
	"miniblog/app/post/rpc/postclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发布微型博客
func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishLogic) Publish(req *types.PublishReq) (resp *types.PublishResp, err error) {
	// 1. 从 JWT 解析后的 Context 中提取当前登录用户的 userId
	// 注意：go-zero 解析出来的数据默认是 json.Number 类型
	userIdNumber, ok := l.ctx.Value("userId").(json.Number)
	if !ok {
		return nil, fmt.Errorf("非法请求：未获取到有效的用户凭证")
	}

	userId, err := userIdNumber.Int64()
	if err != nil {
		return nil, fmt.Errorf("用户凭证解析失败")
	}

	// ==================== 新增：序列化图片数组逻辑 ====================
	var imagesBytes []byte
	if len(req.Images) > 0 {
		imagesBytes, err = json.Marshal(req.Images)
		if err != nil {
			return nil, fmt.Errorf("图片数据处理失败: %v", err)
		}
	}
	// ==================================================================

	// 2. 调用底层 Post RPC 服务，传入从 Token 解析出的 UserId 和前端传来的内容
	rpcRes, err := l.svcCtx.PostRpc.Publish(l.ctx, &postclient.PublishRequest{
		UserId:  userId,
		Content: req.Content,
		Images:  string(imagesBytes), // 新增：将序列化后的图片 JSON 字符串传给底层的 RPC 服务
	})
	if err != nil {
		return nil, err
	}

	// 3. 组装 HTTP 响应
	return &types.PublishResp{
		PostId: rpcRes.PostId,
	}, nil
}
