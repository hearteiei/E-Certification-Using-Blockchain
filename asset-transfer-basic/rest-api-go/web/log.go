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

	otp := generateOTP()

	// Send OTP to email
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

func sendEmail(to, subject, body string) error {

	from := "test259492@gmail.com"
	password := "rijq jocq csnq lhmq"

	host := "smtp.gmail.com"
	port := "587"

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
	users[loginData.Mail] = user
	fmt.Fprintf(w, `{"message": "User %s logged in successfully"}`, loginData.Mail)
}

type LoginResponse struct {
	User        User   `json:"user"`
	AccessToken string `json:"access_token"`
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {

	var allUsers []User

	for _, user := range users {
		allUsers = append(allUsers, user)
	}

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

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["mail"] = mail
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

	accessToken, err := token.SignedString([]byte("259492")) // Replace "secret_key" with your own secret key
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func generateOTP() string {
	return strconv.Itoa(1000 + rand.Intn(9000))
}
