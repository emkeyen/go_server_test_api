package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

var (
	Users  = make(map[int]User)
	Mu     sync.RWMutex
	NextID = 1
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("This is a simple Go http server :)\n"))
}

func GetHello(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Hello, HTTP!\n"))
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetUser(w, r)
	case http.MethodPost:
		handleCreateUser(w, r)
	case http.MethodPatch:
		handleUpdateUser(w, r)
	case http.MethodDelete:
		handleDeleteUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	Mu.RLock()
	user, exists := Users[id]
	Mu.RUnlock()

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	if user.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	Mu.Lock()
	defer Mu.Unlock()

	if user.ID == 0 {
		user.ID = NextID
		NextID++
	}

	if _, exists := Users[user.ID]; exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	Users[user.ID] = user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	Mu.Lock()
	defer Mu.Unlock()

	_, exists := Users[user.ID]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	Users[user.ID] = user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	Mu.Lock()
	defer Mu.Unlock()

	if _, exists := Users[id]; !exists {
		http.Error(w, "User not found :<", http.StatusNotFound)
		return
	}

	delete(Users, id)
	w.WriteHeader(http.StatusNoContent)
}
