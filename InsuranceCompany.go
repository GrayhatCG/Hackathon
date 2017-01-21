/*
Copyright Capgemini India. 2016 All Rights Reserved.
*/

package main

import (
	"errors"
	"fmt"
	//"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/golang/protobuf/ptypes/timestamp"
)

var insuranceIndexKey = "_insuranceIndexKey"

type InsuranceCompany struct {
	insuranceCompanyID string `json:"INSURANCE_COMPANY_ID"`
	insuranceName      string `json:"INSURANCE_NAME"`
}

type InsuranceDetails struct {
	insuranceID        string `json:"INSURANCE_ID"`
	totalSumInsured    string `json:"TOTAL_SUM_INSURED"`
	startDate          string `json:"START_DATE"`	
	endDate            string `json:"END_DATE"`
	eligibleSumInsured string `json:"ELIGIBLE_SUM_INSURED"`
	//insuranceCompany1   InsuranceCompany `json:"INSURANCE_COMPANY"`
}


func main() {
	err := shim.Start(new(InsuranceDetails))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

func (t *InsuranceDetails) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error
	// Initialize the chaincode
	fmt.Printf("Deployment of InsuranceDetails :completed\n")

	var insuranceDetails []InsuranceDetails
	jsonAsBytes, _ := json.Marshal(insuranceDetails)
	err = stub.PutState(insuranceIndexKey, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Add region data for the policy
func (t *InsuranceDetails) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == insuranceIndexKey {
		return t.RegisterPolicy(stub, args)
	}
	return nil, nil
}

func (t *InsuranceDetails)  RegisterPolicy(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var InsuranceDetailsObj InsuranceDetails
	var InsuranceDetailsList []InsuranceDetails
	var err error

	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Need 6 arguments")
	}

	// Initialize the chaincode
	//InsuranceDetailsObj.INSURANCE_ID = args[0]
	//InsuranceDetailsObj.TOTAL_SUM_INSURED = args[1]
	//InsuranceDetailsObj.START_DATE = args[2]
//	InsuranceDetailsObj.END_DATE = args[3]
//	InsuranceDetailsObj.ELIGIBLE_SUM_INSURED = args[4]
	//InsuranceDetailsObj.INSURANCE_COMPANY = args[5]

	fmt.Printf("Input from user:%s\n", InsuranceDetailsObj)

	insuranceDetailsAsBytes, err := stub.GetState(insuranceIndexKey)
	if err != nil {
		return nil, errors.New("Failed to get Insurance Transactions")
	}
	json.Unmarshal(insuranceDetailsAsBytes, &InsuranceDetailsList)

	InsuranceDetailsList = append(InsuranceDetailsList, InsuranceDetailsObj)
	jsonAsBytes, _ := json.Marshal(InsuranceDetailsList)

	err = stub.PutState(insuranceIndexKey, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *InsuranceDetails) Query(stub shim.ChaincodeStubInterface,function string, args []string) ([]byte, error) {

	var insuranceID string // Entities
	var err error
	var resAsBytes []byte

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the policy ID to query")
	}

	insuranceID = args[0]

	resAsBytes, err = t.GetPolicyDetails(stub, insuranceID)
	fmt.Printf("Query Response:%s\n", resAsBytes)
	if err != nil {
		return nil, err
	}

	return resAsBytes, nil
}

func (t *InsuranceDetails)  GetPolicyDetails(stub shim.ChaincodeStubInterface, insuranceID string) ([]byte, error) {

	//var requiredObj RegionData
	var objFound bool
	InsurancePolicyAsBytes, err := stub.GetState(insuranceIndexKey)
	if err != nil {
		return nil, errors.New("Failed to get Merchant Transactions")
	}
	var InsurancePolicyObjects []InsuranceDetails
	var InsurancePolicyObjects1 []InsuranceDetails
	json.Unmarshal(InsurancePolicyAsBytes, &InsurancePolicyObjects)
	length := len(InsurancePolicyObjects)
	fmt.Printf("Output from chaincode: %s\n", InsurancePolicyAsBytes)

	if insuranceID == "" {
		res, err := json.Marshal(InsurancePolicyObjects)
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	}

	objFound = false
	// iterate
	for i := 0; i < length; i++ {
		obj := InsurancePolicyObjects[i]
		if insuranceID == obj.insuranceID {
			InsurancePolicyObjects1 = append(InsurancePolicyObjects1,obj)
			//requiredObj = obj
			objFound = true
		}
	}

	if objFound {
		res, err := json.Marshal(InsurancePolicyObjects1)
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	} else {
		res, err := json.Marshal("No Data found")
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	}
}


