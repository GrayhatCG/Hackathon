// PHRTCR
package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var regionIndexTxStr = "_regionIndexTxStr"

type PatientRecord struct {
	PATIENT_ID                   string `json:"PATIENT_ID"`
	PATIENT_NAME                 string `json:"PATIENT_NAME"`
	PATIENT_DOB                  string `json:"PATIENT_DOB"`
	INSUR_ID                     string `json:"INSUR_ID"`
	INSUR_NAME                   string `json:"INSUR_NAME"`
	PATIENT_INSUR_STARTDATE      string `json:"PATIENT_INSUR_STARTDATE"`
	PATIENT_INSUR_ENDDATE        string `json:"PATIENT_INSUR_ENDDATE"`
	PATIENT_INSUR_TOTALSUM       string `json:"PATIENT_INSUR_TOTALSUM"`
	PATIENT_INSUR_ELIGIBLESUM    string `json:"PATIENT_INSUR_ELIGIBLESUM"`
	CLAIM_ID                     string `json:"CLAIM_ID"`
	HELTH_RECORD_DRNAME          string `json:"HELTH_RECORD_DRNAME"`
	HELTH_RECORD_HOSPITALID      string `json:"HELTH_RECORD_HOSPITALID"`
	HELTH_RECORD_TREATMENTDETAIL string `json:"HELTH_RECORD_TREATMENTDETAIL"`
	DC_ID                        string `json:"DC_ID"`
	HELTH_RECORD_DCREQ           string `json:"HELTH_RECORD_DICREQ"`
	HELTH_RECORD_DCDETAIL        string `json:"HELTH_RECORD_DCDETAIL"`
	HELTH_RECORD_FEE             string `json:"HELTH_RECORD_FEE"`
}

func (t *PatientRecord) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error
	// Initialize the chaincode

	fmt.Printf("Deployment of Patient’s Health Records Tracking & Claim Processing is completed\n")

	var emptyPolicyTxs []PatientRecord
	jsonAsBytes, _ := json.Marshal(emptyPolicyTxs)
	err = stub.PutState(regionIndexTxStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Add region data for the policy
func (t *PatientRecord) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == regionIndexTxStr {
		return t.RegisterInsurPolicy(stub, args)
	}
	return nil, nil
}

func (t *PatientRecord) RegisterInsurPolicy(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var RegionDataObj PatientRecord
	var RegionDataList []PatientRecord
	var err error

	if len(args) != 15 {
		return nil, errors.New("Incorrect number of arguments. Need 14 arguments")
	}

	// Initialize the chaincode
	RegionDataObj.PATIENT_ID = args[0]
	RegionDataObj.PATIENT_NAME = args[1]
	RegionDataObj.PATIENT_DOB = args[2]
	RegionDataObj.INSUR_ID = args[3]
	RegionDataObj.INSUR_NAME = args[4]
	RegionDataObj.PATIENT_INSUR_STARTDATE = args[5]
	RegionDataObj.PATIENT_INSUR_ENDDATE = args[6]
	RegionDataObj.PATIENT_INSUR_TOTALSUM = args[7]
	RegionDataObj.PATIENT_INSUR_ELIGIBLESUM = args[8]

	fmt.Printf("Input from user:%s\n", RegionDataObj)

	regionTxsAsBytes, err := stub.GetState(regionIndexTxStr)
	if err != nil {
		return nil, errors.New("Failed to get consumer Transactions")
	}
	json.Unmarshal(regionTxsAsBytes, &RegionDataList)

	RegionDataList = append(RegionDataList, RegionDataObj)
	jsonAsBytes, _ := json.Marshal(RegionDataList)

	err = stub.PutState(regionIndexTxStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
// Query callback representing the query of a chaincode
func (t *PatientRecord) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var PATIENT_ID string // Entities
	var err error
	var resAsBytes []byte

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	PATIENT_ID = args[0]

	resAsBytes, err = t.GetPatientDetails(stub, PATIENT_ID)

	fmt.Printf("Query Response:%s\n", resAsBytes)

	if err != nil {
		return nil, err
	}

	return resAsBytes, nil
}

func (t *PatientRecord) GetPatientDetails(stub shim.ChaincodeStubInterface, PATIENTID string) ([]byte, error) {

	fmt.Println("PATIENTID " + PATIENTID)
	//var requiredObj PatientRecord
	var objFound bool
	PatientAsBytes, err := stub.GetState(regionIndexTxStr)
	if err != nil {
		return nil, errors.New("Failed to get Patient Details")
	}
	var PatientObjects []PatientRecord
	//var PatientObjects1 []PatientRecord
	json.Unmarshal(PatientAsBytes, &PatientObjects)
	length := len(PatientObjects)
	fmt.Printf("Output from chaincode: %s\n", PatientAsBytes)

	if PATIENTID == "" {
		res, err := json.Marshal(PatientObjects)
		if err != nil {
			return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	}
fmt.Println("Hello World!1")
	objFound = false
	// iterate
	for i := 0; i < length; i++ {
		fmt.Println("Hello World!2")
		obj := PatientObjects[i]
		//fmt.Println("#### " +i+ " obj" + obj)
		if PATIENTID == obj.PATIENT_ID {
			PatientObjects = append(PatientObjects, obj)
			//requiredObj = obj
			objFound = true
		}
		fmt.Println("Hello World!3")
	}

	if objFound {
		res, err := json.Marshal(PatientObjects)
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


func main() {
	fmt.Println("Hello World!")

	err := shim.Start(new(PatientRecord))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

