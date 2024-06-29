package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Abhinav-987/GenArtAI/db"
	"github.com/Abhinav-987/GenArtAI/models"
	"github.com/Abhinav-987/GenArtAI/view/generate"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/replicate/replicate-go"
	"github.com/uptrace/bun"
)

const creditsPerImage = 2

func HandleGenerateImageStatus(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return err
	}
	image, err := db.GetImageById(id)
	if err != nil {
		return err
	}
	slog.Info("checking image status", "id", id)
	return render(r, w, generate.GalleryImage(image))
}

func HandleGenerateCreate(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	amount, _ := strconv.Atoi(r.FormValue("amount"))
	// prompt := "red sportscar in the garden"
	params := generate.FormParams{
		Prompt: r.FormValue("prompt"),
		Amount: amount,
	}
	errors := generate.FormErrors{}
	if amount <= 0 || amount > 8 {
		errors.Amount = "Please enter a valid amount"
		return render(r, w, generate.Form(params, errors))
	}
	if len(params.Prompt) < 10 {
		errors.Prompt = "Prompt must be at least 10 characters long."
	} else if len(params.Prompt) > 100 {
		errors.Prompt = "Prompt must be no more than 100 characters long."
	}
	if errors.Prompt != "" {
		return render(r, w, generate.Form(params, errors))
	}
	creditsNeeded := params.Amount * creditsPerImage
	if user.Account.Credits < creditsNeeded {
		errors.CreditsNeeded = creditsNeeded
		errors.UserCredits = user.Account.Credits
		errors.Credits = true
		return render(r, w, generate.Form(params, errors))
	}
	user.Account.Credits -= creditsNeeded

	if err := db.UpdateCreditsByUserID(user.UserID, user.Account.Credits); err != nil {
		return err
	}

	batchID := uuid.New()
	genParams := GenerateImageParams{
		Prompt:  params.Prompt,
		Amount:  params.Amount,
		UserID:  user.ID,
		BatchID: batchID,
	}
	if err := generateImages(r.Context(), genParams); err != nil {
		return err
	}
	err := db.Bun.RunInTx(r.Context(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for i := 0; i < params.Amount; i++ {
			img := models.Image{
				UserID:  user.ID,
				Status:  models.ImageStatusPending,
				BatchID: batchID,
			}
			if err := db.CreateImage(tx, &img); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return hxRedirect(w, r, "/generate")
}

func HandleGenerateIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	images, err := db.GetImagesByUserId(user.ID)
	if err != nil {
		return err
	}
	data := generate.ViewData{
		Images: images,
	}
	return render(r, w, generate.Index(data))
}

type GenerateImageParams struct {
	Prompt  string
	Amount  int
	BatchID uuid.UUID
	UserID  uuid.UUID
}

func generateImages(ctx context.Context, params GenerateImageParams) error {
	r8, err := replicate.NewClient(replicate.WithTokenFromEnv())
	if err != nil {
		log.Fatal(err)
	}
	//model := "stability-ai/sdxl"
	version := "7762fd07cf82c948538e41f63f77d685e02b063e37e496e96eefd46c929f9bdc"
	input := replicate.PredictionInput{
		"prompt":      params.Prompt,
		"num_outputs": params.Amount,
	}
	webhook := replicate.Webhook{
		URL:    fmt.Sprintf("https://webhook.site/5a02082c-dd00-4484-8a96-95e7e17cc767/%s/%s", params.UserID, params.BatchID),
		Events: []replicate.WebhookEventType{"start", "completed"},
	}
	_, err = r8.CreatePrediction(ctx, version, input, &webhook, false)
	return err
}
