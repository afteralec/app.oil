package handler

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/header"
	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/partial"
	"petrichormud.com/app/internal/player"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/room"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/service"
	"petrichormud.com/app/internal/util"
	"petrichormud.com/app/internal/view"
)

func RoomsPage(i *service.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !util.IsLoggedIn(c) {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(view.Login, view.Bind(c), layout.Standalone)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return c.Render(view.Forbidden, view.Bind(c), layout.Standalone)
		}

		if !perms.HasPermission(player.PermissionViewAllRooms.Name) {
			c.Status(fiber.StatusForbidden)
			return c.Render(view.Forbidden, view.Bind(c), layout.Standalone)
		}

		records, err := i.Queries.ListRooms(context.Background())
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		pageRooms := []fiber.Map{}
		for _, record := range records {
			var sb strings.Builder
			fmt.Fprintf(&sb, "[%d] %s", record.ID, record.Title)
			pageRoom := fiber.Map{
				"Title":      sb.String(),
				"Size":       record.Size,
				"SizeString": room.SizeToString(record.Size),
				"Path":       route.RoomPath(record.ID),
			}

			if perms.HasPermission(player.PermissionCreateRoom.Name) {
				pageRoom["EditPath"] = route.EditRoomPath(record.ID)
			}

			pageRooms = append(pageRooms, pageRoom)
		}

		b := view.Bind(c)
		b["Rooms"] = pageRooms
		b["PageHeader"] = fiber.Map{
			"Title":    "Rooms",
			"SubTitle": "Individual rooms, where their exits and individual properties are assigned",
		}
		return c.Render(view.Rooms, b)
	}
}

func RoomPage(i *service.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(view.Login, view.Bind(c), layout.Standalone)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return c.Render(view.Forbidden, view.Bind(c), layout.Standalone)
		}

		if !perms.HasPermission(player.PermissionViewAllRooms.Name) {
			c.Status(fiber.StatusForbidden)
			return c.Render(view.Forbidden, view.Bind(c), layout.Standalone)
		}

		rmid, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.Render(view.Forbidden, view.Bind(c), layout.Standalone)
		}

		record, err := i.Queries.GetRoom(context.Background(), rmid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return c.Render(view.NotFound, view.Bind(c), layout.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		b := view.Bind(c)
		b["NavBack"] = fiber.Map{
			"Path":  route.Rooms,
			"Label": "Back to Rooms",
		}
		b["PageHeader"] = fiber.Map{
			"Title":    room.TitleWithID(record.Title, record.ID),
			"SubTitle": "Room",
		}
		b["Name"] = "ImageName"
		b["Title"] = record.Title
		b["Size"] = room.SizeToString(record.Size)
		b["Description"] = record.Description
		return c.Render(view.Room, b, layout.Main)
	}
}

