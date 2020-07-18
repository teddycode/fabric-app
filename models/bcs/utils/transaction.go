package utils

import (
	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/pkg/errors"
)

type TxConvFactory struct {
}

func GetReadWriteSetFromEnvelope(envelope *cb.Envelope) []*rwsetutil.NsRwSet {
	payload := ExtractPayloadOrPanic(envelope)
	chdr, err := UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		panic(err)
	}
	if cb.HeaderType(chdr.Type) == cb.HeaderType_ENDORSER_TRANSACTION {
		tx, err := GetTransaction(payload.Data)
		if err != nil {
			panic(errors.Errorf("Bad envelope: %v", err))
		}
		for _, action := range tx.Actions {
			chaPayload, err := GetChaincodeActionPayload(action.Payload)
			if err != nil {
				panic(err)
			}
			//cpp := &peer.ChaincodeProposalPayload{}
			//err = proto.Unmarshal(chaPayload.ChaincodeProposalPayload, cpp)
			//if err != nil {
			//	panic(err)
			//}
			//cis := &peer.ChaincodeInvocationSpec{}
			//err = proto.Unmarshal(cpp.Input, cis)
			//if err != nil {
			//	panic(err)
			//}
			//for key, value := range cpp.TransientMap {
			//	p.Item("Key", key)
			//	p.Field("Value", value)
			//	p.ItemEnd()
			//}
			prp := &peer.ProposalResponsePayload{}
			unmarshalOrPanic(chaPayload.Action.ProposalResponsePayload, prp)
			chaincodeAction := &peer.ChaincodeAction{}
			unmarshalOrPanic(prp.Extension, chaincodeAction)
			if len(chaincodeAction.Results) > 0 {
				txRWSet := &rwsetutil.TxRwSet{}
				if err := txRWSet.FromProtoBytes(chaincodeAction.Results); err != nil {
					panic(err)
				}
				return txRWSet.NsRwSets
			}
		}
	}
	return nil
}

//
//// PrintData prints the block of data formatted according to the given HeaderType
//func ProccessData(headerType cb.HeaderType, data []byte) {
//	if headerType == cb.HeaderType_CONFIG {
//		envelope := &cb.ConfigEnvelope{}
//		if err := proto.Unmarshal(data, envelope); err != nil {
//			panic(errors.Errorf("Bad envelope: %v", err))
//		}
//		p.Print("Config Envelope:")
//		p.PrintConfigEnvelope(envelope)
//	} else if headerType == cb.HeaderType_CONFIG_UPDATE {
//		envelope := &cb.ConfigUpdateEnvelope{}
//		if err := proto.Unmarshal(data, envelope); err != nil {
//			panic(errors.Errorf("Bad envelope: %v", err))
//		}
//		p.Print("Config Update Envelope:")
//		p.PrintConfigUpdateEnvelope(envelope)
//	} else if headerType == cb.HeaderType_ENDORSER_TRANSACTION {
//		tx, err := GetTransaction(data)
//		if err != nil {
//			panic(errors.Errorf("Bad envelope: %v", err))
//		}
//		p.Print("Transaction:")
//		p.PrintTransaction(tx)
//	}
//}
