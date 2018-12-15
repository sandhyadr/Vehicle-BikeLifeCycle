
package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the bike structure, with 6 properties.  Structure tags are used by encoding/json library
type Bike struct {
	Make   string `json:"make"`
	Model  string `json:"model"`
	Colour string `json:"colour"`
	Vin    string `json:"vin"`
	EngineCC string `json:"engcc"`
	Owner  string `json:"owner"`
}

type SmartContract struct {
}

// Define the person structure, with 4 properties.  Structure tags are used by encoding/json library
type Person struct {
	Name   string `json:"name"`
	Address string `json:"address"`
	Email   string `json:"email"`
	Owner  string `json:"owner"`
}




/*
 * The Init method is called when the Smart Contract "fabBike" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabBike"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryBike" {
		return s.queryBike(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createBike" {
		return s.createBike(APIstub, args)
	} else if function == "queryAllBikes" {
		return s.queryAllBikes(APIstub)
	} else if function == "changeBikeOwner" {
		return s.changeBikeOwner(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryBike(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	BikeAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(BikeAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	Bikes := []Bike{
		Bike{Make: "Honda", Model: "Activa", Colour: "gold",Vin: "567SDA12",EngineCC: "100CC",Owner: "Sandhya"},
		Bike{Make: "Vespa", Model: "Activa", Colour: "red",Vin: "345VLS67",EngineCC: "150CC" ,Owner: "Lakshmi"},
		Bike{Make: "Bajaj", Model: "AXv", Colour: "blue",Vin: "135AB79",EngineCC: "150CC" ,Owner: "Satwik"},
		Bike{Make:  "HeroHonda", Model: "Infinit", Colour: "grey",Vin: "145HG45",EngineCC: "100CC" , Owner: "Prabhu"},
		Bike{Make:  "Yamaha", Model: "", Colour: "metallic",Vin: "179NF68",EngineCC: "150CC" , Owner: "Syed"},
		Bike{Make:  "TVS", Model: "", Colour: "green",Vin: "236KJ12",EngineCC: "150CC" , Owner: "Anand"},
		Bike{Make: "Honda" , Model: "Kinetic", Colour: "silver",Vin: "934DS21",EngineCC: "150CC" , Owner: "Seshu"},
		Bike{Make:  "", Model: "", Colour: "yellow",Vin: "689LK12",EngineCC: "100CC" , Owner: "Pavan"},
		Bike{Make:  "Honda", Model: "Kinetic", Colour: "black",Vin: "123MHK45",EngineCC: "100CC" , Owner: "Vani"},
		Bike{Make: "", Model: "Pleasure", Colour: "white",Vin: "4566VMS78",EngineCC: "50CC" , Owner: "Madhuri"},
	    Bike{Make:  "", Model: "", Colour: "pink",Vin: "",EngineCC: "100CC" , Owner: "Yogita"},	
	}

	i := 0
	for i < len(Bikes) {
		fmt.Println("i is ", i)
		BikeAsBytes, _ := json.Marshal(Bikes[i])
		APIstub.PutState("Bike"+strconv.Itoa(i), BikeAsBytes)
		fmt.Println("Added", Bikes[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createBike(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	var Bike = Bike{Make: args[1], Model: args[2], Colour: args[3], Vin: args[4],EngineCC: args[5],Owner: args[6]}

	BikeAsBytes, _ := json.Marshal(Bike)
	APIstub.PutState(args[0], BikeAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllBikes(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "Bike0"
	endKey := "Bike999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllBikes:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (t *SimpleChaincode) getBikeDetails(stub shim.ChaincodeStubInterface, args []string) sc.Response{
	
	return shim.Success(nil)

}






func (s *SmartContract) changeBikeOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	BikeAsBytes, _ := APIstub.GetState(args[0])
	Bike := Bike{}

	json.Unmarshal(BikeAsBytes, &Bike)
	Bike.Owner = args[0]

	BikeAsBytes, _ = json.Marshal(Bike)
	APIstub.PutState(args[0], BikeAsBytes)

	return shim.Success(nil)
}









// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
