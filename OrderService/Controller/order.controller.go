package controller

import (
	"fmt"
	"net/http"
	"time"

	"deneme.com/bng-go/Model"
	service "deneme.com/bng-go/Service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderController struct {
	OrderService service.OrderService
}

func New(orderservice service.OrderService) OrderController {
	return OrderController{OrderService: orderservice}
}

func (oc *OrderController) CreateOrder(ctx *gin.Context) {
	var order Model.Order
	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := oc.OrderService.CreateOrder(&order)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (oc *OrderController) GetOrder(ctx *gin.Context) {
	orderid, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid order ID"})
		return
	}
	order, err := oc.OrderService.GetOrder(&orderid)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success", "order": order})
}

func (oc *OrderController) GetOrders(ctx *gin.Context) {
	status, err, checkStatus := Model.ParseOrderStatus(ctx.Param("status"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	orders, err := oc.OrderService.GetOrders(&status, &checkStatus)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success", "orders": orders})
}

func (oc *OrderController) UpdateOrder(ctx *gin.Context) {
	var update_req struct {
		Quantity     int        `json:"quantity" bson:"quantity,omitempty"`
		Price        float64    `json:"price" bson:"price,omitempty"`
		DeliveryDate *time.Time `json:"delivery_date" bson:"delivery_date,omitempty"`
	}
	orderid, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid order ID"})
		return
	}
	if err := ctx.ShouldBindJSON(&update_req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	fmt.Println(update_req)
	err = oc.OrderService.UpdateOrder(&orderid, update_req)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (oc *OrderController) UpdateStatus(ctx *gin.Context) {
	var status_req struct {
		Id     uuid.UUID `json:"id" bson:"_id,omitempty"`
		Status string    `json:"status" bson:"status,omitempty"`
	}
	if err := ctx.ShouldBindJSON(&status_req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	orderStatus, err, _ := Model.ParseOrderStatus(status_req.Status)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = oc.OrderService.UpdateStatus(&status_req.Id, orderStatus)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (oc *OrderController) DeleteOrder(ctx *gin.Context) {
	orderid, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid order ID"})
		return
	}

	err = oc.OrderService.DeleteOrder(&orderid)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (oc *OrderController) RegisterOrderRoutes(rg *gin.RouterGroup) {
	orderroute := rg.Group("/order")
	orderroute.POST("", oc.CreateOrder)
	orderroute.GET("/:id", oc.GetOrder)
	orderroute.GET("/list/:status", oc.GetOrders)
	orderroute.GET("", oc.GetOrders)
	orderroute.PUT("/:id", oc.UpdateOrder)
	orderroute.PUT("/status", oc.UpdateStatus)
	orderroute.DELETE("/:id", oc.DeleteOrder)
}
