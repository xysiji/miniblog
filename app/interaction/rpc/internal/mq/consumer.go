package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"miniblog/app/interaction/model"
	"miniblog/app/interaction/rpc/internal/svc"
	"miniblog/app/notice/rpc/noticerpc" // 必须引入通知服务

	"github.com/bwmarrin/snowflake"
	"github.com/zeromicro/go-zero/core/logx"
)

type TaskData struct {
	PostId int64 `json:"post_id"`
	UserId int64 `json:"user_id"`
	Action int   `json:"action"` // 1: 点赞, 2: 取消
}

func StartLikeTaskConsumer(svcCtx *svc.ServiceContext) {
	logx.Info("====== 启动后台消费者协程：监听点赞 MQ 任务 ======")
	ctx := context.Background()
	mqKey := "biz:interact:mq_like_tasks"

	node, _ := snowflake.NewNode(4)

	for {
		taskJson, err := svcCtx.BizRedis.RpopCtx(ctx, mqKey)
		if err != nil || taskJson == "" {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		var task TaskData
		if err := json.Unmarshal([]byte(taskJson), &task); err != nil {
			continue
		}

		// 定义缓存 Key
		cacheKey := fmt.Sprintf("cache:post:id:%d", task.PostId)

		if task.Action == 1 {
			// 防止重复点击导致的唯一索引冲突 (Upsert逻辑)
			record, err := svcCtx.LikeRecordModel.FindOneByPostIdUserId(ctx, task.PostId, task.UserId)
			isNewLike := false

			if err == model.ErrNotFound {
				// 没点过，插入
				newRecord := &model.LikeRecord{
					Id:     node.Generate().Int64(),
					PostId: task.PostId,
					UserId: task.UserId,
					Status: 1,
				}
				if _, err := svcCtx.LikeRecordModel.Insert(ctx, newRecord); err == nil {
					isNewLike = true
				}
			} else if record != nil && record.Status == 0 {
				// 曾经取消过，更新为 1
				record.Status = 1
				if err := svcCtx.LikeRecordModel.Update(ctx, record); err == nil {
					isNewLike = true
				}
			}

			if isNewLike {
				// 1. 更新 MySQL 数字
				_ = svcCtx.PostModel.IncrLikeCount(ctx, task.PostId)

				// 2. 【核心修复】：删除旧缓存，保证刷新网页不归零 (两个下划线 _, _)
				_, _ = svcCtx.BizRedis.DelCtx(ctx, cacheKey)

				// 3. 【核心修复】：发送通知
				postData, err := svcCtx.PostModel.FindOne(ctx, task.PostId)
				if err == nil && postData.UserId != task.UserId && svcCtx.NoticeRpc != nil {
					_, _ = svcCtx.NoticeRpc.CreateNotice(ctx, &noticerpc.CreateNoticeReq{
						UserId:  postData.UserId,
						ActorId: task.UserId,
						PostId:  task.PostId,
						Type:    1,
					})
				}
			}
		}
	}
}
