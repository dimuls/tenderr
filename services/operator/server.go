package operator

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

//go:embed ui/dist
var ui embed.FS

type Storage interface {
	Elements() ([]entity.Element, error)

	AddUserError(ue entity.UserError) error
	UserErrors() ([]entity.UserError, error)

	AddErrorNotification(n entity.ErrorNotification) error
	ResolveErrorNotification(id uuid.UUID, message string) error
	ErrorNotifications() ([]entity.ErrorNotification, error)

	AddErrorResolveWaiter(w entity.ErrorResolveWaiter) error
	ErrorResolveWaiterStats() ([]entity.ErrorResolveWaiterStats, error)
	RemoveErrorResolveWaiters(enID uuid.UUID) ([]entity.ErrorResolveWaiter, error)
}

type MessageSender interface {
	SendMessage(c entity.Contact, message string)
}

type Server struct {
	Addr          string
	Storage       Storage
	MessageSender MessageSender
	CORS          CORS
	Logger        *zap.Logger

	rules   map[*regexp.Regexp]*entity.Class
	rulesMx sync.RWMutex
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
		r.Get("/elements", s.getElements)

		r.Get("/user-errors", s.getUserErrors)
		r.Post("/user-errors", s.postUserErrors)

		r.Get("/error-notifications", s.getErrorNotifications)
		r.Post("/error-notifications", s.postErrorNotifications)
		r.Patch("/error-notifications", s.patchErrorNotifications)

		r.Post("/error-resolve-waiter", s.postErrorResolveWaiter)
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

func (s *Server) getElements(c *fiber.Ctx) error {
	es, err := s.Storage.Elements()
	if err != nil {
		return fmt.Errorf("get elements from storage: %w", err)
	}

	return c.JSON(es)
}

func (s *Server) getUserErrors(c *fiber.Ctx) error {
	ues, err := s.Storage.UserErrors()
	if err != nil {
		return fmt.Errorf("get user errors from storage: %w", err)
	}

	return c.JSON(ues)
}

func (s *Server) postUserErrors(c *fiber.Ctx) error {
	var ue entity.UserError

	err := json.Unmarshal(c.Body(), &ue)
	if err != nil {
		return fiber.ErrBadRequest
	}

	ue.ID, err = uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("generate id: %w", err)
	}

	err = s.Storage.AddUserError(ue)
	if err != nil {
		return fmt.Errorf("add user error to storage: %w", err)
	}

	return c.SendStatus(http.StatusCreated)
}

func (s *Server) getErrorNotifications(c *fiber.Ctx) error {
	ns, err := s.Storage.ErrorNotifications()
	if err != nil {
		return fmt.Errorf("get error notifications from storage: %w", err)
	}

	return c.JSON(ns)
}

func (s *Server) postErrorNotifications(c *fiber.Ctx) error {
	var en entity.ErrorNotification

	err := json.Unmarshal(c.Body(), &en)
	if err != nil {
		return fiber.ErrBadRequest
	}

	en.ID, err = uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("generate id: %w", err)
	}

	err = s.Storage.AddErrorNotification(en)
	if err != nil {
		return fmt.Errorf("add error notification to storage: %w", err)
	}

	return c.SendStatus(http.StatusCreated)
}

func (s *Server) patchErrorNotifications(c *fiber.Ctx) error {
	var r struct {
		ID      uuid.UUID `json:"id"`
		Message string    `json:"message"`
	}

	err := json.Unmarshal(c.Body(), &r)
	if err != nil {
		return fiber.ErrBadRequest
	}

	err = s.Storage.ResolveErrorNotification(r.ID, r.Message)
	if err != nil {
		return fmt.Errorf("resolve error notification in storage: %w", err)
	}

	// Запускаем фоновую джообу для обработки. Для примера просто горутину.
	go func() {
		ws, err := s.Storage.RemoveErrorResolveWaiters(r.ID)
		if err != nil {
			s.Logger.Error("remove error resolve waiters", zap.Error(err))
			return
		}

		for _, w := range ws {
			s.MessageSender.SendMessage(w.Contact, r.Message)
		}
	}()

	return c.SendStatus(http.StatusOK)
}

func (s *Server) postErrorResolveWaiter(c *fiber.Ctx) error {
	var w entity.ErrorResolveWaiter

	err := json.Unmarshal(c.Body(), &w)
	if err != nil {
		return fiber.ErrBadRequest
	}

	w.ID, err = uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("generate id: %w", err)
	}

	err = s.Storage.AddErrorResolveWaiter(w)
	if err != nil {
		return fmt.Errorf("add error resolve waiter to storage: %w", err)
	}

	return c.SendStatus(http.StatusCreated)
}
