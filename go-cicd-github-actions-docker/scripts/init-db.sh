#!/bin/bash

# 输出颜色配置
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# MySQL容器名称
MYSQL_CONTAINER="mysql"
SQL_FILE="$(dirname "$0")/../configs/db_init.sql"

echo -e "${YELLOW}正在初始化数据库...${NC}"

# 检查SQL文件是否存在
if [ ! -f "$SQL_FILE" ]; then
    echo -e "${RED}错误: SQL文件 '$SQL_FILE' 不存在!${NC}"
    exit 1
fi

# 检查Docker容器是否在运行
if ! docker ps | grep -q "$MYSQL_CONTAINER"; then
    echo -e "${RED}错误: MySQL容器 '$MYSQL_CONTAINER' 未运行!${NC}"
    echo -e "${YELLOW}请先运行: docker-compose -f docker-compose.dev-env.yml up -d${NC}"
    exit 1
fi

# 将SQL文件复制到容器内
echo -e "${YELLOW}将SQL文件复制到容器...${NC}"
docker cp "$SQL_FILE" "$MYSQL_CONTAINER":/tmp/init.sql

# 在容器内执行SQL文件
echo -e "${YELLOW}在容器内执行SQL文件...${NC}"
docker exec -i "$MYSQL_CONTAINER" bash -c "mysql -u root -p123456 < /tmp/init.sql"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}数据库初始化成功!${NC}"
else
    echo -e "${RED}数据库初始化失败!${NC}"
    exit 1
fi

echo -e "${YELLOW}已创建以下用户:${NC}"
echo -e "  用户名: admin, 邮箱: admin@example.com, 密码: password123"
echo -e "  用户名: test_user, 邮箱: user@example.com, 密码: password123"
