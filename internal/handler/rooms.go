package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/headers"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/partials"
	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/rooms"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/util"
	"petrichormud.com/app/internal/views"
)

func RoomsPage(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !util.IsLoggedIn(c) {
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

		records, err := i.Queries.ListRooms(context.Background())
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		pageRooms := []fiber.Map{}
		for _, record := range records {
			var sb strings.Builder
			fmt.Fprintf(&sb, "[%d] %s", record.ID, record.Title)
			pageRoom := fiber.Map{
				"Title":      sb.String(),
				"Size":       record.Size,
				"SizeString": rooms.SizeToString(record.Size),
				"Path":       routes.RoomPath(record.ID),
			}

			if perms.HasPermission(permissions.PlayerCreateRoomName) {
				pageRoom["EditPath"] = routes.EditRoomPath(record.ID)
			}

			pageRooms = append(pageRooms, pageRoom)
		}

		b := views.Bind(c)
		b["Rooms"] = pageRooms
		b["PageHeader"] = fiber.Map{
			"Title":    "Rooms",
			"SubTitle": "Individual rooms, where their exits and individual properties are assigned",
		}
		return c.Render(views.Rooms, b)
	}
}

func RoomPage(i *interfaces.Shared) fiber.Handler {
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

		rmid, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
		}

		record, err := i.Queries.GetRoom(context.Background(), rmid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return c.Render(views.NotFound, views.Bind(c), layouts.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		b := views.Bind(c)
		b["NavBack"] = fiber.Map{
			"Path":  routes.Rooms,
			"Label": "Back to Rooms",
		}
		b["PageHeader"] = fiber.Map{
			"Title":    rooms.TitleWithID(record.Title, record.ID),
			"SubTitle": "Room",
		}
		b["Name"] = "ImageName"
		b["Title"] = record.Title
		b["Size"] = rooms.SizeToString(record.Size)
		b["Description"] = record.Description
		return c.Render(views.Room, b, layouts.Main)
	}
}

func NewRoom(i *interfaces.Shared) fiber.Handler {
	type input struct {
		Direction string `form:"direction"`
		LinkID    int64  `form:"id"`
		TwoWay    bool   `form:"two-way"`
	}

	const sectionID string = "edit-room-exits-create-error"

	return func(c *fiber.Ctx) error {
		_, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    sectionID,
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
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    sectionID,
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
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    sectionID,
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
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		result, err := qtx.CreateRoom(context.Background(), queries.CreateRoomParams{
			Title:       rooms.DefaultTitle,
			Description: rooms.DefaultDescription,
			Size:        rooms.DefaultSize,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    sectionID,
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
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}

		in := new(input)
		if err := c.BodyParser(in); err != nil && err != fiber.ErrUnprocessableEntity {
			c.Status(fiber.StatusBadRequest)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}

		if in.LinkID != 0 {
			// TODO: Can run validations on the room to be linked here, asserting that:
			// 1. The link-to exit isn't already filled
			// 2. There isn't a setpiece that leads to the proposed destination room
			// etc

			room, err := qtx.GetRoom(context.Background(), in.LinkID)
			if err != nil {
				if err == sql.ErrNoRows {
					c.Status(fiber.StatusNotFound)
					c.Append(headers.HXAcceptable, "true")
					c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
					return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
						SectionID:    sectionID,
						SectionClass: "pt-2",
						NoticeText: []string{
							"Something's gone terribly wrong.",
						},
						RefreshButton: true,
						NoticeIcon:    true,
					}), layouts.None)
				}
				c.Status(fiber.StatusInternalServerError)
				c.Append(headers.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
					SectionID:    sectionID,
					SectionClass: "pt-2",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				}), layouts.None)
			}

			exitRoom, err := qtx.GetRoom(context.Background(), rid)
			if err != nil {
				if err == sql.ErrNoRows {
					c.Status(fiber.StatusNotFound)
					c.Append(headers.HXAcceptable, "true")
					c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
					return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
						SectionID:    sectionID,
						SectionClass: "pt-2",
						NoticeText: []string{
							"Something's gone terribly wrong.",
						},
						RefreshButton: true,
						NoticeIcon:    true,
					}), layouts.None)
				}
				c.Status(fiber.StatusInternalServerError)
				c.Append(headers.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
					SectionID:    sectionID,
					SectionClass: "pt-2",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				}), layouts.None)
			}

			if !rooms.IsDirectionValid(in.Direction) {
				c.Status(fiber.StatusBadRequest)
				c.Append(headers.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
					SectionID:    sectionID,
					SectionClass: "pt-2",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				}), layouts.None)
			}

			if err := rooms.Link(rooms.LinkParams{
				Queries:   qtx,
				ID:        in.LinkID,
				To:        rid,
				Direction: in.Direction,
				TwoWay:    in.TwoWay,
			}); err != nil {
				c.Status(fiber.StatusInternalServerError)
				c.Append(headers.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
					SectionID:    sectionID,
					SectionClass: "pt-2",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				}), layouts.None)
			}

			room, err = qtx.GetRoom(context.Background(), room.ID)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				c.Append(headers.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
					SectionID:    sectionID,
					SectionClass: "pt-2",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				}), layouts.None)
			}
			exitRoom, err = qtx.GetRoom(context.Background(), exitRoom.ID)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				c.Append(headers.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
					SectionID:    sectionID,
					SectionClass: "pt-2",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				}), layouts.None)
			}

			exitGraph, err := rooms.BuildGraph(rooms.BuildGraphParams{
				Queries:  qtx,
				Room:     &room,
				MaxDepth: 1,
			})
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				c.Append(headers.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
					SectionID:    sectionID,
					SectionClass: "pt-2",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				}), layouts.None)
			}

			gridGraph, err := rooms.BuildGraph(rooms.BuildGraphParams{
				Queries: qtx,
				Room:    &room,
			})
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				c.Append(headers.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
					SectionID:    sectionID,
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
				c.Append(headers.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
					SectionID:    sectionID,
					SectionClass: "pt-2",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				}), layouts.None)
			}

			grid := gridGraph.BindMatrix(rooms.BindMatrixParams{
				Matrix:  rooms.EmptyBindMatrix(5),
				Row:     2,
				Col:     2,
				Shallow: false,
			})
			grid = rooms.AnnotateMatrixExits(grid)

			c.Status(fiber.StatusCreated)
			b := exitGraph.BindExit(in.Direction)
			b["Exits"] = exitGraph.BindExits()
			b["RoomGrid"] = grid
			return c.Render(partials.EditRoomExitEdit, b, layouts.EditRoomExitsSelect)
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    sectionID,
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

