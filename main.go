package main

import (
	"fmt"
	cli "github.com/fabric-app/controller"
	"log"
)

//const (
//	org1CfgPath = "./config/conn-fn1.yaml"
//	org2CfgPath = "./config/conn-fn2.yaml"
//)
const (
	org1CfgPath = "./conf-local/conn-fn1.yaml"
	org2CfgPath = "./conf-local/conn-fn2.yaml"
)

var (
	peer0Org1 = "peer0.org1.lzawt.com"
	peer0Org2 = "peer0.org2.lzawt.com"
)


func main() {
	org1Client := cli.New(org1CfgPath, "Org1", "mychannel","users","Admin", "User1")
	//org2Client := cli.New(org2CfgPath, "Org2",  "mychannel","users","Admin", "User1")

	defer org1Client.Close()
	//defer org2Client.Close()

	// Install, instantiate, invoke, query
	Phase1(org1Client)
	// Install, upgrade, invoke, query
	//Phase2(org1Client, org2Client)
}

func Phase1(cli1 *cli.Client) {
	var args []string
	log.Println("=================== Phase 1 begin ===================")
	defer log.Println("=================== Phase 1 end ===================")

	// register users
	str,ok := cli1.RegisterUser("teddy","org1","123456","user")
	fmt.Println(str)

	str,ok = cli1.EnrollUser("teddy","org1","123456","user")
	if !ok {
		log.Panicf(str)
	}
	fmt.Println(str)

	if err := cli1.InstallCC("v1", peer0Org1); err != nil {
		log.Panicf("Intall chaincode error: %v", err)
	}
	log.Println("Chaincode has been installed on org1's peers")

	// InstantiateCC chaincode only need once for each channel
	args = []string{""}
	if _, err := cli1.InstantiateCC(args,"v1", peer0Org1); err != nil {
		log.Panicf("Instantiated chaincode error: %v", err)
	}
	log.Println("Chaincode has been instantiated")

	// add
	args =[]string{"teddy", "123456"}
	if _, err := cli1.InvokeCC("add",args,[]string{peer0Org1}); err != nil {
		log.Panicf("Invoke chaincode error: %v", err)
	}
	log.Println("Invoke chaincode add success")

	// query check ture
	args = []string{"teddy", "123456"}
	if err := cli1.QueryCC("check",args,"peer0.org1.lzawt.com"); err != nil {
		log.Panicf("Query chaincode error: %v", err)
	}
	log.Println("Query chaincode check success on peer0.org1")

	// query check error
	args = []string{"teddy1", "123456"}
	if err := cli1.QueryCC("check",args,"peer0.org1.lzawt.com"); err != nil {
		log.Panicf("Query chaincode error: %v", err)
	}
	log.Println("Query chaincode success on peer0.org1")


	// invoke del
	args = []string{"123456"}
	if _,err := cli1.InvokeCC("del",args,[]string{peer0Org1}); err != nil {
		log.Panicf("Query chaincode error: %v", err)
	}
	log.Println("Invoke chaincode success on peer0.org1.")

	// query check error
	args = []string{"teddy", "123456"}
	if err := cli1.QueryCC("check",args,"peer0.org1.lzawt.com"); err != nil {
		log.Panicf("Query chaincode error: %v", err)
	}
	log.Println("Query chaincode success on peer0.org1")

}


//func Phase1(cli1, cli2 *cli.Client) {
//	var args []string
//	log.Println("=================== Phase 1 begin ===================")
//	defer log.Println("=================== Phase 1 end ===================")
//
//	if err := cli1.InstallCC("v1", peer0Org1); err != nil {
//		log.Panicf("Intall chaincode error: %v", err)
//	}
//	log.Println("Chaincode has been installed on org1's peers")
//
//	if err := cli2.InstallCC("v1", peer0Org2); err != nil {
//		log.Panicf("Intall chaincode error: %v", err)
//	}
//	log.Println("Chaincode has been installed on org2's peers")
//
//	// InstantiateCC chaincode only need once for each channel
//	args = []string{""}
//	if _, err := cli1.InstantiateCC(args,"v1", peer0Org1); err != nil {
//		log.Panicf("Instantiated chaincode error: %v", err)
//	}
//	log.Println("Chaincode has been instantiated")
//
//	// add
//	args =[]string{"teddy", "123456"}
//	if _, err := cli1.InvokeCC("add",args,[]string{peer0Org1}); err != nil {
//		log.Panicf("Invoke chaincode error: %v", err)
//	}
//	log.Println("Invoke chaincode add success")
//
//	// query check ture
//	args = []string{"teddy", "123456"}
//	if err := cli1.QueryCC("check",args,"peer0.org1.lzawt.com"); err != nil {
//		log.Panicf("Query chaincode error: %v", err)
//	}
//	log.Println("Query chaincode check success on peer0.org1")
//
//	// query check error
//	args = []string{"teddy1", "123456"}
//	if err := cli1.QueryCC("check",args,"peer0.org1.lzawt.com"); err != nil {
//		log.Panicf("Query chaincode error: %v", err)
//	}
//	log.Println("Query chaincode success on peer0.org1")
//
//
//	// invoke del
//	args = []string{"123456"}
//	if _,err := cli1.InvokeCC("del",args,[]string{peer0Org1}); err != nil {
//		log.Panicf("Query chaincode error: %v", err)
//	}
//	log.Println("Invoke chaincode success on peer0.org1.")
//
//	// query check error
//	args = []string{"teddy", "123456"}
//	if err := cli1.QueryCC("check",args,"peer0.org1.lzawt.com"); err != nil {
//		log.Panicf("Query chaincode error: %v", err)
//	}
//	log.Println("Query chaincode success on peer0.org1")
//
//}

//func Phase2(cli1, cli2 *cli.Client) {
//	log.Println("=================== Phase 2 begin ===================")
//	defer log.Println("=================== Phase 2 end ===================")
//
//	v := "v2"
//
//	// Install new version chaincode
//	if err := cli1.InstallCC(v, peer0Org1); err != nil {
//		log.Panicf("Intall chaincode error: %v", err)
//	}
//	log.Println("Chaincode has been installed on org1's peers")
//
//	if err := cli2.InstallCC(v, peer0Org2); err != nil {
//		log.Panicf("Intall chaincode error: %v", err)
//	}
//	log.Println("Chaincode has been installed on org2's peers")
//
//	// Upgrade chaincode only need once for each channel
//	if err := cli1.UpgradeCC(v, peer0Org1); err != nil {
//		log.Panicf("Upgrade chaincode error: %v", err)
//	}
//	log.Println("Upgrade chaincode success for channel")
//
//	args :=[]string{"teddy", "123456", "10"}
//	if _, err := cli1.InvokeCC([]string{"peer0.org1.lzawt.com","peer0.org2.lzawt.com"},args); err != nil {
//		log.Panicf("Invoke chaincode error: %v", err)
//	}
//	log.Println("Invoke chaincode success")
//
//	if err := cli1.QueryCC("peer0.org2.lzawt.com", "a"); err != nil {
//		log.Panicf("Query chaincode error: %v", err)
//	}
//	log.Println("Query chaincode success on peer0.org2")
//}
