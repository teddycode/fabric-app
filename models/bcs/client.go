package bcs

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	dsc "github.com/hyperledger/fabric-sdk-go/pkg/fab/discovery"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"log"
	"os"
)

var (
	CCGoPath  = os.Getenv("GOPATH") // GOPATH used for chaincode
)

type Client struct {
	// Fabric network information
	ConfigPath string
	OrgName    string
	OrgAdmin   string
	OrgUser    string
	Channel    string

	// sdk clients
	SDK *fabsdk.FabricSDK
	rc  *resmgmt.Client
	cc  *channel.Client
	mc  *msp.Client
	lc  *ledger.Client

	// discovery
	dc *dsc.Client // discovery client
}

func New(cfg, org, admin, user, chans string) *Client {
	c := &Client{
		ConfigPath: cfg,
		OrgName:    org,
		OrgAdmin:   admin,
		OrgUser:    user,
	}
	c.Channel = chans

	// create sdk
	sdk, err := fabsdk.New(config.FromFile(c.ConfigPath))
	if err != nil {
		log.Panicf("failed to create fabric sdk: %s", err)
	}
	c.SDK = sdk
	log.Println("Initialized fabric sdk")

	c.rc, c.cc, c.lc = NewSDKClient(sdk, chans, c.OrgName, c.OrgAdmin, c.OrgUser)

	// CA
	c.removeUserData()
	c.NewCAClient()

	return c
}

// TODO fix discovery service
// New discovery client
//func NewDiscoveryClient(sdk *fabsdk.FabricSDK, orgName, orgAdmin string) (dc *dsc.Client){
//	ccp := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(orgName))
//	fabContent.Client()
// 	dc, err := dsc.New()
//}

// NewSdkClient create resource client and channel client
func NewSDKClient(sdk *fabsdk.FabricSDK, channelID, orgName,
	orgAdmin, OrgUser string) (rc *resmgmt.Client, cc *channel.Client, lc *ledger.Client) {

	var err error
	// create rc
	rcp := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(orgName))
	rc, err = resmgmt.New(rcp)
	if err != nil {
		log.Panicf("failed to create resource client: %s", err)
	}
	log.Println("Initialized resource client")

	// create cc
	ccp := sdk.ChannelContext(channelID, fabsdk.WithUser(OrgUser))
	cc, err = channel.New(ccp)
	if err != nil {
		log.Panicf("failed to create channel client: %s", err)
	}
	log.Println("Initialized channel client")

	lc, err = ledger.New(ccp)
	if err != nil {
		log.Panicf("failed to create ledger client: %s", err)
	}
	log.Println("Initialized ledger client")

	return rc, cc, lc
}

// RegisterChaincodeEvent more easy than event client to registering chaincode event.
//func (c *Client) RegisterChaincodeEvent(ccid, eventName string) (fab.Registration, <-chan *fab.CCEvent, error) {
//	return c.cc.RegisterChaincodeEvent(ccid, eventName)
//}
