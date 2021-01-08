package usecase

import (
	"github.com/AleksMVP/sybd/pkg/user"
	"github.com/AleksMVP/sybd/pkg/forum"
	"github.com/AleksMVP/sybd/models"
	"github.com/AleksMVP/sybd/pkg/errors"
	"github.com/mailru/easyjson"
)

type UserUseCase struct {
	userRepository user.IUserRepository
	forumRepository forum.IForumRepository
}

func NewUserUseCase(userRepository user.IUserRepository,
					forumRepository forum.IForumRepository) UserUseCase {
		return UserUseCase{
			userRepository: userRepository,
			forumRepository: forumRepository,
		}
	}

func (h *UserUseCase)CreateUser(user models.User) (u easyjson.Marshaler, e error) {
	u, e = h.userRepository.CreateUser(user)
	switch e {
	case errors.ErrUserOrEmailExist:
		q := h.userRepository.GetUsersByNickOrEmail(user)
		return q, e
	}

	return u, e
}

func (h *UserUseCase)GetUserProfile(nickname string) (u models.User, e error) {
	return h.userRepository.GetUserProfile(nickname)
}

func (h *UserUseCase)EditUserProfile(nickname string, user models.UserUpdate) (u models.User, e error) {
	return h.userRepository.EditUserProfile(nickname, user)
}

func (h *UserUseCase)GetUsers(slug, limit, since, desc string) (u models.Users, e error) {
	if _, err := h.forumRepository.GetForum(slug); err != nil {
		return u, errors.ErrForumNotFound
	}

	return h.userRepository.GetUsers(slug, limit, since, desc)
}