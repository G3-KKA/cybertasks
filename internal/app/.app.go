package app

import (
	"cybertask/config"
	"cybertask/internal/logger"
	"os"

	"github.com/prometheus/common/route"
)

// App represents whole aplication.
type App struct {
	// depMgr depmngr.Manager
}

//go:generate mockery --filename=mock_.go --name=Controller --dir=. --structname=MockController --outpkg=mock_
type Controller interface {
}

type TrueCustomError interface {
	Unwrap() error
	Error() string
	Is(err, target error) bool
	As(target any) bool
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
func New(cfg config.Config) (*App, error) {
	l, err := logger.New(cfg.Logger)
	if err != nil {
		return nil, err
	}
	l.Info().Msg("logger construction succeded")
	cfg.Print(l)

	// Сделаю в самом конце, сейчас не принципиально.
	depMgr := depmgr.New(cfg.DepManager)
	defer func() {
		if err != nil {
			err := depMgr.ShutdownAll()
			if err != nil {
				l.Error().AnErr("depmgrerror", err).Send()

			}
		}
	}()
	// pgpool.Conn => bun.Upgrade => wrapper ? pgx or pgpool ?
	// Shutdown-able via RegisterAsDep or something
	pg, err := pgsql.New(l, cfg.Postgres)
	if err != nil {
		return nil, err
	}
	err = depMgr.NewLevel().Register(pg)
	if err != nil {
		return nil, err
	}

	// There should be an error ?
	// I think not
	taskMetrics, err := tmetrics.New(l, cfg.TaskMetrics)
	if err != nil {
		return nil, err
	}
	depMgr.Register(taskMetrics)

	taskUseCase, err := ucase.NewTaskUsecase(l,
		repo.NewTaskRepo(l, pg), // Only model.Task related ops
		taskMetrics,
	) // model.Task related business logic
	taskhandler := handler.NewTaskHandler(l, taskUseCase)
	router := route.New(taskhandler, swaggerhandler)

	httpC, err := httpctl.New(l, cfg.C.HTTPServer, router) //hs ... httpctl.handlers aka attachable
	if err != nil {
		return nil, err
	}
	depMgr.NewLevel().Register(httpC)

	//db, err pgx.Open
	//bunpg := bun.Db
	//tasksUsecase = usecase.New(
	//	repo.New(cfg.Repository,l,bunpg),
	//	tasksrvc.New(cfg.TaskService,l)
	//)
	// ============== ПИШЕМ СНИЗУ ВВЕРХ ! =============

}
func (a *App) Run() error {

	err = ctl.Shutdown()
	if err != nil {
		return nil, err
	}

}
