package backend

import (
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/yellowphil/go-smtp-relay/pkg/db"
	"github.com/yellowphil/go-smtp-relay/pkg/utils"
)

func HandleRegister(ctx fiber.Ctx) error {
	rr := new(RegisterRequest)
	if err := ctx.Bind().Body(rr); err != nil {
		return err
	}
	if err := db.GetUserStorage().AddUser(rr.Username, rr.Password); err != nil {
		log.Errorf("Error registering user: %v\n", err)
		return fmt.Errorf("Could not add user")
	}
	return nil
}

func HandleLogin(ctx fiber.Ctx) error {
	sess := session.FromContext(ctx)
	lr := new(LoginRequest)
	if err := ctx.Bind().Body(lr); err != nil {
		return err
	}
	if pwdHash, err := db.GetUserStorage().GetUser(lr.Username); err != nil {
		log.Errorf("Error getting user: %v\n", err)
		return fmt.Errorf("Incorrect username or password")
	} else {
		if inputHash := utils.Sha3SumString(lr.Password); reflect.DeepEqual(inputHash, pwdHash) {
			return fmt.Errorf("Incorrect username or password")
		}
	}
	return nil
}

func NewBackend() *fiber.App {
	app := fiber.New()
	sessionStore := session.New(session.Config{
		Storage: &db.BadgerSessionStorage{},
	})
	app.Use(sessionStore)
	app.Use(csrf.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	return app
}
