/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

// ====CHAINCODE EXECUTION SAMPLES (CLI) ==================

// ==== Invoke coins ====
// peer chaincode invoke -C myc1 -n coins -c '{"Args":["initCoin","coin1","blue","35","tom"]}'
// peer chaincode invoke -C myc1 -n coins -c '{"Args":["initCoin","coin2","red","50","tom"]}'
// peer chaincode invoke -C myc1 -n coins -c '{"Args":["initCoin","coin3","blue","70","tom"]}'
// peer chaincode invoke -C myc1 -n coins -c '{"Args":["transferCoin","coin2","jerry"]}'
// peer chaincode invoke -C myc1 -n coins -c '{"Args":["transferCoinsBasedOnColor","blue","jerry"]}'
// peer chaincode invoke -C myc1 -n coins -c '{"Args":["delete","coin1"]}'

// ==== Query coins ====
// peer chaincode query -C myc1 -n coins -c '{"Args":["readCoin","coin1"]}'
// peer chaincode query -C myc1 -n coins -c '{"Args":["getCoinsByRange","coin1","coin3"]}'
// peer chaincode query -C myc1 -n coins -c '{"Args":["getHistoryForCoin","coin1"]}'

// Rich Query (Only supported if CouchDB is used as state database):
//   peer chaincode query -C myc1 -n coins -c '{"Args":["queryCoinsByOwner","tom"]}'
//   peer chaincode query -C myc1 -n coins -c '{"Args":["queryCoins","{\"selector\":{\"owner\":\"tom\"}}"]}'

// INDEXES TO SUPPORT COUCHDB RICH QUERIES
//
// Indexes in CouchDB are required in order to make JSON queries efficient and are required for
// any JSON query with a sort. As of Hyperledger Fabric 1.1, indexes may be packaged alongside
// chaincode in a META-INF/statedb/couchdb/indexes directory. Each index must be defined in its own
// text file with extension *.json with the index definition formatted in JSON following the
// CouchDB index JSON syntax as documented at:
// http://docs.couchdb.org/en/2.1.1/api/database/find.html#db-index
//
// This coins02 example chaincode demonstrates a packaged
// index which you can find in META-INF/statedb/couchdb/indexes/indexOwner.json.
// For deployment of chaincode to production environments, it is recommended
// to define any indexes alongside chaincode so that the chaincode and supporting indexes
// are deployed automatically as a unit, once the chaincode has been installed on a peer and
// instantiated on a channel. See Hyperledger Fabric documentation for more details.
//
// If you have access to the your peer's CouchDB state database in a development environment,
// you may want to iteratively test various indexes in support of your chaincode queries.  You
// can use the CouchDB Fauxton interface or a command line curl utility to create and update
// indexes. Then once you finalize an index, include the index definition alongside your
// chaincode in the META-INF/statedb/couchdb/indexes directory, for packaging and deployment
// to managed environments.
//
// In the examples below you can find index definitions that support coins02
// chaincode queries, along with the syntax that you can use in development environments
// to create the indexes in the CouchDB Fauxton interface or a curl command line utility.
//

//Example hostname:port configurations to access CouchDB.
//
//To access CouchDB docker container from within another docker container or from vagrant environments:
// http://couchdb:5984/
//
//Inside couchdb docker container
// http://127.0.0.1:5984/

// Index for docType, owner.
// Note that docType and owner fields must be prefixed with the "data" wrapper
//
// Index definition for use with Fauxton interface
// {"index":{"fields":["data.docType","data.owner"]},"ddoc":"indexOwnerDoc", "name":"indexOwner","type":"json"}
//
// Example curl command line to define index in the CouchDB channel_chaincode database
// curl -i -X POST -H "Content-Type: application/json" -d "{\"index\":{\"fields\":[\"data.docType\",\"data.owner\"]},\"name\":\"indexOwner\",\"ddoc\":\"indexOwnerDoc\",\"type\":\"json\"}" http://hostname:port/myc1_coins/_index
//

