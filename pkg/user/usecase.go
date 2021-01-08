package user

import (
	"github.com/AleksMVP/sybd/models"
	"github.com/mailru/easyjson"
)

type IUserUseCase interface {
	CreateUser(user models.User) (u easyjson.Marshaler, e error)
	GetUserProfile(nickname string) (u models.User, e error) 
	EditUserProfile(nickname string, user models.UserUpdate) (u models.User, e error)
	GetUsers(slug, limit, since, desc string) (u models.Users, e error)
}