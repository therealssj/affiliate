package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/spaco/affiliate/src/config"
	"net/http/httptest"
	"testing"
	"time"
)

//func setAuthHeaders(req *http.Request, teller *config.Teller) {
//	//	timestamp := strconv.Itoa(time.Now().Unix())
//	timestamp := fmt.Sprintf("%d", time.Now().Unix())
//	//	hash := md5.New()
//	//	io.WriteString(hash, timestamp+teller.ApiToken)
//	hash := hmac.New(sha256.New, []byte(teller.ApiToken))
//	hash.Write([]byte(timestamp))
//
//	req.Header.Set("timestamp", timestamp)
//	//	req.Header.Set("auth", fmt.Sprintf("%x", hash.Sum(nil)))
//	req.Header.Set("auth", hex.EncodeToString(hash.Sum(nil)))
//	req.Header.Set("affiliate", "true")
//}

func TestChechAuthHeader(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	conf := config.GetApiForTellerConfig()
	hash := hmac.New(sha256.New, []byte(conf.AuthToken))
	hash.Write([]byte(timestamp))
	r.Header.Set("timestamp", timestamp)
	r.Header.Set("auth", hex.EncodeToString(hash.Sum(nil)))

	if !chechAuthHeader(w, r) {
		t.Errorf("Failed. check error")
	}

	timestamp = fmt.Sprintf("%d", time.Now().Unix()-15)
	conf = config.GetApiForTellerConfig()
	hash = hmac.New(sha256.New, []byte(conf.AuthToken))
	hash.Write([]byte(timestamp))
	r.Header.Set("timestamp", timestamp)
	r.Header.Set("auth", hex.EncodeToString(hash.Sum(nil)))
	if !chechAuthHeader(w, r) {
		t.Errorf("Failed. check error")
	}

	timestamp = fmt.Sprintf("%d", time.Now().Unix()-16)
	conf = config.GetApiForTellerConfig()
	hash = hmac.New(sha256.New, []byte(conf.AuthToken))
	hash.Write([]byte(timestamp))
	r.Header.Set("timestamp", timestamp)
	r.Header.Set("auth", hex.EncodeToString(hash.Sum(nil)))
	if chechAuthHeader(w, r) {
		t.Errorf("Failed. check error")
	}

	timestamp = fmt.Sprintf("%d", time.Now().Unix()+3)
	conf = config.GetApiForTellerConfig()
	hash = hmac.New(sha256.New, []byte(conf.AuthToken))
	hash.Write([]byte(timestamp))
	r.Header.Set("timestamp", timestamp)
	r.Header.Set("auth", hex.EncodeToString(hash.Sum(nil)))
	if !chechAuthHeader(w, r) {
		t.Errorf("Failed. check error")
	}

	timestamp = fmt.Sprintf("%d", time.Now().Unix()+4)
	conf = config.GetApiForTellerConfig()
	hash = hmac.New(sha256.New, []byte(conf.AuthToken))
	hash.Write([]byte(timestamp))
	r.Header.Set("timestamp", timestamp)
	r.Header.Set("auth", hex.EncodeToString(hash.Sum(nil)))
	if chechAuthHeader(w, r) {
		t.Errorf("Failed. check error")
	}

	hash = hmac.New(sha256.New, []byte("test-not-right"))
	hash.Write([]byte(timestamp))
	r.Header.Set("timestamp", timestamp)
	r.Header.Set("auth", hex.EncodeToString(hash.Sum(nil)))

	if chechAuthHeader(w, r) {
		t.Errorf("Failed. check error")
	}
}
