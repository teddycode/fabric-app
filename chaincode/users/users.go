package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strings"
)

/*
	user manage modular
	// admin add user info
	// user login must check
*/

type Auth struct {
}

func (this *Auth) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("Chaincode init success!")
	return shim.Success(nil)
}

func (this *Auth) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, parameters := stub.GetFunctionAndParameters()

	if function == "check" {
		return this.check(stub, parameters)
	} else if function == "add" {
		return this.add(stub, parameters)
	} else if function == "del" {
		return this.del(stub, parameters)
	}

	return shim.Error("Invalid Smart Contract function name")
}

// params: user real name, identity number hash
func (this *Auth) add(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.Expecting 2")
	}

	name := args[0]
	id := args[1]

	err := stub.PutState(id, []byte(name))
	if err != nil {
		shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// params: user real name, identity number hash
func (this *Auth) check(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.Expecting 2")
	}

	name := args[0]
	id := args[1]

	//  query by user name
	data, err := stub.GetState(id)

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("data:%v,name:%v", data, name)
	if data != nil && strings.Compare(string(data), name) == 0 {
		return shim.Success(nil)
	}

	return shim.Error("User not found.")
}

// key para: identity id
func (this *Auth) del(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.Expecting 1")
	}
	id := args[0]
	data, err := stub.GetState(id)
	if err != nil {
		return shim.Error(err.Error())
	}

	if data == nil {
		return shim.Error("User not found.")
	}

	err = stub.DelState(id)
	if err != nil {
		shim.Error(err.Error())
	}

	return shim.Success([]byte("User unAuth ok!"))

}

func main() {
	err := shim.Start(new(Auth))
	if err != nil {
		fmt.Println("chaincode start error!")
	}
}
