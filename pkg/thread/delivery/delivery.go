package delivery

import (
	"github.com/AleksMVP/sybd/pkg/thread"
	"github.com/AleksMVP/sybd/models"
	"github.com/AleksMVP/sybd/pkg/errors"
	"github.com/AleksMVP/sybd/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"encoding/json"
	"log"
)

type ThreadDelivery struct {
	threadUseCase thread.IThreadUseCase
}

func NewThreadDelivery(threadUseCase thread.IThreadUseCase) ThreadDelivery {
	return ThreadDelivery{
		threadUseCase: threadUseCase,
	}
}

func (h *ThreadDelivery)PostForumCreateThread(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	thread := models.Thread{}
	if err := json.NewDecoder(r.Body).Decode(&thread); err != nil {
		log.Println(err)
		utils.WriteJson(w, http.StatusBadRequest, models.Error{Message: "Bad request"})
		return
	}

	thread.Forum = slug

	newThread, err := h.threadUseCase.CreateThread(thread)
	switch err {
	case nil:
		utils.WriteJson(w, http.StatusCreated, newThread)
	case errors.ErrUserOrForumNotFound:
		utils.WriteJson(w, http.StatusNotFound, models.Error{Message: "User or Forum not found"})
	case errors.ErrThreadExist:
		utils.WriteJson(w, http.StatusConflict, newThread)
	default:
		utils.WriteJson(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}

	/*
	 * Имя форума в slug
	 * Параметры Thread в body
	 * 201 - ветка создана -> Thread
	 * 404 - Автор ветки или форум не найдены -> Error
	 * 409 - Ветка уже есть в бд -> Thread
	 */
}

func (h *ThreadDelivery)GetForumThreads(w http.ResponseWriter, r *http.Request) {
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

	threads, err := h.threadUseCase.GetThreads(slug, limit, since, desc)
	switch err {
	case nil:
		utils.WriteJson(w, http.StatusOK, &threads)
	case errors.ErrForumNotFound:
		utils.WriteJson(w, http.StatusNotFound, &models.Error{Message: "Forum not found"})
	default:
		utils.WriteJson(w, http.StatusInternalServerError, &models.Error{Message: err.Error()})
	}

	/*
	 * Name in slug
	 * int limit(min=1, max=10000, default=100)
	 * string date since дата создание ветви 
	 * bool desc
	 * 200 - OK -> Threads
	 * 404 - Not Found -> Error 
	 */
}

func (h *ThreadDelivery)GetThreadDetails(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	_, err := strconv.Atoi(slug)

	thread, err := h.threadUseCase.GetThread(slug)

	switch err {
	case nil:
		utils.WriteJson(w, http.StatusOK, thread)
	case errors.ErrThreadNotFound:
		utils.WriteJson(w, http.StatusNotFound, models.Error{Message: "Forum not found"})
	default:
		utils.WriteJson(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}
	/*
	 * Thread slug or id in slug
	 * 200 -> Thread
	 * 404 -> Error
	 */
}

func (h *ThreadDelivery)PostThreadDetails(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	thread := models.Thread{}
	if err := json.NewDecoder(r.Body).Decode(&thread); err != nil {
		log.Println(err)
		utils.WriteJson(w, http.StatusBadRequest, models.Error{Message: "Bad request"})
		return
	}

	newThread, err := h.threadUseCase.EditThread(slug, thread)
	switch err {
	case nil:
		utils.WriteJson(w, http.StatusOK, newThread)
	case errors.ErrThreadNotFound:
		utils.WriteJson(w, http.StatusNotFound, models.Error{Message: "ThreadNotFound"})
	default:
		utils.WriteJson(w, http.StatusInternalServerError, models.Error{Message: err.Error()})
	}
	/*
	 * Thread slug or id in slug
	 * Thread in body
	 * 200 -> Thread
	 * 404 -> Error
	 */
}