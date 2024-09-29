package it

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"test-edot/constants"
	"test-edot/src/dto"
	"test-edot/src/models"
	"test-edot/util"
	"testing"
)

type e2eTestSuite struct {
	suite.Suite
	users     []models.UserWithJwt
	adminShop models.User
	baseUrl   string
	Log       *zap.Logger
}

func TestE2ETestSuite(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	suite.Run(t, &e2eTestSuite{Log: logger})
}

func (s *e2eTestSuite) SetupTest() {
	s.baseUrl = fmt.Sprintf("http://localhost:%s", util.GetEnvTest("APP_PORT", ""))
	fmt.Println(s.baseUrl)
	s.users = []models.UserWithJwt{
		{
			FullName: "userbiasa1",
			Password: "userbiasa1123",
			Role:     constants.ROLE_USER,
			Email:    "usebiasa1@example.com",
			Phone:    "8382093560",
		},
		{
			FullName: "usersuper1",
			Password: "usersuper1123",
			Role:     constants.ROLE_USER,
			Email:    "usebiasa12@example.com",
			Phone:    "8382093561",
		},
	}

	for _, user := range s.users {
		if err := s.registerUser(user); err != nil {
			return
		}
		s.Log.Info("finish register user", zap.String("user", user.FullName))

		token, err := s.loginUser(user)
		if err != nil {
			return
		}
		s.Log.Info("finish login user", zap.String("user", user.FullName), zap.String("token", token))
	}

}

func (s *e2eTestSuite) registerUser(user models.UserWithJwt) error {
	reqRegister := dto.RegisterUser{
		FullName: user.FullName,
		Email:    user.Email,
		Phone:    user.Phone,
		Role:     user.Role,
		Password: user.Password,
	}

	payloadByte, _ := json.Marshal(reqRegister)
	url := s.baseUrl + "/api/users/register"
	req, err := util.Req("POST", url, bytes.NewBuffer(payloadByte))
	if err != nil {
		s.Log.Error("Error creating request", zap.Error(err))
		return err
	}

	res, err := util.ReqDo(req)
	if err != nil {
		s.Log.Error("Error do req", zap.Error(err))
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 201 {
		var errRes dto.ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
			return err
		}

		if errRes.Error == constants.UserAlreadyInserted.Error() {
			return nil
		}

		s.Log.Error("error registering user", zap.String("status", res.Status), zap.Any("res", util.ResponseBodyToString(res)))
		return errors.New("error registering user")
	}

	return nil
}

func (s *e2eTestSuite) loginUser(user models.UserWithJwt) (string, error) {
	reqRegister := dto.LoginUser{
		Email:    user.Email,
		Password: user.Password,
	}

	payloadByte, _ := json.Marshal(reqRegister)
	url := s.baseUrl + "/api/users/login"
	req, err := util.Req("POST", url, bytes.NewBuffer(payloadByte))
	if err != nil {
		s.Log.Error("Error creating request", zap.Error(err))
		return "", err
	}

	res, err := util.ReqDo(req)
	if err != nil {
		s.Log.Error("Error do req", zap.Error(err))
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		s.Log.Error("error login user", zap.String("status", res.Status), zap.Any("res", util.ResponseBodyToString(res)))
		return "", errors.New("error login user")
	}

	var resLogin dto.Response
	if err := json.NewDecoder(res.Body).Decode(&resLogin); err != nil {
		return "", err
	}

	convert := resLogin.Data.(map[string]interface{})

	return convert["token"].(string), nil
}

func (s *e2eTestSuite) Test_EndToEnd_CreateArticle() {
	fmt.Println("Running e2e tests...2")
}
