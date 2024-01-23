package rooms

import (
	"context"
	"errors"

	_ "github.com/go-sql-driver/mysql"

	"petrichormud.com/app/internal/queries"
)

const (
	errInvalidDirection string = "invalid direction"
	errLinkSelf         string = "cannot link a room to itself"
)

var (
	ErrInvalidDirection error = errors.New(errInvalidDirection)
	ErrLinkSelf         error = errors.New(errLinkSelf)
)

type LinkParams struct {
	Queries   *queries.Queries
	Direction string
	To        int64
	ID        int64
	TwoWay    bool
}

func Link(in LinkParams) error {
	if !IsDirectionValid(in.Direction) {
		return ErrInvalidDirection
	}

	if in.To == in.ID {
		return ErrLinkSelf
	}

	switch in.Direction {
	case DirectionNorth:
		if err := in.Queries.UpdateRoomExitNorth(context.Background(), queries.UpdateRoomExitNorthParams{
			ID:    in.ID,
			North: in.To,
		}); err != nil {
			return err
		}
	case DirectionNortheast:
		if err := in.Queries.UpdateRoomExitNortheast(context.Background(), queries.UpdateRoomExitNortheastParams{
			ID:        in.ID,
			Northeast: in.To,
		}); err != nil {
			return err
		}
	case DirectionEast:
		if err := in.Queries.UpdateRoomExitEast(context.Background(), queries.UpdateRoomExitEastParams{
			ID:   in.ID,
			East: in.To,
		}); err != nil {
			return err
		}
	case DirectionSoutheast:
		if err := in.Queries.UpdateRoomExitSoutheast(context.Background(), queries.UpdateRoomExitSoutheastParams{
			ID:        in.ID,
			Southeast: in.To,
		}); err != nil {
			return err
		}
	case DirectionSouth:
		if err := in.Queries.UpdateRoomExitSouth(context.Background(), queries.UpdateRoomExitSouthParams{
			ID:    in.ID,
			South: in.To,
		}); err != nil {
			return err
		}
	case DirectionSouthwest:
		if err := in.Queries.UpdateRoomExitSouthwest(context.Background(), queries.UpdateRoomExitSouthwestParams{
			ID:        in.ID,
			Southwest: in.To,
		}); err != nil {
			return err
		}
	case DirectionWest:
		if err := in.Queries.UpdateRoomExitWest(context.Background(), queries.UpdateRoomExitWestParams{
			ID:   in.ID,
			West: in.To,
		}); err != nil {
			return err
		}
	case DirectionNorthwest:
		if err := in.Queries.UpdateRoomExitNorthwest(context.Background(), queries.UpdateRoomExitNorthwestParams{
			ID:        in.ID,
			Northwest: in.To,
		}); err != nil {
			return err
		}
	default:
		return ErrInvalidDirection
	}

	if in.TwoWay {
		if err := Link(LinkParams{
			Queries:   in.Queries,
			ID:        in.To,
			To:        in.ID,
			TwoWay:    false,
			Direction: DirectionOpposite(in.Direction),
		}); err != nil {
			return err
		}
	}

	return nil
}

type UnlinkParams struct {
	Queries   *queries.Queries
	Direction string
	ID        int64
}

func Unlink(in UnlinkParams) error {
	if !IsDirectionValid(in.Direction) {
		return ErrInvalidDirection
	}

	switch in.Direction {
	case DirectionNorth:
		if err := in.Queries.UpdateRoomExitNorth(context.Background(), queries.UpdateRoomExitNorthParams{
			ID:    in.ID,
			North: 0,
		}); err != nil {
			return err
		}
	case DirectionNortheast:
		if err := in.Queries.UpdateRoomExitNortheast(context.Background(), queries.UpdateRoomExitNortheastParams{
			ID:        in.ID,
			Northeast: 0,
		}); err != nil {
			return err
		}
	case DirectionEast:
		if err := in.Queries.UpdateRoomExitEast(context.Background(), queries.UpdateRoomExitEastParams{
			ID:   in.ID,
			East: 0,
		}); err != nil {
			return err
		}
	case DirectionSoutheast:
		if err := in.Queries.UpdateRoomExitSoutheast(context.Background(), queries.UpdateRoomExitSoutheastParams{
			ID:        in.ID,
			Southeast: 0,
		}); err != nil {
			return err
		}
	case DirectionSouth:
		if err := in.Queries.UpdateRoomExitSouth(context.Background(), queries.UpdateRoomExitSouthParams{
			ID:    in.ID,
			South: 0,
		}); err != nil {
			return err
		}
	case DirectionSouthwest:
		if err := in.Queries.UpdateRoomExitSouthwest(context.Background(), queries.UpdateRoomExitSouthwestParams{
			ID:        in.ID,
			Southwest: 0,
		}); err != nil {
			return err
		}
	case DirectionWest:
		if err := in.Queries.UpdateRoomExitWest(context.Background(), queries.UpdateRoomExitWestParams{
			ID:   in.ID,
			West: 0,
		}); err != nil {
			return err
		}
	case DirectionNorthwest:
		if err := in.Queries.UpdateRoomExitNorthwest(context.Background(), queries.UpdateRoomExitNorthwestParams{
			ID:        in.ID,
			Northwest: 0,
		}); err != nil {
			return err
		}
	default:
		return ErrInvalidDirection
	}

	return nil
}
