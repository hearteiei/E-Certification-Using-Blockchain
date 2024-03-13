package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"

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

type DiplomaInfo struct {
	Course     string `json:"course"`
	Status     string `json:"status"`
	IssuerDate string `json:"issuerDate"`
}

// GetDiplomasInfoByIssuer returns the course, status, and issuerDate of diplomas issued by a specific issuer
func (s *SmartContract) GetDiplomasInfoByIssuer(ctx contractapi.TransactionContextInterface, issuer string) ([]*DiplomaInfo, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	diplomasInfoMap := make(map[string]*DiplomaInfo) // Change to store pointers
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var diploma Diploma
		err = json.Unmarshal(queryResponse.Value, &diploma)
		if err != nil {
			return nil, err
		}

		if diploma.Issuer == issuer {
			key := diploma.Course + "|" + diploma.IssuedDate
			if existingDiploma, ok := diplomasInfoMap[key]; ok {
				// Update status based on priority
				if existingDiploma.Status == "Pending" {
					if diploma.Status == "Success" {
						existingDiploma.Status = "Success"
					}
				} else if diploma.Status == "Success" {
					existingDiploma.Status = "Success"
				}
			} else {
				diplomasInfoMap[key] = &DiplomaInfo{ // Store pointer to DiplomaInfo
					Course:     diploma.Course,
					Status:     diploma.Status,
					IssuerDate: diploma.IssuedDate,
				}
			}
		}
	}

	// Convert map to slice of DiplomaInfo
	var diplomasInfo []*DiplomaInfo
	for _, info := range diplomasInfoMap {
		diplomasInfo = append(diplomasInfo, info)
	}

	return diplomasInfo, nil
}
