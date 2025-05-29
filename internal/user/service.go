package user

import (
	"context"
	// "database/sql"
	"time"

	"github.com/google/uuid"
	// "github.com/jmoiron/sqlx"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, u *User) error {
	u.ID = uuid.New()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	return s.repo.InsertUser(ctx, u)
}

// func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
// 	var user User
// 	err := s.DB.GetContext(ctx, &user, `SELECT * FROM users WHERE email = $1`, email)
// 	if err == sql.ErrNoRows {
// 		return nil, nil
// 	}
// 	return &user, err
// }
