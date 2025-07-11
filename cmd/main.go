package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Puneet-Vishnoi/kafka-simple/config"
	"github.com/Puneet-Vishnoi/kafka-simple/kafka"
	"github.com/gin-gonic/gin"
)

func main() {
	// database.Init()

	r := gin.Default()

	// ---------------- ROUTES ----------------

	// r.POST("/signup", func(c *gin.Context) {
	// 	var signupData struct {
	// 		Email     string `json:"email"`
	// 		Password  string `json:"password"`
	// 		FirstName string `json:"first_name"`
	// 		LastName  string `json:"last_name"`
	// 		UserType  string `json:"user_type"`
	// 	}

	// 	if err := c.ShouldBindJSON(&signupData); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
	// 		return
	// 	}

	// 	exists, _ := models.UserExists(signupData.Email)
	// 	if exists {
	// 		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
	// 		return
	// 	}

	// 	hashedPassword, err := helpers.HashPassword(signupData.Password)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
	// 		return
	// 	}

	// 	userID := helpers.GenerateUUID()

	// 	err = models.CreateUser(models.User{
	// 		UserID:    userID,
	// 		Email:     signupData.Email,
	// 		Password:  hashedPassword,
	// 		FirstName: signupData.FirstName,
	// 		LastName:  signupData.LastName,
	// 		UserType:  signupData.UserType,
	// 	})
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
	// 		return
	// 	}

	// 	accessToken, refreshToken, err := helpers.GenerateAllTokens(
	// 		signupData.Email, signupData.FirstName, signupData.LastName, signupData.UserType, userID,
	// 	)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
	// 		return
	// 	}

	// 	c.JSON(http.StatusCreated, gin.H{
	// 		"message":       "Signup successful",
	// 		"access_token":  accessToken,
	// 		"refresh_token": refreshToken,
	// 		"user_id":       userID,
	// 		"email":         signupData.Email,
	// 	})
	// })

	// r.POST("/signin", func(c *gin.Context) {
	// 	var loginData struct {
	// 		Email    string `json:"email"`
	// 		Password string `json:"password"`
	// 	}

	// 	if err := c.ShouldBindJSON(&loginData); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
	// 		return
	// 	}

	// 	user, err := models.GetUserByEmail(loginData.Email)
	// 	if err != nil {
	// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
	// 		return
	// 	}

	// 	if !helpers.VerifyPassword(user.Password, loginData.Password) {
	// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
	// 		return
	// 	}

	// 	accessToken, refreshToken, err := helpers.GenerateAllTokens(
	// 		user.Email, user.FirstName, user.LastName, user.UserType, user.UserID,
	// 	)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
	// 		return
	// 	}

	// 	c.JSON(http.StatusOK, gin.H{
	// 		"access_token":  accessToken,
	// 		"refresh_token": refreshToken,
	// 		"user_id":       user.UserID,
	// 		"email":         user.Email,
	// 	})
	// })

	// r.POST("/refresh", func(c *gin.Context) {
	// 	var req struct {
	// 		RefreshToken string `json:"refresh_token"`
	// 	}

	// 	if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing refresh token"})
	// 		return
	// 	}

	// 	claims, err := helpers.ValidateToken(req.RefreshToken)
	// 	if err != nil {
	// 		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	newToken, newRefresh, err := helpers.GenerateAllTokens(
	// 		claims.Email, claims.FirstName, claims.LastName, claims.UserType, claims.Uid,
	// 	)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token regeneration failed"})
	// 		return
	// 	}

	// 	c.JSON(http.StatusOK, gin.H{
	// 		"access_token":  newToken,
	// 		"refresh_token": newRefresh,
	// 	})
	// })

	// r.GET("/protected", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
	// 	user := c.MustGet("user").(string)
	// 	c.JSON(http.StatusOK, gin.H{"message": "Welcome " + user})
	// })

	cfg := config.LoadConfig()
	// // Ensure DLQ topic exists before starting DLQ writer
	// if err := kafka.EnsureTopicExists([]string{cfg.Brokers}, cfg.DLQTopic); err != nil {
	// 	log.Fatalf("DLQ topic verification/creation failed: %v", err)
	// }
	// Kafka DLQ writer
	dlqWriter := kafka.InitDLQWriter(cfg)
	defer dlqWriter.Close()

	// Kafka Producer and Consumer instances
	producer := kafka.NewKafkaProducer(cfg)
	defer producer.Close()

	consumer := kafka.NewKafkaConsumer(cfg)
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// // Start producing sample message asynchronously
	// go func() {
	// 	if err := producer.Produce(ctx, []byte("order-created")); err != nil {
	// 		log.Printf("Producer error: %v", err)
	// 	}
	// }()

	r.POST("/produce", func(c *gin.Context) {
		var messages []string

		if err := c.ShouldBindJSON(&messages); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Error parsing messages",
			})
			return
		}

		for _, msg := range messages {
			if err := producer.Produce(ctx, []byte(msg)); err != nil {
				log.Printf("Producer error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to produce message",
				})
				return
			}
		}

		c.JSON(http.StatusAccepted, gin.H{
			"success": "Messages sent successfully",
		})
	})

	// Start consuming messages with retry and DLQ handling
	go consumer.Consume(ctx, kafka.ProcessMessage, dlqWriter)

	// ---------------- START SERVER WITH GRACEFUL SHUTDOWN ----------------
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Channel to listen for interrupt or terminate signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// | Action            | Signal         | Result                           |
	// | ----------------- | -------------- | -------------------------------- |
	// | Press `Ctrl+C`    | `os.Interrupt` | Sent to `quit` → shutdown starts |
	// | `docker stop app` | `SIGTERM`      | Sent to `quit` → shutdown starts |
	// | No signal         | —              | App keeps running                |

	go func() {
		log.Println("Server running on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %v\n", err)
		}
	}()

	<-quit // Wait here until signal is received
	log.Println("Shutting down server...")

	// Shutdown with timeout context
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
