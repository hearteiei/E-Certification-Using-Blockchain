package chaincode

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
type DiplomaStatus string

var idCounter int = 0

const (
	DiplomaStatusPending DiplomaStatus = "Pending"
	DiplomaStatusSuccess DiplomaStatus = "Success"
)

type Diploma struct {
	ID            string `json:"ID"`
	StudentName   string `json:"studentName"`
	Endorser_name string `json:"endorser_name"`
	Mail          string `json:"mail"`
	Course        string `json:"course"`
	Issuer        string `json:"issuer"`
	IssuedDate    string `json:"issuedDate"`
	Begin_date    string `json:"begin_Date"`
	End_date      string `json:"End_Date"`
	Status        string `json:"Status"`
}

type User struct {
	UserID       string `json:"userid"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Mail         string `json:"mail"`
	Password     string `json:"Password"`
	OTP          string
	OTPTimestamp int64
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	diplomas := []Diploma{
		{
			ID:            "asset1",
			StudentName:   "kunasin techasueb",
			Endorser_name: "Dom pothingan",
			Mail:          "kunasin@gmail.com",
			Course:        "Fullstack Developement",
			Issuer:        "CMU-Eleaning",
			IssuedDate:    "2023-3-2",
			Begin_date:    "2023-2-30",
			End_date:      "2023-3-2",
			Status:        "Success",
		},
		{
			ID:            "asset2",
			StudentName:   "kontakan kamfoo",
			Endorser_name: "kampong woradut",
			Mail:          "kontakan@gmail.com",
			Course:        "Fullstack Developement",
			Issuer:        "CMU-Eleaning",
			IssuedDate:    "2023-12-15",
			Begin_date:    "2023-12-13",
			End_date:      "2023-12-2",
			Status:        "Success",
		},
		{
			ID:            "asset3",
			StudentName:   "Ronaldo",
			Endorser_name: "messi",
			Mail:          "Ronaldo@gmail.com",
			Course:        "How To dribbing",
			Issuer:        "Barcelona FC",
			IssuedDate:    "2022-3-2",
			Begin_date:    "2022-2-30",
			End_date:      "2022-3-2",
			Status:        "Success",
		},
		{
			ID:            "asset4",
			StudentName:   "stephen curry",
			Endorser_name: "jame harden",
			Mail:          "stephen@gmail.com",
			Course:        "How to shoot 3 point",
			Issuer:        "Golden warriors",
			IssuedDate:    "2021-3-2",
			Begin_date:    "2021-2-30",
			End_date:      "2021-3-2",
			Status:        "Success",
		},
		{
			ID:            "asset5",
			StudentName:   "Buakaw Bunchamek",
			Endorser_name: "Rodthang jitmueangnon",
			Mail:          "Buakaw@gmail.com",
			Course:        "How to knock in first round",
			Issuer:        "one championship",
			IssuedDate:    "2020-3-2",
			Begin_date:    "2020-2-30",
			End_date:      "2020-3-2",
			Status:        "Success",
		},
	}

	for _, diploma := range diplomas {
		diplomaJSON, err := json.Marshal(diploma)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(diploma.ID, diplomaJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, studentname string, teacherName string, mail string, subjectTopic string, issuer string, issuedDate string, beginDate string, endDate string, status string) error {
	idCounter++
	id := strconv.Itoa(idCounter)
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	diploma := Diploma{
		ID:            id,
		StudentName:   studentname,
		Endorser_name: teacherName,
		Mail:          mail,
		Course:        subjectTopic,
		Issuer:        issuer,
		IssuedDate:    issuedDate,
		Begin_date:    beginDate,
		End_date:      endDate,
		Status:        status,
	}
	diplomaJSON, err := json.Marshal(diploma)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, diplomaJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Diploma, error) {
	diplomaJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if diplomaJSON == nil {
		return nil, fmt.Errorf("the diploma %s does not exist", id)
	}

	var diploma Diploma
	err = json.Unmarshal(diplomaJSON, &diploma)
	if err != nil {
		return nil, err
	}

	return &diploma, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, studentname string, teacherName string, mail string, subjectTopic string, issuer string, issuedDate string, beginDate string, endDate string, status string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// overwriting original asset with new asset
	diploma := Diploma{
		ID:            id,
		StudentName:   studentname,
		Endorser_name: teacherName,
		Mail:          mail,
		Course:        subjectTopic,
		Issuer:        issuer,
		IssuedDate:    issuedDate,
		Begin_date:    beginDate,
		End_date:      endDate,
		Status:        status,
	}
	diplomaJSON, err := json.Marshal(diploma)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, diplomaJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Diploma, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Diploma
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Diploma
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}
func generateUserID(firstName, lastName string) string {
	// Convert first name and last name to lowercase and remove spaces
	firstName = strings.ToLower(strings.TrimSpace(firstName))
	lastName = strings.ToLower(strings.TrimSpace(lastName))

	// Generate a UUID
	id := uuid.New().String()

	// Combine first name, last name, and UUID to create a unique user ID
	userID := fmt.Sprintf("%s_%s_%s", firstName, lastName, id)

	return userID
}

func (s *SmartContract) UserExists(ctx contractapi.TransactionContextInterface, userID string) (bool, error) {
	// Retrieve user information from the world state
	userJSON, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	// Check if user exists by verifying if userJSON is not nil
	return userJSON != nil, nil
}
func generateOTP(length int) string {
	// Define the character set from which to generate the OTP
	characters := "0123456789"

	// Use time.Now().UnixNano() as the seed for the local random generator
	randSource := rand.NewSource(time.Now().UnixNano())

	// Create a local random generator
	random := rand.New(randSource)

	// Initialize an empty string to store the OTP
	otp := ""

	// Generate the OTP by randomly selecting characters from the character set
	for i := 0; i < length; i++ {
		// Generate a random index to select a character from the character set
		index := random.Intn(len(characters))

		// Append the randomly selected character to the OTP
		otp += string(characters[index])
	}

	return otp
}
func sendOTPByEmail(email, otp string) error {
	// Set up authentication information
	auth := smtp.PlainAuth("", "faaknorn261405@gmail.com", "faaknorn123", "smtp.gmail.com")

	// Set up email content
	from := "faaknorn261405@gmail.com"
	to := []string{email}
	subject := "Your OTP"
	body := "Your OTP is: " + otp

	// Compose the email message
	msg := []byte("To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	// Send the email
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
	if err != nil {
		return err
	}

	return nil
}

// RegisterUserWithTimeLimitedOTP registers a new user with a certificate and sends a time-limited OTP to their email for verification
func (s *SmartContract) RegisterUserWithTimeLimitedOTP(ctx contractapi.TransactionContextInterface, firstName string, lastName string, email string, password string) (string, error) {
	// Generate unique user ID
	userID := generateUserID(firstName, lastName)

	exists, err := s.UserExists(ctx, userID)
	if err != nil {
		return "", err
	}
	if exists {
		return "", fmt.Errorf("the user %s already exists", userID)
	}

	// Generate OTP and timestamp
	otp := generateOTP(6) // You need to implement this function
	timestamp := time.Now().Unix()

	// Send OTP to the user's email
	err2 := sendOTPByEmail(email, otp) // You need to implement this function
	if err2 != nil {
		return "", err2
	}

	// Store user information and OTP details in the world state
	user := User{
		UserID:       userID,
		Firstname:    firstName,
		Lastname:     lastName,
		Mail:         email,
		Password:     password,
		OTP:          otp,
		OTPTimestamp: timestamp,
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	// Store user in the world state
	err = ctx.GetStub().PutState(userID, userJSON)
	if err != nil {
		return "", fmt.Errorf("failed to put user data into world state: %v", err)
	}

	return userID, nil
}

// VerifyTimeLimitedOTP verifies the time-limited OTP sent to the user's email during registration
func (s *SmartContract) VerifyTimeLimitedOTP(ctx contractapi.TransactionContextInterface, userID string, otp string) (bool, error) {
	// Retrieve user information from the world state
	userJSON, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return false, fmt.Errorf("failed to read user from world state: %v", err)
	}
	if userJSON == nil {
		return false, fmt.Errorf("the user %s does not exist", userID)
	}

	var user User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return false, err
	}

	// Check if the provided OTP matches the stored OTP
	if user.OTP == otp {
		// Check if OTP is still valid within the time limit (e.g., 5 minutes)
		currentTime := time.Now().Unix()
		if currentTime-user.OTPTimestamp <= 300 { // 300 seconds = 5 minutes
			return true, nil
		}
	}
	return false, nil
}

// LoginUserWithOTP logs in a user using OTP verification
func (s *SmartContract) LoginUserWithOTP(ctx contractapi.TransactionContextInterface, mail string, otp string) (bool, error) {
	// Retrieve the user ID associated with the provided email
	userID, err := s.getUserIDByEmail(ctx, mail)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve user ID by email: %v", err)
	}
	if userID == "" {
		return false, fmt.Errorf("no user with email %s found", mail)
	}

	// Verify OTP for the user
	verified, err := s.VerifyTimeLimitedOTP(ctx, userID, otp)
	if err != nil {
		return false, fmt.Errorf("failed to verify OTP: %v", err)
	}
	if !verified {
		return false, nil // OTP verification failed
	}

	// OTP verification successful, proceed with login
	return true, nil
}

// Helper function to get user ID by email
func (s *SmartContract) getUserIDByEmail(ctx contractapi.TransactionContextInterface, mail string) (string, error) {
	// Create a composite key for the user based on email
	compositeKey, err := ctx.GetStub().CreateCompositeKey("user", []string{"mail", mail})
	if err != nil {
		return "", err
	}

	// Retrieve the user's ID using the composite key
	userJSON, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return "", fmt.Errorf("failed to read user data from world state: %v", err)
	}
	if userJSON == nil {
		// User with the given email not found
		return "", nil
	}

	// Extract the user ID from the composite key
	_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(compositeKey)
	if err != nil {
		return "", fmt.Errorf("failed to split composite key: %v", err)
	}
	userID := compositeKeyParts[1]

	return userID, nil
}

// Define an additional struct to hold the index data
type CourseIssuerIndex struct {
	Course string `json:"course"`
	Issuer string `json:"issuer"`
}

// Create an index entry mapping course to issuer
func (s *SmartContract) CreateCourseIssuerIndex(ctx contractapi.TransactionContextInterface, course string, issuer string) error {
	indexKey, err := ctx.GetStub().CreateCompositeKey("courseIssuerIndex", []string{course, issuer})
	if err != nil {
		return err
	}
	// Store the index entry in the world state
	return ctx.GetStub().PutState(indexKey, []byte{})
}

// Update the index entry mapping course to issuer
func (s *SmartContract) UpdateCourseIssuerIndex(ctx contractapi.TransactionContextInterface, course string, oldIssuer string, newIssuer string) error {
	// Delete the existing index entry
	oldIndexKey, err := ctx.GetStub().CreateCompositeKey("courseIssuerIndex", []string{course, oldIssuer})
	if err != nil {
		return err
	}
	err = ctx.GetStub().DelState(oldIndexKey)
	if err != nil {
		return err
	}
	// Create a new index entry
	newIndexKey, err := ctx.GetStub().CreateCompositeKey("courseIssuerIndex", []string{course, newIssuer})
	if err != nil {
		return err
	}
	// Store the updated index entry in the world state
	return ctx.GetStub().PutState(newIndexKey, []byte{})
}

// Query the issuer by course
func (s *SmartContract) QueryIssuerByCourse(ctx contractapi.TransactionContextInterface, course string) ([]string, error) {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("courseIssuerIndex", []string{course})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var issuers []string
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Extract issuer from composite key
		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return nil, err
		}
		issuer := compositeKeyParts[1]
		issuers = append(issuers, issuer)
	}

	return issuers, nil
}
