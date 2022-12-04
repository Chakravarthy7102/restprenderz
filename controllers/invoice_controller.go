package controllers

import (
	"golang-resto-management/database"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type InvoiceViewFormat struct {
	Invoice_ID       string
	Payment_Method   string
	Order_Id         string
	Payment_status   *string
	Payment_Due      interface{}
	Table_number     interface{}
	Payment_Due_date time.Time
	Order_details    interface{}
}

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
