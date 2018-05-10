//
// Copyright (c) 2018 Juniper Networks, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller_rest

import (
	"github.com/Juniper/contrail-go-api"
	"github.com/Juniper/contrail-windows-docker-driver/adapters/secondary/controller_rest/api"
	"github.com/Juniper/contrail-windows-docker-driver/adapters/secondary/controller_rest/auth"
)

func NewControllerWithKeystoneAdapter(keys auth.KeystoneParams, apiClient *contrail.Client) (*ControllerAdapter, error) {
	auth, err := auth.NewKeystoneAuth(keys)
	if err != nil {
		return nil, err
	}

	err = auth.Authenticate()
	if err != nil {
		return nil, err
	}

	apiClient.SetAuthenticator(auth)

	impl := NewControllerAdapterImpl(apiClient)
	return newControllerAdapter(impl), nil
}

func NewControllerInsecureAdapter(apiClient contrail.ApiClient) (*ControllerAdapter, error) {
	impl := NewControllerAdapterImpl(apiClient)

	// When keystone is not present only default-project is created in controller.
	// WebUI doesn't treat like regular resource and doesn't permit some operations
	// on it, so a regular project needs to be created.
	_, err := impl.GetOrCreateProject(DomainName, AdminProject)
	if err != nil {
		return nil, err
	}
	return newControllerAdapter(impl), nil
}

func NewFakeControllerAdapter() *ControllerAdapter {
	fakeApiClient := api.NewFakeApiClient()
	impl := NewControllerAdapterImpl(fakeApiClient)
	return newControllerAdapter(impl)
}
