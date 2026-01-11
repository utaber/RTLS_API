package user

import (
	"RTLS_API/pkg/models"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type FirebaseDB interface {
	Get(ctx context.Context, path string, dest interface{}) error
}

type Service struct {
	ctx context.Context
	db  FirebaseDB
}

func NewService(ctx context.Context, db FirebaseDB) *Service {
	return &Service{
		ctx: ctx,
		db:  db,
	}
}

func (s *Service) AuthenticateByEmail(email, password string) (string, error) {
	var users map[string]models.LoginRequest

	err := s.db.Get(s.ctx, "Users", &users)
	if err != nil {
		return "", err
	}

	for username, user := range users {
		if user.Email == email {
			err := bcrypt.CompareHashAndPassword(
				[]byte(user.Password),
				[]byte(password),
			)

			if err == nil {
				return username, nil
			}
		}
	}

	return "", errors.New("invalid credentials")
}
