package routes

import (
	welcome "kiripos/src/controllers"
	"kiripos/src/controllers/auth"
	"kiripos/src/controllers/masters"
	"kiripos/src/controllers/orders"
	"kiripos/src/controllers/reports"
	"kiripos/src/middlewares"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	router.GET("/", welcome.Index)

	// AUTH
	_auth := router.Group("/auth")
	{
		_auth.POST("/login", auth.Login)
	}

	// V1
	authorized := router.Group("/v1")
	authorized.Use(middlewares.Authorization())
	{
		// MASTERS
		_roles := authorized.Group(("/roles"))
		{
			_roles.GET("/", masters.RoleList)
		}

		_users := authorized.Group(("/users"))
		{
			_users.GET("/", masters.UserList)
			_users.GET("/:id", masters.UserDetail)
			_users.POST("/", masters.UserInsert)
			_users.PUT("/", masters.UserUpdate)
			_users.DELETE("/:id", masters.UserDelete)
		}

		_branches := authorized.Group(("/branches"))
		{
			_branches.GET("/", masters.BranchList)
			_branches.GET("/:id", masters.BranchDetail)
			_branches.POST("/", masters.BranchInsert)
			_branches.PUT("/", masters.BranchUpdate)
			_branches.DELETE("/:id", masters.BranchDelete)
		}

		_branche_users := authorized.Group(("/branch_users"))
		{
			_branche_users.GET("/:branch_id", masters.BranchUserList)
			_branche_users.GET("/:branch_id/:id", masters.BranchUserDetail)
			_branche_users.POST("/:branch_id", masters.BranchUserInsert)
			_branche_users.DELETE("/:branch_id/:id", masters.BranchUserDelete)
		}

		_categories := authorized.Group(("/categories"))
		{
			_categories.GET("/", masters.CategoryList)
			_categories.GET("/:id", masters.CategoryDetail)
			_categories.POST("/", masters.CategoryInsert)
			_categories.PUT("/", masters.CategoryUpdate)
			_categories.DELETE("/:id", masters.CategoryDelete)
		}

		_products := authorized.Group(("/products"))
		{
			_products.GET("/", masters.ProductList)
			_products.GET("/:id", masters.ProductDetail)
			_products.POST("/", masters.ProductInsert)
			_products.PUT("/", masters.ProductUpdate)
			_products.DELETE("/:id", masters.ProductDelete)
		}

		_customers := authorized.Group(("/customers"))
		{
			_customers.GET("/", masters.CustomerList)
			_customers.GET("/:id", masters.CustomerDetail)
			_customers.POST("/", masters.CustomerInsert)
			_customers.PUT("/", masters.CustomerUpdate)
			_customers.DELETE("/:id", masters.CustomerDelete)
		}

		_suppliers := authorized.Group(("/suppliers"))
		{
			_suppliers.GET("/", masters.SupplierList)
			_suppliers.GET("/:id", masters.SupplierDetail)
			_suppliers.POST("/", masters.SupplierInsert)
			_suppliers.PUT("/", masters.SupplierUpdate)
			_suppliers.DELETE("/:id", masters.SupplierDelete)
		}

		// TRANSACTIONS
		_trx := authorized.Group(("/orders"))
		{
			_trx.GET("/", orders.OrderList)
			_trx.GET("/:id", orders.OrderDetail)
			_trx.POST("/", orders.OrderCreate)
			_trx.PUT("/", orders.OrderUpdate)
			_trx.DELETE("/:id", orders.OrderDelete)
		}

		_purchase := authorized.Group(("/purchase"))
		{
			_purchase.GET("/", orders.PurchaseList)
			_purchase.GET("/:id", orders.PurchaseDetail)
			_purchase.POST("/", orders.PurchaseCreate)
			_purchase.PUT("/", orders.PurchaseUpdate)
			_purchase.DELETE("/:id", orders.PurchaseDelete)
		}

		// REPORTS
		_reports := authorized.Group(("/reports"))
		{
			_reports.GET("/dashboard", reports.Dashboard)
			_reports.GET("/orders", reports.ReportOrder)
			_reports.GET("/purchase", reports.ReportPurchase)
		}
	}
}
