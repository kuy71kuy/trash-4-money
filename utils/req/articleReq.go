package req

import (
	"app/model"
	"app/model/web"
)

func ArticlePassBody(article web.ArticleResponse) *model.Article {
	return &model.Article{
		Title:     article.Title,
		Text:      article.Text,
		Link:      article.Link,
		Thumbnail: article.Thumbnail,
	}
}

func ArticleAiPassBody(answerFromAi web.ChatResponse) *model.Article {
	return &model.Article{
		Title:     answerFromAi.Title,
		Text:      answerFromAi.Text,
		Link:      "aiCreated",
		Thumbnail: "aiCreated",
	}
}
