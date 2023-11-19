### TODO
- Better error handling
- Better database reading/writing/etc
- Add sig generation
- Add priv/pub key generation
- Add bech32 encoding/decoding
- Fetch other referenced events?
- Check ConnectionList stuff is right
- Sort out tags -> JSON in `writeFilter`
- Fetch other linked to nostr notes / threads?
- Fetch any linked to media?

### To run

First run
```
git submodule update --init
```

As a dev build:
```
go run .
```

As a prod build:
```
go run -tags prod .
```
