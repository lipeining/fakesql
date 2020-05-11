package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
func TestPostTable(t *testing.T) {
	setup()
	router := setupRouter()

	// bodyStr := `{ "tblName": "card", "num": 10000, jsonPath: "./tables/card.json" }`
	// requestData := bytes.NewBuffer([]byte(bodyStr))
	
	w := httptest.NewRecorder()
	param := url.Values{}
	param.Set("tblName", "card")
	param.Add("num", "10000")
	// cols := "[{\"name\":\"id\",\"t\":\"id\",\"def\":\"\"},{\"name\":\"name\",\"t\":\"string\",\"def\":\"\"},{\"name\":\"card\",\"t\":\"string\",\"def\":\"\"},{\"name\":\"create_time\",\"t\":\"datetime\",\"def\":\"\"}]"
	// param.Add("cols", cols)
	cols, _ := json.Marshal([]map[string]interface{}{
		{ "name": "id", "t": "id", "def": ""},
		{ "name": "name", "t": "string", "def": ""},
		{ "name": "card", "t": "string", "def": ""},
		{ "name": "create_time", "t": "datetime", "def": ""},
	})
	param.Add("cols", string(cols))
	// param.Add("jsonPath", "\\.\\/tables\\/card\\.json")
	requestData := strings.NewReader(param.Encode())
	req, _ := http.NewRequest("POST", "/tables", requestData)
	// req.Header.Set("Content-type", "application/json")
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "ok", w.Body.String())
}