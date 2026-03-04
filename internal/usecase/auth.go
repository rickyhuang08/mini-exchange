package usecase

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rickyhuang08/mini-exchange.git/helpers"
	"github.com/rickyhuang08/mini-exchange.git/internal/entity"
	"github.com/rickyhuang08/mini-exchange.git/internal/repository"
	pkg_jwt "github.com/rickyhuang08/mini-exchange.git/pkg/jwt"
	"github.com/rickyhuang08/mini-exchange.git/pkg/logger"
)

type AuthUsecase struct {
	UserRepositoryModule repository.UserInterface
	JwtHelper            *pkg_jwt.JWTHelper
	Logger               *logger.Logger
	PublicKeyPath        string
}

func NewAuthUsecase(
	userRepoM repository.UserInterface,
	jwtHelper *pkg_jwt.JWTHelper,
	logger *logger.Logger,
	publicKeyPath string,
) *AuthUsecase {
	return &AuthUsecase{
		UserRepositoryModule: userRepoM,
		JwtHelper:            jwtHelper,
		Logger:               logger,
		PublicKeyPath:        publicKeyPath,
	}
}

func (uc *AuthUsecase) Login(request entity.LoginRequest) (entity.LoginResponse, error) {
	uc.Logger.LogLevel(logger.LogLevelInfo, "Login is Running")

	var loginResponse entity.LoginResponse
	user, err := uc.UserRepositoryModule.FindByEmail(request.Email)
	if user == nil && err == nil {
		err = errors.New("user not found")
		uc.Logger.LogLevel(logger.LogLevelError, fmt.Sprintf("[Login] Find By Email return error : %s", err.Error()))
		return loginResponse, err

	} else if err != nil {
		uc.Logger.LogLevel(logger.LogLevelError, fmt.Sprintf("[Login] Find By Email return error : %s", err.Error()))
		return loginResponse, err
	}

	if err := helpers.CheckPassword(request.Password, user.PasswordHash); err != nil {
		uc.Logger.LogLevel(logger.LogLevelError, fmt.Sprintf("[Login] password check failed: %v", err))
		return loginResponse, err
	}

	token, err := uc.JwtHelper.GenerateJWT(user.ID, user.Role, user.Email)
	if err != nil {
		uc.Logger.LogLevel(logger.LogLevelError, fmt.Sprintf("[Login] Generate JWT return error : %s", err.Error()))
		return loginResponse, err

	}

	userInfo := entity.UserInfo{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	loginResponse.Token = token
	loginResponse.User = userInfo
	uc.Logger.LogLevel(logger.LogLevelInfo, fmt.Sprintf("[Login] Success, Token : %s", token))
	return loginResponse, nil
}

func (uc *AuthUsecase) ParseJWT(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	uc.Logger.LogLevel(logger.LogLevelInfo, "ParseJWT is Running")

	rsaPublicKey, err := pkg_jwt.LoadPublicKey(uc.PublicKeyPath)
	if err != nil {
		err = fmt.Errorf("unable to load public key: %w", err)
		uc.Logger.LogLevel(logger.LogLevelError, err.Error())
		return nil, nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return rsaPublicKey, nil
	})

	if err != nil || !token.Valid {
		err = fmt.Errorf("invalid token: %w", err)
		uc.Logger.LogLevel(logger.LogLevelError, err.Error())
		return nil, nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("invalid claims structure")
		uc.Logger.LogLevel(logger.LogLevelError, err.Error())
		return nil, nil, err
	}

	uc.Logger.LogLevel(logger.LogLevelInfo, fmt.Sprintf("ParseJWT is Success, token : %s, claims : %v", tokenString, claims))
	return token, claims, nil
}