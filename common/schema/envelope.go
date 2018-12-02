package schema

import (
	"github.com/l-ra/dagverso/common/hash"
)

var EnvelopeMagic = []byte{0x64, 0x76, 0x72, 0x73}

func CreateEnvelope(schemaHashId string, data []byte) map[string]interface{} {
	envelope := make(map[string]interface{})
	envelope["magic"] = EnvelopeMagic
	envelope["schema_hash"] = hash.DecodeHashRaw(schemaHashId)
	envelope["data"] = data
	return envelope
}
