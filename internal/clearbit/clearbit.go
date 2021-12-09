package clearbit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var RetryError = errors.New("please retry later")

type (
	API interface {
		LookupPerson(email string) (Person, error)
	}
)

type clearbitClient struct {
	baseUrl    string
	apiKey     string
	httpClient *http.Client
	// TODO implement caching (would need to subscribe to webhooks)
}

func NewClearbitAPI(apiKey string) API {
	return &clearbitClient{
		baseUrl: "https://person.clearbit.com",
		apiKey:  apiKey,
		// TODO customizable http client
		httpClient: http.DefaultClient,
	}
}

// TODO support optional query parameters
func (c *clearbitClient) LookupPerson(email string) (Person, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/people/find?email=%s", c.baseUrl, email), nil)
	if err != nil {
		return Person{}, err
	}
	req.SetBasicAuth(c.apiKey, "")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Person{}, err
	} else if resp == nil {
		return Person{}, fmt.Errorf("no response")
	}

	if resp.StatusCode >= 300 {
		body, bodyErr := ioutil.ReadAll(resp.Body)
		if bodyErr != nil {
			return Person{}, fmt.Errorf("error reading response body: %v", bodyErr)
		}
		return Person{}, fmt.Errorf("%d: %s", resp.StatusCode, string(body))
	} else if resp.StatusCode == http.StatusAccepted {
		// clearbit returns an immediate http.StatusAccepted when they need to do an asynchronous lookup
		return Person{}, RetryError
	}
	var p Person
	if err = json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return p, fmt.Errorf("error decoding respons json: %v", err)
	}
	return p, nil
}
