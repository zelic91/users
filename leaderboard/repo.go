package leaderboard

import (
	"database/sql"

	"github.com/go-openapi/strfmt"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	DB *sqlx.DB
}

type LeaderboardItem struct {
	ID       int64  `db:"id"`
	UserID   int64  `db:"user_id"`
	Username string `db:"username"`
	Score    int64  `db:"score"`
	Rank     int64  `db:"rank"`
}

type Score struct {
	ID        int64            `db:"id"`
	UserID    int64            `db:"user_id"`
	Score     int64            `db:"score"`
	CreatedAt *strfmt.DateTime `db:"created_at"`
	UpdatedAt *strfmt.DateTime `db:"updated_at"`
}

func NewRepo(db *sqlx.DB) Repo {
	return Repo{
		DB: db,
	}
}

func (r Repo) GetLeaderboard(limit *int32, offset *int32) ([]*LeaderboardItem, error) {
	query := `
	SELECT 
		l.id id, 
		l.user_id user_id, 
		l.score score, 
		u.username username,
		RANK() OVER (
			ORDER BY l.score DESC, l.created_at ASC
		) rank
	FROM leaderboard AS l
	INNER JOIN users AS u ON l.user_id = u.id
	ORDER BY l.score DESC
	`

	if limit != nil {
		query = query + ` LIMIT $1`
	}

	if offset != nil {
		query = query + ` OFFSET $2`
	}

	results := []*LeaderboardItem{}
	err := r.DB.Select(&results, query)

	if err != nil {
		return nil, err
	}

	return results, err
}

func (r Repo) GetScore(user_id int64) (*Score, error) {
	query := `
	SELECT * FROM leaderboard
	WHERE user_id = $1 
	LIMIT 1
	`

	score := Score{}
	err := r.DB.Get(&score, query, user_id)

	if err != nil {
		return nil, err
	}

	return &score, nil
}

func (r Repo) CreateScore(user_id int64, score int64) error {
	query := `
	INSERT INTO leaderboard(
		user_id,
		score
	) VALUES (
		$1,
		$2
	)
	`
	_, err := r.DB.Exec(query, user_id, score)
	return err
}

func (r Repo) UpdateScore(user_id int64, score int64) error {
	query := `
	UPDATE leaderboard
	SET score = $2
	WHERE user_id = $1
	`
	_, err := r.DB.Exec(query, user_id, score)
	return err
}

func (r Repo) SubmitScore(user_id int64, score int64) error {
	currentScore, err := r.GetScore(user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return r.CreateScore(user_id, score)
		} else {
			return err
		}
	}

	// New score must be larger than old score
	if currentScore.Score >= score {
		return nil
	}

	return r.UpdateScore(user_id, score)
}
