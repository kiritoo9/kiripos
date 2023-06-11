package seeders

import (
	"fmt"

	"kiritech/src/configs"
	"kiritech/src/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func UserSeed() {
	configs.Connect()

	password := []byte("admin")
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	password = hash

	users := []*models.Users{
		&models.Users{Id: uuid.New(), Name: "Administrator", Email: "admin@admin.com", Password: string(password)},
	}

	result := configs.DB.Create(users)
	if result.Error != nil {
		fmt.Println("User seed failed ", result.Error)
	} else {
		fmt.Println("User seed successfully created!")
	}
}
