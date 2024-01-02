package util

import (
	"errors"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"
)

var (
	ErrNoPID      error = errors.New("no PID value found")
	ErrInvalidPID error = errors.New("invalid PID value")
)

func GetPID(c *fiber.Ctx) (int64, error) {
	lpid := c.Locals("pid")
	if lpid == nil {
		return 0, nil
	}
	pid, ok := lpid.(int64)
	if !ok {
		return 0, ErrInvalidPID
	}

	return pid, nil
}

var ErrNoID error = errors.New("no ID value found")

func GetID(c *fiber.Ctx) (int64, error) {
	param := c.Params("id")
	if len(param) == 0 {
		return 0, ErrNoID
	}
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}
