package logic

import (
	"context"
	"fmt"

	"miniblog/app/interaction/model"
	"miniblog/app/interaction/rpc/interaction"
	"miniblog/app/interaction/rpc/internal/svc"
	"miniblog/app/notice/rpc/noticerpc"

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

func (l *CommentLogic) Comment(in *interaction.CommentRequest) (*interaction.CommentResponse, error) {
	node, err := snowflake.NewNode(3)
	if err != nil {
		return nil, fmt.Errorf("生成评论ID失败: %v", err)
	}
	commentId := node.Generate().Int64()

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
	// 3. 微服务联动：基于防灾协程的异步通知与统计更新
	// ==========================================
	threading.GoSafe(func() {
		bgCtx := context.Background()

		// 【核心闭环】：评论成功后，原子增加博文表的评论统计数
		countErr := l.svcCtx.PostModel.IncrCommentCount(bgCtx, in.PostId)
		if countErr != nil {
			logx.Errorf("异步更新评论数失败: %v", countErr)
		}

		postData, err := l.svcCtx.PostModel.FindOne(bgCtx, in.PostId)
		if err != nil {
			logx.Errorf("[异步通知异常] 查询博文详情失败, PostId: %d, err: %v", in.PostId, err)
			return
		}

		if postData.UserId == in.UserId {
			return
		}

		_, err = l.svcCtx.NoticeRpc.CreateNotice(bgCtx, &noticerpc.CreateNoticeReq{
			UserId:  postData.UserId,
			ActorId: in.UserId,
			PostId:  in.PostId,
			Type:    2,
		})

		if err != nil {
			logx.Errorf("[异步通知异常] 调用 NoticeRpc 生成通知失败, err: %v", err)
		} else {
			logx.Infof("[异步通知成功] 已推送评论通知", postData.UserId, in.UserId)
		}
	})

	return &interaction.CommentResponse{
		CommentId: commentId,
	}, nil
}
