
package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"time"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var productIndexStr = "_productindex"				//name for the key/value that will store a list of all products

type Product struct{
	Product_Id string `json:"product_id"`
	Category string `json:"category"`
	Product_Description string `json:"product_description"`
	Availability_Start_Date string `json:"availability_start_date"`
	Availability_End_Date string `json:"availability_end_date"`
	List_Price float64 `json:"list_price"`
	Currency string `json:"currency"`
	Price_Start_Date string `json:"price_start_date"`
	Price_End_Date string `json:"price_end_date"`
	User_Type string `json:"user_type"`
}

type Offering struct{
	Offering_ID string `json:"offering_id"`
	Offering_Category string `json:"offering_category"`
	Offering_Description string `json:"offering_description"`
	Availability_Start_Date string `json:"availability_start_date"`
	Availability_End_Date string `json:"availability_end_date"`
	Current_List_Price float64 `json:"current_list_price"`
	Currency string `json:"currency"`
	Price_Start_Date string `json:"price_start_date"`
	Price_End_Date string `json:"price_end_date"`
	Product_ID_01 string `json:"product_id_01"`
	Product_ID_02 string `json:"product_id_02"`
}
var offeringIndexStr = "_offeringindex"



// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	var Aval int
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	// Write the state to the ledger
	err = stub.PutState("abc", []byte(strconv.Itoa(Aval)))				//making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}

	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(productIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	var location []string
	jsonAsBytesForOffering, _ := json.Marshal(location)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(offeringIndexStr, jsonAsBytesForOffering)
	if err != nil {
		return nil, err
	}


	return nil, nil
}
// ============================================================================================================================
// Run - Our entry point for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func (t *SimpleChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}

// ============================================================================================================================
// Run - Our entry point for Invokcations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "delete_product" {										//deletes an entity from its state
		res, err := t.delete_product(stub, args)
		return res, err
	} else if function == "delete_offering" {										//deletes an entity from its state
		res, err := t.delete_offering(stub, args)
		return res, err
	} else if function == "write" {											//writes a value to the chaincode state
		return t.Write(stub, args)
	} else if function == "init_product" {									//create a new product
		return t.init_product(stub, args)
	} else if function == "set_user_type" {										//change user_type of a product
		res, err := t.set_user_type(stub, args)
		return res, err
	} else if function == "read_product_index" {
		return t.read_product_index(stub,args);
	} else if function == "read_offering_index" {
		return t.read_offering_index(stub,args);
	} else if function == "init_offering" {									//create a new product
		return t.init_offering(stub, args)
	} 

	fmt.Println("run did not find func: " + function)						//error

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)
	fmt.Println("Argument " + args[0])
	// Handle different functions
	if function == "read" {													//read a variable
		return t.read(stub, args)
	} else if function == "read_product_index" {
		return t.read_product_index(stub,args);
	} else if function == "read_offering_index" {
		return t.read_offering_index(stub,args);
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query")
}

// ============================================================================================================================
// Read - read a variable from chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	name = args[0]
	fmt.Println("Argument " + name)
	valAsbytes, err := stub.GetState(name)									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil													//send it onward
}
//====================================================

//Read Product index
func (t *SimpleChaincode) read_product_index(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	valAsbytes, err := stub.GetState("_productindex")									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil													//send it onward
}

//====================================================

