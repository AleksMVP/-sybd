package delivery

import (
	"github.com/AleksMVP/sybd/pkg/forum"
	"github.com/AleksMVP/sybd/models"
	"github.com/AleksMVP/sybd/pkg/errors"
	"github.com/AleksMVP/sybd/utils"
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"log"
)

type ForumDelivery struct {
	forumUseCase forum.IForumUseCase
}

func NewForumDelivery(forumUseCase forum.IForumUseCase) ForumDelivery {
	return ForumDelivery{
		forumUseCase: forumUseCase,
	}
}

func (h *ForumDelivery)PostForumCreate(w http.ResponseWriter, r *http.Request) {
	forum := models.Forum{}
	if err := json.NewDecoder(r.Body).Decode(&forum); err != nil {
		log.Println(err)
		utils.WriteJson(w, http.StatusBadRequest, models.Error{Message: "Bad request"})
		return
	}

	newForum, err := h.forumUseCase.CreateForum(forum)
	switch err {
	case nil:
		utils.WriteJson(w, http.StatusCreated, newForum)
	case errors.ErrUserNotFound:
		utils.WriteJson(w, http.StatusNotFound, models.Error{Message: "User not found"})
	case errors.ErrForumExist:
		utils.WriteJson(w, http.StatusConflict, newForum)
	default:
		utils.WriteJson(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}

	/* 
	 * Параметры Forum в body
	 * 201 - успешно создан -> форум Forum
	 * 404 - владелец форума не найден -> ошибку Error
	 * 409 - форум уже существует в бд -> форум Forum
	 */
}

func (h *ForumDelivery)GetForumDetails(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	forum, err := h.forumUseCase.GetForum(slug)

	switch err {
	case nil:
		utils.WriteJson(w, http.StatusOK, forum)
	case errors.ErrForumNotFound:
		utils.WriteJson(w, http.StatusNotFound, models.Error{Message: "Forum not found"})
	default:
		utils.WriteJson(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}
	/* 
	 * Имя форума в slug
	 * 200 - OK -> форум Forum
	 * 404 - форума нет в системе -> Error
	 */
}