package db

import (
	"context"

	"github.com/Abhinav-987/GenArtAI/models"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func CreateImage(tx bun.Tx, image *models.Image) error {
	_, err := tx.NewInsert().Model(image).Exec(context.Background())
	return err
}

func UpdateImage(tx bun.Tx, image *models.Image) error {
	_, err := tx.NewUpdate().Model(image).WherePK().Exec(context.Background())
	return err
}

func GetImagesByBatchId(batchID uuid.UUID) ([]models.Image, error) {
	images := []models.Image{}
	err := Bun.NewSelect().Model(&images).Where("batch_id = ?", batchID).Scan(context.Background())
	return images, err
}

func GetImageById(id int) (models.Image, error) {
	image := models.Image{}
	err := Bun.NewSelect().Model(&image).Where("id = ?", id).Scan(context.Background())
	return image, err
}

func GetImagesByUserId(userID uuid.UUID) ([]models.Image, error) {
	images := []models.Image{}
	err := Bun.NewSelect().Model(&images).Where("deleted = ?", false).Where("user_id = ?", userID).Order("created_at desc").Scan(context.Background())
	return images, err
}

func UpdateAccount(account *models.Account) error {
	_, err := Bun.NewUpdate().Model(account).WherePK().Exec(context.Background())
	return err
}

func GetAccountByUserID(userID uuid.UUID) (models.Account, error) {
	account := models.Account{}
	err := Bun.NewSelect().Model(&account).Where("user_id = ?", userID).Scan(context.Background())
	return account, err
}

func CreateAccount(account *models.Account) error {
	_, err := Bun.NewInsert().Model(account).Exec(context.Background())
	return err
}
