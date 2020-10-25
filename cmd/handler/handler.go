package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"lootjestrekken/cmd/store"
	"lootjestrekken/pkg/lootjestrekken"
	"net/http"
	"strings"
)

type Handler struct {
	Store store.Store
}

func (h *Handler) NewTrekking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["trekking-name"]
	if name == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Errorf("name variable was empty")
		return
	}

	log.Debugf("Creating new trekking with name %s", name)

	if err := h.Store.AddTrekking(name, lootjestrekken.Trekking{}); err != nil {
		if err == store.ErrExists {
			http.Error(w, "Couldn't create trekking because trekking with this name already exists", http.StatusConflict)
		} else {
			http.Error(w, "Couldn't create trekking", http.StatusInternalServerError)
		}
		return
	}

	_, err := w.Write([]byte(fmt.Sprintf("New trekking created with name %s", name)))
	if err != nil {
		log.Errorf("Couldn't write %v", err)
	}
}

func (h *Handler) ListTrekkingen(w http.ResponseWriter, r *http.Request) {
	log.Debug("listing all trekkingen")

	names, err := h.Store.GetTrekkingNames()
	if err != nil {
		http.Error(w, "Couldn't read trekkingen", http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(strings.Join(names, "\n")))
	if err != nil {
		log.Errorf("Couldn't write %v", err)
	}
}

func (h *Handler) RawTrekking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["trekking-name"]
	if name == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Debug("getting raw trekking named %s", name)

	trekking, err := h.Store.GetTrekking(name)

	if err != nil {
		http.Error(w, "Couldn't find trekking", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(trekking)
	if err != nil {
		log.Printf("Couldn't write %v", err)
	}
}

func (h *Handler) GetPeople(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["trekking-name"]
	if name == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Debug("getting people associated with trekking %s", name)

	trekking, err := h.Store.GetTrekking(name)

	if err != nil {
		http.Error(w, "Couldn't find trekking", http.StatusNotFound)
		return
	}

	people := trekking.People

	_, err = w.Write([]byte(strings.Join(people, "\n")))
	if err != nil {
		log.Printf("Couldn't write %v", err)
	}
}

func (h *Handler) AddPerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trekkingname := vars["trekking-name"]
	personname := vars["name"]
	if trekkingname == "" || personname == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Debugf("Adding person %s to trekking %s", personname, trekkingname)

	trekking, err := h.Store.GetTrekking(trekkingname)

	if err != nil {
		http.Error(w, "Couldn't find trekking", http.StatusNotFound)
		return
	}

	if trekking.Getrokken {
		http.Error(w, "This trekking is already getrokken", http.StatusConflict)
		return
	}

	trekking.AddPerson(personname)
	err = h.Store.UpdateTrekking(trekking)
	if err != nil {
		http.Error(w, "Failed to add person to trekking", http.StatusBadRequest)
		return
	}

	_, err = w.Write([]byte("Added succesfully"))
	if err != nil {
		log.Printf("Couldn't write %v", err)
	}
}

func (h *Handler) RemovePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trekkingname := vars["trekking-name"]
	personname := vars["name"]
	if trekkingname == "" || personname == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Debugf("Removing person %s from trekking %s", personname, trekkingname)

	trekking, err := h.Store.GetTrekking(trekkingname)

	if err != nil {
		http.Error(w, "Couldn't find trekking", http.StatusNotFound)
		return
	}

	if trekking.Getrokken {
		http.Error(w, "This trekking is already getrokken", http.StatusConflict)
		return
	}

	trekking.RemovePerson(personname)
	err = h.Store.UpdateTrekking(trekking)
	if err != nil {
		http.Error(w, "Failed to remove person to trekking", http.StatusBadRequest)
		return
	}

	_, err = w.Write([]byte("Removed succesfully"))
	if err != nil {
		log.Printf("Couldn't write %v", err)
	}
}

func (h *Handler) Trek(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["trekking-name"]

	log.Debugf("Initiating trek on trekking with name %s", name)

	trekking, err := h.Store.GetTrekking(name)
	if err != nil {
		http.Error(w, "Couldn't find trekking", http.StatusNotFound)
		return
	}

	trekking.Trek()
	err = h.Store.UpdateTrekking(trekking)
	if err != nil {
		http.Error(w, "Failed to trek trekking", http.StatusBadRequest)
		return
	}

	_, err = w.Write([]byte("Trekking successfully getrokken. "))
	if err != nil {
		log.Printf("Couldn't write %v", err)
	}
}

func (h *Handler) Getrokken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trekkingname := vars["trekking-name"]
	personname := vars["name"]
	if trekkingname == "" || personname == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Debugf("Getting getrokken person for %s in treking %s", personname, trekkingname)

	trekking, err := h.Store.GetTrekking(trekkingname)
	if err != nil {
		http.Error(w, "Couldn't find trekking", http.StatusNotFound)
		return
	}

	if !trekking.Getrokken {
		http.Error(w, "This trekking is not yet getrokken", http.StatusConflict)
		return
	}

	getrokken, err := trekking.GetrokkenPerson(personname)
	if err != nil {
		http.Error(w, "You are not part of this trekking", http.StatusNotFound)
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf("You have getrokken: %s", getrokken)))
	if err != nil {
		log.Printf("Couldn't write %v", err)
	}
}