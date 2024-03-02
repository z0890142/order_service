package controller

import (
	"context"
	"fmt"
	"net/http"
	"order_service/c"
	"order_service/internal/data"
	"order_service/internal/serivce/handler"
	"order_service/pkg/logger"
	"order_service/pkg/models"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/genproto/googleapis/rpc/code"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	dataMgr       data.DataManager
	tokenHandler  handler.TokenHandler
	shuntDownOnce sync.Once
}

func NewController(dataMgr data.DataManager) *Controller {
	tokenHandler := handler.NewTokenHandler(dataMgr)
	return &Controller{
		dataMgr:       dataMgr,
		shuntDownOnce: sync.Once{},
		tokenHandler:  tokenHandler,
	}
}

// @Summary list patinets
// @router /order-service/login [post]
// @param params body models.Doctor true "doctor"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} models.HttpError
func (ctrl *Controller) Login(ginc *gin.Context) {
	doctor := models.Doctor{}
	if err := ginc.BindJSON(&doctor); err != nil {
		ctrl.handleError(ginc, err, http.StatusBadRequest, code.Code_INTERNAL)
		return
	}

	token, err := ctrl.tokenHandler.GetAccessToken(ginc, doctor)
	if err != nil {
		ctrl.handleError(ginc, err, http.StatusBadRequest, code.Code_INTERNAL)
		return
	}

	ginc.JSON(http.StatusOK, token)
}

func (ctrl *Controller) RefreshToken(ginc *gin.Context) {

	authHeader := ginc.GetHeader("Authorization")
	refreshToken := strings.Split(authHeader, " ")[1]

	resp, err := ctrl.tokenHandler.RefreshToken(ginc, refreshToken)
	if err != nil {
		ctrl.handleError(ginc, err, http.StatusBadRequest, code.Code_INTERNAL)
		return
	}
	ginc.JSON(http.StatusOK, resp)
}

// @Summary list patinets
// @router /order-service/api/v1/patients [get]
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @param Authorization header string true "Authorization"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.HttpError
func (ctrl *Controller) ListPatinets(ginc *gin.Context) {
	limitStr := ginc.Query("limit")
	offsetStr := ginc.Query("offset")

	defaultLimit := 20
	defaultOffset := 0

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = defaultLimit
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = defaultOffset
	}

	patients, err := ctrl.dataMgr.ListPatients(ginc, limit, offset)
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

// @Summary list patinet's order
// @router /order-service/api/v1/patients/{patientId}/orders [get]
// @Param patientId path int true "patinet ID"
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.HttpError
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
		}).Error("ListOrders")
		ctrl.handleError(ginc, err, http.StatusBadRequest, code.Code_INTERNAL)
	}

	limitStr := ginc.Query("limit")
	offsetStr := ginc.Query("offset")

	defaultLimit := 20
	defaultOffset := 0

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = defaultLimit
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = defaultOffset
	}

	orders, err := ctrl.dataMgr.ListOrderByPatientId(ginc, patientID, limit, offset)
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

// @Summary update patinet's order
// @router /order-service/api/v1/patients/{patientId}/orders/{orderId} [put]
// @Param patientId path int true "patinet ID"
// @Param orderId path string true "order ID"
// @param params body models.Order true "order"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.HttpError
func (ctrl *Controller) UpdateOrder(ginc *gin.Context) {
	doctorInter, _ := ginc.Get("doctor")
	doctor := doctorInter.(models.Doctor)

	orderIDStr := ginc.Param("orderId")
	if orderIDStr == "" {
		logger.GetLoggerWithKeys(map[string]interface{}{
			"error": "INVALID_ARGUMENT",
		}).Error("UpdateOrder")
		ctrl.handleError(ginc, fmt.Errorf("INVALID_ARGUMENT"), http.StatusBadRequest, code.Code_INVALID_ARGUMENT)
		return
	}

	patientIDStr := ginc.Param("patientId")
	if patientIDStr == "" {
		logger.GetLoggerWithKeys(map[string]interface{}{
			"error": "INVALID_ARGUMENT",
		}).Error("ListPatinets")
		ctrl.handleError(ginc, fmt.Errorf("INVALID_ARGUMENT"), http.StatusBadRequest, code.Code_INVALID_ARGUMENT)
		return
	}

	patientID, err := strconv.Atoi(patientIDStr)
	if err != nil {
		logger.GetLoggerWithKeys(map[string]interface{}{
			"error": err,
		}).Error("ListOrders")
		ctrl.handleError(ginc, err, http.StatusBadRequest, code.Code_INTERNAL)

		return
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
		return
	}

	order.ID = orderId
	order.PatientID = patientID
	order.UpdatedAt = time.Now()
	order.DoctorName = doctor.Username

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

// @Summary create patinet's order
// @router /order-service/api/v1/patients/{patientId}/orders [post]
// @Param patientId path int true "patinet ID"
// @param params body models.Order true "order"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.HttpError
func (ctrl *Controller) CreateOrder(ginc *gin.Context) {
	doctorInter, _ := ginc.Get("doctor")
	doctor := doctorInter.(models.Doctor)

	patientIDStr := ginc.Param("patientId")
	if patientIDStr == "" {
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

	patientID, err := strconv.Atoi(patientIDStr)
	if err != nil {
		ctrl.handleError(ginc, err, http.StatusBadRequest, code.Code_INTERNAL)
		return
	}

	order.PatientID = patientID
	order.CreatedAt = time.Now()
	order.DoctorName = doctor.Username

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
