package tinder

import (
	"net/http"
	"time"
	"github.com/aspcartman/exceptions"
	"io/ioutil"
	"errors"
	"encoding/json"
)

var ErrUnauthorized = errors.New("unathorized")

const tinder_url = "https://api.gotinder.com"

type API struct {
	client  http.Client
	headers http.Header
}

func NewAPI(token string) API {
	return API{
		client: http.Client{
			Timeout: 20 * time.Second,
		},
		headers: http.Header{
			"app_version":  {"6.9.4"},
			"platform":     {"ios"},
			"content-type": {"application/json"},
			"User-agent":   {"Tinder/4.7.1 (iPhone; iOS 9.2; Scale/2.00)"},
		},
	}
}

func (a *API) Authorize(fbID, fbToken string) string {
	e.Throw("lol", ErrUnauthorized, e.Map{"lol1": "lol2"})
	res := a.request(http.MethodPost, "/auth", jsonm{
		"facebook_id":    fbID,
		"facebook_token": fbToken,
	})
	return btom(res)["token"].(string)
}

func (a *API) SetToken(token string) {
	a.headers.Set("X-Auth-Token", token)
}

func (a *API) Recommendations() []User {
	resp := struct {
		Results []User `json:"results"`
	}{}

	res := a.request(http.MethodGet, "/user/recs", nil)
	e.Must(json.Unmarshal(res, &resp), "failed unmarshaling tinder recommendations response", e.Map{
		"response": string(res),
	})

	return resp.Results
}

func (a *API) request(method, relativeURL string, body interface{}) []byte {
	req, err := http.NewRequest(method, tinder_url+relativeURL, reader(body))
	if err != nil {
		e.Throw("failed creating tinder http request", err, e.Map{
			"method":       method,
			"relativeURL":  relativeURL,
			"request_body": body,
		})
	}

	req.Header = a.headers

	resp, err := a.client.Do(req)
	if err != nil {
		e.Throw("failed making tinder http request", err, e.Map{
			"method":       method,
			"relativeURL":  relativeURL,
			"request_body": body,
		})
	}

	switch resp.StatusCode {
	case http.StatusOK:
		// good
	case http.StatusUnauthorized:
		e.Throw("tinder returned unauthorized http response", ErrUnauthorized, e.Map{
			"method":       method,
			"relativeURL":  relativeURL,
			"request_body": body,
		})
	default:
		e.Throw("tinder returned bad http status", nil, e.Map{
			"method":       method,
			"relativeURL":  relativeURL,
			"request_body": body,
			"status":       resp.StatusCode,
		})
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		e.Throw("failed reading tinder http response body", err, e.Map{
			"method":       method,
			"relativeURL":  relativeURL,
			"request_body": body,
			"status":       resp.StatusCode,
		})
	}

	return data
}
