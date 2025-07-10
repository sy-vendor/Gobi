package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gobi/config"
	"io/ioutil"
	"net/http"
)

type DeepSeekRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Search      bool      `json:"search,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type DeepSeekResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func CallDeepSeek(prompt string, withSearch bool) (string, error) {
	apiKey := config.GetConfig().AI.DeepSeekAPIKey
	reqBody := DeepSeekRequest{
		Model: "deepseek-chat",
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
		Search:      withSearch,
		Temperature: 0.7,
		MaxTokens:   2000,
	}
	data, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("DeepSeek prompt:", prompt)
	fmt.Println("DeepSeek raw response:", string(body))

	var dsResp DeepSeekResponse
	json.Unmarshal(body, &dsResp)
	if len(dsResp.Choices) > 0 {
		return dsResp.Choices[0].Message.Content, nil
	}
	return "", nil
}
