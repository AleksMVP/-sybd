package delivery

import (
	"github.com/AleksMVP/sybd/pkg/post"
	"github.com/AleksMVP/sybd/models"
	"github.com/AleksMVP/sybd/pkg/errors"
	"github.com/AleksMVP/sybd/utils"
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"log"
)

type PostDelivery struct {
	postUseCase post.IPostUseCase
}

func NewPostDelivery(postUseCase post.IPostUseCase) PostDelivery {
	return PostDelivery{
		postUseCase: postUseCase,
	}
}

func (h *PostDelivery)GetPostDetails(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["slug"]
	item := r.URL.Query().Get("related")

	post, err := h.postUseCase.GetPostFull(id, item)
	switch err {
	case nil:
		utils.WriteJson(w, http.StatusOK, post)
	case errors.ErrPostNotFound:
		utils.WriteJson(w, http.StatusNotFound, models.Error{Message: "Post not found"})
	default:
		utils.WriteJson(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}

	/*
	 * int id in slug
	 * array related in query
	 * Включение полной информации о соответвующем объекте сообщения.
	 * Если тип объекта не указан, то полная информация об этих объектах не передаётся.
	 * string items enum user, forum, thread
	 * 200 - OK -> PostFull
	 * 404 - ветка обсуждения отсутствует -> Error
	 */

}

func (h *PostDelivery)PostPostDetails(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["slug"]
	postUpdate := models.PostUpdate{}
	if err := json.NewDecoder(r.Body).Decode(&postUpdate); err != nil {
		log.Println(err)
		utils.WriteJson(w, http.StatusBadRequest, models.Error{Message: "Bad request"})
		return
	}

	post, err := h.postUseCase.UpdatePost(id, postUpdate)
	switch err {
	case nil:
		utils.WriteJson(w, http.StatusOK, post)
	case errors.ErrPostNotFound:
		utils.WriteJson(w, http.StatusNotFound, models.Error{Message: "Post not found"})
	default:
		utils.WriteJson(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}

	/*
	 * int id in slug
	 * PostUpdate in body
	 * 200 ОК -> Post
	 * 404 Сообщения нет в форуме -> Error
	 */
}

func (h *PostDelivery)PostCreateNewPosts(w http.ResponseWriter, r *http.Request) {
	slugOrId:= mux.Vars(r)["slug"]
	posts := models.Posts{}
	if err := json.NewDecoder(r.Body).Decode(&posts); err != nil {
		log.Println(err)
		utils.WriteJson(w, http.StatusBadRequest, models.Error{Message: "Bad request"})
		return
	}

	posts, err := h.postUseCase.CreateNewPosts(posts, slugOrId)
	switch err {
	case nil:
		utils.WriteJson(w, http.StatusCreated, &posts)
	case errors.ErrUserNotFound:
		utils.WriteJson(w, http.StatusNotFound, models.Error{Message: "User not found"})
	case errors.ErrThreadNotFound:
		utils.WriteJson(w, http.StatusNotFound, models.Error{Message: "Thread not found"})
	case errors.ErrPostNotFound:
		utils.WriteJson(w, http.StatusConflict, models.Error{Message: "Post not found"})
	default:
		utils.WriteJson(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}

	/*
	 * Thread slug or id in slug
	 * Posts in body
	 * 201 -> Posts
	 * 404 -> Error
	 * 409 -> Error хотя бы один родительский пост отсутствует в текущей ветке
	 */ 
}

func (h *PostDelivery)GetThreadPosts(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "100"
	}

	since := r.URL.Query().Get("since")
	sort := r.URL.Query().Get("sort")
	desc := r.URL.Query().Get("desc")
	if desc == "true" {
		desc = "DESC"
	} else if desc == "false" || desc == "" {
		desc = "ASC"
	}

	posts, err := h.postUseCase.GetPosts(slug, limit, since, desc, sort)
	switch err {
	case nil:
		utils.WriteJson(w, http.StatusOK, &posts)
	case errors.ErrThreadNotFound:
		utils.WriteJson(w, http.StatusNotFound, models.Error{Message: "Thread Not Found"})
	default:
		utils.WriteJson(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}

	/*
	 * Thread slug or id in slug
	 * int limit в query(min=1, max=10000, default=100)
	 * int since 
	 * string enum(flat, tree, parent_tree) sort
	 * bool desc
	 * 200 -> Posts
	 * 404 -> Error
	 */
}