//Read Offering index
func (t *SimpleChaincode) read_offering_index(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	valAsbytes, err := stub.GetState("_offeringindex")									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil													//send it onward
}
//=====================================================
// ============================================================================================================================
// Delete - remove a key/value pair from state
// ============================================================================================================================
func (t *SimpleChaincode) delete_product(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	name := args[0]
	err := stub.DelState(name)													//remove the key from chaincode state
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	//get the product index
	productsAsBytes, err := stub.GetState(productIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get product index")
	}
	var productIndex []string
	json.Unmarshal(productsAsBytes, &productIndex)								//un stringify it aka JSON.parse()

	//remove product from index
	for i,val := range productIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + name)
		if val == name{															//find the correct marble
			fmt.Println("found marble")
			productIndex = append(productIndex[:i], productIndex[i+1:]...)			//remove it
			for x:= range productIndex{											//debug prints...
				fmt.Println(string(x) + " - " + productIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(productIndex)									//save new index
	err = stub.PutState(productIndexStr, jsonAsBytes)
	return nil, nil
}

//=====================================================
// ============================================================================================================================
// Delete an offering
// ============================================================================================================================
func (t *SimpleChaincode) delete_offering(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	name := args[0]
	err := stub.DelState(name)													
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	//get the offering index
	offeringsAsBytes, err := stub.GetState(offeringIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get offering index")
	}
	var offeringIndex []string
	json.Unmarshal(offeringsAsBytes, &offeringIndex)								

	//remove offering from index
	for i,val := range offeringIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + name)
		if val == name{															
			offeringIndex = append(offeringIndex[:i], offeringIndex[i+1:]...)			
			for x:= range offeringIndex{											
				fmt.Println(string(x) + " - " + offeringIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(offeringIndex)									
	err = stub.PutState(offeringIndexStr, jsonAsBytes)
	return nil, nil
}




// ============================================================================================================================
// Write - write variable into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) Write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var name, value string // Entities
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
	}

	name = args[0]															//rename for funsies
	value = args[1]
	err = stub.PutState(name, []byte(value))								//write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// ============================================================================================================================
// Init Marble - create a new marble, store into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) init_product(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error

	//   0       1       2     3
	// "asdf", "blue", "35", "bob"
	if len(args) != 10 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	fmt.Println("- start init marble")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return nil, errors.New("5th argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return nil, errors.New("6th argument must be a non-empty string")
	}
	if len(args[6]) <= 0 {
		return nil, errors.New("7th argument must be a non-empty string")
	}
	if len(args[7]) <= 0 {
		return nil, errors.New("8th argument must be a non-empty string")
	}
	if len(args[8]) <= 0 {
 	 return nil, errors.New("9th argument must be a non-empty string")
  }
	list_price, err := strconv.ParseFloat(args[5],64)
	if err != nil {
		return nil, errors.New("5rd argument must be a numeric string")
	}



	user_type := strings.ToLower(args[9])

	str := `{"product_id": "` + args[0] + `", "category": "` + args[1] +
	 `", "product_description": "` + args[2] + `", "availability_start_date": "` + args[3] +
	 `", "availability_end_date": "` + args[4] + `", "list_price": ` + strconv.FormatFloat(list_price, 'f', -1, 64) +
	 `, "currency": "` + args[6] + `", "price_start_date": "` + args[7] +
	 `", "price_end_date": "` + args[8]+ `", "user_type": "` + user_type +
	  `"}`
	err = stub.PutState(args[0], []byte(str))								//store marble with id as key
	if err != nil {
		return nil, err
	}

	//get the product index
	productsAsBytes, err := stub.GetState(productIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get product index")
	}
	var productIndex []string
	json.Unmarshal(productsAsBytes, &productIndex)							//un stringify it aka JSON.parse()

	//check if the product_id exist
	if(!findProduct(productIndex,args[0]) ) {
	//append
	productIndex = append(productIndex, args[0])								//add product name to index list
	fmt.Println("! product index: ", productIndex)
	jsonAsBytes, _ := json.Marshal(productIndex)
	err = stub.PutState(productIndexStr, jsonAsBytes)						//store name of product

	if err != nil {
			fmt.Println("Error creating Product Index");
			return nil, errors.New("Failed to add product index")
		}

		fmt.Println("New Product index added")
	} else {
	fmt.Println("Modified the existing Product")
	}



	fmt.Println("- end init product")
	return nil, nil
}

func findProduct(productsIndex []string, product_id string) (bool) {

	for _,value:= range productsIndex {
		if value == product_id {
			return true;
		}
	}
	return false;
}


// ============================================================================================================================
// Create a new Offering
// ============================================================================================================================
func (t *SimpleChaincode) init_offering(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error

	//   0       1       2     3
	// "asdf", "blue", "35", "bob"
	if len(args) != 11 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	fmt.Println("- start init marble")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return nil, errors.New("5th argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return nil, errors.New("6th argument must be a non-empty string")
	}
	if len(args[6]) <= 0 {
		return nil, errors.New("7th argument must be a non-empty string")
	}
	if len(args[7]) <= 0 {
		return nil, errors.New("8th argument must be a non-empty string")
	}
	if len(args[8]) <= 0 {
 	 return nil, errors.New("9th argument must be a non-empty string")
  }
  	if len(args[9]) <= 0 {
 	 return nil, errors.New("10th argument must be a non-empty string")
  }
  	if len(args[10]) <= 0 {
 	 return nil, errors.New("11th argument must be a non-empty string")
  }
	list_price, err := strconv.ParseFloat(args[5],64)
	if err != nil {
		return nil, errors.New("5rd argument must be a numeric string")
	}


	str := `{"offering_id": "` + args[0] + `", "offering_category": "` + args[1] +
	 `", "offering_description": "` + args[2] + `", "availability_start_date": "` + args[3] +
	 `", "availability_end_date": "` + args[4] + `", "current_list_price": ` + strconv.FormatFloat(list_price, 'f', -1, 64) +
	 `, "currency": "` + args[6] + `", "price_start_date": "` + args[7] +
	 `", "price_end_date": "` + args[8]+ `", "product_id_01": "` + args[9] +`", "product_id_02": "` + args[10] +
	  `"}`
	err = stub.PutState(args[0], []byte(str))								
	if err != nil {
		return nil, err
	}

	//get the offering index
	offeringsAsBytes, err := stub.GetState(offeringIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get offering index")
	}
	var offeringIndex []string
	json.Unmarshal(offeringsAsBytes, &offeringIndex)							

	//check if the offering_id exist
	if(!findOffering(offeringIndex,args[0]) ) {
	//append
	offeringIndex = append(offeringIndex, args[0])								
	fmt.Println("! offering index: ", offeringIndex)
	jsonAsBytes, _ := json.Marshal(offeringIndex)
	err = stub.PutState(offeringIndexStr, jsonAsBytes)					

	if err != nil {
			fmt.Println("Error creating offering Index");
			return nil, errors.New("Failed to add offering index")
		}

		fmt.Println("New offering index added")
	} else {
	fmt.Println("Modified the existing offering")
	}



	fmt.Println("- end init offering")
	return nil, nil
}

func findOffering(offeringsIndex []string, offering_id string) (bool) {

	for _,value:= range offeringsIndex {
		if value == offering_id {
			return true;
		}
	}
	return false;
}
// ============================================================================================================================
// Set User type Permission on Product
// ============================================================================================================================
func (t *SimpleChaincode) set_user_type(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error

	//   0       1
	// "name", "bob"
	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	fmt.Println("- start set user type")
	fmt.Println(args[0] + " - " + args[1])
	productAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get thing")
	}
	res := Product{}
	json.Unmarshal(productAsBytes, &res)										//un stringify it aka JSON.parse()
	res.User_Type = args[1]														//change the user type

	jsonAsBytes, _ := json.Marshal(res)
	err = stub.PutState(args[0], jsonAsBytes)								//rewrite the Product with id as key
	if err != nil {
		return nil, err
	}

	fmt.Println("- end set user type")
	return nil, nil
}



// ============================================================================================================================
// Make Timestamp - create a timestamp in ms
// ============================================================================================================================
func makeTimestamp() int64 {
    return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
}
