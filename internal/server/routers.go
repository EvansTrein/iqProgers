package server

import services "github.com/EvansTrein/iqProgers/service"

func (s *HttpServer) InitRouters(wallet *services.Wallet) {
	
	s.router.POST("/deposit", Deposit(s.log, wallet))
	s.router.GET("/operations")
	
}
