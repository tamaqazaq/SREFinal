package handlers

import (
	"log"
	"net/http"

	"github.com/gemdivk/Crowdfunding-system/internal/models"
	"github.com/gin-gonic/gin"
)

var AchievementThresholds = map[string]struct {
	Threshold int
	Name      string
	Points    int
}{
	"searcher":       {5, "Searcher", 20},
	"active_donator": {5, "Active Donator", 50},
	"daily_login":    {1, "Daily Login", 10},
	"coin_click":     {1, "Coin Click", 1},
}

func GetLeaderboard(c *gin.Context) {
	users, err := models.GetLeaderboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve leaderboard"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func UpdateUserPoints(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	achievementKey := c.Query("achievement")
	threshold, exists := AchievementThresholds[achievementKey]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid achievement"})
		return
	}

	count, err := models.IncrementUserAction(userIDInt, achievementKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track action"})
		return
	}

	log.Printf("✅ Достижение должно добавиться: User %d | Действие: %s | Количество: %d | Порог: %d", userIDInt, achievementKey, count, threshold.Threshold)

	if count == threshold.Threshold {
		err := models.AddUserAchievement(userIDInt, threshold.Name, threshold.Points)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add achievement"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Achievement unlocked!", "achievement": threshold.Name})
		return
	}

	err = models.UpdateUserPoints(userIDInt, threshold.Points)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update points"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Points updated"})
}
