#!/bin/bash
set -e

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查必要的环境变量
if [ -z "$IMAGE_NAME" ]; then
    log_error "IMAGE_NAME 环境变量未设置"
    exit 1
fi

if [ -z "$IMAGE_TAG" ]; then
    log_error "IMAGE_TAG 环境变量未设置"
    exit 1
fi

# 拉取指定版本的镜像
log_info "拉取镜像: ${IMAGE_NAME}:${IMAGE_TAG}"
docker pull ${IMAGE_NAME}:${IMAGE_TAG} || {
    log_error "拉取镜像失败"
    exit 1
}

# 获取 docker-compose 中所有服务
SERVICES=$(docker-compose config --services)
TOTAL_SERVICES=$(echo "$SERVICES" | wc -l)

log_info "发现 $TOTAL_SERVICES 个服务需要重启"

# 逐个重启服务
for SERVICE in $SERVICES; do
    log_info "处理服务: $SERVICE"

    # 停止当前服务
    log_info "停止服务: $SERVICE"
    docker-compose stop $SERVICE || {
        log_warn "停止服务 $SERVICE 失败，尝试继续"
    }

    # 移除旧容器
    log_info "移除旧容器: $SERVICE"
    docker-compose rm -f $SERVICE || {
        log_warn "移除容器 $SERVICE 失败，尝试继续"
    }

    # 启动新服务
    log_info "启动服务: $SERVICE"
    docker-compose up -d $SERVICE || {
        log_error "启动服务 $SERVICE 失败"
        exit 1
    }

    # 等待服务启动并稳定
    CONTAINER_ID=$(docker-compose ps -q $SERVICE)

    if [ -z "$CONTAINER_ID" ]; then
        log_error "无法获取 $SERVICE 的容器ID"
        exit 1
    fi

    log_info "等待服务 $SERVICE 稳定 (容器ID: $CONTAINER_ID)"

    # 等待容器准备就绪
    RETRIES=30
    for i in $(seq 1 $RETRIES); do
        if docker inspect --format='{{.State.Running}}' $CONTAINER_ID 2>/dev/null | grep -q "true"; then
            log_info "服务 $SERVICE 已启动并运行中"
            # 给服务一些额外时间来完全初始化
            sleep 5
            break
        fi

        if [ $i -eq $RETRIES ]; then
            log_error "服务 $SERVICE 未能在预期时间内启动"
            exit 1
        fi

        log_info "等待服务 $SERVICE 启动 (尝试 $i/$RETRIES)"
        sleep 2
    done
done

log_info "所有服务已成功滚动更新"
exit 0
