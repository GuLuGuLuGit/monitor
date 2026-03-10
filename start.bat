@echo off
REM OpenClaw 监控系统启动脚本 (Windows)

echo 正在启动 OpenClaw 监控系统...

REM 检查配置文件
if not exist "config.yaml" (
    echo 错误: 找不到 config.yaml 配置文件
    exit /b 1
)

REM 创建数据目录
if not exist "data" mkdir data

REM 进入后端目录
cd backend

REM 检查 Go 是否安装
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo 错误: 未安装 Go
    exit /b 1
)

REM 下载依赖
echo 正在下载依赖...
go mod download

REM 运行服务
echo 正在启动服务...
go run cmd\server\main.go -config ..\config.yaml
