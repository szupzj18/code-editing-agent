#!/bin/bash

# Code Editing Agent 启动脚本
# 这个脚本会检查 API keys 并启动适当的 AI 提供商

echo "🤖 Code Editing Agent"
echo "===================="

# 检查 API Keys
if [ -n "$OPENAI_API_KEY" ]; then
    echo "✅ 发现 OpenAI API Key - 将使用 GPT-4o"
elif [ -n "$ANTHROPIC_API_KEY" ]; then
    echo "✅ 发现 Anthropic API Key - 将使用 Claude"
else
    echo "❌ 错误: 未找到 API Key"
    echo ""
    echo "请设置以下环境变量之一："
    echo "  export OPENAI_API_KEY='your-openai-api-key'"
    echo "  export ANTHROPIC_API_KEY='your-anthropic-api-key'"
    echo ""
    echo "获取 API Key:"
    echo "  OpenAI: https://platform.openai.com/api-keys"
    echo "  Anthropic: https://console.anthropic.com/"
    exit 1
fi

echo ""
echo "🚀 正在启动 Agent..."
echo "💡 提示: 你可以要求我读取文件、解释代码或进行其他操作"
echo "⛔ 使用 Ctrl+C 退出"
echo ""

# 构建并运行
go build -o agent && ./agent
