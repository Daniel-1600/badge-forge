package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"tahrir-go/internal/models"
	"tahrir-go/internal/rules"
)

func TestCreatePersonHandler(t *testing.T) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(&models.Person{}))

	handler := CreatePersonHandler(database)
	post := func(body string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		handler(rec, httptest.NewRequest(http.MethodPost, "/persons", strings.NewReader(body)))
		return rec
	}

	first := post(`{"Nickname":"ferris"}`)
	require.Equal(t, http.StatusCreated, first.Code)
	var created models.Person
	require.NoError(t, json.NewDecoder(first.Body).Decode(&created))
	require.NotZero(t, created.ID)

	second := post(`{"Nickname":"ferris"}`)
	require.Equal(t, http.StatusOK, second.Code)
	var existing models.Person
	require.NoError(t, json.NewDecoder(second.Body).Decode(&existing))
	require.Equal(t, created.ID, existing.ID)

	require.Equal(t, http.StatusBadRequest, post(`{}`).Code)
	require.Equal(t, http.StatusBadRequest, post(`not json`).Code)
}

func TestPostAssertionDefaultsIssuedOn(t *testing.T) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(&models.Assertion{}))

	handler := PostAssertion(database, make(chan rules.Event, 2))
	post := func(body string) models.Assertion {
		rec := httptest.NewRecorder()
		handler(rec, httptest.NewRequest(http.MethodPost, "/assertions", strings.NewReader(body)))
		require.Equal(t, http.StatusCreated, rec.Code)
		var created models.Assertion
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&created))
		return created
	}

	before := time.Now()
	defaulted := post(`{"ID":"a1","BadgeID":"b","PersonID":1}`)
	require.False(t, defaulted.IssuedOn.Before(before))

	explicit := post(`{"ID":"a2","BadgeID":"b","PersonID":1,"IssuedOn":"2024-01-02T03:04:05Z"}`)
	require.True(t, explicit.IssuedOn.Equal(time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)))
}
