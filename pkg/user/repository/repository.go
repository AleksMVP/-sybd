package repository

import (
	"github.com/AleksMVP/sybd/models"
	"github.com/AleksMVP/sybd/pkg/errors"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{
		db: db,
	}
}

func (h *UserRepository)GetUsersByNickOrEmail(user models.User) (users models.Users) {
	rows, _ := h.db.Query(
		`SELECT nickname, fullname, about, email
		 FROM users
		 WHERE email = $1 OR nickname = $2`,
		user.Email,
		user.Nickname,
	)

	for rows.Next() {
		us := models.User{}
		_ = rows.Scan(
			&us.Nickname,
			&us.Fullname,
			&us.About,
			&us.Email,
		)

		users = append(users, &us)
	}

	rows.Close()

	return users
}

func (h *UserRepository)CreateUser(user models.User) (u models.User, e error) {
	_, err := h.db.Exec(
		`INSERT INTO users(nickname, fullname, about, email)
		 VALUES($1, $2, $3, $4)`,
		user.Nickname,
		user.Fullname,
		user.About,
		user.Email,
	)

	if err, ok := err.(*pq.Error); ok {
		switch err.Code {
		case "23505": // Already exist
			return u, errors.ErrUserOrEmailExist
		}
	}

	return user, err
}

func (h *UserRepository)GetUserProfile(nickname string) (u models.User, e error) {
	err := h.db.QueryRow(
		`SELECT nickname, fullname, about, email
		 FROM users
		 WHERE nickname = $1`,
		nickname,
	).Scan(
		&u.Nickname,
		&u.Fullname,
		&u.About,
		&u.Email,
	)

	if err == sql.ErrNoRows {
		return u, errors.ErrUserNotFound
	}

	return u, err
}

func (h *UserRepository)EditUserProfile(nickname string, user models.UserUpdate) (u models.User, e error) {
	err := h.db.QueryRow(
		`UPDATE users
		 SET fullname = (CASE WHEN $1 = '' THEN fullname ELSE $1 END), 
			 about = (CASE WHEN $2 = '' THEN about ELSE $2 END), 
			 email = (CASE WHEN $3 = '' THEN email ELSE $3 END)
		 WHERE nickname = $4
		 RETURNING nickname, fullname, about, email`,
		user.Fullname,
		user.About,
		user.Email,
		nickname,
	).Scan(
		&u.Nickname,
		&u.Fullname,
		&u.About,
		&u.Email,
	)

	if err == sql.ErrNoRows {
		return u, errors.ErrUserNotFound
	}

	if err, ok := err.(*pq.Error); ok {
		switch err.Code {
		case "23505": // Already exist
			return u, errors.ErrUserExist
		}
	}

	return u, err
}

func (h *UserRepository)GetUsers(slug, limit, since, desc string) (u models.Users, e error) {
	u = make(models.Users, 0)
	var sort string
	if desc == "ASC" {
		sort = ">"
		if since == "" {
			since = ""
		}
	} else {
		sort = "<"
		if since == "" {
			since = "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"
		}
	}

	rows, err := h.db.Query(
		fmt.Sprintf(
			`SELECT u.nickname, u.fullname, u.about, u.email
			 FROM forum_users AS fr INNER JOIN users AS u ON fr.nickname = u.nickname
			 WHERE forum = $1 AND fr.nickname %s $2
			 ORDER BY fr.nickname %s
			 LIMIT $3`, sort, desc,
		),
		slug,
		since,
		limit,
	)

	if err != nil {
		return u, err
	}

	for rows.Next() {
		t := models.User{}
		rows.Scan(
			&t.Nickname,
			&t.Fullname,
			&t.About,
			&t.Email,
		)

		u = append(u, &t)
	}
	rows.Close()

	return u, err
}