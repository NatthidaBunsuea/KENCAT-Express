// โอม
package util

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	RoleAdmin       = "ADMIN"
	RoleParcelClerk = "PARCEL_CLERK"
	RoleMessenger   = "MESSENGER"

	StatusPending        = "Pending"
	StatusInTransit      = "In Transit"
	StatusDelivered      = "Delivered"
	StatusDeliveryFailed = "Delivery Failed"

	DeliveryStandard = "STANDARD"
	DeliveryExpress  = "EXPRESS"

	ZoneBangkokMetro = "BANGKOK_METRO"
	ZoneUpcountry    = "UPCOUNTRY"

	CookieName = "kencat_token"
)

var phoneRegex = regexp.MustCompile(`^\d{9,10}$`)

type Claims struct {
	Subject      string `json:"sub"`
	EmployeeID   int64  `json:"employee_id"`
	EmployeeCode string `json:"employee_code"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Name         string `json:"name"`
	IssuedAt     int64  `json:"iat"`
	ExpiresAt    int64  `json:"exp"`
}

func DecodeJSON(r *http.Request, dst interface{}) error {
	defer r.Body.Close()
	body, err := io.ReadAll(io.LimitReader(r.Body, 2<<20))
	if err != nil {
		return err
	}
	if len(strings.TrimSpace(string(body))) == 0 {
		return errors.New("request body is required")
	}
	return json.Unmarshal(body, dst)
}

func WriteJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func ErrorJSON(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]interface{}{
		"success": false,
		"message": message,
		"error":   message,
	})
}

func SuccessJSON(w http.ResponseWriter, status int, data map[string]interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	data["success"] = true
	WriteJSON(w, status, data)
}

func NormalizeWhitespace(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func NormalizeRole(s string) string {
	s = strings.TrimSpace(strings.ToUpper(s))
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.Join(strings.Fields(s), " ")
	switch s {
	case "P", "PARCEL CLERK", "CLERK":
		return RoleParcelClerk
	case "M", "MESSENGER", "DELIVERY":
		return RoleMessenger
	case "A", "ADMIN", "ADMINISTRATOR":
		return RoleAdmin
	default:
		return s
	}
}

func RoleToLegacyCode(role string) string {
	switch NormalizeRole(role) {
	case RoleParcelClerk:
		return "P"
	case RoleMessenger:
		return "M"
	case RoleAdmin:
		return "A"
	default:
		return ""
	}
}

func RoleDisplayName(role string) string {
	switch NormalizeRole(role) {
	case RoleParcelClerk:
		return "Parcel Clerk"
	case RoleMessenger:
		return "Messenger"
	case RoleAdmin:
		return "Admin"
	default:
		return role
	}
}

func NormalizeDeliveryType(s string) string {
	s = strings.ToUpper(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, " ", "_")
	switch s {
	case "STANDARD", "NORMAL":
		return DeliveryStandard
	case "EXPRESS", "FAST", "URGENT":
		return DeliveryExpress
	default:
		return s
	}
}

func NormalizeStatus(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.Join(strings.Fields(s), " ")
	switch s {
	case "pending", "wait", "waiting":
		return StatusPending
	case "in transit", "transit", "start", "start delivery", "on delivery", "doing delivery", "doing":
		return StatusInTransit
	case "delivered", "complete", "success", "done":
		return StatusDelivered
	case "delivery failed", "failed", "incomplete", "unsuccessful":
		return StatusDeliveryFailed
	default:
		return strings.TrimSpace(s)
	}
}

func LegacyStatus(status string) string {
	switch NormalizeStatus(status) {
	case StatusDelivered:
		return "Complete"
	case StatusDeliveryFailed:
		return "Incomplete"
	default:
		return status
	}
}

func DefaultTrackingDescription(status string) string {
	switch NormalizeStatus(status) {
	case StatusPending:
		return "Parcel registered at parcel counter"
	case StatusInTransit:
		return "Messenger picked up parcel for delivery"
	case StatusDelivered:
		return "Parcel delivered to receiver successfully"
	case StatusDeliveryFailed:
		return "Delivery attempt was unsuccessful"
	default:
		return "Tracking status updated"
	}
}

func IsValidPhone(phone string) bool {
	return phoneRegex.MatchString(strings.TrimSpace(phone))
}

func ZoneFromProvince(province string) string {
	p := strings.TrimSpace(strings.ToLower(province))
	switch p {
	case "bangkok", "กรุงเทพมหานคร", "nonthaburi", "นนทบุรี", "pathum thani", "ปทุมธานี", "samut prakan", "สมุทรปราการ", "bangkok metropolitan":
		return ZoneBangkokMetro
	default:
		return ZoneUpcountry
	}
}

func GenerateCode(prefix string) string {
	var b [4]byte
	_, _ = rand.Read(b[:])
	stamp := time.Now().UnixNano() % 1000000
	return fmt.Sprintf("%s-%06d%s", strings.ToUpper(prefix), stamp, strings.ToUpper(hex.EncodeToString(b[:]))[:4])
}

func RoundCurrency(v float64) float64 {
	return math.Round(v*100) / 100
}

func HashPassword(raw, salt string) string {
	h := sha256.Sum256([]byte(salt + raw))
	return hex.EncodeToString(h[:])
}

func CheckPassword(raw, hash, salt string) bool {
	return strings.EqualFold(HashPassword(raw, salt), strings.TrimSpace(hash))
}

func ParseBearerToken(authHeader string) string {
	if authHeader == "" {
		return ""
	}
	parts := strings.SplitN(strings.TrimSpace(authHeader), " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func BuildJWT(claims Claims, secret string) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	enc := base64.RawURLEncoding
	headerPart := enc.EncodeToString(headerJSON)
	payloadPart := enc.EncodeToString(payloadJSON)
	signingInput := headerPart + "." + payloadPart
	sig := signHMAC(signingInput, secret)
	return signingInput + "." + sig, nil
}

func ParseJWT(token, secret string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}
	signingInput := parts[0] + "." + parts[1]
	expected := signHMAC(signingInput, secret)
	if !hmac.Equal([]byte(expected), []byte(parts[2])) {
		return nil, errors.New("invalid token signature")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid token payload")
	}
	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, errors.New("invalid token claims")
	}
	if claims.ExpiresAt > 0 && time.Now().Unix() >= claims.ExpiresAt {
		return nil, errors.New("token expired")
	}
	return &claims, nil
}

func signHMAC(data, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
