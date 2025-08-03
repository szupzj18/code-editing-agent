package main

import (
	"os"
	"testing"

	"agent/tools"

	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

// 测试用的结构体
type TestStruct struct {
	Name     string `json:"name" jsonschema_description:"用户姓名"`
	Age      int    `json:"age" jsonschema_description:"用户年龄"`
	Email    string `json:"email" jsonschema_description:"用户邮箱"`
	Optional *bool  `json:"optional,omitempty" jsonschema_description:"可选字段"`
}

type EmptyStruct struct{}

type NestedStruct struct {
	Basic  TestStruct `json:"basic" jsonschema_description:"基本信息"`
	Status string     `json:"status" jsonschema_description:"状态"`
}

func TestGenerateSchema(t *testing.T) {
	t.Run("基本功能测试", func(t *testing.T) {
		schema := tools.GenerateSchema[tools.ReadFileInput]()

		assert.NotNil(t, schema)
		assert.NotNil(t, schema.Properties)

		// 验证 Properties 是正确的类型
		if properties, ok := schema.Properties.(*orderedmap.OrderedMap[string, *jsonschema.Schema]); ok {
			// 验证 Properties 包含 path 字段
			pathProperty, exists := properties.Get("path")
			assert.True(t, exists, "应该包含 path 属性")
			assert.NotNil(t, pathProperty)
		} else {
			t.Errorf("Properties 的实际类型是: %T", schema.Properties)
		}
	})

	t.Run("复杂结构体测试", func(t *testing.T) {
		schema := tools.GenerateSchema[TestStruct]()

		require.NotNil(t, schema)
		require.NotNil(t, schema.Properties)

		if properties, ok := schema.Properties.(*orderedmap.OrderedMap[string, *jsonschema.Schema]); ok {
			// 验证所有预期的字段都存在
			expectedFields := []string{"name", "age", "email", "optional"}
			for _, field := range expectedFields {
				_, exists := properties.Get(field)
				assert.True(t, exists, "应该包含字段: %s", field)
			}

			// 验证字段数量
			assert.Equal(t, 4, properties.Len(), "应该有4个属性")
		} else {
			t.Errorf("Properties 的实际类型是: %T", schema.Properties)
		}
	})

	t.Run("空结构体测试", func(t *testing.T) {
		schema := tools.GenerateSchema[EmptyStruct]()

		assert.NotNil(t, schema)
		// 空结构体应该有空的 Properties 或者长度为0
		if properties, ok := schema.Properties.(*orderedmap.OrderedMap[string, *jsonschema.Schema]); ok {
			assert.Equal(t, 0, properties.Len(), "空结构体应该没有属性")
		}
	})

	t.Run("嵌套结构体测试", func(t *testing.T) {
		schema := tools.GenerateSchema[NestedStruct]()

		require.NotNil(t, schema)
		require.NotNil(t, schema.Properties)

		if properties, ok := schema.Properties.(*orderedmap.OrderedMap[string, *jsonschema.Schema]); ok {
			// 验证顶层字段
			basicProperty, exists := properties.Get("basic")
			assert.True(t, exists, "应该包含 basic 字段")
			assert.NotNil(t, basicProperty)

			statusProperty, exists := properties.Get("status")
			assert.True(t, exists, "应该包含 status 字段")
			assert.NotNil(t, statusProperty)
		} else {
			t.Errorf("Properties 的实际类型是: %T", schema.Properties)
		}
	})

	t.Run("原有测试保持兼容", func(t *testing.T) {
		schema := tools.GenerateSchema[tools.ReadFileInput]()
		assert.NotNil(t, schema)
		t.Log("Generated schema:", schema)
	})
}

func TestGenerateSchemaProperties(t *testing.T) {
	t.Run("验证ReadFileInput的具体属性", func(t *testing.T) {
		schema := tools.GenerateSchema[tools.ReadFileInput]()

		require.NotNil(t, schema.Properties)

		if properties, ok := schema.Properties.(*orderedmap.OrderedMap[string, *jsonschema.Schema]); ok {
			// 验证 path 字段的具体属性
			pathProperty, exists := properties.Get("path")
			require.True(t, exists)
			require.NotNil(t, pathProperty)

			// 验证 path 字段是字符串类型
			assert.NotNil(t, pathProperty.Type)
			t.Logf("Path property type: %v", pathProperty.Type)
			t.Logf("Path property description: %s", pathProperty.Description)
		} else {
			t.Errorf("Properties 的实际类型是: %T", schema.Properties)
		}
	})
}

func Test_ReadFileTool(t *testing.T) {
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

// 基准测试
func BenchmarkGenerateSchema(b *testing.B) {
	b.Run("ReadFileInput", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = tools.GenerateSchema[tools.ReadFileInput]()
		}
	})

	b.Run("TestStruct", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = tools.GenerateSchema[TestStruct]()
		}
	})
}
