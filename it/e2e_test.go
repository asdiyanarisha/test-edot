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
	var orderIds []int
	var jwtMaps []string

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
		s.processTestOrder(products[0], &orderIds, &jwtMaps)
	})

	s.T().Run("Test_Payment", func(t *testing.T) {
		for i, id := range orderIds {
			statusCode, err := s.paymentOrder(id, jwtMaps[i])
			s.NoError(err, "payment order failed")
			s.Equal(200, statusCode, "payment get failed")
		}
	})

	s.T().Run("Test_Stock_After_Payment", func(t *testing.T) {
		product := s.getProductDetail(products[0].Id)

		s.Equal(0, product.ReservedStock, "check product reserved stock failed")
		s.Equal(5, product.Stock, "check product stock failed")
	})

}

func (s *e2eTestSuite) processTestOrder(product dto.ProductResponse, orderIds *[]int, jwtMaps *[]string) {
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

			s.T().Run("Test_Order_"+test.name, func(t *testing.T) {
				statusCode, res, err := s.createOrder(test.user, test.qty, product)
				s.NoError(err, "create order failed")
				s.Equal(http.StatusCreated, statusCode, "create order not 200 code")

				*orderIds = append(*orderIds, res.Id)
				*jwtMaps = append(*jwtMaps, test.user.Jwt)
			})
		}()
	}

	wg.Wait()
}
