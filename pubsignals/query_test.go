package pubsignals

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"testing"

	core "github.com/Coursant/Core_origin"
	"github.com/stretchr/testify/require"
)

var (
	issuerDID   = "did:polygonid:polygon:mumbai:2qHSHBGWGJ68AosMKcLCTp8FYdVrtYE6MtNHhq8xpK"
	iid, _      = new(big.Int).SetString("22638457188543025296541325416907897762715008870723718557276875842936181250", 10)
	issuerID, _ = core.IDFromInt(iid)

	schemaHashKYCEmployee, _ = big.NewInt(0).SetString("40507426258035268209384482717564233104", 10)
	KYCEmployeeSchema        = core.NewSchemaHashFromInt(schemaHashKYCEmployee)

	schemaHashKYCCountry, _ = big.NewInt(0).SetString("336615423900919464193075592850483704600", 10)
	KYCCountrySchema        = core.NewSchemaHashFromInt(schemaHashKYCCountry)

	bigIntTrueHash, _ = big.NewInt(0).SetString("18586133768512220936620570745912940619677854269274689475585506675881198879027", 10)
)

type mockMemorySchemaLoader struct {
}

func (r *mockMemorySchemaLoader) Load(_ context.Context, _ string) (schema []byte, ext string, err error) {
	return []byte(`{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "$metadata": {
    "uris": {
      "jsonLdContext": "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
      "jsonSchema": "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json/KYCCountryOfResidenceCredential-v2.json"
    },
	"serialization": {
		"valueDataSlotB": "countryCode"
	}
  },
  "@context": [
    {
      "@version": 1.1,
      "@protected": true,
      "id": "@id",
      "type": "@type",
      "KYCAgeCredential": {
        "@id": "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld#KYCAgeCredential",
        "@context": {
          "@version": 1.1,
          "@protected": true,
          "id": "@id",
          "type": "@type",
          "kyc-vocab": "https://github.com/iden3/claim-schema-vocab/blob/main/credentials/kyc.md#",
          "xsd": "http://www.w3.org/2001/XMLSchema#",
          "birthday": {
            "@id": "kyc-vocab:birthday",
            "@type": "xsd:integer"
          },
          "documentType": {
            "@id": "kyc-vocab:documentType",
            "@type": "xsd:integer"
          }
        }
      },
      "KYCCountryOfResidenceCredential": {
        "@id": "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld#KYCCountryOfResidenceCredential",
        "@context": {
          "@version": 1.1,
          "@protected": true,
          "id": "@id",
          "type": "@type",
          "kyc-vocab": "https://github.com/iden3/claim-schema-vocab/blob/main/credentials/kyc.md#",
          "xsd": "http://www.w3.org/2001/XMLSchema#",
          "countryCode": {
            "@id": "kyc-vocab:countryCode",
            "@type": "xsd:integer"
          },
          "documentType": {
            "@id": "kyc-vocab:documentType",
            "@type": "xsd:integer"
          }
        }
      },
	  "KYCEmployee": {
        "@id": "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld#KYCEmployee",
        "@context": {
          "@version": 1.1,
          "@protected": true,
          "id": "@id",
          "type": "@type",
          "kyc-vocab": "https://github.com/iden3/claim-schema-vocab/blob/main/credentials/kyc.md#",
          "xsd": "http://www.w3.org/2001/XMLSchema#",
          "documentType": {
            "@id": "kyc-vocab:documentType",
            "@type": "xsd:integer"
          },
          "ZKPexperiance": {
            "@id": "kyc-vocab:hasZKPexperiance",
            "@type": "xsd:boolean"
          },
          "hireDate": {
            "@id": "kyc-vocab:hireDate",
            "@type": "xsd:dateTime"
          },
          "position": {
            "@id": "kyc-vocab:position",
            "@type": "xsd:string"
          },
          "salary": {
            "@id": "kyc-vocab:salary",
            "@type": "xsd:double"
          }
        }
      }
    }
  ]
}
`), "json-ld", nil
}

