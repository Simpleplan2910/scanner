package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"scanner/pkg/db"
	"scanner/pkg/repos"
	"scanner/pkg/services/git"
	"scanner/pkg/services/scanner"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Server struct {
	ListenAddr    *net.TCPAddr
	RootRouter    *mux.Router
	Router        *mux.Router
	Logger        *logrus.Entry
	ListenAddress string
	Server        *http.Server
	Store         *db.Store

	ReposService repos.Service

	GitService     git.Service
	ScannerService scanner.Service
}

// Option : server option
type Option func(s *Server) error

// NewServer : new server
func NewServer(options ...Option) (*Server, error) {
	s := &Server{}
	for _, option := range options {
		if err := option(s); err != nil {
			return nil, errors.Wrap(err, "failed to apply option")
		}
	}
	// setup Router for mux
	s.Router = mux.NewRouter()
	s.RootRouter = s.Router.PathPrefix("/v1").Subrouter()
	s.initService()
	return s, nil
}

// Start : start server
func (s *Server) Start() error {
	s.Logger.Infof("Starting Server at %s....", s.ListenAddress)
	recoverhandler := handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(s.Router)
	loghandler := handlers.LoggingHandler(os.Stdout, recoverhandler)
	s.Server = &http.Server{
		Handler:      loghandler,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 30,
	}
	listener, err := net.Listen("tcp", s.ListenAddress)
	if err != nil {
		return fmt.Errorf("start server error %w", err)
	}
	s.ListenAddr = listener.Addr().(*net.TCPAddr)
	go s.Server.Serve(listener)

	return nil
}

// TODO: read options from file
func ServerOpts(store *db.Store) Option {
	return func(s *Server) error {
		logger := logrus.New()
		gService := git.New("samples")
		sService := scanner.NewLocalService(s.Store.Result, gService, s.Store.Scan, 8)
		s.ReposService = repos.NewService(gService, s.Store.Repos, s.Store.Result, sService, s.Store.Scan)
		s.GitService = gService
		s.ScannerService = sService
		s.Logger = logger.WithField("service", "scanner")
		return nil
	}
}

// Shutdown : shutdown server
func (s *Server) Shutdown(ctx context.Context) {
	s.Logger.Infoln("Stopping Server...")
	if s.Server != nil {
		if err := s.Server.Shutdown(ctx); err != nil {
			s.Logger.Warn("Unable to shutdown server", err)
		}
		s.Server.Close()
		s.Server = nil
	}
}

func (s *Server) initService() {
	s.intiHealthCheck()
	s.initReposService()
}

func (s *Server) intiHealthCheck() {
	s.Router.HandleFunc("/healthCheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
}

func (s *Server) initReposService() {
	loggingService := repos.NewLoggingService(s.Logger, s.ReposService)
	repos.NewHandler(s.Logger, loggingService).AddRoutes(s.RootRouter)
}
