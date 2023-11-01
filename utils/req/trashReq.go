package req

import (
	"app/model"
	"app/model/web"
)

func TrashPassBody(trash web.TrashResponse) *model.Trash {
	return &model.Trash{
		UserId:  trash.UserId,
		Type:    trash.Type,
		Weight:  trash.Weight,
		Address: trash.Address,
		Image:   trash.Image,
		Note:    trash.Note,
		Status:  trash.Status,
	}
}
