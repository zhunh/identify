package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//事迹结构体
type record struct {
	SidTs string
	StudentID string
	Timestamp string
	Note []string
}

type RecordChaincode struct {

}

func main() {
	err := shim.Start(new(RecordChaincode))
	if err != nil {
		fmt.Printf("Error starting StudentChaincode:%S", err)
	}
}
//调用
func (i *RecordChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "CreateRecord" {
		return i.CreateRecord(stub, args)
	}else if function =="queryAllRecord" {
		return i.queryAllRecord(stub, args)
	}
	return shim.Error("operation failed.")
}
//初始化
func (i *RecordChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	rc := `{"SidTs":"110","StudentID":"6120170105","Timestamp":"71285945","Note":["apple","banana","orange"]}`
	_ = stub.PutState("FirstRecordData", []byte(rc))
	return  shim.Success([]byte("init suc."))
}

//创建记录
func (i *RecordChaincode) CreateRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//1.参数检验
	if len(args) != 2 {
		return  shim.Error("Incorrect number of argument. Expecting 2")
	}
	//2.取参
	key := args[0]
	value := []byte(args[1])
	//3.检查合法性
	var record record
	err := json.Unmarshal(value,&record)
	if err != nil{
		return shim.Error("json.Unmarsh err. record struct err.")
	}
	//4.存入区块链
	err = stub.PutState(key, value)
	if err != nil{
		return shim.Error("stub.PutState err.")
	}
	return shim.Success([]byte("CreateRecord success"))
}
//查询
func (i *RecordChaincode) queryAllRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {

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