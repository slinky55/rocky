package controllers

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/labstack/echo/v4"
	"github.com/slinky55/rocky/logic"
	"log/slog"
	"net/http"
)

var challenges = make(map[string][]byte)
var sessions = make(map[string]string)

func AuthChallenge(c echo.Context) error {
	var data struct {
		UserId string `json:"userId"`
	}

	err := c.Bind(&data)
	if err != nil {
		slog.Error("error binding json", "error", err.Error())
		return c.JSON(http.StatusBadRequest, JSON{
			"error": "malformed json",
		})
	}

	if data.UserId == "" {
		slog.Error("userId is empty")
		return c.JSON(http.StatusBadRequest, JSON{
			"error": "userId is empty",
		})
	}

	phoneNumber, err := logic.FormatPhoneNumber(data.UserId)
	if err != nil {
		slog.Error("FormatPhoneNumber", "error", err.Error())
		return c.JSON(http.StatusInternalServerError, JSON{
			"error": "malformed phone number",
		})
	}

	challenge, err := logic.GenerateChallenge()
	if err != nil {
		slog.Error("GenerateChallenge", "error", err.Error())
		return c.JSON(http.StatusInternalServerError, JSON{
			"error": "internal server error",
		})
	}

	challenges[phoneNumber] = challenge

	return c.JSON(http.StatusOK, JSON{
		"challenge": base64.StdEncoding.EncodeToString(challenge),
	})
}

func AuthLogin(c echo.Context) error {
	var data struct {
		UserId    string `json:"userId"`
		Signature string `json:"signature"`
	}

	err := c.Bind(&data)
	if err != nil {
		slog.Error("error binding json: ", err.Error())
		return c.JSON(http.StatusBadRequest, JSON{
			"error": "malformed json",
		})
	}

	if data.UserId == "" || data.Signature == "" {
		slog.Error("userId or signedChallenge is empty")
		return c.JSON(http.StatusBadRequest, JSON{
			"error": "userId or signedChallenge is empty",
		})
	}

	signature, err := base64.StdEncoding.DecodeString(data.Signature)
	if err != nil {
		slog.Error("decoding signature failed", "error", err.Error())
		return c.JSON(http.StatusBadRequest, JSON{
			"error": "internal server error",
		})
	}

	phoneNumber, err := logic.FormatPhoneNumber(data.UserId)
	if err != nil {
		slog.Error("FormatPhoneNumber", "error", err.Error())
		return c.JSON(http.StatusInternalServerError, JSON{
			"error": "malformed phone number",
		})
	}

	publicKey, err := logic.GetUserPublicKey(phoneNumber)
	if err != nil {
		slog.Error("GetUserPublicKey", "error", err.Error())
		return c.JSON(http.StatusInternalServerError, JSON{
			"error": "internal server error",
		})
	}

	challenge := challenges[phoneNumber]
	hashed := sha256.Sum256(challenge)

	err = logic.VerifyChallenge(hashed[:], signature, publicKey)
	if err != nil {
		slog.Error("VerifyChallenge", "error", err.Error())
		return c.JSON(http.StatusUnauthorized, JSON{
			"error": "unauthorized",
		})
	}

	tkn, err := logic.GenerateSessionToken()
	if err != nil {
		slog.Error("GenerateSessionToken", "error", err.Error())
		return c.JSON(http.StatusInternalServerError, JSON{
			"error": "internal server error",
		})
	}

	sessions[phoneNumber] = tkn
	delete(challenges, phoneNumber)

	return c.JSON(http.StatusOK, JSON{
		"token": tkn,
	})
}
