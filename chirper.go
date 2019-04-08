package main

import (
  "crypto/tls"
  "encoding/json"
  "flag"
  "fmt"
  "log"
  "net/http"

  "github.com/thoj/go-ircevent"
  "github.com/gorilla/mux"

)

var ircChan = make(chan string, 32)

type ircConfig struct {
    channel string // if doesnt start with # then direct nick
    server string // addr:port
    nick string // joe
    ssl bool
    ircObj *irc.Connection
}


var ircCfg = ircConfig{}

type apiConfig struct {
    host string
    port int // port to bind to
    prefix string // /chirper
    magickey string
}

var apiCfg = apiConfig{}


func runChirper(irccon *irc.Connection, c chan string) {
    var msg string
    for {
        select {
        case msg = <-c:
            irccon.Privmsg(ircCfg.channel, msg)
        }
    }
}


func runIrc(ircOpts ircConfig, c chan string) {
  irccon := irc.IRC(ircOpts.nick, "chirperIRC")
  irccon.VerboseCallbackHandler = true
  irccon.Debug = true
  irccon.UseTLS = ircOpts.ssl
  irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
  irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(ircOpts.channel) })
  irccon.AddCallback("366", func(e *irc.Event) {  })
  err := irccon.Connect(ircOpts.server)
  if err != nil {
    fmt.Printf("Err %s", err )
    return
  }
  ircOpts.ircObj = irccon
  go runChirper(irccon, c)
  irccon.Loop()
}


func Chirp(w http.ResponseWriter, r *http.Request) {
    params, ok := r.URL.Query()["msg"]
    if !ok {
        json.NewEncoder(w).Encode("missing msg param")
        return
    }
    if ircCfg.ircObj == nil {
        x := string(params[0])
        // ircCfg.ircObj.Privmsg(ircCfg.channel, x)
        ircChan <- x
    }
    json.NewEncoder(w).Encode("ok")
    return
}


func runHttpd(apiOpts apiConfig) {
    router := mux.NewRouter()

    router.HandleFunc("/chirp", Chirp).Methods("GET")

    bindhost := fmt.Sprintf("%s:%d", apiOpts.host, apiOpts.port)
    log.Fatal(http.ListenAndServe(bindhost, router))
}


func setupCfg() {

    flag.StringVar(&ircCfg.server, "server", "", "host:port to connect to" )
    flag.BoolVar(&ircCfg.ssl, "ssl", true, "ssl to server")
    flag.StringVar(&ircCfg.channel, "channel",
        "#chirper", "channel to connect to" )
    flag.StringVar(&ircCfg.nick, "nick",
        "chirper", "nick for chirper to use" )

    flag.StringVar(&apiCfg.host, "host",
        "localhost", "host to listen on" )
    flag.IntVar(&apiCfg.port, "port",
        8890, "port to listen on" )

    flag.Parse()
}


func main() {
  setupCfg()

  go runIrc(ircCfg, ircChan)
  runHttpd(apiCfg)
  // setupIrc()
}
