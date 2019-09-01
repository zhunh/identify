package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//学生结构体
type student struct {
	StudentID string
	StudentName string
	StudentPhone string
	StudentAddress string
	StudentStatus int
	SchoolName string
}

type StudentChaincode struct {

}

func main() {
	err := shim.Start(new(StudentChaincode))
	if err != nil {
		fmt.Printf("Error starting StudentChaincode:%S", err)
	}
}

func (i *StudentChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "CreateStudent" {
		return i.CreateStudent(stub, args)
	}else if function == "GetStudentById" {
		return i.GetStudentById(stub, args)
	}else if function =="queryAllStudent" {
		return i.queryAllStudent(stub, args)
	}
	return shim.Error("operation failed.")
}

func (i *StudentChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	st01 := `{"StudentID":"6120170105","StudentName":"znh","StudentPhone":"18146688393","StudentAddress":"Street89,room3","StudentStatus":1,"SchoolName":"JXLG"}`
	_ = stub.PutState("FirstStudentData", []byte(st01))
	return shim.Success([]byte("init success"))
}

//创建学生
func (i *StudentChaincode) CreateStudent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//1.参数检验
	if len(args) != 2 {
		return  shim.Error("Incorrect number of argument. Expecting 2")
	}
	//2.取参
	key := args[0]
	value := []byte(args[1])
	//3.检查合法性
	var student student
	err := json.Unmarshal(value,&student)
	if err != nil{
		return shim.Error("json.Unmarsh err. student struct err.")
	}
	//4.存入区块链
	err = stub.PutState(key, value)
	if err != nil{
		return shim.Error("stub.PutState err.")
	}
	return shim.Success([]byte("CreateStudent success"))
}
//查找学生
func (i *StudentChaincode) GetStudentById(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
func (i *StudentChaincode) queryAllStudent(stub shim.ChaincodeStubInterface, args []string) pb.Response {

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