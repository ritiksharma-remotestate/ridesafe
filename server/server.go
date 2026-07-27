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

	//  routes for driver protected

	mux.Handle("POST /drivers", middleware.AuthMiddleware(middleware.ShouldHaveRole(http.HandlerFunc(handlers.CreateDriver), models.RoleDriver)))
	mux.Handle("GET /drivers/me", middleware.AuthMiddleware(middleware.ShouldHaveRole(http.HandlerFunc(handlers.GetMyDriverProfile), models.RoleDriver)))
	mux.Handle("PATCH /drivers/location", middleware.AuthMiddleware(middleware.ShouldHaveRole(http.HandlerFunc(handlers.UpdateDriverLocation), models.RoleDriver)))
	mux.Handle("PATCH /drivers/online", middleware.AuthMiddleware(middleware.ShouldHaveRole(http.HandlerFunc(handlers.SetDriverOnline), models.RoleDriver)))
	mux.Handle("PATCH /drivers/available", middleware.AuthMiddleware(middleware.ShouldHaveRole(http.HandlerFunc(handlers.SetDriverAvailable), models.RoleDriver)))
	mux.Handle("PATCH /rides/{id}/accept", middleware.AuthMiddleware(middleware.ShouldHaveRole(http.HandlerFunc(handlers.AcceptRide), models.RoleDriver)))
	mux.Handle("PATCH /rides/{id}/arrive", middleware.AuthMiddleware(middleware.ShouldHaveRole(http.HandlerFunc(handlers.ArriveRide), models.RoleDriver)))
	mux.Handle("PATCH /rides/{id}/start", middleware.AuthMiddleware(middleware.ShouldHaveRole(http.HandlerFunc(handlers.StartRide), models.RoleDriver)))
	mux.Handle("PATCH /rides/{id}/complete", middleware.AuthMiddleware(middleware.ShouldHaveRole(http.HandlerFunc(handlers.CompleteRide), models.RoleDriver)))
	mux.Handle("PATCH /rides/{id}/verify-otp", middleware.AuthMiddleware(middleware.ShouldHaveRole(http.HandlerFunc(handlers.VerifyRideOTP), models.RoleDriver)))
	//   routes for passengers
	mux.Handle("GET /rides/{id}/drivers-available", middleware.AuthMiddleware(middleware.ShouldHaveRole(http.HandlerFunc(handlers.GetAvailableDrivers), models.RolePassenger)))
	mux.Handle("POST /rides", middleware.AuthMiddleware(middleware.ShouldHaveRole(http.HandlerFunc(handlers.CreateRide), models.RolePassenger)))
	mux.Handle("GET /location", middleware.AuthMiddleware(middleware.ShouldHaveRole(http.HandlerFunc(handlers.GetDriverLocation), models.RolePassenger)))
	// routes for both passenger and driver
	mux.Handle("GET /rides/{id}", middleware.AuthMiddleware(middleware.ShouldHaveRole(http.HandlerFunc(handlers.GetRideByID), models.RolePassenger, models.RoleDriver)))
	mux.Handle("PATCH /rides/{id}/cancel", middleware.AuthMiddleware(middleware.ShouldHaveRole(http.HandlerFunc(handlers.CancelRide), models.RolePassenger, models.RoleDriver)))

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
