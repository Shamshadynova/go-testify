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

func TestMainHandlerCorrectRequest(t *testing.T) {
	count := 3
	city := "moscow"
	req := httptest.NewRequest("GET", "/?count="+strconv.Itoa(count)+"&city="+city, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Проверка статуса ответа
	require.Equal(t, http.StatusOK, responseRecorder.Code, "Expected status OK")

	// Проверка тела ответа
	expected := strings.Join(cafeList[city][:count], ",")
	responseBody := responseRecorder.Body.String()

	assert.NotEmpty(t, responseBody, "Response body should not be empty")
	assert.Equal(t, expected, responseBody, "Response body should match expected value")
}

func TestMainHandlerCityNotSupported(t *testing.T) {
	count := 3
	city := "unknown"
	req := httptest.NewRequest("GET", "/?count="+strconv.Itoa(count)+"&city="+city, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Проверка статуса ответа
	require.Equal(t, http.StatusBadRequest, responseRecorder.Code, "Expected status BadRequest")

	// Проверка тела ответа
	expected := "wrong city value"
	responseBody := responseRecorder.Body.String()

	assert.NotEmpty(t, responseBody, "Response body should not be empty")
	assert.Equal(t, expected, responseBody, "Response body should match expected value")
}