func EditRoomPage(i *interfaces.Shared) fiber.Handler {
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

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		room, err := qtx.GetRoom(context.Background(), rmid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return c.Render(views.NotFound, views.Bind(c), layouts.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		graph, err := rooms.BuildGraph(rooms.BuildGraphParams{
			Queries: qtx,
			Room:    &room,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		grid := graph.BindMatrix(rooms.BindMatrixParams{
			Matrix:  rooms.EmptyBindMatrix(5),
			Row:     2,
			Col:     2,
			Shallow: false,
		})
		grid = rooms.AnnotateMatrixExits(grid)
		exits := graph.BindExits()

		b := views.Bind(c)
		// TODO: Get a bind function for this
		b["NavBack"] = fiber.Map{
			"Path":  routes.Rooms,
			"Label": "Back to Rooms",
		}
		// TODO: Get a bind function for this too
		b["PageHeader"] = fiber.Map{
			"Title":    rooms.TitleWithID(room.Title, room.ID),
			"SubTitle": "Update room properties here",
		}
		b["RoomGrid"] = grid
		b["Title"] = room.Title
		b["TitlePath"] = routes.RoomTitlePath(room.ID)
		b["Description"] = room.Description
		b["DescriptionPath"] = routes.RoomDescriptionPath(room.ID)
		b["Size"] = room.Size
		b["SizePath"] = routes.RoomSizePath(room.ID)
		b = rooms.BindSizeRadioGroup(b, &room)
		// TODO: I don't think these individual dirs are needed
		b["North"] = room.North
		b["Northeast"] = room.Northeast
		b["East"] = room.East
		b["Southeast"] = room.Southeast
		b["South"] = room.South
		b["Southwest"] = room.Southwest
		b["West"] = room.West
		b["Northwest"] = room.Northwest
		b["Exits"] = exits
		return c.Render(views.EditRoom, b)
	}
}

func RoomGrid(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !util.IsLoggedIn(c) {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if !perms.HasPermission(permissions.PlayerCreateRoomName) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		rid, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		qtx := i.Queries.WithTx(tx)

		room, err := qtx.GetRoom(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		graph, err := rooms.BuildGraph(rooms.BuildGraphParams{
			Queries: qtx,
			Room:    &room,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		grid := graph.BindMatrix(rooms.BindMatrixParams{
			Matrix:  rooms.EmptyBindMatrix(5),
			Row:     2,
			Col:     2,
			Shallow: false,
		})
		grid = rooms.AnnotateMatrixExits(grid)

		b := fiber.Map{}
		b["RoomGrid"] = grid
		return c.Render(partials.RoomGrid, b, layouts.None)
	}
}

func EditRoomExit(i *interfaces.Shared) fiber.Handler {
	type input struct {
		Direction string `form:"direction"`
		LinkID    int64  `form:"id"`
		TwoWay    bool   `form:"two-way"`
	}

	// TODO: Get this in a shared constant
	const sectionID string = "edit-room-exits-edit-error"

	internalServerErrorNoticeParams := partials.BindNoticeSectionParams{
		SectionID:    sectionID,
		SectionClass: "pt-2",
		NoticeText: []string{
			"Something's gone terribly wrong.",
		},
		RefreshButton: true,
		NoticeIcon:    true,
	}

	sessionExpiredNoticeParams := partials.BindNoticeSectionParams{
		SectionID:    sectionID,
		SectionClass: "pt-2",
		NoticeText: []string{
			"It looks like your session may have expired.",
		},
		RefreshButton: true,
		NoticeIcon:    true,
	}

	noPermissionNoticeParams := partials.BindNoticeSectionParams{
		SectionID:    sectionID,
		SectionClass: "pt-2",
		NoticeText: []string{
			"You don't have the permission required to edit this room.",
		},
		RefreshButton: true,
		NoticeIcon:    true,
	}

	notFoundNoticeParams := partials.BindNoticeSectionParams{
		SectionID:    sectionID,
		SectionClass: "pt-2",
		NoticeText: []string{
			"The room you're looking for no longer exists.",
		},
		RefreshButton: true,
		NoticeIcon:    true,
	}

	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil && err != fiber.ErrUnprocessableEntity {
			c.Status(fiber.StatusBadRequest)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(internalServerErrorNoticeParams), layouts.None)
		}

		if !rooms.IsDirectionValid(in.Direction) {
			c.Status(fiber.StatusBadRequest)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(internalServerErrorNoticeParams), layouts.None)
		}

		if in.LinkID == 0 {
			c.Status(fiber.StatusBadRequest)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(internalServerErrorNoticeParams), layouts.None)
		}

		_, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(sessionExpiredNoticeParams), layouts.None)
		}

		rid, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(internalServerErrorNoticeParams), layouts.None)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(internalServerErrorNoticeParams), layouts.None)
		}

		if !perms.HasPermission(permissions.PlayerCreateRoomName) {
			c.Status(fiber.StatusForbidden)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(noPermissionNoticeParams), layouts.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(internalServerErrorNoticeParams), layouts.None)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		room, err := qtx.GetRoom(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				c.Append(headers.HXAcceptable, "true")
				return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(notFoundNoticeParams), layouts.None)
			}
			c.Status(fiber.StatusInternalServerError)
			c.Append(headers.HXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(internalServerErrorNoticeParams), layouts.None)
		}
		exitRoom, err := qtx.GetRoom(context.Background(), in.LinkID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				c.Append(headers.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(notFoundNoticeParams), layouts.None)
			}
			c.Status(fiber.StatusInternalServerError)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(internalServerErrorNoticeParams), layouts.None)
		}

		if err := rooms.Link(rooms.LinkParams{
			Queries:   qtx,
			ID:        rid,
			To:        in.LinkID,
			Direction: in.Direction,
			TwoWay:    in.TwoWay,
		}); err != nil {
			if err == rooms.ErrLinkSelf {
				c.Status(fiber.StatusBadRequest)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		room, err = qtx.GetRoom(context.Background(), room.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				c.Append(headers.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
					SectionID:    sectionID,
					SectionClass: "pt-2",
					NoticeText: []string{
						"The room you're looking for no longer exists.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				}), layouts.None)
			}
			c.Status(fiber.StatusInternalServerError)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}
		exitRoom, err = qtx.GetRoom(context.Background(), exitRoom.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				c.Append(headers.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
					SectionID:    sectionID,
					SectionClass: "pt-2",
					NoticeText: []string{
						"The room you're looking for no longer exists.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				}), layouts.None)
			}
			c.Status(fiber.StatusInternalServerError)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}

		graph, err := rooms.BuildGraph(rooms.BuildGraphParams{
			Queries: qtx,
			Room:    &room,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    sectionID,
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
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}

		grid := graph.BindMatrix(rooms.BindMatrixParams{
			Matrix:  rooms.EmptyBindMatrix(5),
			Row:     2,
			Col:     2,
			Shallow: false,
		})
		grid = rooms.AnnotateMatrixExits(grid)

		c.Status(fiber.StatusOK)
		b := graph.BindExit(in.Direction)
		b["Exits"] = graph.BindExits()
		b["RoomGrid"] = grid
		return c.Render(partials.EditRoomExitEdit, b, layouts.EditRoomExitsSelect)
	}
}

func ClearRoomExit(i *interfaces.Shared) fiber.Handler {
	// TODO: Get constants for common section IDs
	const sectionID string = "edit-room-exits-edit-error"

	internalServerErrorNoticeParams := partials.BindNoticeSectionParams{
		SectionID:    sectionID,
		SectionClass: "pt-2",
		NoticeText: []string{
			"Something's gone terribly wrong.",
		},
		RefreshButton: true,
		NoticeIcon:    true,
	}

	sessionExpiredNoticeParams := partials.BindNoticeSectionParams{
		SectionID:    sectionID,
		SectionClass: "pt-2",
		NoticeText: []string{
			"It looks like your session may have expired.",
		},
		RefreshButton: true,
		NoticeIcon:    true,
	}

	noPermissionNoticeParams := partials.BindNoticeSectionParams{
		SectionID:    sectionID,
		SectionClass: "pt-2",
		NoticeText: []string{
			"You don't have the permission required to edit this room.",
		},
		RefreshButton: true,
		NoticeIcon:    true,
	}

	notFoundNoticeParams := partials.BindNoticeSectionParams{
		SectionID:    sectionID,
		SectionClass: "pt-2",
		NoticeText: []string{
			"The room you're looking for no longer exists.",
		},
		RefreshButton: true,
		NoticeIcon:    true,
	}

	return func(c *fiber.Ctx) error {
		_, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(sessionExpiredNoticeParams), layouts.None)
		}

		rid, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(internalServerErrorNoticeParams), layouts.None)
		}

		dir := c.Params("exit")
		if !rooms.IsDirectionValid(dir) {
			c.Status(fiber.StatusBadRequest)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(internalServerErrorNoticeParams), layouts.None)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(internalServerErrorNoticeParams), layouts.None)
		}

		if !perms.HasPermission(permissions.PlayerCreateRoomName) {
			c.Status(fiber.StatusForbidden)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(noPermissionNoticeParams), layouts.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(internalServerErrorNoticeParams), layouts.None)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		room, err := qtx.GetRoom(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				c.Append(headers.HXAcceptable, "true")
				return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(notFoundNoticeParams), layouts.None)
			}
			c.Status(fiber.StatusInternalServerError)
			c.Append(headers.HXAcceptable, "true")
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(internalServerErrorNoticeParams), layouts.None)
		}

		exitID := rooms.ExitID(&room, dir)
		exitRoom, err := qtx.GetRoom(context.Background(), exitID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				c.Append(headers.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(notFoundNoticeParams), layouts.None)
			}
			c.Status(fiber.StatusInternalServerError)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(internalServerErrorNoticeParams), layouts.None)
		}
		exitDir, err := rooms.ExitDirection(&exitRoom, rid)
		if err != nil && err != rooms.ErrExitIDNotFound {
			c.Status(fiber.StatusInternalServerError)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(internalServerErrorNoticeParams), layouts.None)
		}
		if err != rooms.ErrExitIDNotFound {
			if err := rooms.Unlink(rooms.UnlinkParams{
				Queries:   qtx,
				ID:        exitID,
				Direction: exitDir,
			}); err != nil {
				c.Status(fiber.StatusInternalServerError)
				c.Append(headers.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(internalServerErrorNoticeParams), layouts.None)
			}
		}

		if err := rooms.Unlink(rooms.UnlinkParams{
			Queries:   qtx,
			ID:        rid,
			Direction: dir,
		}); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		room, err = qtx.GetRoom(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				c.Append(headers.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(notFoundNoticeParams), layouts.None)
			}
		}

		graph, err := rooms.BuildGraph(rooms.BuildGraphParams{
			Queries: qtx,
			Room:    &room,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(internalServerErrorNoticeParams), layouts.None)
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(headers.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partials.NoticeSectionError, partials.BindNoticeSection(partials.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layouts.None)
		}

		grid := graph.BindMatrix(rooms.BindMatrixParams{
			Matrix:  rooms.EmptyBindMatrix(5),
			Row:     2,
			Col:     2,
			Shallow: false,
		})
		grid = rooms.AnnotateMatrixExits(grid)

		c.Status(fiber.StatusOK)
		b := graph.BindEmptyExit(dir)
		b["Exits"] = graph.BindExits()
		b["RoomGrid"] = grid
		return c.Render(partials.EditRoomExitEdit, b, layouts.EditRoomExitsSelect)
	}
}

