//go:build !prod

package main

import "github.com/mmyoungman/nostr_backup/internal/json"

func UNUSED(x ...interface{}) {}

func DevBuildValidJson(str string) {
	if !json.IsValidJson(str) {
		panic("Json is not valid!")
	}
}