// Index for docType, owner, size (descending order).
// Note that docType, owner and size fields must be prefixed with the "data" wrapper
//
// Index definition for use with Fauxton interface
// {"index":{"fields":[{"data.size":"desc"},{"data.docType":"desc"},{"data.owner":"desc"}]},"ddoc":"indexSizeSortDoc", "name":"indexSizeSortDesc","type":"json"}
//
// Example curl command line to define index in the CouchDB channel_chaincode database
// curl -i -X POST -H "Content-Type: application/json" -d "{\"index\":{\"fields\":[{\"data.size\":\"desc\"},{\"data.docType\":\"desc\"},{\"data.owner\":\"desc\"}]},\"ddoc\":\"indexSizeSortDoc\", \"name\":\"indexSizeSortDesc\",\"type\":\"json\"}" http://hostname:port/myc1_coins/_index

// Rich Query with index design doc and index name specified (Only supported if CouchDB is used as state database):
//   peer chaincode query -C myc1 -n coins -c '{"Args":["queryCoins","{\"selector\":{\"docType\":\"coin\",\"owner\":\"tom\"}, \"use_index\":[\"_design/indexOwnerDoc\", \"indexOwner\"]}"]}'

// Rich Query with index design doc specified only (Only supported if CouchDB is used as state database):
//   peer chaincode query -C myc1 -n coins -c '{"Args":["queryCoins","{\"selector\":{\"docType\":{\"$eq\":\"coin\"},\"owner\":{\"$eq\":\"tom\"},\"size\":{\"$gt\":0}},\"fields\":[\"docType\",\"owner\",\"size\"],\"sort\":[{\"size\":\"desc\"}],\"use_index\":\"_design/indexSizeSortDoc\"}"]}'

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type coin struct {
	Name	string `json:"Name"`
	Amount       string `json:"amount"`    //the fieldtags are needed to keep case from bouncing around
	Owner      string `json:"owner"`
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "initCoin" { //create a new coin
		return t.initCoin(stub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "transferCoin" { //change owner of a specific coin
		return t.transferCoin(stub, args)
	} else if function == "delete" { //delete a coin
		return t.delete(stub, args)
	} else if function == "readCoin" { //read a coin
		return t.readCoin(stub, args)
	} else if function == "queryCoinsByOwner" { //find coins for owner X using rich query
		return t.queryCoinsByOwner(stub, args)
	} else if function == "queryCoins" { //find coins based on an ad hoc rich query
		return t.queryCoins(stub, args)
	} else if function == "getHistoryForCoin" { //get history of values for a coin
		return t.getHistoryForCoin(stub, args)
	} else if function == "getCoinsByRange" { //get coins based on range query
		return t.getCoinsByRange(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// initCoin - create a new coin, store into chaincode state
// ============================================================
func (t *SimpleChaincode) initCoin(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	//   0       1       2  
	// "coin1",  "aCent",   "bob"
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// ==== Input sanitation ====
	fmt.Println("- start init coin")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	coinName := args[0]
	owner := strings.ToLower(args[2])
	amount := strings.ToLower(args[1])


	// ==== Check if coin already exists ====
	coinAsBytes, err := stub.GetState(coinName)
	if err != nil {
		return shim.Error("Failed to get coin: " + err.Error())
	} else if coinAsBytes != nil {
		fmt.Println("This coin already exists: " + coinName)
		return shim.Error("This coin already exists: " + coinName)
	}

	// ==== Create coin object and marshal to JSON ====
	//objectType := "coin"
	coin := &coin{coinName, amount , owner}
	coinJSONasBytes, err := json.Marshal(coin)
	if err != nil {
		return shim.Error(err.Error())
	}
	//Alternatively, build the coin json string manually if you don't want to use struct marshalling
	//coinJSONasString := `{"docType":"Coin",  "name": "` + coinName + `", "amount": ` + strconv.Itoa(amount) + `, "owner": "` + owner + `"}`
	//coinJSONasBytes := []byte(str)

	// === Save coin to state ===
	err = stub.PutState(coinName, coinJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//  ==== Index the coin to enable color-based range queries, e.g. return all blue coins ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on indexName~color~name.
	//  This will enable very efficient state range queries based on composite keys matching indexName~color~*
	indexName := "amount~name"
	amountNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{coin.Amount, coin.Name})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the coin.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(amountNameIndexKey, value)

	// ==== Coin saved and i