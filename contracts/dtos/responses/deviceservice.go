//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package responses

import (
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/common"
)

// DeviceServiceResponse defines the Response Content for GET DeviceService DTOs.
// This object and its properties correspond to the DeviceServiceResponse object in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.x#/DeviceServiceResponse
type DeviceServiceResponse struct {
	common.BaseResponse `json:",inline"`
	Service             dtos.DeviceService `json:"service"`
}

func NewDeviceServiceResponse(requestId string, message string, statusCode int, deviceService dtos.DeviceService) DeviceServiceResponse {
	return DeviceServiceResponse{
		BaseResponse: common.NewBaseResponse(requestId, message, statusCode),
		Service:      deviceService,
	}
}

// MultiDeviceServicesResponse defines the Response Content for GET multiple DeviceService DTOs.
// This object and its properties correspond to the MultiDeviceServicesResponse object in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.x#/MultiDeviceServicesResponse
type MultiDeviceServicesResponse struct {
	common.BaseResponse `json:",inline"`
	Total               uint32               `json:"total"`
	Services            []dtos.DeviceService `json:"services"`
}

func NewMultiDeviceServicesResponse(requestId string, message string, statusCode int, deviceServices []dtos.DeviceService, total uint32) MultiDeviceServicesResponse {
	return MultiDeviceServicesResponse{
		BaseResponse: common.NewBaseResponse(requestId, message, statusCode),
		Total:        total,
		Services:     deviceServices,
	}
}
