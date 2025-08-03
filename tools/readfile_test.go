package tools

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 测试清理辅助函数
func cleanupTestDir(t *testing.T, dir string) {
	t.Helper()
	if err := os.RemoveAll(dir); err != nil {
		t.Logf("Failed to remove test directory %s: %v", dir, err)
	}
}

func restoreWorkingDir(t *testing.T, originalWd string) {
	t.Helper()
	if err := os.Chdir(originalWd); err != nil {
		t.Logf("Failed to restore working directory to %s: %v", originalWd, err)
	}
}

func TestReadFileInput(t *testing.T) {
	t.Run("JSON序列化和反序列化", func(t *testing.T) {
		input := ReadFileInput{
			Path: "test/file.txt",
		}

		// 测试序列化
		jsonData, err := json.Marshal(input)
		require.NoError(t, err)

		// 测试反序列化
		var unmarshaled ReadFileInput
		err = json.Unmarshal(jsonData, &unmarshaled)
		require.NoError(t, err)
		assert.Equal(t, input.Path, unmarshaled.Path)
	})

	t.Run("空路径处理", func(t *testing.T) {
		input := ReadFileInput{
			Path: "",
		}

		jsonData, err := json.Marshal(input)
		require.NoError(t, err)

		var unmarshaled ReadFileInput
		err = json.Unmarshal(jsonData, &unmarshaled)
		require.NoError(t, err)
		assert.Equal(t, "", unmarshaled.Path)
	})
}

