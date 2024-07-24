package chatgpt

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	chatgpt_errors "github.com/ayush6624/go-chatgpt/utils"
)

type ChatGPTModel string

const (
	GPT35Turbo        ChatGPTModel = "gpt-3.5-turbo"

	// Deprecated: Use gpt-3.5-turbo-0613 instead, model will discontinue on 09/13/2023
	GPT35Turbo0301    ChatGPTModel = "gpt-3.5-turbo-0301"
	
	GPT35Turbo0613    ChatGPTModel = "gpt-3.5-turbo-0613"
	GPT35Turbo16k     ChatGPTModel = "gpt-3.5-turbo-16k"
	GPT35Turbo16k0613 ChatGPTModel = "gpt-3.5-turbo-16k-0613"
	GPT4              ChatGPTModel = "gpt-4"
	
	// Deprecated: Use gpt-4-0613 instead, model will discontinue on 09/13/2023
	GPT4_0314         ChatGPTModel = "gpt-4-0314"
	
	GPT4_0613         ChatGPTModel = "gpt-4-0613"
	GPT4_32k          ChatGPTModel = "gpt-4-32k"
	
	// Deprecated: Use gpt-4-32k-0613 instead, model will discontinue on 09/13/2023
	GPT4_32k_0314     ChatGPTModel = "gpt-4-32k-0314"
	
	GPT4_32k_0613     ChatGPTModel = "gpt-4-32k-0613"
)

type ChatGPTModelRole string

const (
	ChatGPTModelRoleUser      ChatGPTModelRole = "user"
	ChatGPTModelRoleSystem    ChatGPTModelRole = "system"
	ChatGPTModelRoleAssistant ChatGPTModelRole = "assistant"
)

