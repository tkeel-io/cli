/*
Copyright 2021 The tKeel Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kubernetes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/yaml"
)

const pluginCore = "core"

type CoreStruct struct {
	APIVersion string                 `yaml:"apiVersion"` //nolint
	Type       string                 `yaml:"type"`
	ID         string                 `yaml:"id"`
	Source     string                 `yaml:"source"`
	Owner      string                 `yaml:"owner"`
	Properties map[string]interface{} `yaml:"properties"`
}

func file2CoreStruct(filename string) (*CoreStruct, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read create file")
	}

	coreStruct := CoreStruct{}
	err = yaml.Unmarshal(bytes, &coreStruct)

	if err != nil {
		return nil, errors.Wrap(err, "unable to parse create file")
	}

	return &coreStruct, nil
}

func CoreApply(filenames []string) error {
	for _, filename := range filenames {
		coreStruct, err := file2CoreStruct(filename)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		method := fmt.Sprintf("v1/entities/%s?source=%s&owner=%s&type=%s", coreStruct.ID, coreStruct.Source, coreStruct.Owner, coreStruct.Type)
		data, _ := json.Marshal(coreStruct.Properties)
		if resp, err := InvokeByPortForward(pluginCore, method, data, http.MethodPut); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(resp)
		}
	}
	return nil
}

func CoreCreate(filenames []string) error {
	for _, filename := range filenames {
		coreStruct, err := file2CoreStruct(filename)
		if err != nil {
			fmt.Println(err)
			continue
		}

		method := fmt.Sprintf("v1/entities?source=%s&owner=%s&type=%s", coreStruct.Source, coreStruct.Owner, coreStruct.Type)
		if coreStruct.ID != "" {
			method += fmt.Sprintf("&id=%s", coreStruct.ID)
		}

		data, _ := json.Marshal(coreStruct.Properties)

		if resp, err := InvokeByPortForward(pluginCore, method, data, http.MethodPost); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(resp)
		}
	}
	return nil
}

type SearchCondition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

type SearchRequest struct {
	Query     string             `json:"query"`
	Condition []*SearchCondition `json:"condition"`
}

var operatorMap = map[string]string{"=": "$eq"}

func selector2SearchConditions(selector string) []*SearchCondition {
	resp := make([]*SearchCondition, 0)
	items := strings.Split(selector, ",")
	for _, item := range items {
		for k, v := range operatorMap {
			iiterms := strings.Split(item, k)
			if len(iiterms) == 2 {
				searchCondition := &SearchCondition{
					Field:    iiterms[0],
					Operator: v,
					Value:    iiterms[1],
				}
				resp = append(resp, searchCondition)
			}
		}
	}
	return resp
}

func CoreList(search, selector string) error {
	searchRequest := &SearchRequest{}
	searchRequest.Query = search
	searchRequest.Condition = selector2SearchConditions(selector)

	method := "v1/entities/search"
	data, _ := json.Marshal(searchRequest)
	if resp, err := InvokeByPortForward(pluginCore, method, data, http.MethodPost); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
	return nil
}

func CoreWatch(entityID string) error {
	method := "v1/ws"
	data, _ := json.Marshal(map[string]string{"id": entityID})
	if resp, err := WebsocketByPortForward("core-broker", method, data); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
	return nil
}

func CoreGet(entityID string) error {
	method := fmt.Sprintf("v1/entities/%s", entityID)
	if resp, err := InvokeByPortForward(pluginCore, method, nil, http.MethodGet); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
	return nil
}
