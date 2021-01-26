package handler

import "github.com/gofiber/fiber/v2"

type GeneralHandler struct {
	//instance of database repostitory or service
}

// GET returns message "hi"
func (gh *GeneralHandler) SayHi(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"messgae" : "Hi"})
}

// GET person matching id
func (gh *GeneralHandler) GetPersonMatchingId(c *fiber.Ctx) error {
	id := c.Params("id")

	data := fiber.Map{
		"name" : "david",
		"id" : id,
	}

	// Set header
	c.SendStatus(200)
	// Set JSON response 
	return c.JSON(data)
}

// POST person matching id
func (gh *GeneralHandler) PersonCreds(c *fiber.Ctx) error {
	type Person struct {
			Name string `json:"name" xml:"name" form:"name"`
			Id int `json:"id" xml:"id" form:"id"`
	}
	person := new(Person)

	err := c.BodyParser(person)
	if err != nil {
		c.SendStatus(500)
		return c.JSON(fiber.Map{"message" : err})
	}

	data := fiber.Map{
		"name" : person.Name,
		"id" : person.Id,
	}

	// Set header
	c.SendStatus(200)
	// Set JSON response 
	return c.JSON(data)
}