func NewRoom(i *service.Interfaces) fiber.Handler {
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
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"It looks like your session may have expired.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layout.None)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layout.None)
		}

		if !perms.HasPermission(player.PermissionCreateRoom.Name) {
			c.Status(fiber.StatusForbidden)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"You don't have the permission(s) necessary to create a Room Image.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layout.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layout.None)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		result, err := qtx.CreateRoom(context.Background(), query.CreateRoomParams{
			Title:       room.DefaultTitle,
			Description: room.DefaultDescription,
			Size:        room.DefaultSize,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layout.None)
		}

		rid, err := result.LastInsertId()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layout.None)
		}

		in := new(input)
		if err := c.BodyParser(in); err != nil && err != fiber.ErrUnprocessableEntity {
			c.Status(fiber.StatusBadRequest)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layout.None)
		}

		if in.LinkID != 0 {
			// TODO: Can run validations on the room to be linked here, asserting that:
			// 1. The link-to exit isn't already filled
			// 2. There isn't a setpiece that leads to the proposed destination room
			// etc

			rm, err := qtx.GetRoom(context.Background(), in.LinkID)
			if err != nil {
				if err == sql.ErrNoRows {
					c.Status(fiber.StatusNotFound)
					c.Append(header.HXAcceptable, "true")
					c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
					return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
						SectionID:    sectionID,
						SectionClass: "pt-2",
						NoticeText: []string{
							"Something's gone terribly wrong.",
						},
						RefreshButton: true,
						NoticeIcon:    true,
					}), layout.None)
				}
				c.Status(fiber.StatusInternalServerError)
				c.Append(header.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
					SectionID:    sectionID,
					SectionClass: "pt-2",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				}), layout.None)
			}

			exitrm, err := qtx.GetRoom(context.Background(), rid)
			if err != nil {
				if err == sql.ErrNoRows {
					c.Status(fiber.StatusNotFound)
					c.Append(header.HXAcceptable, "true")
					c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
					return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
						SectionID:    sectionID,
						SectionClass: "pt-2",
						NoticeText: []string{
							"Something's gone terribly wrong.",
						},
						RefreshButton: true,
						NoticeIcon:    true,
					}), layout.None)
				}
				c.Status(fiber.StatusInternalServerError)
				c.Append(header.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
					SectionID:    sectionID,
					SectionClass: "pt-2",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				}), layout.None)
			}

			if !room.IsDirectionValid(in.Direction) {
				c.Status(fiber.StatusBadRequest)
				c.Append(header.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
					SectionID:    sectionID,
					SectionClass: "pt-2",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				}), layout.None)
			}

			if err := room.Link(room.LinkParams{
				Queries:   qtx,
				ID:        in.LinkID,
				To:        rid,
				Direction: in.Direction,
				TwoWay:    in.TwoWay,
			}); err != nil {
				c.Status(fiber.StatusInternalServerError)
				c.Append(header.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
					SectionID:    sectionID,
					SectionClass: "pt-2",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				}), layout.None)
			}

			rm, err = qtx.GetRoom(context.Background(), rm.ID)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				c.Append(header.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
					SectionID:    sectionID,
					SectionClass: "pt-2",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				}), layout.None)
			}
			exitrm, err = qtx.GetRoom(context.Background(), exitrm.ID)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				c.Append(header.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
					SectionID:    sectionID,
					SectionClass: "pt-2",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				}), layout.None)
			}

			exitGraph, err := room.BuildGraph(room.BuildGraphParams{
				Queries:  qtx,
				Room:     &rm,
				MaxDepth: 1,
			})
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				c.Append(header.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
					SectionID:    sectionID,
					SectionClass: "pt-2",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				}), layout.None)
			}

			gridGraph, err := room.BuildGraph(room.BuildGraphParams{
				Queries: qtx,
				Room:    &rm,
			})
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				c.Append(header.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
					SectionID:    sectionID,
					SectionClass: "pt-2",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				}), layout.None)
			}

			if err := tx.Commit(); err != nil {
				c.Status(fiber.StatusInternalServerError)
				c.Append(header.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
					SectionID:    sectionID,
					SectionClass: "pt-2",
					NoticeText: []string{
						"Something's gone terribly wrong.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				}), layout.None)
			}

			grid := gridGraph.BindMatrix(room.BindMatrixParams{
				Matrix:  room.EmptyBindMatrix(5),
				Row:     2,
				Col:     2,
				Shallow: false,
			})
			grid = room.AnnotateMatrixExits(grid)

			c.Status(fiber.StatusCreated)
			b := exitGraph.BindExit(in.Direction)
			b["Exits"] = exitGraph.BindExits()
			b["RoomGrid"] = grid
			return c.Render(partial.EditRoomExitEdit, b, layout.EditRoomExitsSelect)
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layout.None)
		}

		c.Status(fiber.StatusCreated)
		c.Append("HX-Redirect", route.EditRoomPath(rid))
		c.Append("HX-Reswap", "none")
		return nil
	}
}

