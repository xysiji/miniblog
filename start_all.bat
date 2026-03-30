@echo off
chcp 65001 >nul
echo ===================================================
echo      Miniblog 微服务集群一键启动脚本 (Windows)
echo                 当前节点数: 11
echo ===================================================
echo.

echo [1/11] 正在启动 User RPC 服务...
start "User RPC (Port:8080)" cmd /k "go run app/user/rpc/user.go -f app/user/rpc/etc/user.yaml"

echo [2/11] 正在启动 User API 网关...
start "User API (Port:8888)" cmd /k "go run app/user/api/user.go -f app/user/api/etc/user-api.yaml"

echo [3/11] 正在启动 Post RPC 服务 (多级缓存)...
start "Post RPC (Port:8081)" cmd /k "go run app/post/rpc/post.go -f app/post/rpc/etc/post.yaml"

echo [4/11] 正在启动 Post API 网关...
start "Post API (Port:8889)" cmd /k "go run app/post/api/post.go -f app/post/api/etc/post-api.yaml"

echo [5/11] 正在启动 Interaction RPC 服务 (含 MQ 消费者)...
start "Interaction RPC (Port:8082)" cmd /k "go run app/interaction/rpc/interaction.go -f app/interaction/rpc/etc/interaction.yaml"

echo [6/11] 正在启动 Interaction API 网关...
start "Interaction API (Port:8890)" cmd /k "go run app/interaction/api/interaction.go -f app/interaction/api/etc/interaction-api.yaml"

echo [7/11] 正在启动 Notice RPC 服务 (通知内网直连)...
start "Notice RPC (Port:8083)" cmd /k "go run app/notice/rpc/notice.go -f app/notice/rpc/etc/notice.yaml"

echo [8/11] 正在启动 Notice API 网关 (通知前端接口)...
start "Notice API (Port:8891)" cmd /k "go run app/notice/api/notice.go -f app/notice/api/etc/notice-api.yaml"

echo [9/11] 正在启动 OSS API 网关 (对象存储直传服务)...
start "OSS API (Port:8004)" cmd /k "go run app/oss/api/oss.go -f app/oss/api/etc/oss-api.yaml"

:: --- 新增 Relation 关注模块 ---
echo [10/11] 正在启动 Relation RPC 服务 (关系图谱底座)...
start "Relation RPC (Port:8084)" cmd /k "go run app/relation/rpc/relation.go -f app/relation/rpc/etc/relation.yaml"

echo [11/11] 正在启动 Relation API 网关 (关注前端接口)...
start "Relation API (Port:8892)" cmd /k "go run app/relation/api/relation.go -f app/relation/api/etc/relation-api.yaml"

echo.
echo ===================================================
echo 所有 11 个微服务节点已触发启动！
echo 请检查弹出的 11 个控制台窗口，若无报错则集群启动成功。
echo (提示: 确保 Docker 中的 MySQL, Redis, MinIO 已提前启动)
echo ===================================================
pause