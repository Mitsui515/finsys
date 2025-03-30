package service

import (
	"regexp"

	"github.com/Mitsui515/finsys/config"
	"github.com/Mitsui515/finsys/middleware"
	"github.com/Mitsui515/finsys/model"
	"github.com/Mitsui515/finsys/repository"
	"github.com/Mitsui515/finsys/utils"
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepository: repository.NewUserRepository(config.DB),
	}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *UserService) Register(req *RegisterRequest) (uint, error) {
	if len(req.Username) < 3 || len(req.Username) > 20 {
		return 0, utils.ErrInvalidUsername
	}
	if len(req.Password) < 6 || len(req.Password) > 20 {
		return 0, utils.ErrInvalidPassword
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return 0, utils.ErrInvalidEmail
	}
	_, err := s.userRepository.FindByUsername(req.Username)
	if err == nil {
		return 0, utils.ErrExistedUsername
	}
	_, err = s.userRepository.FindByEmail(req.Email)
	if err == nil {
		return 0, utils.ErrExistedEmail
	}
	user := model.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}
	if err := s.userRepository.Create(&user); err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (s *UserService) Login(req *LoginRequest) (string, error) {
	user, err := s.userRepository.FindByUsername(req.Username)
	if err != nil {
		return "", utils.ErrFalseUsername
	}
	if !user.CheckPassword(req.Password) {
		return "", utils.ErrFalsePassword
	}
	token, err := middleware.GenerateToken(user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *UserService) GetByID(id uint) (*model.User, error) {
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
