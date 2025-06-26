#!/bin/bash

# 输出颜色配置
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# MySQL容器名称
MYSQL_CONTAINER="mysql"

echo -e "${YELLOW}检查数据库初始化状态...${NC}"

# 检查Docker容器是否在运行
if ! docker ps | grep -q "$MYSQL_CONTAINER"; then
    echo -e "${RED}错误: MySQL容器 '$MYSQL_CONTAINER' 未运行!${NC}"
    echo -e "${YELLOW}请先运行: docker-compose -f docker-compose.dev-env.yml up -d${NC}"
    exit 1
fi

# 检查数据库表是否存在
TABLE_COUNT=$(docker exec -i "$MYSQL_CONTAINER" bash -c "mysql -u root -p123456 -e 'USE practice; SHOW TABLES;' | wc -l")

if [ "$TABLE_COUNT" -gt 1 ]; then
    echo -e "${GREEN}数据库已正确初始化!${NC}"

    # 检查用户表中的记录数
    USER_COUNT=$(docker exec -i "$MYSQL_CONTAINER" bash -c "mysql -u root -p123456 -e 'SELECT COUNT(*) FROM practice.user;' | tail -1")
    echo -e "${GREEN}用户表中有 $USER_COUNT 条记录${NC}"

    # 显示用户信息
    echo -e "${YELLOW}用户信息:${NC}"
    docker exec -i "$MYSQL_CONTAINER" bash -c "mysql -u root -p123456 -e 'SELECT id, username, email, first_name, last_name FROM practice.user;'"

    # 显示角色信息
    echo -e "${YELLOW}角色信息:${NC}"
    docker exec -i "$MYSQL_CONTAINER" bash -c "mysql -u root -p123456 -e 'SELECT * FROM practice.role;'"
else
    echo -e "${RED}数据库未正确初始化!${NC}"
    echo -e "${YELLOW}请运行: make init-db${NC}"
    exit 1
fi
