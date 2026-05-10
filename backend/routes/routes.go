package routes

import (
	"ChiQuoc/HocGolang/handlers"
	"ChiQuoc/HocGolang/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(route *gin.Engine) {
	api := route.Group("/api")

	api.POST("/auth/login", handlers.LoginHandler)

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddlewareJWT())
	{
		// thông tin bản thân user
		protected.GET("/auth/me", handlers.MeHandler)

		// Thông tin dashboard
		protected.GET("/dashboard", handlers.GetDashboardStats)

		// Mọi user xem được phòng ban và chức vụ
		protected.GET("/departments", handlers.GetDepartments)
		protected.GET("/departments/:id", handlers.GetDepartmentByID)
		protected.GET("/departments/:id/positions", handlers.GetPositionsByDepartment)

		// Thêm sửa xoá phòng ban chỉ admin được làm
		adminDept := protected.Group("/departments")
		adminDept.Use(middleware.AdminOnlyMiddleware())
		{
			adminDept.POST("", handlers.CreateDepartment)
			adminDept.PUT("/:id", handlers.UpdateDepartment)
			adminDept.DELETE("/:id", handlers.DeleteDepartment)
		}

		// Thêm sửa xoá chức vụ chỉ admin được làm
		adminPos := protected.Group("/positions")
		adminPos.Use(middleware.AdminOnlyMiddleware())
		{
			adminPos.POST("", handlers.CreatePosition)
			adminPos.PUT("/:id", handlers.UpdatePosition)
			adminPos.DELETE("/:id", handlers.DeletePosition)
		}

		// Mọi user có thể xem được nhân viên
		protected.GET("/employees", handlers.GetEmployeeList)
		protected.GET("/employees/:id", handlers.GetEmployeeID)

		// Thêm sửa xoá nhân viên chỉ admin được làm
		adminEmp := protected.Group("/employees")
		adminEmp.Use(middleware.AdminOnlyMiddleware())
		adminEmp.POST("", handlers.CreateEmployee)
		adminEmp.PUT("/:id", handlers.UpdateEmployee)
		adminEmp.DELETE("/:id", handlers.DeleteEmployee)

		// Thêm sửa xoá tài khoản chỉ admin được làm
		users := protected.Group("/users")
		users.Use(middleware.AdminOnlyMiddleware())
		{
			users.GET("", handlers.GetUsers)
			users.GET("/:id", handlers.GetUserByID)
			users.POST("", handlers.CreateUser)
			users.PUT("/:id", handlers.UpdateUser)
			users.DELETE("/:id", handlers.DeleteUser)
		}
	}
}
