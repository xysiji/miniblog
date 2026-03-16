@echo off
chcp 65001 >nul
echo ===================================================
echo      Miniblog 微服务集群一键启动脚本 (Windows)
echo ===================================================
echo.

echo [1/6] 正在启动 User RPC 服务...
start "User RPC (Port:8080)" cmd /k "go run app/user/rpc/user.go -f app/user/rpc/etc/user.yaml"

echo [2/6] 正在启动 User API 网关...
start "User API (Port:8888)" cmd /k "go run app/user/api/user.go -f app/user/api/etc/user-api.yaml"

echo [3/6] 正在启动 Post RPC 服务 (多级缓存)...
start "Post RPC (Port:8081)" cmd /k "go run app/post/rpc/post.go -f app/post/rpc/etc/post.yaml"

echo [4/6] 正在启动 Post API 网关...
start "Post API (Port:8889)" cmd /k "go run app/post/api/post.go -f app/post/api/etc/post-api.yaml"

echo [5/6] 正在启动 Interaction RPC 服务 (含 MQ 消费者)...
start "Interaction RPC (Port:8082)" cmd /k "go run app/interaction/rpc/interaction.go -f app/interaction/rpc/etc/interaction.yaml"

echo [6/6] 正在启动 Interaction API 网关...
start "Interaction API (Port:8890)" cmd /k "go run app/interaction/api/interaction.go -f app/interaction/api/etc/interaction-api.yaml"

echo.
echo 所有微服务已触发启动！
echo 请检查弹出的 6 个控制台窗口，若无报错则启动成功。
echo (要关闭服务，直接关闭对应的黑窗口即可)
echo ===================================================
pause