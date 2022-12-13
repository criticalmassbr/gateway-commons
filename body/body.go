package body

import (
	"encoding/json"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Uses the fiber body parser to parse and validate the request body to struct T,
// marshaling the resulting struct to a byte slice to be used as a queue service message body.
func ParseValidateBodyToQueueBody[T any](ctx *fiber.Ctx, validator *validator.Validate, bodyType *T) (queueBody []byte, statusCode int, err error) {
	if err := ctx.BodyParser(bodyType); err != nil {
		return nil, fiber.StatusBadRequest, errors.New("invalid request body")
	}

	if err := validator.Struct(*bodyType); err != nil {
		return nil, fiber.StatusBadRequest, err
	}

	// Its usefull to marshal the struct back instead of using ctx.Body()
	// to remove json fields that might be outside T but are in the req body
	body, err := json.Marshal(bodyType)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.New("failed to marshal request json")
	}

	return body, fiber.StatusOK, nil
}
