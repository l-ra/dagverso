package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/l-ra/dagverso/common/hash"
	"github.com/l-ra/dagverso/common/schema"
	goavro "gopkg.in/linkedin/goavro.v2"
)

var title string
var creator string
var subject string
var description string
var publisher string
var contributor string
var date string
var type_ string
var format string
var identifier string
var source string
var language string
var coverage string
var rights string

var replaces_hash []byte
var entity_hash []byte
var relation_hash [][]byte

type DvrsHash struct {
	Hash []byte
}

func (v DvrsHash) String() string {
	if v.Hash != nil {
		return hash.EncodeHashRaw(v.Hash)
	}
	return ""
}

func (v DvrsHash) Set(s string) error {
	h, err := hash.DecodeHashRawWithErr(s)
	if err == nil {
		v.Hash = h
	}
	return nil
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s [flags] inputFile outdir:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&title, "title", "", "title according to dublin core")
	flag.StringVar(&creator, "creator", "", "creator according to dublin core")
	flag.StringVar(&subject, "subject", "", "subject according to dublin core")
	flag.StringVar(&description, "description", "", "description according to dublin core")
	flag.StringVar(&publisher, "publisher", "", "publisher according to dublin core")
	flag.StringVar(&contributor, "contributor", "", "contributor according to dublin core")
	flag.StringVar(&date, "date", "", "date according to dublin core")
	flag.StringVar(&type_, "type", "", "type according to dublin core")
	flag.StringVar(&format, "format", "", "format according to dublin core  - when empty file utility is called to get mime type (FIXME: no impl yet")
	flag.StringVar(&identifier, "identifier", "", "identifier according to dublin core")
	flag.StringVar(&source, "source", "", "source according to dublin core")
	flag.StringVar(&language, "language", "", "language according to dublin core")
	flag.StringVar(&coverage, "coverage", "", "coverage according to dublin core")
	flag.StringVar(&rights, "rights", "", "rights according to dublin core")

	flag.Var(DvrsHash{replaces_hash}, "replaces", "entity according to dublin core")
	flag.Var(DvrsHash{entity_hash}, "entity", "entity according to dublin core")

	//FIXME: not implemented yet
	//flag.StringVar(&relation, "relation", "", "relation according to dublin core")
}

func main() {
	flag.Parse()
	args := flag.Args()
	var err error

	envelopeSchema, err := schema.GetSchemaFromCache(hash.DecodeHashIdOrNil(schema.EnvelopeSchemaHashID))
	if err != nil {
		log.Fatalf("Failed to get envelope schema. Problem with schema cache? %s", err.Error())
	}

	basicmetaSchema, err := schema.GetSchemaFromCache(hash.DecodeHashIdOrNil(schema.BasicMetaSchemaHashID))
	if err != nil {
		log.Fatalf("Failed to get basicmeta schema. Problem with schema cache? %s", err.Error())
	}

	envelopeCodec, err := goavro.NewCodec(string(envelopeSchema))
	if err != nil {
		log.Fatalf("Failed to parse envelope schema. Corrupted schema cache? %s", err.Error())
	}

	basicmetaCodec, err := goavro.NewCodec(string(basicmetaSchema))
	if err != nil {
		log.Fatalf("Failed to parse basicmeta schema. Corrupted schema cache? %s", err.Error())
	}

	if len(args) < 2 {
		log.Fatalf("missing inputFile and outDir use -help flag to show usage")
	}

	inputFile := args[0]
	outDir := args[1]

	inp, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Faield to read input file: %s", err.Error())
	}
	defer func() { inp.Close() }()
	fileHash, err := hash.ComputeHashFromReader(inp)
	if err != nil {
		log.Fatalf("Failed to compute hash of input: %s", err.Error())
	}

	outFileVersionHash := hash.HashToVersionHash(fileHash)
	outFile := path.Join(outDir, hash.EncodeHashId(fileHash))
	outFileWr, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to write file %s", err.Error())
	}
	defer func() { outFileWr.Close() }()

	// copy stream to out file
	inp.Seek(0, 0)
	buf := make([]byte, 4096)
	for {
		nr, err := inp.Read(buf)
		if err != nil {
			if err == io.EOF {
				nw, err := outFileWr.Write(buf[:nr])
				if nw != nr {
					log.Fatalf("failed to write to out file %s", err.Error())
				}
				if err != nil {
					os.Remove(outFile)
					log.Fatalf("failed to write to out file %s", err.Error())
				}
				break
			}

			log.Fatalf("failed to write to out file %s", err.Error())
		}
		nw, err := outFileWr.Write(buf[:nr])
		if nw != nr {
			log.Fatalf("failed to write to out file %s", err.Error())
		}
		if err != nil {
			os.Remove(outFile)
			log.Fatalf("failed to write to out file %s", err.Error())
		}
	}

	if title == "" {
		title = path.Base(inputFile)
	}

	data := make(map[string]interface{})

	setMapString(data, "title", title)
	setMapString(data, "creator", creator)
	setMapString(data, "subject", subject)
	setMapString(data, "description", description)
	setMapString(data, "publisher", publisher)
	setMapString(data, "contributor", contributor)
	setMapString(data, "date", date)
	setMapString(data, "type", type_)
	setMapString(data, "format", format)
	setMapString(data, "identifier", identifier)
	setMapString(data, "source", source)
	setMapString(data, "language", language)
	setMapString(data, "coverage", coverage)
	setMapString(data, "rights", rights)

	data["relation_hash"] = nil
	data["replaces_hash"] = nil
	data["entity_hash"] = outFileVersionHash

	dataAvro, err := basicmetaCodec.BinaryFromNative(nil, data)
	if err != nil {
		log.Fatal("failed to avro serialize basic metadata: ", err.Error())
	}

	envelope := schema.CreateEnvelope(schema.BasicMetaSchemaHashID, dataAvro)
	envelopeAvro, err := envelopeCodec.BinaryFromNative(nil, envelope)
	if err != nil {
		log.Fatal("failed to avro serialize envelope", err.Error())
	}

	metaOutFile := path.Join(outDir, hash.EncodeHashId(hash.ComputeHashFromBytes(envelopeAvro)))
	err = ioutil.WriteFile(metaOutFile, envelopeAvro, os.ModePerm)
	if err != nil {
		log.Fatal("failed to write basic metadata")
	}

	log.Printf("INFO: All OK. data: %s meta: %s", outFile, metaOutFile)
}

func setMapString(m map[string]interface{}, key string, value string) {
	if value != "" {
		v := make(map[string]interface{})
		v["string"] = value
		m[key] = v
	} else {
		m[key] = nil
	}
}

func setMapIfc(m map[string]interface{}, key string, value interface{}) {
	if value != nil {
		v := make(map[string]interface{})
		v["bytes"] = value
		m[key] = v
	} else {
		m[key] = nil
	}
}
