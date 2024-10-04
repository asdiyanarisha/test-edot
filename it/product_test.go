package it

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"strconv"
	"test-edot/src/dto"
	"test-edot/util"
)

type (
	ResponseListProduct struct {
		Data []dto.ProductResponse
	}

	ResponseDetailProduct struct {
		Data dto.ProductDetailResponse
	}
)

func (s *e2eTestSuite) getListProduct() []dto.ProductResponse {
	url := s.baseUrl + "/api/products"
	req, err := util.Req("GET", url, nil)
	if err != nil {
		s.Log.Error("Error creating request", zap.Error(err))
		return []dto.ProductResponse{}
	}

	res, err := util.ReqDo(req)
	if err != nil {
		s.Log.Error("Error do req", zap.Error(err))
		return []dto.ProductResponse{}
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		var errRes dto.ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
			s.Log.Error("erorr get decode", zap.Error(err))
			return []dto.ProductResponse{}
		}

		s.Log.Error("error add product", zap.String("status", res.Status), zap.Any("res", util.ResponseBodyToString(res)))
		return []dto.ProductResponse{}
	}

	var products ResponseListProduct
	if err := json.NewDecoder(res.Body).Decode(&products); err != nil {
		s.Log.Error("erorr get decode", zap.Error(err))
		return []dto.ProductResponse{}
	}

	return products.Data
}

func (s *e2eTestSuite) getProductDetail(productId int) dto.ProductDetailResponse {
	url := s.baseUrl + "/api/products/" + strconv.Itoa(productId) + "/detail"
	req, err := util.Req("GET", url, nil)
	if err != nil {
		s.Log.Error("Error creating request", zap.Error(err))
		return dto.ProductDetailResponse{}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+s.adminShop.Jwt)

	res, err := util.ReqDo(req)
	if err != nil {
		s.Log.Error("Error do req", zap.Error(err))
		return dto.ProductDetailResponse{}
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		var errRes dto.ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
			s.Log.Error("erorr get decode", zap.Error(err))
			return dto.ProductDetailResponse{}
		}

		s.Log.Error("error get product", zap.String("status", res.Status), zap.Any("res", util.ResponseBodyToString(res)))
		return dto.ProductDetailResponse{}
	}

	var products ResponseDetailProduct
	if err := json.NewDecoder(res.Body).Decode(&products); err != nil {
		s.Log.Error("erorr get decode", zap.Error(err))
		return dto.ProductDetailResponse{}
	}

	return products.Data
}

func (s *e2eTestSuite) transferProduct(product dto.ProductResponse) (int, error) {
	url := s.baseUrl + "/api/products/" + strconv.Itoa(product.Id) + "/transfer"
	payload := dto.TransferProductWarehouse{
		FromWarehouseId: s.warehouses[0].ID,
		ToWarehouseId:   s.warehouses[1].ID,
		Qty:             20,
	}

	payloadByte, _ := json.Marshal(payload)

	req, err := util.Req("POST", url, bytes.NewBuffer(payloadByte))
	if err != nil {
		s.Log.Error("Error creating request", zap.Error(err))
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+s.adminShop.Jwt)

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

		s.Log.Error("error transfer product", zap.String("status", res.Status), zap.Any("res", util.ResponseBodyToString(res)))
		return res.StatusCode, err
	}

	return res.StatusCode, nil
}
