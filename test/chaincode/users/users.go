package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

/*
	user manage modular
	// admin add user info
	// user login must check
	//  name~id = role
*/
const (
	ROLE_ADMIN   = "0"
	ROLE_USER    = "1"
	ROLE_REVOKED = "2"
)

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

// params: user real name, role ,identity number hash
// add a user identity
func (this *Auth) add(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments.Expecting 3")
	}

	name := args[0]
	role := args[1]
	id := args[2]
	key := name + "~" + id

	err := stub.PutState(key, []byte(role))
	if err != nil {
		shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// params: user real name, identity number hash
// check if exist or is admin user
func (this *Auth) check(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.Expecting 2")
	}

	name := args[0]
	id := args[1]
	key := name + "~" + id

	//  query by user name
	data, err := stub.GetState(key)

	if err != nil {
		return shim.Error(err.Error())
	}
	if data == nil {
		return shim.Success([]byte(ROLE_REVOKED))
	}
	return shim.Success(data)
}

// key para: user name , identity id
func (this *Auth) del(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.Expecting 1")
	}
	name := args[0]
	id := args[1]
	key := name + "~" + id
	data, err := stub.GetState(key)
	if err != nil {
		return shim.Error(err.Error())
	}

	if data == nil {
		return shim.Error("User not found.")
	}

	err = stub.PutState(key, []byte(ROLE_REVOKED))
	if err != nil {
		shim.Error(err.Error())
	}
	return shim.Success([]byte("User unAuth ok!"))
}

func main() {
	err := shim.Start(new(Auth))
	if err != nil {
		fmt.Println("Chaincode start error!")
	}
}