var vp = []byte(`{
	"@context": [
		"https://www.w3.org/2018/credentials/v1"
	],
	"@type": "VerifiablePresentation",
	"verifiableCredential": {
		"@context": [
			"https://www.w3.org/2018/credentials/v1",
			"https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld"
		],
		"@type": ["VerifiableCredential","KYCCountryOfResidenceCredential"],
		"credentialSubject": {
			"type": "KYCCountryOfResidenceCredential",
			"countryCode": 800
		}
	}
}`)

var vpEmployee = []byte(`{
	"@type": "VerifiablePresentation",
	"@context": [
		"https://www.w3.org/2018/credentials/v1"
	],
	"verifiableCredential": {
		"@context": [
			"https://www.w3.org/2018/credentials/v1",
			"https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v101.json-ld"
		],
		"@type": [
			"VerifiableCredential",
			"KYCEmployee"
		],
		"credentialSubject": {
			"@type": "KYCEmployee",
			"position": "SSI Consultant"
		}
	}
}`)

func TestCheckRequest_Success(t *testing.T) {
	tests := []struct {
		name   string
		query  Query
		pubSig *CircuitOutputs
		vp     json.RawMessage
	}{
		{
			name: "Check merkalized query",
			query: Query{
				AllowedIssuers: []string{"*"},
				CredentialSubject: map[string]interface{}{
					"countryCode": map[string]interface{}{
						"$nin": []interface{}{float64(800)},
					},
				},
				Context: "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
				Type:    "KYCCountryOfResidenceCredential",
			},
			pubSig: &CircuitOutputs{
				IssuerID:    &issuerID,
				ClaimSchema: KYCCountrySchema,
				ClaimPathKey: func() *big.Int {
					v, _ := big.NewInt(0).SetString("17002437119434618783545694633038537380726339994244684348913844923422470806844", 10)
					return v
				}(),
				Operator:            5,
				Value:               []*big.Int{big.NewInt(800)},
				Merklized:           1,
				IsRevocationChecked: 1,
			},
		},
		{
			name: "Selective disclosure",
			query: Query{
				AllowedIssuers: []string{"*"},
				CredentialSubject: map[string]interface{}{
					"countryCode": map[string]interface{}{},
				},
				Context: "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
				Type:    "KYCCountryOfResidenceCredential",
			},
			pubSig: &CircuitOutputs{
				IssuerID:    &issuerID,
				ClaimSchema: KYCCountrySchema,
				ClaimPathKey: func() *big.Int {
					v, _ := big.NewInt(0).SetString("17002437119434618783545694633038537380726339994244684348913844923422470806844", 10)
					return v
				}(),
				Operator:            1,
				Value:               []*big.Int{big.NewInt(800)},
				Merklized:           1,
				IsRevocationChecked: 1,
			},
			vp: vp,
		},
		{
			name: "Query with boolean type",
			query: Query{
				AllowedIssuers: []string{"*"},
				CredentialSubject: map[string]interface{}{
					"ZKPexperiance": map[string]interface{}{
						"$eq": true,
					},
				},
				Context: "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
				Type:    "KYCEmployee",
			},
			pubSig: &CircuitOutputs{
				IssuerID:    &issuerID,
				ClaimSchema: KYCEmployeeSchema,
				ClaimPathKey: func() *big.Int {
					v, _ := big.NewInt(0).SetString("1944808975288007371356450257872165609440470546066507760733183342797918372827", 10)
					return v
				}(),
				Operator:            1,
				Value:               []*big.Int{bigIntTrueHash},
				Merklized:           1,
				IsRevocationChecked: 1,
			},
		},
		{
			name: "Selective disclosure with xsd:string type",
			query: Query{
				AllowedIssuers: []string{"*"},
				CredentialSubject: map[string]interface{}{
					"position": map[string]interface{}{},
				},
				Context: "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
				Type:    "KYCEmployee",
			},
			pubSig: &CircuitOutputs{
				IssuerID:    &issuerID,
				ClaimSchema: KYCEmployeeSchema,
				ClaimPathKey: func() *big.Int {
					v, _ := big.NewInt(0).SetString("15406634529806189041952040954758558497189093183268091368437514469450172572054", 10)
					return v
				}(),
				Operator: 1,
				Value: func() []*big.Int {
					v, _ := big.NewInt(0).SetString("957410455271905675920624030785024750144198809104092676617070098470852489834", 10)
					return []*big.Int{v}
				}(),
				Merklized:           1,
				IsRevocationChecked: 1,
			},
			vp: vpEmployee,
		},
		{
			name: "EQ operator for xsd:string type",
			query: Query{
				AllowedIssuers: []string{"*"},
				CredentialSubject: map[string]interface{}{
					"position": map[string]interface{}{
						"$eq": "Software Engineer",
					},
				},
				Context: "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
				Type:    "KYCEmployee",
			},
			pubSig: &CircuitOutputs{
				IssuerID:    &issuerID,
				ClaimSchema: KYCEmployeeSchema,
				ClaimPathKey: func() *big.Int {
					v, _ := big.NewInt(0).SetString("15406634529806189041952040954758558497189093183268091368437514469450172572054", 10)
					return v
				}(),
				Operator: 1,
				Value: func() []*big.Int {
					v, _ := big.NewInt(0).SetString("7481731651336040098616464366227645531920423822088928207225802836605991806542", 10)
					return []*big.Int{v}
				}(),
				Merklized:           1,
				IsRevocationChecked: 1,
			},
			vp: vpEmployee,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.query.Check(context.Background(), &mockMemorySchemaLoader{}, tt.pubSig, tt.vp)
			require.NoError(t, err)
		})
	}
}

