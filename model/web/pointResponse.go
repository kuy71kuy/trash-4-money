package web

type PointResponse struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	UserId int    `json:"user_id"`
	Amount int    `json:"amount" form:"amount"`
}

type RankPointResponse struct {
	Position int    `json:"position"`
	Id       int    `json:"id"`
	Name     string `json:"name"`
	UserId   int    `json:"user_id"`
	Amount   int    `json:"amount" form:"amount"`
}
