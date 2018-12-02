package main

import (
	"fmt"
	"strconv"

	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	UserCompositeKey  = "userCompositeKey"
	EntryCompositeKey = "entryCompositeKey"
)

// user
type UserType struct { // entry
	UserTid         string `json:"userTid"`
	UserName        string `json:"userName"`
	EntryDate       string `json:"entryDate"`
	RequestedAmount int64  `json:"requestedAmount"`
}

// entry
type EntryType struct { // entry
	UserTid          string `json:"userTid"`
	SourceId         string `json:"sourceId"`
	EntryId          string `json:"entryId"`
	EntryAccountType string `json:"entityAccountType"`
	EntryType        string `json:"entityType"`
	EntryName        string `json:"entityName"`
	EntryDate        string `json:"entityDate"`
	Amount           int64  `json:"amount"`
}

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Init chaincode")

	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "createUser" {
		// Create Entry
		return t.createUser(stub, args)
	} else if function == "retrieveUser" {
		// Retrieve Entry
		return t.retrieveUser(stub, args)
	} else if function == "createEntry" {
		// Create Entry
		return t.createEntry(stub, args)
	} else if function == "retrieveEntries" {
		// Retrieve Entries
		return t.retrieveEntries(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

func (t *SimpleChaincode) createUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var userKey string

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}
	fmt.Printf("- createEntity(FirstName=%q,LastName=%q,PhoneNumber=%q,Email=%q,TokenAmount=%q,IsIssuer=%q)\n", args[0], args[1], args[2], args[3], args[4], args[5])

	//
	user := UserType{}
	user.UserTid = args[0]
	user.UserName = args[1]
	user.EntryDate = args[2]
	user.RequestedAmount, err = strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		return shim.Error("Expecting integer value for RequestedAmount")
	}

	userKey, _ = stub.CreateCompositeKey(UserCompositeKey, []string{user.UserTid})

	userJsonAsBytes, _ := json.Marshal(user)
	// Write the state to the ledger
	err = stub.PutState(userKey, userJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("- created user=%s\n", userJsonAsBytes)
	return shim.Success(userJsonAsBytes)
}

// Retrieve Entity
func (t *SimpleChaincode) retrieveUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var userKey string

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	fmt.Printf("- retrieveUser(userTId=%q)\n", args[0])
	var userTid = args[0]

	userKey, _ = stub.CreateCompositeKey(UserCompositeKey, []string{userTid})
	userAsBytes, err := stub.GetState(userKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("- returned user=%s\n", userAsBytes)
	return shim.Success(userAsBytes)
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) createEntry(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var entryKey string

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}
	fmt.Printf("- createEntity(FirstName=%q,LastName=%q,PhoneNumber=%q,Email=%q,TokenAmount=%q,IsIssuer=%q)\n", args[0], args[1], args[2], args[3], args[4], args[5])
	//
	entry := EntryType{}
	entry.UserTid = args[0]
	entry.EntryId = args[1]
	entry.EntryAccountType = args[2]
	entry.EntryType = args[3]
	entry.EntryName = args[4]
	entry.EntryDate = args[5]
	entry.Amount, err = strconv.ParseInt(args[6], 10, 64)
	if err != nil {
		return shim.Error("Expecting integer value for Entry")
	}

	entryKey, _ = stub.CreateCompositeKey(EntryCompositeKey, []string{entry.UserTid, entry.EntryId})

	entryJsonAsBytes, _ := json.Marshal(entry)
	// Write the state to the ledger
	err = stub.PutState(entryKey, entryJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("- return=%s\n", entryJsonAsBytes)
	return shim.Success(entryJsonAsBytes)
}

// Deletes an entity from state
func (t *SimpleChaincode) retrieveEntries(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
