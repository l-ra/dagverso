package schema

import (
	"log"

	"github.com/l-ra/dagverso/common/hash"

	"io/ioutil"
	"os"
	"path"
)

const DefaultSchemaCacheDir = ".dvrs/schema-cache"

var schemaCacheDir = ""

func initSchemaCacheDir() {
	if schemaCacheDir != "" {
		return
	}
	var success bool
	var homedir string
	schemaCacheDir, success = os.LookupEnv("DVRS_SCHEMA_CACHE_DIR")

	if !success {
		homedir, success = os.LookupEnv("HOME")
		if success {
			schemaCacheDir = path.Join(homedir, DefaultSchemaCacheDir)
		}
	}
	if !success {
		homedir, success = os.LookupEnv("USERPROFILE")
		if success {
			schemaCacheDir = path.Join(homedir, DefaultSchemaCacheDir)
		}
	}

	if !success {
		schemaCacheDir = path.Join(".", DefaultSchemaCacheDir)
	}

	log.Printf("INFO: using schema cache dir: %s", schemaCacheDir)
}

func GetSchemaFromCache(hash_ []byte) ([]byte, error) {
	initSchemaCacheDir()
	file, err := os.Open(path.Join(schemaCacheDir, hash.EncodeHashId(hash_)))
	if err != nil {
		return nil, err
	}
	defer func() {
		file.Close()
	}()
	data, err := ioutil.ReadAll(file)
	return data, err
}

func GetSchemaFromCacheById(hashId string) ([]byte, error) {
	hash_, err := hash.DecodeHashId(hashId)
	if err != nil {
		return nil, err
	}
	return GetSchemaFromCache(hash_)
}
