package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"ogomez/mkt-export/pkg/config"
	"path/filepath"
	"strings"
)

type RestClient struct {
	Client http.Client
	Bearer string
}

func New(url string, credentials config.Credentials) *RestClient {
	var client *http.Client
	tls := credentials.Certificates != config.Certificates{}
	if tls {
		client = &http.Client{
			Transport: getTransport(credentials.Certificates),
		}
	} else {
		client = &http.Client{}
	}
	bearer := ""
	keySecretAuth := (credentials.Key != "") && (credentials.Secret != "")
	if keySecretAuth {
		user := credentials.Key + ":" + credentials.Secret
		bearer = b64.StdEncoding.EncodeToString([]byte(user))
	}
	bearerAuth := credentials.Bearer != ""
	if bearerAuth {
		bearer = credentials.Bearer
	}

	return &RestClient{
		Client: *client,
		Bearer: bearer,
	}
}

// POST request
func (RClient *RestClient) Post(requestURL string, requestBody []byte, headers map[string]string) ([]interface{}, error) {

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req = setHeaders(req, headers)
	resp, err := RClient.buildArrayRequest(req)
	if err != nil {
		log.Println(fmt.Sprint(err))
		return nil, err
	}

	return resp, nil
}

func (RClient *RestClient) PostMultipart(requestURL string, headers map[string]string, requestBody []byte, files map[string]string) ([]interface{}, error) {

	body := &bytes.Buffer{}
	// Creates a new multipart Writer with a random boundary
	// writing to the empty buffer
	writer := multipart.NewWriter(body)

	bodyPart, err := writer.CreateFormField("eventRegistration")
  bodyString := string(requestBody)
  fmt.Printf("bodyString: %v\n", bodyString)
  io.Copy(bodyPart, strings.NewReader(bodyString))

  // Write the part body
	if schemaFilepath, ok := files["schema"]; ok {
		_, schemaFileName := filepath.Split(schemaFilepath)
		schemaFile, err := ioutil.ReadFile(schemaFilepath)
		if err != nil {
			log.Println(fmt.Sprint(err))
			return nil, err
		}
		schemaFilePart, err := writer.CreateFormFile("fileDocumentation", schemaFileName)
		if schemaFile != nil {
			io.Copy(schemaFilePart, bytes.NewReader(schemaFile))
		}
	} else {
		log.Println("No Documentation file present")
		writer.CreateFormField("fileDocumentation")
	}

	if exampleFilepath, ok := files["example"]; ok {
		_, exampleFileName := filepath.Split(exampleFilepath)
		exampleFile, err := ioutil.ReadFile(exampleFilepath)
		if err != nil {
			log.Println(fmt.Sprint(err))
			return nil, err
		}
		exampleFilePart, err := writer.CreateFormFile("fileExample", exampleFileName)
		if exampleFile != nil {
			io.Copy(exampleFilePart, bytes.NewReader(exampleFile))
		}
	} else {
		log.Println("No Example File present")
		writer.CreateFormField("fileExample")
	}
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(body.Bytes()))
	if err != nil {
		log.Println(fmt.Sprint(err))
		return nil, err
	}

	req = setHeaders(req, headers)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	log.Println(req.Header)
	resp, err := RClient.buildArrayRequest(req)
	if err != nil {
		log.Println(fmt.Println(err))
		return nil, err
	}

	return resp, nil
}

func (RClient *RestClient) Get(requestURL string, headers map[string]string) (interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req = setHeaders(req, headers)
	resp, err := RClient.buildRequest(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Get Request
// Expect results --> data:[]
func (RClient *RestClient) GetList(requestURL string, headers map[string]string) ([]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req = setHeaders(req, headers)
	resp, err := RClient.buildArrayRequest(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (RClient *RestClient) buildArrayRequest(req *http.Request) ([]interface{}, error) {
	result, err := RClient.build(req)
	if err != nil {
		return nil, err
	}
	switch v := result.(type) {
	case map[string]interface{}:
		if v["data"] != nil {
			return v["data"].([]interface{}), nil
		}
	default:
		return result.([]interface{}), nil
	}

	return nil, errors.New("No data result")
}

func (RClient *RestClient) buildRequest(req *http.Request) (interface{}, error) {
	return RClient.build(req)
}

// Build request - Client Do
func (RClient *RestClient) build(req *http.Request) (interface{}, error) {
	res, err := RClient.Client.Do(req)
	if err != nil {
		log.Printf("Rest client: error making http request: %s\n", err)
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		errorString := fmt.Sprintf("Rest client:: %d - %s : %v \n", res.StatusCode, req.Method, req.URL)
		return nil, errors.New(errorString)
	}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Rest client:: could not read response body: %s\n", err)
		return nil, err
	}

	var result interface{}
	json.Unmarshal([]byte(resBody), &result)
	return result, nil
}

// Get Transport from certificates
func getTransport(certificates config.Certificates) *http.Transport {
	certFile := certificates.CertFile
	keyFile := certificates.KeyFile
	caFile := certificates.CAFile

	// Load client cert
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Printf("Error loading cert files")
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Printf("Error reading the CA cert")
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	return &http.Transport{TLSClientConfig: tlsConfig}

}

func setHeaders(req *http.Request, headers map[string]string) *http.Request {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req
}
