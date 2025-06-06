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

package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/koderover/zadig/v2/pkg/microservice/systemconfig/core/features/service"
	internalhandler "github.com/koderover/zadig/v2/pkg/shared/handler"
)

type feature struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

func GetFeature(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	name := c.Param("name")
	enabled, err := service.FeatureEnabled(name, ctx.Logger)

	if err != nil {
		ctx.RespErr = err
		return
	}

	ctx.Resp = &feature{
		Name:    c.Param("name"),
		Enabled: enabled,
	}
}

func UpdateOrCreateFeature(c *gin.Context) {
	ctx, err := internalhandler.NewContextWithAuthorization(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	if err != nil {
		ctx.RespErr = fmt.Errorf("authorization Info Generation failed: err %s", err)
		ctx.UnAuthorized = true
		return
	}

	// authorization checks
	if !ctx.Resources.IsSystemAdmin {
		ctx.UnAuthorized = true
		return
	}

	req := new(service.FeatureReq)
	if err := c.ShouldBindJSON(req); err != nil {
		ctx.RespErr = err
		return
	}
	ctx.RespErr = service.UpdateOrCreateFeature(req, ctx.Logger)
}
