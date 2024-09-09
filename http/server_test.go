package http_test

import (
	"testing"

	fwthttp "github.com/maliByatzes/fwt/http"
	"github.com/maliByatzes/fwt/mock"
	"github.com/maliByatzes/fwt/postgres"
	"github.com/stretchr/testify/require"
)

const (
	TestPort      = "8080"
	TestSecretKey = "00000000000000000000000000000000"
)

type Server struct {
	*fwthttp.Server
	UserService mock.UserService
}

func MustOpenServer(tb testing.TB) *Server {
	tb.Helper()

	srv, err := fwthttp.NewServer(&postgres.DB{}, TestSecretKey)
	require.NoError(tb, err)
	s := &Server{Server: srv}

	s.Server.UserService = &s.UserService

	err = s.Run(TestPort)
	require.NoError(tb, err)

	return s
}

func MustCloseServer(tb testing.TB, s *Server) {
	tb.Helper()
	err := s.Close()
	require.NoError(tb, err)
}
