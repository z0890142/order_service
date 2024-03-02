package models

import "github.com/dgrijalva/jwt-go"

type CustomerClaims struct {
	DoctorName string `json:"doctor_name"`
	DoctorId   uint   `json:"doctor_id"`
	jwt.StandardClaims
}
