package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cafeList = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	count := totalCount + 1
	city := "moscow"
	req := httptest.NewRequest("GET", "/?count="+strconv.Itoa(count)+"&city="+city, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Проверка статуса ответа
	require.Equal(t, http.StatusOK, responseRecorder.Code, "Expected status OK")

	// Проверка тела ответа
	expected := strings.Join(cafeList[city], ",")
	responseBody := responseRecorder.Body.String()

	assert.NotEmpty(t, responseBody, "Response body should not be empty")
	assert.Equal(t, expected, responseBody, "Response body should match expected value")
	assert.Len(t, strings.Split(responseBody, ","), totalCount, "Response should contain the expected number of elements")
}
