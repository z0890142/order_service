package app

import (
	"fmt"
	"order_service/internal/data"
	"order_service/internal/serivce/controller"
	"order_service/internal/serivce/middleware"
	"order_service/pkg/database"

	"github.com/gin-contrib/cors"
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

	mid := middleware.NewMiddleware(dataMgr)
	ctrl := controller.NewController(dataMgr)

	r.POST("/login", ctrl.Login)
	r.POST("/refresh-token", ctrl.RefreshToken)

	v1Group := r.Group("order-service/api/v1")
	v1Group.Use(mid.AuthMiddleware)
	v1Group.GET("/patients", ctrl.ListPatinets)
	v1Group.GET("/patients/:patientId/orders", ctrl.ListOrders)
	v1Group.PUT("/patients/:patientId/orders/:orderId", ctrl.UpdateOrder)
	v1Group.POST("/patients/:patientId/orders", ctrl.CreateOrder)

	return nil
}

func InitGinApplicationHook(app *Application) error {
	if app.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.EnableJsonDecoderUseNumber()

	r := gin.New()
	config := cors.DefaultConfig()
	config.AddAllowHeaders("Authorization")
	config.AllowCredentials = true
	config.AllowAllOrigins = false
	config.AllowOriginFunc = func(origin string) bool {
		return true
	}

	r.Use(cors.New(config))
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
