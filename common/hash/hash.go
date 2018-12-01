package hash

import (
	sha2_256 "crypto/sha256"
	"fmt"
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
	return &Hash{
		hash: sha2_256.New(),
	}
}

func (hash *Hash) Update(data []byte) error {
	_, err := hash.hash.Write(data)
	return err
}

func (hash *Hash) FinalId() string {
	return EncodeHash(hash.FinalHash())
}

func (hash *Hash) FinalHash() []byte {
	sha2Hash := hash.hash.Sum(nil)
	fmt.Printf("---sha2 hash %s\n", base58.Encode(sha2Hash[:]))
	sha3Hash := sha3.Sum256(sha2Hash)
	return sha3Hash[:]
}

// Computes hash using combination sha2/sha3
// Output=sha3_256(sha2_256(data))
func ComputeHashFromBytes(data []byte) []byte {
	sha2Hash := sha2_256.Sum256(data)
	fmt.Printf("sha2 hash %s\n", base58.Encode(sha2Hash[:]))
	sha3Hash := sha3.Sum256(sha2Hash[:])
	return sha3Hash[:]
}

//Encodes hash to string identifier
func EncodeHash(hash []byte) string {
	return "D" + base58.Encode(hash)
}
