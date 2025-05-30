package models

import (
	"database/sql"
	"fmt"
	"github.com/gemdivk/Crowdfunding-system/internal/db"
	"log"
	"time"
)

type Campaign struct {
	CampaignID   int       `json:"campaign_id" form:"campaign_id"`
	UserID       int       `json:"user_id" form:"user_id"`
	Title        string    `json:"title" form:"title"`
	Description  string    `json:"description" form:"description"`
	TargetAmount float64   `json:"target_amount" form:"target_amount"`
	AmountRaised float64   `json:"amount_raised" form:"amount_raised"`
	Status       string    `json:"status" form:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	MediaPath    string    `json:"media_path"`
	Category     string    `json:"category"`
	Email        interface{}
}

type CampaignWithCreator struct {
	Campaign Campaign `json:"campaign"`
	Creator  User     `json:"creator"`
}

var allowedCategories = map[string]bool{
	"Social Impact":                true,
	"Education & Research":         true,
	"Creative Arts":                true,
	"Technology & Innovation":      true,
	"Environment & Sustainability": true,
	"General & Miscellaneous":      true,
}

func IsValidCategory(category string) bool {
	return allowedCategories[category]
}
func CreateCampaign(campaign Campaign) error {
	// Check if the user exists.
	var userExists bool
	checkUserQuery := `SELECT EXISTS(SELECT 1 FROM "User" WHERE user_id = $1)`
	err := db.DB.QueryRow(checkUserQuery, campaign.UserID).Scan(&userExists)
	if err != nil {
		log.Printf("Error checking user existence: %v", err)
		return err
	}
	if !userExists {
		log.Printf("User with ID %d does not exist.", campaign.UserID)
		return fmt.Errorf("user with ID %d does not exist", campaign.UserID)
	}
	query := `INSERT INTO "Campaign" (user_id, title, description, target_amount, status,category, media_path) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING campaign_id, created_at, updated_at`
	err = db.DB.QueryRow(query,
		campaign.UserID,
		campaign.Title,
		campaign.Description,
		campaign.TargetAmount,
		campaign.Status, campaign.Category,
		campaign.MediaPath).
		Scan(&campaign.CampaignID, &campaign.CreatedAt, &campaign.UpdatedAt)
	if err != nil {
		log.Printf("Error inserting campaign: %v", err)
		return err
	}
	return nil
}

func GetAllCampaigns(category, search string, targetAmount, amountRaised float64) ([]Campaign, error) {
	query := `SELECT campaign_id, user_id, title, description, target_amount, amount_raised, status, media_path, category, created_at, updated_at FROM "Campaign" where 1=1`
	args := []interface{}{}
	argIndex := 1
	if category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIndex)
		args = append(args, category)
		argIndex++
	}

	if targetAmount > 0 {
		query += fmt.Sprintf(" AND target_amount <= $%d", argIndex)
		args = append(args, targetAmount)
		argIndex++
	}

	if amountRaised > 0 {
		query += fmt.Sprintf(" AND amount_raised <= $%d", argIndex)
		args = append(args, amountRaised)
		argIndex++
	}
	if search != "" {
		query += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex+1)
		args = append(args, "%"+search+"%", "%"+search+"%")
		argIndex += 2
	}
	log.Printf("Executing SQL: %s with args: %+v", query, args)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	var campaigns []Campaign

	defer rows.Close()
	for rows.Next() {
		var campaign Campaign
		var mediaPath sql.NullString
		if err := rows.Scan(&campaign.CampaignID, &campaign.UserID, &campaign.Title, &campaign.Description,
			&campaign.TargetAmount, &campaign.AmountRaised, &campaign.Status, &mediaPath, &campaign.Category, &campaign.CreatedAt, &campaign.UpdatedAt); err != nil {
			return nil, err
		}
		if mediaPath.Valid {
			campaign.MediaPath = mediaPath.String
		} else {
			campaign.MediaPath = ""
		}
		campaigns = append(campaigns, campaign)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	log.Printf("Returning campaigns: %+v", campaigns)
	return campaigns, nil
}

func GetCampaignById(campaignid int) (*CampaignWithCreator, error) {
	var result CampaignWithCreator
	query := `
		SELECT 
			c.campaign_id, 
			c.user_id, 
			c.title, 
			c.description, 
			c.target_amount, 
			c.amount_raised, 
			c.status, 
			c.created_at, 
			c.updated_at,
			c.category, 
			c.media_path, 
			u.name, 
			u.email
		FROM "Campaign" c
		JOIN "User" u ON c.user_id = u.user_id
		WHERE c.campaign_id = $1
	`
	err := db.DB.QueryRow(query, campaignid).Scan(
		&result.Campaign.CampaignID,
		&result.Campaign.UserID,
		&result.Campaign.Title,
		&result.Campaign.Description,
		&result.Campaign.TargetAmount,
		&result.Campaign.AmountRaised,
		&result.Campaign.Status,
		&result.Campaign.CreatedAt,
		&result.Campaign.UpdatedAt,
		&result.Campaign.Category,
		&result.Campaign.MediaPath,
		&result.Creator.Name,
		&result.Creator.Email,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Error retrieving campaign by ID: %v", err)
		return nil, fmt.Errorf("failed to retrieve campaign: %v", err)
	}
	return &result, nil
}

func UpdateCampaign(campaignID string, campaign Campaign) error {
	query := `
		UPDATE "Campaign" 
		SET title = $1, description = $2, target_amount = $3, status = $4, media_path = $5, category = $6,updated_at = CURRENT_TIMESTAMP
		WHERE campaign_id = $7`
	_, err := db.DB.Exec(query,
		campaign.Title,
		campaign.Description,
		campaign.TargetAmount,
		campaign.Status,
		campaign.MediaPath, campaign.Category,
		campaignID)
	if err != nil {
		log.Printf("Error updating campaign: %v", err)
		return fmt.Errorf("failed to update campaign: %v", err)
	}
	return nil
}

func DeleteCampaign(campaignID int) error {
	query := `DELETE FROM "Campaign" WHERE campaign_id = $1`
	_, err := db.DB.Exec(query, campaignID)
	if err != nil {
		log.Printf("Error deleting campaign: %v", err)
		return fmt.Errorf("failed to delete campaign: %v", err)
	}
	return nil
}

func SearchCampaigns(queryStr string) ([]Campaign, error) {
	sqlQuery := `
		SELECT campaign_id, user_id, title, description, target_amount, amount_raised, status, created_at, updated_at
		FROM "Campaign"
		WHERE title ILIKE $1 OR status ILIKE $1 OR user_id::text ILIKE $1`
	rows, err := db.DB.Query(sqlQuery, "%"+queryStr+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var campaigns []Campaign
	for rows.Next() {
		var campaign Campaign
		if err := rows.Scan(&campaign.CampaignID, &campaign.UserID, &campaign.Title, &campaign.Description,
			&campaign.TargetAmount, &campaign.AmountRaised, &campaign.Status, &campaign.CreatedAt, &campaign.UpdatedAt); err != nil {
			return nil, err
		}
		campaigns = append(campaigns, campaign)
	}
	return campaigns, nil
}

func GetCampaignByuser(userid any) ([]Campaign, error) {
	query := `SELECT campaign_id, user_id, title, description, target_amount, amount_raised, status, created_at, updated_at, media_path, category
              FROM "Campaign" WHERE user_id = $1`
	rows, err := db.DB.Query(query, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var campaigns []Campaign
	for rows.Next() {
		var campaign Campaign
		var mediaPath sql.NullString
		if err := rows.Scan(&campaign.CampaignID, &campaign.UserID, &campaign.Title, &campaign.Description,
			&campaign.TargetAmount, &campaign.AmountRaised, &campaign.Status, &campaign.CreatedAt, &campaign.UpdatedAt, &mediaPath, &campaign.Category); err != nil {
			return nil, err
		}
		campaign.MediaPath = mediaPath.String
		campaigns = append(campaigns, campaign)
	}
	return campaigns, nil
}
func GetUserEmailByID(userID int) (string, error) {
	var email string
	query := `SELECT email FROM "User" WHERE user_id = $1`
	err := db.DB.QueryRow(query, userID).Scan(&email)
	if err != nil {
		log.Printf("Error fetching user email: %v", err)
		return "", err
	}
	return email, nil
}
func GetMediaByID(campaignID string) (string, error) {
	var media string
	query := `SELECT media_path FROM "Campaign" WHERE campaign_id = $1`
	err := db.DB.QueryRow(query, campaignID).Scan(&media)
	if err != nil {
		log.Printf("Error fetching a media path of the campaign: %v", err)
		return "", err
	}
	return media, nil
}
