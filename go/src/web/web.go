/* Implement a web service. Features
 * Manage google calendar authorization
 * Customize presets
 * Just change the lights to some setting now.
 */

package web

import "net/http"
import (
	"alarm"
	"config"
	"github.com/Joker/jade"
	"google"
	"hw"
	"io/ioutil"
	"encoding/json"
	"log"
	"path"
	"strconv"
)

/* What state do we care about? This is pretty much the entire application. */

type WebState struct {
	Config *config.Config
	Hw     hw.HWInterface
	Alarm  *alarm.Alarm
	GCal   *google.CalendarState

	/* What directory do our web resouces live in?
		 * We expect the subdirectory "views" to contain jade files
	       * and the subdirectory "resources" to contain scripts, etc.
	*/
	Webroot string
}

func StartServer(w *WebState) {

	sm := http.NewServeMux()
	sm.HandleFunc("/", w.MainPage)
	sm.HandleFunc("/views/", w.ViewHandler)
	sm.HandleFunc("/resources/", w.ResourceHandler)

	sm.HandleFunc("/api/setlights", w.APISetLights)
	sm.HandleFunc("/api/setoauth", w.APISetOauth)
	sm.HandleFunc("/api/test", w.Test)
	sm.HandleFunc("/api/getoauth", w.APIGetAuthLink)
	sm.HandleFunc("/api/config", w.APIGetConfig)
	sm.HandleFunc("/api/setconfig", w.APISetConfig)

	http.ListenAndServe(":9090", sm)

}

/* Handle resources. Just return with the raw file */
func (ws *WebState) ResourceHandler(w http.ResponseWriter, r *http.Request) {
	file := path.Base(r.URL.Path)
	data, _ := ioutil.ReadFile(path.Join(ws.Webroot, "resources", file))
	w.Write(data)

}

/* Handle views. parse the JADE file and return that */
func (ws *WebState) ViewHandler(w http.ResponseWriter, r *http.Request) {
	file := path.Base(r.URL.Path)
	data, _ := ioutil.ReadFile(path.Join(ws.Webroot, "views", file))

	output, _ := jade.Parse(file, string(data))

	w.Write([]byte(output))
}

func (ws *WebState) MainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World\n"))
}

/* API to set light intensities */
func (ws *WebState) APISetLights(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	red, _ := strconv.Atoi(q.Get("red"))
	green, _ := strconv.Atoi(q.Get("green"))
	blue, _ := strconv.Atoi(q.Get("blue"))
	ws.Alarm.UIChangeLights(red, green, blue)

}

/* API to set complete an OAUTH transaction */
func (ws *WebState) APISetOauth(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	log.Printf("DO IT: %v", q)
	code := q.Get("oauth")

	token := ws.GCal.AuthCallback(code)

	ws.Config.GoogleAuth = token
	config.SaveConfig(ws.Config)
	w.Write([]byte("OK"))
}

func (ws *WebState) Test(w http.ResponseWriter, r *http.Request) {

	ws.GCal.GetEvents()
}

/* API to get an oauth link */
func (ws *WebState) APIGetAuthLink(w http.ResponseWriter, r *http.Request) {

	u := ws.GCal.GetAuthURL()

	w.Write([]byte(u))
}

func (ws * WebState) APIGetConfig(w http.ResponseWriter, r * http.Request) {
	j, _ := json.Marshal(ws.Config);
	w.Write(j);
}

func (ws * WebState) APISetConfig(w http.ResponseWriter, r * http.Request) {
	q := r.URL.Query();
	
	ws.Config.WakeupMP3 = q.Get("WakeupMP3");
	ws.Config.WakeupMinsToMeeting, _ = strconv.Atoi(q.Get("WakeupMinsToMeeting"))
	ws.Config.MorningMeetingHourStart, _ = strconv.Atoi(q.Get("MorningMeetingHourStart"))
	ws.Config.MorningMeetingHourEnd, _ = strconv.Atoi(q.Get("MorningMeetingHourEnd"))
	config.SaveConfig(ws.Config);

}


