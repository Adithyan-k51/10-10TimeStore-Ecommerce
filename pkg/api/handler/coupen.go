package handler

import (
	"ecommerce/pkg/api/utilhandler"

	"ecommerce/pkg/commonhelp/requests.go"
	"ecommerce/pkg/commonhelp/response"
	"ecommerce/pkg/domain"
	services "ecommerce/pkg/usecase/interface"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type CouponHandler struct {
	CouponUsecase services.CouponUseCase
}

func NewCouponHandler(CouponUsecase services.CouponUseCase) *CouponHandler {
	return &CouponHandler{
		CouponUsecase: CouponUsecase,
	}
}

// AddCoupon godoc
// @summary api for add Coupons for ecommerce
// @description Admin can add coupon
// @security ApiKeyAuth
// @id AddCoupon
// @tags Coupon
// @Param input body requests.Coupon true "Input true info"
// @Router /admin/coupon/AddCoupons [post]
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
func (cr *CouponHandler) AddCoupon(ctx *gin.Context) {
	var body requests.Coupon
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "Failed to bind request body",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	var coupon domain.Coupon
	if err := copier.Copy(&coupon, &body); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			StatusCode: 500,
			Message:    "Failed to copy coupon data",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	if err := cr.CouponUsecase.CreateCoupon(ctx, coupon); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "Failed to create coupon",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "Successfully created coupon",
		Data:       body,
		Errors:     nil,
	})
}

// UpdateCoupon godoc
// @Summary Admin can update existing coupon
// @ID update-coupon
// @Description Admin can update existing coupon
// @Tags Coupon
// @Accept json
// @Produce json
// @Param CouponID path int true "CouponID"
// @Param coupon_details body requests.Coupon true "details of coupon to be updated"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /admin/coupon/Update/{CouponID} [patch]
func (cr *CouponHandler) UpdateCoupon(ctx *gin.Context) {
	id := ctx.Param("CouponID")
	couponID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "Invalid coupon ID",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	var body requests.Coupon
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "Failed to bind request body",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	updatedCoupon, err := cr.CouponUsecase.UpdateCouponById(ctx, couponID, body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "Failed to update coupon",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "Successfully updated coupon",
		Data:       updatedCoupon,
		Errors:     nil,
	})
}

// DeleteCoupon godoc
// @Summary Admin can delete a coupon
// @ID delete-coupon
// @Description Admin can delete a coupon
// @Tags Coupon
// @Accept json
// @Produce json
// @Param CouponID path string true "CouponID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /admin/coupon/Delete/{CouponID} [delete]
func (cr *CouponHandler) DeleteCoupon(ctx *gin.Context) {
	id := ctx.Param("CouponID")
	couponID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "Invalid coupon ID",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	if err := cr.CouponUsecase.DeleteCoupon(ctx, couponID); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "Failed to delete coupon",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "Successfully deleted coupon",
		Data:       nil,
		Errors:     nil,
	})
}

// ViewCoupon godoc
// @Summary Admins can see Coupons with coupon_id
// @ID find-Coupon-by-id
// @Description Admins can see Coupons with coupon_id
// @Tags Coupon
// @Accept json
// @Produce json
// @Param id path string true "CouponID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /admin/coupon/Viewcoupon/{id} [get]
func (cr *CouponHandler) ViewCoupon(ctx *gin.Context) {
	paramsId := ctx.Param("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "Invalid coupon ID",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	coupon, err := cr.CouponUsecase.ViewCoupon(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "Failed to fetch coupon",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "Coupon details",
		Data:       coupon,
		Errors:     nil,
	})
}

// Coupons godoc
// @Summary Get all coupons
// @ID List-all-coupons
// @Description Endpoint for getting all coupons
// @Tags Coupon
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /admin/coupon/couponlist [get]
func (cr *CouponHandler) Coupons(ctx *gin.Context) {
	coupons, err := cr.CouponUsecase.ViewCoupons(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "Failed to fetch coupons",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "List of coupons",
		Data:       coupons,
		Errors:     nil,
	})
}

// ApplyCoupon godoc
// @Summary User can apply a coupon to the cart
// @ID apply-coupon-to-cart
// @Description User can apply coupon to the cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param code path string true "Coupon code"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /coupon/apply [patch]
func (cr *CouponHandler) ApplyCoupon(ctx *gin.Context) {
	userID, err := utilhandler.GetUserIdFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "Failed to get user ID",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	code := ctx.Param("code") // <-- Use path parameter
	if code == "" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "Coupon code is required",
			Data:       nil,
			Errors:     "empty coupon code",
		})
		return
	}

	discount, err := cr.CouponUsecase.ApplyCoupontoCart(ctx, userID, code)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "Failed to apply coupon",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "Coupon applied successfully",
		Data:       map[string]interface{}{"discount": discount},
		Errors:     nil,
	})
}

// UserCoupons godoc
// @Summary Get all coupons for users
// @ID List-all-coupons-user
// @Description Endpoint for getting all coupons in user side
// @Tags Coupon
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /coupon/coupons [get]
func (cr *CouponHandler) UserCoupons(ctx *gin.Context) {
	coupons, err := cr.CouponUsecase.ViewCoupons(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "Failed to fetch coupons",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "List of coupons",
		Data:       coupons,
		Errors:     nil,
	})
}
