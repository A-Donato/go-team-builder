package handlers

import (
	"net/http"
	"soccer-service/internal/core/models"
	"soccer-service/internal/core/services"

	"github.com/labstack/echo/v4"
)

type PlayerHandler struct {
	service *services.PlayerService
}

func NewPlayerHandler(service *services.PlayerService) *PlayerHandler {
	return &PlayerHandler{service: service}
}

func (h *PlayerHandler) Create(c echo.Context) error {
	var player models.Player
	if err := c.Bind(&player); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if err := h.service.CreatePlayer(player); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, player)
}

func (h *PlayerHandler) Get(c echo.Context) error {
	id := c.Param("id")
	player, err := h.service.GetPlayer(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Player not found"})
	}

	return c.JSON(http.StatusOK, player)
}

func (h *PlayerHandler) Update(c echo.Context) error {
	var player models.Player
	if err := c.Bind(&player); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if err := h.service.UpdatePlayer(player); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, player)
}

func (h *PlayerHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.DeletePlayer(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *PlayerHandler) List(c echo.Context) error {
	players, err := h.service.ListPlayers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, players)
} 