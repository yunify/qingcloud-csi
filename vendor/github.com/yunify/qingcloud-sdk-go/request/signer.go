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
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/yunify/qingcloud-sdk-go/logger"
	"github.com/yunify/qingcloud-sdk-go/utils"
)

// Signer is the http request signer for IaaS service.
type Signer struct {
	AccessKeyID     string
	SecretAccessKey string

	BuiltURL  string
	BuiltForm string
}

// WriteSignature calculates signature and write it to http request.
func (is *Signer) WriteSignature(request *http.Request) error {
	_, err := is.BuildSignature(request)
	if err != nil {
		return err
	}

	newRequest, err := http.NewRequest(request.Method,
		request.URL.Scheme+"://"+request.URL.Host+is.BuiltURL, strings.NewReader(is.BuiltForm))
	if err != nil {
		return err
	}
	request.URL = newRequest.URL
	request.Body = newRequest.Body

	logger.Info(fmt.Sprintf(
		"Signed QingCloud request: [%d] %s",
		utils.StringToUnixInt(request.Header.Get("Date"), "RFC 822"),
		request.URL.String()))

	return nil
}

// BuildSignature calculates the signature string.
func (is *Signer) BuildSignature(request *http.Request) (string, error) {
	stringToSign, err := is.BuildStringToSign(request)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, []byte(is.SecretAccessKey))
	h.Write([]byte(stringToSign))

	signature := strings.TrimSpace(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	signature = strings.Replace(signature, " ", "+", -1)
	signature = url.QueryEscape(signature)

	logger.Debug(fmt.Sprintf(
		"QingCloud signature: [%d] %s",
		utils.StringToUnixInt(request.Header.Get("Date"), "RFC 822"),
		signature))
	if request.Method == "GET" {
		is.BuiltURL += "&signature=" + signature
	} else if request.Method == "POST" {
		is.BuiltForm += "&signature=" + signature
	}

	return signature, nil
}

// BuildStringToSign build the string to sign.
func (is *Signer) BuildStringToSign(request *http.Request) (string, error) {
	if request.Method == "GET" {
		return is.BuildStringToSignByValues(request.Header.Get("Date"), request.Method, request.URL.Path, request.URL.Query())
	} else if request.Method == "POST" {
		return is.BuildStringToSignByValues(request.Header.Get("Date"), request.Method, request.URL.Path, request.Form)
	}
	return "", fmt.Errorf("Requset Type Not Support For Sign ")
}

// BuildStringToSignByValues build the string to sign.
func (is *Signer) BuildStringToSignByValues(requestDate string, requestMethod string, requestPath string, requestParams url.Values) (string, error) {
	requestParams.Set("access_key_id", is.AccessKeyID)
	requestParams.Set("signature_method", "HmacSHA256")
	requestParams.Set("signature_version", "1")

	var timeValue time.Time
	if requestDate != "" {
		var err error
		timeValue, err = utils.StringToTime(requestDate, "RFC 822")
		if err != nil {
			return "", err
		}
	} else {
		timeValue = time.Now()
	}
	requestParams.Set("time_stamp", utils.TimeToString(timeValue, "ISO 8601"))

	keys := []string{}
	for key := range requestParams {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	parts := []string{}
	for _, key := range keys {
		values := requestParams[key]
		if len(values) > 0 {
			if values[0] != "" {
				value := strings.TrimSpace(strings.Join(values, ""))
				value = url.QueryEscape(value)
				value = strings.Replace(value, "+", "%20", -1)
				parts = append(parts, key+"="+value)
			} else {
				parts = append(parts, key+"=")
			}
		} else {
			parts = append(parts, key+"=")
		}
	}

	urlParams := strings.Join(parts, "&")

	stringToSign := requestMethod + "\n" + requestPath + "\n" + urlParams

	logger.Debug(fmt.Sprintf(
		"QingCloud string to sign: %s",
		stringToSign))

	if requestMethod == "GET" {
		is.BuiltURL = requestPath + "?" + urlParams
		is.BuiltForm = ""
	} else if requestMethod == "POST" {
		is.BuiltURL = requestPath
		is.BuiltForm = urlParams
	}

	return stringToSign, nil
}
