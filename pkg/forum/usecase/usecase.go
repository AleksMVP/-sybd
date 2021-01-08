package usecase

import (
	"github.com/AleksMVP/sybd/pkg/forum"
	"github.com/AleksMVP/sybd/models"
	"github.com/AleksMVP/sybd/pkg/errors"
)

type ForumUseCase struct {
	forumRepository forum.IForumRepository
}

func NewForumUseCase(forumRepository forum.IForumRepository) ForumUseCase {
	return ForumUseCase{
		forumRepository: forumRepository,
	}
}

func (h *ForumUseCase)CreateForum(forum models.Forum) (f models.Forum, e error) {
	f, err := h.forumRepository.CreateForum(forum)
	switch err {
	case errors.ErrForumExist:
		f, _ := h.forumRepository.GetForum(forum.Slug)
		return f, errors.ErrForumExist
	}

	return f, err
}

func (h *ForumUseCase)GetForum(slug string) (f models.Forum, e error) {
	return h.forumRepository.GetForum(slug)
}	