package http

func (s *Server) routes() {
	s.router.Use(CORSMiddleware())

	apiRouter := s.router.Group("/api/v1")
	{
		apiRouter.GET("/healthchecker", healthCheck())
	}
}
