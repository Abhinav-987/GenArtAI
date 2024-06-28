package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/Abhinav-987/GenArtAI/db"
	"github.com/Abhinav-987/GenArtAI/models"
	"github.com/Abhinav-987/GenArtAI/pkg/sb"
	"github.com/Abhinav-987/GenArtAI/pkg/util"
	"github.com/Abhinav-987/GenArtAI/view/auth"

	"github.com/gorilla/sessions"
	"github.com/nedpals/supabase-go"
)

const (
	sessionUserKey        = "user"
	sessionAccessTokenKey = "accessToken"
)

func HandleResetPasswordIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.ResetPassword())
}

func HandleResetPasswordCreate(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	params := map[string]any{
		"email":      user.Email,
		"redirectTo": "http://localhost:3000/auth/reset-password",
	}
	b, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", sb.BaseAuthURL, bytes.NewReader(b))
	req.Header.Set("apikey", os.Getenv("SUPABASE_SECRET"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("supabase password recovery responded with a non 200 status code: %d => %s", resp.StatusCode, string(b))
	}
	return render(r, w, auth.ResetPasswordSuccess(user.Email))
}

func HandleResetPasswordUpdate(w http.ResponseWriter, r *http.Request) error {
	//return render(r, w, auth.ResetPassword())
	user := getAuthenticatedUser(r)
	params := map[string]any{
		"password": r.FormValue("password"),
	}
	_, err := sb.Client.Auth.UpdateUser(r.Context(), user.AccessToken, params)
	errors := auth.ResetPasswordErrors{
		NewPassword: "Please enter a valid password",
	}
	if err != nil {
		return render(r, w, auth.ResetPasswordForm(errors))
	}
	// fmt.Printf("%+v\n", resp)
	return hxRedirect(w, r, "/")
}

func HandleAccountSetupIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.AccountSetup())
}

func HandleAccountSetupCreate(w http.ResponseWriter, r *http.Request) error {
	params := auth.AccountSetupParams{
		Username: r.FormValue("username"),
	}
	errors := auth.AccountSetupErrors{}
	username := strings.TrimSpace(params.Username)
	if len(username) < 2 {
		errors.Username = "Username must be at least 2 characters long."
	} else if len(username) > 50 {
		errors.Username = "Username must be no more than 50 characters long."
	}
	if errors.Username != "" {
		return render(r, w, auth.AccountSetupForm(params, errors))
	}
	user := getAuthenticatedUser(r)
	account := models.Account{
		UserID:   user.ID,
		Username: params.Username,
	}

	if err := db.CreateAccount(&account); err != nil {
		return err
	}
	return hxRedirect(w, r, "/")
}

func HandleLoginIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.Login())
}

func HandleSignupIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.Signup())
}

func HandleSignupCreate(w http.ResponseWriter, r *http.Request) error {
	params := auth.SignupParams{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirmpassword"),
	}
	errors := auth.SignupErrors{}

	if params.Email == "" || !util.IsValidEmail(params.Email) {
		errors.Email = "Please enter a valid email id"
	}
	if params.Password == "" {
		errors.Password = "Password cannot be empty"
	} else {
		if reason, ok := util.ValidatePassword(params.Password); !ok {
			errors.Password = reason
		}
	}
	if params.Password != params.ConfirmPassword {
		errors.ConfirmPassword = "Passwords do not match"
	}
	if errors.Email != "" || errors.Password != "" || errors.ConfirmPassword != "" {
		return render(r, w, auth.SignupForm(params, errors))
	}
	user, err := sb.Client.Auth.SignUp(r.Context(), supabase.UserCredentials{
		Email:    params.Email,
		Password: params.Password,
	})
	if err != nil {
		return err
	}
	return render(r, w, auth.SignupSuccess(user.Email))
}
func HandleLoginWithGoogle(w http.ResponseWriter, r *http.Request) error {
	res, err := sb.Client.Auth.SignInWithProvider(supabase.ProviderSignInOptions{
		Provider:   "google",
		RedirectTo: "http://localhost:3000/auth/callback",
	})
	if err != nil {
		return err
	}
	http.Redirect(w, r, res.URL, http.StatusSeeOther)
	return nil

}

func HandleLoginCreate(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	resp, err := sb.Client.Auth.SignIn(r.Context(), credentials)
	if err != nil {
		slog.Error("login error", "err", err)
		return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
			InvalidCredentials: "The credentials you have entered are invalid",
		}))
	}
	if err := setAuthSession(w, r, resp.AccessToken); err != nil {
		return err
	}
	return hxRedirect(w, r, "/")
}

func HandleAuthCallback(w http.ResponseWriter, r *http.Request) error {
	accessToken := r.URL.Query().Get("access_token")
	if len(accessToken) == 0 {
		return render(r, w, auth.CallbackScript())
	}
	//fmt.Println(accessToken)
	if err := setAuthSession(w, r, accessToken); err != nil {
		return err
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func HandleLogoutCreate(w http.ResponseWriter, r *http.Request) error {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	session, _ := store.Get(r, sessionUserKey)
	session.Values[sessionAccessTokenKey] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return nil
}

func setAuthSession(w http.ResponseWriter, r *http.Request, accessToken string) error {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	session, _ := store.Get(r, sessionUserKey)
	session.Values[sessionAccessTokenKey] = accessToken
	return sessions.Save(r, w)

}