func TestCheckRequest_SelectiveDisclosure_Error(t *testing.T) {
	tests := []struct {
		name   string
		query  Query
		pubSig *CircuitOutputs
		vp     json.RawMessage
		expErr error
	}{
		{
			name: "Empty disclosure value for disclosure request",
			query: Query{
				AllowedIssuers: []string{"*"},
				CredentialSubject: map[string]interface{}{
					"countryCode": map[string]interface{}{},
				},
				Context: "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
				Type:    "KYCCountryOfResidenceCredential",
			},
			vp: nil,
			pubSig: &CircuitOutputs{
				IssuerID:    &issuerID,
				ClaimSchema: KYCCountrySchema,
				ClaimPathKey: func() *big.Int {
					v, _ := big.NewInt(0).SetString("17002437119434618783545694633038537380726339994244684348913844923422470806844", 10)
					return v
				}(),
				Operator:            1,
				Value:               []*big.Int{big.NewInt(800)},
				Merklized:           1,
				IsRevocationChecked: 1,
			},
			expErr: errors.New("selective disclosure value is missed"),
		},
		{
			name: "Not EQ operation for disclosure request",
			query: Query{
				AllowedIssuers: []string{"*"},
				CredentialSubject: map[string]interface{}{
					"countryCode": map[string]interface{}{},
				},
				Context: "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
				Type:    "KYCCountryOfResidenceCredential",
			},
			vp: vp,
			pubSig: &CircuitOutputs{
				IssuerID:    &issuerID,
				ClaimSchema: KYCCountrySchema,
				ClaimPathKey: func() *big.Int {
					v, _ := big.NewInt(0).SetString("17002437119434618783545694633038537380726339994244684348913844923422470806844", 10)
					return v
				}(),
				Operator:            5,
				Value:               []*big.Int{big.NewInt(800)},
				Merklized:           1,
				IsRevocationChecked: 1,
			},
			expErr: errors.New("selective disclosure available only for equal operation"),
		},
		{
			name: "Not array of values for disclosure request",
			query: Query{
				AllowedIssuers: []string{"*"},
				CredentialSubject: map[string]interface{}{
					"countryCode": map[string]interface{}{},
				},
				Context: "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
				Type:    "KYCCountryOfResidenceCredential",
			},
			vp: vp,
			pubSig: &CircuitOutputs{
				IssuerID:    &issuerID,
				ClaimSchema: KYCCountrySchema,
				ClaimPathKey: func() *big.Int {
					v, _ := big.NewInt(0).SetString("17002437119434618783545694633038537380726339994244684348913844923422470806844", 10)
					return v
				}(),
				Operator:            1,
				Value:               []*big.Int{big.NewInt(800), big.NewInt(801)},
				Merklized:           1,
				IsRevocationChecked: 1,
			},
			expErr: errors.New("selective disclosure not available for array of values"),
		},
		{
			name: "Proof was generated for another disclosure value",
			query: Query{
				AllowedIssuers: []string{"*"},
				CredentialSubject: map[string]interface{}{
					"countryCode": map[string]interface{}{},
				},
				Context: "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
				Type:    "KYCCountryOfResidenceCredential",
			},
			vp: vp,
			pubSig: &CircuitOutputs{
				IssuerID:    &issuerID,
				ClaimSchema: KYCCountrySchema,
				ClaimPathKey: func() *big.Int {
					v, _ := big.NewInt(0).SetString("17002437119434618783545694633038537380726339994244684348913844923422470806844", 10)
					return v
				}(),
				Operator:            1,
				Value:               []*big.Int{big.NewInt(1)},
				Merklized:           1,
				IsRevocationChecked: 1,
			},
			expErr: errors.New("different value between proof and disclosure value"),
		},
		{
			name: "Different key between proof and disclosure response",
			query: Query{
				AllowedIssuers: []string{"*"},
				CredentialSubject: map[string]interface{}{
					"documentType": map[string]interface{}{},
				},
				Context: "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
				Type:    "KYCCountryOfResidenceCredential",
			},
			vp: vp,
			pubSig: &CircuitOutputs{
				IssuerID:    &issuerID,
				ClaimSchema: KYCCountrySchema,
				ClaimPathKey: func() *big.Int {
					v, _ := big.NewInt(0).SetString("17002437119434618783545694633038537380726339994244684348913844923422470806844", 10)
					return v
				}(),
				Operator:            1,
				Value:               []*big.Int{big.NewInt(800)},
				Merklized:           1,
				IsRevocationChecked: 1,
			},
			expErr: errors.New("path '[https://www.w3.org/2018/credentials#verifiableCredential https://www.w3.org/2018/credentials#credentialSubject https://github.com/iden3/claim-schema-vocab/blob/main/credentials/kyc.md#documentType]' doesn't exist in document"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.query.Check(context.Background(), &mockMemorySchemaLoader{}, tt.pubSig, tt.vp)
			require.EqualError(t, err, tt.expErr.Error())
		})
	}
}

