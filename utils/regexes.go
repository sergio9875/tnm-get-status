package utils

import (
	"encoding/base64"
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
	"time"
)

func RegisterValidations() *validator.Validate {
	validatorObj := validator.New()

	_ = validatorObj.RegisterValidation("token", isToken)
	_ = validatorObj.RegisterValidation("transactionToken", isTransactionToken)
	_ = validatorObj.RegisterValidation("amount", isAmount)
	_ = validatorObj.RegisterValidation("textOnly", isTextOnly)
	_ = validatorObj.RegisterValidation("name", isName)
	_ = validatorObj.RegisterValidation("phonePrefix", isPhonePrefix)
	_ = validatorObj.RegisterValidation("boolean", isBoolean)
	_ = validatorObj.RegisterValidation("boolean2", isBoolean2)
	_ = validatorObj.RegisterValidation("boolean3", isBoolean3)
	_ = validatorObj.RegisterValidation("creditCardNumber", isCreditCardNumber)
	_ = validatorObj.RegisterValidation("cvv", isCvv)
	_ = validatorObj.RegisterValidation("iata", isIata)
	_ = validatorObj.RegisterValidation("password", isPassword)
	_ = validatorObj.RegisterValidation("iso", isIso)
	_ = validatorObj.RegisterValidation("creditCardExpiry", isCreditCardExpiry)
	_ = validatorObj.RegisterValidation("text", isText)
	_ = validatorObj.RegisterValidation("unicodeText", isUnicodeText)
	_ = validatorObj.RegisterValidation("phone", isPhone)
	_ = validatorObj.RegisterValidation("ptlType", isPTLtype)
	_ = validatorObj.RegisterValidation("datetime", isDateTime)
	_ = validatorObj.RegisterValidation("base64Check", isCustomBase64)
	_ = validatorObj.RegisterValidation("fileExt", isFileExt)
	_ = validatorObj.RegisterValidation("required_when", requiredWhen)
	_ = validatorObj.RegisterValidation("monthNum", isMonthNum)
	_ = validatorObj.RegisterValidation("yearNum", isYearNum)

	return validatorObj
}

const (
	token            = `(?i)^\{?[A-Z0-9]{8}-[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{12}\}?$`
	amount           = `^[0-9]+(\.[0-9]{1,4})?$`
	textNotEmpty     = ``
	textOnly         = `(?i)^[a-z\s]*$`
	text             = `(?i)^[a-z0-9 .\-\/\[\]\(\):,&%]+$`
	unicodeText      = `(?i)^[a-z0-9ßäöüÄÖÜ .\-\/\[\]:,&%]+$`
	name             = `^(?:[,.\p{L}\p{Mn}\p{Pd}\'\x{2019}]+\s?[,.\p{L}\p{Mn}\p{Pd}\'\x{2019}]+\s?)+$` //Name include special chars like French
	phone            = `^[0-9]{6,20}?$`
	phonePrefix      = `^[0-9]{1,3}?`                  //PhonePrefix(int) - between 1 and 3 digits
	boolean          = `^(1|0)$`                       //Boolean with true/false or 1/0
	boolean2         = `^(1|0)$`                       //Boolean with 1/0 only
	boolean3         = `^(true|false|1|0)$`            //Boolean with true/false or 1/0
	creditCardNumber = `^[0-9]{15,19}$`                //CreditCardNumber
	cvv              = `^[0-9]{3,4}$`                  //CVV
	iata             = `(?i)^[a-z\s]{3}$`              //3 letters of iata code
	password         = `(?i)^([0-9a-z!@#$%&()?]){6,}$` //chars for password
	iso              = `(?i)^[a-z]{2}$`                //2 letters code
	creditCardExpiry = `(?i)^[0-9]{4}$`                //4 numbers of expiry card
	monthNum         = `(?i)(^[0]{1}[1-9]{1})|^([1]{1}[0-2]{1})`
	yearNum          = `(?i)^[0-9]{4}$`
)

var (
	tokenRegex            = regexp.MustCompile(token)
	amountRegex           = regexp.MustCompile(amount)
	textOnlyRegex         = regexp.MustCompile(textOnly)
	textRegex             = regexp.MustCompile(text)
	unicodeTextRegex      = regexp.MustCompile(unicodeText)
	nameRegex             = regexp.MustCompile(name)
	phoneRegex            = regexp.MustCompile(phone)
	phonePrefixRegex      = regexp.MustCompile(phonePrefix)
	booleanRegex          = regexp.MustCompile(boolean)
	boolean2Regex         = regexp.MustCompile(boolean2)
	boolean3Regex         = regexp.MustCompile(boolean3)
	creditCardNumberRegex = regexp.MustCompile(creditCardNumber)
	cvvRegex              = regexp.MustCompile(cvv)
	iataRegex             = regexp.MustCompile(iata)
	passwordRegex         = regexp.MustCompile(password)
	isoRegex              = regexp.MustCompile(iso)
	creditCardExpiryRegex = regexp.MustCompile(creditCardExpiry)
	monthNumRegex         = regexp.MustCompile(monthNum)
	yearNumRegex          = regexp.MustCompile(yearNum)
	rules                 = []string{"2006-01-02 15:04", "2006-01-02 15:04:05", "2006-01-02", "2006/01/02 15:04", "2006/01/02 15:04:05", "2006/01/02"}
)

