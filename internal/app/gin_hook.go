package app

import (
	"fmt"
	"order_service/internal/data"
	"order_service/internal/serivce/controller"
	"order_service/pkg/database"

	"github.com/gin-gonic/gin"
)

var ctrl *controller.Controller

func initCtrl(app *Application, r *gin.Engine) error {

	gormCli, err := database.InitGormClient(app.GetDatabase())
	if err != nil {
		fmt.Errorf("initCtrl: %s", err.Error())
	}

	mongoCli, err := database.MongoConnect()
	if err != nil {
		fmt.Errorf("initCtrl: %s", err.Error())
	}

	dataMgr := data.NewDataManager(gormCli, mongoCli)
	ctrl := controller.NewController(dataMgr)

	r.GET("order-service/api/v1/patients", ctrl.ListPatinets)
	r.GET("order-service/api/v1/patients/:patientId/orders", ctrl.ListOrders)
	r.PUT("order-service/api/v1/patients/:patientId/orders/:orderId", ctrl.UpdateOrder)
	r.POST("order-service/api/v1/patients/:patientId/orders", ctrl.CreateOrder)

	return nil
}

func InitGinApplicationHook(app *Application) error {
	if app.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.EnableJsonDecoderUseNumber()
	r := gin.New()
	r.Use(gin.Recovery())

	initCtrl(app, r)
	addr := fmt.Sprintf("%s:%s", app.GetConfig().Service.Host, app.GetConfig().Service.Port)

	app.SetAddr(addr)
	app.SetSrv(r)

	return nil
}

func DestroyGinApplicationHook(app *Application) error {
	ctrl.Shutdown()
	return nil
}
