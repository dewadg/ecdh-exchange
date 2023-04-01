package server

import (
	"errors"
	"net/http"

	"github.com/dewadg/ecdh-exchange/internal/exchange"
	"github.com/gofiber/fiber/v2"
)

type keyExchangeRequest struct {
	PublicKey string `json:"publicKey"`
}

type keyExchangeResponse struct {
	SessionID string `json:"sessionId"`
	PublicKey string `json:"publicKey"`
}

func handleKeyExchange() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody keyExchangeRequest
		if err := c.BodyParser(&reqBody); err != nil {
			return c.Status(http.StatusBadRequest).JSON(err.Error())
		}

		peerPublicKey := []byte(reqBody.PublicKey)
		sessionID, publicKey, err := exchange.Exchange(c.UserContext(), peerPublicKey)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(err.Error())
		}

		return c.JSON(keyExchangeResponse{
			SessionID: sessionID,
			PublicKey: string(publicKey),
		})
	}
}

type getSecretResponse struct {
	Secret string `json:"secret"`
}

func handleGetSecret() fiber.Handler {
	return func(c *fiber.Ctx) error {
		secret, err := exchange.GetSecret(c.UserContext(), c.Params("sessionID"))
		if err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, exchange.ErrSecretNotFound) {
				status = http.StatusNotFound
			}

			return c.Status(status).JSON(err.Error())
		}

		return c.JSON(getSecretResponse{
			Secret: string(secret),
		})
	}
}
