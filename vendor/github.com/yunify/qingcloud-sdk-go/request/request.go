// +-------------------------------------------------------------------------
// | Copyright (C) 2016 Yunify, Inc.
// +-------------------------------------------------------------------------
// | Licensed under the Apache License, Version 2.0 (the "License");
// | you may not use this work except in compliance with the License.
// | You may obtain a copy of the License in the LICENSE file, or at:
// |
// | http://www.apache.org/licenses/LICENSE-2.0
// |
// | Unless required by applicable law or agreed to in writing, software
// | distributed under the License is distributed on an "AS IS" BASIS,
// | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// | See the License for the specific language governing permissions and
// | limitations under the License.
// +-------------------------------------------------------------------------

package request

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"github.com/yunify/qingcloud-sdk-go/logger"
	"github.com/yunify/qingcloud-sdk-go/request/data"
	"github.com/yunify/qingcloud-sdk-go/utils"
)

// A Request can build, sign, send and unpack API request.
type Request struct {
	Operation *data.Operation
	Input     *reflect.Value
	Output    *reflect.Value

	HTTPRequest  *http.Request
	HTTPResponse *http.Response
}

// DefaultCredentialProxyHost is default credential proxy host
const DefaultCredentialProxyHost = "169.254.169.254"

// DefaultCredentialProxyPort is default credential proxy port
const DefaultCredentialProxyPort = 80

// DefaultCredentialProxyProtocol is default credential proxy protocol
const DefaultCredentialProxyProtocol = "http"

// DefaultCredentialProxyURI is default credential proxy URI
const DefaultCredentialProxyURI = "/latest/meta-data/security-credentials"

// TokenOutput is the structure of token when retrieving it
type TokenOutput struct {
	Jti          string `json:"jti"`
	Token        string `json:"id_token"`
	AccessKey    string `json:"access_key"`
	SecretAccess string `json:"secret_key"`
	Expiration   int64  `json:"expiration"`
	Action       string `json:"action,omitempty"`
	RetCode      string `json:"ret_code"`
}

// New create a Request from given Operation, Input and Output.
// It returns a Request.
func New(o *data.Operation, i data.Input, x interface{}) (*Request, error) {
	input := reflect.ValueOf(i)
	if input.Elem().IsValid() {
		err := i.Validate()
		if err != nil {
			return nil, err
		}
	}
	output := reflect.ValueOf(x)

	return &Request{
		Operation: o,
		Input:     &input,
		Output:    &output,
	}, nil
}

// Send sends API request.
// It returns error if error occurred.
func (r *Request) Send() error {
	err := r.check()
	if err != nil {
		return err
	}

	err = r.build()
	if err != nil {
		return err
	}

	err = r.sign()
	if err != nil {
		return err
	}

	err = r.send()
	if err != nil {
		return err
	}

	err = r.unpack()
	if err != nil {
		return err
	}

	return nil
}

func (r *Request) check() error {
	if r.Operation.Config.AccessKeyID == "" && r.Operation.Config.SecretAccessKey == "" || r.Operation.Config.URI == "/iam" && r.isTokenExpired() {
		t := TokenOutput{}
		err := t.GetToken(r.getCredentialProxyURL())

		if err != nil {
			return err
		}
		r.Operation.Config.AccessKeyID = t.AccessKey
		r.Operation.Config.SecretAccessKey = t.SecretAccess
		r.Operation.Config.URI = "/iam"
		r.Operation.Config.Token = t.Token
		r.Operation.Config.Expiration = t.Expiration
	}

	if r.Operation.Config.AccessKeyID == "" {
		return errors.New("access key not provided")
	}

	if r.Operation.Config.SecretAccessKey == "" {
		return errors.New("secret access key not provided")
	}

	return nil
}

func (r *Request) build() error {
	b := &Builder{}
	httpRequest, err := b.BuildHTTPRequest(r.Operation, r.Input)
	if err != nil {
		return err
	}

	r.HTTPRequest = httpRequest
	return nil
}

func (r *Request) sign() error {
	s := &Signer{
		AccessKeyID:     r.Operation.Config.AccessKeyID,
		SecretAccessKey: r.Operation.Config.SecretAccessKey,
	}
	err := s.WriteSignature(r.HTTPRequest)
	if err != nil {
		return err
	}

	return nil
}

func (r *Request) send() error {
	var response *http.Response
	var err error

	if r.Operation.Config.Connection == nil {
		return errors.New("connection not initialized")
	}

	retries := r.Operation.Config.ConnectionRetries + 1
	for {
		if retries > 0 {
			logger.Info(fmt.Sprintf(
				"Sending request: [%d] %s",
				utils.StringToUnixInt(r.HTTPRequest.Header.Get("Date"), "RFC 822"),
				r.HTTPRequest.Host))

			response, err = r.Operation.Config.Connection.Do(r.HTTPRequest)
			if err == nil {
				retries = 0
			} else {
				retries--
				time.Sleep(time.Second)
			}
		} else {
			break
		}
	}
	if err != nil {
		return err
	}

	r.HTTPResponse = response

	return nil
}

func (r *Request) unpack() error {
	u := &Unpacker{}
	err := u.UnpackHTTPRequest(r.Operation, r.HTTPResponse, r.Output)
	if err != nil {
		return err
	}

	return nil
}

// GetToken is used to get token from credential proxy server
func (t *TokenOutput) GetToken(credentialProxyURL string) error {
	response, err := http.Get(credentialProxyURL)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return err
	}

	_, err = utils.JSONDecode(content, t)

	return err
}

func (r *Request) isTokenExpired() bool {
	if r.Operation.Config.Token == "" {
		return true
	}

	now := time.Now().UTC().Unix()

	return now >= r.Operation.Config.Expiration
}

func (r *Request) getCredentialProxyURL() string {
	var credentialProxyProtocol string
	var credentialProxyHost string
	var credentialProxyPort int
	var credentialProxyURI string

	if r.Operation.Config.CredentialProxyProtocol != "" {
		credentialProxyProtocol = r.Operation.Config.CredentialProxyProtocol
	} else {
		credentialProxyProtocol = DefaultCredentialProxyProtocol
	}

	if r.Operation.Config.CredentialProxyHost != "" {
		credentialProxyHost = r.Operation.Config.CredentialProxyHost
	} else {
		credentialProxyHost = DefaultCredentialProxyHost
	}

	if r.Operation.Config.CredentialProxyPort != 0 {
		credentialProxyPort = r.Operation.Config.CredentialProxyPort
	} else {
		credentialProxyPort = DefaultCredentialProxyPort
	}

	if r.Operation.Config.CredentialProxyURI != "" {
		credentialProxyURI = r.Operation.Config.CredentialProxyURI
	} else {
		credentialProxyURI = DefaultCredentialProxyURI
	}

	credentialProxyURL := fmt.Sprintf("%s://%s:%d%s", credentialProxyProtocol, credentialProxyHost, credentialProxyPort, credentialProxyURI)

	return credentialProxyURL
}
