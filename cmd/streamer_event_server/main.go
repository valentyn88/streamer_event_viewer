package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"

	sev "github.com/valentyn88/streamer_event_viewer"
	"github.com/valentyn88/streamer_event_viewer/pkg/http/rest"
	"github.com/valentyn88/streamer_event_viewer/storage"
)

var (
	clientID     = "z6xtn9ohqrz8ar5ecnn6390b2j54px"
	clientSecret = "j2kme1vp9gaf32npjjs0nrjh9gkzn8"
	scopes       = []string{"user:read:email"}
	redirectURL  = "http://localhost:7001/redirect"
	oauth2Config *oauth2.Config
)

func main() {
	log.SetOutput(os.Stdout)

	gob.Register(&oauth2.Token{})
	gob.Register(&sev.User{})
	gob.Register(&sev.Streamer{})

	oauth2Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		Endpoint:     twitch.Endpoint,
		RedirectURL:  redirectURL,
	}

	var h rest.Handler
	h.Storage = storage.NewInMemory()
	h.CookieStore = sessions.NewCookieStore([]byte("mSGZBGUrufn9L5y79Qas"))
	h.Oauth2Cnfg = oauth2Config

	http.HandleFunc("/", h.HandleRoot)
	http.HandleFunc("/login", h.HandleLogin)
	http.HandleFunc("/redirect", h.HandleOAuth2Callback)
	http.HandleFunc("/subscribe-form", h.HandleSubscribeForm)
	http.HandleFunc("/subscribe", h.HandleSubscribe)
	http.HandleFunc("/subscription-events", h.HandleSubscriptionEvents)
	http.HandleFunc("/livestream", h.HandleLivestream)
	http.HandleFunc("/logout", h.HandleLogout)

	if err := http.ListenAndServe(":7001", nil); err != nil {
		log.Fatalf("couldn't start password server error: %s", err.Error())
	}
}
