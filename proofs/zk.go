package proofs

import (
	"github.com/Coursant/Core_common/protocol"
	"github.com/Coursant/Core_rapidsnark/verifier"
	"github.com/pkg/errors"
)

// VerifyProof performs protocol proof response verification
func VerifyProof(resp protocol.ZeroKnowledgeProofResponse, verificationKey []byte) (err error) {
	switch resp.Proof.Protocol {
	case "groth16":
		return verifier.VerifyGroth16(resp.ZKProof, verificationKey)
	default:
		return errors.Errorf("%s protocol is not supported", resp.Proof.Protocol)
	}
}