func EditRoomPage(i *service.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(view.Login, view.Bind(c), layout.Standalone)
		}

		rmid, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			return c.Render(view.Forbidden, view.Bind(c), layout.Standalone)
		}

		if !perms.HasPermission(player.PermissionCreateRoom.Name) {
			c.Status(fiber.StatusForbidden)
			return c.Render(view.Forbidden, view.Bind(c), layout.Standalone)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		rm, err := qtx.GetRoom(context.Background(), rmid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return c.Render(view.NotFound, view.Bind(c), layout.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		graph, err := room.BuildGraph(room.BuildGraphParams{
			Queries: qtx,
			Room:    &rm,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		grid := graph.BindMatrix(room.BindMatrixParams{
			Matrix:  room.EmptyBindMatrix(5),
			Row:     2,
			Col:     2,
			Shallow: false,
		})
		grid = room.AnnotateMatrixExits(grid)
		exits := graph.BindExits()

		b := view.Bind(c)
		// TODO: Get a bind function for this
		b["NavBack"] = fiber.Map{
			"Path":  route.Rooms,
			"Label": "Back to Rooms",
		}
		// TODO: Get a bind function for this too
		b["PageHeader"] = fiber.Map{
			"Title":    room.TitleWithID(rm.Title, rm.ID),
			"SubTitle": "Update room properties here",
		}
		b["RoomGrid"] = grid
		b["Title"] = rm.Title
		b["TitlePath"] = route.RoomTitlePath(rm.ID)
		b["Description"] = rm.Description
		b["DescriptionPath"] = route.RoomDescriptionPath(rm.ID)
		b["Size"] = rm.Size
		b["SizePath"] = route.RoomSizePath(rm.ID)
		b = room.BindSizeRadioGroup(b, &rm)
		// TODO: I don't think these individual dirs are needed
		b["North"] = rm.North
		b["Northeast"] = rm.Northeast
		b["East"] = rm.East
		b["Southeast"] = rm.Southeast
		b["South"] = rm.South
		b["Southwest"] = rm.Southwest
		b["West"] = rm.West
		b["Northwest"] = rm.Northwest
		b["Exits"] = exits
		return c.Render(view.EditRoom, b)
	}
}

func RoomGrid(i *service.Interfaces) fiber.Handler {
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

		if !perms.HasPermission(player.PermissionCreateRoom.Name) {
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

		rm, err := qtx.GetRoom(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		graph, err := room.BuildGraph(room.BuildGraphParams{
			Queries: qtx,
			Room:    &rm,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		grid := graph.BindMatrix(room.BindMatrixParams{
			Matrix:  room.EmptyBindMatrix(5),
			Row:     2,
			Col:     2,
			Shallow: false,
		})
		grid = room.AnnotateMatrixExits(grid)

		b := fiber.Map{}
		b["RoomGrid"] = grid
		return c.Render(partial.RoomGrid, b, layout.None)
	}
}

func EditRoomExit(i *service.Interfaces) fiber.Handler {
	type input struct {
		Direction string `form:"direction"`
		LinkID    int64  `form:"id"`
		TwoWay    bool   `form:"two-way"`
	}

	// TODO: Get this in a shared constant
	const sectionID string = "edit-room-exits-edit-error"

	internalServerErrorNoticeParams := partial.BindNoticeSectionParams{
		SectionID:    sectionID,
		SectionClass: "pt-2",
		NoticeText: []string{
			"Something's gone terribly wrong.",
		},
		RefreshButton: true,
		NoticeIcon:    true,
	}

	sessionExpiredNoticeParams := partial.BindNoticeSectionParams{
		SectionID:    sectionID,
		SectionClass: "pt-2",
		NoticeText: []string{
			"It looks like your session may have expired.",
		},
		RefreshButton: true,
		NoticeIcon:    true,
	}

	noPermissionNoticeParams := partial.BindNoticeSectionParams{
		SectionID:    sectionID,
		SectionClass: "pt-2",
		NoticeText: []string{
			"You don't have the permission required to edit this room.",
		},
		RefreshButton: true,
		NoticeIcon:    true,
	}

	notFoundNoticeParams := partial.BindNoticeSectionParams{
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
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(internalServerErrorNoticeParams), layout.None)
		}

		if !room.IsDirectionValid(in.Direction) {
			c.Status(fiber.StatusBadRequest)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(internalServerErrorNoticeParams), layout.None)
		}

		if in.LinkID == 0 {
			c.Status(fiber.StatusBadRequest)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(internalServerErrorNoticeParams), layout.None)
		}

		_, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(sessionExpiredNoticeParams), layout.None)
		}

		rid, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(internalServerErrorNoticeParams), layout.None)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(internalServerErrorNoticeParams), layout.None)
		}

		if !perms.HasPermission(player.PermissionCreateRoom.Name) {
			c.Status(fiber.StatusForbidden)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(noPermissionNoticeParams), layout.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(internalServerErrorNoticeParams), layout.None)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		rm, err := qtx.GetRoom(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				c.Append(header.HXAcceptable, "true")
				return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(notFoundNoticeParams), layout.None)
			}
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(internalServerErrorNoticeParams), layout.None)
		}
		exitrm, err := qtx.GetRoom(context.Background(), in.LinkID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				c.Append(header.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(notFoundNoticeParams), layout.None)
			}
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(internalServerErrorNoticeParams), layout.None)
		}

		if err := room.Link(room.LinkParams{
			Queries:   qtx,
			ID:        rid,
			To:        in.LinkID,
			Direction: in.Direction,
			TwoWay:    in.TwoWay,
		}); err != nil {
			if err == room.ErrLinkSelf {
				c.Status(fiber.StatusBadRequest)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		rm, err = qtx.GetRoom(context.Background(), rm.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				c.Append(header.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
					SectionID:    sectionID,
					SectionClass: "pt-2",
					NoticeText: []string{
						"The room you're looking for no longer exists.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				}), layout.None)
			}
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layout.None)
		}
		exitrm, err = qtx.GetRoom(context.Background(), exitrm.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				c.Append(header.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
					SectionID:    sectionID,
					SectionClass: "pt-2",
					NoticeText: []string{
						"The room you're looking for no longer exists.",
					},
					RefreshButton: true,
					NoticeIcon:    true,
				}), layout.None)
			}
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layout.None)
		}

		graph, err := room.BuildGraph(room.BuildGraphParams{
			Queries: qtx,
			Room:    &rm,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layout.None)
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layout.None)
		}

		grid := graph.BindMatrix(room.BindMatrixParams{
			Matrix:  room.EmptyBindMatrix(5),
			Row:     2,
			Col:     2,
			Shallow: false,
		})
		grid = room.AnnotateMatrixExits(grid)

		c.Status(fiber.StatusOK)
		b := graph.BindExit(in.Direction)
		b["Exits"] = graph.BindExits()
		b["RoomGrid"] = grid
		return c.Render(partial.EditRoomExitEdit, b, layout.EditRoomExitsSelect)
	}
}

