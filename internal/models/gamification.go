package models

import (
	"log"

	"github.com/gemdivk/Crowdfunding-system/internal/db"
)

type LeaderboardUser struct {
	Name         string `json:"name"`
	Achievements string `json:"achievements"`
	Points       int    `json:"points"`
}

func GetLeaderboard() ([]LeaderboardUser, error) {
	var users []LeaderboardUser
	query := `
		SELECT u.name, 
		       COALESCE(string_agg(DISTINCT ua.achievement, ', '), 'No achievements yet') AS achievements, 
		       u.points 
		FROM "User" u
		LEFT JOIN "UserAchievements" ua ON u.user_id = ua.user_id
		GROUP BY u.user_id, u.name, u.points
		ORDER BY u.points DESC`
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Printf("Error fetching leaderboard: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user LeaderboardUser
		if err := rows.Scan(&user.Name, &user.Achievements, &user.Points); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func AddUserAchievement(userID int, achievementName string, points int) error {
	log.Printf("🏆 Добавляем достижение в БД: User %d | %s (%d points)", userID, achievementName, points)

	query := `
        INSERT INTO "UserAchievements" (user_id, achievement, points)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id, achievement) DO NOTHING
        RETURNING achievement`

	var addedAchievement string
	err := db.DB.QueryRow(query, userID, achievementName, points).Scan(&addedAchievement)
	if err != nil {
		log.Printf("Error adding achievement: %v", err)
		return err
	}

	log.Printf("Achievement added successfully: %s", addedAchievement)
	return UpdateUserPoints(userID, points)
}

func IncrementUserAction(userID int, action string) (int, error) {
	var count int
	query := `
		INSERT INTO "UserActions" (user_id, action, count)
		VALUES ($1, $2, 1)
		ON CONFLICT (user_id, action) DO UPDATE
		SET count = "UserActions".count + 1 RETURNING count`
	err := db.DB.QueryRow(query, userID, action).Scan(&count)
	if err != nil {
		log.Printf("Error incrementing user action: %v", err)
		return 0, err
	}
	return count, nil
}

func UpdateUserPoints(userID int, points int) error {
	query := `UPDATE "User" SET points = points + $1 WHERE user_id = $2 RETURNING points`
	var newPoints int
	err := db.DB.QueryRow(query, points, userID).Scan(&newPoints)
	if err != nil {
		log.Printf("Error updating user points: %v", err)
		return err
	}
	log.Printf("User %d now has %d points", userID, newPoints)
	return nil
}
