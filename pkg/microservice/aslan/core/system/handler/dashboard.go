/*
Copyright 2022 The KodeRover Authors.

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
	"github.com/koderover/zadig/v2/pkg/microservice/aslan/core/system/service"
	internalhandler "github.com/koderover/zadig/v2/pkg/shared/handler"
	e "github.com/koderover/zadig/v2/pkg/tool/errors"
)

func CreateOrUpdateDashboardConfiguration(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	args := new(service.DashBoardConfig)
	if err := c.BindJSON(args); err != nil {
		ctx.RespErr = e.ErrInvalidParam.AddDesc("invalid dashboard config")
		return
	}

	ctx.RespErr = service.CreateOrUpdateDashboardConfiguration(ctx.UserName, ctx.UserID, args, ctx.Logger)
}

func GetDashboardConfiguration(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	ctx.Resp, ctx.RespErr = service.GetDashboardConfiguration(ctx.UserName, ctx.UserID, ctx.Logger)
}

func GetRunningWorkflow(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	ctx.Resp, ctx.RespErr = service.GetRunningWorkflow(ctx.Logger)
}

func GetMyWorkflow(c *gin.Context) {
	ctx, err := internalhandler.NewContextWithAuthorization(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	if err != nil {
		ctx.RespErr = fmt.Errorf("authorization Info Generation failed: err %s", err)
		ctx.UnAuthorized = true
		return
	}

	ctx.Resp, ctx.RespErr = service.GetMyWorkflow(c.Request.Header, ctx.UserName, ctx.UserID, ctx.Resources.IsSystemAdmin, c.Query("card_id"), ctx.Logger)
}

func GetMyEnvironment(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	projectName := c.Query("projectName")
	envName := c.Param("name")
	productionStr := c.Query("production")

	production := false
	if productionStr == "true" {
		production = true
	}

	ctx.Resp, ctx.RespErr = service.GetMyEnvironment(projectName, envName, production, ctx.UserName, ctx.UserID, ctx.Logger)
}
