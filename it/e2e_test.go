package it

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"test-edot/src/models"
	"test-edot/util"
	"testing"
)

type e2eTestSuite struct {
	suite.Suite
	users     []models.UserWithJwt
	adminShop models.UserWithJwt
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

	s.registerUsers()
	s.registerAdminShop()

	s.createShop()
}

func (s *e2eTestSuite) Test_EndToEnd_CreateArticle() {
	fmt.Println("Running e2e tests...2")
}
