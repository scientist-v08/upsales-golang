package model

import "gorm.io/gorm"

type Movies struct {
	gorm.Model
	Title         	string
	Type          	string
	Director      	string
	Budget        	string // Assuming string for simplicity; change to int/float if needed
	Location      	string
	Duration      	string // Assuming string; change to int if needed
	YearOfRelease 	string
	IsAdminApproved bool
	Image         	[]byte // Store image as byte slice
	MimeType		string
}