func TestCheckRequest_Error(t *testing.T) {
	tests := []struct {
		name   string
		query  Query
		pubSig *CircuitOutputs
		expErr error
	}{
		{
			name: "Invalid issuer",
			query: Query{
				AllowedIssuers: []string{"123"},
			},
			pubSig: &CircuitOutputs{
				IssuerID: &issuerID,
			},
			expErr: ErrUnavailableIssuer,
		},
		{
			name: "Invalid Schema ID",
			query: Query{
				AllowedIssuers: []string{issuerDID},
				Context:        "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
				Type:           "KYCAgeCredential",
			},
			pubSig: &CircuitOutputs{
				IssuerID:    &issuerID,
				ClaimSchema: KYCCountrySchema,
			},
			expErr: ErrSchemaID,
		},
		{
			name: "Multiply query",
			query: Query{
				AllowedIssuers: []string{issuerDID},
				Context:        "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
				Type:           "KYCCountryOfResidenceCredential",
				CredentialSubject: map[string]interface{}{
					"req1": struct{}{},
					"req2": struct{}{},
				},
			},
			pubSig: &CircuitOutputs{
				IssuerID:    &issuerID,
				ClaimSchema: KYCCountrySchema,
			},
			expErr: errors.New("multiple requests not supported"),
		},
		{
			name: "Failed params in request",
			query: Query{
				AllowedIssuers: []string{issuerDID},
				Context:        "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
				Type:           "KYCCountryOfResidenceCredential",
				CredentialSubject: map[string]interface{}{
					"req1": 1,
				},
			},
			pubSig: &CircuitOutputs{
				IssuerID:    &issuerID,
				ClaimSchema: KYCCountrySchema,
			},
			expErr: errors.New("failed cast type map[string]interface"),
		},
		{
			name: "Multiple predicates in one request",
			query: Query{
				AllowedIssuers: []string{issuerDID},
				Context:        "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
				Type:           "KYCCountryOfResidenceCredential",
				CredentialSubject: map[string]interface{}{
					"countryCode": map[string]interface{}{
						"$eq":  20,
						"$nin": 21,
					},
				},
			},
			pubSig: &CircuitOutputs{
				IssuerID:    &issuerID,
				ClaimSchema: KYCCountrySchema,
			},
			expErr: errors.New("multiple predicates for one field not supported"),
		},
		{
			name: "Proof was generated for another query operator",
			query: Query{
				AllowedIssuers: []string{issuerDID},
				Context:        "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
				Type:           "KYCCountryOfResidenceCredential",
				CredentialSubject: map[string]interface{}{
					"countryCode": map[string]interface{}{
						"$eq": []interface{}{float64(20)},
					},
				},
			},
			pubSig: &CircuitOutputs{
				IssuerID:    &issuerID,
				ClaimSchema: KYCCountrySchema,
				Operator:    3,
			},
			expErr: ErrRequestOperator,
		},
		{
			name: "Proof was generated for another values",
			query: Query{
				AllowedIssuers: []string{issuerDID},
				Context:        "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
				Type:           "KYCCountryOfResidenceCredential",
				CredentialSubject: map[string]interface{}{
					"countryCode": map[string]interface{}{
						"$nin": []interface{}{float64(20)},
					},
				},
			},
			pubSig: &CircuitOutputs{
				IssuerID:    &issuerID,
				ClaimSchema: KYCCountrySchema,
				Operator:    5,
				Value:       []*big.Int{big.NewInt(40)},
			},
			expErr: ErrInvalidValues,
		},
		{
			name: "Proof was generated for another path",
			query: Query{
				AllowedIssuers: []string{issuerDID},
				Context:        "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
				Type:           "KYCCountryOfResidenceCredential",
				CredentialSubject: map[string]interface{}{
					"documentType": map[string]interface{}{
						"$nin": []interface{}{float64(20)},
					},
				},
			},
			pubSig: &CircuitOutputs{
				IssuerID:            &issuerID,
				ClaimSchema:         KYCCountrySchema,
				ClaimPathKey:        big.NewInt(0),
				Operator:            5,
				Value:               []*big.Int{big.NewInt(20)},
				Merklized:           1,
				IsRevocationChecked: 1,
			},
			expErr: errors.New("proof was generated for another path"),
		},
		{
			name: "Different slot index",
			query: Query{
				AllowedIssuers: []string{issuerDID},
				Context:        "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
				Type:           "KYCCountryOfResidenceCredential",
				CredentialSubject: map[string]interface{}{
					"countryCode": map[string]interface{}{
						"$nin": []interface{}{float64(20)},
					},
				},
			},
			pubSig: &CircuitOutputs{
				IssuerID:            &issuerID,
				ClaimSchema:         KYCCountrySchema,
				Operator:            5,
				Value:               []*big.Int{big.NewInt(20)},
				Merklized:           0,
				SlotIndex:           0,
				IsRevocationChecked: 1,
			},
			expErr: errors.New("different slot index for claim"),
		},
		{
			name: "Check revocation is required",
			query: Query{
				AllowedIssuers: []string{issuerDID},
				Context:        "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
				Type:           "KYCCountryOfResidenceCredential",
				CredentialSubject: map[string]interface{}{
					"countryCode": map[string]interface{}{
						"$nin": []interface{}{float64(20)},
					},
				},
				SkipClaimRevocationCheck: false,
			},
			pubSig: &CircuitOutputs{
				IssuerID:            &issuerID,
				ClaimSchema:         KYCCountrySchema,
				Operator:            5,
				Value:               []*big.Int{big.NewInt(20)},
				Merklized:           0,
				SlotIndex:           0,
				IsRevocationChecked: 0,
			},
			expErr: errors.New("check revocation is required"),
		},
		{
			name: "Unsupported lt operator for xsd:boolean",
			query: Query{
				AllowedIssuers: []string{issuerDID},
				Context:        "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
				Type:           "KYCEmployee",
				CredentialSubject: map[string]interface{}{
					"ZKPexperiance": map[string]interface{}{
						"$lt": 20,
					},
				},
				SkipClaimRevocationCheck: false,
			},
			pubSig: &CircuitOutputs{
				IssuerID:            &issuerID,
				ClaimSchema:         KYCEmployeeSchema,
				Operator:            2,
				Value:               []*big.Int{big.NewInt(20)},
				Merklized:           0,
				SlotIndex:           0,
				IsRevocationChecked: 0,
			},
			expErr: errors.New("invalid operation '$lt' for field type 'http://www.w3.org/2001/XMLSchema#boolean'"),
		},
		{
			name: "Negative value in request",
			query: Query{
				AllowedIssuers: []string{issuerDID},
				Context:        "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
				Type:           "KYCEmployee",
				CredentialSubject: map[string]interface{}{
					"documentType": map[string]interface{}{
						"$eq": -1,
					},
				},
				SkipClaimRevocationCheck: false,
			},
			pubSig: &CircuitOutputs{
				IssuerID:            &issuerID,
				ClaimSchema:         KYCEmployeeSchema,
				Operator:            1,
				Value:               []*big.Int{big.NewInt(-1)},
				Merklized:           0,
				SlotIndex:           0,
				IsRevocationChecked: 0,
			},
			expErr: ErrNegativeValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.query.Check(context.Background(), &mockMemorySchemaLoader{}, tt.pubSig, nil)
			require.EqualError(t, err, tt.expErr.Error())
		})
	}
}

