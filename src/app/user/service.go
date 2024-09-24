package user

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"regexp"
	"test-edot/constants"
	"test-edot/src/dto"
	"test-edot/src/factory"
	"test-edot/src/models"
	"test-edot/src/repository"
	"test-edot/util"
	"time"
)

type Service interface {
	Register(ctx context.Context, user dto.RegisterUser) error
	Login(ctx context.Context, payload dto.LoginUser) (string, error)
}

type service struct {
	Log            *zap.Logger
	UserRepository repository.UserRepositoryInterface
}

func NewService(f *factory.Factory) Service {
	return &service{
		Log: f.Log,

		UserRepository: f.UserRepository,
	}
}

func (s service) validateRegister(user dto.RegisterUser) error {
	if !constants.MapRoleAvail[user.Role] {
		return constants.RolePayloadInvalid
	}

	reEmail := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !reEmail.MatchString(user.Email) {
		return constants.FormatEmailInvalid
	}

	rePhone := regexp.MustCompile(`^[0-9]{8,15}$`)
	if !rePhone.MatchString(user.Phone) {
		return constants.FormatPhoneInvalid
	}

	return nil
}

func (s service) checkPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (s service) Login(ctx context.Context, payload dto.LoginUser) (string, error) {
	userTrack, err := s.UserRepository.FindOne(ctx, "id,role,email,phone,password", "email = ? or phone = ?", payload.Email, payload.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	if userTrack == (models.User{}) {
		return "", constants.UserNotFound
	}

	if !s.checkPasswordHash(payload.Password, userTrack.Password) {
		return "", constants.InvalidPassword
	}

	token, err := util.GenerateJWT(userTrack)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s service) Register(ctx context.Context, user dto.RegisterUser) error {
	userTrack, err := s.UserRepository.FindOne(ctx, "email,phone", "email = ? or phone = ?", user.Email, user.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if userTrack != (models.User{}) {
		return constants.UserAlreadyInserted
	}

	if err := s.validateRegister(user); err != nil {
		return err
	}

	passwordHashed, err := util.HashPassword(user.Password)
	if err != nil {
		return err
	}

	userData := models.User{
		FullName:  user.FullName,
		Password:  passwordHashed,
		Role:      user.Role,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: time.Now().In(util.LocationTime),
		UpdatedAt: time.Now().In(util.LocationTime),
	}
	if err := s.UserRepository.Create(ctx, &userData); err != nil {
		s.Log.Error("error creating user", zap.Error(err))
		return err
	}

	return nil
}
