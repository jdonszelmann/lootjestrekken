package store

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"lootjestrekken/pkg/lootjestrekken"
	"sync"
)

type InMemoryStore struct {
	trekkingen map[string]lootjestrekken.Trekking
	sync.Mutex
}

func (i *InMemoryStore) UpdateTrekking(trekking lootjestrekken.Trekking) error {
	i.Lock()
	defer i.Unlock()

	log.Debugf("updating trekking with name %s in store", trekking.Name)

	i.trekkingen[trekking.Name] = trekking
	return nil
}

func (i *InMemoryStore) GetTrekking(name string) (lootjestrekken.Trekking, error) {
	i.Lock()
	defer i.Unlock()

	log.Debugf("getting trekking with name %s from store", name)

	res, ok := i.trekkingen[name]
	if !ok {
		return lootjestrekken.Trekking{}, errors.New("trekking name not found")
	}

	return res, nil
}

func (i *InMemoryStore) GetTrekkingNames() ([]string, error) {
	i.Lock()
	defer i.Unlock()

	log.Debug("getting all trekking names from store")

	keys := make([]string, 0, len(i.trekkingen))
	for k := range i.trekkingen {
		keys = append(keys, k)
	}
	return keys, nil
}

func (i *InMemoryStore) GetTrekkingInfos() ([]string, error) {
	i.Lock()
	defer i.Unlock()

	log.Debug("getting all trekking infos from store")

	keys := make([]string, 0, len(i.trekkingen))
	for _, v := range i.trekkingen {
		keys = append(keys, v.GetInfo())
	}
	return keys, nil
}

func (i *InMemoryStore) AddTrekking(name string, trekking lootjestrekken.Trekking) error {
	i.Lock()
	defer i.Unlock()

	log.Debugf("Adding trekking with name %s to store", name)

	trekking.Name = name

	if _, ok := i.trekkingen[name]; ok {
		return ErrExists
	}

	i.trekkingen[name] = trekking

	return nil
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore {
		trekkingen: map[string]lootjestrekken.Trekking{},
	}
}


