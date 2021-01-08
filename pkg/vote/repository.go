package vote

import (
	"github.com/AleksMVP/sybd/models"
)

type IVoteRepository interface {
	VoteForThread(t models.Thread, vote models.Vote) (thread models.Thread, e error)
}