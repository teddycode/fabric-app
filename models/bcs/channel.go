package bcs

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/fabric-app/models/bcs/utils"
	"github.com/fabric-app/pkg/logging"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"strconv"
	"strings"
)

type Responser struct {
	Type      string `json:"type"`
	Point     string `json:"point"`
	TimeStamp string `json:"time_stamp"`
	Raw       string `json:"raw"`
}
type PeerInfo struct {
	ID      string `json:"msp_id"`
	Address string `json:"address"`
}

func (c *Client) GetBlockHeight() (string, error) {
	chainInfo, err := c.lc.QueryInfo()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(int64(chainInfo.BCI.Height), 10), nil
}

func (c *Client) QueryTxByID(txid string, endpoint string) (string, error) {
	transactions, err := c.lc.QueryTransaction(fab.TransactionID(txid), ledger.WithTargetEndpoints(endpoint))
	if err != nil {
		logging.Error("QueryTransaction error: " + err.Error())
		return "", err
	}
	rws := utils.GetReadWriteSetFromEnvelope(transactions.TransactionEnvelope)
	var writes []string
	bf := bytes.Buffer{}
	for _, u := range rws {
		for _, w := range u.KvRwSet.Writes {
			if w.IsDelete {
				continue
			}
			value := base64.StdEncoding.EncodeToString(w.GetValue())
			fmt.Println(string(w.GetValue()))
			keys := strings.Split(w.GetKey(), "~")
			resp := Responser{
				Type:      keys[0],
				Point:     keys[1],
				TimeStamp: keys[2],
				Raw:       value,
			}
			json1, _ := json.Marshal(&resp)
			writes = append(writes, string(json1))
		}
	}

	bf.WriteString("[")
	for i, x := range writes {
		bf.WriteString(x)
		if i < len(writes)-1 {
			bf.WriteString(",")
		}
	}
	bf.WriteString("]")
	fmt.Println(bf.String())

	return bf.String(), nil
}

//func (c *Client) QueryTxByID(txid string, endpoint string) ([]byte, error) {
//	transactions, err := c.lc.QueryTransaction(fab.TransactionID(txid), ledger.WithTargetEndpoints(endpoint))
//	if err != nil {
//		logging.Error("QueryTransaction error: " + err.Error())
//		return nil, err
//	}
//
//	//var reads, writes []string
//	var writes []string
//	bf := bytes.Buffer{}
//	rws := transh.GetReadWriteSet(transactions.TransactionEnvelope)
//	for _, u := range rws {
//		for _, w := range u.KVRWSet.Writes {
//			if w.IsDelete {
//				continue
//			}
//			bw, _ := json.Marshal(w)
//			writes = append(writes, string(bw))
//		}
//	}
//
//	bf.WriteString("{\"writes\":[")
//	for i, x := range writes {
//		bf.WriteString(x)
//		if i < len(writes)-1 {
//			bf.WriteString(",")
//		}
//	}
//	bf.WriteString("]}")
//
//	return bf.Bytes(), nil
//}
//func (c *Client) QueryPeers() ([]string, error) {
//	local, err := contextImpl.NewLocal(c.SDK.Context())
//	if err != nil {
//		return nil, err
//	}
//	peers, err := local.LocalDiscoveryService().GetPeers()
//	if err != nil {
//		return nil, err
//	}
//	pss := []string{}
//	for _, peer := range peers {
//		pss = append(pss, peer.URL())
//	}
//	return pss, nil
//}

func (c *Client) QueryPeers() ([]string, error) {

	chProvider := c.SDK.ChannelContext(c.Channel, fabsdk.WithUser(c.OrgUser))
	chContent, err := chProvider()
	if err != nil {
		return nil, err
	}

	//chContext, err := contextImpl.NewChannel(c.SDK.Context(), c.Channel)
	//if err != nil {
	//	return nil, err
	//}

	discovery, err := chContent.ChannelService().Discovery()
	if err != nil {
		return nil, err
	}

	peers, err := discovery.GetPeers()
	if err != nil {
		return nil, err
	}
	pss := []string{}
	for _, peer := range peers {
		peer := PeerInfo{
			ID:      peer.MSPID(),
			Address: peer.URL(),
		}
		str, _ := json.Marshal(peer)
		pss = append(pss, string(str))
	}
	return pss, nil
}