func TestVerifyQuery_Success(t *testing.T) {
	tests := []struct {
		name   string
		query  Query
		pubSig *CircuitOutputs
	}{
		{
			name: "Issuer DID is valid ID from public signals",
			query: Query{
				AllowedIssuers: []string{"did:polygonid:polygon:mumbai:2qHSHBGWGJ68AosMKcLCTp8FYdVrtYE6MtNHhq8xpK"},
			},
			pubSig: &CircuitOutputs{
				IssuerID: func() *core.ID {
					i, _ := big.NewInt(0).SetString("22638457188543025296541325416907897762715008870723718557276875842936181250", 10)
					userID, _ := core.IDFromInt(i)
					return &userID
				}(),
			},
		},
		{
			name: "All issuers are allowed",
			query: Query{
				AllowedIssuers: []string{"*"},
			},
			pubSig: &CircuitOutputs{
				IssuerID: func() *core.ID {
					i, _ := big.NewInt(0).SetString("22638457188543025296541325416907897762715008870723718557276875842936181250", 10)
					userID, _ := core.IDFromInt(i)
					return &userID
				}(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.query.verifyIssuer(tt.pubSig)
			require.NoError(t, err)
		})
	}
}

func TestVerifyQuery_Error(t *testing.T) {
	tests := []struct {
		name   string
		query  Query
		pubSig *CircuitOutputs
		err    error
	}{
		{
			name: "Invalid issuer in public signals",
			query: Query{
				AllowedIssuers: []string{"did:polygonid:polygon:mumbai:2qHSHBGWGJ68AosMKcLCTp8FYdVrtYE6MtNHhq8xpK"},
			},
			pubSig: &CircuitOutputs{
				IssuerID: func() *core.ID {
					i, _ := big.NewInt(0).SetString("42", 10)
					userID, _ := core.IDFromInt(i)
					return &userID
				}(),
			},
		},
		{
			name: "Issuer not found",
			query: Query{
				AllowedIssuers: []string{
					"did:polygonid:polygon:mumbai:2qMe71smt9D591WQdKvkbJBSQfEUQXtvyPmzoCqDd7",
					"did:polygonid:polygon:mumbai:2qPvUBjKWgqADsAr2diJabs1NkqNJii482E2y1ZciQ",
				},
			},
			pubSig: &CircuitOutputs{
				IssuerID: func() *core.ID {
					i, _ := big.NewInt(0).SetString("22638457188543025296541325416907897762715008870723718557276875842936181250", 10)
					userID, _ := core.IDFromInt(i)
					return &userID
				}(),
			},
			err: ErrUnavailableIssuer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.query.verifyIssuer(tt.pubSig)
			fmt.Println("err", err)
			require.Error(t, err)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			}
		})
	}
}
