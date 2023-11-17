package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func EventsEndpointGet(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		option := c.Query("option", "")
		if option == "" {
			relationId, err := strconv.ParseInt(c.Query("relationId", "0"), 10, 64)
			if err != nil {
				log.Println(err.Error())
				return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
			}
			referenceId, err := strconv.ParseInt(c.Query("referenceId", "0"), 10, 64)
			if err != nil {
				log.Println(err.Error())
				return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
			}
			var tags []string
			if c.Query("tags", "") != "" {
				tags = strings.Split(c.Query("tags", ""), ",")
			} else {
				tags = []string{}
			}
			timeRange := []string{}
			if c.Query("timeRange", "") != "" {
				timeRange = strings.Split(c.Query("timeRangeBegin", ""), ",")
				if len(timeRange) != 2 {
					return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
				}
			}
			skip, err := strconv.ParseInt(c.Query("skip", "0"), 10, 64)
			if err != nil {
				log.Println(err.Error())
				return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
			}
			take, err := strconv.Atoi(c.Query("take", "10"))
			if err != nil {
				log.Println(err.Error())
				return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
			}
			result, err := EventsFilter(relationId, referenceId, tags, c.Query("detail", ""), timeRange, skip, take)
			if err != nil {
				log.Println(err.Error())
				return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
			}
			return c.JSON(result)
		}
	}
	if id > 0 {
		return c.SendStatus(200)
	}
	return c.Status(200).SendString("")
}
