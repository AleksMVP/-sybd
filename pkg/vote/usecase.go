package vote

import (
	"github.com/AleksMVP/sybd/models"
)

type IVoteUseCase interface {
	VoteForThread(slugOrId string, vote models.Vote) (thread models.Thread, e error)
}