func ClearRoomExit(i *service.Interfaces) fiber.Handler {
	// TODO: Get constant for common section IDs
	const sectionID string = "edit-room-exits-edit-error"

	internalServerErrorNoticeParams := partial.BindNoticeSectionParams{
		SectionID:    sectionID,
		SectionClass: "pt-2",
		NoticeText: []string{
			"Something's gone terribly wrong.",
		},
		RefreshButton: true,
		NoticeIcon:    true,
	}

	sessionExpiredNoticeParams := partial.BindNoticeSectionParams{
		SectionID:    sectionID,
		SectionClass: "pt-2",
		NoticeText: []string{
			"It looks like your session may have expired.",
		},
		RefreshButton: true,
		NoticeIcon:    true,
	}

	noPermissionNoticeParams := partial.BindNoticeSectionParams{
		SectionID:    sectionID,
		SectionClass: "pt-2",
		NoticeText: []string{
			"You don't have the permission required to edit this room.",
		},
		RefreshButton: true,
		NoticeIcon:    true,
	}

	notFoundNoticeParams := partial.BindNoticeSectionParams{
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
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(sessionExpiredNoticeParams), layout.None)
		}

		rid, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(internalServerErrorNoticeParams), layout.None)
		}

		dir := c.Params("exit")
		if !room.IsDirectionValid(dir) {
			c.Status(fiber.StatusBadRequest)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(internalServerErrorNoticeParams), layout.None)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusForbidden)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(internalServerErrorNoticeParams), layout.None)
		}

		if !perms.HasPermission(player.PermissionCreateRoom.Name) {
			c.Status(fiber.StatusForbidden)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(noPermissionNoticeParams), layout.None)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(internalServerErrorNoticeParams), layout.None)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		rm, err := qtx.GetRoom(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				c.Append(header.HXAcceptable, "true")
				return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(notFoundNoticeParams), layout.None)
			}
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(internalServerErrorNoticeParams), layout.None)
		}

		exitID := room.ExitID(&rm, dir)
		exitrm, err := qtx.GetRoom(context.Background(), exitID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				c.Append(header.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(notFoundNoticeParams), layout.None)
			}
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(internalServerErrorNoticeParams), layout.None)
		}
		exitDir, err := room.ExitDirection(&exitrm, rid)
		if err != nil && err != room.ErrExitIDNotFound {
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(internalServerErrorNoticeParams), layout.None)
		}
		if err != room.ErrExitIDNotFound {
			if err := room.Unlink(room.UnlinkParams{
				Queries:   qtx,
				ID:        exitID,
				Direction: exitDir,
			}); err != nil {
				c.Status(fiber.StatusInternalServerError)
				c.Append(header.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(internalServerErrorNoticeParams), layout.None)
			}
		}

		if err := room.Unlink(room.UnlinkParams{
			Queries:   qtx,
			ID:        rid,
			Direction: dir,
		}); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		rm, err = qtx.GetRoom(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				c.Append(header.HXAcceptable, "true")
				c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
				return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(notFoundNoticeParams), layout.None)
			}
		}

		graph, err := room.BuildGraph(room.BuildGraphParams{
			Queries: qtx,
			Room:    &rm,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(internalServerErrorNoticeParams), layout.None)
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Append(header.HXAcceptable, "true")
			c.Append("HX-Retarget", util.PrependHTMLID(sectionID))
			return c.Render(partial.NoticeSectionError, partial.BindNoticeSection(partial.BindNoticeSectionParams{
				SectionID:    sectionID,
				SectionClass: "pt-2",
				NoticeText: []string{
					"Something's gone terribly wrong.",
				},
				RefreshButton: true,
				NoticeIcon:    true,
			}), layout.None)
		}

		grid := graph.BindMatrix(room.BindMatrixParams{
			Matrix:  room.EmptyBindMatrix(5),
			Row:     2,
			Col:     2,
			Shallow: false,
		})
		grid = room.AnnotateMatrixExits(grid)

		c.Status(fiber.StatusOK)
		b := graph.BindEmptyExit(dir)
		b["Exits"] = graph.BindExits()
		b["RoomGrid"] = grid
		return c.Render(partial.EditRoomExitEdit, b, layout.EditRoomExitsSelect)
	}
}

