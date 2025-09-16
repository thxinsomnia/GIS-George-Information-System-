package attendances

import (
	"errors"
	"log"
	"GIS/config"
	"GIS/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Custom participant struct for the payload
type ParticipantPayloadByName struct {
	UserName   string `json:"user_name" binding:"required"`
	BonusValor int    `json:"bonus_valor"`
}

// Main payload now uses names
type ManualEventPayloadByName struct {
	EventName    string                     `json:"event_name" binding:"required"`
	EventDate    string                     `json:"event_date" binding:"required"`
	EventTypeName  string                   `json:"event_type" binding:"required"`
	BaseValor    int                        `json:"base_valor" binding:"required"`
	Participants []ParticipantPayloadByName `json:"participants" binding:"required"`
}

func Attendance(c *gin.Context) {
	var payload ManualEventPayloadByName
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload: " + err.Error()})
		return
	}

	parsedDate, err := time.Parse("2006-01-02", payload.EventDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
		return
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Get EventType to determine points
		var eventType models.Type
		if err := tx.First(&eventType, "type_name = ?", payload.EventTypeName).Error; err != nil {
			return errors.New("event type not found")
		}
		

		newEvent := models.Event{
			EventName: payload.EventName,
			EventTime: parsedDate,
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"), // or use time.Now() if your model expects time.Time
			TypeId:  eventType.TypeId,
		}
		
		if err := tx.Create(&newEvent).Error; err != nil {
			return err
		}

		// 3. Process each participant
		for _, participant := range payload.Participants {
			// 1. Find the user by name to get their UUID
			var userProfile models.User
			if err := tx.Where("name = ?", participant.UserName).First(&userProfile).Error; err != nil {
				// If a user is not found, you can choose to skip or fail the entire transaction
				log.Printf("User with name '%s' not found, skipping.", participant.UserName)
				continue
			}
			userID := userProfile.Id // <-- This is the safe UUID

			// 2. Proceed with the logic using the UUID
			valorTotal := payload.BaseValor + participant.BonusValor
			
			// Create attendance record
			attendanceRecord := models.Attendance{
				SoldierId:  userID,
				EventId:    int64(newEvent.EventId),
				CreatedAt:  time.Now().Format("2006-01-02 15:04:05"),
				ValorEarned: int64(valorTotal),
			}
			if err := tx.Create(&attendanceRecord).Error; err != nil {
				log.Printf("Could not record attendance for user %s: %v", participant.UserName, err)
				continue
			}

			// Add points to user's profile
			result := tx.Model(&models.User{}).Where("id = ?", userID).
				Update("total_valor", gorm.Expr("total_valor + ?", valorTotal))
			if result.Error != nil {
				return result.Error
			}

			// Check for rank-up
			// if err := checkAndUpdateUserRank(tx, participant.UserID); err != nil {
			// 	return err
			// }
		}

		return nil // Commit transaction
	})

	if err != nil {
		// Handle transaction errors...
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event and record attendance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event created and attendance recorded successfully"})
}
