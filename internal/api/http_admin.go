package api

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/avpetkun/the-prime/internal/common"
	"github.com/avpetkun/the-prime/pkg/tgu"
)

func (api *HTTP) adminRequiredMiddleware(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	isAdmin, err := api.svc.ch.CheckUserAdmin(c.UserContext(), userID)
	if err != nil {
		api.log.Error().Err(err).
			Int64("user_id", userID).
			Msg("[api] failed to check admin")
		return c.SendStatus(http.StatusInternalServerError)
	}
	if !isAdmin {
		return c.SendStatus(http.StatusForbidden)
	}
	return c.Next()
}

func (api *HTTP) handlerAdminGetOverview() fiber.Handler {
	type Response struct {
		Chats    []tgu.Chat         `json:"chats"`
		Tasks    []*common.FullTask `json:"tasks"`
		Products []*common.Product  `json:"products"`
	}
	return func(c *fiber.Ctx) error {
		chats, err := api.svc.GetAllBotChats(c.UserContext())
		if err != nil {
			api.log.Warn().Err(err).Msg("[api] failed to get bot joined chats")
			return c.Status(http.StatusInternalServerError).JSON(err.Error())
		}
		tasks, err := api.svc.GetAllTasksWithClicks(c.UserContext())
		if err != nil {
			api.log.Warn().Err(err).Msg("[api] failed to get all tasks")
			return c.Status(http.StatusInternalServerError).JSON(err.Error())
		}
		products, err := api.svc.GetAllProducts(c.UserContext())
		if err != nil {
			api.log.Warn().Err(err).Msg("[api] failed to get all products")
			return c.Status(http.StatusInternalServerError).JSON(err.Error())
		}
		return c.JSON(Response{
			Chats:    chats,
			Tasks:    tasks,
			Products: products,
		})
	}
}

func (api *HTTP) handlerAdminSaveProduct(c *fiber.Ctx) error {
	p := new(common.Product)
	if err := c.BodyParser(p); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}
	if err := api.svc.SaveProduct(c.UserContext(), p); err != nil {
		api.log.Warn().Err(err).Msg("[api] failed to save product")
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}
	return c.JSON(p)
}

func (api *HTTP) handlerAdminDeleteProduct(c *fiber.Ctx) error {
	productID, err := strconv.ParseInt(c.Params("product_id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}
	if err := api.svc.DeleteProduct(c.UserContext(), productID); err != nil {
		api.log.Warn().Err(err).Msg("[api] failed to delete product")
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}
	return c.SendStatus(http.StatusNoContent)
}

func (api *HTTP) handlerAdminSaveTask(c *fiber.Ctx) error {
	task := new(common.FullTask)
	if err := c.BodyParser(task); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}
	if err := api.svc.SaveTask(c.UserContext(), task); err != nil {
		api.log.Warn().Err(err).Msg("[api] failed to save task")
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}
	return c.JSON(task)
}

func (api *HTTP) handlerAdminDeleteTask(c *fiber.Ctx) error {
	taskID, err := strconv.ParseInt(c.Params("task_id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}
	if err := api.svc.DeleteTask(c.UserContext(), taskID); err != nil {
		api.log.Warn().Err(err).Msg("[api] failed to delete task")
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}
	return c.SendStatus(http.StatusNoContent)
}

func (api *HTTP) handlerAdminRewardUser(c *fiber.Ctx) error {
	userID, err := strconv.ParseInt(c.Query("user"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}
	points, err := strconv.ParseInt(c.Query("points"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	if err = api.svc.ForceRewardUser(c.UserContext(), userID, points); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}
	return c.SendStatus(http.StatusOK)
}
