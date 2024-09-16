package bsdweb

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"sync"
	"time"

	T "bsvapi/types"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type wssingle struct {
	Type      string
	Conn      *websocket.Conn
	IsAlive   bool
	PingTimer *time.Timer
	Mlock     *sync.Mutex
}

type wsdb_struct struct {
	Pool map[string][]*wssingle
}

type API_obj struct {
	Type     string           `json:"type,omitempty"`
	Session  *T.BSD_Session   `json:"session,omitempty"`
	Sessions *[]T.BSD_Session `json:"sessions,omitempty"`
	Song     *T.BSD_Song      `json:"song,omitempty"`
	Songs    *[]T.BSD_Song    `json:"songs,omitempty"`
}

type ErrorStruct struct {
	Error   string
	Code    int
	Message string
}

var g_settings *T.WEB_Settings

var pb []string = []string{
	"The machine with a base-plate of prefabulated aluminite, surmounted by a malleable logarithmic casing in such a way that the two main spurving bearings were in a direct line with the pentametric fan",
	"IKEA battery supplies",
	"Probably not you...",
	"php 4.0.1",
	"The smallest brainfuck interpreter written using Piet",
	"8192 monkeys with typewriters",
	"16 dumplings and one chicken nuggie",
	"Imaginary cosmic duck",
	"13 space chickens",
	" // TODO: Fill this field in",
	"Marshmallow on a stick",
	"Two sticks and a duct tape",
	"Multipolygonal eternal PNGs",
	"40 potato batteries. Embarrassing. Barely science, really.",
	"Aperture Science computer-aided enrichment center",
	"A cluster*** of protogens",
	"Fifteen Hundred Megawatt Aperture Science Heavy Duty Super-Colliding Super Button",
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var wsdb *wsdb_struct = &wsdb_struct{
	Pool: make(map[string][]*wssingle),
}

func reader(ws_obj *wssingle) {
	for {
		_, p, err := ws_obj.Conn.ReadMessage()
		if err != nil {
			return
		}
		api_obj := API_obj{}
		json.Unmarshal(p, &api_obj)

		if api_obj.Type == "pong" {
			ws_obj.IsAlive = true
			ws_obj.PingTimer.Reset(30 * time.Second)
			continue
		}
		if api_obj.Type == "ping" {
			ws_obj.IsAlive = true
			ws_obj.PingTimer.Reset(30 * time.Second)
			ws_obj.Conn.WriteMessage(1, []byte(`{"type":"pong"}`))
			continue
		}

		// @TODO: WS API parsing

	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "404: I couldn't find my cookie")
	return
}

func broadcastApiStruct(ao API_obj) {
	ao_bytes, err := json.Marshal(ao)
	if err != nil {
		log.Printf("Broadcast API error: %s", err.Error())
		return
	}
	broadcastMsg(string(ao_bytes))
}

func broadcastMsg(msg string) {
	for _, wsarr := range wsdb.Pool {
		for _, ws := range wsarr {
			if ws != nil {
				ws.Mlock.Lock()
				ws.Conn.WriteMessage(1, []byte(msg))
				ws.Mlock.Unlock()
			}
		}
	}
}

func inarr(s1 string, arr []string) bool {
	for _, v := range arr {
		if v == s1 {
			return true
		}
	}
	return false
}

// Blackhole
func denyIncoming(w http.ResponseWriter, r *http.Request) {
	rd, e := rand.Int(rand.Reader, big.NewInt(int64(len(pb))))
	if e != nil {
		rd = big.NewInt(int64(0))
	}
	w.Header().Add("X-Powered-By", pb[rd.Int64()])
	w.Header().Add("content-type", "text/plain")
	// w.Header().Add("access-control-allow-origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With, Access-Key, API-usr, Token, ref-key, lu-key")
	w.WriteHeader(403)
	fmt.Fprintf(w, "403: Access denied")
}

func wsHandler(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	usr := r.RemoteAddr

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	ws_obj := &wssingle{Conn: ws, IsAlive: true, PingTimer: time.NewTimer(1 * time.Second), Mlock: &sync.Mutex{}}
	ws_obj.PingTimer = time.AfterFunc(30*time.Second, func() { pingShit(ws_obj) })

	pingShit(ws_obj)
	wsdb.Pool[usr] = append(wsdb.Pool[usr], ws_obj)

	ws.SetCloseHandler(func(code int, text string) error {
		wsdb.Pool[usr] = removeConn(wsdb.Pool[usr], ws)
		return nil
	})

	err = ws.WriteMessage(1, []byte(`{"type":"conncheck"}`))
	if err != nil {
		log.Println(err)
	}

	reader(ws_obj)
}

func removeConn(arr []*wssingle, dead *websocket.Conn) []*wssingle {
	idx, found := 0, false
	for s, v := range arr {
		if v.Conn == dead {
			v.PingTimer.Stop()
			v = nil
			found = true
			idx = s
			break
		}
	}
	if !found {
		return arr
	}
	return append(arr[:idx], arr[idx+1:]...)
}

func pingShit(ws_obj *wssingle) {
	if !ws_obj.IsAlive {
		ws_obj.Conn.CloseHandler()(1001, "Ping was not received")
		return
	}
	ws_obj.IsAlive = false
	ws_obj.Mlock.Lock()
	ws_obj.Conn.WriteMessage(1, []byte(`{"type":"ping"}`))
	ws_obj.Mlock.Unlock()
	ws_obj.PingTimer.Reset(30 * time.Second)
}

func preflight(w http.ResponseWriter, r *http.Request) {
	rd, e := rand.Int(rand.Reader, big.NewInt(int64(len(pb))))
	if e != nil {
		rd = big.NewInt(int64(0))
	}
	w.Header().Add("X-Powered-By", pb[rd.Int64()])
	w.Header().Add("content-type", "text/plain")
	w.Header().Add("access-control-allow-origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With, Access-Key, API-usr, Token, ref-key, lu-key")
	w.WriteHeader(204)
	fmt.Fprintf(w, "204: Access denied, but with love to the poor browser that for some reason wanted to access this page.")
}

func apiGlobalRouter(w http.ResponseWriter, r *http.Request) {
	// api thingamajig
	if r.Method == "OPTIONS" {
		preflight(w, r)
		return
	}

	vars := mux.Vars(r)
	group, ok := vars["group"]
	if !ok {
		fmt.Println("group is missing in parameters")
		homePage(w, r)
		return
	}
	endpoint, ok := vars["endpoint"]
	if !ok {
		fmt.Println("endpoint is missing in parameters")
		homePage(w, r)
		return
	}
	operation, ok := vars["operation"]
	if !ok {
		fmt.Println("operation is missing in parameters")
		homePage(w, r)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.Header().Set("access-control-allow-origin", "*")

	if r.Method != "POST" {
		homePage(w, r)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		e := ErrorStruct{Error: err.Error(), Message: "Could not read request body", Code: 500}
		j, _ := json.Marshal(e)
		fmt.Fprint(w, string(j))
		return
	}
	var apireq API_obj
	err = json.Unmarshal(body, &apireq)
	if err != nil {
		// invalid params
		e := ErrorStruct{Error: err.Error(), Message: "Could not parse JSON from the request body", Code: 500}
		j, _ := json.Marshal(e)
		fmt.Fprint(w, string(j))
		return
	}

	switch group {
	case "session":
		apiSession(w, r, endpoint, operation, &apireq)
	default:
		fmt.Println("invalid API group")
		homePage(w, r)
		return
	}
}

func apiSession(w http.ResponseWriter, r *http.Request, endpoint string, operation string, apireq *API_obj) {
	switch endpoint {
	case "latest":

	default:
		fmt.Println("invalid API endpoint")
		homePage(w, r)
	}
}

func apiRespond(w http.ResponseWriter, apiobj *API_obj) {
	j, err := json.Marshal(apiobj)
	if err != nil {
		ThrowApiErr(w, "Failed to respond with data", err, 500)
		return
	}
	fmt.Fprint(w, string(j))
}

func ThrowApiErr(w http.ResponseWriter, custom string, err error, code int) {
	if custom != "" && err == nil {
		e := ErrorStruct{Message: custom, Error: custom, Code: code}
		j, _ := json.Marshal(e)
		fmt.Fprint(w, string(j))
		return
	}

	if err != nil && custom != "" {
		e := ErrorStruct{Message: custom, Error: err.Error(), Code: 403}
		j, _ := json.Marshal(e)
		fmt.Fprint(w, string(j))
		return
	}
	custom = "You cannot perform this action. Access Denied."

	e := ErrorStruct{Message: custom, Error: custom, Code: 403}
	j, _ := json.Marshal(e)
	fmt.Fprint(w, string(j))
	return
}

func setupRoutes() {
	router := mux.NewRouter().StrictSlash(false)

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { wsHandler(w, r) })
	// router.HandleFunc("/api/{group}/{endpoint}/{operation}", apiGlobalRouter)

	router.NotFoundHandler = router.NewRoute().HandlerFunc(homePage).GetHandler()
	log.Fatal(http.ListenAndServe(":1337", router))
}

func onNewSession(c chan *T.BSD_Session) {
	for {
		session := <-c
		if session == nil {
			continue
		}
		ao := API_obj{
			Type:    "newsession",
			Session: session,
		}
		go broadcastApiStruct(ao)
	}
}
func onNewSongData(c chan *T.BSD_Song) {
	for {
		song := <-c
		if song == nil {
			continue
		}

		// strip deep tracker data
		song.DeepTrackers = nil

		ao := API_obj{
			Type: "newsong",
			Song: song,
		}
		go broadcastApiStruct(ao)
	}
}

func Init(settings *T.WEB_Settings) error {
	log.Println("Starting Web API...")

	go func() {
		setupRoutes()
		// log.Fatal(http.ListenAndServe(":1337", nil))
	}()
	g_settings = settings

	if g_settings.OnNewSongData != nil {
		go onNewSongData(*g_settings.OnNewSongData)
	}
	if g_settings.OnNewSession != nil {
		go onNewSession(*g_settings.OnNewSession)
	}

	log.Println("Web ready")
	return nil
}

func Close() error {
	// @TODO: Implement API shutdown
	return nil
}
