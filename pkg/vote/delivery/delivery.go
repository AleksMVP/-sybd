package delivery

import (
	"github.com/AleksMVP/sybd/pkg/vote"
	"github.com/AleksMVP/sybd/models"
	"github.com/AleksMVP/sybd/pkg/errors"
	"github.com/AleksMVP/sybd/utils"
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"log"
)

type VoteDelivery struct {
	voteUseCase vote.IVoteUseCase
}

func NewVoteDelivery(voteUseCase vote.IVoteUseCase) VoteDelivery {
	return VoteDelivery{
		voteUseCase: voteUseCase,
	}
}

func (h *VoteDelivery)PostThreadVote(w http.ResponseWriter, r *http.Request) {
	slugOrId := mux.Vars(r)["slug"]
	vote := models.Vote{}
	if err := json.NewDecoder(r.Body).Decode(&vote); err != nil {
		log.Println(err)
		utils.WriteJson(w, http.StatusBadRequest, models.Error{Message: "Bad request"})
		return
	}

	thread, err := h.voteUseCase.VoteForThread(slugOrId, vote)
	switch err {
	case nil:
		utils.WriteJson(w, http.StatusOK, thread)
	case errors.ErrThreadNotFound, errors.ErrUserNotFound:
		utils.WriteJson(w, http.StatusNotFound, models.Error{Message: "Thread not found"})
	default:
		utils.WriteJson(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}

	/*
	 * Голосование за ветвь
	 * Thread slug or id in slug
	 * Vote in body
	 * 200 -> Thread
	 * 404 -> Error
	 */
}