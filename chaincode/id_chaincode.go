package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type IdChaincode struct {

}
type ID struct {
	Name string
	SchoolId string
	Sex string
	Age int
}
//'{"Args":["init","a","100","b","200"]}'
func (i *IdChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	//参数A是第一个值， Aval是A的值
	key := args[0]
	value := []byte(args[1])
	//检验能否转化为ID结构体
	var znh ID
	err := json.Unmarshal(value, &znh)
	if err != nil {
		return shim.Error("json.Unmarshal err,ID struct false.")
	}
	fmt.Println("ID struct true.")
	//验证成功后存入区块链
	err = stub.PutState(key,value)
	if err != nil {
		return shim.Error("stub.PutState err")
	}
	return shim.Success([]byte("success"))
}

func (i *IdChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "addId" {
		return i.addId(stub, args)
	}else if function == "queryId" {
		return i.queryId(stub, args)
	}else if function == "deleteId" {
		return i.deleteId(stub, args)
	}else if function == "queryAllIds" {
		return i.queryAllIds(stub)
	}
	return shim.Error("operation failed.")
}

func main() {
	err := shim.Start(new(IdChaincode))
	if err != nil {
		fmt.Printf("Error starting IdChaincode:%S", err)
	}
}
//增加ID
//	'{"Args":["addId","znh","{'Name':'znh','SchoolId':'123','Sex':'男','Age':22”}"]}'
//                      0                               1
func (i* IdChaincode) addId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return  shim.Error("Incorrect number of arguments. Expecting name of the person to query.")
	}
	var newId = args[0]
	var jsonstr = []byte(args[1])
	//检验能否转化为ID结构体
	var name ID
	err := json.Unmarshal(jsonstr, &name)
	if err != nil {
		return shim.Error("json.Unmarshal err,ID struct false.")
	}
	fmt.Println("ID struct true.")
	//验证成功后存入区块链
	err = stub.PutState(newId,jsonstr)
	if err != nil {
		return shim.Error("stub.PutState err")
	}
	return shim.Success(jsonstr)
}
//查询ID
func (i* IdChaincode) queryId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return  shim.Error("Incorrect number of arguments. Expecting name of ")
	}
	key := args[0]
	value, err := stub.GetState(key)
	if err != nil {
		return shim.Error("Fail to get state for：" + key)
	}
	if value == nil{
		return shim.Error("Fail to get state for" + key)
	}
	jsonResp := "{\"Name\":\"" + key + "\"value:\"" + string(value) + "}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(value)
}
//查询
func (s *IdChaincode) queryAllIds(APIstub shim.ChaincodeStubInterface) pb.Response {

	startKey := ""
	endKey := ""

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

	fmt.Printf("- queryAllCars:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
//删除ID
func (i* IdChaincode) deleteId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return  shim.Error("Incorrect number of arguments. Expecting 1")
	}
	key := args[0]
	//删除key
	err := stub.DelState(key)
	if err != nil {
		return shim.Error("Failed to delete state.")
	}
	return shim.Success([]byte("deleteId success."))
}