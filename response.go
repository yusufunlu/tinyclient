package tinyclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Response is a wrapper of http.Response which provide some extra features
type Response struct {
	client     *Client
	Request    *Request
	Response   *http.Response
	bodyBytes  []byte
	ReceivedAt time.Time
}

// ReadBody reads the http.Response bodyBytes and assigns it to the r.Body
func (response *Response) ReadBody() ([]byte, error) {

	// If r.bodyBytes already set then return r.bodyBytes
	if response.bodyBytes != nil {
		return response.bodyBytes, nil
	}

	// Check if Response.resp (*http.Response) is nil
	if response.Response == nil {
		err := fmt.Errorf("http.Response is nil")
		response.client.ErrorLogger.Println(err)
		return nil, err
	}

	// Check if Response.resp.Body (*http.Response.Body) is nil
	if response.Response.Body == nil {
		err := fmt.Errorf("http.Response's Body is nil")
		response.client.ErrorLogger.Println(err)
		return nil, err
	}

	// Read response bodyBytes
	b, err := ioutil.ReadAll(response.Response.Body)
	if err != nil {
		response.client.ErrorLogger.Printf("Can't read http.Response bodyBytes Error: %v!", err)
		return nil, err
	}

	// Set response readBody
	response.bodyBytes = b
	if len(response.bodyBytes) == 0 {
		response.client.InfoLogger.Println("Response body is empty")
	}

	// Close response bodyBytes
	err = response.Response.Body.Close()
	if err != nil {
		response.client.ErrorLogger.Printf("Can't close http.Response body Error: %v!", err)
		return nil, err
	}

	return b, nil

}

func (response *Response) BodyUnmarshall(v interface{}) error {
	resBody, err := response.ReadBody()

	if err != nil {
		return err
	}

	//if len(resBody) == 0 can be handled later
	err = json.Unmarshal(resBody, v)
	if err != nil {
		response.client.ErrorLogger.Println(err)
		return err
	}
	return nil
}