func isToken(fl validator.FieldLevel) bool {
	tokenValue := fl.Field().String()
	return tokenRegex.MatchString(tokenValue) && len(tokenValue) == 36
}

func isTransactionToken(fl validator.FieldLevel) bool {
	tokenValue := fl.Field().String()
	return (tokenRegex.MatchString(tokenValue) && len(tokenValue) == 36) || len(tokenValue) < 10
}

func isAmount(fl validator.FieldLevel) bool {
	return amountRegex.MatchString(fl.Field().String())
}

func isTextOnly(fl validator.FieldLevel) bool {
	return textOnlyRegex.MatchString(fl.Field().String())
}

func isName(fl validator.FieldLevel) bool {
	return nameRegex.MatchString(fl.Field().String())
}

func isPhonePrefix(fl validator.FieldLevel) bool {
	return phonePrefixRegex.MatchString(fl.Field().String())
}

func isBoolean(fl validator.FieldLevel) bool {
	return booleanRegex.MatchString(fl.Field().String())
}

func isBoolean2(fl validator.FieldLevel) bool {
	return boolean2Regex.MatchString(fl.Field().String())
}

func isBoolean3(fl validator.FieldLevel) bool {
	return boolean3Regex.MatchString(fl.Field().String())
}

func isCreditCardNumber(fl validator.FieldLevel) bool {
	return creditCardNumberRegex.MatchString(fl.Field().String())
}

func isCvv(fl validator.FieldLevel) bool {
	return cvvRegex.MatchString(fl.Field().String())
}

func isIata(fl validator.FieldLevel) bool {
	return iataRegex.MatchString(fl.Field().String())
}

func isPassword(fl validator.FieldLevel) bool {
	return passwordRegex.MatchString(fl.Field().String())
}

func isIso(fl validator.FieldLevel) bool {
	return isoRegex.MatchString(fl.Field().String())
}

func isCreditCardExpiry(fl validator.FieldLevel) bool {
	return creditCardExpiryRegex.MatchString(fl.Field().String())
}

func isMonthNum(fl validator.FieldLevel) bool {
	return monthNumRegex.MatchString(fl.Field().String())
}

func isYearNum(fl validator.FieldLevel) bool {
	return yearNumRegex.MatchString(fl.Field().String())
}

func removeNewline(text string) string {
	return strings.ReplaceAll(text, "\n", " ")
}

func removeCarriageReturn(value string) string {
	return strings.ReplaceAll(value, "\r", " ")
}

func removePlusSign(value string) string {
	return strings.ReplaceAll(value, "+", "")
}

func isText(fl validator.FieldLevel) bool {
	textValue := removeNewline(removeCarriageReturn(fl.Field().String()))
	return textRegex.MatchString(textValue)
}
func isUnicodeText(fl validator.FieldLevel) bool {
	textValue := removeNewline(removeCarriageReturn(fl.Field().String()))
	return unicodeTextRegex.MatchString(textValue)
}

func isPhone(fl validator.FieldLevel) bool {
	value := removePlusSign(fl.Field().String())
	return phoneRegex.MatchString(value)
}

// Non Regex Validation

func isPTLtype(fl validator.FieldLevel) bool {
	value := strings.ToLower(fl.Field().String())
	return value == "hours" || value == "minutes"
}

func isDateTime(fl validator.FieldLevel) bool {
	for _, rule := range rules {
		_, err := time.Parse(rule, fl.Field().String())
		if err == nil {
			return true
		}
	}
	return false
}

func isCustomBase64(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	decodeValue, _ := base64.StdEncoding.DecodeString(value)

	encodedValue := base64.StdEncoding.EncodeToString(decodeValue)

	return encodedValue == value
}

func isFileExt(fl validator.FieldLevel) bool {
	value := strings.ToLower(fl.Field().String())
	return value == "jpeg" || value == "jpg" || value == "png" || value == "gif" || value == "pdf"
}

func requiredWhen(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	splitValue := strings.Split(fl.Param(), "=")
	data := fl.Parent().Elem().FieldByName(splitValue[0])

	if data.String() == splitValue[1] && len(value) == 0 {
		return false
	}
	return true
}
