package http

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"scrum-poker/backend/internal/config"
	"scrum-poker/backend/internal/service"
	"scrum-poker/backend/internal/ws"
)

type Dependencies struct {
	Config  config.Config
	Service *service.RoomService
	Hub     *ws.Hub
}

func Register(app *fiber.App, deps Dependencies) {
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: deps.Config.AppOrigin,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,OPTIONS",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"ok": true})
	})

	api := app.Group("/api")

	api.Post("/rooms", func(c *fiber.Ctx) error {
		var body struct {
			DisplayName string `json:"displayName"`
		}
		if err := c.BodyParser(&body); err != nil {
			return writeError(c, fiber.StatusBadRequest, "درخواست نامعتبر است.")
		}
		session, err := deps.Service.CreateRoom(c.UserContext(), body.DisplayName)
		if err != nil {
			return handleError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(session)
	})

	api.Post("/rooms/join", func(c *fiber.Ctx) error {
		var body struct {
			DisplayName string `json:"displayName"`
			Code        string `json:"code"`
		}
		if err := c.BodyParser(&body); err != nil {
			return writeError(c, fiber.StatusBadRequest, "درخواست نامعتبر است.")
		}
		session, err := deps.Service.JoinRoom(c.UserContext(), body.Code, body.DisplayName)
		if err != nil {
			return handleError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(session)
	})

	api.Get("/rooms/:code", func(c *fiber.Ctx) error {
		snapshot, err := deps.Service.GetRoomSnapshot(c.UserContext(), c.Params("code"), bearerToken(c))
		if err != nil {
			return handleError(c, err)
		}
		return c.JSON(snapshot)
	})

	api.Get("/rooms/:code/history", func(c *fiber.Ctx) error {
		history, err := deps.Service.GetRoomHistory(c.UserContext(), c.Params("code"), bearerToken(c))
		if err != nil {
			return handleError(c, err)
		}
		return c.JSON(fiber.Map{"items": history})
	})

	api.Post("/rooms/:code/tasks", func(c *fiber.Ctx) error {
		var body struct {
			Title string `json:"title"`
		}
		if err := c.BodyParser(&body); err != nil {
			return writeError(c, fiber.StatusBadRequest, "درخواست نامعتبر است.")
		}
		snapshot, err := deps.Service.CreateTask(c.UserContext(), c.Params("code"), bearerToken(c), body.Title)
		if err != nil {
			return handleError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(snapshot)
	})

	api.Post("/rooms/:code/tasks/:taskId/votes", func(c *fiber.Ctx) error {
		taskID, err := strconv.ParseInt(c.Params("taskId"), 10, 64)
		if err != nil {
			return writeError(c, fiber.StatusBadRequest, "شناسه تسک نامعتبر است.")
		}
		var body struct {
			Value string `json:"value"`
		}
		if err := c.BodyParser(&body); err != nil {
			return writeError(c, fiber.StatusBadRequest, "درخواست نامعتبر است.")
		}
		snapshot, err := deps.Service.SubmitVote(c.UserContext(), c.Params("code"), taskID, bearerToken(c), body.Value)
		if err != nil {
			return handleError(c, err)
		}
		return c.JSON(snapshot)
	})

	app.Use("/ws/rooms/:code", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("roomCode", c.Params("code"))
			c.Locals("token", c.Query("token"))
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/rooms/:code", websocket.New(func(conn *websocket.Conn) {
		roomCode, _ := conn.Locals("roomCode").(string)
		token, _ := conn.Locals("token").(string)

		ctx := context.Background()
		if err := deps.Service.HandlePresenceConnected(ctx, roomCode, token); err != nil {
			log.Printf("presence connect failed: %v", err)
			_ = conn.WriteJSON(fiber.Map{"type": "error", "message": "جلسه معتبر نیست."})
			_ = conn.Close()
			return
		}

		client := deps.Hub.Register(roomCode, token, conn)
		defer func() {
			if err := deps.Service.HandlePresenceDisconnected(context.Background(), roomCode, client.MemberID); err != nil {
				log.Printf("presence disconnect failed: %v", err)
			}
			deps.Hub.Unregister(client)
		}()

		_ = conn.WriteJSON(fiber.Map{
			"type":      "system.connected",
			"roomCode":  roomCode,
			"timestamp": time.Now().UTC(),
		})

		for {
			var message map[string]any
			if err := conn.ReadJSON(&message); err != nil {
				return
			}
			if message["type"] == "heartbeat" {
				if err := deps.Service.HandlePresenceHeartbeat(context.Background(), roomCode, token); err != nil {
					log.Printf("presence heartbeat failed: %v", err)
					return
				}
			}
		}
	}))
}

func bearerToken(c *fiber.Ctx) string {
	header := c.Get("Authorization")
	if strings.HasPrefix(header, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	}
	return ""
}

func handleError(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*service.AppError); ok {
		return writeError(c, appErr.Status, appErr.Message)
	}
	log.Printf("unexpected error: %v", err)
	return writeError(c, fiber.StatusInternalServerError, "خطای غیرمنتظره‌ای رخ داد.")
}

func writeError(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"message": message,
	})
}
