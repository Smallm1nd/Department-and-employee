package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Smallm1nd/Department-and-employee/internal/handlers"
	"github.com/Smallm1nd/Department-and-employee/internal/repository"
	"github.com/Smallm1nd/Department-and-employee/internal/routes"
	"github.com/Smallm1nd/Department-and-employee/internal/service"
	"github.com/pressly/goose/v3"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := getDSN()

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := conn.DB()
	if err != nil {
		log.Fatal("failed to get underlying sql.DB: ", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal("failed to set dialect: ", err)
	}

	if err := goose.Up(sqlDB, "migrations"); err != nil {
		log.Fatal("failed to run migrations: ", err)
	}

	deptRepo := repository.NewDepartmentRepo(conn)
	empRepo := repository.NewEmployeeRepo(conn)

	deptService := service.NewDepartmentService(deptRepo)
	empService := service.NewEmployeeService(empRepo, deptRepo)

	deptHandler := handlers.NewDepartmentHandler(deptService)
	empHandler := handlers.NewEmployeeHandler(empService)

	mux := routes.NewRouter(deptHandler, empHandler)

	log.Println("Сервер запущен на порту :6969")
	err = http.ListenAndServe(":6969", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func getDSN() string {
	dbServer := flag.String("dbserver", "localhost", "Database server address")
	dbName := flag.String("dbname", "test_department", "Database name")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbUser := flag.String("dbuser", "postgres", "User for database")
	dbPass := flag.String("dbpass", "", "Password for database")

	flag.Parse()

	if *dbPass == "" {
		log.Fatal("Error: Password for database is required")
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", *dbServer, *dbUser, *dbPass, *dbName, *dbPort)
}
