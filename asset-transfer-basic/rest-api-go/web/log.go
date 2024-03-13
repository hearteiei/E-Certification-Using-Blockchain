package web

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"
	"time"
)

var users = make(map[string]User)

type User struct {
	UserID       string `json:"userid"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Mail         string `json:"mail"`
	Password     string `json:"password"`
	OTP          string
	OTPTimestamp int64
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, exists := users[newUser.UserID]; exists {
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

	users[newUser.UserID] = newUser
	fmt.Fprintf(w, "User %s registered successfully. OTP sent to %s\n", newUser.UserID, newUser.Mail)
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
		UserID string `json:"userid"`
		OTP    string `json:"otp"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := users[loginData.UserID]
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

	fmt.Fprintf(w, "User %s logged in successfully\n", loginData.UserID)
}

func resendOTP(userID string) error {
	// Generate a new OTP
	newOTP := generateOTP()

	// Update OTP and timestamp for the user
	user, ok := users[userID]
	if !ok {
		return fmt.Errorf("User not found")
	}
	user.OTP = newOTP
	user.OTPTimestamp = time.Now().Unix()
	users[userID] = user

	// Send the new OTP to the user's email address (replace with your email sending logic)
	fmt.Printf("New OTP sent to %s: %s\n", user.Mail, newOTP)

	return nil
}

func resendOTPHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		UserID string `json:"userid"`
	}

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Resend OTP for the specified user ID
	if err := resendOTP(requestData.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	fmt.Fprintf(w, "New OTP sent to user %s\n", requestData.UserID)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		UserID   string `json:"userid"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, exists := users[loginData.UserID]
	if !exists || user.Password != loginData.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "User %s logged in successfully\n", loginData.UserID)
}

func generateOTP() string {
	return strconv.Itoa(1000 + rand.Intn(9000))
}
