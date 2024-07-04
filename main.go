package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/Abhinav-987/GenArtAI/db"
	"github.com/Abhinav-987/GenArtAI/handler"
	cld "github.com/Abhinav-987/GenArtAI/pkg/cloudinary"
	"github.com/Abhinav-987/GenArtAI/pkg/sb"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := initEverything(); err != nil {
		log.Fatal(err)
	}
	router := chi.NewMux()
	router.Use(handler.WithUser)

	fileServer := http.FileServer(http.Dir("./public"))
	router.Handle("/public/*", http.StripPrefix("/public", fileServer))
	router.Get("/", handler.MakeHandler(handler.HandleHomeIndex))
	router.Get("/login", handler.MakeHandler(handler.HandleLoginIndex))
	router.Get("/login/provider/google", handler.MakeHandler(handler.HandleLoginWithGoogle))
	router.Get("/signup", handler.MakeHandler(handler.HandleSignupIndex))
	router.Post("/login", handler.MakeHandler(handler.HandleLoginCreate))
	router.Post("/logout", handler.MakeHandler(handler.HandleLogoutCreate))
	router.Post("/signup", handler.MakeHandler(handler.HandleSignupCreate))
	router.Get("/auth/callback", handler.MakeHandler(handler.HandleAuthCallback))
	router.Post("/replicate/callback/{userID}/{batchID}", handler.MakeHandler(handler.HandleReplicateCallback))

	router.Group(func(r chi.Router) {
		r.Use(handler.WithAuth)
		r.Get("/account/setup", handler.MakeHandler(handler.HandleAccountSetupIndex))
		r.Post("/account/setup", handler.MakeHandler(handler.HandleAccountSetupCreate))
	})

	router.Group(func(r chi.Router) {
		r.Use(handler.WithAuth, handler.WithAccountSetup)
		r.Get("/settings", handler.MakeHandler(handler.HandleSettingIndex))
		r.Put("/settings/account/profile", handler.MakeHandler(handler.HandleSettingUsernameUpdate))
		r.Get("/auth/reset-password", handler.MakeHandler(handler.HandleResetPasswordIndex))
		r.Post("/auth/reset-password", handler.MakeHandler(handler.HandleResetPasswordCreate))
		r.Put("/auth/reset-password", handler.MakeHandler(handler.HandleResetPasswordUpdate))

		r.Get("/generate", handler.MakeHandler(handler.HandleGenerateIndex))
		r.Post("/generate", handler.MakeHandler(handler.HandleGenerateCreate))
		r.Get("/buy-credits", handler.MakeHandler(handler.HandleCreditsIndex))
		r.Get("/generate/image/status/{id}", handler.MakeHandler(handler.HandleGenerateImageStatus))
	})

	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	slog.Info("application running", "port", port)
	log.Fatal(http.ListenAndServe(addr, router))

}

func initEverything() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	if err := db.Init(); err != nil {
		return err
	}
	if err := cld.Init(); err != nil {
		return err
	}
	return sb.Init()
}
