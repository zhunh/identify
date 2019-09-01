package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)
//学校结构体
type school struct {
	SchoolName string
	SchoolLocation string
	SchoolPhone string
}

type SchoolChaincode struct {

}

func main() {
	err := shim.Start(new(SchoolChaincode))
	if err != nil {
		fmt.Printf("Error starting AccChaincode:%S", err)
	}
}

func (i *SchoolChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "CreateSchool" {
		return i.CreateSchool(stub, args)
	}else if function == "GetSchoolByName" {
		return i.GetSchoolByName(stub, args)
	}else if function =="queryAllSchool" {
		return i.queryAllSchool(stub, args)
	}
	return shim.Error("operation failed.")
}

func (i *SchoolChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	sc := `{"SchoolName":"JXLG","SchoolLocation":"GanZhou Street176","SchoolPhone":"630078"}`
	_ = stub.PutState("FirstSchoolData", []byte(sc))
	return shim.Success([]byte("success"))
}
//创建学校
func (i *SchoolChaincode) CreateSchool(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//1.参数检验
	if len(args) != 2 {
		return  shim.Error("Incorrect number of argument. Expecting 2")
	}
	//2.取参
	key := args[0]
	value := []byte(args[1])
	//3.检查合法性
	var school school
	err := json.Unmarshal(value,&school)
	if err != nil{
		return shim.Error("json.Unmarsh err. school struct err.")
	}
	//4.存入区块链
	err = stub.PutState(key, value)
	if err != nil{
		return shim.Error("stub.PutState err.")
	}
	return shim.Success([]byte("CreateSchool success"))
}
//查找学校
func (i *SchoolChaincode) GetSchoolByName(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return  shim.Error("Incorrect number of arguments. Expecting 1")
	}
	key := args[0]
	value, err := stub.GetState(key)
	if err != nil {
		return shim.Error("Fail to get state for：" + key)
	}
	if value == nil{
		return shim.Error("Fail to get state for" + key)
	}
	return shim.Success(value)
}
//查询
func (i *SchoolChaincode) queryAllSchool(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	startKey := ""
	endKey := ""

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
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

	fmt.Printf("- queryAllCars:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}