package data_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vicebe/following-service/data"
)

func TestToJson(t *testing.T) {

	type simpleResponse struct {
		Message string `json:"message"`
	}

	sr := &simpleResponse{Message: "test"}

	var b bytes.Buffer
	if err := data.ToJson(sr, &b); err != nil {
		t.Fatal(err)
	}

	got := strings.TrimSpace(b.String())
	wanted := fmt.Sprintf("{\"message\":\"%s\"}", sr.Message)

	if got != wanted {
		t.Fatalf("wanted \"%v\" got \"%v\"", wanted, got)
	}

}
