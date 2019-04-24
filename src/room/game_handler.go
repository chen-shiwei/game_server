package room

import (
	"errors"
)

type GameRoomEntry struct {
	Name        string
	GameId      string
	Handler     GameHandler
	MemberLimit uint16
}

var (
	entries             = make(map[string]*GameRoomEntry)
	ErrorDupGameId      = errors.New("gameId dup")
	ErrorEmptyHandler   = errors.New("empty handler")
	ErrorEntryNotExists = errors.New("entry does not exists")
)

func RegisterGameEntry(name string, gameId string, handler GameHandler, limit uint16) error {
	if _, ok := entries[gameId]; ok {
		return ErrorDupGameId
	}
	if handler == nil {
		return ErrorEmptyHandler
	}
	if limit == 0 {
		limit = 3
	}
	entries[gameId] = &GameRoomEntry{Name: name, GameId: gameId, Handler: handler, MemberLimit: limit}
	return nil
}

func GetGameRoomEntry(gamdId string) (*GameRoomEntry, error) {
	if entry, ok := entries[gamdId]; ok {
		return entry, nil
	}
	return nil, ErrorEntryNotExists
}
