/*
Copyright IBM Corp. 2016 All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
                 http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
        "bytes"
        "fmt"
        "encoding/json"
        "strconv"
        "time"

        "github.com/hyperledger/fabric/core/chaincode/shim"
        pb "github.com/hyperledger/fabric/protos/peer"
)

// TsohueChainCode - Chaincode type
type TsohueChainCode struct {
}



// Groupbuy - Groupbuy product. ProductID=NULL is total amount for the product
type Groupbuy struct { // Key: ClientID, ProductID
        ObjectType string `json : objectType`
        Timestamp string `json : timestamp`
	    GroupbuyID string `json :groupbuyID`
        ProductID string `json : productID`
        Currency string `json : currency`
	    Target_amount float64 `json : target_amount`
        Share float64 `json : share`
	    Status string `json : status`
	    Dividend string `json : dividend`
        Expiry_date string `json : expiry date`
	    Capital string `json : capital`
}


// Groupbuy_record - Groupbuy_record  infoemation
type Groupbuy_record struct { // Key: ClientID, ProductID
        ObjectType string `json : objectType`
        Timestamp string `json : timestamp`
        GroupbuyID string `json :groupbuyID`
        TransactionID  string `json : transactionID`
        Currency string `json : currency`
        Target_amount float64 `json : target_amount`
        Share float64 `json : share`
        Accrued_interest float64 `json : accrued_interest`
        Status string `json : status`
}


// User - User infoemation
type User struct { // Key: ClientID, ProductID
        ObjectType string `json : objectType`
	Timestamp string `json : timestamp`
        ClientID string `json : clientID`
        UserID string `json : userID`
        Bank string `json : bank`
        Bank_account string `json : bank_account`
        Address string `json : address`
        Phone string `json : phone`
        Status_open string `json : status_open`
        Status_bind string `json : status_bind`
}


// User - User infoemation
type User_record struct { // Key: ClientID, ProductID
        ObjectType string `json : objectType`
        Timestamp string `json : timestamp`
        ClientID string `json : clientID`
        UserID string `json : userID`
        Bank string `json : bank`
        Bank_account string `json : bank_account`
        Address string `json : address`
        Phone string `json : phone`
        Status_open string `json : status_open`
        Status_bind string `json : status_bind`
}


// Transaction - Transaction infoemation
type Transaction_contract struct { // Key: ClientID, ProductID
        ObjectType string `json : objectType`
        Timestamp string `json : timestamp`
	    TransactionID string `json : transactionID`
        ClientID string `json : clientID`
	    GroupbuyID string `json : groupbuyID`
	    Currency string `json : currency`
        Amount float64 `json : amount`
	    Status string `json : status`
}


// Transaction - Transaction infoemation
type Transaction_contract_record struct { // Key: ClientID, ProductID
        ObjectType string `json : objectType`
        Timestamp string `json : timestamp`
        TransactionID string `json : transactionID`
        ClientID string `json : clientID`
        GroupbuyID string `json : groupbuyID`
        Currency string `json : currency`
        Amount float64 `json : amount`
        Status string `json : status`
}


// Transaction - Transaction infoemation
type Transaction_blocked struct { // Key: ClientID, ProductID
        ObjectType string `json : objectType`
        Timestamp string `json : timestamp`
        TransactionID string `json : transactionID`
        ClientID string `json : clientID`
        GroupbuyID string `json : groupbuyID`
        Currency string `json : currency`
        Amount float64 `json : amount`
        Status string `json : status`
}


// Transaction - Transaction infoemation
type Transaction_blocked_record struct { // Key: ClientID, ProductID
        ObjectType string `json : objectType`
        Timestamp string `json : timestamp`
        TransactionID string `json : transactionID`
        ClientID string `json : clientID`
        GroupbuyID string `json : groupbuyID`
        Currency string `json : currency`
        Amount float64 `json : amount`
        Status string `json : status`
}


// Transaction - Transaction infoemation
type Transaction_action struct { // Key: ClientID, ProductID
        ObjectType string `json : objectType`
        Timestamp string `json : timestamp`
        TransactionID string `json : transactionID`
        ClientID string `json : clientID`
        GroupbuyID string `json : groupbuyID`
	    Status string `json : status`
}


// Transaction - Transaction infoemation
type Transaction_action_record struct { // Key: ClientID, ProductID
        ObjectType string `json : objectType`
        Timestamp string `json : timestamp`
        TransactionID string `json : transactionID`
        ClientID string `json : clientID`
        GroupbuyID string `json : groupbuyID`
        Status string `json : status`
}


// Transaction - Transaction infoemation
type Transaction_takeover struct { // Key: ClientID, ProductID
        ObjectType string `json : objectType`
        Timestamp string `json : timestamp`
        TransactionID string `json : transactionID`
        GroupbuyID string `json : groupbuyID`
        ClientID_sell string `json : clientID_sell`
	    ClientID_buy string `json : clientID_buy`
	    Status string `json : status`
}


// Transaction - Transaction infoemation
type Transaction_takeover_record struct { // Key: ClientID, ProductID
        ObjectType string `json : objectType`
        Timestamp string `json : timestamp`
        TransactionID string `json : transactionID`
        GroupbuyID string `json : groupbuyID`
        ClientID_sell string `json : clientID_sell`
        ClientID_buy string `json : clientID_buy`
        Status string `json : status`
}


// Use nested structs to parse and process JSON inputs.
// JSONInput - struct for JSON string inputs.
type JSONInput struct {
        ProductID string `json : productID`
        Currency string `json : currency`
        ClientInfo []ClientInfo `json : clientInfo` // Nested struct
}

// ClientInfo - struct type passed in through JSON. Used in deduct, dividend functions to update all clients within one groupbuy
type ClientInfo struct {
        ClientID string `json : clientID`
        Amount string `json : amount` // The amount passed in to change to Inventory
}


func main() {
        err := shim.Start(new(TsohueChainCode))
        if err != nil {
                fmt.Printf("Error starting cash chaincode: %s", err)
        }
}

// Init - initializes chaincode
func (t *TsohueChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
        fmt.Println("Please enter invoke functions.")
        return shim.Success(nil)
}

// Invoke - invoke functions to update/query chaincode
func (t *TsohueChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
        function, args := stub.GetFunctionAndParameters()
        if function == "query" || function == "invoke" {
                function = args[0]
                args = args[1:]
        }
	fmt.Println("Chaincode Invoked " + function)
               if function == "groupbuyRaise" {
                return t.groupbuyRaise(stub, args)
        } else if function == "userRaise" {
                return t.userRaise(stub, args)
        } else if function == "groupbuyJoin" {
                return t.groupbuyJoin(stub, args)
	} else if function == "groupbuyJoin_failed" {
                return t.groupbuyJoin_failed(stub, args)
        } else if function == "groupbuyLeave" {
                return t.groupbuyLeave(stub, args)
	} else if function == "groupbuyLeave_failed" {
                return t.groupbuyLeave_failed(stub, args)
        } else if function == "groupbuyBlocked" {
                return t.groupbuyBlocked(stub, args)
	} else if function == "groupbuyBlocked_failed" {
                return t.groupbuyBlocked_failed(stub, args)
        } else if function == "groupbuyUnblocked" {
                return t.groupbuyUnblocked(stub, args)
        } else if function == "groupbuyUnblocked_failed"{
                return t.groupbuyUnblocked_failed(stub, args)
        } else if function == "groupbuySuccess" {
                return t.groupbuySuccess(stub, args)
	} else if function == "groupbuyNoGoodPrice" {
                return t.groupbuyNoGoodPrice(stub, args)
	} else if function == "groupbuyInvesting" {
                return t.groupbuyInvesting(stub, args)
        } else if function == "groupbuyMataned" {
                return t.groupbuyMataned(stub, args)
	} else if function == "groupbuyLiquidate" {
                return t.groupbuyLiquidate(stub, args)
	} else if function == "groupbuyEnded" {
                return t.groupbuyEnded(stub, args)
	} else if function == "groupbuyFailed"{
                return t.groupbuyFailed(stub, args)
        } else if function == "dividend" {
                return t.dividend(stub, args)
        } else if function == "capital" {
                return t.capital(stub, args)
	} else if function == "liquidated" {
                return t.liquidated(stub, args)
	} else if function == "contract_failed" {
                return t.contract_failed(stub, args)
	} else if function == "takeOver_raise" {
                return t.takeOver_raise(stub, args)
        } else if function == "takeOver_failed" {
                return t.takeOver_failed(stub, args)
	} else if function == "takeOver" {
                return t.takeOver(stub, args)



	} else if function == "queryTransaction_blockedByTransactionID" {
                return t.queryTransaction_blockedByTransactionID(stub, args)
        } else if function == "queryTransaction_blockedByClientID" {
                return t.queryTransaction_blockedByClientID(stub, args)
        } else if function == "queryTransaction_blocked_recordByTransactionID" {
                return t.queryTransaction_blocked_recordByTransactionID(stub, args)
	} else if function == "queryTransaction_contractByGroupbuyID" {
                return t.queryTransaction_contractByGroupbuyID(stub, args)
        } else if function == "queryTransaction_contractByClientID" {
                return t.queryTransaction_contractByClientID(stub, args)
        } else if function == "queryTransaction_contractByTransactionID" {
                return t.queryTransaction_contractByTransactionID(stub, args)
        } else if function == "queryTransaction_contract_recordByTransactionID" {
                return t.queryTransaction_contract_recordByTransactionID(stub, args)
	} else if function == "queryTransaction_action_joinByClientID" {
                return t.queryTransaction_action_joinByClientID(stub, args)
        } else if function == "queryTransaction_action_leaveByClientID" {
                return t.queryTransaction_action_leaveByClientID(stub, args)
        } else if function == "queryTransaction_action_joinByGroupbuyID" {
                return t.queryTransaction_action_joinByGroupbuyID(stub, args)
        } else if function == "queryTransaction_action_leaveByGroupbuyID" {
                return t.queryTransaction_action_leaveByGroupbuyID(stub, args)
	} else if function == "queryTransaction_actionByTransactionID" {
                return t.queryTransaction_actionByTransactionID(stub, args)
        } else if function == "queryTransaction_action_recordByTransactionID" {
                return t.queryTransaction_action_recordByTransactionID(stub, args)
	} else if function == "queryTransaction_takeoverByClientID_sell" {
                return t.queryTransaction_takeoverByClientID_sell(stub, args)
        } else if function == "queryTransaction_takeoverByClientID_buy" {
                return t.queryTransaction_takeoverByClientID_buy(stub, args)
        } else if function == "queryTransaction_takeoverByGroupbuyID" {
                return t.queryTransaction_takeoverByGroupbuyID(stub, args)
        } else if function == "queryTransaction_takeoverByTransactionID" {
                return t.queryTransaction_takeoverByTransactionID(stub, args)
        } else if function == "queryTransaction_takeover_recordByTransactionID" {
                return t.queryTransaction_takeover_recordByTransactionID(stub, args)
	} else if function == "queryGroupbuyByGroupbuyID" {
                return t.queryGroupbuyByGroupbuyID(stub, args)
        } else if function == "queryGroupbuy_recordByGroupbuyID" {
                return t.queryGroupbuy_recordByGroupbuyID(stub, args)
        } else if function == "queryGroupbuy_recordByTransactionID" {
                return t.queryGroupbuy_recordByTransactionID(stub, args)
        } else if function == "queryUserByClientID" {
                return t.queryUserByClientID(stub, args)
        } else if function == "queryAssetHistory" {
                return t.queryAssetHistory(stub, args)
        } else if function == "queryGroupbuyHistory" {
                return t.queryGroupbuyHistory(stub, args) 
	

	
	} else if function == "queryUser" {
                return t.queryUser(stub, args)
        } else if function == "queryUser_record" {
                return t.queryUser_record(stub, args)
	} else if function == "queryGroupbuy" {
                return t.queryGroupbuy(stub, args)
	} else if function == "queryGroupbuy_record" {
                return t.queryGroupbuy_record(stub, args)
        } else if function == "queryTransaction_blocked" {
                return t.queryTransaction_blocked(stub, args)
        } else if function == "queryTransaction_blocked_record" {
                return t.queryTransaction_blocked_record(stub, args)
	} else if function == "queryTransaction_contract" {
                return t.queryTransaction_contract(stub, args)
	} else if function == "queryTransaction_contract_record" {
                return t.queryTransaction_contract_record(stub, args)
        } else if function == "queryTransaction_action" {
                return t.queryTransaction_action(stub, args)
        } else if function == "queryTransaction_action_record" {
                return t.queryTransaction_action_record(stub, args)
        } else if function == "queryTransaction_takeover" {
                return t.queryTransaction_takeover(stub, args)
        } else if function == "queryTransaction_takeover_record" {
                return t.queryTransaction_takeover_record(stub, args)
        } else if function == "queryByString" {
                return t.queryByString(stub, args)
        }
        fmt.Println("Error! Invoke did not find function: " + function)
        return shim.Error("Received unknown function invocation")
}



//// Groupbuy Related /////
// =======================================================
// groupbuyRaise
// Raising(init) a groupbuy product,
// 1. Adds product to groupBuy (status=open),
// 2. Adds record to groupBuy_record,
// input: transactionID, groupbuyID, productID, currency, target_amount, dividend, expiry_date, capital
// =======================================================
func (t *TsohueChainCode) groupbuyRaise(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error
        if len(args) != 8 {
                return shim.Error("Incorrect number of arguments. Expecting 8")
        }
	transactionID := args[0]
	groupbuyID := args[1]
        productID := args[2]
        currency := args[3]
	target_amount := args[4] // will be checked in called functions
	share := "0"
        status := "open"
	dividend := args[5]
	expiry_date := args[6]
	capital := args[7]
        //update groupbuy 
        groupbuyRaiseArgs := []string{groupbuyID, productID, currency, target_amount, share, status, dividend, expiry_date, capital}
        err = updateGroupbuyCalled(stub, groupbuyRaiseArgs)
        if err != nil {
                return shim.Error(err.Error())
        }
	//update groupbuy_record 
	groupbuy_recordArgs := []string{transactionID, groupbuyID, currency, target_amount, share, status}
        err = updateGroupbuy_recordCalled(stub, groupbuy_recordArgs)
        if err != nil {
                return shim.Error(err.Error())
        }

        return shim.Success(nil)
}
// =======================================================
// userRaise
// Client data raise,
// 1.Adds client to User,
// input: clientID, userID, bank, bank_account, address, phone, status_open, status_bind
// =======================================================
func (t *TsohueChainCode) userRaise(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var clientID, userID, bank, bank_account, address, phone, status_open, status_bind string
	var err error
        if len(args) != 8 {
                return shim.Error("Incorrect number of arguments. Expecting 8")
        }
        clientID = args[0]
        userID = args[1]
        bank = args[2]
        bank_account = args[3]
        address = args[4]
        phone = args[5]
        status_open = args[6]
        status_bind = args[7]
        userArgs := []string{clientID, userID, bank, bank_account, address, phone, status_open, status_bind}
        err = updateUserCalled(stub, userArgs)
        if err != nil {
                return shim.Error(err.Error())
        }


	user_recordArgs := []string{clientID, userID, bank, bank_account, address, phone, status_open, status_bind}
        err = updateUser_recordCalled(stub, user_recordArgs)
        if err != nil {
                return shim.Error(err.Error())
        }



        return shim.Success(nil)
}
// =======================================================
// groupbuyJoin
// Client joining a groupbuy product,
// 1. Update record to groupbuy.share (groupbuy.share=groupbuy.share+share)
// 2. Adds record to transcation_action (status=join),
// 3. Adds record to transcation_action_record ,
// 4. Adds record to transcation_blocked (status=blocked),
// 5. Adds record to transcation_blocked_record ,
// input: TransactionID, ClientID, ProductID, currency, share
// =======================================================
func (t * TsohueChainCode) groupbuyJoin(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error
	var share float64

        if len(args) != 5 {
                return shim.Error("Incorrect number of arguments. Expecting 5")
        }
        transactionID := args[0]
        clientID := args[1]
	groupbuyID := args[2]
	currency := args[3]
	share, err = strconv.ParseFloat(args[4], 64)
	status_action := "join"
	status_blocked := "blocked"


        //Update groupbuy
        //query groupbuy 
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"groupbuy\",\"GroupbuyID\":\"%s\"}}", groupbuyID)
        // Get Groupbuy results from queryString
        resultsIterator, err := stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
        count := 0
        // Loop through all records of Groupbuy,
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                groupbuyTemp := &Groupbuy{}
                groupbuyKey := responseRange.Key

                // Unmarshal responseRange.Value to JSON
                err = json.Unmarshal(responseRange.Value, &groupbuyTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemp
                groupbuyTemp.Share = groupbuyTemp.Share + share

                // Marshal and Put groupbuy to state
                groupbuyJSONasBytes, _ := json.Marshal(groupbuyTemp)
                err = stub.PutState(groupbuyKey, groupbuyJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())
                }
        count ++
        }



        // Add to Transaction_action
        transaction_actionArgs := []string{transactionID, clientID, groupbuyID, status_action}
        err = updateTransaction_actionCalled(stub, transaction_actionArgs)
        if err != nil {
                return shim.Error(err.Error())
        }



	// Add to Transaction_action_record
        transaction_action_recordArgs := []string{transactionID, clientID, groupbuyID, status_action}
        err = updateTransaction_action_recordCalled(stub, transaction_action_recordArgs)
        if err != nil {
                return shim.Error(err.Error())
        }



        // Add to Transaction_blocked
        transaction_blockedArgs := []string{transactionID, clientID, groupbuyID, currency, fmt.Sprintf("%f",share), status_blocked}
        err = updateTransaction_blockedCalled(stub, transaction_blockedArgs)
        if err != nil {
                return shim.Error(err.Error())
        }



        // Add to Transaction_blocked_record
        transaction_blocked_recordArgs := []string{transactionID, clientID, groupbuyID, currency, fmt.Sprintf("%f",share), status_blocked}
        err = updateTransaction_blocked_recordCalled(stub, transaction_blocked_recordArgs)
        if err != nil {
                return shim.Error(err.Error())
        }



        fmt.Println("Groupbuy joined success!")
        return shim.Success(nil)
}
// =======================================================
// groupbuyJoin_failed
// Client Join a groupbuy product failed,
// 1. Update record to transaction_action (status=fail),
// 2. Adds record to transaction_action_record ,
// 3. Update record to Transaction_blocked (status=fail),
// 4. Adds record to Transaction_blocked_record ,
// 5. Update record to groupbuy.share (groupbuy.share=groupbuy.share-share)
// input: TransactionID, ClientID, ProductID
// =======================================================
func (t * TsohueChainCode) groupbuyJoin_failed(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error
	var share float64
	var groupbuyID string
        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }
        transactionID := args[0]



        //Update Transaction_action
        //query Transaction_action
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_action\",\"TransactionID\":\"%s\"}}", transactionID)
        // Get Groupbuy results from queryString
        resultsIterator, err := stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
        // Loop through all records of Groupbuy,
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                transaction_actionTemp := &Transaction_action{}
                transaction_actionKey := responseRange.Key

                // Unmarshal responseRange
                err = json.Unmarshal(responseRange.Value, &transaction_actionTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemp
                transaction_actionTemp.Status = "fail"
                transaction_actionTemp.Timestamp = time.Now().String()

                // Marshal and Put groupbuy to state
                transaction_actionJSONasBytes, _ := json.Marshal(transaction_actionTemp)
                err = stub.PutState(transaction_actionKey, transaction_actionJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())

                }
                // Add to Transaction_action_record
                transaction_action_recordArgs := []string{transactionID, transaction_actionTemp.ClientID, transaction_actionTemp.GroupbuyID, transaction_actionTemp.Status}
                err = updateTransaction_action_recordCalled(stub, transaction_action_recordArgs)
                if err != nil {
                        return shim.Error(err.Error())
                }


        }



        //Update Transaction_blocked
        //query Transaction_blocked
        queryString = fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_blocked\",\"TransactionID\":\"%s\"}}", transactionID)
        // Get Groupbuy results from queryString
        resultsIterator, err = stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
        // Loop through all records of Groupbuy,
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                transaction_blockedTemp := &Transaction_blocked{}
                transaction_blockedKey := responseRange.Key

                // Unmarshal responseRange
                err = json.Unmarshal(responseRange.Value, &transaction_blockedTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemp
                transaction_blockedTemp.Status = "fail"
                transaction_blockedTemp.Timestamp = time.Now().String()
		share = transaction_blockedTemp.Amount
		groupbuyID = transaction_blockedTemp.GroupbuyID

                // Marshal and Put groupbuy to state
                transaction_blockedJSONasBytes, _ := json.Marshal(transaction_blockedTemp)
                err = stub.PutState(transaction_blockedKey, transaction_blockedJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())

                }
	        // Add to Transaction_blocked_record
                transaction_blocked_recordArgs := []string{transactionID, transaction_blockedTemp.ClientID, transaction_blockedTemp.GroupbuyID, transaction_blockedTemp.Currency, fmt.Sprintf("%f",transaction_blockedTemp.Amount), transaction_blockedTemp.Status}
                err = updateTransaction_blocked_recordCalled(stub, transaction_blocked_recordArgs)
                if err != nil {
                        return shim.Error(err.Error())
                }


        }



        //Update groupbuy
        //query groupbuy
        queryString = fmt.Sprintf("{\"selector\":{\"ObjectType\":\"groupbuy\",\"GroupbuyID\":\"%s\"}}", groupbuyID)
        // Get Groupbuy results from queryString
        resultsIterator, err = stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
        // Loop through all records of Groupbuy,
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                groupbuyTemp := &Groupbuy{}
                groupbuyKey := responseRange.Key

                // Unmarshal responseRange.Value to JSON
                err = json.Unmarshal(responseRange.Value, &groupbuyTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemp
                groupbuyTemp.Share = groupbuyTemp.Share - share

                // Marshal and Put groupbuy to state
                groupbuyJSONasBytes, _ := json.Marshal(groupbuyTemp)
                err = stub.PutState(groupbuyKey, groupbuyJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())
                }

        }



        fmt.Println("Groupbuy leaved success!")
        return shim.Success(nil)
}
// =======================================================
// groupbuyLeave
// Client Leaving a groupbuy product,
// 1. Adds record to transaction_action (status=leave),
// 2. Adds record to transaction_action_record ,
// 3. Update record to Transaction_blocked (status=fail),
// 4. Adds record to Transaction_blocked_record ,
// 5. Update record to groupbuy.share (groupbuy.share=groupbuy.share-share)
// input: TransactionID, ClientID, ProductID
// =======================================================
func (t * TsohueChainCode) groupbuyLeave(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error
	var share float64
	var groupbuyID string
        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }
        transactionID := args[0]



        //Update Transaction_action
        //query Transaction_action
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_action\",\"TransactionID\":\"%s\"}}", transactionID)
        // Get Groupbuy results from queryString
        resultsIterator, err := stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
        // Loop through all records of Groupbuy,
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                transaction_actionTemp := &Transaction_action{}
                transaction_actionKey := responseRange.Key

                // Unmarshal responseRange
                err = json.Unmarshal(responseRange.Value, &transaction_actionTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemp
                transaction_actionTemp.Status = "leave"
                transaction_actionTemp.Timestamp = time.Now().String()

                // Marshal and Put groupbuy to state
                transaction_actionJSONasBytes, _ := json.Marshal(transaction_actionTemp)
                err = stub.PutState(transaction_actionKey, transaction_actionJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())

                }
		// Add to Transaction_action_record
                transaction_action_recordArgs := []string{transactionID, transaction_actionTemp.ClientID, transaction_actionTemp.GroupbuyID, transaction_actionTemp.Status}
                err = updateTransaction_action_recordCalled(stub, transaction_action_recordArgs)
                if err != nil {
                        return shim.Error(err.Error())
                }

        }



        // Add to Transaction_contract
        //query Transaction_contract
        queryString = fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_blocked\",\"TransactionID\":\"%s\"}}", transactionID)
        // Get Groupbuy results from queryString
        resultsIterator, err = stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
        // Loop through all records of Groupbuy,
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                transaction_blockedTemp := &Transaction_blocked{}
                transaction_blockedKey := responseRange.Key

                // Unmarshal responseRange
                err = json.Unmarshal(responseRange.Value, &transaction_blockedTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemp
                transaction_blockedTemp.Status = "fail"
                transaction_blockedTemp.Timestamp = time.Now().String()
		share = transaction_blockedTemp.Amount
		groupbuyID = transaction_blockedTemp.GroupbuyID

                // Marshal and Put groupbuy to state
                transaction_blockedJSONasBytes, _ := json.Marshal(transaction_blockedTemp)
                err = stub.PutState(transaction_blockedKey, transaction_blockedJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())

                }

                // Add to Transaction_blocked_record
                transaction_blocked_recordArgs := []string{transactionID, transaction_blockedTemp.ClientID, transaction_blockedTemp.GroupbuyID, transaction_blockedTemp.Currency, fmt.Sprintf("%f",transaction_blockedTemp.Amount), transaction_blockedTemp.Status}
                err = updateTransaction_blocked_recordCalled(stub, transaction_blocked_recordArgs)
                if err != nil {
                        return shim.Error(err.Error())
                }


        }



        //Update groupbuy
        //query groupbuy
        queryString = fmt.Sprintf("{\"selector\":{\"ObjectType\":\"groupbuy\",\"GroupbuyID\":\"%s\"}}", groupbuyID)
        // Get Groupbuy results from queryString
        resultsIterator, err = stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
        // Loop through all records of Groupbuy,
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                groupbuyTemp := &Groupbuy{}
                groupbuyKey := responseRange.Key

                // Unmarshal responseRange.Value to JSON
                err = json.Unmarshal(responseRange.Value, &groupbuyTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemp
                groupbuyTemp.Share = groupbuyTemp.Share - share

                // Marshal and Put groupbuy to state
                groupbuyJSONasBytes, _ := json.Marshal(groupbuyTemp)
                err = stub.PutState(groupbuyKey, groupbuyJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())
                }

        }



        fmt.Println("Groupbuy leaved success!")
        return shim.Success(nil)
}
// =======================================================
// groupbuyLeave_failed
// Client Leaved a groupbuy product fail,
// 1. Update record to transaction_action (status=join),
// 2. Adds record to transaction_action_record ,
// 3. Update record to Transaction_blocked (status=blocked),
// 4. Adds record to Transaction_blocked_record ,
// 5. Update record to groupbuy.share (groupbuy.share=groupbuy.share+share)
// input: TransactionID, ClientID, ProductID
// =======================================================
func (t * TsohueChainCode) groupbuyLeave_failed(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error
	var share float64
	var groupbuyID string
        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }
        transactionID := args[0]



        // Update Transaction_action
        //query Transaction_action
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_action\",\"ClientID\":\"%s\",\"GroupbuyID\":\"%s\"}}", transactionID)
        // Get Groupbuy results from queryString
        resultsIterator, err := stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
        // Loop through all records of Groupbuy,
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                transaction_actionTemp := &Transaction_action{}
                transaction_actionKey := responseRange.Key

                // Unmarshal responseRange
                err = json.Unmarshal(responseRange.Value, &transaction_actionTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update transaction_actionTemp
                transaction_actionTemp.Status = "join"
                transaction_actionTemp.Timestamp = time.Now().String()

                // Marshal and Put groupbuy to state
                transaction_actionJSONasBytes, _ := json.Marshal(transaction_actionTemp)
                err = stub.PutState(transaction_actionKey, transaction_actionJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())

                }
		// Add to Transaction_action_record
                transaction_action_recordArgs := []string{transactionID, transaction_actionTemp.ClientID, transaction_actionTemp.GroupbuyID, transaction_actionTemp.Status}
                err = updateTransaction_action_recordCalled(stub, transaction_action_recordArgs)
                if err != nil {
                        return shim.Error(err.Error())
                }

        }



        // Update Transaction_contract
        //query Transaction_contract
        queryString = fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_blocked\",\"TransactionID\":\"%s\"}}", transactionID)
        // Get Groupbuy results from queryString
        resultsIterator, err = stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
        // Loop through all records of Groupbuy,
        for resultsIterator.HasNext() {
		responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                transaction_blockedTemp := &Transaction_blocked{}
                transaction_blockedKey := responseRange.Key

                // Unmarshal responseRange
                err = json.Unmarshal(responseRange.Value, &transaction_blockedTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemp
                transaction_blockedTemp.Status = "blocked"
                transaction_blockedTemp.Timestamp = time.Now().String()
		share = transaction_blockedTemp.Amount
		groupbuyID = transaction_blockedTemp.GroupbuyID

                // Marshal and Put groupbuy to state
                transaction_blockedJSONasBytes, _ := json.Marshal(transaction_blockedTemp)
                err = stub.PutState(transaction_blockedKey, transaction_blockedJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())

                }
		// Add to Transaction_blocked_record
                transaction_blocked_recordArgs := []string{transactionID, transaction_blockedTemp.ClientID, transaction_blockedTemp.GroupbuyID, transaction_blockedTemp.Currency, fmt.Sprintf("%f",transaction_blockedTemp.Amount), transaction_blockedTemp.Status}
                err = updateTransaction_blocked_recordCalled(stub, transaction_blocked_recordArgs)
                if err != nil {
                        return shim.Error(err.Error())
                }

        }



        //Update groupbuy
        //query groupbuy
        queryString = fmt.Sprintf("{\"selector\":{\"ObjectType\":\"groupbuy\",\"GroupbuyID\":\"%s\"}}", groupbuyID)
        // Get Groupbuy results from queryString
        resultsIterator, err = stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
        // Loop through all records of Groupbuy,
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                groupbuyTemp := &Groupbuy{}
                groupbuyKey := responseRange.Key

                // Unmarshal responseRange.Value to JSON
                err = json.Unmarshal(responseRange.Value, &groupbuyTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemp
                groupbuyTemp.Share = groupbuyTemp.Share + share

                // Marshal and Put groupbuy to state
                groupbuyJSONasBytes, _ := json.Marshal(groupbuyTemp)
                err = stub.PutState(groupbuyKey, groupbuyJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())
                }

        }



        fmt.Println("Groupbuy leaved success!")
        return shim.Success(nil)
}
// =======================================================
// groupbuyBlock
// Blocked this transactionID
// 1. Update Transaction_blocked (status = blocked),
// 2. Add to Transaction_blocked_record,
// input: transactionID, clientID, groupbuyID, currency, amount,
// =======================================================
func (t *TsohueChainCode) groupbuyBlocked(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 3")
        }
        fmt.Println("- start groupbuyBlock")
	transactionID := args[0]
	clientID := args[1]
	groupbuyID := args[2]
	currency := args[3]
	amount := args[4]
        status := "blocked"



        //update Transaction_blocked
        transaction_blocked_Args := []string{transactionID, clientID, groupbuyID, currency, amount, status}
        err = updateTransaction_takeoverCalled(stub, transaction_blocked_Args)
        if err != nil {
        return shim.Error(err.Error())
        }



        // Add to Transaction_blocked_record
        transaction_blocked_recordArgs := []string{transactionID, clientID, groupbuyID, currency, fmt.Sprintf("%f", amount), status}
        err = updateTransaction_blocked_recordCalled(stub, transaction_blocked_recordArgs)
        if err != nil {
                return shim.Error(err.Error())
        }



        fmt.Println("Groupbuy blocked success!")
        return shim.Success(nil)

}
// =======================================================
// groupbuyBlocked_failed
// Blocked this transactionID failed,
// 1. Update Transaction_blocked (status=fail),
// 2. Adds record to transaction_blocked_record ,
// input: transactionID,
// =======================================================
func (t * TsohueChainCode) groupbuyBlocked_failed(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error
        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }
        transactionID := args[0]



        // Update Transaction_blocked
        //query Transaction_blocked
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_blocked\",\"TransactionID\":\"%s\"}}", transactionID)
        // Get Groupbuy results from queryString
        resultsIterator, err := stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
        // Loop through all records of Groupbuy,
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                transaction_blockedTemp := &Transaction_blocked{}
                transaction_blockedKey := responseRange.Key

                // Unmarshal responseRange.Value to JSON
                err = json.Unmarshal(responseRange.Value, &transaction_blockedTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemup
                transaction_blockedTemp.Status = "blocked_fail"
                transaction_blockedTemp.Timestamp = time.Now().String()

                // Marshal and Put groupbuy to state
                transaction_blockedJSONasBytes, _ := json.Marshal(transaction_blockedTemp)
                err = stub.PutState(transaction_blockedKey, transaction_blockedJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())

                }
		// Add to Transaction_blocked_record
                transaction_blocked_recordArgs := []string{transactionID, transaction_blockedTemp.ClientID, transaction_blockedTemp.GroupbuyID, transaction_blockedTemp.Currency, fmt.Sprintf("%f",transaction_blockedTemp.Amount), transaction_blockedTemp.Status}
                err = updateTransaction_blocked_recordCalled(stub, transaction_blocked_recordArgs)
                if err != nil {
                        return shim.Error(err.Error())
                }

        }



        fmt.Println("Groupbuy leaved success!")
        return shim.Success(nil)
}
// =======================================================
// groupbuyUnblocked
// Unblocked this transactionID,
// 1. Update Transaction_blocked (status=Unblocked),
// 2. Adds record to transaction_blocked_record ,
// 3. Adds record toTransaction_contract (status="join_amount"),
// 4. Adds record toTransaction_contract_record,
// input: transactionID,
// =======================================================
func (t * TsohueChainCode) groupbuyUnblocked(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error
        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }
	transactionID := args[0]
	status_contract := "join_amount"



	//Update Transaction_blocked
        //query Transaction_blocked
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_blocked\",\"transactionID\":\"%s\"}}", transactionID)
        // Get Groupbuy results from queryString
        resultsIterator, err := stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
        // Loop through all records of Groupbuy,
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                transaction_blockedTemp := &Transaction_blocked{}
                transaction_blockedKey := responseRange.Key

                // Unmarshal responseRange.Value to JSON
                err = json.Unmarshal(responseRange.Value, &transaction_blockedTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemup
                transaction_blockedTemp.Status = "Unblocked"
                transaction_blockedTemp.Timestamp = time.Now().String()

                // Marshal and Put groupbuy to state
                transaction_blockedJSONasBytes, _ := json.Marshal(transaction_blockedTemp)
                err = stub.PutState(transaction_blockedKey, transaction_blockedJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())

                }
                // Add to Transaction_contract
                transaction_contractArgs := []string{transactionID, transaction_blockedTemp.ClientID, transaction_blockedTemp.GroupbuyID, transaction_blockedTemp.Currency, fmt.Sprintf("%f", transaction_blockedTemp.Amount), status_contract}
                err = updateTransaction_contractCalled(stub, transaction_contractArgs)
                if err != nil {
                        return shim.Error(err.Error())
                }
		// Add to Transaction_contract_record
                transaction_contract_recordArgs := []string{transactionID, transaction_blockedTemp.ClientID, transaction_blockedTemp.GroupbuyID, transaction_blockedTemp.Currency, fmt.Sprintf("%f", transaction_blockedTemp.Amount), status_contract}
                err = updateTransaction_contract_recordCalled(stub, transaction_contract_recordArgs)
                if err != nil {
                        return shim.Error(err.Error())
                }
		// Add to Transaction_blocked_record
                transaction_blocked_recordArgs := []string{transactionID, transaction_blockedTemp.ClientID, transaction_blockedTemp.GroupbuyID, transaction_blockedTemp.Currency, fmt.Sprintf("%f",transaction_blockedTemp.Amount), transaction_blockedTemp.Status}
                err = updateTransaction_blocked_recordCalled(stub, transaction_blocked_recordArgs)
                if err != nil {
                        return shim.Error(err.Error())
                }

        }



        fmt.Println("Groupbuy leaved success!")
        return shim.Success(nil)
}
// =======================================================
// groupbuyUnblocked_failed
// Unblocked this transactionID failed,
// 1. Updates Transaction_blocked (status=blocked),
// 2. Adds record to transaction_blocked_record,
// 3. Updates record to Transaction_contract (status="fail"),
// 4. Adds record toTransaction_contract_record,
// input: TransactionID, ClientID, ProductID
// =======================================================
func (t * TsohueChainCode) groupbuyUnblocked_failed(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error
        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }
        transactionID := args[0]



        // Update to Transaction_blocked
        //query Transaction_blocked
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_blocked\",\"TransactionID\":\"%s\"}}", transactionID)
        // Get Groupbuy results from queryString
        resultsIterator, err := stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
        // Loop through all records of Groupbuy,
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                transaction_blockedTemp := &Transaction_blocked{}
                transaction_blockedKey := responseRange.Key

                // Unmarshal responseRange.Value to JSON
                err = json.Unmarshal(responseRange.Value, &transaction_blockedTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemup
                transaction_blockedTemp.Status = "blocked"
                transaction_blockedTemp.Timestamp = time.Now().String()

                // Marshal and Put groupbuy to state
                transaction_blockedJSONasBytes, _ := json.Marshal(transaction_blockedTemp)
                err = stub.PutState(transaction_blockedKey, transaction_blockedJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())

                }
		// Add to Transaction_blocked_record
                transaction_blocked_recordArgs := []string{transactionID, transaction_blockedTemp.ClientID, transaction_blockedTemp.GroupbuyID, transaction_blockedTemp.Currency, fmt.Sprintf("%f",transaction_blockedTemp.Amount), transaction_blockedTemp.Status}
                err = updateTransaction_blocked_recordCalled(stub, transaction_blocked_recordArgs)
                if err != nil {
                        return shim.Error(err.Error())
                }

        }



        // Update to Transaction_contract
        //query Transaction_contract
        queryString = fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_contract\",\"TransactionID\":\"%s\"}}", transactionID)
        // Get Groupbuy results from queryString
        resultsIterator, err = stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
        // Loop through all records of Groupbuy,
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                transaction_contractTemp := &Transaction_contract{}
                transaction_contractKey := responseRange.Key

                // Unmarshal responseRange.Value to JSON
                err = json.Unmarshal(responseRange.Value, &transaction_contractTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemup
                transaction_contractTemp.Status = "fail"
                transaction_contractTemp.Timestamp = time.Now().String()

                // Marshal and Put groupbuy to state
                transaction_contractJSONasBytes, _ := json.Marshal(transaction_contractTemp)
                err = stub.PutState(transaction_contractKey, transaction_contractJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())

                }
                // Add to Transaction_contract_record
                transaction_contract_recordArgs := []string{transactionID, transaction_contractTemp.ClientID, transaction_contractTemp.GroupbuyID, transaction_contractTemp.Currency, fmt.Sprintf("%f", transaction_contractTemp.Amount),  transaction_contractTemp.Status}
                err = updateTransaction_contract_recordCalled(stub, transaction_contract_recordArgs)
                if err != nil {
                        return shim.Error(err.Error())
                }



        }



        fmt.Println("Groupbuy leaved success!")
        return shim.Success(nil)
}
// =======================================================
// groupbuySuccess
// Changing groupbuy status to success,
// 1. Updates Groupbuy (status=success),
// 2. Add record to groupbuy_record,
// input: transaction_id, groupbuyID
// =======================================================
func (t *TsohueChainCode) groupbuySuccess(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error

        if len(args) != 2 {
                return shim.Error("Incorrect number of arguments. Expecting 2")
        }
	transaction_id := args[0]
        groupbuyID := args[1]

	//query groupbuy
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"groupbuy\",\"GroupbuyID\":\"%s\"}}", groupbuyID)
        // Get Groupbuy results from queryString
        resultsIterator, err := stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
        count := 0
        // Loop through all records of Groupbuy,
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                groupbuyTemp := &Groupbuy{}
                groupbuyKey := responseRange.Key

                // Unmarshal responseRange.Value to JSON
                err = json.Unmarshal(responseRange.Value, &groupbuyTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemp
                groupbuyTemp.Status = "success"
                groupbuyTemp.Timestamp = time.Now().String()

                // Marshal and Put groupbuy to state
                groupbuyJSONasBytes, _ := json.Marshal(groupbuyTemp)
                err = stub.PutState(groupbuyKey, groupbuyJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())
                }

	        //update Groupbuy_record
                status := "success"
                groupbuy_recordArgs := []string{transaction_id, groupbuyTemp.GroupbuyID, groupbuyTemp.Currency, fmt.Sprintf("%f", groupbuyTemp.Target_amount), fmt.Sprintf("%f", groupbuyTemp.Share), status}
                err = updateGroupbuy_recordCalled(stub, groupbuy_recordArgs)
                if err != nil {
                        return shim.Error(err.Error())

                } 

        }

        fmt.Printf("Successfully added %d asset records\n", count)
        return shim.Success(nil) 
}
// =======================================================
// groupbuyNoGoodPrice
// Changing groupbuy status to no good price,
// 1. Updates Groupbuy (status=no good price),
// 2. Add record to groupbuy_record,
// input: transaction_id, groupbuyID
// =======================================================
func (t *TsohueChainCode) groupbuyNoGoodPrice(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error

        if len(args) != 2 {
                return shim.Error("Incorrect number of arguments. Expecting 2")
        }
        transaction_id := args[0]
        groupbuyID := args[1]

        //query groupbuy
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"groupbuy\",\"GroupbuyID\":\"%s\"}}", groupbuyID)

        // Get Groupbuy results from queryString
        resultsIterator, err := stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
        count := 0

        // Loop through all records of Groupbuy,
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                groupbuyTemp := &Groupbuy{}
                groupbuyKey := responseRange.Key

                // Unmarshal responseRange.Value to JSON
                err = json.Unmarshal(responseRange.Value, &groupbuyTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemp
                groupbuyTemp.Status = "no good price"
                groupbuyTemp.Timestamp = time.Now().String()

                // Marshal and Put groupbuy to state
                groupbuyJSONasBytes, _ := json.Marshal(groupbuyTemp)
                err = stub.PutState(groupbuyKey, groupbuyJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())
                }

                //update Groupbuy_record
                status := "no good price"
                groupbuy_recordArgs := []string{transaction_id, groupbuyTemp.GroupbuyID, groupbuyTemp.Currency, fmt.Sprintf("%f", groupbuyTemp.Target_amount), fmt.Sprintf("%f", groupbuyTemp.Share), status}
                err = updateGroupbuy_recordCalled(stub, groupbuy_recordArgs)
                if err != nil {
                        return shim.Error(err.Error())

                }

        }

        fmt.Printf("Successfully added %d asset records\n", count)
        return shim.Success(nil)

}

// =======================================================
// groupbuyInvesting
// Changing groupbuy status to investing,
// 1. Update Groupbuy (status=Investing),
// 2. Add record to groupbuy_record,
// input: transaction_id, groupbuyID
// =======================================================
func (t *TsohueChainCode) groupbuyInvesting(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error

        if len(args) != 2 {
                return shim.Error("Incorrect number of arguments. Expecting 2")
        }
        transaction_id := args[0]
        groupbuyID := args[1]

        //query groupbuy
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"groupbuy\",\"GroupbuyID\":\"%s\"}}", groupbuyID)

        // Get Groupbuy results from queryString
        resultsIterator, err := stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
        count := 0

        // Loop through all records of Groupbuy,
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                groupbuyTemp := &Groupbuy{}
                groupbuyKey := responseRange.Key

                // Unmarshal responseRange.Value to JSON
                err = json.Unmarshal(responseRange.Value, &groupbuyTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemp
                groupbuyTemp.Status = "investing"
                groupbuyTemp.Timestamp = time.Now().String()

                // Marshal and Put groupbuy to state
                groupbuyJSONasBytes, _ := json.Marshal(groupbuyTemp)
                err = stub.PutState(groupbuyKey, groupbuyJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())
                }

                //update Groupbuy_record
                status := "investing"
                groupbuy_recordArgs := []string{transaction_id, groupbuyTemp.GroupbuyID, groupbuyTemp.Currency, fmt.Sprintf("%f", groupbuyTemp.Target_amount), fmt.Sprintf("%f", groupbuyTemp.Share), status}
                err = updateGroupbuy_recordCalled(stub, groupbuy_recordArgs)
                if err != nil {
                        return shim.Error(err.Error())

                }

        }

        fmt.Printf("Successfully added %d asset records\n", count)
        return shim.Success(nil)

}
// =======================================================
// groupbuyMataned
// Changing groupbuy status to mataned,
// 1. Update Groupbuy (status=mataned),
// 2. Add record to groupbuy_record,
// input: transaction_id, groupbuyID
// =======================================================
func (t *TsohueChainCode) groupbuyMataned(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error

        if len(args) != 2 {
                return shim.Error("Incorrect number of arguments. Expecting 2")
        }
        transaction_id := args[0]
        groupbuyID := args[1]

        //query groupbuy
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"groupbuy\",\"GroupbuyID\":\"%s\"}}", groupbuyID)

        // Get Groupbuy results from queryString
        resultsIterator, err := stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
        count := 0

        // Loop through all records of Groupbuy,
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                groupbuyTemp := &Groupbuy{}
                groupbuyKey := responseRange.Key

                // Unmarshal responseRange.Value to JSON
                err = json.Unmarshal(responseRange.Value, &groupbuyTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemp
                groupbuyTemp.Status = "mataned"
                groupbuyTemp.Timestamp = time.Now().String()

                // Marshal and Put groupbuy to state
                groupbuyJSONasBytes, _ := json.Marshal(groupbuyTemp)
                err = stub.PutState(groupbuyKey, groupbuyJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())
                }

                //update Groupbuy_record
                status := "mataned"
                groupbuy_recordArgs := []string{transaction_id, groupbuyTemp.GroupbuyID, groupbuyTemp.Currency, fmt.Sprintf("%f", groupbuyTemp.Target_amount), fmt.Sprintf("%f", groupbuyTemp.Share), status}
                err = updateGroupbuy_recordCalled(stub, groupbuy_recordArgs)
                if err != nil {
                        return shim.Error(err.Error())

                }

        }

        fmt.Printf("Successfully added %d asset records\n", count)
        return shim.Success(nil)

}
// =======================================================
// groupbuyLiquidate
// Changing groupbuy status to Liquidate,
// 1. Update Groupbuy (status=liquidate),
// 2. Add record to groupbuy_record,
// input: transaction_id, groupbuyID
// =======================================================
func (t *TsohueChainCode) groupbuyLiquidate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error

        if len(args) != 2 {
                return shim.Error("Incorrect number of arguments. Expecting 2")
        }
        transaction_id := args[0]
        groupbuyID := args[1]

        //query groupbuy
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"groupbuy\",\"GroupbuyID\":\"%s\"}}", groupbuyID)

        // Get Groupbuy results from queryString
        resultsIterator, err := stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
        count := 0

        // Loop through all records of Groupbuy,
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                groupbuyTemp := &Groupbuy{}
                groupbuyKey := responseRange.Key

                // Unmarshal responseRange.Value to JSON
                err = json.Unmarshal(responseRange.Value, &groupbuyTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemp
                groupbuyTemp.Status = "liquidate"
                groupbuyTemp.Timestamp = time.Now().String()

                // Marshal and Put groupbuy to state
                groupbuyJSONasBytes, _ := json.Marshal(groupbuyTemp)
                err = stub.PutState(groupbuyKey, groupbuyJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())
                }

                //update Groupbuy_record
                status := "liquidate"
                groupbuy_recordArgs := []string{transaction_id, groupbuyTemp.GroupbuyID, groupbuyTemp.Currency, fmt.Sprintf("%f", groupbuyTemp.Target_amount), fmt.Sprintf("%f", groupbuyTemp.Share), status}
                err = updateGroupbuy_recordCalled(stub, groupbuy_recordArgs)
                if err != nil {
                        return shim.Error(err.Error())

                }

        }

        fmt.Printf("Successfully added %d asset records\n", count)
        return shim.Success(nil)

}
// =======================================================
// groupbuyEnded
// Changing groupbuy status to ended,
// 1. Update Groupbuy (status=ended),
// 2. Add record to groupbuy_record,
// input: transaction_id, groupbuyID
// =======================================================
func (t *TsohueChainCode) groupbuyEnded(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error

        if len(args) != 2 {
                return shim.Error("Incorrect number of arguments. Expecting 2")
        }
        transaction_id := args[0]
        groupbuyID := args[1]

        //query groupbuy
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"groupbuy\",\"GroupbuyID\":\"%s\"}}", groupbuyID)

        // Get Groupbuy results from queryString
        resultsIterator, err := stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
        count := 0

        // Loop through all records of Groupbuy,
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                groupbuyTemp := &Groupbuy{}
                groupbuyKey := responseRange.Key

                // Unmarshal responseRange.Value to JSON
                err = json.Unmarshal(responseRange.Value, &groupbuyTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemp
                groupbuyTemp.Status = "ended"
                groupbuyTemp.Timestamp = time.Now().String()

                // Marshal and Put groupbuy to state
                groupbuyJSONasBytes, _ := json.Marshal(groupbuyTemp)
                err = stub.PutState(groupbuyKey, groupbuyJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())
                }

                //update Groupbuy_record
                status := "ended"
                groupbuy_recordArgs := []string{transaction_id, groupbuyTemp.GroupbuyID, groupbuyTemp.Currency, fmt.Sprintf("%f", groupbuyTemp.Target_amount), fmt.Sprintf("%f", groupbuyTemp.Share), status}
                err = updateGroupbuy_recordCalled(stub, groupbuy_recordArgs)
                if err != nil {
                        return shim.Error(err.Error())

                }

        }

        fmt.Printf("Successfully added %d asset records\n", count)
        return shim.Success(nil)

}
// =======================================================
// groupbuyFailed
// Changing groupbuy status to fail,
// 1.Update Groupbuy (status=fail)
// 2.Add record to groupbuy_record,
// input: transaction_id, groupbuyID
// =======================================================
func (t *TsohueChainCode) groupbuyFailed(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error

        if len(args) != 2 {
                return shim.Error("Incorrect number of arguments. Expecting 2")
        }
        transaction_id := args[0]
        groupbuyID := args[1]

        //query groupbuy
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"groupbuy\",\"GroupbuyID\":\"%s\"}}", groupbuyID)

        // Get Groupbuy results from queryString
        resultsIterator, err := stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
        count := 0

        // Loop through all records of Groupbuy,
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                groupbuyTemp := &Groupbuy{}
                groupbuyKey := responseRange.Key

                // Unmarshal responseRange.Value to JSON
                err = json.Unmarshal(responseRange.Value, &groupbuyTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemp
                groupbuyTemp.Status = "fail"
                groupbuyTemp.Timestamp = time.Now().String()

                // Marshal and Put groupbuy to state
                groupbuyJSONasBytes, _ := json.Marshal(groupbuyTemp)
                err = stub.PutState(groupbuyKey, groupbuyJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())
                }

                //update Groupbuy_record
                status := "fail"
                groupbuy_recordArgs := []string{transaction_id, groupbuyTemp.GroupbuyID, groupbuyTemp.Currency, fmt.Sprintf("%f", groupbuyTemp.Target_amount), fmt.Sprintf("%f", groupbuyTemp.Share), status}
                err = updateGroupbuy_recordCalled(stub, groupbuy_recordArgs)
                if err != nil {
                        return shim.Error(err.Error())

                }

        }

        fmt.Printf("Successfully added %d asset records\n", count)
        return shim.Success(nil)

}
// =======================================================
// dividend
// 1.Adds records to Groupbuy_contract
// 2.Adds records tp Groupbuy_contract_records
// status = dividend
// input: transactionID, clientID, groupbuyID, currency, amount,
// =======================================================
func (t *TsohueChainCode) dividend(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error

        if len(args) != 5 {
                return shim.Error("Incorrect number of arguments. Expecting 5")
        }
        transactionID := args[0]
        clientID := args[1]
        groupbuyID := args[2]
        currency := args[3]
	amount := args[4]
        status := "dividend"



        // Add to Transaction_contract
        transaction_contractArgs := []string{transactionID, clientID, groupbuyID, currency, amount, status}
        err = updateTransaction_contractCalled(stub, transaction_contractArgs)
        if err != nil {
                return shim.Error(err.Error())
        }



        // Add to Transaction_contract_record
        transaction_contract_recordArgs := []string{transactionID, clientID, groupbuyID, currency, amount, status}
        err = updateTransaction_contract_recordCalled(stub, transaction_contract_recordArgs)
        if err != nil {
                return shim.Error(err.Error())
        }



        fmt.Println("Groupbuy joined success!")
        return shim.Success(nil)

}
// =======================================================
// Capital
// 1.Adds records to Groupbuy_contract
// 2.Adds records tp Groupbuy_contract_records
// status = capital
// input: transactionID, clientID, groupbuyID, currency, amount
// =======================================================
func (t *TsohueChainCode) capital(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error

        if len(args) != 5 {
                return shim.Error("Incorrect number of arguments. Expecting 5")
        }
        transactionID := args[0]
        clientID := args[1]
        groupbuyID := args[2]
        currency := args[3]
        amount := args[4]
        status := "capital"



        // Add to Transaction_contract
        transaction_contractArgs := []string{transactionID, clientID, groupbuyID, currency, amount, status}
        err = updateTransaction_contractCalled(stub, transaction_contractArgs)
        if err != nil {
                return shim.Error(err.Error())
        }



        // Add to Transaction_contract_record
        transaction_contract_recordArgs := []string{transactionID, clientID, groupbuyID, currency, amount, status}
        err = updateTransaction_contract_recordCalled(stub, transaction_contract_recordArgs)
        if err != nil {
                return shim.Error(err.Error())
        }



        fmt.Println("Groupbuy joined success!")
        return shim.Success(nil)

}
// =======================================================
// liquidated
// 1.Adds records to Groupbuy_contract
// 2.Adds records tp Groupbuy_contract_records
// status = liquidate
// input: transactionID, clientID, groupbuyID, currency, amount
// =======================================================
func (t *TsohueChainCode) liquidated(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error

        if len(args) != 5 {
                return shim.Error("Incorrect number of arguments. Expecting 5")
        }
        transactionID := args[0]
        clientID := args[1]
        groupbuyID := args[2]
        currency := args[3]
        amount := args[4]
        status := "liquidate"



        // Add to Transaction_contract
        transaction_contractArgs := []string{transactionID, clientID, groupbuyID, currency, amount, status}
        err = updateTransaction_contractCalled(stub, transaction_contractArgs)
        if err != nil {
                return shim.Error(err.Error())
        }



        // Add to Transaction_contract_record
        transaction_contract_recordArgs := []string{transactionID, clientID, groupbuyID, currency, amount, status}
        err = updateTransaction_contract_recordCalled(stub, transaction_contract_recordArgs)
        if err != nil {
                return shim.Error(err.Error())
        }



        fmt.Println("Groupbuy joined success!")
        return shim.Success(nil)

}
// =======================================================
// contract_failed
// 1.update Transaction_concract Status = fail
// 2.Add to Transaction_contract_record
// input: transactionID
// =======================================================
func (t *TsohueChainCode) contract_failed(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }
        transactionID := args[0]



        //updateTransaction_concract
	queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_contract\",\"TransactionID\":\"%s\"}}", transactionID)

        // search record by transactionID
	resultsIterator, err := stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()

        // Loop through all records
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                transaction_contractTemp := &Transaction_contract{}
                transaction_contractKey := responseRange.Key

                // Unmarshal responseRange.Value to JSON
                err = json.Unmarshal(responseRange.Value, &transaction_contractTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemp
                transaction_contractTemp.Status = "fail"
                transaction_contractTemp.Timestamp = time.Now().String()

                // Marshal and Put groupbuy to state
                transaction_contractJSONasBytes, _ := json.Marshal(transaction_contractTemp)
                err = stub.PutState(transaction_contractKey, transaction_contractJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())

                }

	        // Add to Transaction_contract_record
                transaction_contract_recordArgs := []string{transactionID, transaction_contractTemp.ClientID, transaction_contractTemp.GroupbuyID, transaction_contractTemp.Currency, fmt.Sprintf("%f", transaction_contractTemp.Amount), transaction_contractTemp.Status}
                err = updateTransaction_contract_recordCalled(stub, transaction_contract_recordArgs)
                if err != nil {
                        return shim.Error(err.Error())
                }



        }
        fmt.Printf("Successfully update transaction_contract\n")
        fmt.Printf("TakeOver failed success!")
        return shim.Success(nil)

}
// =======================================================
// takeover_raise
// 1.Adds records to Transaction_takeover (status=open),
// 2.Adds records to Transaction_takeover_record,
// 3.Adds records to Transaction_blocked (status=blocked),
// 4.Adds records to Transaction_blocked_record,
// input: transactionID, groupbuyID, currency, amount, clientID_sell
// =======================================================

func (t *TsohueChainCode) takeOver_raise(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error

        if len(args) != 5 {
                return shim.Error("Incorrect number of arguments. Expecting 5")
        }
        transactionID := args[0]
        groupbuyID := args[1]
	currency := args[2]
	amount := args[3]
        clientID_sell := args[4]
	clientID_buy := "null"
	status_takeover := "open"
	status_blocked := "blocked"



        //updateTransaction_takeover
        transaction_takeover_Args := []string{transactionID, groupbuyID, clientID_sell, clientID_buy, status_takeover}
        err = updateTransaction_takeoverCalled(stub, transaction_takeover_Args)
        if err != nil {
        return shim.Error(err.Error())
        }

        fmt.Println("Takeover raise  success!")



        //updateTransaction_takeover_record
        transaction_takeover_record_Args := []string{transactionID, groupbuyID, clientID_sell, clientID_buy, status_takeover}
        err = updateTransaction_takeover_recordCalled(stub, transaction_takeover_record_Args)
        if err != nil {
        return shim.Error(err.Error())
        }

        fmt.Println("Takeover raise  success!")



	//updateTransaction_blocked
        transaction_blocked_Args := []string{transactionID, clientID_sell, groupbuyID, currency, amount, status_blocked}
        err = updateTransaction_blockedCalled(stub, transaction_blocked_Args)
        if err != nil {
        return shim.Error(err.Error())
        }

        fmt.Println("Takeover raise  success!")



        //updateTransaction_blocked_record
        transaction_blocked_record_Args := []string{transactionID, clientID_sell, groupbuyID, currency, amount, status_blocked}
        err = updateTransaction_blocked_recordCalled(stub, transaction_blocked_record_Args)
        if err != nil {
        return shim.Error(err.Error())
        }

        fmt.Println("Takeover raise  success!")
        return shim.Success(nil)



}
// =======================================================
// takeover_failed
// 1.Updates Transaction_takeover (status=fail)
// 2.Adds records to Transaction_takeover_record,
// 3.Updates Transaction_blocked (status=fail)
// 4.Add buyer transaction record to transaction
// input: ClientID_sell, ClientID_buy, ProductID.
// =======================================================

func (t *TsohueChainCode) takeOver_failed(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }
        transactionID := args[0]


	//updateTransaction_takeover
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_takeover\",\"TransactionID\":\"%s\"}}", transactionID)
        // search transactionID
        resultsIterator, err := stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
	// Loop through all records
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                transaction_takeoverTemp := &Transaction_takeover{}
                transaction_takeoverKey := responseRange.Key

                // Unmarshal responseRange.Value to JSON
                err = json.Unmarshal(responseRange.Value, &transaction_takeoverTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update transaction_takeoverTemp
                transaction_takeoverTemp.Status = "fail"
                transaction_takeoverTemp.Timestamp = time.Now().String()

                // Marshal and Put groupbuy to state
                transaction_takeoverJSONasBytes, _ := json.Marshal(transaction_takeoverTemp)
                err = stub.PutState(transaction_takeoverKey, transaction_takeoverJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())

                }

		//updateTransaction_takeover_record
                transaction_takeover_record_Args := []string{transactionID, transaction_takeoverTemp.GroupbuyID, transaction_takeoverTemp.ClientID_sell, transaction_takeoverTemp.ClientID_buy, "fail"}
                err = updateTransaction_takeover_recordCalled(stub, transaction_takeover_record_Args)
                if err != nil {
                return shim.Error(err.Error())
                }


        }

        fmt.Printf("Successfully added %d asset records\n")



	//updateTransaction_blocked
        queryString = fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_blocked\",\"TransactionID\":\"%s\"}}", transactionID)
	// search record by transactionID
        resultsIterator, err = stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
        // Loop through all records
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                transaction_blockedTemp := &Transaction_blocked{}
                transaction_blockedKey := responseRange.Key

                // Unmarshal responseRange.Value to JSON
                err = json.Unmarshal(responseRange.Value, &transaction_blockedTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemp
                transaction_blockedTemp.Status = "fail"
                transaction_blockedTemp.Timestamp = time.Now().String()

                // Marshal and Put groupbuy to state
                transaction_blockedJSONasBytes, _ := json.Marshal(transaction_blockedTemp)
                err = stub.PutState(transaction_blockedKey, transaction_blockedJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())

                }

	        //updateTransaction_blocked_record
                transaction_blocked_record_Args := []string{transactionID, transaction_blockedTemp.ClientID, transaction_blockedTemp.GroupbuyID, transaction_blockedTemp.Currency, fmt.Sprintf("%f",transaction_blockedTemp.Amount), "fail"}
                err = updateTransaction_blocked_recordCalled(stub, transaction_blocked_record_Args)
                if err != nil {
                        return shim.Error(err.Error())
                }

        }
        fmt.Printf("Successfully added %d asset records\n")



        //updateTransaction_concract
        queryString = fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_contract\",\"TransactionID\":\"%s\"}}", transactionID)
        // search record by transactionID
        resultsIterator, err = stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()

        count := 0
        // Loop through all records
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                transaction_contractTemp := &Transaction_contract{}
                transaction_contractKey := responseRange.Key

                // Unmarshal responseRange.Value to JSON
                err = json.Unmarshal(responseRange.Value, &transaction_contractTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update transaction_contractTemp
                transaction_contractTemp.Status = "fail"
                transaction_contractTemp.Timestamp = time.Now().String()

                // Marshal and Put groupbuy to state
                transaction_contractJSONasBytes, _ := json.Marshal(transaction_contractTemp)
                err = stub.PutState(transaction_contractKey, transaction_contractJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())

                }

	        // Add clientID_sell to Transaction_contract_record
                transaction_contract_recordArgs := []string{transactionID, transaction_contractTemp.ClientID, transaction_contractTemp.GroupbuyID, transaction_contractTemp.Currency, fmt.Sprintf("%f",transaction_contractTemp.Amount), transaction_contractTemp.Status}
                err = updateTransaction_contract_recordCalled(stub, transaction_contract_recordArgs)
                if err != nil {
                        return shim.Error(err.Error())
                }
                count ++ 

        }
        fmt.Printf("Successfully added %d asset records\n")
        fmt.Printf("TakeOver failed success!")
        return shim.Success(nil)
}

// =======================================================
// takeover 
// 1.Updates Transaction_takeover (status=success),
// 2.Adds Transaction_takeover_record,
// 3.Adds clientID_sell to Transaction_contract
// 4.Adds clientID_sell to Transaction_contract_record,
// 5.Adds clientID_buy to Transaction_contract,
// 6.Adds clientID_buy to Transaction_contract_record,
// 7.Updates Transaction_blocked (status=unblocked),
// 8.Adds Transaction_blocked_record,
// input: transactionID, groupbuyID, currency, amount, clientID_sell, clientID_buy
// =======================================================

func (t *TsohueChainCode) takeOver(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error

        if len(args) != 6 {
                return shim.Error("Incorrect number of arguments. Expecting 6")
        }
	transactionID := args[0]
	groupbuyID := args[1]
        currency := args[2]
        amount := args[3]
	clientID_sell := args[4]
	clientID_buy := args[5]



        //updateTransaction_takeover
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_takeover\",\"TransactionID\":\"%s\"}}", transactionID)

        // search transactionID
        resultsIterator, err := stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()

        // Loop through all records
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                transaction_takeoverTemp := &Transaction_takeover{}
                transaction_takeoverKey := responseRange.Key

                // Unmarshal responseRange.Value to JSON
                err = json.Unmarshal(responseRange.Value, &transaction_takeoverTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemp
                transaction_takeoverTemp.Status = "success"
                transaction_takeoverTemp.Timestamp = time.Now().String()

                // Marshal and Put groupbuy to state
                transaction_takeoverJSONasBytes, _ := json.Marshal(transaction_takeoverTemp)
                err = stub.PutState(transaction_takeoverKey, transaction_takeoverJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())

                }

        }
        fmt.Printf("Successfully added %d asset records\n")



        //updateTransaction_takeover_record
        transaction_takeover_record_Args := []string{transactionID, groupbuyID, clientID_sell, clientID_buy, "success"}
        err = updateTransaction_takeover_recordCalled(stub, transaction_takeover_record_Args)
        if err != nil {
        return shim.Error(err.Error())
        }



	//Update Transaction_contract
        // Add clientID_sell to Transaction_contract
        transaction_contractArgs_sell := []string{transactionID, clientID_sell, groupbuyID, currency, amount, "takeover_sell"}
        err = updateTransaction_contractCalled(stub, transaction_contractArgs_sell)
        if err != nil {
                return shim.Error(err.Error())
        }



        //Adds clientID_sell to Transaction_contract_record
        transaction_contract_recordArgs_sell := []string{transactionID, clientID_sell, groupbuyID, currency, amount, "takeover_sell"}
        err = updateTransaction_contract_recordCalled(stub, transaction_contract_recordArgs_sell)
        if err != nil {
                return shim.Error(err.Error())
        }



	// Add clientID_buy to Transaction_contract
	transaction_contractArgs_buy := []string{transactionID, clientID_buy, groupbuyID, currency, amount, "takeover_buy"}
        err = updateTransaction_contractCalled(stub, transaction_contractArgs_buy)
        if err != nil {
                return shim.Error(err.Error())
        }



        // Add clientID_buy to Transaction_contract_record
	transaction_contract_recordArgs_buy := []string{transactionID, clientID_buy, groupbuyID, currency, amount, "takeover_buy"}
        err = updateTransaction_contract_recordCalled(stub, transaction_contract_recordArgs_buy)
        if err != nil {
                return shim.Error(err.Error())
        }



        //updateTransaction_blocked
        queryString = fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_blocked\",\"TransactionID\":\"%s\"}}", transactionID)
        // search record by transactionID
        resultsIterator, err = stub.GetQueryResult(queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()
        // Loop through all records
        for resultsIterator.HasNext() {
                responseRange, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                transaction_blockedTemp := &Transaction_blocked{}
                transaction_blockedKey := responseRange.Key

                // Unmarshal responseRange.Value to JSON
                err = json.Unmarshal(responseRange.Value, &transaction_blockedTemp)
                if err != nil {
                        return shim.Error(err.Error())
                }

                // Update groupbuyTemp
		transaction_blockedTemp.Status = "unblocked"
                transaction_blockedTemp.Timestamp = time.Now().String()

                // Marshal and Put groupbuy to state
                transaction_blockedJSONasBytes, _ := json.Marshal(transaction_blockedTemp)
                err = stub.PutState(transaction_blockedKey, transaction_blockedJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())

                }

        }
        fmt.Printf("Successfully added %d asset records\n")



        //updateTransaction_blocked_record
        transaction_blocked_record_Args := []string{transactionID, clientID_sell, groupbuyID, currency, amount, "unblocked"}
        err = updateTransaction_blocked_recordCalled(stub, transaction_blocked_record_Args)
        if err != nil {
        return shim.Error(err.Error())
        }



        fmt.Printf("Takeover success!")
        return shim.Success(nil)



}
//// Query functions ////
// =======================================================
// queryGroupbuyByProductID - Query all groupbuy transactions based on productID
// input: groupbuyID
// =======================================================
func (t *TsohueChainCode) queryGroupbuyByGroupbuyID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        groupbuyID := args[0]

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"groupbuy\",\"GroupbuyID\":\"%s\"}}", groupbuyID)
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryAssetHistory - query history of asset of client
// input: clientID, productID, currency
// =======================================================
func (t *TsohueChainCode) queryGroupbuyHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var objectType, groupbuyID string
        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }
        objectType = "gorupbuy"
        groupbuyID = args[0]

        // Create composite key
        groupbuyKeyString := objectType + "_" + groupbuyID
        groupbuyKey, err := stub.CreateCompositeKey(groupbuyKeyString, []string{objectType, groupbuyID})
        if err != nil {
                return shim.Error(err.Error())
        }

        // Get history states
        historyIer, err := stub.GetHistoryForKey(groupbuyKey)
        if err != nil {
                return shim.Error(err.Error())
        } else if historyIer == nil {
                return shim.Error("Cannot find groupbuy: " + groupbuyKeyString)
        }
        var historySliceBytes []byte

        for historyIer.HasNext() {
                modification, err := historyIer.Next()
                if err != nil {
                        return shim.Error("Error reading history for client: " + err.Error())
                }
                historySliceBytes = append(historySliceBytes, modification.Value...)
                fmt.Println(string(modification.Value))
        }
        fmt.Println("Get history finished.")
        fmt.Println("================================")
        return shim.Success(historySliceBytes)
}
// =======================================================
// queryGroupbuyByProductID - Query all groupbuy transactions based on productID
// input: ProductID
// =======================================================
func (t *TsohueChainCode) queryGroupbuy_recordByGroupbuyID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        groupbuyID := args[0]

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"groupbuy_record\",\"GroupbuyID\":\"%s\"}}", groupbuyID)
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryGroupbuyByProductID - Query all groupbuy transactions based on productID
// input: ProductID
// =======================================================
func (t *TsohueChainCode) queryGroupbuy_recordByTransactionID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        transactionID := args[0]

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"groupbuy_record\",\"TransactionID\":\"%s\"}}", transactionID)
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryTransactionByProductID - Query all transactions records based on productID
// input: ProductID
// =======================================================
func (t *TsohueChainCode) queryTransaction_blockedByClientID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        clientID := args[0]
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_blocked\",\"ClientID\":\"%s\"}}", clientID)
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryTransactionByProductID - Query all transactions records based on productID
// input: ProductID
// =======================================================
func (t *TsohueChainCode) queryTransaction_blockedByTransactionID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        transactionID := args[0]
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_blocked\",\"TransactionID\":\"%s\"}}", transactionID)
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryTransactionByProductID - Query all transactions records based on productID
// input: ProductID
// =======================================================
func (t *TsohueChainCode) queryTransaction_blocked_recordByTransactionID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        transactionID := args[0]
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_blocked_record\",\"TransactionID\":\"%s\"}}", transactionID)
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryTransactionByProductID - Query all transactions records based on productID
// input: ProductID
// =======================================================
func (t *TsohueChainCode) queryTransaction_contractByGroupbuyID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        groupbuyID := args[0]

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_contract\",\"GroupbuyID\":\"%s\"}}", groupbuyID)
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryTransactionByClientID - Query all transactions records based on clientID
// input: ClientID
// =======================================================
func (t *TsohueChainCode) queryTransaction_contractByClientID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        clientID := args[0]

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_contract\",\"ClientID\":\"%s\"}}", clientID)
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryTransactionByClientID - Query all transactions records based on clientID
// input: ClientID
// =======================================================
func (t *TsohueChainCode) queryTransaction_contractByTransactionID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        clientID := args[0]

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_contract\",\"TransactionID\":\"%s\"}}", clientID)
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryTransactionByClientID - Query all transactions records based on clientID
// input: ClientID
// =======================================================
func (t *TsohueChainCode) queryTransaction_contract_recordByTransactionID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        clientID := args[0]

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_contract_record\",\"TransactionID\":\"%s\"}}", clientID)
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryTransactionByProductID - Query all transactions records based on productID
// input: ProductID
// =======================================================
func (t *TsohueChainCode) queryTransaction_action_joinByClientID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        clientID := args[0]
        status := "join"
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_action\",\"ClientID\":\"%s\",\"Status\":\"%s\"}}", clientID, status)
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryTransactionByProductID - Query all transactions records based on productID
// input: ProductID
// =======================================================
func (t *TsohueChainCode) queryTransaction_action_leaveByClientID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        clientID := args[0]
        status := "leave"
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_action\",\"ClientID\":\"%s\",\"Status\":\"%s\"}}", clientID, status)
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryTransactionByProductID - Query all transactions records based on productID
// input: ProductID
// =======================================================
func (t *TsohueChainCode) queryTransaction_action_joinByGroupbuyID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        groupbuyID := args[0]
        status := "join"
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_action\",\"GroupbuyID\":\"%s\",\"Status\":\"%s\"}}", groupbuyID, status)
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryTransactionByProductID - Query all transactions records based on productID
// input: ProductID
// =======================================================
func (t *TsohueChainCode) queryTransaction_action_leaveByGroupbuyID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        groupbuyID := args[0]
        status := "leave"
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_action\",\"GroupbuyID\":\"%s\",\"Status\":\"%s\"}}", groupbuyID, status)
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryTransactionByProductID - Query all transactions records based on productID
// input: ProductID
// =======================================================
func (t *TsohueChainCode) queryTransaction_actionByTransactionID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        transactionID := args[0]
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_action\",\"TransactionID\":\"%s\"}}", transactionID)
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryTransactionByProductID - Query all transactions records based on productID
// input: ProductID
// =======================================================
func (t *TsohueChainCode) queryTransaction_action_recordByTransactionID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        transactionID := args[0]
        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_action_record\",\"TransactionID\":\"%s\"}}", transactionID)
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryTransactionByProductID - Query all transactions records based on productID
// input: ProductID
// =======================================================
func (t *TsohueChainCode) queryTransaction_takeoverByClientID_sell(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        clientID := args[0]

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_takeover\",\"ClientID_sell\":\"%s\"}}", clientID)
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryTransactionByProductID - Query all transactions records based on productID
// input: ProductID
// =======================================================
func (t *TsohueChainCode) queryTransaction_takeoverByClientID_buy(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        clientID := args[0]

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_takeover\",\"ClientID_buy\":\"%s\"}}", clientID)
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryTransactionByProductID - Query all transactions records based on productID
// input: ProductID
// =======================================================
func (t *TsohueChainCode) queryTransaction_takeoverByGroupbuyID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        groupbuyID := args[0]

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_takeover\",\"GroupbuyID\":\"%s\"}}", groupbuyID)
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryTransactionByProductID - Query all transactions records based on productID
// input: ProductID
// =======================================================
func (t *TsohueChainCode) queryTransaction_takeoverByTransactionID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        transactionID := args[0]

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_takeover\",\"TransactionID\":\"%s\"}}", transactionID)
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryTransactionByProductID - Query all transactions records based on productID
// input: ProductID
// =======================================================
func (t *TsohueChainCode) queryTransaction_takeover_recordByTransactionID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        transactionID := args[0]

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_takeover_record\",\"TransactionID\":\"%s\"}}", transactionID)
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryUserByClientID - Query all user records based on clientID
// input: ClientID
// =======================================================
func (t *TsohueChainCode) queryUserByClientID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        clientID := args[0]

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"user\",\"ClientID\":\"%s\"}}", clientID)
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryAssetHistory - query history of asset of client
// input: clientID, productID, currency
// =======================================================
func (t *TsohueChainCode) queryAssetHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var objectType, clientID, productID, currency string
        if len(args) != 3 {
                return shim.Error("Incorrect number of arguments. Expecting 3")
        }

        objectType = "asset"
        clientID = args[0]
        productID = args[1]
        currency = args[2]


        // Create composite key
        assetKeyString := objectType + "_" + clientID + "_" + productID + "_" + currency
        assetKey, err := stub.CreateCompositeKey(assetKeyString, []string{objectType, clientID, productID, currency})
        if err != nil {
                return shim.Error(err.Error())
        }

        // Get history states
        historyIer, err := stub.GetHistoryForKey(assetKey)
        if err != nil {
                return shim.Error(err.Error())
        } else if historyIer == nil {
                return shim.Error("Cannot find asset: " + assetKeyString)
        }
        var historySliceBytes []byte

        for historyIer.HasNext() {
                modification, err := historyIer.Next()
                if err != nil {
                        return shim.Error("Error reading history for client: " + err.Error())
                }
                historySliceBytes = append(historySliceBytes, modification.Value...)
                fmt.Println(string(modification.Value))
        }
        fmt.Println("Get history finished.")
        fmt.Println("================================")
        return shim.Success(historySliceBytes)
}
// =======================================================
// queryUserByClientID - Query all user records based on clientID
// input: ClientID
// =======================================================
func (t *TsohueChainCode) queryUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 0 {
                return shim.Error("Incorrect number of arguments. Expecting ")
        }

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"user\"}}")
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryUserByClientID - Query all user records based on clientID
// input: ClientID
// =======================================================
func (t *TsohueChainCode) queryUser_record(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 0 {
                return shim.Error("Incorrect number of arguments. Expecting ")
        }

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"user_record\"}}")
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryUserByClientID - Query all user records based on clientID
// input: ClientID
// =======================================================
func (t *TsohueChainCode) queryGroupbuy(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 0 {
                return shim.Error("Incorrect number of arguments. Expecting ")
        }

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"groupbuy\"}}")
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryUserByClientID - Query all user records based on clientID
// input: ClientID
// =======================================================
func (t *TsohueChainCode) queryGroupbuy_record(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 0 {
                return shim.Error("Incorrect number of arguments. Expecting ")
        }

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"groupbuy_record\"}}")
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryUserByClientID - Query all user records based on clientID
// input: ClientID
// =======================================================
func (t *TsohueChainCode) queryTransaction_blocked(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 0 {
                return shim.Error("Incorrect number of arguments. Expecting ")
        }

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_blocked\"}}")
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryUserByClientID - Query all user records based on clientID
// input: ClientID
// =======================================================
func (t *TsohueChainCode) queryTransaction_blocked_record(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 0 {
                return shim.Error("Incorrect number of arguments. Expecting ")
        }

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_blocked_record\"}}")
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryUserByClientID - Query all user records based on clientID
// input: ClientID
// =======================================================
func (t *TsohueChainCode) queryTransaction_contract(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 0 {
                return shim.Error("Incorrect number of arguments. Expecting ")
        }

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_contract\"}}")
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryUserByClientID - Query all user records based on clientID
// input: ClientID
// =======================================================
func (t *TsohueChainCode) queryTransaction_contract_record(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 0 {
                return shim.Error("Incorrect number of arguments. Expecting ")
        }

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_contract_record\"}}")
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryUserByClientID - Query all user records based on clientID
// input: ClientID
// =======================================================
func (t *TsohueChainCode) queryTransaction_action(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 0 {
                return shim.Error("Incorrect number of arguments. Expecting ")
        }

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_action\"}}")
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryUserByClientID - Query all user records based on clientID
// input: ClientID
// =======================================================
func (t *TsohueChainCode) queryTransaction_action_record(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 0 {
                return shim.Error("Incorrect number of arguments. Expecting ")
        }

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_action_record\"}}")
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryUserByClientID - Query all user records based on clientID
// input: ClientID
// =======================================================
func (t *TsohueChainCode) queryTransaction_takeover(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 0 {
                return shim.Error("Incorrect number of arguments. Expecting ")
        }

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_takeover\"}}")
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryUserByClientID - Query all user records based on clientID
// input: ClientID
// =======================================================
func (t *TsohueChainCode) queryTransaction_takeover_record(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 0 {
                return shim.Error("Incorrect number of arguments. Expecting ")
        }

        queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"transaction_takeover\"}}")
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        //fmt.Println(string(queryResults))
        return shim.Success(queryResults)
}
// =======================================================
// queryByString - Query by input string
// input: queryString
// Examples:
// {\"selector\":{\"ClientID\":\"Jim\"}, \"use_index\":[\"indexAssetDoc\", \"indexAsset\"]}
// {\"selector\":{\"ObjectType\":\"asset\",\"ClientID\":\"Jim\"}}
// =======================================================
func (t *TsohueChainCode) queryByString(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var queryString string
        var err error

        if len(args) != 1 {
                return shim.Error("Incorrct number of arguments. Expecting 1")
        }

        queryString = args[0]
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        }
        fmt.Println("My result: " + string(queryResults))
        return shim.Success(queryResults)
}
// ===================================
// Function for queries
// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
        fmt.Printf("- getQueryResultForQueryString:\n%s\n", queryString)

        resultsIterator, err := stub.GetQueryResult(queryString)
        if err != nil {
                return nil, err
        }
        defer resultsIterator.Close()

        buffer, err := constructQueryResponseFromIterator(resultsIterator)
        if err != nil {
                return nil, err
        }

        //fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
        return buffer.Bytes(), nil
}
// ===================================
// Function for queries
// ===========================================================================================
// constructQueryResponseFromIterator constructs a JSON array containing query results from
// a given result iterator
// ===========================================================================================
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
        // buffer is a JSON array containing QueryResults
        var buffer bytes.Buffer
        buffer.WriteString("[")

        bArrayMemberAlreadyWritten := false
        for resultsIterator.HasNext() {
                queryResponse, err := resultsIterator.Next()
                if err != nil {
                        return nil, err
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
                //Record is a JSON object, so we write as-is
                buffer.WriteString(string(queryResponse.Value))
                buffer.WriteString("}")
                bArrayMemberAlreadyWritten = true

                fmt.Println(string(queryResponse.Value))
        }
        buffer.WriteString("]")
        return &buffer, nil
}
// ======= Non INVOKE functions =======================
// ====================================================




// =======================================================
// Called buy other functions to update Groupbuy. Not an invoke function!!!
// If update status, input share=0
//
// updateGroupbuyCalled - Update groupbuy, adds up the share. Client=NULL is total amount for groupbuy
// input: ClientID, ProductID, currency, share, status.
// =======================================================
func updateGroupbuyCalled(stub shim.ChaincodeStubInterface, args []string) error {
        var objectType, groupbuyID, productID, currency, status, dividend, expiry_date, capital string
        var share float64
	var target_amount float64
        var err error

        if len(args) != 9 {
                return fmt.Errorf("Incorrect number of arguments. Expecting 9")
        }
        objectType = "groupbuy"
	groupbuyID = args[0]
        productID = args[1]
        currency = args[2]
	target_amount, err = strconv.ParseFloat(args[3], 64)
        share, err = strconv.ParseFloat(args[4], 64)
        if err != nil {
                return fmt.Errorf("4th argument must be a numeric string")
        }
        status = args[5]
	dividend = args[6]
	expiry_date = args[7]
	capital = args[8]

        // Create composite key
        groupbuyKeyString := objectType + "_" + productID 
        groupbuyKey, err := stub.CreateCompositeKey(groupbuyKeyString, []string{objectType, groupbuyID})
        if err != nil {
                fmt.Errorf(err.Error())
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- start updateGroupbuyCalled " + groupbuyKeyString)

        // Get state with composite key
        groupbuyTemp := &Groupbuy{}
        groupbuyAsBytes, err := stub.GetState(groupbuyKey)
        if err != nil {
                return fmt.Errorf("Failed to get state for :" + groupbuyKeyString)
        } else if groupbuyAsBytes == nil {
                // We don't need to check if client has joined groupbuy already.
                // But if client doesn't have existing groupbuy, create one
                groupbuyTemp = &Groupbuy{
                        ObjectType : objectType,
                        Timestamp : time.Now().String(),
			GroupbuyID : groupbuyID,
                        ProductID : productID,
                        Currency : currency,
			Target_amount : 0, 
                        Share : 0,
                        Status : status,
			Dividend : dividend,
			Expiry_date : expiry_date,
			Capital : capital,
                }
                fmt.Println("Groupbuy created: " + groupbuyKeyString)
        } else {
                // Unmarshal client and update value
                err = json.Unmarshal(groupbuyAsBytes, &groupbuyTemp)
                if err != nil {
                        return fmt.Errorf(err.Error())
                }
        }

        // Add the share. (quantity can be a negative number), update timestamp and status
        if groupbuyTemp.Share + share >= 0 {
                groupbuyTemp.Share += share
                groupbuyTemp.Timestamp = time.Now().String()
                groupbuyTemp.Status = status
        } else {
                return fmt.Errorf("Not sufficient shares for " + productID)
        }

        // Add the target_amount. (quantity can be a negative number), update timestamp and status
        if groupbuyTemp.Target_amount + target_amount >= 0 {
                groupbuyTemp.Target_amount += target_amount
                groupbuyTemp.Timestamp = time.Now().String()
                groupbuyTemp.Status = status
        } else {
                return fmt.Errorf("Not sufficient shares for " + productID)
        }


        groupbuyJSONasBytes, err := json.Marshal(groupbuyTemp)
        if err != nil {
                return fmt.Errorf(err.Error())
        }
        // Save asset to state
        err = stub.PutState(groupbuyKey, groupbuyJSONasBytes)
        if err != nil {
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- end updateGroupbuyCalled (success)")
        fmt.Println("==============================")
        return nil
}
// =======================================================
// Called buy other functions to update Groupbuy. Not an invoke function!!!
// If update status, input share=0
//
// updateGroupbuyCalled - Update groupbuy, adds up the share. Client=NULL is total amount for groupbuy
// input: ClientID, ProductID, currency, share, status.
// =======================================================
func updateGroupbuy_recordCalled(stub shim.ChaincodeStubInterface, args []string) error {
        var objectType, groupbuyID, transactionID, currency, status string
	var share float64
	var target_amount float64
        var err error

        if len(args) != 6 {
                return fmt.Errorf("Incorrect number of arguments. Expecting 6")
        }
        objectType = "groupbuy_record"
        transactionID = args[0]
	groupbuyID = args[1]
        currency = args[2]
	target_amount, err = strconv.ParseFloat(args[3], 64)
	if err != nil {
                return fmt.Errorf("3th argument must be a numeric string")
        }
        share, err = strconv.ParseFloat(args[4], 64)
        if err != nil {
                return fmt.Errorf("4th argument must be a numeric string")
        }
        status = args[5]

        // Create composite key
        groupbuy_recordKeyString := objectType + "_" + groupbuyID + "_" + transactionID
        groupbuy_recordKey, err := stub.CreateCompositeKey(groupbuy_recordKeyString, []string{objectType, groupbuyID, transactionID})
        if err != nil {
                fmt.Errorf(err.Error())
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- start updateGroupbuyCalled " + groupbuy_recordKeyString)

        // Get state with composite key
        groupbuy_recordTemp := &Groupbuy_record{}
        groupbuy_recordAsBytes, err := stub.GetState(groupbuy_recordKey)
        if err != nil {
                return fmt.Errorf("Failed to get state for :" + groupbuy_recordKeyString)
        } else if groupbuy_recordAsBytes == nil {
                // We don't need to check if client has joined groupbuy already.
                // But if client doesn't have existing groupbuy, create one
                groupbuy_recordTemp = &Groupbuy_record{
                        ObjectType : objectType,
                        Timestamp : time.Now().String(),
		        TransactionID :transactionID,
			GroupbuyID : groupbuyID,
                        Currency : currency,
			Target_amount : target_amount,
                        Share : 0,
                        Status : status,
                }
                fmt.Println("Groupbuy_record created: " + groupbuy_recordKeyString)
        } else {
                // Unmarshal client and update value
                err = json.Unmarshal(groupbuy_recordAsBytes, &groupbuy_recordTemp)
                                                                                            if err != nil {
                        return fmt.Errorf(err.Error())
                }
        }

        // Add the share. (quantity can be a negative number), update timestamp and status
        if groupbuy_recordTemp.Share + share >= 0 {
                groupbuy_recordTemp.Share += share
                groupbuy_recordTemp.Timestamp = time.Now().String()
                groupbuy_recordTemp.Status = status
        } else {
                return fmt.Errorf("Not sufficient shares for " + groupbuyID)
        }

        groupbuy_recordJSONasBytes, err := json.Marshal(groupbuy_recordTemp)
        if err != nil {
                return fmt.Errorf(err.Error())
        }
        // Save asset to state
        err = stub.PutState(groupbuy_recordKey, groupbuy_recordJSONasBytes)
        if err != nil {
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- end updateGroupbuy_recordCalled (success)")
        fmt.Println("==============================")
        return nil
}
// =======================================================
// Called buy other functions to update User. Not an invoke function!!!
//
// updateUserCalled - Update User,
// input: ClientID, UserID.
// =======================================================
func updateUserCalled(stub shim.ChaincodeStubInterface, args []string) error {
        var objectType, userID, clientID, bank, bank_account, address, phone, status_open, status_bind string
        var err error
        if len(args) != 8 {
                return fmt.Errorf("Incorrect number of arguments. Expecting 8")
        }
        objectType = "user"
        clientID = args[0]
        userID = args[1]
        bank = args[2]
        bank_account = args[3]
        address = args[4]
        phone = args[5]
        status_open = args[6]
        status_bind = args[7]
        // Create composite key
        userKeyString := objectType  + "_" + clientID + "_" + userID
        userKey, err := stub.CreateCompositeKey(userKeyString, []string{objectType, clientID, userID})
        if err != nil {
                fmt.Errorf(err.Error())
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- start updateUserCalled " + userKeyString)

        // Get state with composite key
        userTemp := &User{}
        userAsBytes, err := stub.GetState(userKey)
        if err != nil {
                return fmt.Errorf("Failed to get state for :" + userKeyString)
        } else if userAsBytes == nil {
                // We don't need to check if client has joined groupbuy already.
                // But if client doesn't have existing groupbuy, create one
                userTemp = &User{
                        ObjectType : objectType,
			Timestamp : time.Now().String(),
                        ClientID : clientID,
                        UserID : userID,
			Bank : bank,
			Bank_account : bank_account,
			Address : address,
			Phone : phone,
                        Status_open : status_open,
                        Status_bind : status_bind,
                }
                fmt.Println("User created: " + userKeyString)
        } else {
                // Unmarshal client and update value
                err = json.Unmarshal(userAsBytes, &userTemp)
                if err != nil {
                        return fmt.Errorf(err.Error())
                }
        }

        userJSONasBytes, err := json.Marshal(userTemp)
        if err != nil {
                return fmt.Errorf(err.Error())
        }

        // Save asset to state
        err = stub.PutState(userKey, userJSONasBytes)
        if err != nil {
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- end Called (success)")
        fmt.Println("==============================")
        return nil


}
// =======================================================
// Called buy other functions to update User. Not an invoke function!!!
//
// updateUserCalled - Update User,
// input: ClientID, UserID.
// =======================================================
func updateUser_recordCalled(stub shim.ChaincodeStubInterface, args []string) error {
        var objectType, userID, clientID, bank, bank_account, address, phone, status_open, status_bind string
        var err error
        if len(args) != 8 {
                return fmt.Errorf("Incorrect number of arguments. Expecting 8")
        }
        objectType = "user"
        clientID = args[0]
        userID = args[1]
        bank = args[2]
        bank_account = args[3]
        address = args[4]
        phone = args[5]
        status_open = args[6]
        status_bind = args[7]
        // Create composite key
        user_recordKeyString := objectType  + "_" + clientID + "_" + userID + "_" + bank + "_" + bank_account + "_" + address + "_" + phone + "_" + status_open + "_" + status_bind
        user_recordKey, err := stub.CreateCompositeKey(user_recordKeyString, []string{objectType, clientID, userID, bank, bank_account, address, phone, status_open, status_bind})
        if err != nil {
                fmt.Errorf(err.Error())
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- start updateUserCalled " + user_recordKeyString)

        // Get state with composite key
        user_recordTemp := &User_record{}
        user_recordAsBytes, err := stub.GetState(user_recordKey)
        if err != nil {
                return fmt.Errorf("Failed to get state for :" + user_recordKeyString)
        } else if user_recordAsBytes == nil {
                // We don't need to check if client has joined groupbuy already.
                // But if client doesn't have existing groupbuy, create one
                user_recordTemp = &User_record{
                        ObjectType : objectType,
                        Timestamp : time.Now().String(),
                        ClientID : clientID,
                        UserID : userID,
                        Bank : bank,
                        Bank_account : bank_account,
                        Address : address,
                        Phone : phone,
                        Status_open : status_open,
                        Status_bind : status_bind,
                }
                fmt.Println("User created: " + user_recordKeyString)
                                                                           } else {
                // Unmarshal client and update value
                err = json.Unmarshal(user_recordAsBytes, &user_recordTemp)
                if err != nil {
                        return fmt.Errorf(err.Error())
                }
        }

        user_recordJSONasBytes, err := json.Marshal(user_recordTemp)
        if err != nil {
                return fmt.Errorf(err.Error())
        }

        // Save asset to state
        err = stub.PutState(user_recordKey, user_recordJSONasBytes)
        if err != nil {
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- end Called (success)")
        fmt.Println("==============================")
        return nil


}
// =======================================================
// Called buy other functions to update Transaction. Not an invoke function!!!
//
// updateTransactionCalled - Update Transaction,
// input: TransactionID, Record_type, ClientID, ProductID, Currency, Amount, Opposite.
// =======================================================
func updateTransaction_actionCalled(stub shim.ChaincodeStubInterface, args []string) error {
        var objectType, transactionID, clientID, groupbuyID, status string
        var err error
        if len(args) != 4 {
                return fmt.Errorf("Incorrect number of arguments. Expecting 4")
        }
        objectType = "transaction_action"
        transactionID = args[0]
        clientID = args[1]
        groupbuyID = args[2]
	status = args[3]
        // Create composite key
        transaction_actionKeyString := objectType + "_" +  transactionID  +  "_" + clientID + "_" + groupbuyID + "_" + status
        transaction_actionKey, err := stub.CreateCompositeKey(transaction_actionKeyString, []string{objectType, transactionID, clientID, groupbuyID, status})
        if err != nil {
                fmt.Errorf(err.Error())
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- start updateTransaction_actionCalled " + transaction_actionKeyString)

        // Get state with composite key
        transaction_actionTemp := &Transaction_action{}
        transaction_actionAsBytes, err := stub.GetState(transaction_actionKey)
        if err != nil {
                return fmt.Errorf("Failed to get state for :" + transaction_actionKeyString)
        } else if transaction_actionAsBytes == nil {
                // We don't need to check if client has joined groupbuy already.
                // But if client doesn't have existing groupbuy, create one
                transaction_actionTemp = &Transaction_action{
                        ObjectType : objectType,
                        Timestamp : time.Now().String(),
                        TransactionID : transactionID,
                        ClientID : clientID,
                        GroupbuyID : groupbuyID,
			Status : status,
                }
                fmt.Println("Transaction_action created: " + transaction_actionKeyString)
        } else {
                // Unmarshal client and update value
                err = json.Unmarshal(transaction_actionAsBytes, &transaction_actionTemp)
                if err != nil {
                        return fmt.Errorf(err.Error())
                }
        }

        transaction_actionJSONasBytes, err := json.Marshal(transaction_actionTemp)
        if err != nil {
                return fmt.Errorf(err.Error())
        }

        // Save asset to state
        err = stub.PutState(transaction_actionKey, transaction_actionJSONasBytes)
        if err != nil {
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- end Called (success)")
        fmt.Println("==============================")
        return nil


}
// =======================================================
// Called buy other functions to update Transaction. Not an invoke function!!!
//
// updateTransactionCalled - Update Transaction,
// input: TransactionID, Record_type, ClientID, ProductID, Currency, Amount, Opposite.
// =======================================================
func updateTransaction_action_recordCalled(stub shim.ChaincodeStubInterface, args []string) error {
        var objectType, transactionID, clientID, groupbuyID, status string
        var err error
        if len(args) != 4 {
                return fmt.Errorf("Incorrect number of arguments. Expecting 4")
        }
        objectType = "transaction_action_record"
        transactionID = args[0]
        clientID = args[1]
        groupbuyID = args[2]
        status = args[3]
        // Create composite key
        transaction_action_recordKeyString := objectType + "_" +  transactionID  +  "_" + clientID + "_" + groupbuyID + "_" + status
        transaction_action_recordKey, err := stub.CreateCompositeKey(transaction_action_recordKeyString, []string{objectType, transactionID, clientID, groupbuyID, status})
        if err != nil {
                fmt.Errorf(err.Error())
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- start updateTransaction_action_recordCalled " + transaction_action_recordKeyString)

        // Get state with composite key
        transaction_action_recordTemp := &Transaction_action_record{}
        transaction_action_recordAsBytes, err := stub.GetState(transaction_action_recordKey)
        if err != nil {
                return fmt.Errorf("Failed to get state for :" + transaction_action_recordKeyString)
        } else if transaction_action_recordAsBytes == nil {
                // We don't need to check if client has joined groupbuy already.
                // But if client doesn't have existing groupbuy, create one
                transaction_action_recordTemp = &Transaction_action_record{
                        ObjectType : objectType,
                        Timestamp : time.Now().String(),
                        TransactionID : transactionID,
                        ClientID : clientID,
                        GroupbuyID : groupbuyID,
                        Status : status,
                }
                fmt.Println("Transaction_action_record created: " + transaction_action_recordKeyString)
        } else {
                // Unmarshal client and update value
                err = json.Unmarshal(transaction_action_recordAsBytes, &transaction_action_recordTemp)
                if err != nil {
                        return fmt.Errorf(err.Error())
                }
        }

        transaction_action_recordJSONasBytes, err := json.Marshal(transaction_action_recordTemp)
        if err != nil {
                return fmt.Errorf(err.Error())
        }

        // Save asset to state
        err = stub.PutState(transaction_action_recordKey, transaction_action_recordJSONasBytes)
        if err != nil {
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- end Called (success)")
        fmt.Println("==============================")
        return nil


}
// =======================================================
// Called buy other functions to update Transaction. Not an invoke function!!!
//
// updateTransactionCalled - Update Transaction,
// input: TransactionID, Record_type, ClientID, ProductID, Currency, Amount, Opposite.
// =======================================================
func updateTransaction_takeoverCalled(stub shim.ChaincodeStubInterface, args []string) error {
        var objectType, transactionID, groupbuyID, clientID_sell, clientID_buy, status string
	var err error
        if len(args) != 5 {
                return fmt.Errorf("Incorrect number of arguments. Expecting 5")
        }
        objectType = "transaction_takeover"
	transactionID = args[0]
	groupbuyID = args[1]
	clientID_sell = args[2]
	clientID_buy = args[3]
	status = args[4] 
        // Create composite key
        transaction_takeoverKeyString := objectType + "_" +  transactionID + "_" + groupbuyID + "_" + clientID_sell + "_" + clientID_buy + "_" + status
        transaction_takeoverKey, err := stub.CreateCompositeKey(transaction_takeoverKeyString, []string{objectType, transactionID, groupbuyID, clientID_sell, clientID_buy, status})
        if err != nil {
                fmt.Errorf(err.Error())
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- start updateTransaction_takeoverCalled " + transaction_takeoverKeyString)

        // Get state with composite key
        transaction_takeoverTemp := &Transaction_takeover{}
        transaction_takeoverAsBytes, err := stub.GetState(transaction_takeoverKey)
        if err != nil {
                return fmt.Errorf("Failed to get state for :" + transaction_takeoverKeyString)
        } else if transaction_takeoverAsBytes == nil {
                // We don't need to check if client has joined groupbuy already.
                // But if client doesn't have existing groupbuy, create one
                transaction_takeoverTemp = &Transaction_takeover{
                        ObjectType : objectType,
                        Timestamp : time.Now().String(),
                        TransactionID : transactionID,
                        GroupbuyID : groupbuyID,
                        ClientID_sell : clientID_sell,
                        ClientID_buy : clientID_buy,
			Status : status,
                }
                fmt.Println("Transaction_takeover created: " + transaction_takeoverKeyString)
        } else {
                // Unmarshal client and update value
                err = json.Unmarshal(transaction_takeoverAsBytes, &transaction_takeoverTemp)
                if err != nil {
                        return fmt.Errorf(err.Error())
                }
        }

        transaction_takeoverJSONasBytes, err := json.Marshal(transaction_takeoverTemp)
        if err != nil {
                return fmt.Errorf(err.Error())
        }

        // Save asset to state
        err = stub.PutState(transaction_takeoverKey, transaction_takeoverJSONasBytes)
        if err != nil {
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- end Called (success)")
        fmt.Println("==============================")
        return nil


}
// =======================================================
// Called buy other functions to update Transaction. Not an invoke function!!!
//
// updateTransactionCalled - Update Transaction,
// input: TransactionID, Record_type, ClientID, ProductID, Currency, Amount, Opposite.
// =======================================================
func updateTransaction_takeover_recordCalled(stub shim.ChaincodeStubInterface, args []string) error {
        var objectType, transactionID, groupbuyID, clientID_sell, clientID_buy, status string
        var err error
        if len(args) != 5 {
                return fmt.Errorf("Incorrect number of arguments. Expecting 5")
        }
        objectType = "transaction_takeover_record"
        transactionID = args[0]
        groupbuyID = args[1]
        clientID_sell = args[2]
        clientID_buy = args[3]
        status = args[4]
        // Create composite key
        transaction_takeover_recordKeyString := objectType + "_" +  transactionID + "_" + groupbuyID + "_" + clientID_sell + "_" + clientID_buy + "_" + status
        transaction_takeover_recordKey, err := stub.CreateCompositeKey(transaction_takeover_recordKeyString, []string{objectType, transactionID, groupbuyID, clientID_sell, clientID_buy, status})
        if err != nil {
                fmt.Errorf(err.Error())
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- start updateTransaction_takeover_recordCalled " + transaction_takeover_recordKeyString)

        // Get state with composite key
        transaction_takeover_recordTemp := &Transaction_takeover_record{}
        transaction_takeover_recordAsBytes, err := stub.GetState(transaction_takeover_recordKey)
        if err != nil {
                return fmt.Errorf("Failed to get state for :" + transaction_takeover_recordKeyString)
        } else if transaction_takeover_recordAsBytes == nil {
                // We don't need to check if client has joined groupbuy already.
                // But if client doesn't have existing groupbuy, create one
                transaction_takeover_recordTemp = &Transaction_takeover_record{
                        ObjectType : objectType,
                        Timestamp : time.Now().String(),
                        TransactionID : transactionID,
                        GroupbuyID : groupbuyID,
                        ClientID_sell : clientID_sell,
                        ClientID_buy : clientID_buy,
                        Status : status,
                }
                fmt.Println("Transaction_takeover_record created: " + transaction_takeover_recordKeyString)
        } else {
                // Unmarshal client and update value
                err = json.Unmarshal(transaction_takeover_recordAsBytes, &transaction_takeover_recordTemp)
                if err != nil {
                        return fmt.Errorf(err.Error())
                }
        }

        transaction_takeover_recordJSONasBytes, err := json.Marshal(transaction_takeover_recordTemp)
        if err != nil {
                return fmt.Errorf(err.Error())
        }

        // Save asset to state
        err = stub.PutState(transaction_takeover_recordKey, transaction_takeover_recordJSONasBytes)
        if err != nil {
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- end Called (success)")
        fmt.Println("==============================")
        return nil


}
// =======================================================
// Called buy other functions to update Transaction. Not an invoke function!!!
//
// updateTransactionCalled - Update Transaction,
// input: TransactionID, Record_type, ClientID, ProductID, Currency, Amount, Opposite.
// =======================================================
func updateTransaction_contractCalled(stub shim.ChaincodeStubInterface, args []string) error {
        var objectType, transactionID, clientID, groupbuyID, currency  string
	var amount float64
        var err error
        if len(args) != 6 {
                return fmt.Errorf("Incorrect number of arguments. Expecting 6")
        }
        objectType = "transaction_contract"
        transactionID = args[0]
        clientID = args[1]
        groupbuyID = args[2]
        currency = args[3]
        amount, err = strconv.ParseFloat(args[4], 64)
	status := args[5]
        if err != nil {
                return fmt.Errorf("4th argument must be a numeric string")
        }


        // Create composite key
        transaction_contractKeyString := objectType + "_" +  transactionID  +  "_" + clientID + "_" + groupbuyID + "_" + currency + "_" + status
        transaction_contractKey, err := stub.CreateCompositeKey(transaction_contractKeyString, []string{objectType, transactionID, clientID, groupbuyID, currency, status})
        if err != nil {
                fmt.Errorf(err.Error())
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- start updateTransaction_contractCalled " + transaction_contractKeyString)

        // Get state with composite key
        transaction_contractTemp := &Transaction_contract{}
        transaction_contractAsBytes, err := stub.GetState(transaction_contractKey)
        if err != nil {
                return fmt.Errorf("Failed to get state for :" + transaction_contractKeyString)
        } else if transaction_contractAsBytes == nil {
                // We don't need to check if client has joined groupbuy already.
                // But if client doesn't have existing groupbuy, create one
                transaction_contractTemp = &Transaction_contract{
                        ObjectType : objectType,
                        Timestamp : time.Now().String(),
                        TransactionID : transactionID,
                        ClientID : clientID,
                        GroupbuyID : groupbuyID,
                        Currency : currency,
                        Amount : amount,
			Status : status,
                }
                fmt.Println("Transaction_contract created: " + transaction_contractKeyString)
        } else {
                // Unmarshal client and update value
                err = json.Unmarshal(transaction_contractAsBytes, &transaction_contractTemp)
                if err != nil {
                        return fmt.Errorf(err.Error())
                }
        }

        transaction_contractJSONasBytes, err := json.Marshal(transaction_contractTemp)
        if err != nil {
                return fmt.Errorf(err.Error())
        }

        // Save asset to state
        err = stub.PutState(transaction_contractKey, transaction_contractJSONasBytes)
        if err != nil {
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- end Called (success)")
        fmt.Println("==============================")
        return nil


}
// =======================================================
// Called buy other functions to update Transaction. Not an invoke function!!!
//
// updateTransactionCalled - Update Transaction,
// input: TransactionID, Record_type, ClientID, ProductID, Currency, Amount, Opposite.
// =======================================================
func updateTransaction_contract_recordCalled(stub shim.ChaincodeStubInterface, args []string) error {
        var objectType, transactionID, clientID, groupbuyID, currency  string
        var amount float64
        var err error
        if len(args) != 6 {
                return fmt.Errorf("Incorrect number of arguments. Expecting 6")
        }
        objectType = "transaction_contract_record"
        transactionID = args[0]
        clientID = args[1]
        groupbuyID = args[2]
        currency = args[3]
        amount, err = strconv.ParseFloat(args[4], 64)
        status := args[5]
        if err != nil {
                return fmt.Errorf("4th argument must be a numeric string")
        }


        // Create composite key
        transaction_contract_recordKeyString := objectType + "_" +  transactionID  +  "_" + clientID + "_" + groupbuyID + "_" + currency + "_" + status
        transaction_contract_recordKey, err := stub.CreateCompositeKey(transaction_contract_recordKeyString, []string{objectType, transactionID, clientID, groupbuyID, currency, status})
        if err != nil {
                fmt.Errorf(err.Error())
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- start updateTransaction_contract_recordCalled " + transaction_contract_recordKeyString)

        // Get state with composite key
        transaction_contract_recordTemp := &Transaction_contract_record{}
        transaction_contract_recordAsBytes, err := stub.GetState(transaction_contract_recordKey)
        if err != nil {
                return fmt.Errorf("Failed to get state for :" + transaction_contract_recordKeyString)
        } else if transaction_contract_recordAsBytes == nil {
                // We don't need to check if client has joined groupbuy already.
                // But if client doesn't have existing groupbuy, create one
                transaction_contract_recordTemp = &Transaction_contract_record{
                        ObjectType : objectType,
                        Timestamp : time.Now().String(),
                        TransactionID : transactionID,
                        ClientID : clientID,
                        GroupbuyID : groupbuyID,
                        Currency : currency,
                        Amount : amount,
                        Status : status,
                }
                fmt.Println("Transaction_contract_record created: " + transaction_contract_recordKeyString)
        } else {
                // Unmarshal client and update value
                err = json.Unmarshal(transaction_contract_recordAsBytes, &transaction_contract_recordTemp)
                if err != nil {
                        return fmt.Errorf(err.Error())
                }
        }

        transaction_contract_recordJSONasBytes, err := json.Marshal(transaction_contract_recordTemp)
        if err != nil {
                return fmt.Errorf(err.Error())
        }

        // Save asset to state
        err = stub.PutState(transaction_contract_recordKey, transaction_contract_recordJSONasBytes)
        if err != nil {
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- end Called (success)")
        fmt.Println("==============================")
        return nil


}
// =======================================================
// Called buy other functions to update Transaction. Not an invoke function!!!
//
// updateTransactionCalled - Update Transaction,
// input: TransactionID, Record_type, ClientID, ProductID, Currency, Amount, Opposite.
// =======================================================
func updateTransaction_blockedCalled(stub shim.ChaincodeStubInterface, args []string) error {
        var objectType, transactionID, clientID, groupbuyID, currency  string
        var amount float64
        var err error
        if len(args) != 6 {
                return fmt.Errorf("Incorrect number of arguments. Expecting 6")
        }
        objectType = "transaction_blocked"
        transactionID = args[0]
        clientID = args[1]
        groupbuyID = args[2]
        currency = args[3]
        amount, err = strconv.ParseFloat(args[4], 64)
        status := args[5]
        if err != nil {
                return fmt.Errorf("4th argument must be a numeric string")
        }


        // Create composite key
        transaction_blockedKeyString := objectType + "_" +  transactionID  +  "_" + clientID + "_" + groupbuyID + "_" + currency + "_" + status
        transaction_blockedKey, err := stub.CreateCompositeKey(transaction_blockedKeyString, []string{objectType, transactionID, clientID, groupbuyID, currency, status})
        if err != nil {
                fmt.Errorf(err.Error())
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- start updateTransaction_blockedCalled " + transaction_blockedKeyString)

        // Get state with composite key
        transaction_blockedTemp := &Transaction_blocked{}
        transaction_blockedAsBytes, err := stub.GetState(transaction_blockedKey)
        if err != nil {
                return fmt.Errorf("Failed to get state for :" + transaction_blockedKeyString)
        } else if transaction_blockedAsBytes == nil {
                // We don't need to check if client has joined groupbuy already.
                // But if client doesn't have existing groupbuy, create one
                transaction_blockedTemp = &Transaction_blocked{
                        ObjectType : objectType,
                        Timestamp : time.Now().String(),
                        TransactionID : transactionID,
                        ClientID : clientID,
                        GroupbuyID : groupbuyID,
                        Currency : currency,
                        Amount : amount,
                        Status : status,
                }
                fmt.Println("Transaction_blocked created: " + transaction_blockedKeyString)
        } else {
                // Unmarshal client and update value
                err = json.Unmarshal(transaction_blockedAsBytes, &transaction_blockedTemp)
                if err != nil {
                        return fmt.Errorf(err.Error())
                }
        }

        transaction_blockedJSONasBytes, err := json.Marshal(transaction_blockedTemp)
        if err != nil {
                return fmt.Errorf(err.Error())
        }

        // Save asset to state
        err = stub.PutState(transaction_blockedKey, transaction_blockedJSONasBytes)
        if err != nil {
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- end Called (success)")
        fmt.Println("==============================")
        return nil


}
// =======================================================
// Called buy other functions to update Transaction. Not an invoke function!!!
//
// updateTransactionCalled - Update Transaction,
// input: TransactionID, Record_type, ClientID, ProductID, Currency, Amount, Opposite.
// =======================================================
func updateTransaction_blocked_recordCalled(stub shim.ChaincodeStubInterface, args []string) error {
        var objectType, transactionID, clientID, groupbuyID, currency  string
        var amount float64
        var err error
        if len(args) != 6 {
                return fmt.Errorf("Incorrect number of arguments. Expecting 6")
        }
        objectType = "transaction_blocked_record"
        transactionID = args[0]
        clientID = args[1]
        groupbuyID = args[2]
        currency = args[3]
        amount, err = strconv.ParseFloat(args[4], 64)
        status := args[5]
        if err != nil {
                return fmt.Errorf("4th argument must be a numeric string")
        }


        // Create composite key
        transaction_blocked_recordKeyString := objectType + "_" +  transactionID  +  "_" + clientID + "_" + groupbuyID + "_" + currency + "_" + status
        transaction_blocked_recordKey, err := stub.CreateCompositeKey(transaction_blocked_recordKeyString, []string{objectType, transactionID, clientID, groupbuyID, currency, status})
        if err != nil {
                fmt.Errorf(err.Error())
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- start updateTransaction_blocked_recordCalled " + transaction_blocked_recordKeyString)

        // Get state with composite key
        transaction_blocked_recordTemp := &Transaction_blocked_record{}
        transaction_blocked_recordAsBytes, err := stub.GetState(transaction_blocked_recordKey)
        if err != nil {
                return fmt.Errorf("Failed to get state for :" + transaction_blocked_recordKeyString)
        } else if transaction_blocked_recordAsBytes == nil {
                // We don't need to check if client has joined groupbuy already.
                // But if client doesn't have existing groupbuy, create one
                transaction_blocked_recordTemp = &Transaction_blocked_record{
                        ObjectType : objectType,
                        Timestamp : time.Now().String(),
                        TransactionID : transactionID,
                        ClientID : clientID,
                        GroupbuyID : groupbuyID,
                        Currency : currency,
                        Amount : amount,
                        Status : status,
                }
                fmt.Println("Transaction_blocked_record created: " + transaction_blocked_recordKeyString)
        } else {
                // Unmarshal client and update value
                err = json.Unmarshal(transaction_blocked_recordAsBytes, &transaction_blocked_recordTemp)
                if err != nil {
                        return fmt.Errorf(err.Error())
                }
        }

        transaction_blocked_recordJSONasBytes, err := json.Marshal(transaction_blocked_recordTemp)
        if err != nil {
                return fmt.Errorf(err.Error())
        }

        // Save asset to state
        err = stub.PutState(transaction_blocked_recordKey, transaction_blocked_recordJSONasBytes)
        if err != nil {
                return fmt.Errorf(err.Error())
        }
        fmt.Println("- end Called (success)")
        fmt.Println("==============================")
        return nil


}

