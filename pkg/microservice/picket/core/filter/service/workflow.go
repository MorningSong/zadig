/*
Copyright 2021 The KodeRover Authors.

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

package service

import (
	"net/http"
	"net/url"

	"github.com/koderover/zadig/v2/pkg/shared/client/user"
	"github.com/koderover/zadig/v2/pkg/types"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/koderover/zadig/v2/pkg/microservice/picket/client/aslan"
	"github.com/koderover/zadig/v2/pkg/microservice/picket/client/opa"
	"github.com/koderover/zadig/v2/pkg/microservice/picket/config"
)

type rule struct {
	method   string
	endpoint string
}

func ListWorkflows(header http.Header, qs url.Values, uid string, logger *zap.SugaredLogger) ([]byte, error) {
	projects, _, err := user.New().ListAuthorizedProjectsByResourceAndVerb(uid, types.ResourceTypeWorkflow, types.WorkflowActionView)
	if err != nil {
		logger.Errorf("Failed to get allowed project names, err: %s", err)
		return nil, err
	}
	if len(projects) == 0 {
		return nil, nil
	}
	for _, name := range projects {
		qs.Add("projects", name)
	}
	return aslan.New().ListWorkflows(header, qs)
}

func ListAllWorkflows(header http.Header, qs url.Values, uid string, logger *zap.SugaredLogger) ([]byte, error) {
	projects, _, err := user.New().ListAuthorizedProjectsByResourceAndVerb(uid, types.ResourceTypeWorkflow, types.WorkflowActionView)
	if err != nil {
		logger.Errorf("Failed to get allowed project names, err: %s", err)
		return nil, err
	}
	if len(projects) == 0 {
		return nil, nil
	}
	for _, name := range projects {
		qs.Add("projects", name)
	}
	return aslan.New().ListAllWorkflows(header, qs)
}

func ListTestWorkflows(testName, uid string, header http.Header, qs url.Values, logger *zap.SugaredLogger) ([]byte, error) {
	authorizedProjects, _, err := user.New().ListAuthorizedProjects(uid)
	if err != nil {
		logger.Errorf("Failed to get allowed project names, err: %s", err)
		return nil, err
	}
	for _, project := range authorizedProjects {
		qs.Add("projects", project)
	}
	return aslan.New().ListTestWorkflows(testName, header, qs)
}

// getAllowedProjects
// rulesLogicalOperator@ OR:satisfy any one of rules / AND:satisfy all rules
func getAllowedProjects(headers http.Header, rules []*rule, rulesLogicalOperator config.RulesLogicalOperator, logger *zap.SugaredLogger) (projects []string, err error) {
	var res [][]string
	for _, v := range rules {
		allowedProjects := &allowedProjectsData{}
		opaClient := opa.NewDefault()
		err := opaClient.Evaluate("rbac.user_allowed_projects", allowedProjects, func() (*opa.Input, error) {
			return generateOPAInput(headers, v.method, v.endpoint), nil
		})
		if err != nil {
			logger.Errorf("opa evaluation failed, err: %s", err)
			return nil, err
		}
		res = append(res, allowedProjects.Result)
	}
	if rulesLogicalOperator == config.OR {
		return union(res), nil
	}
	return intersect(res), nil
}

func union(s [][]string) []string {
	if len(s) == 0 {
		return nil
	}
	tmp := sets.NewString(s[0]...)
	for _, v := range s[1:] {
		t := sets.NewString(v...)
		tmp = t.Union(tmp)
	}
	return tmp.List()
}

func intersect(s [][]string) []string {
	if len(s) == 0 {
		return nil
	}
	tmp := sets.NewString(s[0]...)
	for _, v := range s[1:] {
		t := sets.NewString(v...)
		tmp = t.Intersection(tmp)
	}
	return tmp.List()
}
