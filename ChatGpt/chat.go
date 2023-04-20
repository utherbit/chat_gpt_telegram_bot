package ChatGpt

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"time"
)

var updateMessageSleep = time.Second

func (app *AppChatGpt) CallOpenAI(prompt string) (*OpenAICompletionsResponse, error) {
	fmt.Printf("\nCallOpenAI %s", prompt)
	requestData := OpenAICompletionsRequest{
		Model:  app.model,
		Prompt: prompt,

		MaxTokens:   app.maxTokens,
		Temperature: app.temperature,
		//TopP:        app.topP,
		N: app.n,

		Stream: app.stream,
		//Stop:   app.stop,
	}

	requestBytes, _ := json.Marshal(requestData)
	reader := bytes.NewReader(requestBytes)

	request, _ := http.NewRequest("POST", app.urlApi+"/completions", reader)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", app.token))

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var responseData OpenAICompletionsResponse
	fmt.Printf("Respinse %s", body)

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Printf("\nError %v", err)
		return nil, err
	}

	fmt.Printf("\n%v", responseData)
	return &responseData, nil
}

func (app *AppChatGpt) ChatOpenAI(messages []OpenAIChatMessage) (*OpenAIChatResponse, error) {
	fmt.Printf("\nCallOpenAI ")
	requestData := OpenAIChatRequest{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
	}

	requestBytes, _ := json.Marshal(requestData)
	reader := bytes.NewReader(requestBytes)

	request, _ := http.NewRequest("POST", app.urlApi+"/chat/completions", reader)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", app.token))

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var responseData *OpenAIChatResponse
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		//fmt.Printf("\nError %v", err)
		return nil, err
	}

	//fmt.Printf("\n%v", responseData)
	return responseData, nil
}

func (app *AppChatGpt) StreamChatOpenAI(messages []OpenAIChatMessage, writer func(inp string) error) (string, error) {
	fmt.Printf("\nCallOpenAI ")
	requestData := OpenAIChatRequest{
		Stream:   true,
		Model:    "gpt-3.5-turbo",
		Messages: messages,
	}

	requestBytes, _ := json.Marshal(requestData)
	reader := bytes.NewReader(requestBytes)

	request, _ := http.NewRequest("POST", app.urlApi+"/chat/completions", reader)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", app.token))

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	scanner := bufio.NewScanner(response.Body)

	fullText := ""
	sendText := ""
	whileWriting := true
	go func() {
		for whileWriting {
			if fullText != sendText {
				sendText = fullText
				writer(fullText)
			}
			time.Sleep(updateMessageSleep)
		}
		if fullText != sendText {
			sendText = fullText
			writer(sendText)
		}
		fmt.Printf("\nSuc Send")
	}()

	for scanner.Scan() {
		if len(scanner.Bytes()) > 6 {
			//fmt.Printf("\n\n %s", scanner.Text()[6:])
			var pac PacketStream

			err2 := jsoniter.Unmarshal(scanner.Bytes()[6:], &pac)
			if err2 != nil {
				fmt.Printf("Error: %v", err2)
			}

			//fmt.Printf("\nScan %v \n", scanner.Text())

			for _, choice := range pac.Choices {
				//fmt.Printf("content: %s;", choice.Delta.Content)
				fullText += choice.Delta.Content
			}
			//go writer(fullText)
		}
	}
	whileWriting = false
	fmt.Printf("\nSendText %s", fullText)
	//reader := bufio.NewReaderSize(, 4096) // Создание буферизованного Reader.
	//
	//body, err := ioutil.ReadAll(response.Body)
	//if err != nil {
	//	return  err
	//}
	//
	//var responseData *PacketStream
	//err = json.Unmarshal(body, &responseData)
	//if err != nil {
	//	//fmt.Printf("\nError %v", err)
	//	return err
	//}

	//fmt.Printf("\n%v", responseData)
	return fullText, nil
}

type PacketStream struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		} `json:"delta"`
		Index        int         `json:"index"`
		FinishReason interface{} `json:"finish_reason"`
	} `json:"choices"`
}
