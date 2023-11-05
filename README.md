### TODO
- Better error handling
- Better database reading/writing/etc
- Add sig generation
- Add priv/pub key generation
- Add bech32 encoding/decoding
- Fetch other referenced events?
- Check ConnectionList stuff is right
- Sort out tags -> JSON in `writeFilter`

### To run

First
```
cd path/to/nostr_backup
git submodule init
git submodule update
```

As a dev build:
```
go run .
```

As a prod build:
```
go run -tags prod .
```
