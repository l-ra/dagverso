package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	handlers "github.com/gorilla/handlers"
	mux "github.com/gorilla/mux"
	dagversehash "github.com/l-ra/dagverso/common/hash"
)

/*
Req
Access-Control-Request-Method: POST
Access-Control-Request-Headers: X-PINGOTHER, Content-Type

Res
Access-Control-Allow-Origin: http://foo.example
Access-Control-Allow-Methods: POST, GET, OPTIONS
Access-Control-Allow-Headers: X-PINGOTHER, Content-Type
*/

func setCors(wr http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Origin") != "" {
		wr.Header().Add("Access-Control-Allow-Origin", req.Header.Get("Origin"))
	}
	if req.Header.Get("Access-Control-Request-Headers") != "" {
		wr.Header().Add("Access-Control-Allow-Headers", req.Header.Get("Access-Control-Request-Headers"))
	}
	wr.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
}

func failResponse(wr http.ResponseWriter, httpStatus int, reason string) {
	wr.WriteHeader(httpStatus)
	wr.Write([]byte("Error: "))
	wr.Write([]byte(reason))
	wr.Write([]byte("\n"))
}

func storePostTemp(inp io.Reader) (retFile *os.File, retErr error) {
	file, err := ioutil.TempFile("", "dagverso-temp")
	buffer := make([]byte, 255)
	log.Printf("DEBUG: temp file for post: %s", file.Name())
	defer func() {
		if file != nil && err != nil {
			file.Close()
			os.Remove(file.Name())
			retErr = err
			retFile = nil
		}
	}()

	for {
		nr, err := inp.Read(buffer)
		if err != nil {
			if err == io.EOF {
				_, err = retFile.Write(buffer[:nr])
			}
			retErr = err
			return
		}
		_, err = retFile.Write(buffer[:nr])
		if err != nil {
			retErr = err
			return
		}
	}
}

func processPostBody(inp *os.File) (retHashId string, retErr error) {
	inp.Seek(0, io.SeekStart)
	defer func() {
		if inp != nil {
			inp.Close()
			os.Remove(inp.Name())
		}
	}()

	hashId, err := computeHashId(inp)
	if err != nil {
		return "", err
	}
	return hashId, nil
}

func handleTargetedPost(wr http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	expectedHashId := vars["hashId"]
	if expectedHashId == "" {
		failResponse(wr, http.StatusBadRequest, "missing hash value in the targeted post")
		return
	}
	tmpFile, err := storePostTemp(req.Body)
	if err != nil {
		failResponse(wr, http.StatusInternalServerError, "failed to store post to temp")
	}
	hashId, err := processPostBody(tmpFile)
	if err != nil {
		failResponse(wr, http.StatusInternalServerError, err.Error())
	}
	wr.Header().Add("Content-type", "text/plain;charset=utf-8")
	if hashId != expectedHashId {
		failResponse(wr, http.StatusBadRequest, fmt.Sprintf("hashId mismatch. Target hasdId %s different from computed hashId: %s", expectedHashId, hashId))
		return
	}
}

func handleBlindPost(wr http.ResponseWriter, req *http.Request) {
	tmpFile, err := storePostTemp(req.Body)
	if err != nil {
		failResponse(wr, http.StatusInternalServerError, fmt.Sprintf("failed to store post to temp: %s", err.Error()))
	}
	hashId, err := processPostBody(tmpFile)
	if err != nil {
		failResponse(wr, http.StatusInternalServerError, fmt.Sprintf("failed to process post data: %s", err.Error()))
	}
	wr.Header().Add("Content-type", "text/plain;charset=utf-8")
	wr.Write([]byte(hashId))
	wr.WriteHeader(http.StatusOK)
}

func computeHashId(inp io.Reader) (string, error) {
	hash := dagversehash.InitHash()
	buffer := make([]byte, 255)
	for {
		nr, err := inp.Read(buffer)
		if err != nil {
			if err == io.EOF {
				hash.Update(buffer[:nr])
				return hash.FinalId(), nil
			}
			return "", err
		}
		hash.Update(buffer)
	}
}

func configure() {

}
func main() {
	configure()
	router := mux.NewRouter()
	apiV1Router := router.PathPrefix("/api/v1").Subrouter()

	apiV1Router.Methods(http.MethodPost).Path("/{hashId:[a-zA-Z0-9_-]{44}}").HandlerFunc(handleTargetedPost)
	apiV1Router.Methods(http.MethodPost).Path("/").HandlerFunc(handleBlindPost)

	http.ListenAndServe(":8087", handlers.CORS()(router))

}
