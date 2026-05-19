// ยีน
package domain

import "time"

type Employee struct {
	ID           int64      `json:"id"`
	EmployeeCode string     `json:"employeeCode"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	FirstName    string     `json:"firstName"`
	LastName     string     `json:"lastName"`
	Phone        string     `json:"phone,omitempty"`
	Role         string     `json:"role"`
	BirthDate    *time.Time `json:"birthDate,omitempty"`
	IsActive     bool       `json:"isActive"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

type User struct {
	ID        int64        `json:"id"`
	FirstName string       `json:"firstName"`
	LastName  string       `json:"lastName"`
	Phone     string       `json:"phone"`
	Email     string       `json:"email,omitempty"`
	Address   *UserAddress `json:"address,omitempty"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}

type UserAddress struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"userId"`
	HomeNumber  string    `json:"homeNumber"`
	Soi         string    `json:"soi,omitempty"`
	Road        string    `json:"road,omitempty"`
	Subdistrict string    `json:"subdistrict"`
	District    string    `json:"district"`
	Province    string    `json:"province"`
	Zipcode     string    `json:"zipcode"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Vehicle struct {
	ID                 int64     `json:"id"`
	VehicleCode        string    `json:"vehicleCode"`
	Type               string    `json:"type"`
	LicensePlate       string    `json:"licensePlate"`
	Status             string    `json:"status"`
	AssignedEmployeeID *int64    `json:"assignedEmployeeId,omitempty"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

type ShippingRate struct {
	ID           int64   `json:"id"`
	ZoneCode     string  `json:"zoneCode"`
	DeliveryType string  `json:"deliveryType"`
	WeightMin    float64 `json:"weightMin"`
	WeightMax    float64 `json:"weightMax"`
	BasePrice    float64 `json:"basePrice"`
	ExtraPerKg   float64 `json:"extraPerKg"`
}

type Parcel struct {
	ID                  int64      `json:"id"`
	ParcelCode          string     `json:"parcelCode"`
	TrackingCode        string     `json:"trackingCode"`
	SenderUserID        int64      `json:"senderUserId"`
	SenderAddressID     int64      `json:"senderAddressId"`
	ReceiverUserID      int64      `json:"receiverUserId"`
	ReceiverAddressID   int64      `json:"receiverAddressId"`
	ClerkEmployeeID     *int64     `json:"clerkEmployeeId,omitempty"`
	MessengerEmployeeID *int64     `json:"messengerEmployeeId,omitempty"`
	VehicleID           *int64     `json:"vehicleId,omitempty"`
	DeliveryType        string     `json:"deliveryType"`
	Weight              float64    `json:"weight"`
	ShippingCost        float64    `json:"shippingCost"`
	OriginZone          string     `json:"originZone"`
	DestinationZone     string     `json:"destinationZone"`
	Status              string     `json:"status"`
	Notes               string     `json:"notes,omitempty"`
	DepositedAt         time.Time  `json:"depositedAt"`
	DeliveredAt         *time.Time `json:"deliveredAt,omitempty"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
}

type TrackingEvent struct {
	ID                  int64     `json:"id"`
	ParcelID            int64     `json:"parcelId"`
	TrackingCode        string    `json:"trackingCode"`
	Status              string    `json:"status"`
	Location            string    `json:"location,omitempty"`
	Description         string    `json:"description,omitempty"`
	UpdatedByEmployeeID *int64    `json:"updatedByEmployeeId,omitempty"`
	CreatedAt           time.Time `json:"createdAt"`
}

type ParcelListItem struct {
	ParcelID            string    `json:"ParcelID"`
	TrackID             string    `json:"TrackID"`
	ReceiverName        string    `json:"ReceiverName"`
	ReceiverSurname     string    `json:"ReceiverSurname"`
	ReceiverTel         string    `json:"ReceiverTel"`
	HomeNumber          string    `json:"HomeNumber"`
	Soi                 string    `json:"Soi"`
	Road                string    `json:"Road"`
	Subdistrict         string    `json:"Subdistrict"`
	DistrictName        string    `json:"DistrictName"`
	ProvinceName        string    `json:"ProvinceName"`
	Zipcode             string    `json:"Zipcode"`
	Status              string    `json:"Status"`
	DepositDate         time.Time `json:"DepositDate"`
	AssignedMessengerID *int64    `json:"AssignedMessengerID,omitempty"`
}

type ParcelDetail struct {
	Parcel          Parcel          `json:"parcel"`
	Sender          User            `json:"sender"`
	SenderAddress   UserAddress     `json:"senderAddress"`
	Receiver        User            `json:"receiver"`
	ReceiverAddress UserAddress     `json:"receiverAddress"`
	Events          []TrackingEvent `json:"events"`
}

type ParcelTrackingView struct {
	TrackID     string          `json:"TrackID"`
	ParcelID    string          `json:"ParcelID"`
	Status      string          `json:"Status"`
	Sender      string          `json:"Sender"`
	Receiver    string          `json:"Receiver"`
	Address     string          `json:"Address"`
	Deliverer   string          `json:"Deliverer"`
	TypeCar     string          `json:"Typecar"`
	License     string          `json:"License"`
	DeliveredAt *time.Time      `json:"DeliveredAt,omitempty"`
	UpdatedAt   time.Time       `json:"UpdatedAt"`
	Events      []TrackingEvent `json:"Events"`
}

type PersonAddress struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Phone       string `json:"phone"`
	Email       string `json:"email,omitempty"`
	HomeNumber  string `json:"homeNumber"`
	Soi         string `json:"soi,omitempty"`
	Road        string `json:"road,omitempty"`
	District    string `json:"district"`
	Subdistrict string `json:"subdistrict"`
	Province    string `json:"province"`
	Zipcode     string `json:"zipcode"`
}

type ParcelCreateInput struct {
	Sender          PersonAddress `json:"sender"`
	Receiver        PersonAddress `json:"receiver"`
	ParcelCode      string        `json:"parcelCode"`
	TrackingCode    string        `json:"trackingCode"`
	ClerkEmployeeID *int64        `json:"clerkEmployeeId,omitempty"`
	DeliveryType    string        `json:"deliveryType"`
	Weight          float64       `json:"weight"`
	ShippingCost    float64       `json:"shippingCost"`
	OriginZone      string        `json:"originZone"`
	DestinationZone string        `json:"destinationZone"`
	Status          string        `json:"status"`
	Notes           string        `json:"notes,omitempty"`
	DepositedAt     time.Time     `json:"depositedAt"`
	InitialEvent    TrackingEvent `json:"initialEvent"`
}

type ParcelCreated struct {
	ParcelID     int64   `json:"parcelId"`
	ParcelCode   string  `json:"parcelCode"`
	TrackingCode string  `json:"trackingCode"`
	ShippingCost float64 `json:"shippingCost"`
}

type TrackingUpdateInput struct {
	TrackID             string     `json:"trackId"`
	Status              string     `json:"status"`
	Location            string     `json:"location,omitempty"`
	Description         string     `json:"description,omitempty"`
	UpdatedByEmployeeID *int64     `json:"updatedByEmployeeId,omitempty"`
	DeliveredAt         *time.Time `json:"deliveredAt,omitempty"`
}

type VehicleAssignmentInput struct {
	TrackID             string `json:"trackId"`
	VehicleID           int64  `json:"vehicleId"`
	MessengerEmployeeID int64  `json:"messengerEmployeeId"`
}

type Report struct {
	ID                int64     `json:"id"`
	ReportCode        string    `json:"reportCode"`
	ParcelID          *int64    `json:"parcelId,omitempty"`
	TrackingCode      string    `json:"trackingCode"`
	ReporterFirstName string    `json:"reporterFirstName"`
	ReporterLastName  string    `json:"reporterLastName"`
	ReporterPhone     string    `json:"reporterPhone,omitempty"`
	ReporterEmail     string    `json:"reporterEmail,omitempty"`
	IssueType         string    `json:"issueType"`
	Subject           string    `json:"subject"`
	Description       string    `json:"description"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
