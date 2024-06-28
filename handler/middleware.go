package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/Abhinav-987/GenArtAI/db"
	"github.com/Abhinav-987/GenArtAI/models"
	"github.com/Abhinav-987/GenArtAI/pkg/sb"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

func WithAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}
		user := getAuthenticatedUser(r)
		if !user.LoggedIn {
			http.Redirect(w, r, "/login?to=", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func WithAccountSetup(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := getAuthenticatedUser(r)
		account, err := db.GetAccountByUserID(user.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Redirect(w, r, "/account/setup", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
			return
		}
		user.Account = account
		ctx := context.WithValue(r.Context(), models.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}
		store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
		session, err := store.Get(r, sessionUserKey)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		accessToken, ok := session.Values[sessionAccessTokenKey].(string)
		if !ok || accessToken == "" {
			next.ServeHTTP(w, r)
			return
		}
		resp, err := sb.Client.Auth.User(r.Context(), accessToken)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user := models.AuthenticatedUser{
			ID:          uuid.MustParse(resp.ID),
			Email:       resp.Email,
			LoggedIn:    true,
			AccessToken: accessToken,
		}
		ctx := context.WithValue(r.Context(), models.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
