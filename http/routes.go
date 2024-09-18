package http

func (s *Server) routes() {
	s.Router.Use(CORSMiddleware())

	apiRouter := s.Router.Group("/api/v1")
	{
		apiRouter.GET("/healthchecker", healthCheck())
		apiRouter.POST("/users/register", s.createUser())
		apiRouter.POST("/users/login", s.loginUser())
		apiRouter.POST("/users/logout", s.logoutUser())

		apiRouter.Use(s.authenticate())
		{
			apiRouter.GET("/users/me", s.getCurrentUser())
			apiRouter.PATCH("/users/update", s.updateUser())
			apiRouter.DELETE("/users/delete", s.deleteUser())

			apiRouter.POST("/profile/create", s.createProfile())
			apiRouter.GET("/profile", s.getUserProfile())
			apiRouter.PATCH("/profile/update", s.updateProfile())
			apiRouter.DELETE("/profile/delete", s.deleteProfile())

			apiRouter.POST("/workout/create", s.createWorkout())
			apiRouter.GET("/workout/all", s.getAllWorkouts())
			apiRouter.GET("/workout/:id", s.getOneWorkout())
			apiRouter.PATCH("/workout/:id", s.updateWorkout())
			apiRouter.DELETE("/workout/exercises/:id", s.removeExercisesFromWorkout())
			apiRouter.DELETE("/workout/:id", s.deleteWorkout())
		}
	}
}
