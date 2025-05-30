package models

import (
	"fmt"
	"log"

	"github.com/gemdivk/Crowdfunding-system/internal/db"
)

func GetAllUsers() ([]User, error) {
	query := `SELECT user_id, name, email, role, is_verified, created_at FROM "User"`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.UserID, &user.Name, &user.Email, &user.Role, &user.IsVerified, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func DeleteUser(userID int) error {
	query := `DELETE FROM "User" WHERE user_id = $1`
	_, err := db.DB.Exec(query, userID)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}

func UpdateCampaignStatus(campaignID int, status string) error {
	query := `UPDATE "Campaign" SET status = $1 WHERE campaign_id = $2`
	_, err := db.DB.Exec(query, status, campaignID)
	if err != nil {
		log.Printf("Error updating campaign status: %v", err)
		return fmt.Errorf("failed to update campaign status: %v", err)
	}
	return nil
}
func GetTotalCampaigns() (int, error) {
	query := `SELECT COUNT(*) FROM "Campaign"`
	var total int
	err := db.DB.QueryRow(query).Scan(&total)
	if err != nil {
		log.Println("Error fetching total campaigns:", err)
		return 0, err
	}
	return total, nil
}

func GetCampaignStats() (map[string]int, error) {
	query := `
		SELECT 
			COUNT(CASE WHEN status = 'active' THEN 1 END) AS active_campaigns,
			COUNT(CASE WHEN status = 'inactive' THEN 1 END) AS inactive_campaigns,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) AS completed_campaigns
		FROM "Campaign"`
	var active, inactive, completed int
	err := db.DB.QueryRow(query).Scan(&active, &inactive, &completed)
	if err != nil {
		log.Printf("Error getting campaign stats: %v", err)
		return nil, fmt.Errorf("failed to get campaign stats: %v", err)
	}

	stats := map[string]int{
		"active_campaigns":    active,
		"inactive_campaigns":  inactive,
		"completed_campaigns": completed,
	}
	return stats, nil
}

func GetTopDonatedCampaigns() ([]map[string]interface{}, error) {
	query := `
		SELECT c.title, COALESCE(SUM(d.amount), 0) AS total_donated
		FROM "Donation" d
		JOIN "Campaign" c ON d.campaign_id = c.campaign_id
		GROUP BY c.campaign_id, c.title
		ORDER BY total_donated DESC
		LIMIT 5`
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Printf("Error getting top donated campaigns: %v", err)
		return nil, fmt.Errorf("failed to get top donated campaigns: %v", err)
	}
	defer rows.Close()

	var campaigns []map[string]interface{}
	for rows.Next() {
		var title string
		var totalDonated float64
		if err := rows.Scan(&title, &totalDonated); err != nil {
			log.Printf("Error scanning top donated campaigns: %v", err)
			return nil, fmt.Errorf("failed to scan top donated campaigns: %v", err)
		}
		campaigns = append(campaigns, map[string]interface{}{
			"title":         title,
			"total_donated": totalDonated,
		})
	}
	return campaigns, nil
}

func GetTotalDonations() (float64, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM "Donation"`
	var totalDonated float64
	err := db.DB.QueryRow(query).Scan(&totalDonated)
	if err != nil {
		log.Printf("Error getting total donations: %v", err)
		return 0, fmt.Errorf("failed to get total donations: %v", err)
	}
	return totalDonated, nil
}

func GetTopDonors() ([]map[string]interface{}, error) {
	query := `
		SELECT u.name, COALESCE(SUM(d.amount), 0) AS total_donated
		FROM "Donation" d
		JOIN "User" u ON d.user_id = u.user_id
		GROUP BY u.name
		ORDER BY total_donated DESC
		LIMIT 5`
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Printf("Error getting top donors: %v", err)
		return nil, fmt.Errorf("failed to get top donors: %v", err)
	}
	defer rows.Close()

	var donors []map[string]interface{}
	for rows.Next() {
		var name string
		var totalDonated float64
		if err := rows.Scan(&name, &totalDonated); err != nil {
			log.Printf("Error scanning top donors: %v", err)
			return nil, fmt.Errorf("failed to scan top donors: %v", err)
		}
		donors = append(donors, map[string]interface{}{
			"name":          name,
			"total_donated": totalDonated,
		})
	}
	return donors, nil
}

func GetTotalUsers() (int, error) {
	query := `SELECT COUNT(*) FROM "User"`
	var total int
	err := db.DB.QueryRow(query).Scan(&total)
	if err != nil {
		log.Printf("Error getting total users: %v", err)
		return 0, fmt.Errorf("failed to get total users: %v", err)
	}
	return total, nil
}
