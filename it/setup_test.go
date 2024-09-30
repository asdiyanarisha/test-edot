package it

import (
	"bytes"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"test-edot/constants"
	"test-edot/src/dto"
	"test-edot/src/models"
	"test-edot/util"
)

type (
	ResponseCreateShop struct {
		Data dto.ResponseCreateShop
	}
	ResponseRegisterUser struct {
		Data models.User
	}
	ResponseAddWarehouse struct {
		Data dto.ResponseWarehouse
	}
)

func (s *e2eTestSuite) registerAdminShop() {
	s.adminShop = models.UserWithJwt{
		FullName: "admin shop",
		Password: "adminshop123",
		Role:     constants.ROLE_ADMIN_SHOP,
		Email:    "adminshop@example.com",
		Phone:    "8382093562",
	}

	res, err := s.registerUser(s.adminShop)
	if err != nil {
		return
	}
	s.Log.Info("finish register user", zap.String("user", s.adminShop.FullName))

	token, err := s.loginUser(s.adminShop)
	if err != nil {
		return
	}

	s.adminShop.Jwt = token
	s.adminShop.Id = res.Data.Id
	s.Log.Info("finish login admin shop", zap.String("user", s.adminShop.FullName))
}

func (s *e2eTestSuite) createShop() {
	payload := dto.PayloadCreateShop{
		Name:     "Toko Makmur",
		Location: "Jakarta, Indonesia",
	}

	payloadByte, _ := json.Marshal(payload)
	url := s.baseUrl + "/api/shops"
	req, err := util.Req("POST", url, bytes.NewBuffer(payloadByte))
	if err != nil {
		s.Log.Error("Error creating request", zap.Error(err))
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+s.adminShop.Jwt)

	res, err := util.ReqDo(req)
	if err != nil {
		s.Log.Error("Error do req", zap.Error(err))
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 201 {
		var errRes dto.ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
			s.Log.Error("Error decoding response", zap.Error(err))
			return
		}

		if errRes.Error == constants.ShopAlreadyInserted.Error() {
			return
		}

		s.Log.Error("error create shop", zap.String("status", res.Status), zap.Any("res", util.ResponseBodyToString(res)))
		return
	}

	var response ResponseCreateShop
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		s.Log.Error("Error decoding response", zap.Error(err))
		return
	}

	s.shop = models.Shop{
		ID:       response.Data.Id,
		Name:     response.Data.Name,
		Location: response.Data.Location,
	}

	s.Log.Info("success create shop", zap.String("status", res.Status))
}

func (s *e2eTestSuite) createWarehouses() {
	warehouses := []models.Warehouse{
		{Name: "Gudang Jakarta", Location: "Jakarta, Indonesia"},
		{Name: "Gudang Bandung", Location: "Bandung, Indonesia"},
	}

	for _, warehouse := range warehouses {
		res, err := s.createWarehouse(warehouse)
		if err != nil {
			s.Log.Error("Error creating warehouse", zap.Error(err))
			return
		}

		s.warehouses = append(s.warehouses, models.Warehouse{
			ID:       res.ID,
			Name:     res.Name,
			Location: res.Location,
			UserId:   res.UserId,
			IsActive: res.IsActive,
		})
	}

	s.Log.Info("finish setup warehouse")
}

func (s *e2eTestSuite) createWarehouse(warehouse models.Warehouse) (dto.ResponseWarehouse, error) {
	reqRegister := dto.PayloadAddWarehouse{
		Name:     warehouse.Name,
		Location: warehouse.Location,
	}

	payloadByte, _ := json.Marshal(reqRegister)
	url := s.baseUrl + "/api/warehouses"
	req, err := util.Req("POST", url, bytes.NewBuffer(payloadByte))
	if err != nil {
		s.Log.Error("Error creating request", zap.Error(err))
		return dto.ResponseWarehouse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+s.adminShop.Jwt)

	res, err := util.ReqDo(req)
	if err != nil {
		s.Log.Error("Error do req", zap.Error(err))
		return dto.ResponseWarehouse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 201 {
		var errRes dto.ErrorResponse

		if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
			return dto.ResponseWarehouse{}, err
		}

		if errRes.Error == constants.WarehouseAlreadyExisted.Error() {
			return dto.ResponseWarehouse{}, nil
		}

		s.Log.Error("error add warehouse", zap.String("status", res.Status), zap.Any("res", util.ResponseBodyToString(res)))
		return dto.ResponseWarehouse{}, errors.New("error add warehouse")
	}

	var response ResponseAddWarehouse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		s.Log.Error("Error decoding response", zap.Error(err))
		return dto.ResponseWarehouse{}, err
	}

	return response.Data, nil
}

func (s *e2eTestSuite) registerUsers() {
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

	for i, user := range s.users {
		res, err := s.registerUser(user)
		if err != nil {
			return
		}
		s.Log.Info("finish register user", zap.String("user", user.FullName))

		token, err := s.loginUser(user)
		if err != nil {
			return
		}

		s.users[i].Jwt = token
		s.users[i].Id = res.Data.Id

		s.Log.Info("finish login user", zap.String("user", user.FullName))
	}
}

func (s *e2eTestSuite) registerUser(user models.UserWithJwt) (ResponseRegisterUser, error) {
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
		return ResponseRegisterUser{}, err
	}

	res, err := util.ReqDo(req)
	if err != nil {
		s.Log.Error("Error do req", zap.Error(err))
		return ResponseRegisterUser{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 201 {
		var errRes dto.ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
			return ResponseRegisterUser{}, err
		}

		if errRes.Error == constants.UserAlreadyInserted.Error() {
			return ResponseRegisterUser{}, nil
		}

		s.Log.Error("error registering user", zap.String("status", res.Status), zap.Any("res", util.ResponseBodyToString(res)))
		return ResponseRegisterUser{}, errors.New("error registering user")
	}

	var response ResponseRegisterUser
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		s.Log.Error("Error decoding response", zap.Error(err))
		return ResponseRegisterUser{}, err
	}

	return response, nil
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
