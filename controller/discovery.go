package controller

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/msp"
	discoveryclient "github.com/hyperledger/fabric-sdk-go/internal/github.com/hyperledger/fabric/discovery/client"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/discovery"
)

var (
	channels   = "mychannel"
	server     = "peer0.org1.lzawt.com"
	configFile = "discover-conf.yaml"
)

type localPeer struct {
	MSPID    string
	Endpoint string
	Identity string
}

type channelPeer struct {
	MSPID        string
	LedgerHeight uint64
	Endpoint     string
	Identity     string
	Chaincodes   []string
}

// PeerCmd executes channelPeer listing command
type Finder struct {
	stub    Stub
	server  *string
	channel *string
}

type peerLister interface {
	Peers() ([]*discoveryclient.Peer, error)
}

// Stub represents the remote discovery service
type Stub interface {
	// Send sends the request, and receives a response
	Send(server string, conf string, req *discovery.Request) (ServiceResponse, error)
}

type simpleChannelResponse struct {
	discoveryclient.ChannelResponse
}

func (scr *simpleChannelResponse) Peers() ([]*discoveryclient.Peer, error) {
	return scr.ChannelResponse.Peers()
}

// ServiceResponse represents a response sent from the discovery service
type ServiceResponse interface {
	// ForChannel returns a ChannelResponse in the context of a given channel
	ForChannel(string) discoveryclient.ChannelResponse

	// ForLocal returns a LocalResponse in the context of no channel
	ForLocal() discoveryclient.LocalResponse

	// Raw returns the raw response from the server
	Raw() *discovery.Response
}

// Execute executes the command
func (pc *Finder) Execute() ([]byte, error) {
	req := discovery.NewRequest()
	req = req.OfChannel(channels)
	req = req.AddPeersQuery()
	res, err := pc.stub.Send(server, configFile, req)
	if err != nil {
		return nil, err
	}
	return ParseResponse(channels, res)
}

// ParseResponse parses the given response about the given channel
func ParseResponse(channel string, res ServiceResponse) ([]byte, error) {
	var listPeers peerLister
	if channel == "" {
		listPeers = res.ForLocal()
	} else {
		listPeers = &simpleChannelResponse{res.ForChannel(channel)}
	}
	peers, err := listPeers.Peers()
	if err != nil {
		return nil, err
	}

	channelState := channel != ""
	b, _ := json.MarshalIndent(assemblePeers(peers, channelState), "", "\t")
	return (b), nil
}

func assemblePeers(peers []*discoveryclient.Peer, withChannelState bool) interface{} {
	if withChannelState {
		var peerSlices []channelPeer
		for _, p := range peers {
			peerSlices = append(peerSlices, rawPeerToChannelPeer(p))
		}
		return peerSlices
	}
	var peerSlices []localPeer
	for _, p := range peers {
		peerSlices = append(peerSlices, rawPeerToLocalPeer(p))
	}
	return peerSlices
}
func rawPeerToChannelPeer(p *discoveryclient.Peer) channelPeer {
	var ledgerHeight uint64
	var ccs []string
	if p.StateInfoMessage != nil && p.StateInfoMessage.GetStateInfo() != nil && p.StateInfoMessage.GetStateInfo().Properties != nil {
		properties := p.StateInfoMessage.GetStateInfo().Properties
		ledgerHeight = properties.LedgerHeight
		for _, cc := range properties.Chaincodes {
			if cc == nil {
				continue
			}
			ccs = append(ccs, cc.Name)
		}
	}
	var endpoint string
	if p.AliveMessage != nil && p.AliveMessage.GetAliveMsg() != nil && p.AliveMessage.GetAliveMsg().Membership != nil {
		endpoint = p.AliveMessage.GetAliveMsg().Membership.Endpoint
	}
	sID := &msp.SerializedIdentity{}
	proto.Unmarshal(p.Identity, sID)
	return channelPeer{
		MSPID:        p.MSPID,
		Endpoint:     endpoint,
		LedgerHeight: ledgerHeight,
		Identity:     string(sID.IdBytes),
		Chaincodes:   ccs,
	}
}

func rawPeerToLocalPeer(p *discoveryclient.Peer) localPeer {
	var endpoint string
	if p.AliveMessage != nil && p.AliveMessage.GetAliveMsg() != nil && p.AliveMessage.GetAliveMsg().Membership != nil {
		endpoint = p.AliveMessage.GetAliveMsg().Membership.Endpoint
	}
	sID := &msp.SerializedIdentity{}
	proto.Unmarshal(p.Identity, sID)
	return localPeer{
		MSPID:    p.MSPID,
		Endpoint: endpoint,
		Identity: string(sID.IdBytes),
	}
}
