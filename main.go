package main

import (
	"genartai/db"
	"genartai/handler"
	"genartai/pkg/sb"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

// func foo() {
// 	ctx := context.Background()
// 	r8, err := replicate.NewClient(replicate.WithTokenFromEnv())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	model := "stability-ai/sdxl"
// 	version := "7762fd07cf82c948538e41f63f77d685e02b063e37e496e96eefd46c929f9bdc"

// 	input := replicate.PredictionInput{
// 		"prompt": "An astronaut riding a rainbow unicorn",
// 	}
// 	webhook := replicate.Webhook{
// 		URL:    "https://webhook.site/5a02082c-dd00-4484-8a96-95e7e17cc767",
// 		Events: []replicate.WebhookEventType{"start", "completed"},
// 	}
// 	output, err := r8.Run(ctx, fmt.Sprintf("%s:%s", model, version), input, &webhook)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("output :", output)
// 	prediction, _ := r8.CreatePrediction(ctx, version, input, &webhook, false)
// 	_ = r8.Wait(ctx, prediction)
// 	fmt.Println("Prediction :", prediction)
// }

func main() {
	if err := initEverything(); err != nil {
		log.Fatal(err)
	}
	//foo()
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
	slog.Info("application running", "port", port)
	log.Fatal(http.ListenAndServe(port, router))

}

func initEverything() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	if err := db.Init(); err != nil {
		return err
	}
	return sb.Init()
}
