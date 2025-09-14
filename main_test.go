package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int
		want  int
		url   string
	}{
		{0, 0, "/cafe?count=0&city=moscow"},
		{1, 1, "/cafe?count=1&city=moscow"},
		{2, 2, "/cafe?count=2&city=moscow"},
		{100, min(len(cafeList["moscow"]), 100), "/cafe?count=100&city=moscow"},
	}

	for _, v := range requests {
		req := httptest.NewRequest("GET", v.url, nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, req)
		require.Equal(t, http.StatusOK, response.Code)
		gotCount := 0
		if len(response.Body.String()) != 0{
			splittedSlc := strings.Split(response.Body.String(), ",")
			gotCount = len(splittedSlc)
		}

		assert.Equal(t, v.want, gotCount)
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	requests := []struct {
		search    string
		wantCount int
		url       string
	}{
		{"", 0, "/cafe?city=moscow&search=фасоль"},
		{"Мир кофе,Кофе и завтраки", 2, "/cafe?city=moscow&search=кофе"},
		{"Ложка и вилка", 1, "/cafe?city=moscow&search=вилка"},
	}

	for _, v := range requests {
		req := httptest.NewRequest("GET", v.url, nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		respToLower := strings.ToLower(response.Body.String())
		searchToLower := strings.ToLower(v.search)

		respSlc := strings.Split(respToLower, ",")
		searchSlc := strings.Split(searchToLower, ",")

		assert.Equal(t,len(respSlc), len(searchSlc)) 
			for _, v := range searchSlc {
				assert.Contains(t, respSlc, v)
			}
		}
	}

