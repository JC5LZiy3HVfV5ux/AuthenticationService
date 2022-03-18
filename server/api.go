package server

func InitApiV1() {
	apiRouter := Srv.Router.PathPrefix("/api/v1").Subrouter()

	InitAuth(apiRouter)
}
