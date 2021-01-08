package delivery 

import (
	"github.com/AleksMVP/sybd/pkg/user"
	"github.com/AleksMVP/sybd/models"
	"github.com/AleksMVP/sybd/pkg/errors"
	"github.com/AleksMVP/sybd/utils"
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"log"
)

type UserDelivery struct {
	userUseCase user.IUserUseCase
}

func NewUserDelivery(userUseCase user.IUserUseCase) UserDelivery {
	return UserDelivery{
		userUseCase: userUseCase,
	}
}

func (h *UserDelivery)PostUserCreate(w http.ResponseWriter, r *http.Request) {
	nickname := mux.Vars(r)["slug"]
	user := &models.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		log.Println(err)
		utils.WriteJson(w, http.StatusBadRequest, models.Error{Message: "Bad request"})
		return
	}

	user.Nickname = nickname
	users, err := h.userUseCase.CreateUser(*user)
	switch err {
	case nil:
		utils.WriteJson(w, http.StatusCreated, users)
	case errors.ErrUserOrEmailExist:
		utils.WriteJson(w, http.StatusConflict, users)
	default:
		utils.WriteJson(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}

	/*
	 * nickname in slug
	 * User in body
	 * 201 -> User
	 * 409 -> Users с таким же nickname or email
	 */
}

func (h *UserDelivery)GetUserProfile(w http.ResponseWriter, r *http.Request) {
	nickname := mux.Vars(r)["slug"]
	user, err := h.userUseCase.GetUserProfile(nickname)
 
	switch err {
	case nil:
		utils.WriteJson(w, http.StatusOK, user)
	case errors.ErrUserNotFound:
		utils.WriteJson(w, http.StatusNotFound, models.Error{Message: "User not found"})
	default:
		utils.WriteJson(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}

	/*
	 * nickname in slug
	 * 200 -> User
	 * 404 -> Error
	 */
}

func (h *UserDelivery)PostUserProfile(w http.ResponseWriter, r *http.Request) {
	nickname := mux.Vars(r)["slug"]
	userUpdate := models.UserUpdate{}
	if err := json.NewDecoder(r.Body).Decode(&userUpdate); err != nil {
		log.Println(err)
		utils.WriteJson(w, http.StatusBadRequest, models.Error{Message: "Bad request"})
		return
	}

	user, err := h.userUseCase.EditUserProfile(nickname, userUpdate)
	switch err {
	case nil:
		utils.WriteJson(w, http.StatusOK, user)
	case errors.ErrUserNotFound:
		utils.WriteJson(w, http.StatusNotFound, models.Error{Message: "User not found"})
	case errors.ErrUserExist:
		utils.WriteJson(w, http.StatusConflict, models.Error{Message: "User exist"})
	default:
		utils.WriteJson(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}

	/*
	 * nickname in slug 
	 * User in body
	 * 200 -> User 
	 * 404 -> Error
	 * 409 -> Error Конфлик уже с имеющимися
	 */

}

func (h *UserDelivery)GetForumUsers(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "100"
	}

	since := r.URL.Query().Get("since")
	desc := r.URL.Query().Get("desc")
	if desc == "true" {
		desc = "DESC"
	} else if desc == "false" || desc == "" {
		desc = "ASC"
	}

	users, err := h.userUseCase.GetUsers(slug, limit, since, desc)
	switch err {
	case nil:
		utils.WriteJson(w, http.StatusOK, &users)
	case errors.ErrForumNotFound:
		utils.WriteJson(w, http.StatusNotFound, &models.Error{Message: "Forum not found"})
	default:
		utils.WriteJson(w, http.StatusInternalServerError, &models.Error{Message: err.Error()})
	}

	/* 
	 * Имя форума в slug
	 * int limit в query(min=1, max=10000, default=100) в query
	 * since имя с которого будут выводится пользователи в query
	 * bool desc флаг сортировки по убыванию в query
	 * 200 - OK -> Users
	 * 404 - форум отсутствует в системе -> Error
	 */
}