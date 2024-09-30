package it

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"net/http"
	"sync"
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

	s.T().Run("Test_Order_MultiUser", func(t *testing.T) {
		s.processTestOrder(products[0])
	})

}

func (s *e2eTestSuite) processTestOrder(product dto.ProductResponse) {
	dataTests := []struct {
		name string
		user models.UserWithJwt
		qty  int
	}{
		{
			name: "Test User 1",
			user: s.users[0],
			qty:  3,
		},
		{
			name: "Test User 3",
			user: s.users[2],
			qty:  12,
		},
		{
			name: "Test User 4",
			user: s.users[3],
			qty:  20,
		},
	}

	var wg sync.WaitGroup
	for _, test := range dataTests {
		wg.Add(1)
		test := test
		go func() {
			defer wg.Done()
			s.Log.Info("run test", zap.String("name", test.name), zap.Int("qty", test.qty))

			statusCode, err := s.createOrder(test.user, test.qty, product)
			s.NoError(err, "create order failed")
			s.Equal(http.StatusCreated, statusCode, "create order not 200 code")
		}()
	}

	wg.Wait()
}
