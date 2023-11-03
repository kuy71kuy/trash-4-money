package res

import (
	"app/model"
	"app/model/web"
)

func ConvertIndex(users []model.User) []web.GetUserResponse {
	var results []web.GetUserResponse
	for _, user := range users {
		userResponse := web.GetUserResponse{
			Id:    int(user.ID),
			Name:  user.Name,
			Email: user.Email,
		}
		results = append(results, userResponse)
	}

	return results
}

func ConvertGeneral(user *model.User) web.UserResponse {
	return web.UserResponse{
		Id:       int(user.ID),
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
}

func GetConvertGeneral(user *model.User) web.GetUserResponse {
	return web.GetUserResponse{
		Id:    int(user.ID),
		Name:  user.Name,
		Email: user.Email,
	}
}
