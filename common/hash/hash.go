// Hash as currently used is cumputed from the bytes stored in DagVerso
// using sha3_256(sha256(bytes)) producing 32 bytes hash
// the hash is prepened with 0x42 as a versioning mark
//
// base58 (bitcoin flavor) is used for string representation of the hash
// func DagVersoHash(bytes) { returns sha3_256(sha256(bytes)) }
// hash = sha3_256(sha256(bytes))
// versionHash = 0x43|hash
// hashId - base58(versionHash)
package hash

import (
	sha2_256 "crypto/sha256"
	"errors"
	"hash"
	"io"

	"github.com/btcsuite/btcutil/base58"
	sha3 "golang.org/x/crypto/sha3"
)

// DagVerso version mark
const DvrsVersion = 0x42

// DagVerso Hash implementation
type Hash struct {
	hash hash.Hash
}

// Creates new DagVersohash
func New() *Hash {
	return &Hash{
		hash: sha2_256.New(),
	}
}

// updates current state of the hash
func (hash *Hash) Update(data []byte) error {
	//log.Printf("DEBUG: hashing data %d", len(data))
	_, err := hash.hash.Write(data)
	return err
}

// generates hashId after all data are passed to update
func (hash *Hash) FinalId() string {
	return EncodeHashId(hash.FinalHash())
}

// generates hash after all data are passed to update
func (hash *Hash) FinalHash() []byte {
	//log.Printf("DEBUG: finishing hash")
	sha2Hash := hash.hash.Sum(nil)
	//dst := make([]byte, hex.EncodedLen(len(sha2Hash)))
	//hex.Encode(dst, sha2Hash[:])
	//fmt.Printf("---sha2 hash b58 %s\n", base58.Encode(sha2Hash[:]))
	//fmt.Printf("---sha2 hash hex %s\n", dst)
	sha3Hash := sha3.Sum256(sha2Hash)
	return sha3Hash[:]
}

// computes dagverso hash for bytes in memory
func ComputeHashFromBytes(data []byte) []byte {
	sha2Hash := sha2_256.Sum256(data)
	//fmt.Printf("sha2 hash %s\n", base58.Encode(sha2Hash[:]))
	sha3Hash := sha3.Sum256(sha2Hash[:])
	return sha3Hash[:]
}

// computes dagverso hash for input reader
func ComputeHashFromReader(inp io.Reader) ([]byte, error) {
	hash := New()
	buffer := make([]byte, 4096)
	for {
		nr, err := inp.Read(buffer)
		if err != nil {
			if err == io.EOF {
				hash.Update(buffer[:nr])
				return hash.FinalHash(), nil
			}
			return nil, err
		}
		hash.Update(buffer[:nr])
	}
}

// hash -> versionHash
func HashToVersionHash(hash []byte) []byte {
	final := make([]byte, len(hash)+1)
	copy(final[1:], hash)
	final[0] = DvrsVersion
	return final
}

func VersionHashToHash(versionHash []byte) []byte {
	return versionHash[1:]
}

// versionHash -> hashId, panick when version mark not present
func EncodeHashRaw(versionHash []byte) string {
	if versionHash[0] != DvrsVersion {
		panic("failed to encode raw hash - missing version byte")
	}
	return base58.Encode(versionHash)
}

// hashId -> versionHash returns error when version mark not present
func DecodeHashRawWithErr(hashId string) ([]byte, error) {
	versionHash := base58.Decode(hashId)
	if versionHash[0] != DvrsVersion {
		return nil, errors.New("failed to decode raw hash - missing version byte")
	}
	return versionHash, nil
}

// hashId -> versionHash returns panicks when version mark not present
func DecodeHashRaw(hashId string) []byte {
	versionHash, err := DecodeHashRawWithErr(hashId)
	if err != nil {
		panic(err.Error())
	}
	return versionHash
}

//Encodes hash -> hashId
func EncodeHashId(hash []byte) string {
	final := make([]byte, len(hash)+1)
	copy(final[1:], hash)
	final[0] = DvrsVersion
	return base58.Encode(final)
}

//Decodes hashId -> hash, reports erro in case of wrong encoding
func DecodeHashId(hashId string) ([]byte, error) {
	decoded := base58.Decode(hashId)
	if decoded[0] != DvrsVersion {
		return nil, errors.New("bad dvrs hash version byte. decoded hashId shoild start with 0x42")
	}
	return decoded[1:], nil
}

//Decodes hashId -> hash, returns nil in case of error
func DecodeHashIdOrNil(hashId string) []byte {
	ret, err := DecodeHashId(hashId)
	if err != nil {
		return nil
	}
	return ret
}
