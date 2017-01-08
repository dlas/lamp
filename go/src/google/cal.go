/* Code for interfacing with the google canedar API */

package google

import (
	"config"
	"encoding/json"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"log"
	"net/http"
	"time"
)

/* What state do we have to maintain? */
type CalendarState struct {
	ctx    context.Context
	client *http.Client
	config *oauth2.Config
}

/* Build a new calendar state object given a config. This tries to grab
 * the auth and secret tokens out of the config information */
func NewCS(c *config.Config) *CalendarState {
	var cs CalendarState

	var err error
	log.Printf("CL %v", string(c.GoogleSecret))
	cs.config, err = google.ConfigFromJSON(c.GoogleSecret, calendar.CalendarReadonlyScope)

	if err != nil {
		panic(err)
	}
	cs.ctx = context.Background()

	/* If we have an auth token, try to use it */
	if len(c.GoogleAuth) > 0 {
		var tok oauth2.Token
		err := json.Unmarshal(c.GoogleAuth, &tok)
		cs.client = cs.config.Client(cs.ctx, &tok)
		if err != nil {
			panic(err)
		}
	}

	return &cs
}

/* Do we still need to authenticate? */
func (cs *CalendarState) NeedsAuth() bool {
	return cs.client == nil

}

/* Get the URL to go to do the OAUTH2 authentication */
func (cs *CalendarState) GetAuthURL() string {
	authURL := cs.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("AU: %v", authURL)
	return authURL
}

/* Process an a token retrieved from the getauthurl() URL.
 * this will attempt to get an oauth token and will setup a
 * working calendar client. It returns a token to put in the GoogleAuth
 * field so that this can work automatically next time.
 */
func (cs *CalendarState) AuthCallback(code string) []byte {
	tok, err := cs.config.Exchange(oauth2.NoContext, code)
	if err != nil {
		panic(err)
	}
	cs.client = cs.config.Client(cs.ctx, tok)

	m, _ := json.Marshal(tok)
	return m
}

/* This function downloads the calendar and gets the next eligable
 * event.  It returns the start time for that event.
 * XXX: obviosly not... Implement me!
 */
func (cs *CalendarState) GetEvents() {
	t := time.Now().Format(time.RFC3339)
	srv, err := calendar.New(cs.client)
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	log.Printf("EV: %v ERR: %v", events, err)
}

/*
 * Find the next apointment that starts after s and before eAND has a
 * daystar-hour between minwake and maxwake. Return the time for that
 * apointment.
 * If there is no such apointment, we return the zero-time
 */
func (cs *CalendarState) GetNextWakeup(s time.Time, e time.Time, minwake int, maxwake int) (time.Time, error) {
	var z time.Time

	/* Download apointments from google */
	google_t := s.Format(time.RFC3339)
	google_e := e.Format(time.RFC3339)
	srv, err := calendar.New(cs.client)
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(google_t).TimeMax(google_e).MaxResults(100).
		OrderBy("startTime").Do()

	if err != nil {
		return z, err
	}

	/* Loop over all appointments */
	for i := range events.Items {
		/* Parse the start hour out of this apointment */
		ev := events.Items[i]
		when := ev.Start.DateTime
		go_when, err := time.Parse(time.RFC3339, when)
		ev_hour := go_when.Hour()

		if err != nil {
			return z, err
		}
		/* Does it match? Then return it. Otherwise, we'll try the next one */
		if ev_hour >= minwake && ev_hour <= maxwake {
			return go_when, nil
		}
	}

	return z, nil
}
