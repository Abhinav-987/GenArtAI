package handler

import (
	"net/http"
	"strings"

	"github.com/Abhinav-987/GenArtAI/db"
	"github.com/Abhinav-987/GenArtAI/view/settings"
)

func HandleSettingIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	return render(r, w, settings.Index(user))
}

func HandleSettingUsernameUpdate(w http.ResponseWriter, r *http.Request) error {
	params := settings.ProfileParams{
		Username: r.FormValue("username"),
	}
	errors := settings.ProfileErrors{}
	username := strings.TrimSpace(params.Username)
	if len(username) < 2 {
		errors.Username = "Username must be at least 2 characters long."
	} else if len(username) > 50 {
		errors.Username = "Username must be no more than 50 characters long."
	}
	if errors.Username != "" {
		return render(r, w, settings.ProfileForm(params, errors))
	}
	user := getAuthenticatedUser(r)
	user.Account.Username = params.Username

	if err := db.UpdateAccount(&user.Account); err != nil {
		return err
	}
	params.Success = true
	return render(r, w, settings.ProfileForm(params, settings.ProfileErrors{}))
}
