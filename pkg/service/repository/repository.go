package repository

import (
	"database/sql"
	"github.com/AleksMVP/sybd/models"
)

type ServiceRepository struct {
	db *sql.DB
}

func NewServiceRepository(db *sql.DB) ServiceRepository {
	return ServiceRepository{
		db: db,
	}
}

func (h *ServiceRepository)GetStatus() (status models.Status) {
	h.db.QueryRow(
		`SELECT COUNT(*) 
		 FROM users`,
	).Scan(&status.User)
	h.db.QueryRow(
		`SELECT COUNT(*) 
		 FROM forums`,
	).Scan(&status.Forum)
	h.db.QueryRow(
		`SELECT COUNT(*) 
		 FROM threads`,
	).Scan(&status.Thread)
	h.db.QueryRow(
		`SELECT COUNT(*) 
		 FROM posts`,
	).Scan(&status.Post)

	return status
}

func (h *ServiceRepository)Clear() {
	h.db.Exec(
		`TRUNCATE TABLE posts;
		 TRUNCATE TABLE votes;
		 TRUNCATE TABLE forum_users;
		 TRUNCATE TABLE threads CASCADE;
		 TRUNCATE TABLE forums CASCADE;
		 TRUNCATE TABLE users CASCADE;`,
	)
}
