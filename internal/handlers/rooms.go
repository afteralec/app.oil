package handlers

import (
	"context"
	"log"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/partials"
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

		page_room_images := []fiber.Map{}
		for _, room_image := range room_images {
			page_room_images = append(page_room_images, fiber.Map{
				"RoomTitle": room_image.Title,
				"ImageName": room_image.Name,
				"Size":      room_image.Size,
				"Path":      routes.RoomImagePath(room_image.ID),
			})
		}

		b := views.Bind(c)
		b["PageHeader"] = fiber.Map{
			"Title":    "Room Images",
			"SubTitle": "Room Images are what a room assumes its title, description, and other properties from",
		}
		b["RoomImages"] = page_room_images
		b["NewRoomImagePath"] = routes.NewRoomImage
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
		// TODO: Generalize this bind into a function
		b["NavBack"] = fiber.Map{
			"Path":  routes.RoomImages,
			"Label": "Back to Room Images",
		}
		b["SizeRadioGroup"] = []fiber.Map{
			{
				"ID":       "new-room-image-size-tiny",
				"Name":     "size",
				"Variable": "size",
				"Value":    "0",
				"Active":   "false",
				"Label":    "Tiny",
			},
			{
				"ID":       "new-room-image-size-small",
				"Name":     "size",
				"Variable": "size",
				"Value":    "1",
				"Active":   "false",
				"Label":    "Small",
			},
			{
				"ID":       "new-room-image-size-medium",
				"Name":     "size",
				"Variable": "size",
				"Value":    "2",
				"Active":   "true",
				"Label":    "Medium",
			},
			{
				"ID":       "new-room-image-size-large",
				"Name":     "size",
				"Variable": "size",
				"Value":    "3",
				"Active":   "false",
				"Label":    "Large",
			},
			{
				"ID":       "new-room-image-size-huge",
				"Name":     "size",
				"Variable": "size",
				"Value":    "4",
				"Active":   "false",
				"Label":    "Huge",
			},
		}
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
			c.Append(shared.HeaderHXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    "new-room-image-error",
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}

		_, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(shared.HeaderHXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    "new-room-image-error",
				SectionClass: "pt-2",
				NoticeText: []string{
					"It looks like your session may have expired.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			c.Append(shared.HeaderHXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    "new-room-image-error",
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}

		if !perms.HasPermission(permissions.PlayerCreateRoomImageName) {
			c.Status(fiber.StatusForbidden)
			c.Append(shared.HeaderHXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    "new-room-image-error",
				SectionClass: "pt-2",
				NoticeText: []string{
					"You don't have the permission(s) necessary to create a Room Image.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}

		if !rooms.IsImageNameValid(in.Name) {
			c.Status(fiber.StatusBadRequest)
			c.Append(shared.HeaderHXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    "new-room-image-error",
				SectionClass: "pt-2",
				NoticeText: []string{
					"The Image Name you entered isn't valid.",
					"Please use only lowercase letters and dashes.",
				},
				NoticeIcon: true,
			}), layouts.None)
		}

		if !rooms.IsTitleValid(in.Title) {
			c.Status(fiber.StatusBadRequest)
			c.Append(shared.HeaderHXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    "new-room-image-error",
				SectionClass: "pt-2",
				NoticeText: []string{
					"The Room Title you entered isn't valid.",
					"Please try again.",
				},
				NoticeIcon: true,
			}), layouts.None)
		}

		if !rooms.IsDescriptionValid(in.Description) {
			c.Status(fiber.StatusBadRequest)
			c.Append(shared.HeaderHXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    "new-room-image-error",
				SectionClass: "pt-2",
				NoticeText: []string{
					"The Room Description you entered isn't valid.",
					"Please try again.",
				},
				NoticeIcon: true,
			}), layouts.None)
		}

		if !rooms.IsSizeValid(in.Size) {
			c.Status(fiber.StatusBadRequest)
			c.Append(shared.HeaderHXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    "new-room-image-error",
				SectionClass: "pt-2",
				NoticeText: []string{
					"The Room Size you entered isn't valid.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
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
			c.Append(shared.HeaderHXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    "new-room-image-error",
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}

		c.Status(fiber.StatusCreated)
		c.Append("HX-Redirect", routes.RoomImages)
		c.Append("HX-Reswap", "none")
		return nil
	}
}
