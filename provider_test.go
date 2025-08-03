package main

import (
	"context"
	"os"
	"testing"

	"agent/tools"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnthropicProvider(t *testing.T) {
	t.Skip("需要 ANTHROPIC_API_KEY 环境变量")

	provider := NewAnthropicProvider()
	require.NotNil(t, provider)

	// 测试基本对话
	conversation := []Message{
		{Role: "user", Content: "Hello"},
	}

	tools := []tools.ToolDefinition{tools.ReadFileDefinition}

	response, err := provider.RunInference(context.Background(), conversation, tools)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	t.Logf("Anthropic response: %+v", response)
}

func TestOpenAIProvider(t *testing.T) {
	t.Skip("需要 OPENAI_API_KEY 环境变量")

	provider := NewOpenAIProvider("test-key")
	require.NotNil(t, provider)

	// 测试基本对话
	conversation := []Message{
		{Role: "user", Content: "Hello"},
	}

	tools := []tools.ToolDefinition{tools.ReadFileDefinition}

	response, err := provider.RunInference(context.Background(), conversation, tools)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	t.Logf("OpenAI response: %+v", response)
}

func TestReadFileTool(t *testing.T) {
	// 创建一个测试文件
	testContent := "Hello, World!"
	testFile := "/tmp/test_file.txt"

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)
	defer func() {
		if err := os.Remove(testFile); err != nil {
			t.Logf("Failed to remove test file: %v", err)
		}
	}()

	// 测试 ReadFile 工具
	input := `{"path": "/tmp/test_file.txt"}`
	result, err := tools.ReadFile([]byte(input))

	assert.NoError(t, err)
	assert.Equal(t, testContent, result)
}

func TestReadFileToolError(t *testing.T) {
	// 测试读取不存在的文件
	input := `{"path": "/nonexistent/file.txt"}`
	result, err := tools.ReadFile([]byte(input))

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "failed to read file")
}

func TestUnifiedMessageFormat(t *testing.T) {
	msg := Message{
		Role:    "user",
		Content: "Test message",
	}

	assert.Equal(t, "user", msg.Role)
	assert.Equal(t, "Test message", msg.Content)
}

func TestToolCallFormat(t *testing.T) {
	toolCall := ToolCall{
		ID:    "call_123",
		Name:  "test_tool",
		Input: []byte(`{"param": "value"}`),
	}

	assert.Equal(t, "call_123", toolCall.ID)
	assert.Equal(t, "test_tool", toolCall.Name)
	assert.NotNil(t, toolCall.Input)
}
