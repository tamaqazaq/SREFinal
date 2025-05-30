package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gemdivk/Crowdfunding-system/internal/models"
	"github.com/gin-gonic/gin"
)

func UpdateCampaignStatusHandler(c *gin.Context) {
	campaignID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := models.UpdateCampaignStatus(campaignID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update campaign status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign status updated successfully"})
}

func GetUsersHandler(c *gin.Context) {
	users, err := models.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func DeleteUserHandler(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := models.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func GetAdminDashboardHandler(c *gin.Context) {
	log.Println("Fetching Admin Dashboard Data...")

	// Fetch total campaigns
	totalCampaigns, err := models.GetTotalCampaigns()
	if err != nil {
		log.Println("Error fetching total campaigns:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve total campaigns"})
		return
	}

	campaignStats, err := models.GetCampaignStats()
	if err != nil {
		log.Println("Error fetching campaign stats:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve campaign stats"})
		return
	}

	topDonatedCampaigns, err := models.GetTopDonatedCampaigns()
	if err != nil {
		log.Println("Error fetching top donated campaigns:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve top donated campaigns"})
		return
	}

	totalDonations, err := models.GetTotalDonations()
	if err != nil {
		log.Println("Error fetching total donations:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve total donations"})
		return
	}

	topDonors, err := models.GetTopDonors()
	if err != nil {
		log.Println("Error fetching top donors:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve top donors"})
		return
	}

	totalUsers, err := models.GetTotalUsers()
	if err != nil {
		log.Println("Error fetching total users:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve total users"})
		return
	}

	log.Println("Dashboard Data Fetched Successfully")

	c.JSON(http.StatusOK, gin.H{
		"total_campaigns":       totalCampaigns,
		"campaign_stats":        campaignStats,
		"top_donated_campaigns": topDonatedCampaigns,
		"total_donations":       totalDonations,
		"top_donors":            topDonors,
		"total_users":           totalUsers,
	})
}
