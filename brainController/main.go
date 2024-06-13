package main

import (
	"cloud/brainController/common"
	"cloud/brainController/controller"
	"cloud/brainController/middleware"
	"cloud/brainController/service/informer"
	"cloud/brainController/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	Init()
	r := gin.Default()
	r.Use(middleware.Cors())
	initRouter(r)
	err := r.Run(":8080")
	if err != nil {
		return
	}
}

func Init() {
	utils.LoadConf()
	common.TransformConfig()
	utils.ClientSetInit()
	utils.MysqlDBInit()
	utils.OSSInit()
	utils.DockerClientInit()
	utils.InitWS()
	utils.InfluxInit()
	go informer.InitInformer()
}
func initRouter(r *gin.Engine) {
	utils.M.HandleMessage(controller.HandlerMessage)
	phmGroup := r.Group("/api")
	{
		brainControllerGroup := phmGroup.Group("/brainController")
		{
			//jupyter算法建模平台
			brainControllerGroup.POST("/createJupyterResource", controller.CreateJupyterResource)
			brainControllerGroup.POST("/deleteJupyterResource", controller.DeleteJupyterResource)

			//jupyter建模平台输出的phm资源
			brainControllerGroup.POST("/pushPHMResource", controller.PushPHMResource)
			brainControllerGroup.GET("/deletePHMResource", controller.DeletePHMResource)
			brainControllerGroup.GET("/getIaiImageList", middleware.Cookie(), controller.GetIaiImageList)
			brainControllerGroup.POST("/editIaiImageRecord", controller.EditIaiImageRecord)

			//云边协同节点资源
			brainControllerGroup.GET("/getNodeList", middleware.Cookie(), controller.GetNodeList)
			brainControllerGroup.GET("/getNodeInfo", middleware.Cookie(), controller.GetNodeInfo)
			brainControllerGroup.GET("/getNodeMetrics", middleware.Cookie(), controller.GetNodeListWithMetrics)
			brainControllerGroup.POST("/getIaiPodStatus", middleware.Cookie(), controller.GetIaiPodStatus)
			brainControllerGroup.POST("/getIaiPodName", middleware.Cookie(), controller.GetIaiPodName)
			brainControllerGroup.POST("/getIaiPodHostName", middleware.Cookie(), controller.GetIaiPodHostName)

			//云边协同phm任务
			brainControllerGroup.POST("/createTask", controller.CreateTask)
			brainControllerGroup.POST("/createTaskNew", middleware.Cookie(), controller.CreateTaskNew)
			brainControllerGroup.POST("/deleteTask", middleware.Cookie(), controller.DeleteTask)
			brainControllerGroup.GET("/getTaskList", middleware.Cookie(), controller.GetTaskList)
			brainControllerGroup.GET("/getTask", middleware.Cookie(), controller.GetTask)
			brainControllerGroup.GET("/getTasks", middleware.Cookie(), controller.GetAllTasks)
			brainControllerGroup.POST("/podMigration", middleware.Cookie(), controller.PodMigration)
			brainControllerGroup.POST("/podImagesUpdate", middleware.Cookie(), controller.PodImagesUpdate)
			brainControllerGroup.POST("/nodeAffinityMigration", middleware.Cookie(), controller.NodeAffinityMigration)

			//获取phm pod中的信息
			brainControllerGroup.POST("/getPhmPodFileList", controller.GetPhmFileList)
			brainControllerGroup.POST("/getPhmPodFile", controller.GetPhmFile)
			brainControllerGroup.POST("/getPhmPodExecResult", controller.GetPhmExecResult)

			//边缘小脑websocket
			brainControllerGroup.GET("/ws", controller.HandleWSReq)
		}
		platformControllerGroup := phmGroup.Group("/platform")
		{
			platformControllerGroup.GET("/getDataMapperConfig", controller.GetMapper)
			platformControllerGroup.GET("/getSchemaConfig", controller.GetBlocks)
			platformControllerGroup.GET("/searchImageWarehouse", middleware.Cookie(), controller.SearchImage)
			platformControllerGroup.GET("/getBaseCount", controller.GetBaseCount)
			platformControllerGroup.GET("/user", middleware.Cookie(), controller.GetUserDetail)
		}
	}
}
