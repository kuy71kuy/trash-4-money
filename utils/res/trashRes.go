package res

import (
	"app/model"
	"app/model/web"
)

func TrashConvertIndex(trashes []model.Trash) []web.TrashResponse {
	var results []web.TrashResponse
	for _, trash := range trashes {
		trashResponse := web.TrashResponse{
			ID:      int(trash.ID),
			UserId:  trash.UserId,
			Type:    trash.Type,
			Weight:  trash.Weight,
			Address: trash.Address,
			Image:   trash.Image,
			Note:    trash.Note,
			Status:  trash.Status,
		}
		results = append(results, trashResponse)
	}

	return results
}

func PassTrashBody(trash *model.Trash) web.TrashResponse {
	return web.TrashResponse{
		ID:      int(trash.ID),
		UserId:  trash.UserId,
		Type:    trash.Type,
		Weight:  trash.Weight,
		Address: trash.Address,
		Image:   trash.Image,
		Note:    trash.Note,
		Status:  trash.Status,
	}
}

func PassTrashStatusBody(trash *model.Trash) web.TrashStatusResponse {
	return web.TrashStatusResponse{
		Status: trash.Status,
	}
}
