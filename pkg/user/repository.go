package user

import (
	"github.com/AleksMVP/sybd/models"
)

type IUserRepository interface {
	CreateUser(user models.User) (u models.User, e error)
	GetUsersByNickOrEmail(user models.User) (us models.Users)
	GetUserProfile(nickname string) (u models.User, e error) 
	EditUserProfile(nickname string, user models.UserUpdate) (u models.User, e error)
	GetUsers(slug, limit, since, desc string) (u models.Users, e error)
}