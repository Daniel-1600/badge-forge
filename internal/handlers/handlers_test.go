package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"tahrir-go/internal/models"
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
