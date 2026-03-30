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

	// 嵌套逻辑：判断是否为顶级主评论
	rootId := in.RootId
	if rootId == 0 {
		rootId = 0
	}

	newComment := &model.Comment{
		Id:            commentId,
		PostId:        in.PostId,
		UserId:        in.UserId,
		RootId:        rootId,           // 补充嵌套字段
		ParentId:      in.ParentId,      // 补充嵌套字段
		ReplyToUserId: in.ReplyToUserId, // 补充嵌套字段
		Content:       in.Content,
		Status:        1,
	}

	// ⚠️架构师修改：使用你定义的 InsertShard 才能正确写入分表
	_, err = l.svcCtx.CommentModel.InsertShard(l.ctx, newComment)
	if err != nil {
		l.Logger.Errorf("评论写入分表失败: %v", err)
		return nil, fmt.Errorf("评论发布失败，请稍后重试")
	}

	l.Logger.Infof("用户 %d 对博文 %d 发表了评论", in.UserId, in.PostId)

	// ==========================================
	// 3. 微服务联动：基于防灾协程的异步通知与统计更新 (严格保留源码)
	// ==========================================
	threading.GoSafe(func() {
		bgCtx := context.Background()

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
