package handlers

import (
	"context"
	"database/sql"
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

		if !perms.HasPermission(permissions.PlayerViewAllRoomImagesName) {
			c.Status(fiber.StatusForbidden)
			return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
		}

		roomImages, err := i.Queries.ListRooms(context.Background())
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		pageRoomImages := []fiber.Map{}
		for _, roomImage := range roomImages {
			pageRoomImage := fiber.Map{
				"RoomTitle": roomImage.Title,
				"ImageName": "ImageName",
				"Size":      roomImage.Size,
				"Path":      routes.RoomImagePath(roomImage.ID),
			}

			if perms.HasPermission(permissions.PlayerEditRoomImageName) {
				pageRoomImage["EditPath"] = routes.EditRoomImagePath(roomImage.ID)
			}
			pageRoomImages = append(pageRoomImages, pageRoomImage)
		}

		b := views.Bind(c)
		b["PageHeader"] = fiber.Map{
			"Title":    "Room Images",
			"SubTitle": "Room Images are what a room assumes its title, description, and other properties from",
		}
		b["RoomImages"] = pageRoomImages
		b["NewRoomImagePath"] = routes.NewRoomImage
		return c.Render(views.RoomImages, b)
	}
}

func RoomImagePage(i *shared.Interfaces) fiber.Handler {
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

		if !perms.HasPermission(permissions.PlayerViewAllRoomImagesName) {
			c.Status(fiber.StatusForbidden)
			return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
		}

		rmid, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
		}

		roomImage, err := i.Queries.GetRoom(context.Background(), rmid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return c.Render(views.NotFound, views.Bind(c), layouts.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		sizes := []string{
			"Tiny",
			"Small",
			"Medium",
			"Large",
			"Huge",
		}

		b := views.Bind(c)
		b["NavBack"] = fiber.Map{
			"Path":  routes.RoomImages,
			"Label": "Back to Room Images",
		}
		b["PageHeader"] = fiber.Map{
			"Title":    roomImage.Title,
			"SubTitle": "Room Image",
		}
		b["Name"] = "ImageName"
		b["Title"] = roomImage.Title
		b["Size"] = sizes[roomImage.Size]
		b["Description"] = roomImage.Description
		return c.Render(views.RoomImage, b, layouts.Main)
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

func EditRoomImagePage(i *shared.Interfaces) fiber.Handler {
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

		if !perms.HasPermission(permissions.PlayerEditRoomImageName) {
			c.Status(fiber.StatusForbidden)
			return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
		}

		id, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusNotFound)
			return c.Render(views.NotFound, views.Bind(c), layouts.Standalone)
		}

		roomImage, err := i.Queries.GetRoom(context.Background(), id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return c.Render(views.NotFound, views.Bind(c), layouts.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.NotFound, views.Bind(c), layouts.Standalone)
		}

		b := views.Bind(c)
		// TODO: Generalize this bind into a function
		b["NavBack"] = fiber.Map{
			"Path":  routes.RoomImages,
			"Label": "Back to Room Images",
		}
		// TODO: Put this in a helper function
		b["SizeRadioGroup"] = []fiber.Map{
			{
				"ID":       "edit-room-image-size-tiny",
				"Name":     "size",
				"Variable": "size",
				"Value":    "0",
				"Active":   roomImage.Size == 0,
				"Label":    "Tiny",
			},
			{
				"ID":       "edit-room-image-size-small",
				"Name":     "size",
				"Variable": "size",
				"Value":    "1",
				"Active":   roomImage.Size == 1,
				"Label":    "Small",
			},
			{
				"ID":       "edit-room-image-size-medium",
				"Name":     "size",
				"Variable": "size",
				"Value":    "2",
				"Active":   roomImage.Size == 2,
				"Label":    "Medium",
			},
			{
				"ID":       "edit-room-image-size-large",
				"Name":     "size",
				"Variable": "size",
				"Value":    "3",
				"Active":   roomImage.Size == 3,
				"Label":    "Large",
			},
			{
				"ID":       "edit-room-image-size-huge",
				"Name":     "size",
				"Variable": "size",
				"Value":    "4",
				"Active":   roomImage.Size == 4,
				"Label":    "Huge",
			},
		}
		b["PageHeader"] = fiber.Map{
			"Title":    roomImage.Title,
			"SubTitle": "ImageName",
		}
		b["RoomImagePath"] = routes.RoomImagePath(roomImage.ID)
		b["Name"] = "ImageName"
		b["Title"] = roomImage.Title
		b["Description"] = roomImage.Description
		b["Size"] = roomImage.Size
		return c.Render(views.EditRoomImage, b)
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

		_, err = i.Queries.CreateRoom(context.Background(), queries.CreateRoomParams{
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

func EditRoomImage(i *shared.Interfaces) fiber.Handler {
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
				SectionID:    "edit-room-image-error",
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
				SectionID:    "edit-room-image-error",
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
				SectionID:    "edit-room-image-error",
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}

		if !perms.HasPermission(permissions.PlayerEditRoomImageName) {
			c.Status(fiber.StatusForbidden)
			c.Append(shared.HeaderHXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    "edit-room-image-error",
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
				SectionID:    "edit-room-image-error",
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
				SectionID:    "edit-room-image-error",
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
				SectionID:    "edit-room-image-error",
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
				SectionID:    "edit-room-image-error",
				SectionClass: "pt-2",
				NoticeText: []string{
					"The Room Size you entered isn't valid.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}

		riid, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusNotFound)
			c.Append(shared.HeaderHXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    "edit-room-image-error",
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}

		if err := i.Queries.UpdateRoom(context.Background(), queries.UpdateRoomParams{
			ID:          riid,
			Title:       in.Title,
			Description: in.Description,
			Size:        in.Size,
		}); err != nil {
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

		c.Append("HX-Redirect", routes.RoomImagePath(riid))
		c.Append("HX-Reswap", "none")
		return nil
	}
}

func NewRoom(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(shared.HeaderHXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    "new-room-error",
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
				SectionID:    "new-room-error",
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}

		if !perms.HasPermission(permissions.PlayerCreateRoomName) {
			c.Status(fiber.StatusForbidden)
			c.Append(shared.HeaderHXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    "new-room-error",
				SectionClass: "pt-2",
				NoticeText: []string{
					"You don't have the permission(s) necessary to create a Room Image.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
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
		qtx := i.Queries.WithTx(tx)

		result, err := qtx.CreateRoom(context.Background(), queries.CreateRoomParams{
			Title:       rooms.DefaultTitle,
			Description: rooms.DefaultDescription,
			Size:        rooms.DefaultSize,
		})
		if err != nil {
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

		rid, err := result.LastInsertId()
		if err != nil {
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

		if err := tx.Commit(); err != nil {
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
		c.Append("HX-Redirect", routes.EditRoomPath(rid))
		c.Append("HX-Reswap", "none")
		return nil
	}
}

func EditRoomPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(views.Login, views.Bind(c), layouts.Standalone)
		}

		rmid, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
		}

		if !perms.HasPermission(permissions.PlayerCreateRoomName) {
			c.Status(fiber.StatusForbidden)
			return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
		}

		_, err = i.Queries.GetRoom(context.Background(), rmid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return c.Render(views.NotFound, views.Bind(c), layouts.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		roomGrid := []fiber.Map{
			{
				"ID": "test-room-grid-row-one",
				"Rooms": []fiber.Map{
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
				},
			},
			{
				"ID": "test-room-grid-row-two",
				"Rooms": []fiber.Map{
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 1},
					{"ID": 5},
				},
			},
			{
				"ID": "test-room-grid-row-three",
				"Rooms": []fiber.Map{
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 2},
					{"ID": 0},
				},
			},
			{
				"ID": "test-room-grid-row-four",
				"Rooms": []fiber.Map{
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
				},
			},
			{
				"ID": "test-room-grid-row-five",
				"Rooms": []fiber.Map{
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
				},
			},
		}

		b := views.Bind(c)
		// TODO: Generalize this bind into a function
		b["RoomGrid"] = roomGrid
		b["NavBack"] = fiber.Map{
			"Path":  routes.Rooms,
			"Label": "Back to Rooms",
		}
		b["PageHeader"] = fiber.Map{
			"Title":    "New Room",
			"SubTitle": "Create a new room, using a Room Image as a template",
		}
		b["RoomsPath"] = routes.RoomImages
		return c.Render(views.EditRoom, b)
	}
}

func RoomGrid() fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, err := util.GetPID(c)
		if err != nil {
			return nil
		}

		roomGridOne := []fiber.Map{
			{
				"ID": "test-room-grid-row-one",
				"Rooms": []fiber.Map{
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
				},
			},
			{
				"ID": "test-room-grid-row-two",
				"Rooms": []fiber.Map{
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 1},
					{"ID": 5},
				},
			},
			{
				"ID": "test-room-grid-row-three",
				"Rooms": []fiber.Map{
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 2},
					{"ID": 0},
				},
			},
			{
				"ID": "test-room-grid-row-four",
				"Rooms": []fiber.Map{
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
				},
			},
			{
				"ID": "test-room-grid-row-five",
				"Rooms": []fiber.Map{
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
				},
			},
		}

		roomGridTwo := []fiber.Map{
			{
				"ID": "test-room-grid-row-one",
				"Rooms": []fiber.Map{
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
				},
			},
			{
				"ID": "test-room-grid-row-two",
				"Rooms": []fiber.Map{
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 1},
					{"ID": 6},
				},
			},
			{
				"ID": "test-room-grid-row-three",
				"Rooms": []fiber.Map{
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 2},
					{"ID": 0},
				},
			},
			{
				"ID": "test-room-grid-row-four",
				"Rooms": []fiber.Map{
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
				},
			},
			{
				"ID": "test-room-grid-row-five",
				"Rooms": []fiber.Map{
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
					{"ID": 0},
				},
			},
		}

		id, err := util.GetID(c)
		if err != nil {
			return nil
		}

		selected, err := util.GetParamID(c, "selected")
		if err != nil {
			return nil
		}

		if id != 0 {
			return nil
		}

		b := views.Bind(c)
		if selected == 2 {
			b["RoomGrid"] = roomGridTwo
		} else {
			b["RoomGrid"] = roomGridOne
		}

		return c.Render(partials.RoomGrid, b, layouts.None)
	}
}
