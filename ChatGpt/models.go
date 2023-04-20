package ChatGpt

// OpenAICompletionsRequest Структура для отправки POST-запроса
type OpenAICompletionsRequest struct {
	Model       string      `json:"model"`
	Prompt      string      `json:"prompt"`
	MaxTokens   int         `json:"max_tokens"`
	Temperature int         `json:"temperature"`
	TopP        int         `json:"top_p"`
	N           int         `json:"n"`
	Stream      bool        `json:"stream"`
	Logprobs    interface{} `json:"logprobs"`
	Stop        string      `json:"stop"`
}

// OpenAICompletionsResponse Структура для получения ответа от OpenAI API
type OpenAICompletionsResponse struct {
	Choices []struct {
		Text      string  `json:"text"`
		Index     int     `json:"index"`
		Logprobs  *string `json:"logprobs"`
		Finish    *string `json:"finish_reason"`
		Prompting *string `json:"prompt"`
	} `json:"choices"`
	Error *string `json:"error"`
}

type OpenAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type OpenAIChatRequest struct {
	Stream   bool                `json:"stream,omitempty"`
	Model    string              `json:"model" json:"model"`
	Messages []OpenAIChatMessage `json:"messages" json:"messages"`
}
type OpenAIChatResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}
