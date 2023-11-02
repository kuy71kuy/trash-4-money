package web

type ArticleResponse struct {
	ID        int    `json:"id" form:"id"`
	Title     string `json:"title" form:"title"`
	Text      string `json:"text" form:"text"`
	Link      string `json:"link" form:"link"`
	Thumbnail string `json:"thumbnail" form:"thumbnail"`
}
