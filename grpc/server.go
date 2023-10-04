package grpc

import (
	"github.com/dylanpeng/golib/logger"
	"google.golang.org/grpc"
	"net"
	"strconv"
)

type Config struct {
	Host string `toml:"host" json:"host" yaml:"host"`
	Port int    `toml:"port" json:"port" yaml:"port"`
}

func (c *Config) GetAddr() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

type Router interface {
	RegGrpcService(server *grpc.Server)
}

type Server struct {
	cgf    *Config
	logger *logger.Logger
	router Router
	opts   []grpc.ServerOption
	server *grpc.Server
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.cgf.GetAddr())

	if err != nil {
		s.logger.Errorf("start grpc fail. | err: %s", err)
		return err
	}

	s.server = grpc.NewServer(s.opts...)
	s.router.RegGrpcService(s.server)

	go func() {
		e := s.server.Serve(lis)

		if e != nil {
			s.logger.Errorf("grpc server start fail. | err: %s", e)
		}
	}()

	return nil
}

func (s *Server) AddOpt(opt grpc.ServerOption) {
	if opt == nil {
		return
	}

	if s.opts == nil {
		s.opts = make([]grpc.ServerOption, 8)
	}

	s.opts = append(s.opts, opt)
}

func (s *Server) Close() {
	s.server.Stop()
}

func NewServer(cgf *Config, router Router, logger *logger.Logger, opts ...grpc.ServerOption) *Server {
	return &Server{
		cgf:    cgf,
		logger: logger,
		opts:   opts,
		router: router,
	}
}
