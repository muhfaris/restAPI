package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/globalsign/mgo"
	"gitlab.com/muhfaris/restAPI/router"

	"gitlab.com/muhfaris/restAPI/internal/pkg/logging"
)

var (
	logger *logging.Logger
)

func TestDatabase(t *testing.T) {
	host := "127.0.0.1"
	port := 27017
	dbName := "blog"
	uri := fmt.Sprintf("%s:%d", host, port)

	session, err := mgo.Dial(uri)
	if err != nil {
		t.Errorf("Database is offline: %s", err)
	}

	dbPool := session.DB(dbName)

	// close MongoDB connections when we're finished
	//defer session.Close()

	router.Init(logger, dbPool)
}

func TestHandlerArticleList(t *testing.T) {
	//test database
	TestDatabase(t)

	req, err := http.NewRequest("GET", "/articles", nil)
	if err != nil {
		t.Errorf("The URL is not found 404 : /articles: %v", err)
	}

	// response recorder
	rr := httptest.NewRecorder()
	h := router.HandlerFunc(router.HandlerArticleList)

	h.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler structure wrong status code : got %v want %v", status, http.StatusOK)
	}

	if http.StatusOK != rr.Code {
		t.Errorf("Expected: %v, But got : %v", rr.Code, http.StatusOK)
	}
}

func TestHandlerArticleDetail(t *testing.T) {
	//test database
	TestDatabase(t)

	req, err := http.NewRequest("GET", "/articles/5c4148d50897ac0b072b18a3", nil)
	if err != nil {
		t.Errorf("The URL is not found 404 : /articles: %v", err)
	}

	// response recorder
	rr := httptest.NewRecorder()
	h := router.HandlerFunc(router.HandlerArticleDetail)

	h.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler structure wrong status code : got %v want %v", status, http.StatusOK)
	}

	if http.StatusOK != rr.Code {
		t.Errorf("Expected: %v, But got : %v", rr.Code, http.StatusOK)
	}
}

func TestHandlerArticleCreate(t *testing.T) {
	//test database
	TestDatabase(t)

	req, err := http.NewRequest("POST", "/articles", nil)
	if err != nil {
		t.Errorf("The URL is not found 404 : /articles: %v", err)
	}

	// response recorder
	rr := httptest.NewRecorder()
	h := router.HandlerFunc(router.HandlerArticleCreate)

	h.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler structure wrong status code : got %v want %v", status, http.StatusOK)
	}

	if http.StatusOK != rr.Code {
		t.Errorf("Expected: %v, But got : %v", rr.Code, http.StatusOK)
	}
}
