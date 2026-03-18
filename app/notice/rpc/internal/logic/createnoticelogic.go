package logic

import (
	"context"
	"fmt"
	"time"

	"miniblog/app/notice/model"
	"miniblog/app/notice/rpc/internal/svc"
	"miniblog/app/notice/rpc/notice"

	"github.com/bwmarrin/snowflake"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateNoticeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateNoticeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNoticeLogic {
	return &CreateNoticeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateNotice 内部接口：供 Interaction 服务异步调用来生成通知
func (l *CreateNoticeLogic) CreateNotice(in *notice.CreateNoticeReq) (*notice.CreateNoticeResp, error) {
	// 1. 【业务防呆拦截】：自己给自己点赞/评论，不需要发通知
	if in.UserId == in.ActorId {
		l.Logger.Infof("用户 %d 自己操作了自己的博文，跳过生成通知", in.UserId)
		return &notice.CreateNoticeResp{}, nil
	}

	// 2. 生成雪花算法全局唯一 ID (为通知服务分配节点 ID: 3)
	node, err := snowflake.NewNode(3)
	if err != nil {
		return nil, fmt.Errorf("生成通知分布式ID失败: %v", err)
	}
	noticeId := node.Generate().Int64()

	// 3. 组装数据库结构体数据
	newNotice := &model.Notice{
		Id:         noticeId,
		UserId:     in.UserId,      // 接收方
		ActorId:    in.ActorId,     // 操作方 (谁点的赞)
		PostId:     in.PostId,      // 哪篇博文
		Type:       int64(in.Type), // 1-点赞 2-评论
		IsRead:     0,              // 默认 0 为未读
		CreateTime: time.Now(),
	}

	// 4. 执行写入 MySQL
	_, err = l.svcCtx.NoticeModel.Insert(l.ctx, newNotice)
	if err != nil {
		l.Logger.Errorf("=> [Notice-RPC] 通知落库失败: %v", err)
		return nil, fmt.Errorf("通知落库失败")
	}

	l.Logger.Infof("=> [Notice-RPC] 成功生成一条新通知，已落库！通知ID: %d, 接收人: %d", noticeId, in.UserId)

	return &notice.CreateNoticeResp{}, nil
}
