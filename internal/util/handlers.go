package util

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	fiber "github.com/gofiber/fiber/v2"
	"petrichormud.com/app/internal/permissions"
)

var (
	ErrNoPID      error = errors.New("no PID value found")
	ErrInvalidPID error = errors.New("invalid PID value")
)

func GetPID(c *fiber.Ctx) (int64, error) {
	// TODO: Get this locals key into a constant
	lpid := c.Locals("pid")
	if lpid == nil {
		return 0, ErrNoPID
	}
	pid, ok := lpid.(int64)
	if !ok {
		return 0, ErrInvalidPID
	}

	return pid, nil
}

var ErrNoID error = errors.New("no ID value found")

func GetID(c *fiber.Ctx) (int64, error) {
	// TODO: Get this locals key into a constant
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

func GetParamID(c *fiber.Ctx, key string) (int64, error) {
	param := c.Params(key)
	if len(param) == 0 {
		return 0, ErrNoID
	}
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}

var ErrNoPermissions error = errors.New("no permissions found")

func GetPermissions(c *fiber.Ctx) (permissions.PlayerGranted, error) {
	// TODO: Get this locals key into a constant
	lperms := c.Locals("perms")
	if lperms == nil {
		return permissions.PlayerGranted{}, nil
	}
	perms, ok := lperms.(permissions.PlayerGranted)
	if !ok {
		return permissions.PlayerGranted{}, nil
	}
	return perms, nil
}

func PrependHTMLID(id string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "#%s", id)
	return sb.String()
}
