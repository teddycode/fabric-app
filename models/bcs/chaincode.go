package bcs

import (
	"log"
	"net/http"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/pkg/errors"
)

// InstallCC install chaincode for target peer
func (c *Client) InstallCC(cc, path, ver, peer string) error {
	targetPeer := resmgmt.WithTargetEndpoints(peer)

	// pack the chaincode
	ccPkg, err := gopackager.NewCCPackage(path, CCGoPath)
	if err != nil {
		return errors.WithMessage(err, "pack chaincode error")
	}

	// new request of installing chaincode
	req := resmgmt.InstallCCRequest{Name: cc,
		Path:    path,
		Version: ver,
		Package: ccPkg,
	}

	resps, err := c.rc.InstallCC(req, targetPeer)
	if err != nil {
		return errors.WithMessage(err, "installCC error")
	}

	// check other errors
	var errs []error
	for _, resp := range resps {
		log.Printf("Install  response status: %v", resp.Status)
		if resp.Status != http.StatusOK {
			errs = append(errs, errors.New(resp.Info))
		}
		if resp.Info == "already installed" {
			log.Printf("Chaincode %s already installed on peer: %s.\n",
				cc+"-"+ver, resp.Target)
			return nil
		}
	}

	if len(errs) > 0 {
		log.Printf("InstallCC errors: %v", errs)
		return errors.WithMessage(errs[0], "installCC first error")
	}
	return nil
}

// args := packArgs([]string{"init", "a", "100", "b", "200"})
func (c *Client) InstantiateCC(cc, path string, args []string, ver, peer string) (fab.TransactionID, error) {
	// endorser policy
	//org1OrOrg2 := "OR('Org1MSP.member','Org2MSP.member')"
	//	ccPolicy, err := c.genPolicy(org1OrOrg2)
	//	if err != nil {
	//		return "", errors.WithMessage(err, "gen policy from string error")
	//	}

	// new request
	// Attention: args should include `init` for Request not
	// have a method term to call init
	req := resmgmt.InstantiateCCRequest{
		Name:    cc,
		Path:    path,
		Version: ver,
		Args:    packArgs(args),
		Policy:  nil,
	}

	// send request and handle response
	reqPeers := resmgmt.WithTargetEndpoints(peer)
	resp, err := c.rc.InstantiateCC(c.Channel, req, reqPeers)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return "", nil
		}
		return "", errors.WithMessage(err, "instantiate chaincode error")
	}

	log.Printf("Instantitate chaincode tx: %s", resp.TransactionID)
	return resp.TransactionID, nil
}

//func (c *Client) genPolicy(p string) (*common.SignaturePolicyEnvelope, error) {
//	// TODO bug, this any leads to endorser invalid
//	if p == "ANY" {
//		return cauthdsl.SignedByAnyMember([]string{c.OrgName}), nil
//	}
//	return cauthdsl.FromString(p)
//}

func (c *Client) InvokeCC(cc, fun string, args [][]byte, peers []string) (fab.TransactionID, error) {
	// new channel request for invoke
	req := channel.Request{
		ChaincodeID: cc,
		Fcn:         fun,
		Args:        args,
	}
	// send request and handle response, peers is needed
	reqPeers := channel.WithTargetEndpoints(peers...)
	resp, err := c.cc.Execute(req, reqPeers)
	//log.Printf("Invoke chaincode response:\n id: %v\nvalidate: %v\nchaincode status: %v\n\n",
	//	resp.TransactionID,
	//	resp.TxValidationCode,
	//	resp.ChaincodeStatus)
	if err != nil {
		return "", errors.WithMessage(err, "invoke chaincode error")
	}
	return resp.TransactionID, nil
}

func (c *Client) QueryCC(cc, fun string, keys []string, peer string) (string, error) {
	// new channel request for query
	req := channel.Request{
		ChaincodeID: cc,
		Fcn:         fun,
		Args:        packArgs(keys),
	}

	// send request and handle response
	reqPeers := channel.WithTargetEndpoints(peer)
	resp, err := c.cc.Query(req, reqPeers)
	if err != nil {
		return "", errors.WithMessage(err, "query chaincode error")
	}
	//log.Printf("Query chaincode tx response:\ntx: %s\nresult: %v\n\n",
	//	resp.TransactionID,
	//	string(resp.Payload))
	return string(resp.Payload), nil
}

//	args := packArgs([]string{"init", "a", "1000", "b", "2000"})
func (c *Client) UpgradeCC(cc, path, ver string, args []string, peer string) error {
	// endorser policy
	//org1AndOrg2 := "AND('Org1MSP.member','Org2MSP.member')"
	//ccPolicy, err := c.genPolicy(org1AndOrg2)
	//if err != nil {
	//	return errors.WithMessage(err, "gen policy from string error")
	//}

	// new request
	// Attention: args should include `init` for Request not
	// have a method term to call init
	// Reset a b's value to test the upgrade
	req := resmgmt.UpgradeCCRequest{
		Name:    cc,
		Path:    path,
		Version: ver,
		Args:    packArgs(args),
		Policy:  nil,
	}

	// send request and handle response
	reqPeers := resmgmt.WithTargetEndpoints(peer)
	resp, err := c.rc.UpgradeCC(c.Channel, req, reqPeers)
	if err != nil {
		return errors.WithMessage(err, "instantiate chaincode error")
	}

	log.Printf("Instantitate chaincode tx: %s", resp.TransactionID)
	return nil
}

func (c *Client) Close() {
	c.SDK.Close()
}

func packArgs(paras []string) [][]byte {
	var args [][]byte
	for _, k := range paras {
		args = append(args, []byte(k))
	}
	return args
}