func TestReadFile(t *testing.T) {
	// 创建临时目录用于测试
	tempDir, err := os.MkdirTemp("", "readfile_test")
	require.NoError(t, err)
	defer cleanupTestDir(t, tempDir)

	// 保存当前工作目录
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer restoreWorkingDir(t, originalWd)

	// 切换到临时目录
	err = os.Chdir(tempDir)
	require.NoError(t, err)

	t.Run("成功读取文件", func(t *testing.T) {
		// 创建测试文件
		testContent := "Hello, World!\nThis is a test file."
		testFile := "test.txt"
		err := os.WriteFile(testFile, []byte(testContent), 0644)
		require.NoError(t, err)

		// 准备输入参数
		input := ReadFileInput{Path: testFile}
		inputJSON, err := json.Marshal(input)
		require.NoError(t, err)

		// 调用ReadFile函数
		result, err := ReadFile(inputJSON)
		require.NoError(t, err)
		assert.Equal(t, testContent, result)
	})

	t.Run("读取空文件", func(t *testing.T) {
		// 创建空文件
		emptyFile := "empty.txt"
		err := os.WriteFile(emptyFile, []byte(""), 0644)
		require.NoError(t, err)

		// 准备输入参数
		input := ReadFileInput{Path: emptyFile}
		inputJSON, err := json.Marshal(input)
		require.NoError(t, err)

		// 调用ReadFile函数
		result, err := ReadFile(inputJSON)
		require.NoError(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("读取包含特殊字符的文件", func(t *testing.T) {
		// 创建包含特殊字符的文件
		specialContent := "Hello 世界! 🌍\n\t制表符\n\"引号\"\n'单引号'\n\\反斜杠"
		specialFile := "special.txt"
		err := os.WriteFile(specialFile, []byte(specialContent), 0644)
		require.NoError(t, err)

		// 准备输入参数
		input := ReadFileInput{Path: specialFile}
		inputJSON, err := json.Marshal(input)
		require.NoError(t, err)

		// 调用ReadFile函数
		result, err := ReadFile(inputJSON)
		require.NoError(t, err)
		assert.Equal(t, specialContent, result)
	})

	t.Run("读取子目录中的文件", func(t *testing.T) {
		// 创建子目录
		subDir := "subdir"
		err := os.Mkdir(subDir, 0755)
		require.NoError(t, err)

		// 在子目录中创建文件
		subContent := "Content in subdirectory"
		subFile := filepath.Join(subDir, "sub.txt")
		err = os.WriteFile(subFile, []byte(subContent), 0644)
		require.NoError(t, err)

		// 准备输入参数
		input := ReadFileInput{Path: subFile}
		inputJSON, err := json.Marshal(input)
		require.NoError(t, err)

		// 调用ReadFile函数
		result, err := ReadFile(inputJSON)
		require.NoError(t, err)
		assert.Equal(t, subContent, result)
	})

	t.Run("文件不存在", func(t *testing.T) {
		// 准备输入参数（不存在的文件）
		input := ReadFileInput{Path: "nonexistent.txt"}
		inputJSON, err := json.Marshal(input)
		require.NoError(t, err)

		// 调用ReadFile函数
		result, err := ReadFile(inputJSON)
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "failed to read file")
		assert.Contains(t, err.Error(), "nonexistent.txt")
	})

	t.Run("尝试读取目录", func(t *testing.T) {
		// 创建目录
		dirName := "testdir"
		err := os.Mkdir(dirName, 0755)
		require.NoError(t, err)

		// 准备输入参数（目录路径）
		input := ReadFileInput{Path: dirName}
		inputJSON, err := json.Marshal(input)
		require.NoError(t, err)

		// 调用ReadFile函数
		result, err := ReadFile(inputJSON)
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "failed to read file")
	})

	t.Run("空路径", func(t *testing.T) {
		// 准备输入参数（空路径）
		input := ReadFileInput{Path: ""}
		inputJSON, err := json.Marshal(input)
		require.NoError(t, err)

		// 调用ReadFile函数
		result, err := ReadFile(inputJSON)
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("无效的JSON输入", func(t *testing.T) {
		// 无效的JSON
		invalidJSON := json.RawMessage(`{"path": 123}`) // path应该是字符串，不是数字

		// 调用ReadFile函数
		result, err := ReadFile(invalidJSON)
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "failed to parse input")
	})

	t.Run("格式错误的JSON", func(t *testing.T) {
		// 格式错误的JSON
		malformedJSON := json.RawMessage(`{"path": "test.txt"`) // 缺少结束括号

		// 调用ReadFile函数
		result, err := ReadFile(malformedJSON)
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "failed to parse input")
	})
}

func TestReadFileDefinition(t *testing.T) {
	t.Run("验证工具定义结构", func(t *testing.T) {
		def := ReadFileDefinition

		// 验证基本字段
		assert.Equal(t, "read_file", def.Name)
		assert.NotEmpty(t, def.Description)
		assert.Contains(t, def.Description, "Read the contents")
		assert.Contains(t, def.Description, "relative file path")
		assert.NotNil(t, def.InputSchema)
		assert.NotNil(t, def.Function)
	})

	t.Run("验证输入模式", func(t *testing.T) {
		def := ReadFileDefinition
		schema := def.InputSchema

		// 验证模式不为空
		assert.NotNil(t, schema.Properties)
		assert.NotEmpty(t, schema.Required)
	})

	t.Run("验证函数可调用性", func(t *testing.T) {
		// 创建临时文件用于测试
		tempDir, err := os.MkdirTemp("", "readfile_def_test")
		require.NoError(t, err)
		defer cleanupTestDir(t, tempDir)

		// 保存当前工作目录
		originalWd, err := os.Getwd()
		require.NoError(t, err)
		defer restoreWorkingDir(t, originalWd)

		// 切换到临时目录
		err = os.Chdir(tempDir)
		require.NoError(t, err)

		// 创建测试文件
		testContent := "Function test content"
		testFile := "func_test.txt"
		err = os.WriteFile(testFile, []byte(testContent), 0644)
		require.NoError(t, err)

		// 准备输入
		input := ReadFileInput{Path: testFile}
		inputJSON, err := json.Marshal(input)
		require.NoError(t, err)

		// 通过定义调用函数
		result, err := ReadFileDefinition.Function(inputJSON)
		require.NoError(t, err)
		assert.Equal(t, testContent, result)
	})
}

func TestReadFileBehaviorEdgeCases(t *testing.T) {
	// 保存当前工作目录
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer restoreWorkingDir(t, originalWd)

	t.Run("读取大文件", func(t *testing.T) {
		// 创建临时目录
		tempDir, err := os.MkdirTemp("", "readfile_large_test")
		require.NoError(t, err)
		defer cleanupTestDir(t, tempDir)

		// 切换到临时目录
		err = os.Chdir(tempDir)
		require.NoError(t, err)

		// 创建一个相对较大的文件（1MB）
		largeContent := make([]byte, 1024*1024)
		for i := range largeContent {
			largeContent[i] = byte('A' + (i % 26))
		}
		largeFile := "large.txt"
		err = os.WriteFile(largeFile, largeContent, 0644)
		require.NoError(t, err)

		// 准备输入参数
		input := ReadFileInput{Path: largeFile}
		inputJSON, err := json.Marshal(input)
		require.NoError(t, err)

		// 调用ReadFile函数
		result, err := ReadFile(inputJSON)
		require.NoError(t, err)
		assert.Equal(t, string(largeContent), result)
		assert.Len(t, result, 1024*1024)
	})

	t.Run("包含换行符的文件", func(t *testing.T) {
		// 创建临时目录
		tempDir, err := os.MkdirTemp("", "readfile_newline_test")
		require.NoError(t, err)
		defer cleanupTestDir(t, tempDir)

		// 切换到临时目录
		err = os.Chdir(tempDir)
		require.NoError(t, err)

		// 创建包含各种换行符的文件
		newlineContent := "Line 1\nLine 2\r\nLine 3\rLine 4\n\n\nLine 7"
		newlineFile := "newlines.txt"
		err = os.WriteFile(newlineFile, []byte(newlineContent), 0644)
		require.NoError(t, err)

		// 准备输入参数
		input := ReadFileInput{Path: newlineFile}
		inputJSON, err := json.Marshal(input)
		require.NoError(t, err)

		// 调用ReadFile函数
		result, err := ReadFile(inputJSON)
		require.NoError(t, err)
		assert.Equal(t, newlineContent, result)
	})
}
