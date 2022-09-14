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

type NotificationService struct {
	Config     *config.Config
	Properties *NotificationServiceProperties
}

type NotificationServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) Notification(zone string) (*NotificationService, error) {
	properties := &NotificationServiceProperties{
		Zone: &zone,
	}

	return &NotificationService{Config: s.Config, Properties: properties}, nil
}

func (s *NotificationService) DescribeNotificationLists(i *DescribeNotificationListsInput) (*DescribeNotificationListsOutput, error) {
	if i == nil {
		i = &DescribeNotificationListsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeNotificationLists",
		RequestMethod: "GET",
	}

	x := &DescribeNotificationListsOutput{}
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

type DescribeNotificationListsInput struct {
	Limit             *int      `json:"limit" name:"limit" default:"10" location:"params"`
	NotificationLists []*string `json:"notification_lists" name:"notification_lists" location:"params"` // Required
	Offset            *int      `json:"offset" name:"offset" default:"0" location:"params"`
	Owner             *string   `json:"owner" name:"owner" location:"params"`
}

func (v *DescribeNotificationListsInput) Validate() error {

	if len(v.NotificationLists) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "NotificationLists",
			ParentName:    "DescribeNotificationListsInput",
		}
	}

	return nil
}

type DescribeNotificationListsOutput struct {
	Message             *string             `json:"message" name:"message"`
	Action              *string             `json:"action" name:"action" location:"elements"`
	NotificationListSet []*NotificationList `json:"notification_list_set" name:"notification_list_set" location:"elements"`
	RetCode             *int                `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount          *int                `json:"total_count" name:"total_count" location:"elements"`
}

func (s *NotificationService) SendAlarmNotification(i *SendAlarmNotificationInput) (*SendAlarmNotificationOutput, error) {
	if i == nil {
		i = &SendAlarmNotificationInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "SendAlarmNotification",
		RequestMethod: "GET",
	}

	x := &SendAlarmNotificationOutput{}
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

type SendAlarmNotificationInput struct {
	NotificationData   []*NotificationData `json:"notification_data" name:"notification_data" location:"params"`       // Required
	NotificationListID *string             `json:"notification_list_id" name:"notification_list_id" location:"params"` // Required
	ResourceID         *string             `json:"resource_id" name:"resource_id" location:"params"`
	ResourceName       *string             `json:"resource_name" name:"resource_name" location:"params"`
	ResourceType       *string             `json:"resource_type" name:"resource_type" location:"params"`
	UserID             *string             `json:"user_id" name:"user_id" location:"params"` // Required
}

func (v *SendAlarmNotificationInput) Validate() error {

	if len(v.NotificationData) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "NotificationData",
			ParentName:    "SendAlarmNotificationInput",
		}
	}

	if len(v.NotificationData) > 0 {
		for _, property := range v.NotificationData {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.NotificationListID == nil {
		return errors.ParameterRequiredError{
			ParameterName: "NotificationListID",
			ParentName:    "SendAlarmNotificationInput",
		}
	}

	if v.UserID == nil {
		return errors.ParameterRequiredError{
			ParameterName: "UserID",
			ParentName:    "SendAlarmNotificationInput",
		}
	}

	return nil
}

type SendAlarmNotificationOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}
