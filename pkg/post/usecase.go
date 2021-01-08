package post

import (
	"github.com/AleksMVP/sybd/models"
)

type IPostUseCase interface {
	CreateNewPosts(posts models.Posts, slugOrId string) (p models.Posts, e error)
	GetPosts(slug, limit, since, desc, sort string) (p models.Posts, e error)
	UpdatePost(id string, post models.PostUpdate) (p models.Post, e error)
	GetPostFull(id, item string) (p models.PostFull, e error) 
}