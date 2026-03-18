package logic

import (
	"context"
	"encoding/json"
	"errors" // 【修复1】：引入 Go 语言标准的 errors 包

	"miniblog/app/notice/api/internal/svc"
	"miniblog/app/notice/api/internal/types"

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

func (l *ListLogic) List(req *types.NoticeListReq) (resp *types.NoticeListResp, err error) {
	// 1. 【安全核心】：从 JWT Token 中解析出当前用户的 userId (Go-zero 默认解析为 json.Number)
	userIdNumber, ok := l.ctx.Value("userId").(json.Number)
	if !ok {
		// 【修复2】：使用标准库的 errors.New 替代 logx.ApiError
		return nil, errors.New("未授权的访问")
	}
	userId, _ := userIdNumber.Int64()

	// 2. 调用 Model 层，执行数据库查询
	notices, err := l.svcCtx.NoticeModel.FindListByUserId(l.ctx, userId, req.Page, req.PageSize)
	if err != nil {
		l.Logger.Errorf("查询通知列表失败: %v", err)
		// 【修复3】：使用标准库的 errors.New 替代 logx.ApiError
		return nil, errors.New("获取通知失败，请稍后再试")
	}

	// 3. 查询此用户的总通知数 (用于前端分页)
	total, _ := l.svcCtx.NoticeModel.CountByUserId(l.ctx, userId)

	// 4. 将数据库的 Model 转换为返回给前端的 API 结构体类型 (数据脱敏与映射)
	var list []types.NoticeItem
	for _, n := range notices {
		list = append(list, types.NoticeItem{
			Id:         n.Id,
			ActorId:    n.ActorId,
			PostId:     n.PostId,
			Type:       int(n.Type),
			IsRead:     int(n.IsRead),
			CreateTime: n.CreateTime.Unix(), // 转为时间戳给前端，方便前端格式化
		})
	}

	// 5. 完美组装，吐给前端
	return &types.NoticeListResp{
		List:  list,
		Total: total,
	}, nil
}