func EditRoomTitle(i *interfaces.Shared) fiber.Handler {
	type input struct {
		Title string `form:"title"`
	}

	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !rooms.IsTitleValid(in.Title) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !util.IsLoggedIn(c) {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return nil
		}
		if !perms.HasPermission(permissions.PlayerCreateRoomName) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		rmid, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		room, err := qtx.GetRoom(context.Background(), rmid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if room.Title == in.Title {
			c.Status(fiber.StatusConflict)
			return nil
		}

		if err := qtx.UpdateRoomTitle(context.Background(), queries.UpdateRoomTitleParams{
			ID:    room.ID,
			Title: in.Title,
		}); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		room, err = qtx.GetRoom(context.Background(), rmid)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		b := fiber.Map{}
		b["PageHeader"] = fiber.Map{
			"Title":    rooms.TitleWithID(room.Title, room.ID),
			"SubTitle": "Update room properties here",
		}
		b["Title"] = room.Title
		b["TitlePath"] = routes.RoomTitlePath(rmid)
		b["NoticeSection"] = partials.BindNoticeSection(partials.BindNoticeSectionParams{
			Success:      true,
			SectionID:    "room-edit-title-notice",
			SectionClass: "pb-2",
			NoticeText: []string{
				"Success! The room title has been updated.",
			},
			NoticeIcon: true,
		})
		return c.Render(partials.RoomEditTitle, b, layouts.PageHeader)
	}
}

