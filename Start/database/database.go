package database

import (
	"fmt"
	"log"
	"os"

	"goexperiments/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func ConnectDb() {
	dsn := fmt.Sprintf(
		"host=db user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}

	log.Println("connected")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("running migrations")
	db.AutoMigrate(&models.Fact{})

	DB = Dbinstance{
		Db: db,
	}
}

// func ConnectDB2() {
// 	var err error
// 	p := config.Config("DB_PORT")
// 	port, err := strconv.ParseUint(p, 10, 32)
// 	if err != nil {
// 		panic("failed to parse database port")
// 	}

// 	dsn := fmt.Sprintf(
// 		"host=db port=%d user=%s password=%s dbname=%s sslmode=disable",
// 		port,
// 		config.Config("DB_USER"),
// 		config.Config("DB_PASSWORD"),
// 		config.Config("DB_NAME"),
// 	)
// 	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

// 	if err != nil {
// 		panic("failed to connect database")
// 	}

// 	fmt.Println("Connection Opened to Database")
// 	DB.AutoMigrate(&model.Product{}, &model.User{})
// 	fmt.Println("Database Migrated")
// }
