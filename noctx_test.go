// +build -go1.7

package handlers

import (
	"net/http"
	"testing"

	uuid "github.com/satori/go.uuid"
)

func TestGetRequestID(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	u := uuid.NewV4()
	_, ok := GetRequestID(req)
	if ok != false {
		t.Error("expected request id get to return false, got true")
	}
	req = SetRequestID(req, u)
	uid, ok := GetRequestID(req)
	if !ok {
		t.Error("expected request id get to return true, got false")
	}
	if uid.String() != u.String() {
		t.Errorf("expected %s (from context) to equal %s", uid.String(), u.String())
	}
}
