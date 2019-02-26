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

type AppService struct {
	Config     *config.Config
	Properties *AppServiceProperties
}

type AppServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) App(zone string) (*AppService, error) {
	properties := &AppServiceProperties{
		Zone: &zone,
	}

	return &AppService{Config: s.Config, Properties: properties}, nil
}

// Documentation URL: https://docs.qingcloud.com/api/bot/DeployAppVersion.html
func (s *AppService) DeployAppVersion(i *DeployAppVersionInput) (*DeployAppVersionOutput, error) {
	if i == nil {
		i = &DeployAppVersionInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeployAppVersion",
		RequestMethod: "GET",
	}

	x := &DeployAppVersionOutput{}
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

type DeployAppVersionInput struct {
	AppID      *string `json:"app_id" name:"app_id" location:"params"`
	AppType    *string `json:"app_type" name:"app_type" location:"params"`
	ChargeMode *string `json:"charge_mode" name:"charge_mode" location:"params"`
	Conf       *string `json:"conf" name:"conf" location:"params"`
	Debug      *int    `json:"debug" name:"debug" location:"params"`
	Owner      *string `json:"owner" name:"owner" location:"params"`
	VersionID  *string `json:"version_id" name:"version_id" location:"params"`
}

func (v *DeployAppVersionInput) Validate() error {

	return nil
}

type DeployAppVersionOutput struct {
	Message     *string   `json:"message" name:"message"`
	Action      *string   `json:"action" name:"action" location:"elements"`
	AppID       *string   `json:"app_id" name:"app_id" location:"elements"`
	AppVersion  *string   `json:"app_version" name:"app_version" location:"elements"`
	ClusterID   *string   `json:"cluster_id" name:"cluster_id" location:"elements"`
	ClusterName *string   `json:"cluster_name" name:"cluster_name" location:"elements"`
	JobID       *string   `json:"job_id" name:"job_id" location:"elements"`
	NodeCount   *int      `json:"node_count" name:"node_count" location:"elements"`
	NodeIDs     []*string `json:"node_ids" name:"node_ids" location:"elements"`
	RetCode     *int      `json:"ret_code" name:"ret_code" location:"elements"`
	VxNetID     *string   `json:"vxnet_id" name:"vxnet_id" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/bot/describe_app_version_attachments.html
func (s *AppService) DescribeAppVersionAttachments(i *DescribeAppVersionAttachmentsInput) (*DescribeAppVersionAttachmentsOutput, error) {
	if i == nil {
		i = &DescribeAppVersionAttachmentsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeAppVersionAttachments",
		RequestMethod: "GET",
	}

	x := &DescribeAppVersionAttachmentsOutput{}
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

type DescribeAppVersionAttachmentsInput struct {
	AttachmentIDs []*string `json:"attachment_ids" name:"attachment_ids" location:"params"`
	// ContentKeys's available values: config.json, locale/zh-cn.json, locale/en.json, cluster.json.mustache
	ContentKeys []*string `json:"content_keys" name:"content_keys" location:"params"`
	VersionID   *string   `json:"version_id" name:"version_id" location:"params"`
}

func (v *DescribeAppVersionAttachmentsInput) Validate() error {

	return nil
}

type DescribeAppVersionAttachmentsOutput struct {
	Message    *string                 `json:"message" name:"message"`
	Action     *string                 `json:"action" name:"action" location:"elements"`
	RetCode    *int                    `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount *int                    `json:"total_count" name:"total_count" location:"elements"`
	VersionSet []*AppVersionAttachment `json:"version_set" name:"version_set" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/bot/describe_app_versions.html
func (s *AppService) DescribeAppVersions(i *DescribeAppVersionsInput) (*DescribeAppVersionsOutput, error) {
	if i == nil {
		i = &DescribeAppVersionsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeAppVersions",
		RequestMethod: "GET",
	}

	x := &DescribeAppVersionsOutput{}
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

type DescribeAppVersionsInput struct {
	AppIDs  []*string `json:"app_ids" name:"app_ids" location:"params"`
	Limit   *int      `json:"limit" name:"limit" location:"params"`
	Name    *string   `json:"name" name:"name" location:"params"`
	Offset  *int      `json:"offset" name:"offset" location:"params"`
	Owner   *string   `json:"owner" name:"owner" location:"params"`
	Reverse *string   `json:"reverse" name:"reverse" location:"params"`
	SortKey *string   `json:"sort_key" name:"sort_key" location:"params"`
	Status  []*string `json:"status" name:"status" location:"params"`
	// Verbose's available values: 1, 0
	Verbose    *int      `json:"verbose" name:"verbose" location:"params"`
	VersionIDs []*string `json:"version_ids" name:"version_ids" location:"params"`
}

func (v *DescribeAppVersionsInput) Validate() error {

	if v.Verbose != nil {
		verboseValidValues := []string{"1", "0"}
		verboseParameterValue := fmt.Sprint(*v.Verbose)

		verboseIsValid := false
		for _, value := range verboseValidValues {
			if value == verboseParameterValue {
				verboseIsValid = true
			}
		}

		if !verboseIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Verbose",
				ParameterValue: verboseParameterValue,
				AllowedValues:  verboseValidValues,
			}
		}
	}

	return nil
}

type DescribeAppVersionsOutput struct {
	Message    *string       `json:"message" name:"message"`
	Action     *string       `json:"action" name:"action" location:"elements"`
	RetCode    *int          `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount *int          `json:"total_count" name:"total_count" location:"elements"`
	VersionSet []*AppVersion `json:"version_set" name:"version_set" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/bot/describe_apps.html
func (s *AppService) DescribeApps(i *DescribeAppsInput) (*DescribeAppsOutput, error) {
	if i == nil {
		i = &DescribeAppsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeApps",
		RequestMethod: "GET",
	}

	x := &DescribeAppsOutput{}
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

type DescribeAppsInput struct {
	App        *string   `json:"app" name:"app" location:"params"`
	AppName    *string   `json:"app_name" name:"app_name" location:"params"`
	AppType    []*string `json:"app_type" name:"app_type" location:"params"`
	Category   *string   `json:"category" name:"category" location:"params"`
	Limit      *int      `json:"limit" name:"limit" location:"params"`
	Offset     *int      `json:"offset" name:"offset" location:"params"`
	SearchWord *string   `json:"search_word" name:"search_word" location:"params"`
	Status     []*string `json:"status" name:"status" location:"params"`
	Tags       []*string `json:"tags" name:"tags" location:"params"`
	// Verbose's available values: 1, 0
	Verbose *int      `json:"verbose" name:"verbose" location:"params"`
	Zones   []*string `json:"zones" name:"zones" location:"params"`
}

func (v *DescribeAppsInput) Validate() error {

	if v.Verbose != nil {
		verboseValidValues := []string{"1", "0"}
		verboseParameterValue := fmt.Sprint(*v.Verbose)

		verboseIsValid := false
		for _, value := range verboseValidValues {
			if value == verboseParameterValue {
				verboseIsValid = true
			}
		}

		if !verboseIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Verbose",
				ParameterValue: verboseParameterValue,
				AllowedValues:  verboseValidValues,
			}
		}
	}

	return nil
}

type DescribeAppsOutput struct {
	Message    *string `json:"message" name:"message"`
	Action     *string `json:"action" name:"action" location:"elements"`
	AppSet     []*App  `json:"app_set" name:"app_set" location:"elements"`
	RetCode    *int    `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount *int    `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/bot/describe_app_version_attachments.html
func (s *AppService) GetGlobalUniqueId(i *GetGlobalUniqueIdInput) (*GetGlobalUniqueIdOutput, error) {
	if i == nil {
		i = &GetGlobalUniqueIdInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "GetGlobalUniqueId",
		RequestMethod: "GET",
	}

	x := &GetGlobalUniqueIdOutput{}
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

type GetGlobalUniqueIdInput struct {
	UserID *string `json:"user_id" name:"user_id" location:"params"`
}

func (v *GetGlobalUniqueIdInput) Validate() error {

	return nil
}

type GetGlobalUniqueIdOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
	UUID    *string `json:"uuid" name:"uuid" location:"elements"`
}
