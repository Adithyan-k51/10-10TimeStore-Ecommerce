package http

import (
	"ecommerce/pkg/api/handler"
	"ecommerce/pkg/api/middleware"
	"log"
	"net/http"

	_ "ecommerce/cmd/api/docs" // Importing docs for Swagger

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(
	userHandler *handler.UserHandler,
	otpHandler *handler.OtpHandler,
	adminHandler *handler.AdminHandler,
	ProductHandler *handler.ProductHandler,
	CartHandler *handler.CartHandler,
	CouponHandler *handler.CouponHandler,
	OrderHandler *handler.OrderHandler,
) *ServerHTTP {
	engine := gin.Default()
	// engine.Use(gin.Logger())

	// Load HTML templates
	engine.LoadHTMLGlob("templates/*")

	// Swagger endpoint
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Validate required handlers
	if userHandler == nil || otpHandler == nil || adminHandler == nil {
		log.Fatal("handler dependencies cannot be nil")
	}

	// ==================== User Routes ====================
	user := engine.Group("/")
	{
		user.POST("signup", userHandler.UserSignup)
		user.GET("login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.html", nil)
		})
		user.POST("login", userHandler.UserLogin)
		user.POST("otp/send", otpHandler.SendOtp)
		user.POST("otp/verify", otpHandler.ValidateOtp)
		user.GET("home", userHandler.Home)
	}

	user.Use(middleware.UserAuth)
	{
		user.POST("SaveAddress", userHandler.AddAdress)
		user.PATCH("UpdateAddress", userHandler.UpdateAdress)
		user.GET("viewAddress", userHandler.VeiwAddress)
		user.POST("Addwishlist/:id", userHandler.AddToWishList)
		user.DELETE("/Removewishlist/:id", userHandler.RemoveFromWishList)
		user.GET("wishlist", userHandler.GetWishList)
		user.POST("logout", userHandler.UserLogout)

		category := user.Group("/category")
		{
			category.GET("showall/", ProductHandler.ListCategories)
			category.GET("disply/:id", ProductHandler.DisplayCategory)
		}

		product := user.Group("/product")
		{
			product.GET("ViewAllProducts", ProductHandler.ViewAllProducts)
			product.GET("/products/search", ProductHandler.SearchProducts)
			product.GET("/filter", ProductHandler.FilterProductsByPrice)
		}

		cart := user.Group("/cart")
		{
			cart.POST("/AddToCart", CartHandler.AddCartItem)
			cart.DELETE("/RemoveFromCart", CartHandler.RemoveFromCart)
			cart.PUT("/Addcount", CartHandler.Addcount)
			cart.GET("/viewcart", CartHandler.ViewCartItems)
		}

		coupon := user.Group("/coupon")
		{
			coupon.GET("/coupons", CouponHandler.UserCoupons)
			coupon.PATCH("/apply/:code", CouponHandler.ApplyCoupon)
		}

		order := user.Group("/order")
		{
			order.POST("/orderAll/:payment_id", OrderHandler.CashonDElivery)
			order.GET("/razor", OrderHandler.RazorpayCheckout)
			order.POST("/razor/success", OrderHandler.RazorpayVerify)
			order.PATCH("/cancel/:orderId", OrderHandler.CancelOrder)
			order.GET("/view/:order_id", OrderHandler.ListOrder)
			order.GET("/listall", OrderHandler.ListAllOrders)
			order.PATCH("/return/:orderId", OrderHandler.ReturnOrder)
		}
	}

	// ==================== Admin Routes ====================
	admin := engine.Group("/admin")
	{
		admin.POST("/signup", adminHandler.SaveAdmin)
		admin.POST("/login", adminHandler.LoginAdmin)
		admin.POST("/logout", adminHandler.AdminLogout)

		// Protected admin routes
		admin.Use(middleware.AdminAuth)
		{
			admin.GET("/findall", adminHandler.FindAllUser)
			admin.GET("/finduser/:user_id", adminHandler.FindUserByID)
			admin.POST("/finduser", adminHandler.FindUserByID)
			admin.PATCH("/block", adminHandler.BlockUser)
			admin.PATCH("/unblock/:user_id", adminHandler.UnblockUser)
		}

		// Category
		category := admin.Group("/category")
		{
			category.POST("add", ProductHandler.Addcategory)
			category.PATCH("update/:id", ProductHandler.UpdateCategory)
			category.DELETE("delete/:category_id", ProductHandler.DeleteCategory)
			category.GET("showall/", ProductHandler.ListCategories)
			category.GET("disply/:id", ProductHandler.DisplayCategory)
		}

		// Product
		product := admin.Group("/product")
		{
			product.POST("save", ProductHandler.SaveProduct)
			product.PATCH("updateproduct/:id", ProductHandler.UpdateProduct)
			product.DELETE("delete/:product_id", ProductHandler.DeleteProduct)
			product.GET("ViewAllProducts", ProductHandler.ViewAllProducts)
			product.GET("ViewProduct/:id", ProductHandler.VeiwProduct)
			// product.GET("/products/search", ProductHandler.SearchProducts)
		}

		// Order
		order := admin.Group("/order")
		{
			order.GET("/Status", OrderHandler.Statuses)
			order.GET("/Allorders", OrderHandler.AllOrders)
			order.PATCH("/UpdateStatus", OrderHandler.UpdateOrderStatus)
		}

		// Coupon
		coupon := admin.Group("/coupon")
		{
			coupon.POST("/AddCoupons", CouponHandler.AddCoupon)
			coupon.PATCH("/Update/:CouponID", CouponHandler.UpdateCoupon)
			coupon.DELETE("/Delete/:CouponID", CouponHandler.DeleteCoupon)
			coupon.GET("/Viewcoupon/:id", CouponHandler.ViewCoupon)
			coupon.GET("/couponlist", CouponHandler.Coupons)
		}
	}

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	if sh.engine == nil {
		log.Fatal("server engine is nil - initialization failed")
	}
	err := sh.engine.Run(":3002")
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
