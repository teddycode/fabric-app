package main

import (
	"bytes"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"sync"
	"time"
)

// 定义智能合约结构体
type SmartContract struct {
}

// 传感器数据
type Sensor struct {
	PointID string `json:"point_id"` // 采集点ID
	TypeID  string `json:"type_id"`  // 指标类型ID
	Value   string `json:"value"`    // 指标数值
	Unit    string `json:"unit"`     // 单位
}

// 图片数据格式
type Picture struct {
	PointID string `json:"point_id"` // 采集点ID
	Type    string `json:"type"`     // 类型
	FName   string `json:"f_name"`   // 文件名
	Size    string `json:"size"`     //	大小
}

// 农事记录数据格式
type Farm struct {
	Name     string `json:"name"`     // 操作名称
	Time     string `json:"time"`     // 操作时间
	Behavior string `json:"behavior"` // 操作行为
	Info     string `json:"info"`     // 具体信息
}

var Cnt = 0
var Lock sync.Mutex

// 全局计数器，避免高并发引起主键冲突
func getCounter() string {
	Lock.Lock()
	defer Lock.Unlock()
	Cnt++
	if Cnt > 16 {
		Cnt = 0
	}
	return strconv.FormatInt(int64(Cnt), 16)
}

// 时间戳转换
func timeStampToUnixNanoStr(sec, nano int64) string {
	return strconv.FormatInt(time.Unix(sec, nano).UnixNano(), 10)
}

// 初始化合约
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success([]byte("success"))
}

// 合约逻辑调用处理函数
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	// 检索合约函数
	function, args := APIstub.GetFunctionAndParameters()
	if function == "add" { // 添加记录
		return s.add(APIstub, args)
	} else if function == "query" { // 查询记录
		return s.query(APIstub, args)
	} else if function == "queryOne" { // 查询记录
		return s.queryOne(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

// 传感器example：{"Args":["add","s","point01",\"{...}\"]}'
// 图片example：{"Args":["add","p","point01",\"{...}\"]}'
// 农事example：{"Args":["add","f","user1",\"{...}\"]}'
// 以采集点&时间戳为主键，保存传感器数据、图片数据、农事管理数据。
func (s *SmartContract) add(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	ts, err := APIstub.GetTxTimestamp() // 获取交易的时间戳
	if err != nil {
		return shim.Error(err.Error())
	}
	unixNano := timeStampToUnixNanoStr(ts.Seconds, int64(ts.Nanos))
	index := args[0] + "~" + args[1] + "~" + unixNano
	//fmt.Printf("index:%s\n", index)
	APIstub.PutState(index, []byte(args[2]))
	return shim.Success(nil)
}

// 查询指定key的值
func (s *SmartContract) queryOne(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Println("args: ", args)
	byte1, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	//fmt.Println("query results: ", string(byte1))
	return shim.Success(byte1)
}

// 传感器example:'{"Args":["query","s","point01"，"0","time2"]}'
// 按采集点和时间范围查询记录信息，输入采集点、起始时间戳（十六进制）
func (s *SmartContract) query(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	startKey := args[0] + "~" + args[1] + "~" + args[2]
	endKey := args[0] + "~" + args[1] + "~" + args[3]

	fmt.Printf("startKey:%s\nendKey:%s\n", startKey, endKey)

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

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
		buffer.WriteString("{\"K\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"V\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("queryDatas: %s\n", buffer.String())
	return shim.Success(buffer.Bytes())

}

func main() {
	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
