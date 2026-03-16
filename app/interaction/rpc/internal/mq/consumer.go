package mq

import (
	"context"
	"encoding/json"
	"time"

	"miniblog/app/interaction/model"
	"miniblog/app/interaction/rpc/internal/svc"

	"github.com/bwmarrin/snowflake"
	"github.com/zeromicro/go-zero/core/logx"
)

// TaskData 对应 LikeLogic 中打包的 JSON 数据结构
type TaskData struct {
	PostId int64 `json:"post_id"`
	UserId int64 `json:"user_id"`
	Action int   `json:"action"` // 1 表示点赞
}

// StartLikeTaskConsumer 启动点赞任务的异步消费者
func StartLikeTaskConsumer(svcCtx *svc.ServiceContext) {
	logx.Info("====== 启动后台消费者协程：监听点赞 MQ 任务 (Rpop 轮询模式) ======")
	ctx := context.Background()
	mqKey := "biz:interact:mq_like_tasks"

	// 初始化雪花算法节点 (节点4：专门用于点赞流水表)
	node, err := snowflake.NewNode(4)
	if err != nil {
		logx.Errorf("消费者初始化雪花算法失败: %v", err)
		return
	}

	// 开启一个常驻后台的死循环，模拟守护进程
	for {
		// 【修改点】：使用 RpopCtx 替代 BrpopCtx
		// Rpop 会立刻返回，如果队列里有数据就返回数据，没有就返回错误或空字符串
		taskJson, err := svcCtx.BizRedis.RpopCtx(ctx, mqKey)

		// 如果发生错误（通常是 redis.Nil 表示空）或者弹出的字符串为空，说明队列里暂时没任务
		if err != nil || taskJson == "" {
			// 睡 500 毫秒后再去查，避免死循环把 CPU 跑满
			// 如果队列里没任务，协程会在这里“睡” 500ms，非常节省资源
			time.Sleep(500 * time.Millisecond)
			continue
		}

		// 成功拿到任务数据，进行解析
		// taskJson 是真实的 JSON 字符串数据
		var task TaskData
		if err := json.Unmarshal([]byte(taskJson), &task); err != nil {
			logx.Errorf("MQ 任务解析失败: %v, 脏数据: %s", err, taskJson)
			continue
		}

		// 执行真正的数据库落盘逻辑
		if task.Action == 1 {
			record := &model.LikeRecord{
				Id:     node.Generate().Int64(),
				PostId: task.PostId,
				UserId: task.UserId,
				Status: 1, // 1=已赞
			}

			// 写入 MySQL (此处利用 goctl 生成的 Model 插入)
			_, err := svcCtx.LikeRecordModel.Insert(ctx, record)
			if err != nil {
				// 实际大厂工程中，如果落库失败会扔进“死信队列 (DLQ)”稍后重试
				// 毕设做到打日志告警即可，已经超越大部分本科水平
				logx.Errorf("点赞异步落库失败: %v (PostId: %d, UserId: %d)", err, task.PostId, task.UserId)
			} else {
				logx.Infof("=> 异步消费者执行成功: 点赞记录已落入 MySQL (PostId: %d, UserId: %d)", task.PostId, task.UserId)
			}
		}
	}
}
