package database

import (
	"log"

	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/repository"
	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin(repo *repository.UserRepo, email, password string) {
	if _, err := repo.FindByEmail(email); err == nil {
		return // admin already exists
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if _, err := repo.Create(email, string(hashed), "admin"); err != nil {
		log.Printf("[SEED] Could not create admin: %v", err)
		return
	}
	log.Println("[SEED] Admin user created:", email)
}