func EditRoomTitle(i *service.Interfaces) fiber.Handler {
	type input struct {
		Title string `form:"title"`
	}

	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !room.IsTitleValid(in.Title) {
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
		if !perms.HasPermission(player.PermissionCreateRoom.Name) {
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

		rm, err := qtx.GetRoom(context.Background(), rmid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if rm.Title == in.Title {
			c.Status(fiber.StatusConflict)
			return nil
		}

		if err := qtx.UpdateRoomTitle(context.Background(), query.UpdateRoomTitleParams{
			ID:    rm.ID,
			Title: in.Title,
		}); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		rm, err = qtx.GetRoom(context.Background(), rmid)
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
			"Title":    room.TitleWithID(rm.Title, rm.ID),
			"SubTitle": "Update room properties here",
		}
		b["Title"] = rm.Title
		b["TitlePath"] = route.RoomTitlePath(rmid)
		b["NoticeSection"] = partial.BindNoticeSection(partial.BindNoticeSectionParams{
			Success:      true,
			SectionID:    "room-edit-title-notice",
			SectionClass: "pb-2",
			NoticeText: []string{
				"Success! The room title has been updated.",
			},
			NoticeIcon: true,
		})
		return c.Render(partial.RoomEditTitle, b, layout.PageHeader)
	}
}

func EditRoomDescription(i *service.Interfaces) fiber.Handler {
	type input struct {
		Description string `form:"desc"`
	}

	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !room.IsDescriptionValid(in.Description) {
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
		if !perms.HasPermission(player.PermissionCreateRoom.Name) {
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

		rm, err := qtx.GetRoom(context.Background(), rmid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		// TODO: Add Conflict tests for edit room title, description and size
		if rm.Description == in.Description {
			c.Status(fiber.StatusConflict)
			return nil
		}

		if err := qtx.UpdateRoomDescription(context.Background(), query.UpdateRoomDescriptionParams{
			ID:          rm.ID,
			Description: in.Description,
		}); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		rm, err = qtx.GetRoom(context.Background(), rmid)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		b := fiber.Map{}
		b["Description"] = rm.Description
		b["DescriptionPath"] = route.RoomDescriptionPath(rmid)
		b["NoticeSection"] = partial.BindNoticeSection(partial.BindNoticeSectionParams{
			Success:      true,
			SectionID:    "room-edit-title-notice",
			SectionClass: "pb-2",
			NoticeText: []string{
				"Success! The room description has been updated.",
			},
			NoticeIcon: true,
		})
		return c.Render(partial.RoomEditDescription, b, layout.None)
	}
}

func EditRoomSize(i *service.Interfaces) fiber.Handler {
	type input struct {
		Size int32 `form:"size"`
	}

	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !room.IsSizeValid(in.Size) {
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
		if !perms.HasPermission(player.PermissionCreateRoom.Name) {
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

		rm, err := qtx.GetRoom(context.Background(), rmid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		// TODO: Add Conflict tests for edit room title, description and size
		if rm.Size == in.Size {
			c.Status(fiber.StatusConflict)
			return nil
		}

		if err := qtx.UpdateRoomSize(context.Background(), query.UpdateRoomSizeParams{
			ID:   rm.ID,
			Size: in.Size,
		}); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		rm, err = qtx.GetRoom(context.Background(), rmid)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		b := fiber.Map{}
		b["Size"] = rm.Size
		b["SizePath"] = route.RoomSizePath(rmid)
		b = room.BindSizeRadioGroup(b, &rm)
		b["NoticeSection"] = partial.BindNoticeSection(partial.BindNoticeSectionParams{
			Success:      true,
			SectionID:    "room-edit-title-notice",
			SectionClass: "pb-2",
			NoticeText: []string{
				"Success! The room size has been updated.",
			},
			NoticeIcon: true,
		})
		return c.Render(partial.RoomEditSize, b, layout.None)
	}
}
