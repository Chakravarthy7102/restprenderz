package routes

import (
	"golang-resto-management/controllers"

	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(incommingRoutes *gin.Engine) {
	incommingRoutes.GET("/invoices", controllers.GetInvoices())
	incommingRoutes.GET("invoice/:invoice_id", controllers.GetInvoice())
	incommingRoutes.POST("/invoice", controllers.CreateInvoice())
	//there is no editing option for a invoice that's dumb!!
}
