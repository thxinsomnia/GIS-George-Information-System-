package auth

import (
	"GIS/config"
	"GIS/models"
	"errors"
	"net/http"
	"strings"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/gotrue-go/types"
	"gorm.io/gorm"

)

type LoginPayload struct {
    Email    string `json:"email" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var payload LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fill all the required field! " + err.Error()})
		return
	}

	// 2. Hubungi Supabase Auth untuk verifikasi email dan password
	response, err := config.SupabaseClient.Auth.SignInWithEmailPassword(
	 	payload.Email,
		payload.Password,
	)

	// Jika ada error dari Supabase (misal, password salah), kirim response Unauthorized
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong email or password! " + err.Error()})
		return
	}

	// 3. Jika berhasil, Supabase akan memberikan token. Kirim token ini ke client.
	c.JSON(http.StatusOK, gin.H{
		"message":      "Login successful",
		"access_token": response.AccessToken, // Ini adalah JWT dari Supabase
		"user_id":      response.User.ID,
	})
}

// 1. Payload disesuaikan dengan alur baru
type RegistrationPayload struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func ActivateAccount(c *gin.Context) {
	var payload RegistrationPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Gunakan Transaksi untuk menjaga konsistensi data
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Verifikasi NAMA di tabel USERS, pastikan belum aktif
		var user models.User
		if err := tx.Where("name = ? AND status = ?", payload.Name, false).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Nama tidak ada ATAU akun sudah aktif
				return errors.New("name not registered or account already activated")
			}
			return err // Error database lainnya
		}

		requestBody := types.SignupRequest{
			Email:    payload.Email,
			Password: payload.Password,
		}

		// 2. Buat user di Supabase Auth dengan EMAIL dan PASSWORD
		newUser, err := config.SupabaseClient.Auth.Signup(requestBody)
			
		if err != nil {
			// 3. Check if the error is specifically about a duplicate user
			log.Println("Detailed signup error:", err)
    		if strings.Contains(err.Error(), "User already registered") {
        		return errors.New("Email already in use, please choose another email")
    		}
    
			// For all other errors, return a more general message
			return errors.New("Failed to register user, server error occurred")
		}

		// 3. Update baris di tabel UserS yang ditemukan tadi
		updateData := models.User{
			Status: true,          // Set status aktif
			Email:    payload.Email, 
			Id:       newUser.User.ID,    // Sinkronkan ID dengan Supabase Auth
		}
		
		// Gunakan ID unik dari profil yang ditemukan untuk update yang lebih aman
		result := tx.Model(&models.User{}).Where("id = ?", user.Id).Updates(updateData)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("Failed to Update Account Information!")
		}

		// Jika semua berhasil, return nil untuk commit transaksi
		return nil
	})

	// Penanganan error dari hasil transaksi
	if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account successfully activated! Please login."})
}
