package store

import (
	"errors"
	"lootjestrekken/pkg/lootjestrekken"
)

var ErrExists = errors.New("name already exists")

type Store interface {
	AddTrekking(name string, trekking lootjestrekken.Trekking) error
	GetTrekkingNames() ([]string, error)
	GetTrekkingInfos() ([]string, error)
	GetTrekking(name string) (lootjestrekken.Trekking, error)
	UpdateTrekking(trekking lootjestrekken.Trekking) error
}

