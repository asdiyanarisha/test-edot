package it

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"test-edot/src/dto"
	"test-edot/src/models"
	"test-edot/util"
	"testing"
)

type e2eTestSuite struct {
	suite.Suite
	users      []models.UserWithJwt
	adminShop  models.UserWithJwt
	baseUrl    string
	Log        *zap.Logger
	warehouses []models.Warehouse
	shop       models.Shop
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

	s.createWarehouses()

	s.addProduct()
}

func (s *e2eTestSuite) TestOrderProducts() {
	var products []dto.ProductResponse
	s.T().Run("Test_get_product", func(t *testing.T) {
		res := s.getListProduct()
		s.Equal(1, len(res), "get product list failed")
		products = res
	})

	s.T().Run("Test_transfer_stock_by_product", func(t *testing.T) {
		statusCode, err := s.transferProduct(products[0])
		s.NoError(err, "transfer product failed")
		s.Equal(200, statusCode, "transfer product not 200 code")
	})

	fmt.Println("products:", products[0])

}
