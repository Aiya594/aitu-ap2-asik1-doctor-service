package httpdoc

import (
	"net/http"

	usecase "github.com/Aiya594/doctor-service/internal/use-case"
	"github.com/gin-gonic/gin"
)

// Doctor Service Endpoints
//  POST /doctors - create a new doctor
//  GET /doctors/{id} - retrieve a doctor by ID
//  GET /doctors - list all doctors

// type DocHandler interface {
// 	CreateDoctor(ctx *gin.Context)
// 	GetDocByID(id string)
// 	GetAll()
// }

type DocHandler struct {
	us usecase.DocUseCase
}

type request struct {
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	Specialization string `json:"specialization"`
}

func NewDocHandler(us usecase.DocUseCase) *DocHandler {
	return &DocHandler{us: us}
}

func (h *DocHandler) CreateDoctor(ctx *gin.Context) {
	var doc request
	if err := ctx.ShouldBindJSON(&doc); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if err := h.us.CreateDoc(doc.FullName, doc.Email, doc.Specialization); err != nil {
		ctx.JSON(parseError(err), gin.H{
			"message": "failed to create doctor",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "doctor created",
	})
}

func (h *DocHandler) GetDoctor(ctx *gin.Context) {
	id := ctx.Param("id")

	doctor, err := h.us.GetDocbyID(id)
	if err != nil {
		ctx.JSON(parseError(err), gin.H{
			"message": "failed to get doctor",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, doctor)
}

func (h *DocHandler) GetAll(ctx *gin.Context) {
	doctors, err := h.us.ListDoctors()
	if err != nil {
		ctx.JSON(parseError(err), gin.H{
			"message": "failed to get doctors",
			"error":   err.Error(),
		})
		return
	}

	

	ctx.JSON(http.StatusOK, doctors)
}
