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

package service

import (
	"fmt"
	"time"

	"github.com/yunify/qingcloud-sdk-go/config"
	"github.com/yunify/qingcloud-sdk-go/request"
	"github.com/yunify/qingcloud-sdk-go/request/data"
	"github.com/yunify/qingcloud-sdk-go/request/errors"
)

var _ fmt.State
var _ time.Time

type AccesskeyService struct {
	Config     *config.Config
	Properties *AccesskeyServiceProperties
}

type AccesskeyServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) Accesskey(zone string) (*AccesskeyService, error) {
	properties := &AccesskeyServiceProperties{
		Zone: &zone,
	}

	return &AccesskeyService{Config: s.Config, Properties: properties}, nil
}

func (s *AccesskeyService) DeleteAccessKeys(i *DeleteAccessKeysInput) (*DeleteAccessKeysOutput, error) {
	if i == nil {
		i = &DeleteAccessKeysInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteAccessKeys",
		RequestMethod: "GET",
	}

	x := &DeleteAccessKeysOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type DeleteAccessKeysInput struct {
	AccessKeys []*string `json:"access_keys" name:"access_keys" location:"params"` // Required
}

func (v *DeleteAccessKeysInput) Validate() error {

	if len(v.AccessKeys) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "AccessKeys",
			ParentName:    "DeleteAccessKeysInput",
		}
	}

	return nil
}

type DeleteAccessKeysOutput struct {
	Message    *string   `json:"message" name:"message"`
	AccessKeys []*string `json:"access_keys" name:"access_keys" location:"elements"`
	Action     *string   `json:"action" name:"action" location:"elements"`
	JobID      *string   `json:"job_id" name:"job_id" location:"elements"`
	RetCode    *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

func (s *AccesskeyService) DescribeAccessKeys(i *DescribeAccessKeysInput) (*DescribeAccessKeysOutput, error) {
	if i == nil {
		i = &DescribeAccessKeysInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeAccessKeys",
		RequestMethod: "GET",
	}

	x := &DescribeAccessKeysOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type DescribeAccessKeysInput struct {
	AccessKeys []*string `json:"access_keys" name:"access_keys" location:"params"`
	Limit      *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset     *int      `json:"offset" name:"offset" default:"0" location:"params"`
	Owner      *string   `json:"owner" name:"owner" location:"params"`
	SearchWord *string   `json:"search_word" name:"search_word" location:"params"`
	Status     []*string `json:"status" name:"status" location:"params"`
	Verbose    *int      `json:"verbose" name:"verbose" default:"0" location:"params"`
}

func (v *DescribeAccessKeysInput) Validate() error {

	return nil
}

type DescribeAccessKeysOutput struct {
	Message      *string      `json:"message" name:"message"`
	AccessKeySet []*AccessKey `json:"access_key_set" name:"access_key_set" location:"elements"`
	Action       *string      `json:"action" name:"action" location:"elements"`
	RetCode      *int         `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount   *int         `json:"total_count" name:"total_count" location:"elements"`
}
