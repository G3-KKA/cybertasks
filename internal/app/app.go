package app

import (
	"context"
	"cybertask/config"
	"cybertask/internal/controller/handler"
	"cybertask/internal/controller/httpctl"
	"cybertask/internal/db/pgsql"
	"cybertask/internal/logger"
	"cybertask/internal/tmetrics"
	ucase "cybertask/internal/usecase"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:generate mockery --filename=mock_controller.go --name=Controller --dir=. --structname=MockController --outpkg=mock_app
type Controller interface {
	Serve(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

// App represents whole aplication.
type App struct {
	ctl Controller
}

// Run is an app.New.Run() wrapper.
//
// Unhandled errors will be printed to the stderr, followed by exit(1).
func Run(cfg config.Config) {

	app, err := New(cfg)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	err = app.Run()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

}

// App constructor.
//
// Stages:
//   - Build logger.
//   - Connect to databse.
//   - Connect to all external services.
//   - Assemble http server.
func New(cfg config.Config) (*App, error) {
	l, err := logger.New(cfg.L)
	if err != nil {
		return nil, err
	}
	l.Info().Msg("logger construction succeded")
	cfg.Print(l)

	pg, err := pgsql.New(l, cfg.DB)
	if err != nil {
		return nil, err
	}

	taskMetrics, err := tmetrics.New(l, cfg.TM)

	taskUseCase := ucase.NewTaskUsecase(l,
		pg,
		taskMetrics,
	)

	taskhandler := handler.NewTaskHandler(l, taskUseCase)
	root := gin.New()
	routes := func() {
		tasks := func() {
			group := root.Group("task")

			group.GET(":id", taskhandler.Get)
			group.DELETE(":id", taskhandler.Delete)
			group.PUT("", taskhandler.Update)
			group.POST("", taskhandler.Create)
		}
		tasks()
		swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "TURNOFF_SWAGGER_HANDLER")

		root.GET("/swagger/*any", swaggerHandler)
	}
	routes()

	httpC := httpctl.New(l, cfg.C.HTTPServer, root)
	app := &App{
		ctl: httpC,
	}
	return app, nil

}
func (a *App) Run() error {
	a.ctl.Serve(context.TODO())

	err := a.ctl.Shutdown(context.TODO())
	if err != nil {
		return err
	}

	return nil

}
