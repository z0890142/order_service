package controller

import (
	"context"
	"fmt"
	"net/http"
	"order_service/c"
	"order_service/internal/data"
	"order_service/pkg/logger"
	"order_service/pkg/models"
	"strconv"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/genproto/googleapis/rpc/code"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	dataMgr       data.DataManager
	shuntDownOnce sync.Once
}

func NewController(dataMgr data.DataManager) *Controller {
	return &Controller{
		dataMgr:       dataMgr,
		shuntDownOnce: sync.Once{},
	}
}

func (ctrl *Controller) ListPatinets(ginc *gin.Context) {

	patients, err := ctrl.dataMgr.ListPatients(ginc)
	if err != nil {
		logger.GetLoggerWithKeys(map[string]interface{}{
			"error": err,
		}).Error("ListPatinets")
		ctrl.handleError(ginc, err, http.StatusBadRequest, code.Code_INTERNAL)
	}

	response := models.Response{
		Code:    code.Code_OK,
		Message: c.Success,
		Data:    patients,
	}
	ginc.JSON(http.StatusOK, response)
}
func (ctrl *Controller) ListOrders(ginc *gin.Context) {
	patientIDStr := ginc.Param("patientId")
	if patientIDStr == "" {
		logger.GetLoggerWithKeys(map[string]interface{}{
			"error": "INVALID_ARGUMENT",
		}).Error("ListPatinets")
		ctrl.handleError(ginc, fmt.Errorf("INVALID_ARGUMENT"), http.StatusBadRequest, code.Code_INVALID_ARGUMENT)
	}

	patientID, err := strconv.Atoi(patientIDStr)
	if err != nil {
		logger.GetLoggerWithKeys(map[string]interface{}{
			"error": err,
		}).Error("ListPatinets")
		ctrl.handleError(ginc, err, http.StatusBadRequest, code.Code_INTERNAL)
	}

	orders, err := ctrl.dataMgr.ListOrderByPatientId(ginc, patientID)
	if err != nil {
		logger.GetLoggerWithKeys(map[string]interface{}{
			"error": err,
		}).Error("ListPatinets")
		ctrl.handleError(ginc, err, http.StatusBadRequest, code.Code_INTERNAL)
	}

	response := models.Response{
		Code:    code.Code_OK,
		Message: c.Success,
		Data:    orders,
	}
	ginc.JSON(http.StatusOK, response)
}
func (ctrl *Controller) UpdateOrder(ginc *gin.Context) {

	orderIDStr := ginc.Param("oriderId")
	if orderIDStr == "" {
		logger.GetLoggerWithKeys(map[string]interface{}{
			"error": "INVALID_ARGUMENT",
		}).Error("UpdateOrder")
		ctrl.handleError(ginc, fmt.Errorf("INVALID_ARGUMENT"), http.StatusBadRequest, code.Code_INVALID_ARGUMENT)
	}

	order := models.Order{}
	if err := ginc.BindJSON(&order); err != nil {
		ctrl.handleError(ginc, err, http.StatusBadRequest, code.Code_INTERNAL)
		return
	}

	orderId, err := StringToOID(orderIDStr)
	if err != nil {
		logger.GetLoggerWithKeys(map[string]interface{}{
			"error": err,
		}).Error("UpdateOrder")
		ctrl.handleError(ginc, err, http.StatusBadRequest, code.Code_INTERNAL)
	}

	order.ID = orderId

	if err := ctrl.dataMgr.UpdateOrder(ginc, &order); err != nil {
		ctrl.handleError(ginc, err, http.StatusBadRequest, code.Code_INTERNAL)
		return
	}

	response := models.Response{
		Code:    code.Code_OK,
		Message: c.Success,
		Data: []models.Order{
			order,
		},
	}
	ginc.JSON(http.StatusOK, response)

}
func (ctrl *Controller) CreateOrder(ginc *gin.Context) {
	patientID := ginc.Param("patientId")
	if patientID == "" {
		logger.GetLoggerWithKeys(map[string]interface{}{
			"error": "INVALID_ARGUMENT",
		}).Error("CreateOrder")
		ctrl.handleError(ginc, fmt.Errorf("INVALID_ARGUMENT"), http.StatusBadRequest, code.Code_INVALID_ARGUMENT)
	}

	order := models.Order{}
	if err := ginc.BindJSON(&order); err != nil {
		ctrl.handleError(ginc, err, http.StatusBadRequest, code.Code_INTERNAL)
		return
	}
	order.PatientID = patientID

	if err := ctrl.dataMgr.CreateOrder(ginc, &order); err != nil {
		ctrl.handleError(ginc, err, http.StatusBadRequest, code.Code_INTERNAL)
		return
	}

	response := models.Response{
		Code:    code.Code_OK,
		Message: c.Success,
		Data: []models.Order{
			order,
		},
	}
	ginc.JSON(http.StatusOK, response)
}

func (ctrl *Controller) Shutdown() {
	ctrl.shuntDownOnce.Do(func() {
		ctx := context.TODO()
		ctrl.dataMgr.Close(ctx)
	})
}

func (ctrl *Controller) handleError(ginc *gin.Context, err error, httpCode int, errorCode code.Code) {
	ginc.JSON(httpCode, models.HttpError{
		Code:    errorCode,
		Message: err.Error(),
	})
}

func StringToOID(str string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(str)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return oid, nil
}
