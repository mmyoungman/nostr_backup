### TODO
- Better error handling
- Better database reading/writing/etc
- Add sig generation
- Add priv/pub key generation
- Add bech32 encoding/decoding
- Fetch other referenced events?
- Check ConnectionList stuff is right
- Sort out tags -> JSON in `writeFilter`

### contrib

`contrib/secp256k1` currently uses v0.4.0 (or commit `199d27cea32203b224b208627533c2e813cd3b21`)
https://github.com/bitcoin-core/secp256k1/releases/tag/v0.4.0

### To run

As a dev build:
```
go run .
```

As a prod build:
```
go run -tags prod .
```
