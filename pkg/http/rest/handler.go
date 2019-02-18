package rest

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"

	"strings"

	sev "github.com/valentyn88/streamer_event_viewer"
	"github.com/valentyn88/streamer_event_viewer/storage"
)

type Handler struct {
	Storage     storage.Storager
	CookieStore *sessions.CookieStore
	Oauth2Cnfg  *oauth2.Config
}

func (h Handler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	session, err := h.CookieStore.Get(r, oauthSessionName)
	if err != nil {
		log.Printf("session error: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var (
		userPart = `<a href="/subscribe-form">Subscribe</a></br>
					<a href="/livestream">Livestream</a></br>
					<a href="/logout">Logout</a>`
		loginLink = ""
	)

	u, ok := session.Values[userKey].(*sev.User)
	if !ok || u.ID == "" {
		userPart = ""
		loginLink = `<a href="/login">Login using Twitch</a></br>`
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body := `<html>
				<body>
					%s
					%s
				</body>
			</html>`

	w.Write([]byte(fmt.Sprintf(body, loginLink, userPart)))
}

func (h Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var tokenBytes [255]byte
	if _, err := rand.Read(tokenBytes[:]); err != nil {
		log.Printf("couldn't generate a session error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	state := hex.EncodeToString(tokenBytes[:])

	http.Redirect(w, r, h.Oauth2Cnfg.AuthCodeURL(state), http.StatusTemporaryRedirect)

	return
}

func (h Handler) HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	session, err := h.CookieStore.Get(r, oauthSessionName)
	if err != nil {
		log.Printf("session error: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := h.Oauth2Cnfg.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		fmt.Printf("couldn't get token error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		log.Printf("couldn't create request to twitch user API error %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("couldn't make request to twitch user API error %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("couldn't read twitch user API body error: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var data sev.Data
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("couldn't parse twitch API response error: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// add the oauth token to session
	session.Values[oauthTokenKey] = token

	// save user to session
	if len(data.Users) > 0 {
		user := data.Users[0]
		session.Values[userKey] = &user
	}

	if err := session.Save(r, w); err != nil {
		log.Printf("saving session error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("Access token: %s\n", token.AccessToken)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

	return
}

func (h Handler) HandleSubscribeForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	html := `<html>
		<body>
			<form method="GET" action="/subscribe">
				<label for="name">Favorite Twitch streamer name: </label>
    			<input type="text" name="name" required>
				<input type="submit" value="Subscribe!">
			</form>
		</body>
	</html>`
	w.Write([]byte(html))
}

func (h Handler) HandleSubscribe(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Streamer name must be provided"))
		return
	}

	session, err := h.CookieStore.Get(r, oauthSessionName)
	if err != nil {
		log.Printf("session error: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := http.Client{}
	url := fmt.Sprintf("https://api.twitch.tv/helix/users?login=%s", name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("couldn't create request to twitch user API error %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, ok := session.Values[oauthTokenKey].(*oauth2.Token)
	if !ok {
		log.Printf("couldn't get token from session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("couldn't make request to twitch user API error %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var data sev.Data
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("couldn't read twitch user API body error: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("couldn't parse twitch API response error: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(data.Users) == 0 {
		log.Printf("streamer data is empty")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Your favorite streamer user is not exists. Probable you made a mistake in the name?"))
		return
	}

	streamer := data.Users[0]
	session.Values[favoriteStreamerKey] = &sev.Streamer{ID: streamer.ID, Login: streamer.Login}

	if err = session.Save(r, w); err != nil {
		log.Printf("error saving session: %s", err)
		err = nil
		return
	}

	client = http.Client{}
	subscBody := sev.SubscriptionBody{
		Callback: subscriptionEventsURL,
		Mode:     "subscribe",
		Topic:    fmt.Sprintf("https://api.twitch.tv/helix/streams?user_id=%s", streamer.ID),
		Seconds:  864000,
	}
	bb, errc := json.Marshal(&subscBody)
	if errc != nil {
		log.Printf("couldn't Marchal subscription body error: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bReader := bytes.NewReader(bb)
	req, err = http.NewRequest("POST", "https://api.twitch.tv/helix/webhooks/hub", bReader)
	if err != nil {
		log.Printf("couldn't create request to twitch user API error %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	resp, err = client.Do(req)
	if err != nil {
		log.Printf("couldn't make request to twitch user API error %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body = []byte(fmt.Sprintf(`<html>
		<body>
			You were successfully subscribed on %s events! Live stream <a href="/livestream">link</a>
		</body>
		</html>`, name))
	w.Write(body)
}

func (h Handler) HandleSubscriptionEvents(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("couldn't parse request body %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Storage.Save(body)
}

func (h Handler) HandleLivestream(w http.ResponseWriter, r *http.Request) {
	session, err := h.CookieStore.Get(r, oauthSessionName)
	if err != nil {
		log.Printf("session error: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	streamer, ok := session.Values[favoriteStreamerKey].(*sev.Streamer)
	if !ok {
		log.Println("couldn't get streamer from session")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("couldn't get streamer from session"))
		return
	}

	events := h.Storage.Last(10)
	eventsStr := []string{}
	for _, e := range events {
		eventsStr = append(eventsStr, fmt.Sprintf("%v", e))
	}
	if len(eventsStr) == 0 {
		eventsStr = append(eventsStr, "No events")
	}

	table := `<table border="0">
				<tr>
					<td>%s</td>
					<td>%s</td>
					<td>%s</td>
				</tr>
			</table>`

	videoIFrame := `<iframe
    				src="https://player.twitch.tv/?channel=%s&muted=true"
    				height="400"
    				width="300"
    				frameborder="0"
    				scrolling="no"
    				allowfullscreen="true">
				</iframe>`

	chatIframe := `<iframe frameborder="0"
						scrolling="no"
						id="chat_embed"
						src="https://www.twitch.tv/embed/%s/chat"
						height="400"
						width="300">
				</iframe>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(table,
		fmt.Sprintf(videoIFrame, streamer.Login),
		fmt.Sprintf(chatIframe, streamer.Login),
		strings.Join(eventsStr, "\n"))))
	return
}

func (h Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := h.CookieStore.Get(r, oauthSessionName)
	if err != nil {
		log.Printf("session error: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values[userKey] = &sev.User{}
	session.Values[favoriteStreamerKey] = &sev.Streamer{}
	session.Values[oauthTokenKey] = nil

	if err := session.Save(r, w); err != nil {
		log.Printf("couldn't save session error: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
