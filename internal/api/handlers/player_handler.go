package handlers

import (
	"log"
	"net/http"
	"soccer-service/internal/core/models"
	"soccer-service/internal/core/services"
	"strings"

	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
)

const (
	PlayerPrefix = "ply-"
)

type PlayerHandler struct {
	service *services.PlayerService
}

func NewPlayerHandler(service *services.PlayerService) *PlayerHandler {
	return &PlayerHandler{service: service}
}

func generatePlayerID() string {
	return PlayerPrefix + uuid.New().String()
}

func (h *PlayerHandler) Create(c echo.Context) error {
	var player models.Player
	if err := c.Bind(&player); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	player.ID = generatePlayerID()

	if err := h.service.CreatePlayer(player); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, player)
}

func (h *PlayerHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if !strings.HasPrefix(id, PlayerPrefix) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid player ID format. Must start with 'ply-'",
		})
	}

	player, err := h.service.GetPlayer(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Player not found"})
	}

	return c.JSON(http.StatusOK, player)
}

func (h *PlayerHandler) Update(c echo.Context) error {
	id := c.Param("id")
	if !strings.HasPrefix(id, PlayerPrefix) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid player ID format. Must start with 'ply-'",
		})
	}

	var player models.Player
	if err := c.Bind(&player); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	player.ID = id

	if err := h.service.UpdatePlayer(player); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, player)
}

func (h *PlayerHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if !strings.HasPrefix(id, PlayerPrefix) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid player ID format. Must start with 'ply-'",
		})
	}

	if err := h.service.DeletePlayer(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *PlayerHandler) List(c echo.Context) error {
	log.Println("Fetching list of players...")

	players, err := h.service.ListPlayers()
	if err != nil {
		log.Printf("Error fetching players: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	log.Printf("Successfully retrieved %d players", len(players))
	
	return c.JSON(http.StatusOK, players)
}
