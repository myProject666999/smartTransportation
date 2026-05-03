package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password  string         `json:"-" gorm:"size:255;not null"`
	Email     string         `json:"email" gorm:"size:100"`
	Phone     string         `json:"phone" gorm:"size:20"`
	RealName  string         `json:"real_name" gorm:"size:50"`
	IDCard    string         `json:"id_card" gorm:"size:20"`
	Avatar    string         `json:"avatar" gorm:"size:255"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if len(u.Password) > 0 && len(u.Password) < 60 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

type Admin struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password  string         `json:"-" gorm:"size:255;not null"`
	Role      string         `json:"role" gorm:"size:50;default:'admin'"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (a *Admin) BeforeSave(tx *gorm.DB) (err error) {
	if len(a.Password) > 0 && len(a.Password) < 60 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		a.Password = string(hashedPassword)
	}
	return nil
}

func (a *Admin) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	return err == nil
}

type Vehicle struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	PlateNumber  string         `json:"plate_number" gorm:"uniqueIndex;size:20;not null"`
	VehicleType  string         `json:"vehicle_type" gorm:"size:50"`
	OwnerName    string         `json:"owner_name" gorm:"size:50"`
	OwnerPhone   string         `json:"owner_phone" gorm:"size:20"`
	VehicleModel string         `json:"vehicle_model" gorm:"size:100"`
	Status       int            `json:"status" gorm:"default:1"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type VehicleMonitor struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	VehicleID   uint           `json:"vehicle_id"`
	PlateNumber string         `json:"plate_number" gorm:"size:20"`
	Location    string         `json:"location" gorm:"size:255"`
	Longitude   float64        `json:"longitude"`
	Latitude    float64        `json:"latitude"`
	Speed       float64        `json:"speed"`
	Direction   string         `json:"direction" gorm:"size:50"`
	Status      int            `json:"status" gorm:"default:1"`
	MonitorTime time.Time      `json:"monitor_time"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type RoadCondition struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	RoadName    string         `json:"road_name" gorm:"size:100"`
	Location    string         `json:"location" gorm:"size:255"`
	Condition   string         `json:"condition" gorm:"size:50"`
	Severity    string         `json:"severity" gorm:"size:20"`
	Description string         `json:"description" gorm:"size:500"`
	StartTime   time.Time      `json:"start_time"`
	EndTime     time.Time      `json:"end_time"`
	Status      int            `json:"status" gorm:"default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Violation struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserID       uint           `json:"user_id"`
	PlateNumber  string         `json:"plate_number" gorm:"size:20"`
	ViolationType string       `json:"violation_type" gorm:"size:100"`
	Location     string         `json:"location" gorm:"size:255"`
	ViolationTime time.Time     `json:"violation_time"`
	FineAmount   float64        `json:"fine_amount"`
	Points       int            `json:"points"`
	Status       int            `json:"status" gorm:"default:0"`
	PaymentTime  time.Time      `json:"payment_time"`
	Description  string         `json:"description" gorm:"size:500"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type Ticket struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	RouteName   string         `json:"route_name" gorm:"size:100"`
	StartStation string        `json:"start_station" gorm:"size:100"`
	EndStation  string         `json:"end_station" gorm:"size:100"`
	DepartureTime time.Time    `json:"departure_time"`
	ArrivalTime time.Time      `json:"arrival_time"`
	Price       float64        `json:"price"`
	TotalSeats  int            `json:"total_seats"`
	AvailableSeats int          `json:"available_seats"`
	Status      int            `json:"status" gorm:"default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Booking struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	OrderNo     string         `json:"order_no" gorm:"uniqueIndex;size:50"`
	UserID      uint           `json:"user_id"`
	TicketID    uint           `json:"ticket_id"`
	PassengerName string       `json:"passenger_name" gorm:"size:50"`
	PassengerIDCard string     `json:"passenger_id_card" gorm:"size:20"`
	PassengerPhone string      `json:"passenger_phone" gorm:"size:20"`
	Quantity    int            `json:"quantity"`
	TotalAmount float64        `json:"total_amount"`
	Status      int            `json:"status" gorm:"default:0"`
	PaymentTime time.Time      `json:"payment_time"`
	SeatNumbers string         `json:"seat_numbers" gorm:"size:100"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type News struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"size:200"`
	Content   string         `json:"content" gorm:"type:text"`
	Summary   string         `json:"summary" gorm:"size:500"`
	Image     string         `json:"image" gorm:"size:255"`
	Category  string         `json:"category" gorm:"size:50"`
	Author    string         `json:"author" gorm:"size:50"`
	ViewCount int            `json:"view_count" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Announcement struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"size:200"`
	Content   string         `json:"content" gorm:"type:text"`
	PublishTime time.Time    `json:"publish_time"`
	ExpireTime time.Time     `json:"expire_time"`
	Priority  int            `json:"priority" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Carousel struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"size:100"`
	Image     string         `json:"image" gorm:"size:255"`
	Link      string         `json:"link" gorm:"size:255"`
	SortOrder int            `json:"sort_order" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Feedback struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id"`
	Title       string         `json:"title" gorm:"size:200"`
	Content     string         `json:"content" gorm:"type:text"`
	FeedbackType string       `json:"feedback_type" gorm:"size:50"`
	Status      int            `json:"status" gorm:"default:0"`
	Reply       string         `json:"reply" gorm:"type:text"`
	ReplyTime   time.Time      `json:"reply_time"`
	Images      string         `json:"images" gorm:"size:500"`
	Contact     string         `json:"contact" gorm:"size:100"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type TrafficFlow struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	RoadName    string         `json:"road_name" gorm:"size:100"`
	Location    string         `json:"location" gorm:"size:255"`
	VehicleCount int           `json:"vehicle_count"`
	AverageSpeed float64       `json:"average_speed"`
	FlowLevel   string         `json:"flow_level" gorm:"size:20"`
	RecordTime  time.Time      `json:"record_time"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type SignalLight struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Location    string         `json:"location" gorm:"size:255"`
	Intersection string       `json:"intersection" gorm:"size:100"`
	Status      string         `json:"status" gorm:"size:20"`
	GreenDuration int          `json:"green_duration"`
	RedDuration int            `json:"red_duration"`
	YellowDuration int         `json:"yellow_duration"`
	LastUpdated time.Time      `json:"last_updated"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Route struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	RouteName   string         `json:"route_name" gorm:"size:100"`
	RouteNo     string         `json:"route_no" gorm:"size:20"`
	StartStation string        `json:"start_station" gorm:"size:100"`
	EndStation  string         `json:"end_station" gorm:"size:100"`
	Stations    string         `json:"stations" gorm:"type:text"`
	Distance    float64        `json:"distance"`
	Duration    int            `json:"duration"`
	Price       float64        `json:"price"`
	Status      int            `json:"status" gorm:"default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
