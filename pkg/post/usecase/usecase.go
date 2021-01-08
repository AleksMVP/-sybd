package usecase

import (
	"github.com/AleksMVP/sybd/pkg/post"
	"github.com/AleksMVP/sybd/pkg/thread"
	"github.com/AleksMVP/sybd/models"
	"github.com/AleksMVP/sybd/pkg/errors"
)

type PostUseCase struct {
	postDelivery post.IPostRepository
	threadDelivery thread.IThreadRepository
}

func NewPostUseCase(postDelivery post.IPostRepository,
					threadDelivery thread.IThreadRepository) PostUseCase {
	return PostUseCase{
		postDelivery: postDelivery,
		threadDelivery: threadDelivery,
	}
}

func (h *PostUseCase)CreateNewPosts(posts models.Posts, slugOrId string) (p models.Posts, e error) {
	thread, err := h.threadDelivery.GetThread(slugOrId)
	if err != nil {
		return p, err
	}

	return h.postDelivery.CreateNewPosts(posts, thread)
}

func (h *PostUseCase)GetPosts(slug, limit, since, desc, sort string) (p models.Posts, e error) {
	thread, err := h.threadDelivery.GetThread(slug)
	if err != nil {
		return p, errors.ErrThreadNotFound
	}

	return h.postDelivery.GetPosts(limit, since, desc, sort, thread)
}

func (h *PostUseCase)UpdatePost(id string, post models.PostUpdate) (p models.Post, e error) {
	return h.postDelivery.UpdatePost(id, post)
}

func (h *PostUseCase)GetPostFull(id, item string) (p models.PostFull, e error) {
	return h.postDelivery.GetPostFull(id, item)
}

