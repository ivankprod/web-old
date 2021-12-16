package services

import (
	"strconv"

	"ivankprod.ru/src/server/internal/domain"
	"ivankprod.ru/src/server/internal/repositories"
	"ivankprod.ru/src/server/pkg/utils"
)

type UserService interface {
	FindOne(uID uint64, restricted bool) (*domain.User, error)
	IsAuthenticated(uAuth string, uAgent string) (*domain.User, error)
}

type userService struct {
	repository repositories.UserRepository
}

func NewUserService(r repositories.UserRepository) UserService {
	return &userService{
		repository: r,
	}
}

func (s *userService) FindOne(uID uint64, restricted bool) (*domain.User, error) {
	result, err := s.repository.FindOne(uID)

	if result != nil && restricted {
		result.AccessToken = ""
	}

	return result, err
}

func (s *userService) IsAuthenticated(uAuth string, uAgent string) (*domain.User, error) {
	uAuthParsed, err := domain.NewUserAuthFromString(uAuth)
	if err != nil {
		return nil, err
	}
	if uAuthParsed == nil {
		return nil, nil
	}

	result, err := s.repository.FindOne(uAuthParsed.ID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	uSessionHash := utils.HashSHA512(strconv.FormatUint(result.ID, 10) + result.SocialID + result.AccessToken + uAgent)
	if uSessionHash == uAuthParsed.Hash {
		result.AccessToken = "<restricted>"
		result.LastAccessTime = utils.TimeMSK_ToLocaleString()

		return result, nil
	}

	return nil, nil
}
