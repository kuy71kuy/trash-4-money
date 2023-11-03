package res

import (
	"app/model"
	"app/model/web"
)

func ArticleConvertIndex(articles []model.Article) []web.ArticleResponse {
	var results []web.ArticleResponse
	for _, article := range articles {
		articleResponse := web.ArticleResponse{
			ID:        int(article.ID),
			Title:     article.Title,
			Text:      article.Text,
			Link:      article.Link,
			Thumbnail: article.Thumbnail,
		}
		results = append(results, articleResponse)
	}

	return results
}

func xConvertGeneral(user *model.User) web.UserResponse {
	return web.UserResponse{
		Id:       int(user.ID),
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
}

func PassArticleBody(article *model.Article) web.ArticleResponse {
	return web.ArticleResponse{
		ID:        int(article.ID),
		Title:     article.Title,
		Text:      article.Text,
		Link:      article.Link,
		Thumbnail: article.Thumbnail,
	}
}
