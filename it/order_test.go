package it

import (
	"bytes"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"strconv"
	"test-edot/src/dto"
	"test-edot/src/models"
	"test-edot/util"
)

type (
	ResponseCreateOrder struct {
		Data models.Order
	}
)

func (s *e2eTestSuite) paymentOrder(orderId int, jwt string) (int, error) {

	url := s.baseUrl + "/api/orders/" + strconv.Itoa(orderId) + "/payment"

	req, err := util.Req("PUT", url, nil)
	if err != nil {
		s.Log.Error("Error creating request", zap.Error(err))
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+jwt)

	res, err := util.ReqDo(req)
	if err != nil {
		s.Log.Error("Error do req", zap.Error(err))
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		var errRes dto.ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
			s.Log.Error("erorr get decode", zap.Error(err))
			return res.StatusCode, err
		}

		s.Log.Error("error payment", zap.String("status", res.Status), zap.Any("res", errRes))
		return res.StatusCode, errors.New(errRes.Error)
	}

	return res.StatusCode, nil

}
func (s *e2eTestSuite) createOrder(user models.UserWithJwt, qty int, product dto.ProductResponse) (int, models.Order, error) {
	var (
		items        []dto.PayloadCreateOrderItems
		payloadOrder dto.PayloadCreateOrder
	)
	url := s.baseUrl + "/api/orders"
	items = append(items, dto.PayloadCreateOrderItems{
		ProductId: product.Id,
		Qty:       qty,
	})

	payloadOrder = dto.PayloadCreateOrder{Items: items}
	payloadByte, _ := json.Marshal(payloadOrder)

	req, err := util.Req("POST", url, bytes.NewBuffer(payloadByte))
	if err != nil {
		s.Log.Error("Error creating request", zap.Error(err))
		return 0, models.Order{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+user.Jwt)

	res, err := util.ReqDo(req)
	if err != nil {
		s.Log.Error("Error do req", zap.Error(err))
		return 0, models.Order{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 201 {
		var errRes dto.ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
			s.Log.Error("erorr get decode", zap.Error(err))
			return res.StatusCode, models.Order{}, err
		}

		s.Log.Error("error add product", zap.String("status", res.Status), zap.Any("res", errRes))
		return res.StatusCode, models.Order{}, errors.New(errRes.Error)
	}

	if res.StatusCode != 201 {
		var errRes dto.ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
			s.Log.Error("erorr get decode", zap.Error(err))
			return res.StatusCode, models.Order{}, err
		}

		s.Log.Error("error add product", zap.String("status", res.Status), zap.Any("res", errRes))
		return res.StatusCode, models.Order{}, errors.New(errRes.Error)
	}

	var response ResponseCreateOrder
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		s.Log.Error("erorr get decode", zap.Error(err))
		return res.StatusCode, models.Order{}, err
	}

	return res.StatusCode, response.Data, nil
}
