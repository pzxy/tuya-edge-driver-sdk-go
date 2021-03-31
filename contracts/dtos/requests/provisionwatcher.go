//
// Copyright (C) 2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package requests

import (
	"encoding/json"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/errors"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/models"
)

// AddProvisionWatcherRequest defines the Request Content for POST ProvisionWatcher DTO.
// This object and its properties correspond to the AddProvisionWatcherRequest object in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.x#/AddProvisionWatcherRequest
type AddProvisionWatcherRequest struct {
	common.BaseRequest `json:",inline"`
	ProvisionWatcher   dtos.ProvisionWatcher `json:"provisionWatcher"`
}

// Validate satisfies the Validator interface
func (pw AddProvisionWatcherRequest) Validate() error {
	err := contracts.Validate(pw)
	return err
}

// UnmarshalJSON implements the Unmarshaler interface for the AddProvisionWatcherRequest type
func (pw *AddProvisionWatcherRequest) UnmarshalJSON(b []byte) error {
	var alias struct {
		common.BaseRequest
		ProvisionWatcher dtos.ProvisionWatcher
	}
	if err := json.Unmarshal(b, &alias); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "Failed to unmarshal request body as JSON.", err)
	}

	*pw = AddProvisionWatcherRequest(alias)

	// validate AddDeviceRequest DTO
	if err := pw.Validate(); err != nil {
		return err
	}
	return nil
}

// AddProvisionWatcherReqToProvisionWatcherModels transforms the AddProvisionWatcherRequest DTO array to the ProvisionWatcher model array
func AddProvisionWatcherReqToProvisionWatcherModels(addRequests []AddProvisionWatcherRequest) (ProvisionWatchers []models.ProvisionWatcher) {
	for _, req := range addRequests {
		d := dtos.ToProvisionWatcherModel(req.ProvisionWatcher)
		ProvisionWatchers = append(ProvisionWatchers, d)
	}
	return ProvisionWatchers
}

// UpdateProvisionWatcherRequest defines the Request Content for PUT event as pushed DTO.
// This object and its properties correspond to the UpdateProvisionWatcherRequest object in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.x#/UpdateProvisionWatcherRequest
type UpdateProvisionWatcherRequest struct {
	common.BaseRequest `json:",inline"`
	ProvisionWatcher   dtos.UpdateProvisionWatcher `json:"provisionWatcher"`
}

// Validate satisfies the Validator interface
func (pw UpdateProvisionWatcherRequest) Validate() error {
	err := contracts.Validate(pw)
	return err
}

// UnmarshalJSON implements the Unmarshaler interface for the UpdateProvisionWatcherRequest type
func (pw *UpdateProvisionWatcherRequest) UnmarshalJSON(b []byte) error {
	var alias struct {
		common.BaseRequest
		ProvisionWatcher dtos.UpdateProvisionWatcher
	}
	if err := json.Unmarshal(b, &alias); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "Failed to unmarshal request body as JSON.", err)
	}

	*pw = UpdateProvisionWatcherRequest(alias)

	// validate UpdateDeviceRequest DTO
	if err := pw.Validate(); err != nil {
		return err
	}
	return nil
}

// ReplaceProvisionWatcherModelFieldsWithDTO replace existing ProvisionWatcher's fields with DTO patch
func ReplaceProvisionWatcherModelFieldsWithDTO(pw *models.ProvisionWatcher, patch dtos.UpdateProvisionWatcher) {
	if patch.Labels != nil {
		pw.Labels = patch.Labels
	}
	if patch.Identifiers != nil {
		pw.Identifiers = patch.Identifiers
	}
	if patch.BlockingIdentifiers != nil {
		pw.BlockingIdentifiers = patch.BlockingIdentifiers
	}
	if patch.ProfileName != nil {
		pw.ProfileName = *patch.ProfileName
	}
	if patch.ServiceName != nil {
		pw.ServiceName = *patch.ServiceName
	}
	if patch.AdminState != nil {
		pw.AdminState = models.AdminState(*patch.AdminState)
	}
	if patch.AutoEvents != nil {
		pw.AutoEvents = dtos.ToAutoEventModels(patch.AutoEvents)
	}
}
