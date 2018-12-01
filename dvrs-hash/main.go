package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/l-ra/dagverso/common/hash"
)

var outputEncoding string
var inputEncoding string
var convertHash bool

func init() {
	flag.StringVar(&outputEncoding, "encoding", "dvrs", "[hex|bin|dvrs] output encoding")
	flag.StringVar(&outputEncoding, "e", "dvrs", "-encoding shorthand")
	flag.BoolVar(&convertHash, "convert", false, "when present, exactly 1 non flag argument is expected and it is taken as dvrs formated hash a is converted")
	flag.BoolVar(&convertHash, "c", false, "-convert shorthand")
	flag.StringVar(&inputEncoding, "input", "dvrs", "[hex|dvrs] input encoding")
	flag.StringVar(&inputEncoding, "i", "dvrs", "-input shorthand")
}

func main() {
	flag.Parse()
	args := flag.Args()
	var err error
	var hashVal []byte
	reader := os.Stdin
	if convertHash {
		if len(args) == 0 {
			log.Fatalf("Input hash missing")
		}

		switch inputEncoding {
		case "hex":
			hashVal = make([]byte, hex.DecodedLen(len(args[0])))
			_, err = hex.Decode(hashVal, []byte(args[0]))
			break
		case "dvrs":
			hashVal, err = hash.DecodeHash(args[0])
			break
		default:
			fmt.Printf("Unknown encoding %s\n", outputEncoding)
		}

		if err != nil {
			log.Fatalf("Failed to decode input hash: %s", args[0])
		}
	} else {
		if len(args) > 0 {
			reader, err = os.Open(args[0])
			if err != nil {
				log.Fatalf("Error openning file %s: %s", args[0], err.Error())
			}
		}
		hashVal, err = hash.ComputeHashFromReader(reader)
		if err != nil {
			log.Fatalf("Error reading input %s", err.Error())
		}
	}
	switch outputEncoding {
	case "hex":
		dst := make([]byte, hex.EncodedLen(len(hashVal)))
		hex.Encode(dst, hashVal)
		fmt.Printf("%s\n", string(dst))
		break
	case "dvrs":
		fmt.Printf("%s\n", hash.EncodeHash(hashVal))
		break
	case "bin":
		os.Stdout.Write(hashVal)
		break
	default:
		fmt.Printf("Unknown encoding %s\n", outputEncoding)
	}
}
