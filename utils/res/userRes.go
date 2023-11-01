package res

import (
	"app/model"
	"app/model/web"
)

func ConvertIndex(users []model.User) []web.GetUserReponse {
	var results []web.GetUserReponse
	for _, user := range users {
		userResponse := web.GetUserReponse{
			Id:    int(user.ID),
			Name:  user.Name,
			Email: user.Email,
		}
		results = append(results, userResponse)
	}

	return results
}

func ConvertGeneral(user *model.User) web.UserReponse {
	return web.UserReponse{
		Id:       int(user.ID),
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
}

func GetConvertGeneral(user *model.User) web.GetUserReponse {
	return web.GetUserReponse{
		Id:    int(user.ID),
		Name:  user.Name,
		Email: user.Email,
	}
}
