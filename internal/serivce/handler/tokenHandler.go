package handler

import (
	"context"
	"fmt"
	"order_service/config"
	"order_service/internal/data"
	"order_service/pkg/logger"
	"order_service/pkg/models"
	"time"

	"github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"
)

type TokenHandler interface {
	GetAccessToken(ctx context.Context, doctor models.Doctor) (*models.TokenResponse, error)
	RefreshToken(ctx context.Context, refreshToekn string) (*models.TokenResponse, error)
}

type tokenHandler struct {
	dataMgr                data.DataManager
	jwtKey                 []byte
	accessTokenExpireTime  time.Duration
	refreshTokenExpireTime time.Duration
}

func NewTokenHandler(dataMgr data.DataManager) TokenHandler {
	return &tokenHandler{
		dataMgr:                dataMgr,
		jwtKey:                 []byte(config.GetConfig().Jwt.Key),
		accessTokenExpireTime:  time.Duration(config.GetConfig().Jwt.AccessExpire) * time.Second,
		refreshTokenExpireTime: time.Duration(config.GetConfig().Jwt.RefreshExpire) * time.Second,
	}
}

func (t *tokenHandler) GetAccessToken(ctx context.Context, doctor models.Doctor) (*models.TokenResponse, error) {
	findDoctor := models.Doctor{}

	err := t.dataMgr.GetDoctor(ctx, map[string]interface{}{
		"username": doctor.Username,
	}, &findDoctor)
	if err != nil {
		return nil, fmt.Errorf("Login: %s", err.Error())
	}

	if !t.checkPasswordHash(findDoctor.Password, doctor.Password) {
		return nil, fmt.Errorf("Login: password incorrect")
	}

	accessToken, err := t.generateAccessToekn(findDoctor)
	if err != nil {
		return nil, fmt.Errorf("Login: %s", err.Error())
	}

	refreshToken, err := t.generateRefreshToekn(findDoctor)
	if err != nil {
		return nil, fmt.Errorf("Login: %s", err.Error())
	}
	return &models.TokenResponse{
		AccessToken:       accessToken,
		RefreshToken:      refreshToken,
		AccessExpireTime:  config.GetConfig().Jwt.AccessExpire,
		RefreshExpireTime: config.GetConfig().Jwt.RefreshExpire,
	}, nil
}

func (t *tokenHandler) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hashPassword: %s", err.Error())
	}
	return string(bytes), nil
}

func (t *tokenHandler) checkPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		logger.GetLoggerWithKeys(map[string]interface{}{
			"error": err,
		}).Error("checkPasswordHash fail")
	}
	return err == nil
}

func (t *tokenHandler) generateRefreshToekn(doctor models.Doctor) (string, error) {

	refreshTokenExpire := time.Now().Add(t.refreshTokenExpireTime)
	refreshTokenClaims := models.CustomerClaims{
		DoctorName: doctor.Username,
		DoctorId:   doctor.ID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: refreshTokenExpire.Unix(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(t.jwtKey)
	if err != nil {
		return "", fmt.Errorf("generateJwt: %s", err.Error)
	}

	return refreshTokenString, nil
}

func (t *tokenHandler) generateAccessToekn(doctor models.Doctor) (string, error) {
	accessTokenExpire := time.Now().Add(t.accessTokenExpireTime)
	accessTokenClaims := models.CustomerClaims{
		DoctorName: doctor.Username,
		DoctorId:   doctor.ID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: accessTokenExpire.Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(t.jwtKey)
	if err != nil {
		return "", fmt.Errorf("generateJwt: %s", err.Error)
	}

	return accessTokenString, nil
}

func (t *tokenHandler) RefreshToken(ctx context.Context, refreshToekn string) (*models.TokenResponse, error) {
	claims := &models.CustomerClaims{}
	_, err := jwt.ParseWithClaims(refreshToekn, claims, func(token *jwt.Token) (interface{}, error) {
		return t.jwtKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("RefreshToken: %s", err.Error())
	}
	if claims.Valid() != nil {
		return nil, fmt.Errorf("RefreshToken: %s", claims.Valid())

	}

	doctor := models.Doctor{}

	err = t.dataMgr.GetDoctor(ctx, map[string]interface{}{
		"username": claims.DoctorName,
		"id":       claims.DoctorId,
	}, &doctor)
	if err != nil {
		return nil, fmt.Errorf("RefreshToken: %s", err.Error())
	}

	accessTokenString, err := t.generateAccessToekn(doctor)
	if err != nil {
		return nil, fmt.Errorf("RefreshToken: %s", err.Error())
	}

	return &models.TokenResponse{
		AccessToken:      accessTokenString,
		RefreshToken:     refreshToekn,
		AccessExpireTime: config.GetConfig().Jwt.AccessExpire,
	}, nil
}
