package services

import (
	"Perwatch-case/internal/helpers"
	"Perwatch-case/internal/models"
	"Perwatch-case/internal/repositories/database"
	"errors"
)

type UserService struct {
	UserRepository database.UserRepositoryInterface
}

type UserServiceInterface interface {
	Register(request models.UserRegisterRequest) (string, error)
	Login(user models.UserLoginRequest) (string, error)
}

func NewUserService() UserServiceInterface {
	return &UserService{
		UserRepository: database.NewUserRepository(),
	}
}

func (u *UserService) Register(request models.UserRegisterRequest) (string, error) {
	user := models.NewUser()

	firstnameAndLastnameExists, err := u.UserRepository.CheckFirstnameAndLastnameExists(request.Firstname, request.Lastname)
	if err != nil {
		return "", err
	}
	if firstnameAndLastnameExists {
		return "", errors.New("firstname and lastname already exists")
	}

	usernameExists, err := u.UserRepository.CheckUsernameExists(request.Username)
	if err != nil {
		return "", err
	}
	if usernameExists {
		return "", errors.New("username already exists")
	}

	hashedPassword, err := helpers.HashPassword(request.Password)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	user.Firstname = request.Firstname
	user.Lastname = request.Lastname
	user.Username = request.Username
	user.Password = hashedPassword

	token, err := helpers.GenerateJWT(&user)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	if err := u.UserRepository.Create(&user); err != nil {
		return "", errors.New("failed to create user")
	}

	return token, nil
}

func (u *UserService) Login(user models.UserLoginRequest) (string, error) {
	foundUser, err := u.UserRepository.FindByUsername(user.Username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	if !helpers.CheckPasswordHash(foundUser.Password, user.Password) {
		return "", errors.New("invalid username or password")
	}

	token, err := helpers.GenerateJWT(foundUser)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
