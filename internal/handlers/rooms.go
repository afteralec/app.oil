package handlers

import (
	"context"
	"log"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/rooms"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/util"
	"petrichormud.com/app/internal/views"
)

func RoomsPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: IsLoggedIn helper?
		_, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(views.Login, views.Bind(c), layouts.Standalone)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
		}

		if !perms.HasPermission(permissions.PlayerViewAllRoomsName) {
			c.Status(fiber.StatusForbidden)
			return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
		}

		rooms, err := i.Queries.ListRooms(context.Background())
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		b := views.Bind(c)
		b["Rooms"] = rooms
		b["PageHeader"] = fiber.Map{
			"Title":    "Rooms",
			"SubTitle": "Individual rooms, where their exits and individual properties are assigned",
		}
		return c.Render(views.Rooms, b)
	}
}

func RoomImagesPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(views.Login, views.Bind(c), layouts.Standalone)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
		}

		if !perms.HasPermission(permissions.PlayerViewAllRoomsName) {
			c.Status(fiber.StatusForbidden)
			return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
		}

		room_images, err := i.Queries.ListRoomImages(context.Background())
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		b := views.Bind(c)
		b["PageHeader"] = fiber.Map{
			"Title":    "Room Images",
			"SubTitle": "Room Images are what a room assumes its title, description, and other properties from",
		}
		b["RoomImages"] = room_images
		return c.Render(views.RoomImages, b)
	}
}

func NewRoomImagePage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(views.Login, views.Bind(c), layouts.Standalone)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
		}

		if !perms.HasPermission(permissions.PlayerCreateRoomImageName) {
			c.Status(fiber.StatusForbidden)
			return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
		}

		b := views.Bind(c)
		b["SizeTiny"] = 0
		b["SizeSmall"] = 1
		b["SizeMedium"] = 2
		b["SizeLarge"] = 3
		b["SizeHuge"] = 4
		b["SizeIsMedium"] = true
		b["PageHeader"] = fiber.Map{
			"Title":    "New Room Image",
			"SubTitle": "Room Images are what a room assumes its title, description, and other properties from",
		}
		b["RoomImagesPath"] = routes.RoomImages
		return c.Render(views.NewRoomImage, b)
	}
}

func NewRoomImage(i *shared.Interfaces) fiber.Handler {
	type input struct {
		Name        string `form:"name"`
		Title       string `form:"title"`
		Description string `form:"description"`
		Size        int32  `form:"size"`
	}

	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		_, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if !perms.HasPermission(permissions.PlayerCreateRoomImageName) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if !rooms.IsImageNameValid(in.Name) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !rooms.IsTitleValid(in.Title) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !rooms.IsDescriptionValid(in.Description) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !rooms.IsSizeValid(in.Size) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		_, err = i.Queries.CreateRoomImage(context.Background(), queries.CreateRoomImageParams{
			Name:        in.Name,
			Title:       in.Title,
			Description: in.Description,
			Size:        in.Size,
		})
		if err != nil {
			log.Println(err)
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusCreated)
		return nil
	}
}
