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
	"github.com/gin-gonic/gin"
)

type Router struct{}

func (*Router) Inject(router *gin.RouterGroup) {
	// 查看html测试报告不做鉴权
	testReport := router.Group("report")
	{
		testReport.GET("/html/testing/:projectName/:testingName/:taskID/*path", GetTestTaskHtmlReportInfo)
		testReport.GET("/html/workflowv4/:projectName/:workflowName/:jobName/:taskID/*path", GetWorkflowV4HTMLTestReport)
	}

	// sse apis
	sse := router.Group("sse")
	{
		sse.GET("/:testName/tasks/:taskID", GetTestingTaskSSE)
	}

	// ---------------------------------------------------------------------------------------
	// 系统测试接口
	// ---------------------------------------------------------------------------------------
	itReport := router.Group("itreport")
	{
		itReport.GET("/pipelines/:pipelineName/id/:id/names/:testName", GetLocalTestSuite)
		itReport.GET("/workflowv4/:workflowName/id/:id/job/:jobName", GetWorkflowV4LocalTestSuite)
		itReport.GET("/workflow/:pipelineName/id/:id/names/:testName/service/:serviceName", GetWorkflowLocalTestSuite)
		itReport.GET("/latest/service/:testName", GetTestLocalTestSuite)
	}

	// ---------------------------------------------------------------------------------------
	// 测试管理接口
	// ---------------------------------------------------------------------------------------
	tester := router.Group("test")
	{
		tester.POST("", GetTestProductName, CreateTestModule)
		tester.PUT("", GetTestProductName, UpdateTestModule)
		tester.GET("", ListTestModules)
		tester.GET("/:name", GetTestModule)
		tester.DELETE("/:name", DeleteTestModule)
	}

	// ---------------------------------------------------------------------------------------
	// Code scan APIs
	// ---------------------------------------------------------------------------------------
	scanner := router.Group("scanning")
	{
		// code scan config apis
		scanner.POST("", GetScanningProductName, CreateScanningModule)
		scanner.PUT("/:id", GetScanningProductName, UpdateScanningModule)
		scanner.GET("", ListScanningModule)
		scanner.GET("/:id", GetScanningModule)
		scanner.DELETE("/:id", DeleteScanningModule)

		// code scan tasks apis
		scanner.POST("/:id/task", FindScanningProjectNameFromID, CreateScanningTask)
		scanner.GET("/:id/task", FindScanningProjectNameFromID, ListScanningTask)
		scanner.GET("/:id/task/:scan_id", FindScanningProjectNameFromID, GetScanningTask)
		scanner.DELETE("/:id/task/:scan_id", FindScanningProjectNameFromID, CancelScanningTask)
		scanner.GET("/:id/task/:scan_id/sse", FindScanningProjectNameFromID, GetScanningTaskSSE)

		scanner.GET("/:id/task/:scan_id/artifact_info", GetScanningArtifactInfo)
		scanner.GET("/artifact", GetScanningTaskArtifact)
	}

	//testStat := router.Group("teststat")
	//{
	//	// 供aslanx的enterprise模块的数据统计调用
	//	testStat.GET("", ListTestStat)
	//}

	testDetail := router.Group("testdetail")
	{
		testDetail.GET("", ListDetailTestModules)
	}

	// ---------------------------------------------------------------------------------------
	// test 任务接口
	// ---------------------------------------------------------------------------------------
	testTask := router.Group("testtask")
	{
		testTask.POST("", CreateTestTask)
		testTask.GET("", ListTestTask)
		testTask.DELETE("", CancelTestTaskV3)
		testTask.GET("/detail", GetTestTaskInfo)
		testTask.GET("/report", GetTestTaskJUnitReportInfo)
		testTask.POST("/restart", RestartTestTaskV2)
		testTask.GET("/artifact", GetTestingTaskArtifact)
		// TODO:  below is the deprecated apis, remove after 2.2.0
		//testTask.POST("/productName/:productName/id/:id/pipelines/:name/restart", RestartTestTask)
		//testTask.DELETE("/productName/:productName/id/:id/pipelines/:name", CancelTestTaskV2)
	}

	// ---------------------------------------------------------------------------------------
	// Pipeline workspace 管理接口
	// ---------------------------------------------------------------------------------------
	workspace := router.Group("workspace")
	{
		workspace.GET("/workflow/:pipelineName/taskId/:taskId", GetTestArtifactInfo)
		workspace.GET("/testing/:testName/taskId/:taskId", GetTestArtifactInfoV2)

		workspace.GET("/workflowv4/:workflowName/taskId/:taskId/job/:jobName", GetWorkflowV4ArtifactInfo)
	}
}

type QualityCenterRouter struct{}

func (*QualityCenterRouter) Inject(router *gin.RouterGroup) {
	// testing apis
	test := router.Group("tests")
	{
		test.GET("", ListTestingWithStat)
	}
}

type QualityRouter struct{}

func (*QualityRouter) Inject(router *gin.RouterGroup) {
	scan := router.Group("codescan")
	{
		scan.POST("", OpenAPICreateScanningModule)
		scan.POST("/:scanName/task", OpenAPICreateScanningTask)
		scan.GET("/:scanName/task/:taskID", OpenAPIGetScanningTaskDetail)
	}

	test := router.Group("testing")
	{
		test.POST("/task", OpenAPICreateTestTask)
		test.GET("/:testName/task/:taskID", OpenAPIGetTestTaskResult)
	}
}
