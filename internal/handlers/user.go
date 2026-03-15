package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"Assignment5/internal/usecase"
	"Assignment5/pkg/modules"
)

type UserHandler struct {
	usecase *usecase.UserUsecase
}

func NewUserHandler(u *usecase.UserUsecase) *UserHandler {
	return &UserHandler{u}
}

func (h *UserHandler) Health(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	filter := modules.UserFilter{
		Name:      q.Get("name"),
		Email:     q.Get("email"),
		Gender:    q.Get("gender"),
		BirthDate: q.Get("birth_date"),
		OrderBy:   q.Get("order_by"),
		Limit:     limit,
		Offset:    offset,
	}

	if idStr := q.Get("id"); idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err == nil {
			filter.ID = &id
		}
	}

	res, err := h.usecase.GetUsers(filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *UserHandler) CommonFriends(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()

	u1, _ := strconv.Atoi(q.Get("user1"))
	u2, _ := strconv.Atoi(q.Get("user2"))

	users, _ := h.usecase.GetCommonFriends(u1, u2)

	json.NewEncoder(w).Encode(users)
}
