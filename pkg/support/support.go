package support

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/razorpay/razorpay-go"
	"github.com/stebinsabu13/ecommerce-api/pkg/config"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// password authorization
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	} else {
		return true
	}
}

func BuildProfile(details utils.ResponseUsers, address []utils.Address) utils.ResponseUserDetails {
	return utils.ResponseUserDetails{
		FirstName:   details.FirstName,
		LastName:    details.LastName,
		Email:       details.Email,
		MobileNum:   details.MobileNum,
		ReferalCode: details.ReferalCode,
		Address:     address,
	}
}

func Email_validater(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if re.MatchString(email) {
		return true
	} else {
		return false
	}
}

func MobileNum_validater(number string) bool {
	re := regexp.MustCompile(`^[0-9]{10}$`)
	if re.MatchString(number) {
		return true
	} else {
		return false
	}
}

func GenerateRazorpayOrder(razorPayAmount int, recieptIdOptional string) (razorpayOrderID string, err error) {
	// get razor pay key and secret
	razorpayKey := config.GetCofig().RAZORPAYKEY
	razorpaySecret := config.GetCofig().RAZORPAYSECRET

	//create a razorpay client
	client := razorpay.NewClient(razorpayKey, razorpaySecret)

	data := map[string]interface{}{
		"amount":   razorPayAmount,
		"currency": "INR",
		"receipt":  recieptIdOptional,
	}
	// create an order on razor pay
	razorpayRes, err := client.Order.Create(data, nil)
	if err != nil {
		return razorpayOrderID, fmt.Errorf("fadil to create razorpay order for amount %v", razorPayAmount)
	}

	razorpayOrderID = razorpayRes["id"].(string)

	return razorpayOrderID, nil
}

func VeifyRazorpayPayment(razorpayOrderID, razorpayPaymentID, razorpaySignature string) error {

	razorpayKey := config.GetCofig().RAZORPAYKEY
	razorPaySecret := config.GetCofig().RAZORPAYSECRET

	//varify signature
	data := razorpayOrderID + "|" + razorpayPaymentID
	h := hmac.New(sha256.New, []byte(razorPaySecret))
	_, err := h.Write([]byte(data))
	if err != nil {
		return errors.New("faild to veify signature")
	}

	sha := hex.EncodeToString(h.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(sha), []byte(razorpaySignature)) != 1 {
		return errors.New("razorpay signature not match")
	}

	// then vefiy payment
	client := razorpay.NewClient(razorpayKey, razorPaySecret)

	// fetch payment and vefify
	payment, err := client.Payment.Fetch(razorpayPaymentID, nil, nil)

	if err != nil {
		return err
	}

	// check payment status
	if payment["status"] != "captured" {
		return fmt.Errorf("faild to varify payment \nrazorpay payment with payment_id %v", razorpayPaymentID)
	}

	return nil
}

func ReferalCodeGenerator() string {
	uuidObj := uuid.New()
	uuidString := strings.Replace(uuidObj.String(), "-", "", -1)
	return uuidString[:6]
}
