package logger

import (
	"encoding/json"
	"os"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type ActionLog struct {
	Timestamp  string `json:"timestamp"`
	UserId     int64  `json:"user_id"`
	PostId     int64  `json:"post_id"`
	ActionType string `json:"action_type"` // "LIKE" 或 "COMMENT"
	CostTimeMs int64  `json:"cost_time_ms"`
}

// 异步记录行为日志至本地文件（模拟大数据 Flume/Logstash 采集源）
func RecordAction(userId, postId int64, actionType string, costMs int64) {
	go func() {
		logData := ActionLog{
			Timestamp:  time.Now().Format(time.RFC3339),
			UserId:     userId,
			PostId:     postId,
			ActionType: actionType,
			CostTimeMs: costMs,
		}
		bytes, _ := json.Marshal(logData)

		// 追加写入 access_data.log
		f, err := os.OpenFile("access_data.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logx.Errorf("打开日志文件失败: %v", err)
			return
		}
		defer f.Close()
		f.WriteString(string(bytes) + "\n")
	}()
}
