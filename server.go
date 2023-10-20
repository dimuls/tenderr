package tenderr

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"tednerr/entity"
)

const userIDLocal = "userID"

//go:embed ui/dist
var ui embed.FS

type Storage interface {
	Classes() (cs []entity.Class, err error)
	SetClass(c entity.Class) error
	RemoveClass(classID uuid.UUID) error
}

type Server struct {
	Addr    string
	Storage Storage
	CORS    CORS
	Logger  *zap.Logger

	rules   map[*regexp.Regexp]entity.Class
	rulesMx *sync.RWMutex
}

func notStatic(c *fiber.Ctx) bool {
	path := string(c.Request().URI().Path())
	return strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/ws/")
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	ui, err := fs.Sub(ui, "ui/dist")
	if err != nil {
		return fmt.Errorf("init ui fs: %w", err)
	}

	app := fiber.New(fiber.Config{
		AppName:               "tenderr",
		ReadTimeout:           10 * time.Second,
		IdleTimeout:           10 * time.Second,
		BodyLimit:             2 * 1024 * 1024,
		DisableStartupMessage: true,
	})

	app.Use(recover.New())

	app.Use(func(c *fiber.Ctx) error {
		const msg = "request handled"

		var (
			st    = time.Now()
			log   func(msg string, fields ...zap.Field)
			flags []zap.Field
		)

		err := c.Next()

		if err != nil && app.Config().ErrorHandler(c, err) != nil {
			c.SendStatus(http.StatusInternalServerError)
		}

		if userID := c.Locals(userIDLocal); userID != nil {
			flags = append(flags, zap.String("userId", userID.(uuid.UUID).String()))
		}

		flags = append(flags,
			zap.Duration("latency", time.Since(st)),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
		)

		if qp := c.Request().URI().QueryArgs(); qp.Len() > 0 {
			qp.Del("token")
			zap.String("queryParams", qp.String())
		}

		flags = append(flags,
			zap.Int("statusCode", c.Response().StatusCode()),
			zap.String("clientAddr", c.IP()+":"+c.Port()),
			zap.Int("bytesReceived", len(c.Request().Body())),
			zap.Int("bytesSent", len(c.Response().Body())),
		)

		if err != nil {
			logger := s.Logger.With(zap.Error(err))

			switch err.(type) {
			case *fiber.Error:
				log = logger.Warn
			default:
				log = logger.Error
			}
		} else if notStatic(c) {
			log = s.Logger.Info
		} else {
			log = s.Logger.Debug
		}

		log(msg, flags...)

		return nil
	})

	app.Use(recover.New())

	if s.CORS.Enabled {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     s.CORS.AllowedOrigins,
			AllowMethods:     s.CORS.AllowedMethods,
			AllowHeaders:     s.CORS.AllowedHeaders,
			AllowCredentials: s.CORS.AllowCredentials,
		}))
	}

	app.Use("/", filesystem.New(filesystem.Config{
		Next:         notStatic,
		Root:         http.FS(ui),
		Index:        "index.html",
		NotFoundFile: "index.html",
	}))

	app.Route("/api", func(r fiber.Router) {
		r.Get("/classes", s.getClasses)
		r.Post("/classes", s.putClasses)
		r.Delete("/classes/:id", s.deleteClass)
	})

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		var err error

		err = app.Listen(s.Addr)
		if err != nil {
			s.Logger.Error("http server listen and serve", zap.Error(err))
			cancel()
			return
		}

		s.Logger.Info("http server shutdown")
	}()

	s.Logger.Info("server started")

	<-ctx.Done()

	s.Logger.Info("server stopping")

	err = app.ShutdownWithTimeout(10 * time.Second)
	if err != nil {
		s.Logger.Error("http server graceful shutdown", zap.Error(err))
	}

	wg.Wait()

	s.Logger.Info("server stopped")

	return nil
}

func (s *Server) getClasses(c *fiber.Ctx) error {
	classes, err := s.Storage.Classes()
	if err != nil {
		return fmt.Errorf("get classes from storage: %w", err)
	}
	return c.JSON(classes)
}

func (s *Server) putClasses(c *fiber.Ctx) error {
	var class entity.Class

	err := json.Unmarshal(c.Body(), &c)
	if err != nil {
		return fiber.ErrBadRequest
	}

	err = s.Storage.SetClass(class)
	if err != nil {
		return fmt.Errorf("set class in storage: %w", err)
	}

	return c.SendStatus(http.StatusOK)
}

func (s *Server) deleteClass(c *fiber.Ctx) error {
	var id uuid.UUID

	err := json.Unmarshal(c.Body(), &id)
	if err != nil {
		return fiber.ErrBadRequest
	}

	err = s.Storage.RemoveClass(id)
	if err != nil {
		return fmt.Errorf("remove class from storage: %w", err)
	}

	return c.SendStatus(http.StatusOK)
}
