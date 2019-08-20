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

type ProjectService struct {
	Config     *config.Config
	Properties *ProjectServiceProperties
}

type ProjectServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) Project(zone string) (*ProjectService, error) {
	properties := &ProjectServiceProperties{
		Zone: &zone,
	}

	return &ProjectService{Config: s.Config, Properties: properties}, nil
}

func (s *ProjectService) AddProjectResourceItems(i *AddProjectResourceItemsInput) (*AddProjectResourceItemsOutput, error) {
	if i == nil {
		i = &AddProjectResourceItemsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "AddProjectResourceItems",
		RequestMethod: "GET",
	}

	x := &AddProjectResourceItemsOutput{}
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

type AddProjectResourceItemsInput struct {
	ProjectID *string   `json:"project_id" name:"project_id" location:"params"` // Required
	Resources []*string `json:"resources" name:"resources" location:"params"`   // Required
}

func (v *AddProjectResourceItemsInput) Validate() error {

	if v.ProjectID == nil {
		return errors.ParameterRequiredError{
			ParameterName: "ProjectID",
			ParentName:    "AddProjectResourceItemsInput",
		}
	}

	if len(v.Resources) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Resources",
			ParentName:    "AddProjectResourceItemsInput",
		}
	}

	return nil
}

type AddProjectResourceItemsOutput struct {
	Message     *string   `json:"message" name:"message"`
	Action      *string   `json:"action" name:"action" location:"elements"`
	ProjectID   *string   `json:"project_id" name:"project_id" location:"elements"`
	ResourceIDs []*string `json:"resource_ids" name:"resource_ids" location:"elements"`
	RetCode     *int      `json:"ret_code" name:"ret_code" location:"elements"`
	ZoneID      *string   `json:"zone_id" name:"zone_id" location:"elements"`
}

func (s *ProjectService) DeleteProjectResourceItems(i *DeleteProjectResourceItemsInput) (*DeleteProjectResourceItemsOutput, error) {
	if i == nil {
		i = &DeleteProjectResourceItemsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteProjectResourceItems",
		RequestMethod: "GET",
	}

	x := &DeleteProjectResourceItemsOutput{}
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

type DeleteProjectResourceItemsInput struct {
	ProjectID []*string `json:"project_id" name:"project_id" location:"params"` // Required
	Resources []*string `json:"resources" name:"resources" location:"params"`   // Required
}

func (v *DeleteProjectResourceItemsInput) Validate() error {

	if len(v.ProjectID) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "ProjectID",
			ParentName:    "DeleteProjectResourceItemsInput",
		}
	}

	if len(v.Resources) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Resources",
			ParentName:    "DeleteProjectResourceItemsInput",
		}
	}

	return nil
}

type DeleteProjectResourceItemsOutput struct {
	Message     *string   `json:"message" name:"message"`
	Action      *string   `json:"action" name:"action" location:"elements"`
	ProjectID   []*string `json:"project_id" name:"project_id" location:"elements"`
	ResourceIDs []*string `json:"resource_ids" name:"resource_ids" location:"elements"`
	RetCode     *int      `json:"ret_code" name:"ret_code" location:"elements"`
	ZoneID      *string   `json:"zone_id" name:"zone_id" location:"elements"`
}

func (s *ProjectService) DescribeProjectResourceItems(i *DescribeProjectResourceItemsInput) (*DescribeProjectResourceItemsOutput, error) {
	if i == nil {
		i = &DescribeProjectResourceItemsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeProjectResourceItems",
		RequestMethod: "GET",
	}

	x := &DescribeProjectResourceItemsOutput{}
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

type DescribeProjectResourceItemsInput struct {
	InGlobal      *int      `json:"in_global" name:"in_global" location:"params"`
	Limit         *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset        *int      `json:"offset" name:"offset" default:"0" location:"params"`
	Owner         *string   `json:"owner" name:"owner" location:"params"`
	ProjectIDs    []*string `json:"project_ids" name:"project_ids" location:"params"`
	Reserve       *int      `json:"reserve" name:"reserve" location:"params"`
	ResourceTypes []*string `json:"resource_types" name:"resource_types" location:"params"`
	Resources     []*string `json:"resources" name:"resources" location:"params"`
	SortKey       *string   `json:"sort_key" name:"sort_key" location:"params"`
	Verbose       *int      `json:"verbose" name:"verbose" location:"params"`
}

func (v *DescribeProjectResourceItemsInput) Validate() error {

	return nil
}

type DescribeProjectResourceItemsOutput struct {
	Message                *string                `json:"message" name:"message"`
	Action                 *string                `json:"action" name:"action" location:"elements"`
	ProjectResourceItemSet []*ProjectResourceItem `json:"project_resource_item_set" name:"project_resource_item_set" location:"elements"`
	RetCode                *int                   `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount             *int                   `json:"total_count" name:"total_count" location:"elements"`
}

func (s *ProjectService) DescribeProjects(i *DescribeProjectsInput) (*DescribeProjectsOutput, error) {
	if i == nil {
		i = &DescribeProjectsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeProjects",
		RequestMethod: "GET",
	}

	x := &DescribeProjectsOutput{}
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

type DescribeProjectsInput struct {
	Limit      *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset     *int      `json:"offset" name:"offset" default:"0" location:"params"`
	Owner      *string   `json:"owner" name:"owner" location:"params"`
	ProjectIDs []*string `json:"project_ids" name:"project_ids" location:"params"`
	Shared     *string   `json:"shared" name:"shared" default:"False" location:"params"`
	Status     []*string `json:"status" name:"status" location:"params"`
}

func (v *DescribeProjectsInput) Validate() error {

	return nil
}

type DescribeProjectsOutput struct {
	Message    *string    `json:"message" name:"message"`
	Action     *string    `json:"action" name:"action" location:"elements"`
	ProjectSet []*Project `json:"project_set" name:"project_set" location:"elements"`
	RetCode    *int       `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount *int       `json:"total_count" name:"total_count" location:"elements"`
}
