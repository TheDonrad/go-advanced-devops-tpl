package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"strings"
)

func calcGaugeHash(met *Metric, key string) {
	if len(strings.TrimSpace(key)) == 0 {
		return
	}
	data := fmt.Sprintf("%s:gauge:%f", met.ID, met.Value)
	met.Hash = calcHash1(data, key)
}

func calcCounterHash(met *Metric, key string) {
	if len(strings.TrimSpace(key)) == 0 {
		return
	}
	data := fmt.Sprintf("%s:counter:%d", met.ID, met.Delta)
	met.Hash = calcHash1(data, key)
}

func calcHash1(data string, key string) string {
	src := []byte(data)
	h := hmac.New(sha256.New, []byte(key))
	h.Write(src)
	dst := h.Sum(nil)
	return fmt.Sprintf("%x", dst)
}