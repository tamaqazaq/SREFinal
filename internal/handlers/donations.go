package handlers

import (
	"fmt"
	"github.com/gemdivk/Crowdfunding-system/internal/mail"
	"github.com/gemdivk/Crowdfunding-system/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
	"log"
	"net/http"
	"os"
	"strconv"
)

func CreateDonation(c *gin.Context) {
	campaignIDStr := c.Param("id")
	campaignID, err := strconv.Atoi(campaignIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	tokenUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You must be logged in to donate"})
		return
	}
	userID := tokenUserID.(int)

	var req struct {
		Amount          float64 `json:"amount"`
		StripePaymentID string  `json:"stripe_payment_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if req.StripePaymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stripe payment ID is required"})
		return
	}

	donation := models.Donation{
		UserID:          userID,
		CampaignID:      campaignID,
		Amount:          req.Amount,
		StripePaymentID: req.StripePaymentID,
	}

	if err := models.CreateDonation(&donation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process donation"})
		return
	}

	userEmail, err := models.GetUserEmail(userID)
	if err != nil {
		log.Printf("Failed to get user email: %v", err)
	} else {
		subject := "Thanks for donation!"
		body := fmt.Sprintf("You have donated $%.2f to the company ID %d. Thanks for supporting!", req.Amount, campaignID)

		if err := mail.SendEmail(userEmail, subject, body); err != nil {
			log.Printf("Failed to send donation email: %v", err)
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":        "Donation successful",
		"donation_id":    donation.ID,
		"amount":         donation.Amount,
		"campaign_id":    donation.CampaignID,
		"stripe_payment": donation.StripePaymentID,
	})
}

func MyDonationsHandler(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	donations, err := models.GetDonationsByUserWithCampaigns(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve donations"})
		return
	}

	c.JSON(http.StatusOK, donations)
}

func CreatePaymentIntent(c *gin.Context) {
	stripe.Key = os.Getenv("STRIPE_KEY")

	var req struct {
		Amount int64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	pi, err := paymentintent.New(&stripe.PaymentIntentParams{
		Amount:   stripe.Int64(req.Amount),
		Currency: stripe.String("usd"),
	})
	if err != nil {
		log.Println("Stripe error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"client_secret": pi.ClientSecret})
}
func GetDonationsByCampaign(c *gin.Context) {

	campaignIDStr := c.Param("id")
	campaignID, err := strconv.Atoi(campaignIDStr) // Convert string to int

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}
	donations, err := models.GetDonationsForCampaign(campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch donations for campaign"})
		return
	}

	c.JSON(http.StatusOK, donations)
}

func GetDonationsByUser(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr) // Convert string to int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	donations, err := models.GetDonationsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch donations for user"})
		return
	}

	c.JSON(http.StatusOK, donations)
}

func UpdateDonation(c *gin.Context) {

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated "})
		return
	}

	donationIDStr := c.Param("id")
	donationID, err := strconv.Atoi(donationIDStr) // Convert string to int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donation ID"})
		return
	}

	donorid, err1 := models.GetUserByDonationID(donationID)
	if err1 == nil {
		if donorid != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have a permission to update it. Sorry xD"})
			return
		}
	}

	var updatedDonation models.Donation
	if err := c.ShouldBindJSON(&updatedDonation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := models.UpdateDonation(donationID, &updatedDonation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update donation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Donation updated successfully",
		"donation_id": donationID,
		"donor_id":    donorid,
		"user_id":     userID,
	})
}

func DeleteDonation(c *gin.Context) {

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	donationIDStr := c.Param("id")
	donationID, err := strconv.Atoi(donationIDStr) // Convert string to int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donation ID"})
		return
	}

	var donorID int
	donorID, err = models.GetUserByDonationID(donationID)
	if err == nil {
		if donorID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete it"})
			return
		}
	}

	if err := models.DeleteDonation(donationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete donation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Donation deleted successfully",
		"donation_id": donationID,
		"user_id":     userID,
		"donor_id":    donorID,
	})
}
