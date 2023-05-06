package proofs

import (
	"github.com/iden3/go-circuits"
	"os"
	"testing"

	"github.com/Coursant/Core_rapidsnark/types"

	"github.com/Coursant/Core_common/protocol"
	"github.com/stretchr/testify/assert"
)

func TestVerifyProof(t *testing.T) {

	var err error
	proofMessage := protocol.ZeroKnowledgeProofResponse{
		ZKProof: types.ZKProof{
			Proof: &types.ProofData{
				A: []string{
					"261068577516437401613944053873182458364288414130914048345483377226144652651",
					"14191260071695980011679501808453222267520721767757759150101974382053161674611",
					"1",
				},
				B: [][]string{
					{
						"7670847844015116957526183728196977957312627307797919554134684901401436021977",
						"14957845472630017095821833222580194061266186851634053897768738253663253650835",
					},
					{
						"17835642458484628627556329876919077333912011235308758832172880012813397022104",
						"18100861130149678153133025031709897120097098591298817367491920553037011650228",
					},
					{
						"1",
						"0",
					}},
				C: []string{
					"6217865949299990642832523256863048932210546049203189113362851476966824162191",
					"19016949225277755690019647385855936969928994210905992628301967883803670436510",
					"1",
				},
				Protocol: "groth16",
			},
			PubSignals: []string{
				"1",
				"27152676987128542066808591998573000370436464722519513348891049644813718018",
				"23",
				"27752766823371471408248225708681313764866231655187366071881070918984471042",
				"21545768883509657340209171549441005603306012513932221371599501498534807719689",
				"1",
				"21545768883509657340209171549441005603306012513932221371599501498534807719689",
				"1679323038",
				"336615423900919464193075592850483704600",
				"0",
				"17002437119434618783545694633038537380726339994244684348913844923422470806844",
				"0",
				"5",
				"840",
				"120",
				"340",
				"509",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
			},
		},
	}
	proofMessage.CircuitID = string(circuits.AtomicQueryMTPV2CircuitID)

	verificationKey, err := os.ReadFile("../testdata/credentialAtomicQueryMTPV2.json")
	assert.NoError(t, err)

	proofMessage.ID = 1

	err = VerifyProof(proofMessage, verificationKey)
	assert.Nil(t, err)
}
