package forum

import (
	"github.com/AleksMVP/sybd/models"
)

type IForumUseCase interface {
	CreateForum(forum models.Forum) (f models.Forum, e error)
	GetForum(slug string) (f models.Forum, e error)
}