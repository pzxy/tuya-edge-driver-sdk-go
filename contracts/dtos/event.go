//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package dtos

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/models"
)

// Event and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-data/2.x#/Event
type Event struct {
	common.Versionable `json:",inline"`
	Id                 string            `json:"id" validate:"required,uuid"`
	DeviceName         string            `json:"deviceName" validate:"required,edgex-dto-rfc3986-unreserved-chars"`
	ProfileName        string            `json:"profileName" validate:"required,edgex-dto-rfc3986-unreserved-chars"`
	Created            int64             `json:"created"`
	Origin             int64             `json:"origin" validate:"required"`
	Readings           []BaseReading     `json:"readings" validate:"gt=0,dive,required"`
	Tags               map[string]string `json:"tags,omitempty" xml:"-"` // Have to ignore since map not supported for XML
}

// NewEvent creates and returns an initialized Event with no Readings
func NewEvent(profileName string, deviceName string) Event {
	return Event{
		Versionable: common.NewVersionable(),
		Id:          uuid.NewString(),
		DeviceName:  deviceName,
		ProfileName: profileName,
		Origin:      time.Now().UnixNano(),
	}
}

// FromEventModelToDTO transforms the Event Model to the Event DTO
func FromEventModelToDTO(event models.Event) Event {
	var readings []BaseReading
	for _, reading := range event.Readings {
		readings = append(readings, FromReadingModelToDTO(reading))
	}

	tags := make(map[string]string)
	for tag, value := range event.Tags {
		tags[tag] = value
	}

	return Event{
		Versionable: common.NewVersionable(),
		Id:          event.Id,
		DeviceName:  event.DeviceName,
		ProfileName: event.ProfileName,
		Created:     event.Created,
		Origin:      event.Origin,
		Readings:    readings,
		Tags:        tags,
	}
}

// AddSimpleReading adds a simple reading to the Event
func (e *Event) AddSimpleReading(resourceName string, valueType string, value interface{}) error {
	reading, err := NewSimpleReading(e.ProfileName, e.DeviceName, resourceName, valueType, value)
	if err != nil {
		return err
	}
	e.Readings = append(e.Readings, reading)
	return nil
}

// AddBinaryReading adds a binary reading to the Event
func (e *Event) AddBinaryReading(resourceName string, binaryValue []byte, mediaType string) {
	e.Readings = append(e.Readings, NewBinaryReading(e.ProfileName, e.DeviceName, resourceName, binaryValue, mediaType))
}

// ToXML provides a XML representation of the Event as a string
func (e *Event) ToXML() (string, error) {
	eventXml, err := xml.Marshal(e)
	if err != nil {
		return "", err
	}

	// The Tags field is being ignore from XML Marshaling since maps are not supported.
	// We have to provide our own marshaling of the Tags field if it is non-empty
	if len(e.Tags) > 0 {
		tagsXmlElements := []string{"<Tags>"}
		for key, value := range e.Tags {
			tag := fmt.Sprintf("<%s>%s</%s>", key, value, key)
			tagsXmlElements = append(tagsXmlElements, tag)
		}
		tagsXmlElements = append(tagsXmlElements, "</Tags>")
		tagsXml := strings.Join(tagsXmlElements, "")
		eventXml = []byte(strings.Replace(string(eventXml), "</Event>", tagsXml+"</Event>", 1))
	}

	return string(eventXml), nil
}
