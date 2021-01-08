package repository

import (
	"github.com/AleksMVP/sybd/models"
	"github.com/AleksMVP/sybd/pkg/errors"
	"database/sql"
)

type VoteRepository struct {
	db *sql.DB
}

func NewVoteRepository(db *sql.DB) VoteRepository {
	return VoteRepository{
		db: db,
	}
}

func (h *VoteRepository)checkVote(thread int32, nickname string) (vote int32, e error) {
	err := h.db.QueryRow(
		`SELECT voice
		 FROM votes
		 WHERE thread = $1 AND nickname = $2`,
		thread,
		nickname,
	).Scan(&vote)

	if err == sql.ErrNoRows {
		return vote, errors.ErrVoteNotFound
	}

	return vote, nil
}

func (h *VoteRepository)checkUser(nickname string) (err error) {
	err = h.db.QueryRow(
		`SELECT nickname
		 FROM users
		 WHERE nickname = $1`,
		nickname,
	).Scan(&nickname)

	if err == sql.ErrNoRows {
		return errors.ErrUserNotFound
	}

	return err
}

func (h *VoteRepository)VoteForThread(t models.Thread, vote models.Vote) (thread models.Thread, e error) {
	thread = t

	err := h.checkUser(vote.Nickname)
	if err != nil {
		return thread, errors.ErrUserNotFound
	}

	v, err := h.checkVote(thread.ID, vote.Nickname)
	if err != nil {
		_, err = h.db.Exec(
			`INSERT INTO votes(thread, voice, nickname)
			 VALUES($1, $2, $3)`,
			thread.ID,
			vote.Voice,
			vote.Nickname,
		)

		_, err = h.db.Exec(
			`UPDATE threads
			 SET votes = votes + $1
			 WHERE id = $2`,
			vote.Voice,
			thread.ID,
		)

		if vote.Voice < 0 {
			thread.Votes--
		} else {
			thread.Votes++
		}
	} else if v != vote.Voice {
		_, err = h.db.Exec(
			`UPDATE votes
			 SET voice = $1
			 WHERE thread = $2 AND nickname = $3`,
			vote.Voice,
			thread.ID,
			vote.Nickname,
		)

		_, err = h.db.Exec(
			`UPDATE threads
			 SET votes = votes + CASE WHEN $1 < 0 THEN -2 ELSE 2 END
			 WHERE id = $2`,
			vote.Voice,
			thread.ID,
		)
		if vote.Voice < 0 {
			thread.Votes -= 2
		} else {
			thread.Votes += 2
		}
	}

	return thread, err
}

