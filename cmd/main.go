package main

import (
	"context"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	. "lootjestrekken/cmd/handler"
	"lootjestrekken/cmd/store"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var (
	address = flag.String("address", "0.0.0.0", "Address to serve on")
	port = flag.Int("port", 8080, "Port to serve on")
	storetype = flag.String("store", "inmemory", "store type: [inmemory, db]")
	dbloc = flag.String("location", "./data", "db location")
)

func Home(w http.ResponseWriter, r *http.Request) {
	text := `
<head>
<style>
pre {
	font-family: "monospace";
}
</style>
</head>
<body>
<pre>
Welcome to LootjesTrekken!

use /t                                           to list ongoing trekkingen
use /t/{trekking-name}/add                       to start a new trekking with this name
use /t/{trekking-name}/people                    to list people in a trekking
use /t/{trekking-name}/people/{name}/add         to add a person to a trekking with this name
use /t/{trekking-name}/people/{name}/remove      to remove a person from a trekking with this name
use /t/{trekking-name}/trek                      to trek this trekking
use /t/{trekking-name}/people/{name}/getrokken   to see who you have getrokken

</pre>
</body>
`

	_, err := w.Write([]byte(text))
	if err != nil {
		log.Errorf("Couldn't write %v", err)
	}
}

func init() {
	lvlstring := os.Getenv("LOG_LEVEL")

	loglevel, err := log.ParseLevel(lvlstring)
	if err != nil {
		loglevel = log.DebugLevel
	}

	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(loglevel)

	// log error after the logger is initialised
	if err != nil && lvlstring != "" {
		log.Errorf("loglevel string %s could not be parsed, defaulting to Info: %v", lvlstring, err)
	}
}

func getStore(storetype, location string) (store.Store, error) {
	switch storetype {
	default:
		log.Error("Unexpected value for store type: %s", storetype)
		fallthrough
	case "inmemory":
		log.Infof("Using in memory data store")
		return store.NewInMemoryStore(), nil
	case "db":
		log.Infof("Using persistent data store at %s", location)
		return store.NewDbStore(location)
	}
}

func runServer(ctx context.Context, address string, port int, storetype, dbloc string) {
	r := mux.NewRouter()
	s, err := getStore(storetype, dbloc)
	if err != nil {
		log.Fatalf("Couldn't get db connection: %v", err)
	}

	h := Handler{
		Store: s,
	}

	r.StrictSlash(true)

	r.HandleFunc("/", Home)
	r.HandleFunc("/t", h.ListTrekkingen)
	r.HandleFunc("/t/{trekking-name}/add", h.NewTrekking)
	r.HandleFunc("/t/{trekking-name}/raw", h.RawTrekking)
	r.HandleFunc("/t/{trekking-name}/people", h.GetPeople)
	r.HandleFunc("/t/{trekking-name}", h.GetPeople)
	r.HandleFunc("/t/{trekking-name}/people/{name}/add", h.AddPerson)
	r.HandleFunc("/t/{trekking-name}/people/{name}/remove", h.RemovePerson)
	r.HandleFunc("/t/{trekking-name}/trek", h.Trek)
	r.HandleFunc("/t/{trekking-name}/people/{name}/getrokken", h.Getrokken)

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%s:%d", address, port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func(){
		log.Infof("Running server on port %d!", port)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed{
			log.Fatal(err)
		}
	}()



	for {
		select {
		case <-ctx.Done():
			err := srv.Close()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func main() {
	flag.Parse()
	runServer(context.Background(), *address, *port, *storetype, *dbloc)
}
