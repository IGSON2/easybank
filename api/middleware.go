package api

import (
	"easybank/token"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "autorization_payload"
)

func authMiddleware(maker token.Maker) fiber.Handler {
	abort := func(c *fiber.Ctx, err string) error {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse{fmt.Errorf(err)})
	}

	return func(c *fiber.Ctx) error {
		authorizationHeader := c.Get(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			return abort(c, "authorization header is not provided")
		}

		fields := strings.Fields(authorizationHeader)

		if len(fields) < 2 {
			return abort(c, "invalid authorization header format")
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			return abort(c, "unsupported authorization type,"+authorizationType)
		}

		accessToken := fields[1]
		payload, err := maker.VerifyToken(accessToken)
		if err != nil {
			return abort(c, fmt.Sprintf("%v", err))
		}
		c.Locals(authorizationPayloadKey, payload)
		return c.Next()
	}
}
