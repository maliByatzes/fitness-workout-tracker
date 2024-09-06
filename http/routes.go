package http

func (s *Server) routes() {
	s.router.Use(CORSMiddleware())

	apiRouter := s.router.Group("/api/v1")
	{
		apiRouter.GET("/healthchecker", healthCheck())
		apiRouter.POST("/users/register", s.createUser())
		apiRouter.POST("/users/login", s.loginUser())
		apiRouter.POST("/users/logout", s.logoutUser())

		apiRouter.Use(s.authenticate())
		{
			apiRouter.GET("/users/me", s.getCurrentUser())
		}
	}
}
