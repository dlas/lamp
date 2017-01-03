

package web

import "net/http"
import (
	"path"
	"io/ioutil"
	"config"
	"hw"
	"alarm"
	"github.com/Joker/jade"
	"strconv"
	"log"
	"google"

)

type WebState struct {
	Config *config.Config
	Hw hw.HWInterface
	Alarm *alarm.Alarm
	GCal *google.CalendarState
	Webroot string
}

func StartServer( w *WebState) {

	sm := http.NewServeMux();
	sm.HandleFunc("/", w.MainPage);
	sm.HandleFunc("/views/", w.ViewHandler);
	sm.HandleFunc("/resources/", w.ResourceHandler);

	sm.HandleFunc("/api/setlights", w.APISetLights);
	sm.HandleFunc("/api/setoauth", w.APISetOauth)
	sm.HandleFunc("/api/test", w.Test);
	sm.HandleFunc("/api/getoauth", w.APIGetAuthLink);

	http.ListenAndServe(":9090", sm);


}


func (ws * WebState) ResourceHandler(w http.ResponseWriter, r *http.Request) {
	file := path.Base(r.URL.Path);
	data, _ := ioutil.ReadFile(path.Join(ws.Webroot, "resources", file));
	w.Write(data)

}

func (ws * WebState) ViewHandler(w http.ResponseWriter, r *http.Request) {
	file := path.Base(r.URL.Path)
	data, _:= ioutil.ReadFile(path.Join(ws.Webroot,"views",file));

	output, _:= jade.Parse(file, string(data));

	w.Write([]byte(output))
}


func (ws * WebState) MainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World\n"))
}

func (ws * WebState) APISetLights(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	red, _ := strconv.Atoi(q.Get("red"))
	green, _ := strconv.Atoi(q.Get("green"))
	blue, _ := strconv.Atoi(q.Get("blue"))
	ws.Alarm.UIChangeLights(red, green, blue)

}

func (ws * WebState) APISetOauth(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query();
	log.Printf("DO IT: %v", q);
	code := q.Get("oauth")
	
	token := ws.GCal.AuthCallback(code);

	ws.Config.GoogleAuth = token
	config.SaveConfig(ws.Config);
	w.Write([]byte("OK"));
}


func (ws * WebState) Test(w http.ResponseWriter, r *http.Request) {

	ws.GCal.GetEvents();
}
func (ws * WebState) APIGetAuthLink(w http.ResponseWriter, r *http.Request) {

	u := ws.GCal.GetAuthURL();

	w.Write([]byte(u))
}
