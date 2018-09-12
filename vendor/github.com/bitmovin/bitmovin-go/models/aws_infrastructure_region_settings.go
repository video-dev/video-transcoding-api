package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type CreateAWSInfrastructureRegionSettingsRequest struct {
	LimitParallelEncodings                      *int     `json:"limitParallelEncodings,omitempty"`
	MaximumAmountOfCoordinatorAndWorkerInRegion *int     `json:"maximumAmountOfCoordinatorAndWorkerInRegion,omitempty"`
	MaxMoneyToSpentPerMonth                     *int     `json:"maxMoneyToSpentPerMonth,omitempty"`
	SecurityGroupId                             string   `json:"securityGroupId"`
	SubnetId                                    string   `json:"subnetId"`
	MachineTypes                                []string `json:"machineTypes,omitempty"`
	SshPort                                     *string  `json:"sshPort,omitempty"`
}

type AWSInfrastructureRegionSettingsDetail struct {
	LimitParallelEncodings                      int                          `json:"limitParallelEncodings"`
	MaximumAmountOfCoordinatorAndWorkerInRegion int                          `json:"maximumAmountOfCoordinatorAndWorkerInRegion"`
	MaxMoneyToSpentPerMonth                     int                          `json:"maxMoneyToSpentPerMonth"`
	SecurityGroupId                             string                       `json:"securityGroupId"`
	SubnetId                                    string                       `json:"subnetId"`
	MachineTypes                                []string                     `json:"machineTypes"`
	SshPort                                     string                       `json:"sshPort"`
	Region                                      bitmovintypes.AWSCloudRegion `json:"region"`
}

type AWSInfrastructureRegionSettingsResponseData struct {
	Result AWSInfrastructureRegionSettingsDetail `json:"result"`
}

type AWSInfrastructureRegionSettingsResponse struct {
	RequestID *string                                     `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus                `json:"status,omitempty"`
	Data      AWSInfrastructureRegionSettingsResponseData `json:"data,omitempty"`
}
type AWSInfrastructureRegionSettingsListResult struct {
	TotalCount *int64                                  `json:"totalCount,omitempty"`
	Previous   *string                                 `json:"previous,omitempty"`
	Next       *string                                 `json:"next,omitempty"`
	Items      []AWSInfrastructureRegionSettingsDetail `json:"items,omitempty"`
}
type AWSInfrastructureRegionSettingsListResponseData struct {
	Result AWSInfrastructureRegionSettingsListResult `json:"result,omitempty"`
}
type AWSInfrastructureRegionSettingsListResponse struct {
	RequestID *string                                         `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus                    `json:"status,omitempty"`
	Data      AWSInfrastructureRegionSettingsListResponseData `json:"data,omitempty"`
}
