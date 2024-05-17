package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/rs/zerolog"
)

var logger zerolog.Logger

type Message struct {
	Content string `json:"content"`
}

func RenderForm(c *fiber.Ctx) error {
	return c.Render("form", fiber.Map{})
}

func ProcessForm(c *fiber.Ctx) error {
	bodyText := c.FormValue("body")
	confirmation := "Your tale is spun!"
	m := Message{
		Content: bodyText,
	}
	messageJson, err := json.Marshal(m)
	if err != nil {
		logger.Error().Msg("Json Marshal Failed")
		messageJson = []byte(`{
			"body":"Your tale was lost in the Astral Sea! Try again later!"
		}`)
	}
	http.Post("anon-bot-messager:8000/", "application/json", bytes.NewBuffer(messageJson))
	return c.Render("greeting", fiber.Map{"Greeting": confirmation})
}

func main() {
	//Zerolog setup
	logFile, err := os.OpenFile(
		"frontend.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		//Don't run without logging
		panic(err)
	}

	//Clean up after ourselves
	defer logFile.Close()

	logger = zerolog.New(logFile).With().Timestamp().Logger()

	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})
	app.Static("/", "./static")
	app.Get("/", RenderForm)
	app.Post("/submit", ProcessForm)
	app.Listen(":8080")
}
