package application

import (
	"context"
	"errors"

	"github.com/Mluciano1401/go-api/internal/domain"
)

type AuthService struct {
	repo     UserRepository
	hasher   PasswordHasher
	tokenGen TokenGenerator
}

func NewAuthService(repo UserRepository, hasher PasswordHasher, tokenGen TokenGenerator) *AuthService {
	return &AuthService{repo: repo, hasher: hasher, tokenGen: tokenGen}
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (*domain.User, error) {

	existing, err := s.repo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, domain.ErrEmailAlreadyExists
	}

	hashed, err := s.hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	user := domain.NewUser(name, email, hashed)
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		return "", err
	}

	if err := s.hasher.Compare(user.Password, password); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	return s.tokenGen.Generate(user.ID)
}
