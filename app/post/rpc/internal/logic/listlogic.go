package logic

import (
	"context"

	"miniblog/app/post/rpc/internal/svc"
	"miniblog/app/post/rpc/post"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 新增：内部获取列表的 RPC 方法
func (l *ListLogic) List(in *post.ListRequest) (*post.ListResponse, error) {
	// 1. 参数防御与默认值设置 (防呆设计)
	if in.Page < 1 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 10
	}

	// 2. 第一步：先查总条数
	total, err := l.svcCtx.PostModel.Count(l.ctx)
	if err != nil {
		return nil, err
	}

	// 3. 性能优化：如果总数为 0，直接返回空列表，不要再去查数据库了 (短路返回)
	if total == 0 {
		return &post.ListResponse{
			List:  []*post.PostItem{},
			Total: 0,
		}, nil
	}

	// 4. 第二步：查出当前页的博文数据
	posts, err := l.svcCtx.PostModel.FindPageListByPage(l.ctx, in.Page, in.PageSize)
	if err != nil {
		return nil, err
	}

	// 5. 数据转换 (DTO 映射)：将底层的 MySQL Model 结构，转换为网络传输的 Protobuf 结构
	var respList []*post.PostItem
	for _, p := range posts {
		respList = append(respList, &post.PostItem{
			Id:      p.Id,
			UserId:  p.UserId,
			Content: p.Content,
			// 注意：数据库通常会自动生成 create_time，这里转成 Unix 时间戳 (秒) 方便网络传输
			CreateTime: p.CreateTime.Unix(),
		})
	}

	// 6. 返回最终拼装好的结果集
	return &post.ListResponse{
		List:  respList,
		Total: total,
	}, nil
}
