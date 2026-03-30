package logic

import (
	"context"
	"fmt"

	"miniblog/app/post/rpc/internal/svc"
	"miniblog/app/post/rpc/post"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Detail 获取博文详情
func (l *DetailLogic) Detail(in *post.DetailRequest) (*post.DetailResponse, error) {
	postInfo, err := l.svcCtx.PostSlaveModel.FindOne(l.ctx, in.PostId)
	if err != nil {
		return nil, fmt.Errorf("博文不存在或查询失败")
	}

	// ✅ 修复：处理 sql.NullString
	imagesStr := ""
	if postInfo.Images.Valid {
		imagesStr = postInfo.Images.String
	}

	return &post.DetailResponse{
		Post: &post.PostItem{
			Id:         postInfo.Id,
			UserId:     postInfo.UserId,
			Content:    postInfo.Content,
			Images:     imagesStr,
			CreateTime: postInfo.CreateTime.Unix(),
		},
	}, nil
}
