package handler

import (
	"ecommerce/pkg/api/utilhandler"
	"ecommerce/pkg/domain"
	"fmt"
	"net/http"
	"strconv"

	"ecommerce/pkg/commonhelp/requests.go"
	"ecommerce/pkg/commonhelp/response"
	services "ecommerce/pkg/usecase/interface"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	AdminUsecase services.AdminUsecase
}

func NewAdminHandler(Adusecase services.AdminUsecase) *AdminHandler {
	return &AdminHandler{
		AdminUsecase: Adusecase,
	}
}

// @Summary SaveAdmin
// @ID SaveAdmin
// @Description Save admin with details.
// @Tags Admin
// @Accept json
// @Produce json
// @Param   inputs   body     domain.Admin{}   true  "Input Field"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /admin/signup [post]
func (cr *AdminHandler) SaveAdmin(c *gin.Context) {
	var admin domain.Admin
	err := c.Bind(&admin)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.Response{
			StatusCode: 422,
			Message:    "can't bind",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	err = cr.AdminUsecase.SaveAdmin(c.Request.Context(), admin)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "unable signup",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response.Response{
		StatusCode: 201,
		Message:    "signup Successfully",
		Errors:     nil,
	})

}

// / @Summary LoginAdmin
// @ID LoginAdmin
// @Description login admin with details.
// @Tags Admin
// @Accept json
// @Produce json
// @Param   inputs   body     domain.Admin{}   true  "Input Field"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /admin/login  [post]
func (cr *AdminHandler) LoginAdmin(c *gin.Context) {
	var admin domain.Admin
	err := c.Bind(&admin)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.Response{
			StatusCode: 422,
			Message:    "can't bind",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	ss, err := cr.AdminUsecase.LoginAdmin(c.Request.Context(), admin)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "failed to login",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("AdminAuth", ss, 3600*24*1, "", "", false, true)
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "logged in successfuly",
		Data:       nil,
		Errors:     nil,
	})

}

// AdminLogout
// @Summary Adminlogout
// @ID AdminLogout
// @Description Logout as a user exit from the ecommerce site
// @Tags Admin
// @Success 200 "success"
// @Failure 400 "failed"
// @Router /admin/logout [post]
func (cr *AdminHandler) AdminLogout(c *gin.Context) {
	c.SetCookie("AdminAuth", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "logout successfully",
		Data:       "nil",
		Errors:     nil,
	})
}

type PaginationRequest struct {
	Page    uint `json:"page"`
	PerPage uint `json:"perpage"`
}

// @Summary FindAllUsers
// @Id FindAllUsers
// @Discription list of users.
// @tags Admin
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination"
// @Param perPage query int false "Number of items to retrieve per page"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /admin/findall [get]

func (cr *AdminHandler) FindAllUser(c *gin.Context) {
	var req PaginationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "invalid JSON body",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	fmt.Println(req)

	users, err := cr.AdminUsecase.FindAllUser(c.Request.Context(), requests.Pagination{
		Page:    req.Page,
		PerPage: req.PerPage,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "users not found",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "successfully fetched users",
		Data:       users,
		Errors:     nil,
	})
}

// BlockUser
// @Summary Admin can block a user
// @ID block-user
// @Description Admin can block a  user
// @Tags Admin
// @Accept json
// @Produce json
// @Param input body requests.BlockUser{} true "inputs"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 422 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/block [patch]

func (cr *AdminHandler) BlockUser(c *gin.Context) {
	var body requests.BlockUser
	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "bind faild",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	adminId, err := utilhandler.GetAdminIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "Can't find AdminId",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	err = cr.AdminUsecase.BlockUser(body, adminId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "Can't Block",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "User Blocked",
		Data:       nil,
		Errors:     nil,
	})
}

// UnblockUser
// @Summary Admin can unblock a blocked user
// @ID unblock-user
// @Description Admin can unblock a blocked user
// @Tags Admin
// @Accept json
// @Produce json
// @Param user_id path string true "ID of the user to be unblocked"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/unblock/{user_id} [patch]
func (cr *AdminHandler) UnblockUser(c *gin.Context) {
	paramsId := c.Param("user_id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "bind faild",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	err = cr.AdminUsecase.UnblockUser(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "cant unblock user",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "user unblocked",
		Data:       nil,
		Errors:     nil,
	})
}

// FindUserByID
// @Summary Admin can fetch a specific user details using user id
// @ID find-user-by-id
// @Description Admin can fetch a specific user details using user id
// @Tags Admin
// @Accept json
// @Produce json
// @Param user_id path string true "ID of the user to be fetched"
// @Success 200 {object} response.Response
// @Failure 422 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/finduser/{user_id} [get]

func (cr *AdminHandler) FindUserByID(c *gin.Context) {
	paramsID := c.Param("user_id")
	id, err := strconv.Atoi(paramsID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.Response{StatusCode: 422, Message: "failed to parse user id", Data: nil, Errors: err.Error()})
		return
	}
	user, err := cr.AdminUsecase.FindUserbyId(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{StatusCode: 500, Message: "failed fetch user", Errors: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200, Message: "Successfully fetched user details", Data: user, Errors: nil,
	})

}
