package usecase

import (
	"github.com/AleksMVP/sybd/pkg/forum"
	"github.com/AleksMVP/sybd/pkg/thread"
	"github.com/AleksMVP/sybd/models"
	"github.com/AleksMVP/sybd/pkg/errors"
)

type ThreadUseCase struct {
	threadRepository thread.IThreadRepository
	forumRepository forum.IForumRepository
}

func NewThreadUseCase(threadRepository thread.IThreadRepository,
					  forumRepository forum.IForumRepository) ThreadUseCase {
	return ThreadUseCase{
		threadRepository: threadRepository,
		forumRepository: forumRepository,
	}
}

func (h *ThreadUseCase)CreateThread(thread models.Thread) (t models.Thread, e error) {
	t, e = h.threadRepository.CreateThread(thread)

	switch e {
	case errors.ErrThreadExist:
		t, _ = h.GetThread(thread.Slug)
	}

	return t, e
}

func (h *ThreadUseCase)GetThread(slugOrId string) (t models.Thread, e error) {
	return h.threadRepository.GetThread(slugOrId)
}

func (h *ThreadUseCase)GetThreads(slug, limit, since, desc string) (ts models.Threads, e error) {
	ts, e = h.threadRepository.GetThreads(slug, limit, since, desc)

	switch e {
	case errors.ErrForumNotFound:
		if _, err := h.forumRepository.GetForum(slug); err != nil {
			return ts, errors.ErrForumNotFound
		}
			
		return ts, nil
	}

	return ts, e
}

func (h *ThreadUseCase)EditThread(slugOrId string, thread models.Thread) (t models.Thread, e error) {
	return h.threadRepository.EditThread(slugOrId, thread)
}