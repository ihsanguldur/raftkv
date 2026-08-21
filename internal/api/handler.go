package api

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ihsanguldur/raftkv/internal/service"
)

type Handler struct {
	svc        *service.KVService
	httpClient *http.Client
}

func NewHandler(svc *service.KVService) *Handler {
	return &Handler{
		svc:        svc,
		httpClient: &http.Client{Timeout: 3 * time.Second},
	}
}

func (h *Handler) proxyToLeader(c *fiber.Ctx, leaderAddr string) error {
	url := fmt.Sprintf("http://%s%s", leaderAddr, c.OriginalURL())
	req, err := http.NewRequest(c.Method(), url, bytes.NewReader(c.Body()))
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(errorResponse{Error: "leader unreachable"})
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Raft-Proxied", "true")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(errorResponse{Error: "leader unreachable"})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(errorResponse{Error: "leader unreachable"})
	}

	return c.Status(resp.StatusCode).Send(body)
}

func (h *Handler) Get(c *fiber.Ctx) error {
	v, err := h.svc.Get(c.Params("key"))
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(valueResponse{Value: v})
}

func (h *Handler) Put(c *fiber.Ctx) error {
	var req putRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{Error: "invalid body"})
	}
	key := strings.Clone(c.Params("key"))
	if err := h.svc.Put(key, req.Value); err != nil {
		if errors.Is(err, service.ErrNotLeader) && c.Get("X-Raft-Proxied") == "" {
			if leaderAddr := h.svc.LeaderAddr(); leaderAddr != "" {
				return h.proxyToLeader(c, leaderAddr)
			}
		}
		return writeError(c, err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	key := strings.Clone(c.Params("key"))
	if err := h.svc.Delete(key); err != nil {
		if errors.Is(err, service.ErrNotLeader) && c.Get("X-Raft-Proxied") == "" {
			if leaderAddr := h.svc.LeaderAddr(); leaderAddr != "" {
				return h.proxyToLeader(c, leaderAddr)
			}
		}
		return writeError(c, err)
	}
	return c.SendStatus(fiber.StatusOK)
}
