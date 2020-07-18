package utils

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	mspapi "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
)

func OrgUser(sdk *fabsdk.FabricSDK, orgID, username string) (mspapi.SigningIdentity, error) {
	if username == "" {
		return nil, errors.Errorf("no username specified")
	}
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgID))
	if err != nil {
		return nil, errors.Errorf("error creating MSP client: %s", err)
	}

	user, err := mspClient.GetSigningIdentity(username)
	if err != nil {
		return nil, errors.Errorf("GetSigningIdentity returned error: %v", err)
	}
	return user, nil
}