func EditRoomDescription(i *interfaces.Shared) fiber.Handler {
	type input struct {
		Description string `form:"desc"`
	}

	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !rooms.IsDescriptionValid(in.Description) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !util.IsLoggedIn(c) {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return nil
		}
		if !perms.HasPermission(permissions.PlayerCreateRoomName) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		rmid, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		room, err := qtx.GetRoom(context.Background(), rmid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		// TODO: Add Conflict tests for edit room title, description and size
		if room.Description == in.Description {
			c.Status(fiber.StatusConflict)
			return nil
		}

		if err := qtx.UpdateRoomDescription(context.Background(), queries.UpdateRoomDescriptionParams{
			ID:          room.ID,
			Description: in.Description,
		}); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		room, err = qtx.GetRoom(context.Background(), rmid)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		b := fiber.Map{}
		b["Description"] = room.Description
		b["DescriptionPath"] = routes.RoomDescriptionPath(rmid)
		b["NoticeSection"] = partials.BindNoticeSection(partials.BindNoticeSectionParams{
			Success:      true,
			SectionID:    "room-edit-title-notice",
			SectionClass: "pb-2",
			NoticeText: []string{
				"Success! The room description has been updated.",
			},
			NoticeIcon: true,
		})
		return c.Render(partials.RoomEditDescription, b, layouts.None)
	}
}

