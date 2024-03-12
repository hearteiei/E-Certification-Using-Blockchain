package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
type Diploma struct {
	ID            string `json:"ID"`
	StudentName   string `json:"studentName"`
	TeacherName   string `json:"teacherName"`
	DiplomaNumber string `json:"diplomaNumber"`
	SubjectTopic  string `json:"subjectTopic"`
	Issuer        string `json:"issuer"`
	IssuedDate    string `json:"issuedDate"`
	BeginDate     string `json:"beginDate"`
	EndDate       string `json:"endDate"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	diplomas := []Diploma{
		{
			ID:            "asset1",
			StudentName:   "kunasin techasueb",
			TeacherName:   "Dom pothingan",
			DiplomaNumber: "1",
			SubjectTopic:  "Fullstack Developement",
			Issuer:        "CMU-Eleaning",
			IssuedDate:    "2023-3-2",
			BeginDate:     "2023-2-30",
			EndDate:       "2023-3-2"},
		{
			ID:            "asset2",
			StudentName:   "kontakan kamfoo",
			TeacherName:   "kampong woradut",
			DiplomaNumber: "2",
			SubjectTopic:  "Fullstack Developement",
			Issuer:        "CMU-Eleaning",
			IssuedDate:    "2023-12-15",
			BeginDate:     "2023-12-13",
			EndDate:       "2023-12-2",
		},
		{
			ID:            "asset3",
			StudentName:   "Ronaldo",
			TeacherName:   "messi",
			DiplomaNumber: "3",
			SubjectTopic:  "How To dribbing",
			Issuer:        "Barcelona FC",
			IssuedDate:    "2022-3-2",
			BeginDate:     "2022-2-30",
			EndDate:       "2022-3-2",
		},
		{
			ID:            "asset4",
			StudentName:   "stephen curry",
			TeacherName:   "jame harden",
			DiplomaNumber: "4",
			SubjectTopic:  "How to shoot 3 point",
			Issuer:        "Golden warriors",
			IssuedDate:    "2021-3-2",
			BeginDate:     "2021-2-30",
			EndDate:       "2021-3-2",
		},
		{
			ID:            "asset5",
			StudentName:   "Buakaw Bunchamek",
			TeacherName:   "Rodthang jitmueangnon",
			DiplomaNumber: "5",
			SubjectTopic:  "How to knock in first round",
			Issuer:        "one championship",
			IssuedDate:    "2020-3-2",
			BeginDate:     "2020-2-30",
			EndDate:       "2020-3-2",
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
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, studentname string, teacherName string, diplomaNumber string, subjectTopic string, issuer string, issuedDate string, beginDate string, endDate string) error {
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
		TeacherName:   teacherName,
		DiplomaNumber: diplomaNumber,
		SubjectTopic:  subjectTopic,
		Issuer:        issuer,
		IssuedDate:    issuedDate,
		BeginDate:     beginDate,
		EndDate:       endDate,
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
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, studentname string, teacherName string, diplomaNumber string, subjectTopic string, issuer string, issuedDate string, beginDate string, endDate string) error {
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
		TeacherName:   teacherName,
		DiplomaNumber: diplomaNumber,
		SubjectTopic:  subjectTopic,
		Issuer:        issuer,
		IssuedDate:    issuedDate,
		BeginDate:     beginDate,
		EndDate:       endDate,
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
