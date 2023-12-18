package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func EventEndpointGet(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		option := c.Query("option", "")
		if option == "default" {
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
			op := map[string]interface{}{
				"equal":         strings.Split(c.Query("equal", ""), ","),
				"objectContain": strings.Split(c.Query("object-contain", ""), ","),
				"arrayContain":  strings.Split(c.Query("array-contain", ""), ","),
				"like":          strings.Split(c.Query("like", ""), ","),
				"objectLike":    strings.Split(c.Query("object-like", ""), ","),
				"in":            strings.Split(c.Query("in", ""), ","),
				"lesser":        strings.Split(c.Query("lesser", ""), ","),
				"greater":       strings.Split(c.Query("greater", ""), ","),
			}
			result, err := Retrieve("event", skip, take, op)
			if err != nil {
				slogger.Error(err.Error())
				return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
			}
			if len(result) == 0 {
				return c.SendString("[]")
			}
			var response []EventExtended
			for _, it := range result {
				extendedEvent := EventExtended{
					Event: Event{
						Id:          it["id"].(int64),
						RelationId:  it["relation_id"].(int64),
						ReferenceId: it["reference_id"].(int64),
						Tags:        it["tags"].(string),
						Detail:      it["detail"].(string),
						Time:        it["time"].(string),
					},
					Id_:          strconv.FormatInt(it["id"].(int64), 10),
					RelationId_:  strconv.FormatInt(it["relation_id"].(int64), 10),
					ReferenceId_: strconv.FormatInt(it["reference_id"].(int64), 10),
				}
				response = append(response, extendedEvent)
			}
			return c.JSON(response)
		}
	}
	if id > 0 {
		return c.SendStatus(200)
	}
	return c.Status(200).SendString("")
}
