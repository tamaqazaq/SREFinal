package handlers

import (
	"fmt"
	"github.com/gemdivk/Crowdfunding-system/internal/mail"
	"github.com/gemdivk/Crowdfunding-system/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"time"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Email  string `json:"email"`
	UserID int
	jwt.RegisteredClaims
}

func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := models.Register(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := generateVerificationToken(user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verification token"})
		return
	}

	verificationLink := fmt.Sprintf("https://srefinal.onrender.com/verify-email?token=%s", token)

	subject := "Welcome to Qadam! Please verify your email"
	body := fmt.Sprintf(
		"Hello, %s! Your account has been created. Please click the link below to verify your email address:<br><br>"+
			"<a href='%s'>%s</a>",
		user.Name,
		verificationLink,
		verificationLink,
	)

	if err := mail.SendEmail(user.Email, subject, body); err != nil {
		fmt.Println("Error sending email:", err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully. Please check your email for verification.",
		"user_id": user.UserID,
		"email":   user.Email,
		"role":    user.Role,
	})
}
func LoginUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	authUser, err := models.Authenticate(user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	userid, _ := models.GetUserIDbyEmail(user.Email)
	claims := &Claims{
		Email:  authUser.Email,
		UserID: userid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Login successful",
		"token":      tokenString,
		"user_id":    authUser.UserID,
		"role":       authUser.Role,
		"email":      authUser.Email,
		"created_at": authUser.CreatedAt,
		"updated_at": authUser.UpdatedAt,
	})
}
func LogoutUser(c *gin.Context) {
	// Invalidate the token by setting an empty token with an immediate expiration
	c.SetCookie("token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}
func VerifyEmail(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token is required"})
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired verification token"})
		return
	}

	if err := models.VerifyUserEmail(claims.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email"})
		return
	}

	c.Redirect(http.StatusFound, "/static/verification_success.html")
}
func generateVerificationToken(userID int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
