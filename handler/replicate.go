package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Abhinav-987/GenArtAI/db"
	"github.com/Abhinav-987/GenArtAI/models"
	"github.com/Abhinav-987/GenArtAI/store"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

const (
	succeeded  = "succeeded"
	processing = "processing"
)

type ReplicateResp struct {
	Input struct {
		Prompt string `json:"prompt"`
	} `json:"input"`
	Status string   `json:"status"`
	Output []string `json:"output"`
}

func HandleReplicateCallback(w http.ResponseWriter, r *http.Request) error {
	var resp ReplicateResp
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return err
	}
	if resp.Status == processing {
		return nil
	}
	if resp.Status != succeeded {
		return fmt.Errorf("replicate callback responded with a non ok status: %s", resp.Status)
	}
	batchID, err := uuid.Parse(chi.URLParam(r, "batchID"))
	if err != nil {
		return fmt.Errorf("replicate callback invalid batchID: %s", err)
	}
	images, err := db.GetImagesByBatchId(batchID)
	if err != nil {
		return fmt.Errorf("replicate callback failed to find image with batchID %s:%s", batchID, err)
	}
	if len(images) != len(resp.Output) {
		return fmt.Errorf("replicate callback unequal images compared to replicate output")
	}
	err = db.Bun.RunInTx(r.Context(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for i, imageURL := range resp.Output {
			secureURL, err := store.ImageDownloader(imageURL)
			if err != nil {
				return fmt.Errorf("failed to download image: %s", err)
			}
			images[i].Status = models.ImageStatusCompleted
			images[i].ImageLocation = secureURL
			images[i].Prompt = resp.Input.Prompt
			if err := db.UpdateImage(tx, &images[i]); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
