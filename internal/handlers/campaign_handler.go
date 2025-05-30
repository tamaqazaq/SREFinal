package handlers

import (
	"fmt"
	"github.com/gemdivk/Crowdfunding-system/internal/mail"
	"github.com/gemdivk/Crowdfunding-system/internal/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func CreateCampaignHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please log in first"})
		return
	}

	var campaign models.Campaign
	if err := c.ShouldBindJSON(&campaign); err != nil {
		log.Printf("Binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Assign user ID
	campaign.UserID = userID.(int)

	userEmail, err := models.GetUserEmailByID(campaign.UserID)
	if err != nil {
		log.Printf("Failed to get user email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user email"})
		return
	}

	// Save campaign to DB
	if err := models.CreateCampaign(campaign); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create campaign"})
		return
	}

	emailSubject := "Campaign Created Successfully!"
	emailBody := fmt.Sprintf("<h2>Congratulations!</h2><p>Your campaign <b>%s</b> has been created successfully.</p>", campaign.Title)

	go func() {
		if err := mail.SendEmail(userEmail, emailSubject, emailBody); err != nil {
			log.Printf("Failed to send email: %v", err)
		}
	}()

	// Return JSON response
	c.JSON(http.StatusCreated, gin.H{
		"success":  true,
		"message":  "Campaign created successfully",
		"campaign": campaign,
	})
}

func GetCampaignsHandler(c *gin.Context) {
	category := c.Query("category")
	search := c.Query("search")

	targetAmountStr := c.Query("target_amount")
	amountRaisedStr := c.Query("amount_raised")
	var targetAmount, amountRaised float64
	var err error

	if targetAmountStr != "" {
		targetAmount, err = strconv.ParseFloat(targetAmountStr, 64)
		if err != nil {
			log.Println("Error parsing target_amount:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target_amount"})
			return
		}
	}

	if amountRaisedStr != "" {
		amountRaised, err = strconv.ParseFloat(amountRaisedStr, 64)
		if err != nil {
			log.Println("Error parsing amount_raised:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount_raised"})
			return
		}
	}
	campaigns, err := models.GetAllCampaigns(category, search, targetAmount, amountRaised)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve campaigns"})
		return
	}
	log.Println("Campaigns received successfully:", campaigns)
	c.JSON(http.StatusOK, campaigns)
}

func GetCampaignId(c *gin.Context) {
	campaignID := c.Param("id")
	campaignID1, err1 := strconv.Atoi(campaignID)

	if err1 != nil {
		log.Println("Invalid campaign ID:", err1)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	campaign, err := models.GetCampaignById(campaignID1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve campaign"})
		return
	}

	if campaign == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}
	c.JSON(http.StatusOK, campaign)
}

func UpdateCampaignHandler(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please log in first"})
		return
	}

	var campaign models.Campaign
	// Use ShouldBind to allow binding from multipart or JSON
	if err := c.ShouldBind(&campaign); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	campaignID := c.Param("id") // Get the campaign ID from the URL parameter
	if campaign.MediaPath == "" {
		var err1 error
		campaign.MediaPath, err1 = models.GetMediaByID(campaignID)
		if err1 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err1})
		}
	}
	err := models.UpdateCampaign(campaignID, campaign)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update campaign"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "Campaign updated successfully",
		"campaign": campaign,
	})
}

func DeleteCampaignHandler(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please log in first"})
		return
	}
	campaignID := c.Param("id")
	campaignID1, err1 := strconv.Atoi(campaignID)

	if err1 != nil {
		log.Println("Invalid campaign ID:", err1)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}
	err := models.DeleteCampaign(campaignID1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete campaign"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign deleted successfully"})
}
func UploadFileHandler(c *gin.Context) {
	file, err := c.FormFile("media")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
		return
	}

	// Save file locally
	filePath := fmt.Sprintf("./uploads/%s", file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Generate file URL (assuming you're serving files from `/uploads`)
	fileURL := fmt.Sprintf("http://localhost:8080/uploads/%s", file.Filename)

	// Return the URL in JSON response
	c.JSON(http.StatusOK, gin.H{"success": true, "file_url": fileURL})
}
func SearchCampaignsHandler(c *gin.Context) {
	// Get search query from request URL (e.g., ?query=keyword)
	query := c.DefaultQuery("query", "")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	// Call the model function to search campaigns
	campaigns, err := models.SearchCampaigns(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search campaigns"})
		return
	}

	c.JSON(http.StatusOK, campaigns)
}

func GetCampaignsbyUser(c *gin.Context) {
	userid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please log in first"})
		return
	}
	userIDParam := c.Param("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		log.Println("Invalid user ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	if userid != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to do that"})
		return
	}
	campaigns, err := models.GetCampaignByuser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve campaigns"})
		return
	}

	c.JSON(http.StatusOK, campaigns)
}
