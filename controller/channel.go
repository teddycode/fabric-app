package controller

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"strconv"
)

func (c *Client) queryBlockHeight(client *ledger.Client) (string, error) {
	chainInfo, err := client.QueryInfo()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(int64(chainInfo.BCI.Height), 10), nil
}
