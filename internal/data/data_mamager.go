package data

import (
	"context"
	"fmt"
	"order_service/pkg/logger"
	"order_service/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gorm.io/gorm"
)

type DataManager interface {
	ListPatients(context.Context, int, int) ([]models.Patient, error)
	ListOrderByPatientId(context.Context, int, int, int) ([]models.Order, error)
	UpdateOrder(context.Context, *models.Order) error
	CreateOrder(context.Context, *models.Order) error
	GetDoctor(context.Context, map[string]interface{}, *models.Doctor) error
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

func (d *dataMgr) ListPatients(ctx context.Context, limit, offset int) ([]models.Patient, error) {
	var patients []models.Patient

	if err := d.gormClient.
		Offset(offset).Limit(limit).
		Find(&patients).
		Error; err != nil {
		return nil, fmt.Errorf("ListPatients: %s", err.Error())
	}
	return patients, nil
}
func (mgr *dataMgr) ListOrderByPatientId(ctx context.Context, patientId, limit, offset int) ([]models.Order, error) {

	filter := bson.D{{"patient_id", patientId}}
	options := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))

	cursor, err := mgr.mongoClient.Find(ctx, filter, options)
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

	filter := bson.D{
		{"_id", order.ID},
		{"patient_id", order.PatientID},
	}

	// 將醫囑結構轉換為BSON文檔
	update := bson.D{{"$set", bson.D{
		{"doctor_id", order.DoctorID},
		{"content", order.Content},
		{"status", order.Status},
		{"updated_at", order.UpdatedAt},
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
	result, err := mgr.mongoClient.InsertOne(ctx, order)
	if err != nil {
		return fmt.Errorf("CreateOrder: %s", err.Error())
	}

	order.ID = result.InsertedID.(primitive.ObjectID)
	return nil

}

func (mgr *dataMgr) GetDoctor(ctx context.Context, condition map[string]interface{}, doctor *models.Doctor) error {
	if err := mgr.gormClient.
		Where(condition).
		First(doctor).Error; err != nil {
		return fmt.Errorf("GetDoctor: %v", err)
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
