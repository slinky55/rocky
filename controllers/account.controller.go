package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/slinky55/rocky/logic"
	"log/slog"
	"net/http"
)

func AccountRegister(c echo.Context) error {
	var data struct {
		UserId    string `json:"userId"`
		PublicKey string `json:"publicKey"`
	}

	err := c.Bind(&data)
	if err != nil {
		slog.Error("error binding data struct: ", err.Error())
		return c.JSON(http.StatusBadRequest, JSON{"error": "malformed json received"})
	}

	if data.UserId == "" {
		slog.Error("phoneNumber is empty")
		return c.JSON(http.StatusBadRequest, JSON{"error": "phoneNumber is empty"})
	}

	phoneNumber, err := logic.FormatPhoneNumber(data.UserId)
	if err != nil {
		slog.Error("phone number format error: ", err.Error())
		return c.JSON(http.StatusInternalServerError, JSON{"error": "internal server error"})
	}

	if logic.UserExists(phoneNumber) {
		slog.Error("user already exists: ", data.UserId)
		return c.JSON(http.StatusBadRequest, JSON{"error": "user already exists"})
	}

	err = logic.CreateUser(phoneNumber)
	if err != nil {
		slog.Error("error while creating user: ", err)
		return c.JSON(http.StatusInternalServerError, JSON{"error": "error while creating user"})
	}

	err = logic.SendPhoneVerification(phoneNumber)
	if err != nil {
		slog.Error("error while verifying phone number: ", err)
		return c.JSON(http.StatusInternalServerError, JSON{"error": "error while verifying phone number"})
	}

	err = logic.SetUserPublicKey(phoneNumber, data.PublicKey)
	if err != nil {
		slog.Error("error while setting user public key: ", err)
		return c.JSON(http.StatusInternalServerError, JSON{"error": "error while setting user public key"})
	}

	return c.JSON(http.StatusOK, JSON{"message": "account created, please verify phone number"})
}

func AccountVerify(c echo.Context) error {
	var data struct {
		UserId string `json:"userId"`
		Code   string `json:"code"`
	}

	err := c.Bind(&data)
	if err != nil {
		slog.Error("error binding data struct: ", err.Error())
		return c.JSON(http.StatusBadRequest, JSON{"error": "malformed json received"})
	}

	if data.UserId == "" {
		slog.Error("phoneNumber is empty")
		return c.JSON(http.StatusBadRequest, JSON{"error": "phoneNumber is empty"})
	}

	if data.Code == "" {
		slog.Error("code is empty")
		return c.JSON(http.StatusBadRequest, JSON{"error": "code is empty"})
	}

	phoneNumber, err := logic.FormatPhoneNumber(data.UserId)
	if err != nil {
		slog.Error("phone number format error: ", err.Error())
		return c.JSON(http.StatusInternalServerError, JSON{"error": "internal server error"})
	}

	err = logic.CheckPhoneVerification(data.Code, phoneNumber)
	if err != nil {
		slog.Error("error while verifying phone number: ", err.Error())
		return c.JSON(http.StatusInternalServerError, JSON{"error": "error while verifying phone number"})
	}

	err = logic.VerifyUser(phoneNumber)
	if err != nil {
		slog.Error("error while verifying phone number: ", err.Error())
		return c.JSON(http.StatusInternalServerError, JSON{"error": "error while verifying phone number"})
	}

	return c.JSON(http.StatusOK, JSON{"message": "account verified"})
}
