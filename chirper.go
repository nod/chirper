package main

/*
 * SERIOUSLY UGLY HACKY MESSING AROUND
 * USER BEWARE
 * THERE ARE NO TESTS
 */

import (
  "crypto/tls"
  "encoding/json"
  "flag"
  "fmt"
  "log"
  "net/http"
  "os"
  "strings"
  "time"

  "github.com/gorilla/mux"
  "github.com/thoj/go-ircevent"
)

var ircChan = make(chan string, 32)

type ircConfig struct {
    channel string // if doesnt start with # then direct nick
    server string // addr:port
    nick string // joe
    ssl bool
    cmdPrefix string
    ircPrefix string

    stockerPrefix string
    newsPrefix string
}

type apiConfig struct {
    host string
    port int // port to bind to
    prefix string // /chirper
    magickey string
}

var ircCfg = ircConfig{}
var apiCfg = apiConfig{}


func runChirper(irccon *irc.Connection, c chan string) {
    var msg string
    for {
        select {
        case msg = <-c:
            irccon.Privmsg(ircCfg.channel, msg)
            time.Sleep(time.Second)
        }
    }
}

// thank goodness for stackoverflow
// https://stackoverflow.com/questions/17156371/how-to-get-json-response-in-golang
func getJson(url string, target interface{}) error {
    r, err := http.Get(url)
    if err != nil {
        return err
    }
    defer r.Body.Close()
    return json.NewDecoder(r.Body).Decode(target)
}

func stockTicker() {
    time.Sleep(time.Second)
}

func stocker(cmd string, ircChan chan string) {
    ircChan <- fmt.Sprintf("%s %s", ircCfg.stockerPrefix, cmd)
}

// we want a structure similar to 
type Command struct {
    handle func(irc.Event, chan string)
    helpdoc string
}
var cmdHandlers = make(map[string](Command))

func listCmds() {
    for cmdPrefix, c := range cmdHandlers {
        fmt.Println("%s:\t%s", cmdPrefix, c.helpdoc)
    }
}


func routeIRC(event *irc.Event) {
    // completely lame cmd parser
    if strings.HasPrefix(event.Message(), ircCfg.cmdPrefix) {
        cmd := strings.TrimPrefix(event.Message(), ircCfg.cmdPrefix)
        switch {
        case strings.HasPrefix(cmd, "st "):
            stocker(cmd, ircChan)
        }
    }
}

func runIrc(ircOpts ircConfig) {
    irccon := irc.IRC(ircOpts.nick, "chirperIRC")
    irccon.VerboseCallbackHandler = true
    irccon.Debug = true
    irccon.UseTLS = ircOpts.ssl
    irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
    irccon.AddCallback("001",
        func(e *irc.Event) { irccon.Join(ircOpts.channel) })
    irccon.AddCallback("366", func(e *irc.Event) {  })
    irccon.AddCallback("PRIVMSG", routeIRC)
    err := irccon.Connect(ircOpts.server)
    if err != nil {
        fmt.Printf("Err %s", err )
        return
    }
    go runChirper(irccon, ircChan)
    irccon.Loop()
}


func Chirp(w http.ResponseWriter, r *http.Request) {
    params, ok := r.URL.Query()["msg"]
    if !ok {
        json.NewEncoder(w).Encode("missing msg param")
        return
    }
    if ircChan != nil {
        x := string(params[0])
        ircChan <- x
    }
    json.NewEncoder(w).Encode("ok")
}


func runHttpd(apiOpts apiConfig) {
    router := mux.NewRouter()
    router.HandleFunc("/chirp", Chirp).Methods("GET")
    bindhost := fmt.Sprintf("%s:%d", apiOpts.host, apiOpts.port)
    log.Fatal(http.ListenAndServe(bindhost, router))
}


func setupCfg() {
    // irc related config
    flag.StringVar(&ircCfg.server, "server", "", "irc host:port to connect to")
    flag.BoolVar(&ircCfg.ssl, "ssl", true, "irc ssl to server")
    flag.StringVar(&ircCfg.channel, "channel",
        "#chirper", "irc channel to connect to" )
    flag.StringVar(&ircCfg.nick, "nick",
        "chirper", "irc nick for chirper to use" )
    flag.StringVar(&ircCfg.stockerPrefix, "stockprefix",
         "📈 ", "prefix for stocks output")
    flag.StringVar(&ircCfg.newsPrefix, "newsprefix",
         "🏢 ", "prefix for news output")

    flag.StringVar(&ircCfg.cmdPrefix, "cmdprefix",
        ".ch ", "commands to bot must start with this")

    flag.StringVar(&apiCfg.host, "host",
        "localhost", "http host to listen on" )
    flag.IntVar(&apiCfg.port, "port",
        8890, "http port to listen on" )

    showCmds := false
    flag.BoolVar(&showCmds, "list", false,
        "list available irc commands and exit")

    flag.Parse()
    if showCmds {
        listCmds()
        os.Exit(1)
    }
}


func main() {
  setupCfg()

  go runIrc(ircCfg)
  runHttpd(apiCfg)
  // setupIrc()
}
