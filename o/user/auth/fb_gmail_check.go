package auth

import (
	"encoding/json"
	"errors"
	"g/x/web"
	"net/http"
)

type FacebookAuth struct {
	ID    string
	Token string
}

type GmailAuth struct {
	GmId  string
	Token string
}

func (a FacebookAuth) IsAuthenticated() error {
	if a.ID == "" {
		return web.BadRequest("mising facebook id")
	}

	var c = http.Client{}
	req, _ := http.NewRequest("GET", "https://graph.facebook.com/v2.6/me", nil)
	q := req.URL.Query()
	q.Add("access_token", a.Token)
	req.URL.RawQuery = q.Encode()
	var resp, err = c.Do(req)
	if err != nil {
		return web.Unauthorized("check facebook token error")
	}
	var body = struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return web.Unauthorized("parse facebook auth error")
	}
	if body.ID != a.ID {
		return web.Unauthorized("facebook id and token mismatch")
	}
	return nil
}

func (a FacebookAuth) IsPresent() bool {
	return len(a.ID) > 2
}

func (a FacebookAuth) IsValid() error {
	if len(a.Token) < 10 {
		return errors.New("Facebook Token is required")
	}
	return nil
}

func (a GmailAuth) IsAuthenticated() error {
	if a.GmId == "" {
		return web.BadRequest("mising facebook id")
	}
	var c = http.Client{}
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/tokeninfo", nil)
	q := req.URL.Query()
	q.Add("id_token", a.Token)
	req.URL.RawQuery = q.Encode()
	var resp, err = c.Do(req)
	if err != nil {
		return web.Unauthorized("checkgmail token error")
	}
	var body = struct {
		Id      string `json:"sub"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return web.Unauthorized("parse checkgmail auth error")
	}
	if body.Id != a.GmId {
		return web.Unauthorized("checkgmail id and token mismatch")
	}
	return nil
}
