package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/y-koma0130/bookings/pkg/config"
	"github.com/y-koma0130/bookings/pkg/handlers"
	"github.com/y-koma0130/bookings/pkg/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {

	// change this to true when in produciton
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create tremplate cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)

	render.NewTemplates(&app)
	handlers.NewHandlers(repo)

	s := fmt.Sprintf("Starting application on port %s", portNumber)
	fmt.Println(s)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}
