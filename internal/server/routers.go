package server

import services "github.com/EvansTrein/iqProgers/service"

func (s *HttpServer) InitRouters(wallet *services.Wallet) {
	authRouters := s.router.Group("/test")

	authRouters.GET("/ping", HandlerTest(s.log, wallet))
}
