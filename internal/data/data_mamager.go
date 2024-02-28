package data

import (
	"context"
	"fmt"
	"order_service/pkg/logger"
	"order_service/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type DataManager interface {
	ListPatients(context.Context) ([]models.Patient, error)
	ListOrderByPatientId(context.Context, int) ([]models.Order, error)
	UpdateOrder(context.Context, *models.Order) error
	CreateOrder(context.Context, *models.Order) error
	Close(context.Context)
}

type dataMgr struct {
	gormClient  *gorm.DB
	mongoClient *mongo.Collection
}

func NewDataManager(gormClient *gorm.DB, mongoClient *mongo.Collection) DataManager {
	return &dataMgr{
		gormClient:  gormClient,
		mongoClient: mongoClient,
	}
}

func (d *dataMgr) ListPatients(ctx context.Context) ([]models.Patient, error) {
	var patients []models.Patient

	if err := d.gormClient.Find(&patients).Error; err != nil {
		return nil, fmt.Errorf("ListPatients: %s", err.Error())
	}
	return patients, nil
}
func (mgr *dataMgr) ListOrderByPatientId(ctx context.Context, patientId int) ([]models.Order, error) {

	filter := bson.D{{"patient_id", patientId}}
	cursor, err := mgr.mongoClient.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("ListOrderByPatientId: %s", err.Error())
	}

	var orders []models.Order
	if err = cursor.All(context.TODO(), &orders); err != nil {
		return nil, fmt.Errorf("ListOrderByPatientId: %s", err.Error())
	}

	return orders, nil
}
func (mgr *dataMgr) UpdateOrder(ctx context.Context, order *models.Order) error {

	filter := bson.D{{"doctor_id", order.DoctorID}}

	// 將醫囑結構轉換為BSON文檔
	update := bson.D{{"$set", bson.D{
		{"patient_id", order.PatientID},
		{"doctor_id", order.DoctorID},
		{"created_at", order.CreatedAt},
		{"content", order.Content},
		{"status", order.Status},
	}}}
	result, err := mgr.mongoClient.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("UpdateOrder: %s", err.Error())
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("UpdateOrder: No order updated.")
	}

	return nil
}

func (mgr *dataMgr) CreateOrder(ctx context.Context, order *models.Order) error {
	if _, err := mgr.mongoClient.InsertOne(ctx, order); err != nil {
		return fmt.Errorf("CreateOrder: %s", err.Error())
	}

	return nil

}

func (mgr *dataMgr) Close(ctx context.Context) {
	db, err := mgr.gormClient.DB()
	if err != nil {
		logger.Errorf("Close DB : %s", err.Error())
	}
	if err = db.Close(); err != nil {
		logger.Errorf("Close DB : %s", err.Error())
	}

	if err := mgr.mongoClient.Database().Client().Disconnect(ctx); err != nil {
		logger.Errorf("Close DB : %s", err.Error())
	}

}