type ChatCompletionRequest struct {
	// (Required)
	// ID of the model to use.
	Model ChatGPTModel `json:"model"`

	// Required
	// The messages to generate chat completions for
	Messages []ChatMessage `json:"messages"`

	// (Optional - default: 1)
	// What sampling temperature to use, between 0 and 2. Higher values like 0.8 will make the output more random, while lower values like 0.2 will make it more focused and deterministic.
	// We generally recommend altering this or top_p but not both.
	Temperature float64 `json:"temperature,omitempty"`

	// (Optional - default: 1)
	// An alternative to sampling with temperature, called nucleus sampling, where the model considers the results of the tokens with top_p probability mass. So 0.1 means only the tokens comprising the top 10% probability mass are considered.
	// We generally recommend altering this or temperature but not both.
	Top_P float64 `json:"top_p,omitempty"`

	// (Optional - default: 1)
	// How many chat completion choices to generate for each input message.
	N int `json:"n,omitempty"`

	// (Optional - default: infinite)
	// The maximum number of tokens allowed for the generated answer. By default,
	// the number of tokens the model can return will be (4096 - prompt tokens).
	MaxTokens int `json:"max_tokens,omitempty"`

	// (Optional - default: 0)
	// Number between -2.0 and 2.0. Positive values penalize new tokens based on whether they appear in the text so far,
	// increasing the model's likelihood to talk about new topics.
	PresencePenalty float64 `json:"presence_penalty,omitempty"`

	// (Optional - default: 0)
	// Number between -2.0 and 2.0. Positive values penalize new tokens based on their existing frequency in the text so far,
	// decreasing the model's likelihood to repeat the same line verbatim.
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`

	// (Optional)
	// A unique identifier representing your end-user, which can help OpenAI to monitor and detect abuse
	User string `json:"user,omitempty"`

	//Additional fields in fork -- Ben Meeker

	// (Optional - default: false)
	// Return log probabilities of output tokens.
	LogProbs bool `json:"logprobs,omitempty"`

	// (Optional - default: null)
	// An integer between 0 and 20 specifying the most likely tokens to return at each token position.
	// LogProbs MUST be set to TRUE to use this parameter.
	Top_LogProbs int `json:"top_logprobs,omitempty"`

	// (Optional - default: text)
	// Specify the format the ChatGPT returns. Compatible with GPT-4o, GPT-4o mini, GPT-4 Turbo, and all GPT-3.5 Turbo models newer that gpt-3.5-turbo-1106
	// Options
	// Type: "json_object" to enable JSON mode.
	// Type: "text" to enable plain text mode.
	Response_Format *ResponseFormat `json:"response_format,omitempty"`

	// (Optional - default: null)
	// System will try to sample deterministically based on the seed provided. The same seed and parameters should return the same result.
	// Determinism is not guaranteed, refer to system_fingerprint response paramater.
	Seed int `json:"seed,omitempty"`

	// (Optional - default: auto)
	// Specifies latency tier to use for request
	// 'auto' - system will use scale tier credits until exhausted
	// 'default' - request processed using default service tier with lower uptime SLA and no latency guarantee.
	Service_Tier string `json:"service_tier,omitempty"`

	// (Optional - default: false)
	// If set, partial message deltas will be sent. Tokens will be send as data-only server-sent events as they become available.
	// Stream terminated by a data: [DONE] message.
	Stream bool `json:"stream,omitempty"`

	// (Optional - default: null)
	// Only set this when Stream is True
	// Set an additional chunk to stream before data: [DONE] message.
	Stream_Options *StreamOptions `json:"stream_options,omitempty"`

	// (Optional - default: null)
	// A list of tools the model may call
	// Provide a list of functions the model may generate JSON inputs for. 128 functions max supported.
	Tools *[]Tool `json:"tools,omitempty"`

	// (Optional - default: none)
	// Do NOT use this parameter in conjunction with Tool_Choice
	// Options
	// None: No tool will be called and a message will be generated
	// Auto: Any number of tools can be used and/or message generation will take place
	// Required: The model must call one or more tools
	Tool_Choice_Type string `json:"tool_choice,omitempty"`

	// (Optional - default: none)
	// Do NOT use this parameter in conjunction with Tool_Choice_Type
	// Provide a tool object to be called. This forces the model to use that tool.
	Tool_Choice *Tool `json:"tool_choice,omitempty"`

	// (Optional - default: true)
	// Whether to enable parallel function calling during tool use
	Parallel_Tool_Calls bool `json:"parallel_tool_calls,omitempty"`
}

type ChatMessage struct {
	Role    ChatGPTModelRole `json:"role"`
	Content string           `json:"content"`
}

type ChatResponse struct {
	ID        string               `json:"id"`
	Object    string               `json:"object"`
	CreatedAt int64                `json:"created_at"`
	Choices   []ChatResponseChoice `json:"choices"`
	Usage     ChatResponseUsage    `json:"usage"`
}

type ChatResponseChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type ChatResponseUsage struct {
	Prompt_Tokens     int `json:"prompt_tokens"`
	Completion_Tokens int `json:"completion_tokens"`
	Total_Tokens      int `json:"total_tokens"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type StreamOptions struct {
	Include_Usage bool `json:"include_usage"`
}

type Tool struct {
	Type string `json:"type"`
	Function FunctionFormat `json:"function"`
}

type FunctionFormat struct {
	Description string `json:"description"`
	Name string `json:"name"`
	Parameters interface{} `json:"parameters"`
}

func (c *Client) SimpleSend(ctx context.Context, message string) (*ChatResponse, error) {
	req := &ChatCompletionRequest{
		Model: GPT35Turbo,
		Messages: []ChatMessage{
			{
				Role:    ChatGPTModelRoleUser,
				Content: message,
			},
		},
	}

	return c.Send(ctx, req)
}

func (c *Client) Send(ctx context.Context, req *ChatCompletionRequest) (*ChatResponse, error) {
	if err := validate(req); err != nil {
		return nil, err
	}

	reqBytes, _ := json.Marshal(req)

	endpoint := "/chat/completions"
	httpReq, err := http.NewRequest("POST", c.config.BaseURL+endpoint, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}
	httpReq = httpReq.WithContext(ctx)

	res, err := c.sendRequest(ctx, httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var chatResponse ChatResponse
	if err := json.NewDecoder(res.Body).Decode(&chatResponse); err != nil {
		return nil, err
	}

	return &chatResponse, nil
}

func validate(req *ChatCompletionRequest) error {
	if len(req.Messages) == 0 {
		return chatgpt_errors.ErrNoMessages
	}

	isAllowed := false

	allowedModels := []ChatGPTModel{
		GPT35Turbo, GPT35Turbo0301, GPT35Turbo0613, GPT35Turbo16k, GPT35Turbo16k0613, GPT4, GPT4_0314, GPT4_0613, GPT4_32k, GPT4_32k_0314, GPT4_32k_0613,
	}

	for _, model := range allowedModels {
		if req.Model == model {
			isAllowed = true
		}
	}

	if !isAllowed {
		return chatgpt_errors.ErrInvalidModel
	}

	for _, message := range req.Messages {
		if message.Role != ChatGPTModelRoleUser && message.Role != ChatGPTModelRoleSystem && message.Role != ChatGPTModelRoleAssistant {
			return chatgpt_errors.ErrInvalidRole
		}
	}

	if req.Temperature < 0 || req.Temperature > 2 {
		return chatgpt_errors.ErrInvalidTemperature
	}

	if req.PresencePenalty < -2 || req.PresencePenalty > 2 {
		return chatgpt_errors.ErrInvalidPresencePenalty
	}

	if req.FrequencyPenalty < -2 || req.FrequencyPenalty > 2 {
		return chatgpt_errors.ErrInvalidFrequencyPenalty
	}

	return nil
}
