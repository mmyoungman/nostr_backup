package schnorr

// #cgo CFLAGS: -I../../contrib/secp256k1/src/
// #cgo CFLAGS: -DENABLE_MODULE_SCHNORRSIG=1 -DENABLE_MODULE_EXTRAKEYS=1
// #cgo CFLAGS: -DECMULT_WINDOW_SIZE=15 -DECMULT_GEN_PREC_BIT=4
// #cgo amd64 CFLAGS: -DUSE_ASM_X86_64=1
// #include "../../contrib/secp256k1/src/secp256k1.c"
// #include "../../contrib/secp256k1/src/precomputed_ecmult.c"
// #include "../../contrib/secp256k1/src/precomputed_ecmult_gen.c"
import "C"

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
)

var context *C.secp256k1_context

func init() {
	context = C.secp256k1_context_create(C.SECP256K1_CONTEXT_NONE)

	randomize := [32]byte{}
	n, err := rand.Read(randomize[:])
	if err != nil || n != len(randomize) {
		log.Fatal("Failed to generate randomness")
	}

	randomizePtr := (*C.uchar)(&randomize[0])
	return_value := C.secp256k1_context_randomize(context, randomizePtr)
	if return_value != 1 {
		log.Fatal("Failed to randomize the context")
	}
}

func IsSigValid(pubKeyStr string, sig string, serialisedEvent string) bool {
	pubKey := deserializePubKey(pubKeyStr)

	return verifyEventSig(pubKey, sig, serialisedEvent)
}

func deserializePubKey(pubKeyStr string) C.secp256k1_xonly_pubkey {
	serializedPubKey, err := hex.DecodeString(pubKeyStr)
	if err != nil {
		log.Fatal("PubKey is not valid hex. ", err)
	}
	if len(serializedPubKey) != 32 {
		log.Fatal("serializedPubKey incorrect len")
	}

	seralizedPubKeyPtr := (*C.uchar)(&serializedPubKey[0])

	var pubKey C.secp256k1_xonly_pubkey

	return_value := C.secp256k1_xonly_pubkey_parse(C.secp256k1_context_static, &pubKey, seralizedPubKeyPtr)
	if return_value != 1 {
		log.Fatal("Failed to parse serializedPubKey")
	}

	return pubKey
}

func verifyEventSig(pubKey C.secp256k1_xonly_pubkey, sigStr string, serializedEvent string) bool {
	signature, err := hex.DecodeString(sigStr)
	if err != nil {
		log.Fatal("Sig is not valid hex. ", err)
	}
	if len(signature) != 64 {
		log.Fatal("Sig incorrect len.")
	}
	signaturePtr := (*C.uchar)(&signature[0])

	eventHash := sha256.Sum256([]byte(serializedEvent))
	eventHashPtr := (*C.uchar)(&eventHash[0])

	return_value := C.secp256k1_schnorrsig_verify(context, signaturePtr, eventHashPtr, 32, &pubKey)

	if return_value == 1 {
		return true
	}
	return false
}
