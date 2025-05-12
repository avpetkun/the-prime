package webhook

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/avpetkun/the-prime/internal/worker"
)

type Publisher func(subject string, message any) error

func handlerReward(publish Publisher) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		token := c.Params("token")

		rawUserID := c.Query("userid")
		userID := int64(0)
		taskID := int64(0)
		i := strings.IndexByte(rawUserID, '-')
		if i == -1 {
			userID, err = strconv.ParseInt(rawUserID, 10, 64)
		} else {
			userID, err = strconv.ParseInt(rawUserID[:i], 10, 64)
			if err == nil {
				taskID, err = strconv.ParseInt(rawUserID[i+1:], 10, 64)
			}
		}
		if err != nil {
			return c.SendStatus(http.StatusBadRequest)
		}

		err = publish(worker.KeyWebhookReward, worker.WebhookRewardMessage{
			Received: time.Now(),
			Token:    token,
			UserID:   userID,
			TaskID:   taskID,
		})
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		return c.SendStatus(http.StatusOK)
	}
}

func handlerRewardTappAds(publish Publisher) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Params("token")
		userID, err := strconv.ParseInt(c.Query("telegram_user_id"), 10, 64)
		if err != nil {
			return c.SendStatus(http.StatusBadRequest)
		}
		payout, err := strconv.ParseFloat(c.Query("payout"), 64)
		if err != nil {
			return c.SendStatus(http.StatusBadRequest)
		}
		subtaskID, err := strconv.ParseInt(c.Query("offer"), 10, 64)
		if err != nil {
			return c.SendStatus(http.StatusBadRequest)
		}

		err = publish(worker.KeyWebhookReward, worker.WebhookRewardMessage{
			Received: time.Now(),
			Token:    token,
			UserID:   userID,
			SubID:    subtaskID,
			Payout:   payout,
		})
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		return c.SendStatus(http.StatusOK)
	}
}
