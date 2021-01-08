package repository

import (
	"database/sql"
	"github.com/AleksMVP/sybd/models"
	"github.com/AleksMVP/sybd/pkg/errors"
	"github.com/go-openapi/strfmt"
	"context"
	"strings"
	"time"
	"fmt"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return PostRepository{
		db: db,
	}
}

func (h *PostRepository)getParent(threadId int32, id int64) (e error) {
	err := h.db.QueryRow(
		`SELECT id 
		 FROM posts 
		 WHERE thread = $1 AND id = $2`,
		threadId,
		id,
	).Scan(&id)

	if err == sql.ErrNoRows {
		return errors.ErrPostNotFound
	}

	return nil
}

func (h *PostRepository)getAuthor(author string) (e error) {
	err := h.db.QueryRow(
		`SELECT nickname
		 FROM users
		 WHERE nickname = $1`,
		author,
	).Scan(&author)

	if err == sql.ErrNoRows {
		return errors.ErrUserNotFound
	}

	return nil
}

func (h *PostRepository)getPost(id string) (post models.Post, e error) {
	err := h.db.QueryRow(
		`SELECT author, create_date, forum, id, is_edited, msg, parent, thread
		 FROM posts
		 WHERE id = $1`,
		id,
	).Scan(
		&post.Author,
		&post.Created,
		&post.Forum,
		&post.ID,
		&post.IsEdited,
		&post.Message,
		&post.Parent,
		&post.Thread,
	)

	if err == sql.ErrNoRows {
		return post, errors.ErrPostNotFound
	}

	return post, err
}

func (h *PostRepository)CreateNewPosts(posts models.Posts, thread models.Thread) (p models.Posts, e error) {
	template := `(
		'%s', 
		'%s', 
		'%s', 
		'%s', 
		%d, 
		%d, 
		array_append((SELECT path FROM posts WHERE id = %d), (SELECT last_value FROM posts_id_seq))
	)`

	templateUsers := "('%s', '%s')"

	query := "INSERT INTO posts(author, create_date, forum, msg, parent, thread, path) VALUES"
	queryUsers := "INSERT INTO forum_users(forum, nickname) VALUES"

	postsLastIndex := len(posts)

	if postsLastIndex == 0 {
		return p, nil
	}

	time := strfmt.DateTime(time.Now())
	for index, val := range posts {
		val.Created = &time
		val.Thread = thread.ID
		val.Forum = thread.Forum

		e = h.getParent(thread.ID, val.Parent)
		if val.Parent != 0 && e != nil {
			return p, e
		}

		e = h.getAuthor(val.Author)
		if e != nil {
			return p, e
		}

		query += fmt.Sprintf(template, val.Author, val.Created, val.Forum, val.Message, val.Parent, val.Thread, val.Parent)
		queryUsers += fmt.Sprintf(templateUsers, val.Forum, val.Author)

		if index != postsLastIndex - 1 {
			query += ","
			queryUsers += ","
		}
	}
	query += "\nRETURNING id"
	queryUsers += "\nON CONFLICT DO NOTHING"

	ctx := context.Background()
	tx, err := h.db.BeginTx(ctx, nil)

	if err != nil {
		return p, err
	}

	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return p, err
	}

	var i int = 0
	for rows.Next() {
		err = rows.Scan(&posts[i].ID)
		if err != nil  {
			return p, err
		}
		i++	
	}
	err = rows.Err()
	if err != nil {
		return p, err
	}

	tx.ExecContext(
		ctx,
		`UPDATE forums
		 SET posts = posts + $1
		 WHERE slug = $2`,
		len(posts),
		thread.Forum,
	)

	if _, err := tx.ExecContext(ctx, queryUsers); err != nil {
		return p, err
	}

	tx.Commit()

	return posts, err
}

func (h *PostRepository)getFlastSortRows(threadId int32, sort, since, limit, desc string) (rows *sql.Rows, err error) {
	if since == "" {
		rows, err = h.db.Query(
			fmt.Sprintf(
			`SELECT author, create_date, forum, id, is_edited, msg, parent, thread
			 FROM posts
			 WHERE thread = $1
			 ORDER BY id %s
			 LIMIT $2`, desc,
			),
			threadId,
			limit,
		)
	} else {
		rows, err = h.db.Query(
			fmt.Sprintf(
			`SELECT author, create_date, forum, id, is_edited, msg, parent, thread
			 FROM posts
			 WHERE thread = $1 AND id %s $2
			 ORDER BY id %s
			 LIMIT $3`, sort, desc,
			),
			threadId,
			since,
			limit,
		)
	}

	return rows, err
}

func (h *PostRepository)getTreeSortRows(threadId int32, sort, since, limit, desc string) (rows *sql.Rows, err error) {
	if since == "" {
		rows, err = h.db.Query(
			fmt.Sprintf(
				`SELECT author, create_date, forum, id, is_edited, msg, parent, thread
				 FROM posts
				 WHERE thread = $1
				 ORDER BY path %s
				 LIMIT $2`, desc,
			),
			threadId,
			limit,
		)
	} else {
		rows, err = h.db.Query(
			fmt.Sprintf(
				`SELECT author, create_date, forum, id, is_edited, msg, parent, thread
				 FROM posts
				 WHERE thread = $1 AND path %s (SELECT path FROM posts WHERE id = $2)
				 ORDER BY path %s
				 LIMIT $3`, sort, desc,
			),
			threadId,
			since,
			limit,
		)
	}

	return rows, err
}

