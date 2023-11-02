package res

import (
	"app/model"
	"app/model/web"
)

func PassPointResponse(point *model.Point) web.PointResponse {
	return web.PointResponse{
		Id:     int(point.ID),
		Name:   point.Name,
		UserId: point.UserId,
		Amount: point.Amount,
	}
}

func PointConvertIndex(points []model.Point) []web.PointResponse {
	var results []web.PointResponse
	for _, point := range points {
		pointResponse := web.PointResponse{
			Id:     int(point.ID),
			Name:   point.Name,
			UserId: point.UserId,
			Amount: point.Amount,
		}
		results = append(results, pointResponse)
	}

	return results
}

func RankPointConvertIndex(points []model.Point) []web.RankPointResponse {
	var results []web.RankPointResponse
	iPosition := 1
	for _, point := range points {
		pointResponse := web.RankPointResponse{
			Position: iPosition,
			Id:       int(point.ID),
			Name:     point.Name,
			UserId:   point.UserId,
			Amount:   point.Amount,
		}
		results = append(results, pointResponse)
		iPosition++
	}

	return results
}
