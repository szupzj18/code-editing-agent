package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"agent/tools"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/shared"
)

// AI Provider interface for unified handling
type AIProvider interface {
	RunInference(ctx context.Context, conversation []Message, tools []tools.ToolDefinition) (*Response, error)
}

// Unified message structure
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Unified response structure
type Response struct {
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type ToolCall struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Input json.RawMessage `json:"input"`
}

// Anthropic provider implementation
type AnthropicProvider struct {
	client anthropic.Client
}

func NewAnthropicProvider() *AnthropicProvider {
	return &AnthropicProvider{
		client: anthropic.NewClient(),
	}
}

func (ap *AnthropicProvider) RunInference(ctx context.Context, conversation []Message, tools []tools.ToolDefinition) (*Response, error) {
	// Convert unified messages to Anthropic format
	anthropicMessages := make([]anthropic.MessageParam, len(conversation))
	for i, msg := range conversation {
		if msg.Role == "user" {
			anthropicMessages[i] = anthropic.NewUserMessage(anthropic.NewTextBlock(msg.Content))
		} else {
			anthropicMessages[i] = anthropic.NewAssistantMessage(anthropic.NewTextBlock(msg.Content))
		}
	}

	// Convert tools to Anthropic format
	anthropicTools := []anthropic.ToolUnionParam{}
	for _, tool := range tools {
		anthropicTools = append(anthropicTools, anthropic.ToolUnionParam{
			OfTool: &anthropic.ToolParam{
				InputSchema: tool.InputSchema,
				Name:        tool.Name,
				Description: anthropic.String(tool.Description),
			},
		})
	}

	message, err := ap.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_7SonnetLatest,
		MaxTokens: int64(1024),
		Messages:  anthropicMessages,
		Tools:     anthropicTools,
	})
	if err != nil {
		return nil, err
	}

	// Convert response back to unified format
	response := &Response{}
	for _, content := range message.Content {
		switch content.Type {
		case "text":
			response.Content = content.Text
		case "tool_use":
			toolCall := ToolCall{
				ID:    content.ID,
				Name:  content.Name,
				Input: content.Input,
			}
			response.ToolCalls = append(response.ToolCalls, toolCall)
		}
	}

	return response, nil
}

// OpenAI provider implementation
type OpenAIProvider struct {
	client openai.Client
}

func NewOpenAIProvider(apiKey string) *OpenAIProvider {
	return &OpenAIProvider{
		client: openai.NewClient(option.WithAPIKey(apiKey)),
	}
}

func (op *OpenAIProvider) RunInference(ctx context.Context, conversation []Message, tools []tools.ToolDefinition) (*Response, error) {
	// Convert unified messages to OpenAI format
	openaiMessages := make([]openai.ChatCompletionMessageParamUnion, len(conversation))
	for i, msg := range conversation {
		if msg.Role == "user" {
			openaiMessages[i] = openai.UserMessage(msg.Content)
		} else {
			openaiMessages[i] = openai.AssistantMessage(msg.Content)
		}
	}

	// Convert tools to OpenAI format
	openaiTools := []openai.ChatCompletionToolParam{}
	for _, tool := range tools {
		// Convert Anthropic schema to OpenAI schema format
		// Create parameters from schema properties
		params := make(map[string]interface{})
		if tool.InputSchema.Properties != nil {
			if props, ok := tool.InputSchema.Properties.(map[string]interface{}); ok {
				params["type"] = "object"
				params["properties"] = props
				if len(tool.InputSchema.Required) > 0 {
					params["required"] = tool.InputSchema.Required
				}
			}
		}

		openaiTool := openai.ChatCompletionToolParam{
			Function: shared.FunctionDefinitionParam{
				Name:        tool.Name,
				Description: param.NewOpt(tool.Description),
				Parameters:  shared.FunctionParameters(params),
			},
		}
		openaiTools = append(openaiTools, openaiTool)
	}

	params := openai.ChatCompletionNewParams{
		Model:    openai.ChatModelGPT4o,
		Messages: openaiMessages,
	}

	if len(openaiTools) > 0 {
		params.Tools = openaiTools
	}

	completion, err := op.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, err
	}

	// Convert response back to unified format
	response := &Response{}
	if len(completion.Choices) > 0 {
		choice := completion.Choices[0]
		if choice.Message.Content != "" {
			response.Content = choice.Message.Content
		}

		// Handle tool calls
		for _, toolCall := range choice.Message.ToolCalls {
			unifiedToolCall := ToolCall{
				ID:    toolCall.ID,
				Name:  toolCall.Function.Name,
				Input: json.RawMessage(toolCall.Function.Arguments),
			}
			response.ToolCalls = append(response.ToolCalls, unifiedToolCall)
		}
	}

	return response, nil
}

func main() {
	var provider AIProvider

	// 优先使用 OpenAI，如果没有 API key 则使用 Anthropic
	if openaiKey := os.Getenv("OPENAI_API_KEY"); openaiKey != "" {
		provider = NewOpenAIProvider(openaiKey)
		fmt.Println("使用 OpenAI GPT-4o")
	} else {
		provider = NewAnthropicProvider()
		fmt.Println("使用 Anthropic Claude")
	}

	tools := []tools.ToolDefinition{tools.ReadFileDefinition}
	scanner := bufio.NewScanner(os.Stdin)
	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		return scanner.Text(), true
	}

	agent := NewAgent(provider, getUserMessage, tools)
	err := agent.Run(context.TODO())
	if err != nil {
		fmt.Printf("Error: %s\n\n", err)
	}
}

func NewAgent(provider AIProvider, getUserMessage func() (string, bool), tools []tools.ToolDefinition) *Agent {
	return &Agent{
		provider:       provider,
		getUserMessage: getUserMessage,
		tools:          tools,
	}
}

type Agent struct {
	provider       AIProvider
	getUserMessage func() (string, bool)
	tools          []tools.ToolDefinition
}

func (a Agent) Run(ctx context.Context) error {
	conversation := []Message{}

	fmt.Println("Chat with Claude/GPT (use 'ctrl-c' to quit)")
	for {
		fmt.Print("\u001b[94mYou\u001b[0m: ")
		userInput, ok := a.getUserMessage()
		if !ok {
			break
		}

		userMessage := Message{
			Role:    "user",
			Content: userInput,
		}
		conversation = append(conversation, userMessage)

		response, err := a.provider.RunInference(ctx, conversation, a.tools)
		if err != nil {
			return err
		}

		// Handle tool calls first
		if len(response.ToolCalls) > 0 {
			for _, toolCall := range response.ToolCalls {
				// Find and execute the tool
				for _, tool := range a.tools {
					if tool.Name == toolCall.Name {
						result, err := tool.Function(toolCall.Input)
						if err != nil {
							fmt.Printf("\u001b[91mTool Error\u001b[0m: %s\n", err)
						} else {
							fmt.Printf("\u001b[92mTool Result\u001b[0m: %s\n", result)
						}

						// Add tool result to conversation
						toolResultMessage := Message{
							Role:    "assistant",
							Content: fmt.Sprintf("Tool %s executed with result: %s", toolCall.Name, result),
						}
						conversation = append(conversation, toolResultMessage)
						break
					}
				}
			}
		}

		// Display assistant response
		if response.Content != "" {
			fmt.Printf("\u001b[93mAssistant\u001b[0m: %s\n", response.Content)
			assistantMessage := Message{
				Role:    "assistant",
				Content: response.Content,
			}
			conversation = append(conversation, assistantMessage)
		}
	}

	return nil
}
