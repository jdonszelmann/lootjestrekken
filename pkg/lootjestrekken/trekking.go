// + json1

package lootjestrekken

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math/rand"
)

func lpad(s string, pad string, plength int) string {
	for i := len(s); i < plength; i++ {
		s = pad + s
	}
	return s
}

type Trekking struct {
	gorm.Model
	People        []string
	PeopleMapping []string
	Getrokken     bool
	Name          string
}

func (t *Trekking) AddPerson(name string) {
	t.People = append(t.People, name)
}

func (t *Trekking) GetInfo() string {
	if t.Getrokken {
		return fmt.Sprintf("%s getrokken", lpad(t.Name, " ", 30))
	} else {
		return t.Name
	}
}

func (t *Trekking) RemovePerson(name string) {
	found := false
	for index, i := range t.People {
		if found {
			t.People[index-1] = t.People[index]
		} else if i == name {
			found = true
		}
	}

	t.People = t.People[:len(t.People)-1]
}

func (t *Trekking) Trek() {
	t.Getrokken = true

	// First shuffle the people
	rand.Shuffle(len(t.People), func(i, j int) { t.People[i], t.People[j] = t.People[j], t.People[i] })

	// then derange them into a second array
	t.PeopleMapping = Derange(t.People)
}

func (t *Trekking) GetrokkenPerson(name string) (string, error) {
	for index, i := range t.People {
		if i == name {
			return t.PeopleMapping[index], nil
		}
	}

	return "", errors.New("not part of trekking")
}

func Derange(arr []string) []string{
	newarr := make([]string, len(arr))
	copy(newarr, arr)

	for i := len(arr)-1; i >= 1; i-- {
		j := rand.Intn(i)
		newarr[i], newarr[j] = newarr[j], newarr[i]
	}

	return newarr
}