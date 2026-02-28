package database

import (
	"log"

	"greens-co/backend/internal/config"
	"greens-co/backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var defaultAdminPassword = "Admin@greensco1"

func Connect(cfg *config.Config) *gorm.DB {
	// cfg is stored for seeding
	defaultAdminPassword = cfg.AdminDefaultPassword
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Enable pgcrypto for gen_random_uuid()
	db.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto;")

	// AutoMigrate all models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.ProductVariant{},
		&models.Bundle{},
		&models.BundleItem{},
		&models.Order{},
		&models.OrderItem{},
	); err != nil {
		log.Fatalf("failed to automigrate: %v", err)
	}

	Seed(db)
	return db
}

func Seed(db *gorm.DB) {
	// Only seed if tables are empty
	var categoryCount int64
	db.Model(&models.Category{}).Count(&categoryCount)
	if categoryCount > 0 {
		return
	}

	log.Println("Seeding database...")

	// Categories
	categories := []models.Category{
		{Name: "Salad", Slug: "salad"},
		{Name: "Rice Bowl", Slug: "rice-bowl"},
		{Name: "Drinks", Slug: "drinks"},
		{Name: "Snack", Slug: "snack"},
	}
	db.Create(&categories)

	// Map slug → ID
	catMap := map[string]string{}
	for _, c := range categories {
		catMap[c.Slug] = c.ID
	}

	ptr := func(s string) *string { return &s }
	ptrF := func(f float64) *float64 { return &f }
	ptrI := func(i int) *int { return &i }
	ptrI64 := func(i int64) *int64 { return &i }

	// Products
	products := []models.Product{
		{Name: "Garden Fresh Salad", Slug: "garden-fresh-salad", BasePrice: 35000, Badge: ptr("bestseller"), Rating: ptrF(4.9), ReviewCount: ptrI(128), Stock: 10, IsAvailable: true, Image: "https://images.unsplash.com/photo-1512621776951-a57141f2eefd?w=400&q=80", Description: "Mixed greens, cherry tomato, cucumber, house dressing", CategoryID: catMap["salad"]},
		{Name: "Teriyaki Chicken Bowl", Slug: "teriyaki-chicken-bowl", BasePrice: 45000, Badge: ptr("new"), Rating: ptrF(4.8), ReviewCount: ptrI(94), Stock: 8, IsAvailable: true, Image: "https://images.unsplash.com/photo-1547592180-85f173990554?w=400&q=80", Description: "Steamed rice, grilled chicken, teriyaki sauce, sesame", CategoryID: catMap["rice-bowl"]},
		{Name: "Green Detox Juice", Slug: "green-detox-juice", BasePrice: 28000, Rating: ptrF(4.7), ReviewCount: ptrI(76), Stock: 15, IsAvailable: true, Image: "https://images.unsplash.com/photo-1610970881699-44a5587cabec?w=400&q=80", Description: "Spinach, apple, ginger, lemon, cucumber blend", CategoryID: catMap["drinks"]},
		{Name: "Overnight Oats", Slug: "overnight-oats", BasePrice: 32000, Rating: ptrF(4.8), ReviewCount: ptrI(53), Stock: 12, IsAvailable: true, Image: "https://images.unsplash.com/photo-1563805042-7684c019e1cb?w=400&q=80", Description: "Rolled oats, chia seeds, almond milk, mixed berries", CategoryID: catMap["snack"]},
		{Name: "Quinoa Power Bowl", Slug: "quinoa-power-bowl", BasePrice: 48000, OriginalPrice: ptrI64(55000), Badge: ptr("promo"), Rating: ptrF(4.9), ReviewCount: ptrI(112), Stock: 6, IsAvailable: true, Image: "https://images.unsplash.com/photo-1540420773420-3366772f4999?w=400&q=80", Description: "Quinoa, roasted veggies, tahini dressing, seeds", CategoryID: catMap["salad"],
			Variants: []models.ProductVariant{
				{Name: "Regular", AdditionalPrice: 0},
				{Name: "Large (+Protein)", AdditionalPrice: 15000},
			},
		},
		{Name: "Açaí Bowl", Slug: "acai-bowl", BasePrice: 52000, Rating: ptrF(4.9), ReviewCount: ptrI(145), Stock: 0, IsAvailable: false, Image: "https://images.unsplash.com/photo-1490645935967-10de6ba17061?w=400&q=80", Description: "Blended açaí, granola, fresh fruits, honey drizzle", CategoryID: catMap["snack"]},
	}
	db.Create(&products)

	// Admin user
	hashed, _ := bcrypt.GenerateFromPassword([]byte(defaultAdminPassword), 12)
	admin := models.User{
		Name:     "Admin Greens",
		Email:    "admin@greensco.id",
		Password: string(hashed),
		Role:     "ADMIN",
	}
	db.Create(&admin)

	// Bundles
	bundles := []models.Bundle{
		{Name: "Healthy Starter", Slug: "healthy-starter", Price: 79000, OriginalPrice: 95000, IsPopular: false, Image: "https://images.unsplash.com/photo-1540420773420-3366772f4999?w=600&q=80", Description: "Perfect for a light & nutritious meal."},
		{Name: "Full Day Pack", Slug: "full-day-pack", Price: 125000, OriginalPrice: 140000, IsPopular: true, Image: "https://images.unsplash.com/photo-1498837167922-ddd27525d352?w=600&q=80", Description: "Complete nutrition for your entire day."},
		{Name: "Family Pack", Slug: "family-pack", Price: 215000, OriginalPrice: 256000, IsPopular: false, Image: "https://images.unsplash.com/photo-1565299624946-b28f40a0ae38?w=600&q=80", Description: "Feed the whole family with goodness."},
	}
	db.Create(&bundles)

	log.Println("Seeding complete.")
}
