package repository

import "github.com/rickyhuang08/mini-exchange.git/internal/domain"

type UserInterface interface {
	FindByEmail(email string) (*domain.UserInfo, error)
}

type InMemoryUserRepository struct {
	users []*domain.UserInfo
}

func NewInMemoryUserRepository() UserInterface {
	return &InMemoryUserRepository{
		users: []*domain.UserInfo{
			{ID: 1, Name: "YBTech", Email: "ybtech@example.com", Role: 1, PasswordHash: "$2a$12$IhqMPh48LhSE9T/5vPogK.fFsUUDiR7YWNbmWd42fcpNzKn9dc3N."}, // password: "yntech1234"
			{ID: 2, Name: "Alice", Email: "alice@example.com", Role: 2, PasswordHash: "$2a$12$voe9dcZKBr86mLbsD9JGEOtQWMHTbvLmYKoU9oWJtMrxsEsUfG0BO"}, // password: "alice1234"
			{ID: 3, Name: "Bob", Email: "bob@example.com", Role: 2, PasswordHash: "$2a$12$kXhS0jGB7iDRBr.f5K/B4OZB6LKY95xgFwpG61pFVa.FofH2GBrse"}, // password: "bob1234"
			{ID: 4, Name: "Charlie", Email: "charlie@example.com", Role: 2, PasswordHash: "$2a$12$tGN6ilGwsesVZHPf0bAVz.xQ8cPojT9Q6OT9QmVP755rZAAwPwmK2"}, // password: "charlie1234"
		},
	}
}

func (r *InMemoryUserRepository) FindByEmail(email string) (*domain.UserInfo, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}