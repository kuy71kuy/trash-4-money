package utils

import (
	"app/model/web"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ayush6624/go-chatgpt"
	"log"
	"os"
)

func AskAiArticle(exclude string) web.ChatResponse {
	tokenAi := os.Getenv("SECRET_OPENAI_KEY")
	client, err := chatgpt.NewClient(tokenAi)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	promptText := "Buatkan artikel mengenai menjaga lingkungan, dalam bentuk json, title dan text. title adalah judul yang bertipe data string\ntext adalah isi artikel, isi artikel harus dalam bentuk plain teks.\naku ingin kamu memberikan response berupa jsonnya saja.\ncontoh json\n{\n\"title\": \"judul artikel disini\",\n  \"text\": \"teks artikel disini\"\n}\njangan gunakan artikel yang berjudul : " + exclude
	res, _ := client.SimpleSend(ctx, promptText)

	var data web.ChatResponse
	jsonString := res.Choices[0].Message.Content
	err = json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		fmt.Println("Error:", err)
	}
	return data
}

func AiSuggestionHowTo(trashType string) string {
	tokenAi := os.Getenv("SECRET_OPENAI_KEY")
	client, err := chatgpt.NewClient(tokenAi)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	promptText := "Aku memiliki sampah ini : " + trashType + ", bagaimana cara mengelola sampah tersebut agar memiliki nilai jual ?"
	res, _ := client.SimpleSend(ctx, promptText)

	answerString := res.Choices[0].Message.Content

	return answerString
}
