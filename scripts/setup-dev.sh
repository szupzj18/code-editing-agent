#!/bin/bash

# 开发环境设置脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 设置Go项目开发环境...${NC}"

# 检查Go是否安装
echo -e "\n${BLUE}📋 检查环境依赖...${NC}"
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go未安装，请先安装Go 1.21或更高版本${NC}"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo -e "✅ Go版本: $GO_VERSION"

# 检查Git是否安装
if ! command -v git &> /dev/null; then
    echo -e "${RED}❌ Git未安装${NC}"
    exit 1
fi
echo -e "✅ Git已安装"

# 下载依赖
echo -e "\n${BLUE}📦 下载项目依赖...${NC}"
go mod download
go mod verify
echo -e "${GREEN}✅ 依赖下载完成${NC}"

# 安装开发工具
echo -e "\n${BLUE}🔧 安装开发工具...${NC}"

# 安装golangci-lint
if ! command -v golangci-lint &> /dev/null; then
    echo "📥 安装golangci-lint..."
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
    echo -e "${GREEN}✅ golangci-lint安装完成${NC}"
else
    echo -e "✅ golangci-lint已安装"
fi

# 安装gosec
if ! command -v gosec &> /dev/null; then
    echo "📥 安装gosec..."
    go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
    echo -e "${GREEN}✅ gosec安装完成${NC}"
else
    echo -e "✅ gosec已安装"
fi

# 安装gocov-html（用于更好的覆盖率报告）
if ! command -v gocov-html &> /dev/null; then
    echo "📥 安装gocov-html..."
    go install github.com/matm/gocov-html@latest
    echo -e "${GREEN}✅ gocov-html安装完成${NC}"
else
    echo -e "✅ gocov-html已安装"
fi

# 安装gocov
if ! command -v gocov &> /dev/null; then
    echo "📥 安装gocov..."
    go install github.com/axw/gocov/gocov@latest
    echo -e "${GREEN}✅ gocov安装完成${NC}"
else
    echo -e "✅ gocov已安装"
fi

# 设置Git钩子
echo -e "\n${BLUE}🪝 设置Git预提交钩子...${NC}"
if [ -f ".git/hooks/pre-commit" ]; then
    echo -e "${GREEN}✅ Git pre-commit钩子已存在${NC}"
else
    echo -e "${YELLOW}⚠️ Git pre-commit钩子未找到${NC}"
fi

# 创建必要的目录
echo -e "\n${BLUE}📁 创建项目目录...${NC}"
mkdir -p coverage
mkdir -p bin
mkdir -p tmp
echo -e "${GREEN}✅ 目录创建完成${NC}"

# 运行一次测试确保环境正常
echo -e "\n${BLUE}🧪 验证环境设置...${NC}"
echo "运行测试..."
if go test -v ./...; then
    echo -e "${GREEN}✅ 测试通过，环境设置成功${NC}"
else
    echo -e "${RED}❌ 测试失败，请检查代码${NC}"
    exit 1
fi

# 生成初始覆盖率报告
echo -e "\n${BLUE}📊 生成初始覆盖率报告...${NC}"
make test-coverage
echo -e "${GREEN}✅ 覆盖率报告生成完成${NC}"

# 显示使用说明
echo -e "\n${GREEN}🎉 开发环境设置完成！${NC}"
echo -e "\n${BLUE}📖 常用命令:${NC}"
echo -e "  make help          - 查看所有可用命令"
echo -e "  make test          - 运行测试"
echo -e "  make test-coverage - 运行测试并生成覆盖率"
echo -e "  make lint          - 运行代码检查"
echo -e "  make ci            - 运行所有CI检查"
echo -e "  make dev           - 开发模式运行"
echo -e "  make build         - 构建项目"

echo -e "\n${BLUE}🔧 开发工具:${NC}"
echo -e "  ./scripts/coverage.sh - 详细覆盖率分析"
echo -e "  golangci-lint run     - 代码质量检查"
echo -e "  gosec ./...           - 安全扫描"

echo -e "\n${BLUE}📁 重要文件和目录:${NC}"
echo -e "  coverage/             - 覆盖率报告"
echo -e "  .github/workflows/    - CI/CD配置"
echo -e "  .golangci.yml         - 代码检查配置"
echo -e "  Makefile              - 构建任务"

echo -e "\n${GREEN}✨ 环境设置完成，开始开发吧！${NC}"
