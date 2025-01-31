package server

import services "github.com/EvansTrein/iqProgers/service"

func (s *HttpServer) InitRouters(wallet *services.Wallet) {
	
	s.router.GET("/ping", HandlerTest(s.log, wallet))
	
}
