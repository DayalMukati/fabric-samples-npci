package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a boat
type SmartContract struct {
	contractapi.Contract
}

// Boat describes basic details of what makes up a boat
type Boat struct {
	Make   string `json:"make"`
	Model  string `json:"model"`
	Colour string `json:"colour"`
	Owner  string `json:"owner"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Boat
}

// InitLedger adds a base set of boats to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	boats := []Boat{
		Boat{Make: "Toyota", Model: "Prius", Colour: "blue", Owner: "Tomoko"},
		Boat{Make: "Ford", Model: "Mustang", Colour: "red", Owner: "Brad"},
		Boat{Make: "Hyundai", Model: "Tucson", Colour: "green", Owner: "Jin Soo"},
		Boat{Make: "Volkswagen", Model: "Passat", Colour: "yellow", Owner: "Max"},
		Boat{Make: "Tesla", Model: "S", Colour: "black", Owner: "Adriana"},
		Boat{Make: "Peugeot", Model: "205", Colour: "purple", Owner: "Michel"},
		Boat{Make: "Chery", Model: "S22L", Colour: "white", Owner: "Aarav"},
		Boat{Make: "Fiat", Model: "Punto", Colour: "violet", Owner: "Pari"},
		Boat{Make: "Tata", Model: "Nano", Colour: "indigo", Owner: "Valeria"},
		Boat{Make: "Holden", Model: "Barina", Colour: "brown", Owner: "Shotaro"},
	}

	for i, boat := range boats {
		boatAsBytes, _ := json.Marshal(boat)
		err := ctx.GetStub().PutState("BOAT"+strconv.Itoa(i), boatAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateBoat adds a new boat to the world state with given details
func (s *SmartContract) CreateBoat(ctx contractapi.TransactionContextInterface, boatNumber string, make string, model string, colour string, owner string) error {
	boat := Boat{
		Make:   make,
		Model:  model,
		Colour: colour,
		Owner:  owner,
	}
	boatAsBytes, _ := json.Marshal(boat)

	return ctx.GetStub().PutState(boatNumber, boatAsBytes)
}


// QueryBoat returns the boat stored in the world state with given id
func (s *SmartContract) QueryBoat(ctx contractapi.TransactionContextInterface, boatNumber string) (*Boat, error) {
	boatAsBytes, err := ctx.GetStub().GetState(boatNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state: %v", err)
	}
	if boatAsBytes == nil {
		return nil, fmt.Errorf("Boat does not exist: %s", boatNumber)
	}

	boat := new(Boat)
	_ = json.Unmarshal(boatAsBytes, boat)

	return boat, nil
}

// QueryAllBoats returns all boats found in world state
func (s *SmartContract) QueryAllBoats(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := "BOAT0"
	endKey := "BOAT999"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var boats []QueryResult
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var boat Boat
		_ = json.Unmarshal(queryResponse.Value, &boat)

		boatResult := QueryResult{Key: queryResponse.Key, Record: &boat}
		boats = append(boats, boatResult)
	}

	return boats, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create boat chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting boat chaincode: %s", err.Error())
	}
}
