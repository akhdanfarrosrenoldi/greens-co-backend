package models

import "time"

type User struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"default:'CUSTOMER'" json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Category struct {
	ID       string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name     string    `gorm:"not null" json:"name"`
	Slug     string    `gorm:"uniqueIndex;not null" json:"slug"`
	Image    *string   `json:"image"`
	Products []Product `gorm:"foreignKey:CategoryID" json:"-"`
}

type Product struct {
	ID            string           `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name          string           `gorm:"not null" json:"name"`
	Slug          string           `gorm:"uniqueIndex;not null" json:"slug"`
	Description   string           `json:"description"`
	BasePrice     int64            `gorm:"not null" json:"basePrice"`
	OriginalPrice *int64           `json:"originalPrice"`
	Image         string           `json:"image"`
	Stock         int              `gorm:"default:0" json:"stock"`
	IsAvailable   bool             `gorm:"default:true" json:"isAvailable"`
	Badge         *string          `json:"badge"`
	Rating        *float64         `json:"rating"`
	ReviewCount   *int             `json:"reviewCount"`
	CategoryID    string           `gorm:"type:uuid;not null" json:"categoryId"`
	Category      Category         `gorm:"foreignKey:CategoryID" json:"category"`
	Variants      []ProductVariant `gorm:"foreignKey:ProductID" json:"variants"`
	CreatedAt     time.Time        `json:"createdAt"`
	UpdatedAt     time.Time        `json:"updatedAt"`
}

type ProductVariant struct {
	ID              string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProductID       string `gorm:"type:uuid;not null" json:"productId"`
	Name            string `gorm:"not null" json:"name"`
	AdditionalPrice int64  `gorm:"default:0" json:"additionalPrice"`
}

type Bundle struct {
	ID            string       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name          string       `gorm:"not null" json:"name"`
	Slug          string       `gorm:"uniqueIndex;not null" json:"slug"`
	Description   string       `json:"description"`
	Price         int64        `gorm:"not null" json:"price"`
	OriginalPrice int64        `gorm:"not null" json:"originalPrice"`
	Image         string       `json:"image"`
	IsPopular     bool         `gorm:"default:false" json:"isPopular"`
	Items         []BundleItem `gorm:"foreignKey:BundleID" json:"items"`
	CreatedAt     time.Time    `json:"createdAt"`
	UpdatedAt     time.Time    `json:"updatedAt"`
}

type BundleItem struct {
	ID        string  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BundleID  string  `gorm:"type:uuid;not null" json:"bundleId"`
	ProductID string  `gorm:"type:uuid;not null" json:"productId"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
	Qty       int     `gorm:"not null" json:"qty"`
}

type Order struct {
	ID            string      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID        string      `gorm:"type:uuid;not null" json:"userId"`
	User          User        `gorm:"foreignKey:UserID" json:"user"`
	Status        string      `gorm:"default:'PENDING'" json:"status"`
	Type          string      `gorm:"not null" json:"type"`
	TotalPrice    int64       `gorm:"not null" json:"totalPrice"`
	Name          string      `json:"name"`
	Phone         string      `json:"phone"`
	Address       *string     `json:"address"`
	PickupTime    *string     `json:"pickupTime"`
	Notes         *string     `json:"notes"`
	PaymentStatus string      `gorm:"default:'UNPAID'" json:"paymentStatus"`
	MidtransID    *string     `json:"midtransId,omitempty"`
	Items         []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt"`
}

type OrderItem struct {
	ID        string          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrderID   string          `gorm:"type:uuid;not null" json:"orderId"`
	ProductID string          `gorm:"type:uuid;not null" json:"productId"`
	Product   Product         `gorm:"foreignKey:ProductID" json:"product"`
	VariantID *string         `gorm:"type:uuid" json:"variantId"`
	Variant   *ProductVariant `gorm:"foreignKey:VariantID" json:"variant"`
	Qty       int             `gorm:"not null" json:"qty"`
	Price     int64           `gorm:"not null" json:"price"`
	Notes     *string         `json:"notes"`
}
