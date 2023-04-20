package ChatGpt

type AppChatGpt struct {
	urlApi string
	token  string

	model     string
	maxTokens int

	temperature int
	topP        int
	n           int

	stream bool
	stop   string
}

func New(token string) *AppChatGpt {
	return &AppChatGpt{
		urlApi: "https://api.openai.com/v1",
		token:  token,

		model:     "text-davinci-003",
		maxTokens: 2048,

		temperature: 0,
		topP:        1,
		n:           1,

		stream: false,
		stop:   "\n",
	}
}
