package web

type AiSuggestionRequest struct {
	TrashType string `json:"trash_type" form:"trash_type"`
}
