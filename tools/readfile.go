package tools

import (
	"encoding/json"
	"fmt"
	"os"
)

// ReadFileInput 定义读取文件工具的输入参数
type ReadFileInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of a file in the working directory."`
}

// ReadFile 实现文件读取功能
func ReadFile(input json.RawMessage) (string, error) {
	var params ReadFileInput
	err := json.Unmarshal(input, &params)
	if err != nil {
		return "", fmt.Errorf("failed to parse input: %w", err)
	}

	content, err := os.ReadFile(params.Path)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", params.Path, err)
	}

	return string(content), nil
}

// ReadFileDefinition 文件读取工具的完整定义
var ReadFileDefinition = ToolDefinition{
	Name:        "read_file",
	Description: "Read the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names.",
	InputSchema: GenerateSchema[ReadFileInput](),
	Function:    ReadFile,
}
