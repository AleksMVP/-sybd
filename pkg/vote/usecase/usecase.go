package usecase

import (
	"github.com/AleksMVP/sybd/pkg/vote"
	"github.com/AleksMVP/sybd/pkg/thread"
	"github.com/AleksMVP/sybd/models"
	"github.com/AleksMVP/sybd/pkg/errors"
)

type VoteUseCase struct {
	voteRepository vote.IVoteRepository
	threadRepository thread.IThreadRepository
}

func NewVoteUseCase(voteRepository vote.IVoteRepository,
		 			threadRepository thread.IThreadRepository) VoteUseCase {
	return VoteUseCase{
		voteRepository: voteRepository,
		threadRepository: threadRepository,
	}
}

func (h *VoteUseCase)VoteForThread(slugOrId string, vote models.Vote) (thread models.Thread, e error) {
	thread, e = h.threadRepository.GetThread(slugOrId)
	if e != nil {
		return thread, errors.ErrThreadNotFound
	}

	return h.voteRepository.VoteForThread(thread, vote)
}