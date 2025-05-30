package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	Mu.Lock()
	Users = make(map[int]User)
	Users[1] = User{ID: 1, Name: "Test User1"}
	NextID = 2
	Mu.Unlock()
	os.Exit(m.Run())
}

func TestRootEndpoint(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	GetRoot(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "This is a simple Go http server")
}

func TestHelloEndpoint(t *testing.T) {
	req := httptest.NewRequest("GET", "/hello", nil)
	rr := httptest.NewRecorder()
	GetHello(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Hello, HTTP!")
}

func TestCreateUser(t *testing.T) {
	newUser := User{Name: "Integration Test User"}
	jsonData, err := json.Marshal(newUser)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/user", bytes.NewBuffer(jsonData))
	rr := httptest.NewRecorder()
	UserHandler(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var createdUser User
	err = json.Unmarshal(rr.Body.Bytes(), &createdUser)
	require.NoError(t, err)

	assert.NotZero(t, createdUser.ID)
	assert.Equal(t, newUser.Name, createdUser.Name)

	Mu.RLock()
	defer Mu.RUnlock()
	storedUser, exists := Users[createdUser.ID]
	assert.True(t, exists)
	assert.Equal(t, createdUser, storedUser)
}

func TestGetUser(t *testing.T) {
	Mu.Lock()
	testUser := User{ID: 100, Name: "Test Get User"}
	Users[testUser.ID] = testUser
	Mu.Unlock()

	req := httptest.NewRequest("GET", "/user?id="+strconv.Itoa(testUser.ID), nil)
	rr := httptest.NewRecorder()
	UserHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var retrievedUser User
	err := json.Unmarshal(rr.Body.Bytes(), &retrievedUser)
	require.NoError(t, err)

	assert.Equal(t, testUser, retrievedUser)
}

func TestGetUserNotFound(t *testing.T) {
	req := httptest.NewRequest("GET", "/user?id=9999", nil)
	rr := httptest.NewRecorder()
	UserHandler(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestUpdateUser(t *testing.T) {
	Mu.Lock()
	testUser := User{ID: 200, Name: "Before Update"}
	Users[testUser.ID] = testUser
	Mu.Unlock()

	updatedUser := User{ID: 200, Name: "After Update"}
	jsonData, err := json.Marshal(updatedUser)
	require.NoError(t, err)

	req := httptest.NewRequest("PATCH", "/user", bytes.NewBuffer(jsonData))
	rr := httptest.NewRecorder()
	UserHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var responseUser User
	err = json.Unmarshal(rr.Body.Bytes(), &responseUser)
	require.NoError(t, err)

	assert.Equal(t, updatedUser, responseUser)

	Mu.RLock()
	defer Mu.RUnlock()
	storedUser, exists := Users[200]
	assert.True(t, exists)
	assert.Equal(t, updatedUser, storedUser)
}

func TestUpdateUserNotFound(t *testing.T) {
	nonExistentUser := User{ID: 9999, Name: "Non-existent"}
	jsonData, err := json.Marshal(nonExistentUser)
	require.NoError(t, err)

	req := httptest.NewRequest("PATCH", "/user", bytes.NewBuffer(jsonData))
	rr := httptest.NewRecorder()
	UserHandler(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestDeleteUser(t *testing.T) {
	Mu.Lock()
	testUser := User{ID: 300, Name: "To Be Deleted"}
	Users[testUser.ID] = testUser
	Mu.Unlock()

	req := httptest.NewRequest("DELETE", "/user?id="+strconv.Itoa(testUser.ID), nil)
	rr := httptest.NewRecorder()
	UserHandler(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	Mu.RLock()
	defer Mu.RUnlock()
	_, exists := Users[300]
	assert.False(t, exists)
}

func TestDeleteUserNotFound(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/user?id=9999", nil)
	rr := httptest.NewRecorder()
	UserHandler(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), "User not found :<")
}

func TestCreateUserWithExistingID(t *testing.T) {
	Mu.Lock()
	testUser := User{ID: 400, Name: "Existing User"}
	Users[testUser.ID] = testUser
	Mu.Unlock()

	duplicateUser := User{ID: 400, Name: "Duplicate"}
	jsonData, err := json.Marshal(duplicateUser)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/user", bytes.NewBuffer(jsonData))
	rr := httptest.NewRecorder()
	UserHandler(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)
}

func TestCreateUserInvalidData(t *testing.T) {
	invalidUser := User{Name: ""}
	jsonData, err := json.Marshal(invalidUser)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/user", bytes.NewBuffer(jsonData))
	rr := httptest.NewRecorder()
	UserHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest("PUT", "/user", nil)
	rr := httptest.NewRecorder()
	UserHandler(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestUserHandlerInvalidPath(t *testing.T) {
	req := httptest.NewRequest("GET", "/invalid", nil)
	rr := httptest.NewRecorder()
	UserHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetUserMissingID(t *testing.T) {
	req := httptest.NewRequest("GET", "/user", nil)
	rr := httptest.NewRecorder()
	UserHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Missing user ID")
}

func TestDeleteUserMissingID(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/user", nil)
	rr := httptest.NewRecorder()
	UserHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Missing user ID")
}

func TestGetUserInvalidID(t *testing.T) {
	req := httptest.NewRequest("GET", "/user?id=invalid", nil)
	rr := httptest.NewRecorder()
	UserHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid user ID")
}
