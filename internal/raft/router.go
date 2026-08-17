package raft

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, n *Node) {
	app.Post("/raft/request-vote", func(c *fiber.Ctx) error {
		var args RequestVoteArgs
		if err := c.BodyParser(&args); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}

		return c.JSON(n.HandleRequestVote(args))
	})

	app.Post("/raft/append-entries", func(c *fiber.Ctx) error {
		var args AppendEntriesArgs
		if err := c.BodyParser(&args); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}

		return c.JSON(n.HandleAppendEntries(args))
	})
}
