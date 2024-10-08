/*
 * Copyright 2022 The KodeRover Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/koderover/zadig/v2/pkg/microservice/aslan/core/common/service/lark"
	"github.com/koderover/zadig/v2/pkg/setting"
	internalhandler "github.com/koderover/zadig/v2/pkg/shared/handler"
)

func GetLarkDepartment(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	approvalID, departmentID := c.Param("id"), c.Param("department_id")
	userIDType := c.Query("user_id_type")
	if userIDType == "" {
		userIDType = setting.LarkUserOpenID
	}
	if departmentID == "root" {
		ctx.Resp, ctx.Err = lark.GetLarkAppContactRange(approvalID, userIDType)
	} else {
		ctx.Resp, ctx.Err = lark.GetLarkDepartment(approvalID, departmentID, userIDType)
	}
}

func GetLarkUserID(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()
	userIDType := c.Query("user_id_type")
	if userIDType == "" {
		userIDType = setting.LarkUserOpenID
	}
	id, err := lark.GetLarkUserID(c.Param("id"), c.Query("type"), c.Query("value"), userIDType)
	if err != nil {
		ctx.Err = err
		return
	}
	ctx.Resp = map[string]string{"id": id}
}

func LarkEventHandler(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	body, err := c.GetRawData()
	if err != nil {
		ctx.Err = err
		return
	}
	ctx.Resp, ctx.Err = lark.EventHandler(
		c.Param("id"),
		c.GetHeader("X-Lark-Signature"),
		c.GetHeader("X-Lark-Request-Timestamp"),
		c.GetHeader("X-Lark-Request-Nonce"), string(body))
}
