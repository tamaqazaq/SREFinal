package models

import (
	"database/sql"
	"fmt"
	"github.com/gemdivk/Crowdfunding-system/internal/db"
	"log"
	"time"
)

type Donation struct {
	ID              int       `json:"id"`
	CampaignID      int       `json:"campaign_id"` // Changed to int
	UserID          int       `json:"user_id"`     // Changed to int
	Amount          float64   `json:"amount"`
	DonationDate    time.Time `json:"donation_date"`
	StripePaymentID string    `json:"stripe_payment_id"`
}
type DonationWithCampaign struct {
	Donation Donation `json:"donation"`
	Campaign Campaign `json:"campaign"`
}
type DonationWithUser struct {
	Donation Donation `json:"donation"`
	User     User     `json:"user"`
}

func CreateDonation(donation *Donation) error {

	query := `INSERT INTO "Donation" (user_id, campaign_id, amount, donation_date, stripe_payment_id) 
          VALUES ($1, $2, $3, CURRENT_TIMESTAMP, $4) RETURNING donation_id`
	err := db.DB.QueryRow(query, donation.UserID, donation.CampaignID, donation.Amount, donation.StripePaymentID).
		Scan(&donation.ID)
	if err != nil {
		return fmt.Errorf("Failed to insert donation: %v", err)
	}

	err = UpdateAmountRaised(donation.CampaignID, donation.Amount)
	if err != nil {
		return err
	}

	return nil
}

func GetDonationsForCampaign(campaignID int) ([]DonationWithUser, error) {
	var donations []DonationWithUser
	query := `SELECT d.donation_id, d.user_id, d.campaign_id, d.amount, d.donation_date, u.name, u.email 
              FROM "Donation" d join "User" u on d.user_id = u.user_id  WHERE d.campaign_id = $1 `
	rows, err := db.DB.Query(query, campaignID)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch donations for campaign: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var dw DonationWithUser
		if err := rows.Scan(&dw.Donation.ID, &dw.Donation.UserID, &dw.Donation.CampaignID, &dw.Donation.Amount, &dw.Donation.DonationDate, &dw.User.Name, &dw.User.Email); err != nil {
			return nil, fmt.Errorf("Failed to scan donation: %v", err)
		}
		donations = append(donations, dw)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return donations, nil
}

func GetUserByDonationID(donationID int) (int, error) {
	var userID int
	query := `SELECT user_id from "Donation" where donation_id = $1`
	err := db.DB.QueryRow(query, donationID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no donation found for donation_id: %d", donationID)
		}
		return 0, err
	}
	return userID, nil
}

func GetDonationsByUser(userID int) ([]Donation, error) {
	var donations []Donation
	query := `SELECT donation_id, user_id, campaign_id, amount, donation_date 
              FROM "Donation" WHERE user_id = $1`
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch donations for user: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var donation Donation
		if err := rows.Scan(&donation.ID, &donation.UserID, &donation.CampaignID, &donation.Amount, &donation.DonationDate); err != nil {
			return nil, fmt.Errorf("Failed to scan donation: %v", err)
		}
		donations = append(donations, donation)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return donations, nil
}

func UpdateDonation(id int, updatedDonation *Donation) error { // id is now int
	query := `UPDATE "Donation" SET amount = $1, donation_date = CURRENT_TIMESTAMP 
          WHERE donation_id = $2`
	_, err := db.DB.Exec(query, updatedDonation.Amount, id)
	if err != nil {
		return fmt.Errorf("Failed to update donation: %v", err)
	}
	return nil
}

func DeleteDonation(id int) error {

	query := `DELETE FROM "Donation" WHERE donation_id = $1`
	_, err := db.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Failed to delete donation: %v", err)
	}
	return nil
}
func UpdateAmountRaised(campaignID int, donationAmount float64) error {
	var currentAmountRaised float64
	err := db.DB.QueryRow(`SELECT amount_raised FROM "Campaign" WHERE campaign_id = $1`, campaignID).Scan(&currentAmountRaised)
	if err != nil {
		return err
	}
	newAmountRaised := currentAmountRaised + donationAmount
	_, err = db.DB.Exec(`UPDATE "Campaign" SET amount_raised = $1 WHERE campaign_id = $2`,
		newAmountRaised, campaignID)
	return err
}

func GetDonationsByUserWithCampaigns(userID int) ([]DonationWithCampaign, error) {
	query := `
		SELECT 
			d.donation_id, d.campaign_id, d.amount, d.donation_date,
			c.campaign_id, c.title, c.description, c.media_path, c.target_amount, c.amount_raised
		FROM "Donation" d
		JOIN "Campaign" c ON d.campaign_id = c.campaign_id
		WHERE d.user_id = $1
	`
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch donations with campaigns: %v", err)
	}
	defer rows.Close()

	var results []DonationWithCampaign
	for rows.Next() {
		var dwc DonationWithCampaign
		err = rows.Scan(
			&dwc.Donation.ID, &dwc.Donation.CampaignID, &dwc.Donation.Amount, &dwc.Donation.DonationDate,
			&dwc.Campaign.CampaignID, &dwc.Campaign.Title, &dwc.Campaign.Description, &dwc.Campaign.MediaPath,
			&dwc.Campaign.TargetAmount, &dwc.Campaign.AmountRaised,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan donation with campaign: %v", err)
		}
		results = append(results, dwc)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
func GetUserEmail(userID int) (string, error) {
	var email string
	query := `SELECT email FROM "User" WHERE user_id = $1`
	err := db.DB.QueryRow(query, userID).Scan(&email)
	if err != nil {
		log.Printf("Error fetching user email: %v", err)
		return "", err
	}
	return email, nil
}
