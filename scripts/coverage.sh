#!/bin/bash

# 测试覆盖率分析脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
COVERAGE_FILE="coverage.out"
COVERAGE_HTML="coverage.html"
COVERAGE_XML="coverage.xml"
MIN_COVERAGE=80

echo -e "${BLUE}🔍 开始测试覆盖率分析...${NC}"

# 创建覆盖率目录
mkdir -p coverage

# 运行测试并生成覆盖率
echo -e "${BLUE}📊 运行测试并生成覆盖率报告...${NC}"
go test -v -race -coverprofile=coverage/$COVERAGE_FILE -covermode=atomic ./...

if [ ! -f "coverage/$COVERAGE_FILE" ]; then
    echo -e "${RED}❌ 覆盖率文件生成失败${NC}"
    exit 1
fi

# 生成HTML报告
echo -e "${BLUE}📄 生成HTML覆盖率报告...${NC}"
go tool cover -html=coverage/$COVERAGE_FILE -o coverage/$COVERAGE_HTML

# 计算总覆盖率
TOTAL_COVERAGE=$(go tool cover -func=coverage/$COVERAGE_FILE | grep total | awk '{print $3}' | sed 's/%//')

echo -e "${BLUE}📈 覆盖率统计:${NC}"
go tool cover -func=coverage/$COVERAGE_FILE

# 按包显示覆盖率
echo -e "\n${BLUE}📦 按包显示覆盖率:${NC}"
go tool cover -func=coverage/$COVERAGE_FILE | grep -v total | awk '{print $1}' | sed 's/.*\///' | sort | uniq | while read pkg; do
    if [ ! -z "$pkg" ]; then
        FUNC_COUNT=$(go tool cover -func=coverage/$COVERAGE_FILE | grep "/$pkg/" | wc -l)
        if [ $FUNC_COUNT -gt 0 ]; then
            PKG_COVERAGE=$(go tool cover -func=coverage/$COVERAGE_FILE | grep "/$pkg/" | awk '{sum+=$3; count++} END {if(count>0) print sum/count; else print 0}' | sed 's/%//')
            printf "  %-20s: %.1f%%\n" "$pkg" "$PKG_COVERAGE"
        fi
    fi
done

# 检查覆盖率阈值
echo -e "\n${BLUE}🎯 覆盖率检查:${NC}"
if (( $(echo "$TOTAL_COVERAGE >= $MIN_COVERAGE" | bc -l) )); then
    echo -e "${GREEN}✅ 覆盖率 $TOTAL_COVERAGE% 达到要求 (>= $MIN_COVERAGE%)${NC}"
    EXIT_CODE=0
else
    echo -e "${RED}❌ 覆盖率 $TOTAL_COVERAGE% 未达到要求 (>= $MIN_COVERAGE%)${NC}"
    EXIT_CODE=1
fi

# 生成覆盖率徽章
echo -e "\n${BLUE}🏷️  生成覆盖率徽章...${NC}"
if (( $(echo "$TOTAL_COVERAGE >= 80" | bc -l) )); then
    COLOR="brightgreen"
elif (( $(echo "$TOTAL_COVERAGE >= 60" | bc -l) )); then
    COLOR="yellow"
else
    COLOR="red"
fi

BADGE_URL="https://img.shields.io/badge/coverage-${TOTAL_COVERAGE}%25-${COLOR}"
curl -s "$BADGE_URL" > coverage/coverage-badge.svg

# 查找未覆盖的代码
echo -e "\n${BLUE}🔍 未覆盖的代码:${NC}"
UNCOVERED_LINES=$(go tool cover -func=coverage/$COVERAGE_FILE | grep -v "100.0%" | grep -v "total:" | head -10)
if [ ! -z "$UNCOVERED_LINES" ]; then
    echo "$UNCOVERED_LINES"
else
    echo -e "${GREEN}🎉 所有代码都已覆盖!${NC}"
fi

# 生成覆盖率趋势数据
echo -e "\n${BLUE}📊 保存覆盖率历史数据...${NC}"
TIMESTAMP=$(date "+%Y-%m-%d %H:%M:%S")
echo "$TIMESTAMP,$TOTAL_COVERAGE" >> coverage/coverage-history.csv

# 输出文件位置
echo -e "\n${BLUE}📂 覆盖率报告文件:${NC}"
echo "  - HTML报告: coverage/$COVERAGE_HTML"
echo "  - 原始数据: coverage/$COVERAGE_FILE"
echo "  - 覆盖率徽章: coverage/coverage-badge.svg"
echo "  - 历史数据: coverage/coverage-history.csv"

# 如果在交互式环境中，询问是否打开HTML报告
if [ -t 1 ] && command -v open >/dev/null 2>&1; then
    echo -e "\n${YELLOW}是否打开HTML覆盖率报告? (y/N)${NC}"
    read -r response
    if [[ "$response" =~ ^[Yy]$ ]]; then
        open coverage/$COVERAGE_HTML
    fi
fi

echo -e "\n${GREEN}✨ 覆盖率分析完成!${NC}"

exit $EXIT_CODE
