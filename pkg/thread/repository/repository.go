package repository

import (
	"github.com/AleksMVP/sybd/models"
	"github.com/AleksMVP/sybd/pkg/errors"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"strconv"
)

type ThreadRepository struct {
	db *sql.DB
}

func NewThreadRepository(db *sql.DB) ThreadRepository {
	return ThreadRepository{
		db: db,
	}
}

func (h *ThreadRepository)CreateThread(thread models.Thread) (t models.Thread, e error) {
	err := h.db.QueryRow(
		`INSERT INTO threads(author, forum, msg, slug, title, create_date)
		 VALUES(
			 (SELECT nickname FROM users WHERE nickname = $1), 
			 (SELECT slug FROM forums WHERE slug = $2),
			 $3,
			 CASE WHEN $4 = '' THEN NULL ELSE $4 END,
			 $5,
			 $6
		 )
		 RETURNING id, forum, create_date`,
		thread.Author,
		thread.Forum,
		thread.Message,
		thread.Slug,
		thread.Title,
		thread.Created,
	).Scan(
		&thread.ID,
		&thread.Forum,
		&thread.Created,
	)

	if err == nil {
		h.db.Exec(
			`UPDATE forums
			 SET threads = threads + 1
			 WHERE slug = $1`,
			thread.Forum,
		)

		h.db.Exec(
			`INSERT INTO forum_users(forum, nickname)
			 VALUES($1, $2)`,
			thread.Forum,
			thread.Author,
		)
	}

	if err, ok := err.(*pq.Error); ok {
		switch err.Code {
		case "23502": // User not found 404
			return t, errors.ErrUserOrForumNotFound
		case "23505": // Already exist
			return t, errors.ErrThreadExist
		}
	}

	return thread, err
}

func (h *ThreadRepository)GetThread(slugOrId string) (t models.Thread, e error) {
	condition := ""
	if _, err := strconv.Atoi(slugOrId); err != nil {
		condition = " slug = $1 "
	} else {
		condition = " id = $1 "
	}

	err := h.db.QueryRow(
		`SELECT author, 
				create_date, 
				forum, 
				id, 
				msg, 
				CASE WHEN slug IS NULL THEN '' ELSE slug END, 
				title, 
				votes
		 FROM threads
		 WHERE`+condition,
		slugOrId,
	).Scan(
		&t.Author,
		&t.Created,
		&t.Forum,
		&t.ID,
		&t.Message,
		&t.Slug,
		&t.Title,
		&t.Votes,
	)

	if err == sql.ErrNoRows {
		return t, errors.ErrThreadNotFound
	}

	return t, err
}

func (h *ThreadRepository)GetThreads(slug, limit, since, desc string) (ts models.Threads, e error) {
	ts = make(models.Threads, 0)
	var sort string
	if desc == "ASC" {
		sort = ">="
		if since == "" {
			since = "1980-03-01 00:00:00-06"
		}
	} else {
		sort = "<="
		if since == "" {
			since = "2040-03-01 00:00:00-06"
		}
	}

	rows, err := h.db.Query(
		fmt.Sprintf(
			`SELECT author, 
					create_date, 
					forum, 
					id, 
					msg, 
					CASE WHEN slug IS NULL THEN '' ELSE slug END, 
					title, 
					votes 
			 FROM threads
			 WHERE forum = $1 AND create_date %s $2
			 ORDER BY create_date %s
			 LIMIT $3`, sort, desc,
		),
		slug,
		since,
		limit,
	)

	if err != nil {
		return ts, err
	}

	i := 0
	for rows.Next() {
		t := models.Thread{}
		err = rows.Scan(
			&t.Author,
			&t.Created,
			&t.Forum,
			&t.ID,
			&t.Message,
			&t.Slug,
			&t.Title,
			&t.Votes,
		)

		ts = append(ts, &t)
		i++
	}
	rows.Close()

	if i == 0 {
		return ts, errors.ErrForumNotFound
	}

	return ts, nil
}

func (h *ThreadRepository)EditThread(slugOrId string, thread models.Thread) (t models.Thread, e error) {
	condition := ""
	if _, err := strconv.Atoi(slugOrId); err != nil {
		condition = " slug = $3 "
	} else {
		condition = " id = $3 "
	}

	err := h.db.QueryRow(
		`UPDATE threads
		 SET msg = (CASE WHEN $1 = '' THEN msg ELSE $1 END), 
			 title = (CASE WHEN $2 = '' THEN title ELSE $2 END)
		 WHERE`+condition+
		`RETURNING author, 
					create_date, 
					forum, 
					id, 
					msg, 
					CASE WHEN slug IS NULL THEN '' ELSE slug END, 
					title, 
					votes`,
		thread.Message,
		thread.Title,
		slugOrId,
	).Scan(
		&t.Author,
		&t.Created,
		&t.Forum,
		&t.ID,
		&t.Message,
		&t.Slug,
		&t.Title,
		&t.Votes,
	)

	if err == sql.ErrNoRows {
		return t, errors.ErrThreadNotFound
	}

	return t, err
}
