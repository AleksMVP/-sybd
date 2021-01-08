package post

import (
	"github.com/AleksMVP/sybd/models"
)

type IPostRepository interface {
	CreateNewPosts(posts models.Posts, thread models.Thread) (p models.Posts, e error)
	GetPosts(limit, since, desc, sort string, thread models.Thread) (p models.Posts, e error)
	UpdatePost(id string, post models.PostUpdate) (p models.Post, e error)
	GetPostFull(id, item string) (p models.PostFull, e error) 
}