package bcs

import "C"
import (
	"fmt"
	clientMSP "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite"
	"github.com/hyperledger/fabric-sdk-go/pkg/msp"
	"log"
	"math/rand"
	"os"
)

func (c *Client) NewCAClient() {
	log.Print("Enroll registrar")
	ctxProvider := c.SDK.Context()
	mspClient, err := clientMSP.New(ctxProvider)
	//registrarEnrollID, registrarEnrollSecret := c.getRegistrarEnrollmentCredentials(ctxProvider)

	err = mspClient.Enroll("admin", clientMSP.WithSecret("adminpw"))
	if err != nil {
		log.Fatalf("enroll registrar failed: %v", err)
	}
	c.mc = mspClient
}

func (c *Client) removeUserData() {
	configBackend, err := c.SDK.Config()
	if err != nil {
		log.Fatal(err)
	}

	cryptoSuiteConfig := cryptosuite.ConfigFromBackend(configBackend)
	identityConfig, err := msp.ConfigFromBackend(configBackend)
	if err != nil {
		log.Fatal(err)
	}

	keyStorePath := cryptoSuiteConfig.KeyStorePath()
	credentialStorePath := identityConfig.CredentialStorePath()
	c.removePath(keyStorePath)
	c.removePath(credentialStorePath)
}

func (c *Client) removePath(storePath string) {
	err := os.RemoveAll(storePath)
	if err != nil {
		log.Fatalf("Cleaning up directory '%s' failed: %v", storePath, err)
	}
}

//func (c *Client) getRegistrarEnrollmentCredentials(ctxProvider context.ClientProvider) (string, string) {
//
//	ctx, err := ctxProvider()
//	if err != nil {
//		fmt.Printf("failed to get context: %v\n", err)
//	}
//
//	clientConfig := ctx.IdentityConfig().Client()
//	//if err != nil {
//	//	fmt.Printf("config.Client() failed: %v\n", err)
//	//}
//
//	myOrg := clientConfig.Organization
//
//	caConfig, ok := ctx.IdentityConfig().CAConfig(myOrg)
//	if ok {
//		fmt.Printf("CAConfig failed: %v\n", err)
//	}
//
//	return caConfig.Registrar.EnrollID, caConfig.Registrar.EnrollSecret
//}

func GenerateRandomID() string {
	return randomString(10)
}

// Utility to create random string of strlen length
func randomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// Register a new user
func (c *Client) RegisterUser(username, orgName, secret, identityTypeUser string) (string, bool) {
	// Register the new user
	log.Printf("User not found, registering new user: %v", username)
	testAttributes := []clientMSP.Attribute{
		{
			Name:  GenerateRandomID(),
			Value: fmt.Sprintf("%s:ecert", GenerateRandomID()),
			ECert: true,
		},
		{
			Name:  GenerateRandomID(),
			Value: fmt.Sprintf("%s:ecert", GenerateRandomID()),
			ECert: true,
		},
	}
	_, err := c.mc.Register(&clientMSP.RegistrationRequest{
		Name:        username,
		Type:        identityTypeUser,
		Attributes:  testAttributes,
		Affiliation: orgName,
		Secret:      secret, // Is ready to get hash?
	})
	if err != nil {
		return username + err.Error(), false
	}
	//signingIdentity, err := mspClient.GetSigningIdentity(username)
	//log.Printf("%s: %s", signingIdentity.Identifier().ID, string(signingIdentity.EnrollmentCertificate()[:]))
	return username + " register Successfully", true

}

// Enroll a user
func (c *Client) EnrollUser(username, orgName, secret, identityTypeUser string) (string, bool) {
	//ctxProvider := c.SDK.Context(fabsdk.WithOrg(orgName))
	//mspClient, err := clientMSP.New(ctxProvider)
	//if err != nil {
	//	log.Fatalf("Failed to create msp client: %s", err.Error())
	//	return username + " login error :" + err.Error(), false
	//}
	err := c.mc.Enroll(username, clientMSP.WithSecret(secret))
	if err != nil {
		log.Printf("enroll %s failed: %v", username, err)
		return err.Error(), false
	}
	signingIdentity, err := c.mc.GetSigningIdentity(username)
	log.Printf("%s: %s", signingIdentity.Identifier().ID, string(signingIdentity.EnrollmentCertificate()[:]))
	return username + " login success", true
}

// Reovoke a user
func (c *Client) RevokeUser(username, orgName, secret, identityTypeUser string) (string, bool) {
	request := clientMSP.RemoveIdentityRequest{
		ID:     username,
		Force:  true,
		CAName: "ca.org1.lzawt.com",
	}
	idr, err := c.mc.RemoveIdentity(&request)
	if err != nil {
		log.Printf("enroll %s failed: %v", username, err)
		return err.Error(), false
	}
	log.Println(idr)
	return username + " RevokeUser success", true
}

// getRegisteredUser get registered user. If user is not enrolled, enroll new user
func (c *Client) GetRegisteredUser(username, orgName, secret, identityTypeUser string) (string, bool) {
	//ctxProvider := c.SDK.Context(fabsdk.WithOrg(orgName))
	//mspClient, err := clientMSP.New(ctxProvider)
	//if err != nil {
	//	log.Fatalf("Failed to create msp client: %s", err.Error())
	//}
	signingIdentity, err := c.mc.GetSigningIdentity(username)
	if err != nil {
		log.Printf("Check if user %s is enrolled: %s", username, err.Error())
		testAttributes := []clientMSP.Attribute{
			{
				Name:  GenerateRandomID(),
				Value: fmt.Sprintf("%s:ecert", GenerateRandomID()),
				ECert: true,
			},
			{
				Name:  GenerateRandomID(),
				Value: fmt.Sprintf("%s:ecert", GenerateRandomID()),
				ECert: true,
			},
		}

		// Register the new user
		identity, err := c.mc.GetIdentity(username)
		if true {
			log.Printf("User %s does not exist, registering new user", username)
			_, err = c.mc.Register(&clientMSP.RegistrationRequest{
				Name:        username,
				Type:        identityTypeUser,
				Attributes:  testAttributes,
				Affiliation: orgName,
				Secret:      secret,
			})
		} else {
			log.Printf("Identity: %s", identity.Secret)
		}
		//enroll user
		err = c.mc.Enroll(username, clientMSP.WithSecret(secret))
		if err != nil {
			log.Printf("enroll %s failed: %v", username, err)
			return "failed " + err.Error(), false
		}

		return username + " enrolled Successfully", true
	}
	log.Printf("%s: %s", signingIdentity.Identifier().ID, string(signingIdentity.EnrollmentCertificate()[:]))
	return username + " already enrolled", true
}
