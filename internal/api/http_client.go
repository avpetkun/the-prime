package api

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/avpetkun/the-prime/internal/common"
	"github.com/avpetkun/the-prime/pkg/tgu"
)

func (api *HTTP) userAuthMiddleware(c *fiber.Ctx) (err error) {
	var auth *tgu.Auth
	if token := c.Get(fiber.HeaderAuthorization); token != "" {
		auth, err = api.svc.ParseUserAuth(token)
	} else if token := c.Query("auth"); token != "" {
		auth, err = api.svc.ParseUserAuth64(token) // deprecated
	} else {
		api.log.Warn().Str("url", c.BaseURL()).Msg("[api] request without auth")
		return c.SendStatus(http.StatusUnauthorized)
	}
	if err != nil {
		api.log.Warn().Err(err).Str("url", c.BaseURL()).Msg("[api] request invalid auth")
		return c.SendStatus(http.StatusForbidden)
	}
	c.Locals("auth", auth)
	c.Locals("user_id", auth.User.ID)
	return c.Next()
}

func (api *HTTP) handlerUserInit(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	if err := api.svc.ProcessUserInit(c.UserContext(), userID); err != nil {
		api.log.Error().Err(err).
			Int64("user_id", userID).
			Msg("[api] failed to user init")
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.SendStatus(http.StatusOK)
}

func (api *HTTP) handlerGetOverview() fiber.Handler {
	type Response struct {
		IsAdmin bool `json:"isAdmin,omitempty"`

		Points    int64 `json:"points"`
		RefCount  int64 `json:"refCount"`
		RefPoints int64 `json:"refPoints"`

		Products []*common.Product `json:"products"`
		Tasks    []common.UserTask `json:"tasks"`
	}
	return func(c *fiber.Ctx) (err error) {
		ctx := c.UserContext()
		user := c.Locals("auth").(*tgu.Auth).User
		start := c.Query("start")

		userIP := getFiberUserIP(c)
		userAgent := c.Get(fiber.HeaderUserAgent)

		isAdmin, _ := api.svc.ch.CheckUserAdmin(ctx, user.ID)
		refCount, refPoints, err := api.svc.GetUserRefs(ctx, user.ID)
		if err != nil {
			api.log.Warn().Err(err).Int64("user_id", user.ID).Msg("[api] failed to get user refs")
		}
		tasks, points, err := api.svc.GetUserProgress(ctx, user, start, userIP, userAgent)
		if err != nil {
			api.log.Warn().Err(err).Int64("user_id", user.ID).Msg("[api] failed to get user tasks")
			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.JSON(Response{
			IsAdmin:   isAdmin,
			Points:    points,
			RefCount:  refCount,
			RefPoints: refPoints,
			Products:  api.svc.AllProducts(),
			Tasks:     tasks,
		})
	}
}

func (api *HTTP) handlerGetTasksEvents(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	fromTime, err := strconv.ParseInt(c.Query("from"), 10, 64)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	events, err := api.svc.GetTasksEvents(c.UserContext(), userID, fromTime)
	if err != nil {
		api.log.Warn().Err(err).Int64("user_id", userID).Msg("[api] failed to get user tasks events")
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.JSON(events)
}

func (api *HTTP) handlerTaskStart(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	taskID, err := strconv.ParseInt(c.Params("task_id"), 10, 64)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	subID, err := strconv.ParseInt(c.Params("sub_id"), 10, 64)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	userIP := getFiberUserIP(c)
	userAgent := c.Get(fiber.HeaderUserAgent)

	if err := api.svc.TaskStart(c.UserContext(), userID, taskID, subID, userIP, userAgent); err != nil {
		api.log.Warn().Err(err).
			Int64("user_id", userID).
			Int64("task_id", taskID).
			Int64("sub_id", subID).
			Msg("[api] failed to start user task")
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.SendStatus(http.StatusOK)
}

func (api *HTTP) handlerTaskClaim(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	taskID, err := strconv.ParseInt(c.Params("task_id"), 10, 64)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	subID, err := strconv.ParseInt(c.Params("sub_id"), 10, 64)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	if err := api.svc.TaskClaim(c.UserContext(), userID, taskID, subID); err != nil {
		api.log.Warn().Err(err).
			Int64("user_id", userID).
			Int64("task_id", taskID).
			Int64("sub_id", subID).
			Msg("[api] failed to claim user task")
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.SendStatus(http.StatusOK)
}

func (api *HTTP) handlerTaskStarsInvoice(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	taskID, err := strconv.ParseInt(c.Params("task_id"), 10, 64)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	invoiceLink, err := api.svc.GetTaskStarsInvoice(c.UserContext(), userID, taskID)
	if err != nil {
		api.log.Warn().Err(err).
			Int64("user_id", userID).
			Int64("task_id", taskID).
			Msg("[api] failed to get task stars invoice")
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.JSON(invoiceLink)
}

func (api *HTTP) handlerTaskTonInvoice(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	taskID, err := strconv.ParseInt(c.Params("task_id"), 10, 64)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	tx, err := api.svc.GetTaskTonTransaction(c.UserContext(), userID, taskID)
	if err != nil {
		api.log.Warn().Err(err).
			Int64("user_id", userID).
			Int64("task_id", taskID).
			Msg("[api] failed to generate task stars invoice")
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.JSON(tx)
}

func (api *HTTP) handlerProductClaim(c *fiber.Ctx) error {
	productID, err := strconv.ParseInt(c.Params("product_id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}
	user := c.Locals("auth").(*tgu.Auth).User

	spendPoints, err := api.svc.ClaimProduct(c.UserContext(), user, productID)
	if err != nil {
		api.log.Warn().Err(err).
			Int64("product_id", productID).
			Int64("user_id", user.ID).
			Str("username", user.Username).
			Msg("[api] failed to claim product")
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.JSON(spendPoints)
}

func (api *HTTP) handlerGetInviteMessage(c *fiber.Ctx) error {
	user := c.Locals("auth").(*tgu.Auth).User

	msgID, err := api.svc.GetInviteMessage(c.UserContext(), user)
	if err != nil {
		api.log.Warn().Err(err).
			Int64("user_id", user.ID).
			Msg("[api] failed to get invite message")
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.JSON(msgID)
}

func getFiberUserIP(c *fiber.Ctx) string {
	ip := c.Get("X-Real-IP")
	if ip == "" {
		ip = c.Get("X-Forwarded-For")
		if ip == "" {
			ip = c.Context().RemoteAddr().String()
		}
	}
	return ip
}
