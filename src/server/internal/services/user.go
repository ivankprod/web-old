package services

import (
	"strconv"

	"github.com/ivankprod/ivankprod.ru/src/server/internal/domain"
	"github.com/ivankprod/ivankprod.ru/src/server/internal/repositories"
	"github.com/ivankprod/ivankprod.ru/src/server/pkg/utils"
)

type UserService interface {
	Create(uDTO *domain.UserCreateDTO) (*domain.User, error)
	FindOne(uID uint64, restricted bool) (*domain.User, error)
	FindOneBySocialID(uDTO *domain.UserFindOneBySocialIDDTO) (*domain.User, error)
	FindGroup(uGroup uint64) (*domain.Users, error)
	FindAll(uDTO *domain.UserFindAllDTO) (*domain.Users, error)
	Update(uID uint64, uDTO *domain.UserUpdateDTO) (*domain.User, error)
	UpdateLastAccessTime(uID uint64) (*domain.User, error)
	SignIn(uID uint64, uDTO *domain.UserSignInDTO) (*domain.User, error)
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

func (s *userService) Create(uDTO *domain.UserCreateDTO) (*domain.User, error) {
	result, err := s.repository.Create(uDTO)

	if result != nil {
		result.AccessToken = ""
	}

	return result, err
}

func (s *userService) FindOne(uID uint64, restricted bool) (*domain.User, error) {
	result, err := s.repository.FindOne(uID)

	if result != nil && restricted {
		result.AccessToken = ""
	}

	return result, err
}

func (s *userService) FindOneBySocialID(uDTO *domain.UserFindOneBySocialIDDTO) (*domain.User, error) {
	result, err := s.repository.FindOneBySocialID(uDTO)

	if result != nil {
		result.AccessToken = ""
	}

	return result, err
}

func (s *userService) FindGroup(uGroup uint64) (*domain.Users, error) {
	result, err := s.repository.FindGroup(uGroup)

	if result != nil {
		for i := range *result {
			(*result)[i].AccessToken = ""
		}
	}

	return result, err
}

func (s *userService) FindAll(uDTO *domain.UserFindAllDTO) (*domain.Users, error) {
	result, err := s.repository.FindAll(uDTO)

	if result != nil {
		for i := range *result {
			(*result)[i].AccessToken = ""
		}
	}

	return result, err
}

func (s *userService) Update(uID uint64, uDTO *domain.UserUpdateDTO) (*domain.User, error) {
	result, err := s.repository.Update(uID, uDTO)

	if result != nil {
		result.AccessToken = ""
	}

	return result, err
}

func (s *userService) UpdateLastAccessTime(uID uint64) (*domain.User, error) {
	now := utils.TimeMSK_ToLocaleString()

	result, err := s.repository.Update(uID, &domain.UserUpdateDTO{
		LastAccess: &now,
	})

	if result != nil {
		result.AccessToken = ""
	}

	return result, err
}

func (s *userService) SignIn(uID uint64, uDTO *domain.UserSignInDTO) (*domain.User, error) {
	now := utils.TimeMSK_ToLocaleString()

	result, err := s.repository.Update(uID, &domain.UserUpdateDTO{
		NameFirst:   &uDTO.NameFirst,
		NameLast:    &uDTO.NameLast,
		AvatarPath:  &uDTO.AvatarPath,
		Email:       &uDTO.Email,
		AccessToken: &uDTO.AccessToken,
		LastAccess:  &now,
	})

	if result != nil {
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
		result.AccessToken = ""
		result.LastAccessTime = utils.TimeMSK_ToLocaleString()

		return result, nil
	}

	return nil, nil
}
