package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/869413421/chatgpt-web/pkg/logger"
	"github.com/byebyebruce/ghsearch/util"
	gogpt "github.com/sashabaranov/go-gpt3"
)

// Configuration 项目配置
type Configuration struct {
	// gpt apikey
	ApiKey string `json:"api_key"`
	Port   int    `json:"port"`
	//
	BotDesc string `json:"bot_desc"`
	// GPT请求最大字符数
	MaxTokens int `json:"max_tokens"`
	// GPT模型
	Model string `json:"model"`
	// 热度
	Temperature      float64 `json:"temperature"`
	TopP             float32 `json:"top_p"`
	PresencePenalty  float32 `json:"presence_penalty"`
	FrequencyPenalty float32 `json:"frequency_penalty"`
}

var cnf = &Configuration{
	MaxTokens:        2048,
	Port:             8080,
	Model:            "text-davinci-003",
	Temperature:      0.9,
	TopP:             1,
	FrequencyPenalty: 0.0,
	PresencePenalty:  0.6,
}

var bodDesc = "以下是与AI助手的对话。助手乐于助人，富有创造力，聪明且非常友好。"

func main() {
	apiKey := os.Getenv("CHATGPT_KEY")
	client := gogpt.NewClient(apiKey)
	q := bodDesc
	for {
		var prompt string
		fmt.Scanln(&prompt)
		prompt = strings.TrimSpace(prompt)
		if len(prompt) == 0 {
			continue
		}
		q += "\n" + prompt
		//prompt := cnf.BotDesc + "\n" + question.Text
		logger.Info("request prompt is", prompt)
		req := gogpt.CompletionRequest{
			Model:            cnf.Model,
			MaxTokens:        cnf.MaxTokens,
			TopP:             cnf.TopP,
			FrequencyPenalty: cnf.FrequencyPenalty,
			PresencePenalty:  cnf.PresencePenalty,
			Prompt:           q,
		}

		resp, err := util.AsyncTaskAndShowLoadingBar("正在想", func() (gogpt.CompletionResponse, error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
			defer cancel()
			resp, err := client.CreateCompletion(ctx, req)
			return resp, err
		})
		if err != nil {
			fmt.Println("error", err)
		} else {
			fmt.Println(resp.Choices[0].Text)
		}
	}
}
