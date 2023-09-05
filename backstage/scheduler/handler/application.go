package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/model"
	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/model/apperrors"
)

/*
	@function: GetApplictaionAmount
	@description: obtain overall amount of application
*/
func (h *Handler) GetApplictaionAmount(c *gin.Context) {
	// extract user info from middleware
	// temp comment for debugging
	// user := c.MustGet("user")

	// get application type
	applicationType, ok := c.GetQuery("type")
	if !ok {
		log.WithFields(log.Fields{
			// "User email": user.(*model.User).Email,
			"HTTP URL": "application_amount",
		}).Warn("Failed to obtain application type in HTTP request url, abandoned")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewBadRequest("Failed to obtain application type in HTTP request url"),
		})

		return
	}
	if applicationType != model.APPLICATIOON_TYPE_STREAM && applicationType != model.APPLICATIOON_TYPE_CONSOLE {
		log.WithFields(log.Fields{
			// "User email":             user.(*model.User).Email,
			"HTTP URL":               "application_amount",
			"Given Application Type": applicationType,
		}).Warn("Unknown application type, abandoned")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewBadRequest("Unknown application type"),
		})

		return
	}

	ctx := c.Request.Context()

	applicationCount, err := h.ApplicationService.GetStreamApplicationsCount(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": applicationCount,
	})
}

/*
	@function: GetApplictaionList
	@description: obtain application list
*/
func (h *Handler) GetApplictaionList(c *gin.Context) {
	// extract user info from middleware
	// temp comment for debugging
	// user := c.MustGet("user")

	// get application type
	applicationType, ok := c.GetQuery("type")
	if !ok {
		log.WithFields(log.Fields{
			// "User email": user.(*model.User).Email,
			"HTTP URL": "application_list",
		}).Warn("Failed to obtain application type in HTTP request url, abandoned")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewBadRequest("Failed to obtain application type in HTTP request url"),
		})

		return
	}
	if applicationType != model.APPLICATIOON_TYPE_STREAM && applicationType != model.APPLICATIOON_TYPE_CONSOLE {
		log.WithFields(log.Fields{
			// "User email":             user.(*model.User).Email,
			"HTTP URL":               "application_list",
			"Given Application Type": applicationType,
		}).Warn("Unknown application type, abandoned")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewBadRequest("Unknown application type"),
		})

		return
	}

	// get query page
	pageNumberString, ok := c.GetQuery("page")
	if !ok {
		log.WithFields(log.Fields{
			// "User email": user.(*model.User).Email,
			"HTTP URL": "application_list",
		}).Warn("Failed to obtain page query in HTTP request url, abandoned")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewBadRequest("Failed to obtain page query in HTTP request url"),
		})

		return
	}

	// convert page number into int
	pageNumber, err := strconv.Atoi(pageNumberString)
	if err != nil {
		log.WithFields(log.Fields{
			// "User email": user.(*model.User).Email,
			"HTTP URL": "application_list",
		}).Warn("Failed to convert page number into int object, abandoned")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewBadRequest("Failed to convert page number into int object, abandoned"),
		})
	}

	// get query page size
	pageSizeString, ok := c.GetQuery("size")
	if !ok {
		log.WithFields(log.Fields{
			// "User email": user.(*model.User).Email,
			"HTTP URL": "application_list",
		}).Warn("Failed to obtain page size in HTTP request url, abandoned")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewBadRequest("Failed to obtain page size in HTTP request url"),
		})

		return
	}

	// convert page size into int
	pageSize, err := strconv.Atoi(pageSizeString)
	if err != nil {
		log.WithFields(log.Fields{
			// "User email": user.(*model.User).Email,
			"HTTP URL": "application_list",
		}).Warn("Failed to convert page size into int object, abandoned")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewBadRequest("Failed to convert page size into int object, abandoned"),
		})
	}

	// get order method
	orderBy, ok := c.GetQuery("order")
	if !ok {
		log.WithFields(log.Fields{
			// "User email": user.(*model.User).Email,
		}).Warn("Failed to obtain order method in HTTP request url, abandoned")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewBadRequest("Failed to obtain order method in HTTP request url"),
		})
	}
	if orderBy != model.ORDER_BY_UPDATE_TIME && orderBy != model.ORDER_BY_NAME && orderBy != model.ORDER_BY_USAGE_COUNT {
		log.WithFields(log.Fields{
			// "User email":         user.(*model.User).Email,
			"HTTP URL":           "application_list",
			"Given Order Method": orderBy,
		}).Warn("Unknown order method, abandoned")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewBadRequest("Unknown order method"),
		})

		return
	}

	ctx := c.Request.Context()

	// serve for desktop application
	if applicationType == model.APPLICATIOON_TYPE_STREAM {
		// query application list based on page number, page size and order method
		streamApplicationList, err := h.ApplicationService.GetStreamApplications(ctx, pageNumber, pageSize, orderBy)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"err": err,
			})
			return
		}

		// return queried applications if successed
		c.JSON(http.StatusOK, gin.H{
			"applications": streamApplicationList,
		})
		return
	}

	// TODO: serve for shell application
}

/*
	@function: GetApplictaionDetails
	@description: obtain application details
*/
func (h *Handler) GetApplictaionDetails(c *gin.Context) {
	// extract user info from middleware
	// temp comment for debugging
	// user := c.MustGet("user")

	// get application index
	applicationID, ok := c.GetQuery("id")
	if !ok {
		log.WithFields(log.Fields{
			// "User email": user.(*model.User).Email,
			"HTTP URL": "application_details",
		}).Warn("Failed to obtain application type in HTTP request url, abandoned")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewBadRequest("Failed to obtain application type in HTTP request url"),
		})

		return
	}

	// get application type
	applicationType, ok := c.GetQuery("type")
	if !ok {
		log.WithFields(log.Fields{
			// "User email": user.(*model.User).Email,
			"HTTP URL": "application_details",
		}).Warn("Failed to obtain application type in HTTP request url, abandoned")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewBadRequest("Failed to obtain application type in HTTP request url"),
		})

		return
	}
	if applicationType != model.APPLICATIOON_TYPE_STREAM && applicationType != model.APPLICATIOON_TYPE_CONSOLE {
		log.WithFields(log.Fields{
			// "User email":             user.(*model.User).Email,
			"HTTP URL":               "application_details",
			"Given Application Type": applicationType,
		}).Warn("Unknown application type, abandoned")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewBadRequest("Unknown application type"),
		})

		return
	}

	ctx := c.Request.Context()

	// obtain applictaion details and return corresponding http response
	if applicationType == model.APPLICATIOON_TYPE_STREAM {
		streamApplication, err := h.ApplicationService.GetStreamApplicationDetails(ctx, applicationID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"application": streamApplication,
		})
	}
}
