package handlers

import "nctwo/backend/interfaces/middleware"

func (s *Server) InitializeRoutes() {
	s.Router.Use(middleware.CORSMiddleware()) //For CORS

	//user routes
	s.Router.POST("/users", s.Handler.SaveUser)
	s.Router.GET("/users", s.Handler.GetUsers)
	s.Router.GET("/users/:user_id", s.Handler.GetUser)

	//post routes
	s.Router.POST("/post", middleware.AuthMiddleware(), middleware.MaxSizeAllowed(8192000), s.Handler.SavePost)
	s.Router.PUT("/post/:post_id", middleware.AuthMiddleware(), middleware.MaxSizeAllowed(8192000), s.Handler.UpdatePost)
	s.Router.GET("/post/:post_id", s.Handler.GetPostAndCreator)
	s.Router.DELETE("/post/:post_id", middleware.AuthMiddleware(), s.Handler.DeletePost)
	s.Router.GET("/post", s.Handler.GetAllPost)

	//authentication routes
	s.Router.POST("/login", s.Handler.Login)
	s.Router.POST("/logout", s.Handler.Logout)
	s.Router.POST("/refresh", s.Handler.Refresh)
}
