package server

import (
	"net/http"

	"context"
	"ridesafe/handlers"
	"ridesafe/middleware"
	"ridesafe/models"
	"time"
)
type Server struct {
    Router *http.ServeMux
    server *http.Server
}

const (
	readTimeout       = 5 * time.Minute
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 5 * time.Minute
)


func SetupRoutes() *Server {
	mux := http.NewServeMux()

	

	// Public Routes
	mux.HandleFunc("POST /register", handlers.Register)
	mux.HandleFunc("POST /login", handlers.Login)

	// Protected Routes
	mux.Handle("POST /drivers", middleware.AuthMiddleware( middleware.ShouldHaveRole(models.RoleDriver)(http.HandlerFunc(handlers.CreateDriver))))
	mux.Handle("GET /drivers/me", middleware.AuthMiddleware( middleware.ShouldHaveRole(models.RolePassenger)(http.HandlerFunc(handlers.GetMyDriverProfile))))
	mux.Handle("PATCH /drivers/location", middleware.AuthMiddleware( middleware.ShouldHaveRole(models.RoleDriver)(http.HandlerFunc(handlers.UpdateDriverLocation))))
	mux.Handle("PATCH /drivers/online", middleware.AuthMiddleware( middleware.ShouldHaveRole(models.RoleDriver)(http.HandlerFunc(handlers.SetDriverOnline))))
	mux.Handle("PATCH /drivers/available", middleware.AuthMiddleware( middleware.ShouldHaveRole(models.RoleDriver)(http.HandlerFunc(handlers.SetDriverAvailable))))
	mux.Handle("GET /drivers/available", middleware.AuthMiddleware( middleware.ShouldHaveRole(models.RolePassenger)(http.HandlerFunc(handlers.GetAvailableDrivers))))

	
	mux.Handle("POST /rides", middleware.AuthMiddleware( middleware.ShouldHaveRole(models.RolePassenger)(http.HandlerFunc(handlers.CreateRide))))
	mux.Handle("GET /rides/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetRideByID)))
	mux.Handle("PATCH /rides/{id}/accept", middleware.AuthMiddleware(middleware.ShouldHaveRole(models.RoleDriver)(http.HandlerFunc(handlers.AcceptRide))))
	mux.Handle("PATCH /rides/{id}/arrive", middleware.AuthMiddleware(middleware.ShouldHaveRole(models.RoleDriver)(http.HandlerFunc(handlers.ArriveRide))))
	mux.Handle("PATCH /rides/{id}/start", middleware.AuthMiddleware(middleware.ShouldHaveRole(models.RoleDriver)(http.HandlerFunc(handlers.StartRide))))
	mux.Handle("PATCH /rides/{id}/complete", middleware.AuthMiddleware(middleware.ShouldHaveRole(models.RoleDriver)(http.HandlerFunc(handlers.CompleteRide))))
	mux.Handle("PATCH /rides/{id}/cancel", middleware.AuthMiddleware(http.HandlerFunc(handlers.CancelRide)))

	return &Server{
		Router: mux,
	}
}
func (svc *Server) Run(port string) error {
	svc.server = &http.Server{
		Addr:              port,
		Handler:           svc.Router,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
	}
	return svc.server.ListenAndServe()
}
func (svc *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return svc.server.Shutdown(ctx)
}