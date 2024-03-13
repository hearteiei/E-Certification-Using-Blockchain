package web

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var users = make(map[string]User)

type User struct {
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Mail         string `json:"mail"`
	Password     string `json:"password"`
	OTP          string
	OTPTimestamp int64
	Status       string
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, exists := users[newUser.Mail]; exists {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	// Generate OTP
	otp := generateOTP()

	// Send OTP to user's email
	err = sendEmail(newUser.Mail, "OTP for Registration", fmt.Sprintf("Your OTP is: %s", otp))
	if err != nil {
		http.Error(w, "Failed to send OTP via email %s", http.StatusInternalServerError)
		return
	}

	// Store OTP and timestamp in user struct
	newUser.OTP = otp
	newUser.OTPTimestamp = time.Now().Unix()
	newUser.Status = "Pending"

	users[newUser.Mail] = newUser

	fmt.Fprintf(w, `{"message": "User %s registered successfully. OTP sent to %s"}`, newUser.Firstname+newUser.Lastname, newUser.Mail)
}

// This function simulates sending an email. You'll need to replace it with your actual email sending logic.
func sendEmail(to, subject, body string) error {
	// Replace with your actual credentials (store securely!)
	from := "test259492@gmail.com"
	password := "rijq jocq csnq lhmq" // Never store passwords directly in code

	// Use a secure connection (replace with your provider details)
	host := "smtp.gmail.com"
	port := "587" // Common port for TLS

	// Fallback to PlainAuth (less secure, use with caution)
	auth := smtp.PlainAuth("", from, password, host)
	msg := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body
	err := smtp.SendMail(host+":"+port, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}

func checkOTPValidity(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Mail string `json:"mail"`
		OTP  string `json:"otp"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := users[loginData.Mail]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if user.OTP != loginData.OTP {
		http.Error(w, "Invalid OTP", http.StatusUnauthorized)
		return
	}

	// Check OTP expiry (let's say it's valid for 5 minutes)
	if time.Now().Unix()-user.OTPTimestamp > 300 {
		http.Error(w, "OTP expired", http.StatusUnauthorized)
		return
	}
	user.Status = "Success"
<<<<<<< HEAD
=======
	users[loginData.Mail] = user
>>>>>>> 852a34462e80bdadf12da0a137c427b0a059203c
	fmt.Fprintf(w, `{"message": "User %s logged in successfully"}`, loginData.Mail)
}

// func resendOTP(userID string) error {
// 	// Generate a new OTP
// 	newOTP := generateOTP()

// 	// Update OTP and timestamp for the user
// 	user, ok := users[userID]
// 	if !ok {
// 		return fmt.Errorf("User not found")
// 	}
// 	user.OTP = newOTP
// 	user.OTPTimestamp = time.Now().Unix()
// 	users[userID] = user

// 	// Send the new OTP to the user's email address (replace with your email sending logic)
// 	fmt.Printf("New OTP sent to %s: %s\n", user.Mail, newOTP)

// 	return nil
// }

// func resendOTPHandler(w http.ResponseWriter, r *http.Request) {
// 	var requestData struct {
// 		UserID string `json:"userid"`
// 	}

// 	// Decode request body
// 	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// Resend OTP for the specified user ID
// 	if err := resendOTP(requestData.UserID); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Respond with success message
// 	fmt.Fprintf(w, "New OTP sent to user %s\n", requestData.UserID)
// }

type LoginResponse struct {
	User        User   `json:"user"`
	AccessToken string `json:"access_token"`
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	// Create a slice to hold all user data
	var allUsers []User

	// Iterate over the map and append each user to the slice
	for _, user := range users {
		allUsers = append(allUsers, user)
	}

	// Encode the slice of users as JSON and write it to the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allUsers)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Mail     string `json:"mail"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, exists := users[loginData.Mail]
	if !exists || user.Password != loginData.Password || user.Status != "Success" {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT access token
	accessToken, err := generateAccessToken(user.Mail)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	// Return user information and access token
	response := LoginResponse{
		User:        user,
		AccessToken: accessToken,
	}
	json.NewEncoder(w).Encode(response)
}

func generateAccessToken(mail string) (string, error) {
	// Create a new JWT token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["mail"] = mail
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

	// Sign the token with a secret key
	accessToken, err := token.SignedString([]byte("259492")) // Replace "secret_key" with your own secret key
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func generateOTP() string {
	return strconv.Itoa(1000 + rand.Intn(9000))
}
