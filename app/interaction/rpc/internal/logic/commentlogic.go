package logic

import (
	"context"
	"fmt"

	"miniblog/app/interaction/model"
	"miniblog/app/interaction/rpc/interaction"
	"miniblog/app/interaction/rpc/internal/svc"
	"miniblog/app/notice/rpc/noticerpc" // 【关键约束】：严格引入 notice 的 RPC client 包

	"github.com/bwmarrin/snowflake"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

type CommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentLogic {
	return &CommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Comment 发表评论
func (l *CommentLogic) Comment(in *interaction.CommentRequest) (*interaction.CommentResponse, error) {
	// 1. 生成评论专属的雪花算法 ID
	node, err := snowflake.NewNode(3)
	if err != nil {
		return nil, fmt.Errorf("生成评论ID失败: %v", err)
	}
	commentId := node.Generate().Int64()

	// 2. 组装数据，直接同步写入 MySQL
	newComment := &model.Comment{
		Id:      commentId,
		PostId:  in.PostId,
		UserId:  in.UserId,
		Content: in.Content,
	}

	_, err = l.svcCtx.CommentModel.Insert(l.ctx, newComment)
	if err != nil {
		l.Logger.Errorf("评论写入 MySQL 失败: %v", err)
		return nil, fmt.Errorf("评论发布失败，请稍后重试")
	}

	l.Logger.Infof("用户 %d 对博文 %d 发表了评论", in.UserId, in.PostId)

	// ==========================================
	// 3. 微服务联动：基于防灾协程的异步通知
	// ==========================================
	threading.GoSafe(func() {
		// 【关键约束】：为后台任务开辟全新上下文，脱离主请求生命周期
		bgCtx := context.Background()

		// 3.1 异步查询博文详情（查出作者是谁）
		postData, err := l.svcCtx.PostModel.FindOne(bgCtx, in.PostId)
		if err != nil {
			logx.Errorf("[异步通知异常] 查询博文详情失败, PostId: %d, err: %v", in.PostId, err)
			return
		}

		// 3.2 业务防御：A账户评论A账户的博文，不需要发通知
		if postData.UserId == in.UserId {
			return
		}

		// 3.3 联动 NoticeRpc：向作者发送系统通知 (Type=2)
		_, err = l.svcCtx.NoticeRpc.CreateNotice(bgCtx, &noticerpc.CreateNoticeReq{
			UserId:  postData.UserId,
			ActorId: in.UserId,
			PostId:  in.PostId,
			Type:    2,
		})

		if err != nil {
			logx.Errorf("[异步通知异常] 调用 NoticeRpc 生成评论通知失败, PostId: %d, err: %v", in.PostId, err)
		} else {
			logx.Infof("[异步通知成功] 已向用户 %d 推送了来自 %d 的评论通知", postData.UserId, in.UserId)
		}
	})

	return &interaction.CommentResponse{
		CommentId: commentId,
	}, nil
}
