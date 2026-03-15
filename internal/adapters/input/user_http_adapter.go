package input

import (
	"fmt"
	"net/http"

	"github.com/charlires/go-base-backend-service/internal/core/services"
)

type UserHTTPAdapter struct {
	UserService services.UserService
}

func NewUserHTTPAdapter(userService services.UserService) *UserHTTPAdapter {
	return &UserHTTPAdapter{
		UserService: userService,
	}
}

func (a *UserHTTPAdapter) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userid")
	user, err := a.UserService.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "User details: %+v", user)
}
