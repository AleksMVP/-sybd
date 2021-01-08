package repository

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/AleksMVP/sybd/models"
	"github.com/AleksMVP/sybd/pkg/errors"
)

type ForumRepository struct {
	db *sql.DB
}

func NewForumRepository(db *sql.DB) ForumRepository {
	return ForumRepository{
		db: db,
	}
}

func (h *ForumRepository)CreateForum(forum models.Forum) (f models.Forum, e error) {
	err := h.db.QueryRow(
		`INSERT INTO forums(title, nickname, slug) 
		 VALUES($1, (SELECT nickname FROM users WHERE nickname = $2), $3)
		 RETURNING nickname`,
		forum.Title,
		forum.User,
		forum.Slug,
	).Scan(&forum.User)

	if err, ok := err.(*pq.Error); ok {
		switch err.Code {
		case "23502": // User not found 404
			return f, errors.ErrUserNotFound
		case "23505": // Already exist
			return f, errors.ErrForumExist
		}
	}

	return forum, err
}

func (h *ForumRepository)GetForum(slug string) (f models.Forum, e error) {
	err := h.db.QueryRow(
		`SELECT title, nickname, posts, threads, slug
		 FROM forums 
		 WHERE slug = $1`,
		slug,
	).Scan(&f.Title, &f.User, &f.Posts, &f.Threads, &f.Slug)

	if err == sql.ErrNoRows {
		return f, errors.ErrForumNotFound
	}

	return f, err
}