package hash

import (
	sha2_256 "crypto/sha256"
	"hash"

	"github.com/btcsuite/btcutil/base58"
	sha3 "golang.org/x/crypto/sha3"
)

// Computes object identifier based on hash. For current version it starts with D.
// D means - use sha3_256(sha2_256(data))
func ComputeHashId(data []byte) string {
	if data == nil {
		return "D" + EncodeHash(ComputeHashFromBytes([]byte{}))
	} else {
		return "D" + EncodeHash(ComputeHashFromBytes(data))
	}
}

type Hash struct {
	hash hash.Hash
}

func InitHash() *Hash {
	//log.Printf("DEBUG: initing hash")
	return &Hash{
		hash: sha2_256.New(),
	}
}

func (hash *Hash) Update(data []byte) error {
	//log.Printf("DEBUG: hashing data %d", len(data))
	_, err := hash.hash.Write(data)
	return err
}

func (hash *Hash) FinalId() string {
	return EncodeHash(hash.FinalHash())
}

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

// Computes hash using combination sha2/sha3
// Output=sha3_256(sha2_256(data))
func ComputeHashFromBytes(data []byte) []byte {
	sha2Hash := sha2_256.Sum256(data)
	//fmt.Printf("sha2 hash %s\n", base58.Encode(sha2Hash[:]))
	sha3Hash := sha3.Sum256(sha2Hash[:])
	return sha3Hash[:]
}

//Encodes hash to string identifier
func EncodeHash(hash []byte) string {
	return "D" + base58.Encode(hash)
}
