package web

type TrashResponse struct {
	ID      int    `json:"id" form:"id"`
	UserId  int    `json:"user_id" form:"user_id"`
	Type    string `json:"type" form:"type"`
	Weight  int    `json:"weight" form:"weight"`
	Address string `json:"address" form:"address"`
	Image   string `json:"image" form:"image"`
	Note    string `json:"note" form:"note"`
	Status  string `json:"status" form:"status"`
}

type TrashStatusResponse struct {
	Status string `json:"status" form:"status"`
}
