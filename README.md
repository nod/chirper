# http to irc chirper

hacky restful interface to post things to an irc channel

## note on building

```
dep ensure
go build
```

## usage

```
./chirper --server some.irc.server:6697
```

then someplace else you can hit it like

```
curl "http://localhost:8890/chirp?msg=blech"
```

This will cause the message "blech" to show up in the default channel of
`#chirper`.


## options

see `./chirper --help` for all options and defaults

```
./chirper --help
Usage of ./chirper:
  -channel string
    	irc channel to connect to (default "#chirper")
  -host string
    	http host to listen on (default "localhost")
  -nick string
    	irc nick for chirper to use (default "chirper")
  -port int
    	http port to listen on (default 8890)
  -server string
    	irc host:port to connect to
  -ssl
    	irc ssl to server (default true)
```

## what about tests?

Yeah, I feel dirty about that.  todo.
