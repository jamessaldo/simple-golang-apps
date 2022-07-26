package handlers

import "nctwo/backend/interfaces/middleware"

func (s *Server) InitializeRoutes() {
	s.Router.Use(middleware.CORSMiddleware()) //For CORS

	//user routes
	s.Router.POST("/users", s.Handler.SaveUser)
	s.Router.GET("/users", s.Handler.GetUsers)
	s.Router.GET("/users/:user_id", s.Handler.GetUser)

	//post routes
	s.Router.POST("/post", s.Handler.SavePost)
	s.Router.PUT("/post/:post_id", s.Handler.UpdatePost)
	s.Router.GET("/post/:post_id", s.Handler.GetPostAndCreator)
	s.Router.DELETE("/post/:post_id", s.Handler.DeletePost)
	s.Router.GET("/post", s.Handler.GetAllPost)

	//comment routes
	s.Router.POST("/comment", s.Handler.SaveComment)
	s.Router.PUT("/comment/:comment_id", s.Handler.UpdateComment)
	s.Router.GET("/comment/:comment_id", s.Handler.GetCommentAndCreator)
	s.Router.DELETE("/comment/:comment_id", s.Handler.DeleteComment)
	s.Router.GET("/comment", s.Handler.GetAllComment)

	//authentication routes
	s.Router.POST("/login", s.Handler.Login)
	s.Router.POST("/logout", s.Handler.Logout)
	s.Router.POST("/refresh", s.Handler.Refresh)
}
