// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019-2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package autoevent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/OneOfOne/xxhash"
	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"
	"github.com/google/uuid"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/responses"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/errors"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/models"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/command"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/container"
	"github.com/tuya/tuya-edge-driver-sdk-go/logger"
)

type Executor struct {
	deviceName   string
	autoEvent    models.AutoEvent
	lastReadings map[string]interface{}
	duration     time.Duration
	stop         bool
	rwMutex      *sync.RWMutex
}

// Run triggers this Executor executes the handler for the resource periodically
func (e *Executor) Run(ctx context.Context, wg *sync.WaitGroup, dic *di.Container) {
	wg.Add(1)
	defer wg.Done()

	lc := bootstrapContainer.LoggingClientFrom(dic.Get)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(e.duration):
			if e.stop {
				return
			}
			ds := container.DeviceServiceFrom(dic.Get)
			if ds.AdminState == models.Locked {
				lc.Info("AutoEvent - stopped for locked device service")
				return
			}

			lc.Debug(fmt.Sprintf("AutoEvent - executing %v", e.autoEvent))
			correlationID := uuid.NewString()
			er, err := readResource(e, correlationID, dic)
			if err != nil {
				lc.Error(fmt.Sprintf("AutoEvent - error occurs when reading device %s, resource %s, error: %v",
					e.deviceName, e.autoEvent.Resource, err))
				continue
			}

			if len(er.Event.Readings) > 0 {
				// TODO 等compareReadings
				/*
					if e.autoEvent.OnChange {
						if compareReadings(e, er.Event.Readings, lc) {
							lc.Debug(fmt.Sprintf("AutoEvent - readings are the same as previous one %v", e.lastReadings))
							continue
						}
					}
				*/
				/*
					if evt.HasBinaryValue() {
						lc.Debug("AutoEvent - pushing CBOR event")
					} else {
						lc.Debug(fmt.Sprintf("AutoEvent - pushing event %s", evt.String()))
					}
						event := &dsModels.Event{Event: evt.Event}
						// Attach origin timestamp for events if none yet specified
						if event.Origin == 0 {
							event.Origin = common.GetUniqueOrigin()
						}
				*/

				// After the auto event executes a read command, it will create a goroutine to send out events.
				// When the concurrent auto event amount becomes large, core-data might be hard to handle so many HTTP requests at the same time.
				// The device service will get some network errors like EOF or Connection reset by peer.
				// By adding a buffer here, the user can use the Service.AsyncBufferSize configuration to control the goroutine for sending events.
				go func(correlationID string) {
					m.autoeventBuffer <- true
					common.SendEvent(er.Event, lc, container.CoredataEventClientFrom(dic.Get))
					<-m.autoeventBuffer
				}(correlationID)
			} else {
				lc.Debug(fmt.Sprintf("AutoEvent - no event generated when reading resource %s", e.autoEvent.Resource))
			}
		}
	}
}

func readResource(e *Executor, correlationID string, dic *di.Container) (res responses.EventResponse, err errors.EdgeX) {
	vars := make(map[string]string, 2)
	vars[common.NameVar] = e.deviceName
	vars[common.CommandVar] = e.autoEvent.Resource
	return command.CommandHandler(true, false, correlationID, vars, "", dic)
}

func compareReadings(e *Executor, readings []dtos.BaseReading, lc logger.LoggingClient) bool {
	identical := true
	e.rwMutex.RLock()
	defer e.rwMutex.RUnlock()
	for _, r := range readings {
		// TODO
		switch e.lastReadings[r.ResourceName].(type) {
		case uint64:
			checksum := xxhash.Checksum64(r.BinaryValue)
			if e.lastReadings[r.ResourceName] != checksum {
				e.lastReadings[r.ResourceName] = checksum
				identical = false
			}
		case string:
			v, ok := e.lastReadings[r.ResourceName]
			if !ok || v != r.Value {
				e.lastReadings[r.ResourceName] = r.Value
				identical = false
			}
		case nil:
			if r.ValueType == contracts.ValueTypeBinary && len(r.BinaryValue) > 0 {
				e.lastReadings[r.ResourceName] = xxhash.Checksum64(r.BinaryValue)
			} else {
				e.lastReadings[r.ResourceName] = r.Value
			}
			identical = false
		default:
			lc.Error("Error: unsupported reading type (%T) in autoevent - %v\n", e.lastReadings[r.ResourceName], e.autoEvent)
			identical = false
		}
	}
	return identical
}

// Stop marks this Executor stopped
func (e *Executor) Stop() {
	e.stop = true
}

// NewExecutor creates an Executor for an AutoEvent
func NewExecutor(deviceName string, ae models.AutoEvent) (*Executor, error) {
	// check Frequency
	duration, err := time.ParseDuration(ae.Frequency)
	if err != nil {
		return nil, err
	}

	return &Executor{
		deviceName:   deviceName,
		autoEvent:    ae,
		lastReadings: make(map[string]interface{}),
		duration:     duration,
		stop:         false,
		rwMutex:      &sync.RWMutex{}}, nil
}