func (h *PostRepository)getPTreeSortRows(threadId int32, sort, since, limit, desc string) (rows *sql.Rows, err error) {
	if since == "" {
		rows, err = h.db.Query(
			fmt.Sprintf(
				`SELECT author, create_date, forum, id, is_edited, msg, parent, thread
				 FROM posts
				 WHERE thread = $1 and path[1] IN (
					 SELECT path[1]
					 FROM posts
					 WHERE thread = $1 AND array_length(path, 1) = 1
					 ORDER BY path %s
					 LIMIT $2
				 )
				 ORDER BY path[1] %s, path[2:]`, desc, desc,
			),
			threadId,
			limit,
		)
	} else {
		rows, err = h.db.Query(
			fmt.Sprintf(
				`SELECT author, create_date, forum, id, is_edited, msg, parent, thread
				 FROM posts
				 WHERE thread = $1 and path[1] IN (
					 SELECT path[1]
					 FROM posts
					 WHERE thread = $1 AND array_length(path, 1) = 1 AND path[1] %s (SELECT path[1] FROM posts WHERE id = $2)
					 ORDER BY path %s
					 LIMIT $3
				 )
				 ORDER BY path[1] %s, path[2:]`, sort, desc, desc,
			),
			threadId,
			since,
			limit,
		)
	}

	return rows, err
}

func (h *PostRepository)GetPosts(limit, since, desc, sort string, thread models.Thread) (p models.Posts, e error) {
	p = make(models.Posts, 0)

	var sortOperator string
	if desc == "ASC" {
		sortOperator = ">"
	} else {
		sortOperator = "<"
	}

	var rows *sql.Rows

	switch sort {
	case "flat", "":
		rows, e = h.getFlastSortRows(thread.ID, sortOperator, since, limit, desc)
	case "tree":
		rows, e = h.getTreeSortRows(thread.ID, sortOperator, since, limit, desc)
	case "parent_tree":
		rows, e = h.getPTreeSortRows(thread.ID, sortOperator, since, limit, desc)
	}

	if e != nil {
		return p, e
	}

	for rows.Next() {
		ps := models.Post{}
		rows.Scan(
			&ps.Author,
			&ps.Created,
			&ps.Forum,
			&ps.ID,
			&ps.IsEdited,
			&ps.Message,
			&ps.Parent,
			&ps.Thread,
		)

		p = append(p, &ps)
	}
	rows.Close()

	return p, nil
}

func (h *PostRepository)UpdatePost(id string, post models.PostUpdate) (p models.Post, e error) {
	err := h.db.QueryRow(
		`UPDATE posts
		 SET is_edited = (CASE WHEN $1 = '' THEN false WHEN msg = $1 THEN false ELSE true END),
		 	 msg = (CASE WHEN $1 = '' THEN msg ELSE $1 END)
		 WHERE id = $2
		 RETURNING author, create_date, forum, id, is_edited, msg, parent, thread`,
		post.Message,
		id,
	).Scan(
		&p.Author,
		&p.Created,
		&p.Forum,
		&p.ID,
		&p.IsEdited,
		&p.Message,
		&p.Parent,
		&p.Thread,
	)

	if err == sql.ErrNoRows {
		return p, errors.ErrPostNotFound
	}

	return p, err
}

func (h *PostRepository)GetPostFull(id string, item string) (p models.PostFull, e error) {
	post, err := h.getPost(id)
	if err != nil {
		return p, err
	}

	p.Post = &post

	if strings.Contains(item, "user") {
		user := models.User{}
		h.db.QueryRow(
			`SELECT nickname, fullname, about, email
			 FROM users
			 WHERE nickname = $1`,
			post.Author,
		).Scan(
			&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Email,
		)

		p.Author = &user
	}

	if strings.Contains(item, "forum") {
		forum := models.Forum{}
		h.db.QueryRow(
			`SELECT title, nickname, posts, threads, slug
			 FROM forums
			 WHERE slug = $1`,
			post.Forum,
		).Scan(
			&forum.Title,
			&forum.User,
			&forum.Posts,
			&forum.Threads,
			&forum.Slug,
		)

		p.Forum = &forum
	}
	if strings.Contains(item, "thread") {
		thread := &models.Thread{}
		err = h.db.QueryRow(
			`SELECT author, create_date, forum, id, msg, (CASE WHEN slug IS NULL THEN '' ELSE slug END), title, votes
			 FROM threads
			 WHERE id = $1`,
			post.Thread,
		).Scan(
			&thread.Author,
			&thread.Created,
			&thread.Forum,
			&thread.ID,
			&thread.Message,
			&thread.Slug,
			&thread.Title,
			&thread.Votes,
		)
		p.Thread = thread
	}

	return p, err
}