func EditRoomSize(i *interfaces.Shared) fiber.Handler {
	type input struct {
		Size int32 `form:"size"`
	}

	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !rooms.IsSizeValid(in.Size) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !util.IsLoggedIn(c) {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return nil
		}
		if !perms.HasPermission(permissions.PlayerCreateRoomName) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		rmid, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		tx, err := i.Database.Begin()
		if err != nil {
			log.Println(err)
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		room, err := qtx.GetRoom(context.Background(), rmid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			log.Println(err)
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		// TODO: Add Conflict tests for edit room title, description and size
		if room.Size == in.Size {
			c.Status(fiber.StatusConflict)
			return nil
		}

		if err := qtx.UpdateRoomSize(context.Background(), queries.UpdateRoomSizeParams{
			ID:   room.ID,
			Size: in.Size,
		}); err != nil {
			log.Println(err)
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		room, err = qtx.GetRoom(context.Background(), rmid)
		if err != nil {
			log.Println(err)
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err := tx.Commit(); err != nil {
			log.Println(err)
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		b := fiber.Map{}
		b["Size"] = room.Size
		b["SizePath"] = routes.RoomSizePath(rmid)
		b = rooms.BindSizeRadioGroup(b, &room)
		b["NoticeSection"] = partials.BindNoticeSection(partials.BindNoticeSectionParams{
			Success:      true,
			SectionID:    "room-edit-title-notice",
			SectionClass: "pb-2",
			NoticeText: []string{
				"Success! The room size has been updated.",
			},
			NoticeIcon: true,
		})
		return c.Render(partials.RoomEditSize, b, layouts.None)
	}
}
