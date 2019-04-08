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


see `./chirper --help` for all options and defaults
