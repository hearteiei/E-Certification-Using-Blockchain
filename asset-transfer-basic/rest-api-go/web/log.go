package web

import (
	"fmt"
	"net/http"
)

var users = make(map[string]string)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if _, exists := users[username]; exists {
		fmt.Fprintf(w, "Username %s already exists\n", username)
		return
	}

	users[username] = password
	fmt.Fprintf(w, "User %s registered successfully\n", username)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	storedPassword, exists := users[username]
	if !exists || storedPassword != password {
		fmt.Fprintf(w, "Invalid username or password\n")
		return
	}

	fmt.Fprintf(w, "User %s logged in successfully\n", username)
}

// import (
// 	"encoding/json"
// 	"fmt"

// 	"math/rand"
// 	"net/smtp"
// 	"strings"
// 	"time"

// 	"github.com/google/uuid"
// )

// type User struct {
// 	UserID       string `json:"userid"`
// 	Firstname    string `json:"firstname"`
// 	Lastname     string `json:"lastname"`
// 	Mail         string `json:"mail"`
// 	Password     string `json:"Password"`
// 	OTP          string
// 	OTPTimestamp int64
// }

// func generateUserID(firstName, lastName string) string {
// 	// Convert first name and last name to lowercase and remove spaces
// 	firstName = strings.ToLower(strings.TrimSpace(firstName))
// 	lastName = strings.ToLower(strings.TrimSpace(lastName))

// 	// Generate a UUID
// 	id := uuid.New().String()

// 	// Combine first name, last name, and UUID to create a unique user ID
// 	userID := fmt.Sprintf("%s_%s_%s", firstName, lastName, id)

// 	return userID
// }

// func UserExists(userID string) (bool, error) {
// 	// Retrieve user information from the world state
// 	userJSON, err := GetStub().GetState(userID)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to read from world state: %v", err)
// 	}

// 	// Check if user exists by verifying if userJSON is not nil
// 	return userJSON != nil, nil
// }
// func generateOTP(length int) string {
// 	// Define the character set from which to generate the OTP
// 	characters := "0123456789"

// 	// Use time.Now().UnixNano() as the seed for the local random generator
// 	randSource := rand.NewSource(time.Now().UnixNano())

// 	// Create a local random generator
// 	random := rand.New(randSource)

// 	// Initialize an empty string to store the OTP
// 	otp := ""

// 	// Generate the OTP by randomly selecting characters from the character set
// 	for i := 0; i < length; i++ {
// 		// Generate a random index to select a character from the character set
// 		index := random.Intn(len(characters))

// 		// Append the randomly selected character to the OTP
// 		otp += string(characters[index])
// 	}

// 	return otp
// }
// func sendOTPByEmail(email, otp string) error {
// 	// Set up authentication information
// 	auth := smtp.PlainAuth("", "faaknorn261405@gmail.com", "faaknorn123", "smtp.gmail.com")

// 	// Set up email content
// 	from := "faaknorn261405@gmail.com"
// 	to := []string{email}
// 	subject := "Your OTP"
// 	body := "Your OTP is: " + otp

// 	// Compose the email message
// 	msg := []byte("To: " + email + "\r\n" +
// 		"Subject: " + subject + "\r\n" +
// 		"\r\n" +
// 		body + "\r\n")

// 	// Send the email
// 	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // RegisterUserWithTimeLimitedOTP registers a new user with a certificate and sends a time-limited OTP to their email for verification
// func RegisterUserWithTimeLimitedOTP(firstName string, lastName string, email string, password string) (string, error) {
// 	// Generate unique user ID
// 	userID := generateUserID(firstName, lastName)

// 	exists, err := UserExists(userID)
// 	if err != nil {
// 		return "", err
// 	}
// 	if exists {
// 		return "", fmt.Errorf("the user %s already exists", userID)
// 	}

// 	// Generate OTP and timestamp
// 	otp := generateOTP(6) // You need to implement this function
// 	timestamp := time.Now().Unix()

// 	// Send OTP to the user's email
// 	err2 := sendOTPByEmail(email, otp) // You need to implement this function
// 	if err2 != nil {
// 		return "", err2
// 	}

// 	// Store user information and OTP details in the world state
// 	user := User{
// 		UserID:       userID,
// 		Firstname:    firstName,
// 		Lastname:     lastName,
// 		Mail:         email,
// 		Password:     password,
// 		OTP:          otp,
// 		OTPTimestamp: timestamp,
// 	}
// 	userJSON, err := json.Marshal(user)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Store user in the world state
// 	err = ctx.GetStub().PutState(userID, userJSON)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to put user data into world state: %v", err)
// 	}

// 	return userID, nil
// }

// // VerifyTimeLimitedOTP verifies the time-limited OTP sent to the user's email during registration
// func (s *SmartContract) VerifyTimeLimitedOTP(ctx contractapi.TransactionContextInterface, userID string, otp string) (bool, error) {
// 	// Retrieve user information from the world state
// 	userJSON, err := ctx.GetStub().GetState(userID)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to read user from world state: %v", err)
// 	}
// 	if userJSON == nil {
// 		return false, fmt.Errorf("the user %s does not exist", userID)
// 	}

// 	var user User
// 	err = json.Unmarshal(userJSON, &user)
// 	if err != nil {
// 		return false, err
// 	}

// 	// Check if the provided OTP matches the stored OTP
// 	if user.OTP == otp {
// 		// Check if OTP is still valid within the time limit (e.g., 5 minutes)
// 		currentTime := time.Now().Unix()
// 		if currentTime-user.OTPTimestamp <= 300 { // 300 seconds = 5 minutes
// 			return true, nil
// 		}
// 	}
// 	return false, nil
// }

// // LoginUserWithOTP logs in a user using OTP verification
// func (s *SmartContract) LoginUserWithOTP(ctx contractapi.TransactionContextInterface, mail string, otp string) (bool, error) {
// 	// Retrieve the user ID associated with the provided email
// 	userID, err := s.getUserIDByEmail(ctx, mail)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to retrieve user ID by email: %v", err)
// 	}
// 	if userID == "" {
// 		return false, fmt.Errorf("no user with email %s found", mail)
// 	}

// 	// Verify OTP for the user
// 	verified, err := s.VerifyTimeLimitedOTP(ctx, userID, otp)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to verify OTP: %v", err)
// 	}
// 	if !verified {
// 		return false, nil // OTP verification failed
// 	}

// 	// OTP verification successful, proceed with login
// 	return true, nil
// }

// // Helper function to get user ID by email
// func (s *SmartContract) getUserIDByEmail(ctx contractapi.TransactionContextInterface, mail string) (string, error) {
// 	// Create a composite key for the user based on email
// 	compositeKey, err := ctx.GetStub().CreateCompositeKey("user", []string{"mail", mail})
// 	if err != nil {
// 		return "", err
// 	}

// 	// Retrieve the user's ID using the composite key
// 	userJSON, err := ctx.GetStub().GetState(compositeKey)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to read user data from world state: %v", err)
// 	}
// 	if userJSON == nil {
// 		// User with the given email not found
// 		return "", nil
// 	}

// 	// Extract the user ID from the composite key
// 	_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(compositeKey)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to split composite key: %v", err)
// 	}
// 	userID := compositeKeyParts[1]

// 	return userID, nil
// }
