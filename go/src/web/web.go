

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

)

type WebState struct {
	Config *config.Config
	Hw * hw.HW
	Alarm *alarm.Alarm
	Webroot string
}

func StartServer( w *WebState) {

	sm := http.NewServeMux();
	sm.HandleFunc("/", w.MainPage);
	sm.HandleFunc("/views/", w.ViewHandler);
	sm.HandleFunc("/resources/", w.ResourceHandler);

	sm.HandleFunc("/api/setlights", w.APISetLights);

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

