// วุ่น
package dto

type LoginRequest struct {
	EmployeeID string `json:"employeeID"`
	Password   string `json:"password"`
	Role       string `json:"role"`
}

type LoginResponse struct {
	Success bool                   `json:"success"`
	Token   string                 `json:"token"`
	User    map[string]interface{} `json:"user"`
}

type PersonAddressRequest struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
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

type ParcelInfoRequest struct {
	Type   string  `json:"type"`
	Weight float64 `json:"weight"`
	Notes  string  `json:"notes,omitempty"`
}

type ParcelRequest struct {
	Deliver  PersonAddressRequest `json:"deliver"`
	Receiver PersonAddressRequest `json:"receiver"`
	Parcel   ParcelInfoRequest    `json:"parcel"`
}

type ShippingCalculationResponse struct {
	OriginZone      string  `json:"originZone"`
	DestinationZone string  `json:"destinationZone"`
	DeliveryType    string  `json:"deliveryType"`
	Weight          float64 `json:"weight"`
	ShippingCost    float64 `json:"shippingCost"`
}

type ParcelCreateResponse struct {
	Success      bool    `json:"success"`
	Message      string  `json:"message"`
	ParcelID     string  `json:"parcelID"`
	TrackID      string  `json:"trackID"`
	ShippingCost float64 `json:"shippingCost"`
}

type TrackingStatusUpdateRequest struct {
	Status      string `json:"status"`
	TrackID     string `json:"trackID"`
	Description string `json:"description,omitempty"`
	Location    string `json:"location,omitempty"`
}

type LegacyTrackActionRequest struct {
	TrackID string `json:"trackID"`
}

type VehicleAssignRequest struct {
	TrackID     string `json:"trackID"`
	VehicleID   int64  `json:"vehicleID"`
	MessengerID int64  `json:"messengerID"`
}

type VehicleSelectRequest struct {
	TrackID   string `json:"trackID"`
	VehicleID int64  `json:"vehicleID"`
}

type ReportCreateRequest struct {
	TrackID   string `json:"trackID"`
	Issue     string `json:"issue"`
	Reason    string `json:"reason"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type EmployeeProfileResponse struct {
	EmployeeID string `json:"EmployeeID"`
	Role       string `json:"Role"`
	Firstname  string `json:"Firstname"`
	Lastname   string `json:"Lastname"`
	Email      string `json:"Email"`
	Birthdate  string `json:"Birthdate,omitempty"`

	EmployeeCode string `json:"employeeCode"`
	DisplayRole  string `json:"displayRole"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	RoleName     string `json:"roleName"`
}
