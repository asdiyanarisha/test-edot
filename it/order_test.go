package it

import (
	"bytes"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"test-edot/src/dto"
	"test-edot/src/models"
	"test-edot/util"
)

func (s *e2eTestSuite) createOrder(user models.UserWithJwt, qty int, product dto.ProductResponse) (int, error) {
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
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+user.Jwt)

	res, err := util.ReqDo(req)
	if err != nil {
		s.Log.Error("Error do req", zap.Error(err))
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != 201 {
		var errRes dto.ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
			s.Log.Error("erorr get decode", zap.Error(err))
			return res.StatusCode, err
		}

		s.Log.Error("error add product", zap.String("status", res.Status), zap.Any("res", errRes))
		return res.StatusCode, errors.New(errRes.Error)
	}

	return res.StatusCode, nil
}
