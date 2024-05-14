package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

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
		fmt.Println("Json Marshal Failed")
		messageJson = []byte(`{
			"body":"Your tale was lost in the Astral Sea! Try again later!"
		}`)
	}
	http.Post("http://127.0.0.1:3000/", "application/json", bytes.NewBuffer(messageJson))
	return c.Render("greeting", fiber.Map{"Greeting": confirmation})
}

func main() {
	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})
	app.Static("/", "./static")
	app.Get("/", RenderForm)
	app.Post("/submit", ProcessForm)
	app.Listen(":8080")
}
