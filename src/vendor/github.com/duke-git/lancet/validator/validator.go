// Copyright 2021 dudaodong@gmail.com. All rights reserved.
// Use of this source code is governed by MIT license

// Package validator implements some validate function for string.
package validator

import (
	"encoding/json"
	"net"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	alphaMatcher           *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z]+$`)
	letterRegexMatcher     *regexp.Regexp = regexp.MustCompile(`[a-zA-Z]`)
	intStrMatcher          *regexp.Regexp = regexp.MustCompile(`^[\+-]?\d+$`)
	urlMatcher             *regexp.Regexp = regexp.MustCompile(`^((ftp|http|https?):\/\/)?(\S+(:\S*)?@)?((([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(([a-zA-Z0-9]+([-\.][a-zA-Z0-9]+)*)|((www\.)?))?(([a-z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-z\x{00a1}-\x{ffff}]{2,}))?))(:(\d{1,5}))?((\/|\?|#)[^\s]*)?$`)
	dnsMatcher             *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z]([a-zA-Z0-9\-]+[\.]?)*[a-zA-Z0-9]$`)
	emailMatcher           *regexp.Regexp = regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
	chineseMobileMatcher   *regexp.Regexp = regexp.MustCompile(`^1(?:3\d|4[4-9]|5[0-35-9]|6[67]|7[013-8]|8\d|9\d)\d{8}$`)
	chineseIdMatcher       *regexp.Regexp = regexp.MustCompile(`^[1-9]\d{5}(18|19|20|21|22)\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`)
	chineseMatcher         *regexp.Regexp = regexp.MustCompile("[\u4e00-\u9fa5]")
	chinesePhoneMatcher    *regexp.Regexp = regexp.MustCompile(`\d{3}-\d{8}|\d{4}-\d{7}`)
	creditCardMatcher      *regexp.Regexp = regexp.MustCompile(`^(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|(222[1-9]|22[3-9][0-9]|2[3-6][0-9]{2}|27[01][0-9]|2720)[0-9]{12}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\\d{3})\\d{11}|6[27][0-9]{14})$`)
	base64Matcher          *regexp.Regexp = regexp.MustCompile(`^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$`)
	base64URLMatcher       *regexp.Regexp = regexp.MustCompile(`^([A-Za-z0-9_-]{4})*([A-Za-z0-9_-]{2}(==)?|[A-Za-z0-9_-]{3}=?)?$`)
	binMatcher             *regexp.Regexp = regexp.MustCompile(`^(0b)?[01]+$`)
	hexMatcher             *regexp.Regexp = regexp.MustCompile(`^(#|0x|0X)?[0-9a-fA-F]+$`)
	visaMatcher            *regexp.Regexp = regexp.MustCompile(`^4[0-9]{12}(?:[0-9]{3})?$`)
	masterCardMatcher      *regexp.Regexp = regexp.MustCompile(`^5[1-5][0-9]{14}$`)
	americanExpressMatcher *regexp.Regexp = regexp.MustCompile(`^3[47][0-9]{13}$`)
	unionPay               *regexp.Regexp = regexp.MustCompile("^62[0-5]\\d{13,16}$")
	chinaUnionPay          *regexp.Regexp = regexp.MustCompile(`^62[0-9]{14,17}$`)
)

// IsAlpha checks if the string contains only letters (a-zA-Z)
func IsAlpha(str string) bool {
	return alphaMatcher.MatchString(str)
}

// IsAllUpper check if the string is all upper case letters A-Z
func IsAllUpper(str string) bool {
	for _, r := range str {
		if !unicode.IsUpper(r) {
			return false
		}
	}
	return str != ""
}

// IsAllLower check if the string is all lower case letters a-z
func IsAllLower(str string) bool {
	for _, r := range str {
		if !unicode.IsLower(r) {
			return false
		}
	}
	return str != ""
}

// ContainUpper check if the string contain at least one upper case letter A-Z
func ContainUpper(str string) bool {
	for _, r := range str {
		if unicode.IsUpper(r) && unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

// ContainLower check if the string contain at least one lower case letter A-Z
func ContainLower(str string) bool {
	for _, r := range str {
		if unicode.IsLower(r) && unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

// ContainLetter check if the string contain at least one letter
func ContainLetter(str string) bool {
	return letterRegexMatcher.MatchString(str)
}

// IsJSON checks if the string is valid JSON
func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

// IsNumberStr check if the string can convert to a number.
func IsNumberStr(s string) bool {
	return IsIntStr(s) || IsFloatStr(s)
}

// IsFloatStr check if the string can convert to a float.
func IsFloatStr(str string) bool {
	_, e := strconv.ParseFloat(str, 64)
	return e == nil
}

// IsIntStr check if the string can convert to a integer.
func IsIntStr(str string) bool {
	return intStrMatcher.MatchString(str)
}

// IsIp check if the string is a ip address.
func IsIp(ipstr string) bool {
	ip := net.ParseIP(ipstr)
	return ip != nil
}

// IsIpV4 check if the string is a ipv4 address.
func IsIpV4(ipstr string) bool {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return false
	}
	return strings.Contains(ipstr, ".")
}

// IsIpV6 check if the string is a ipv6 address.
func IsIpV6(ipstr string) bool {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return false
	}
	return strings.Contains(ipstr, ":")
}

// IsPort check if the string is a valid net port.
func IsPort(str string) bool {
	if i, err := strconv.ParseInt(str, 10, 64); err == nil && i > 0 && i < 65536 {
		return true
	}
	return false
}

// IsUrl check if the string is url.
func IsUrl(str string) bool {
	if str == "" || len(str) >= 2083 || len(str) <= 3 || strings.HasPrefix(str, ".") {
		return false
	}
	u, err := url.Parse(str)
	if err != nil {
		return false
	}
	if strings.HasPrefix(u.Host, ".") {
		return false
	}
	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return false
	}

	return urlMatcher.MatchString(str)
}

// IsDns check if the string is dns.
func IsDns(dns string) bool {
	return dnsMatcher.MatchString(dns)
}

// IsEmail check if the string is a email address.
func IsEmail(email string) bool {
	return emailMatcher.MatchString(email)
}

// IsChineseMobile check if the string is chinese mobile number.
func IsChineseMobile(mobileNum string) bool {
	return chineseMobileMatcher.MatchString(mobileNum)
}

// IsChineseIdNum check if the string is chinese id number.
func IsChineseIdNum(id string) bool {
	return chineseIdMatcher.MatchString(id)
}

// ContainChinese check if the string contain mandarin chinese.
func ContainChinese(s string) bool {
	return chineseMatcher.MatchString(s)
}

// IsChinesePhone check if the string is chinese phone number.
// Valid chinese phone is xxx-xxxxxxxx or xxxx-xxxxxxx
func IsChinesePhone(phone string) bool {
	return chinesePhoneMatcher.MatchString(phone)
}

// IsCreditCard check if the string is credit card.
func IsCreditCard(creditCart string) bool {
	return creditCardMatcher.MatchString(creditCart)
}

// IsBase64 check if the string is base64 string.
func IsBase64(base64 string) bool {
	return base64Matcher.MatchString(base64)
}

// IsEmptyString check if the string is empty.
func IsEmptyString(str string) bool {
	return len(str) == 0
}

// IsRegexMatch check if the string match the regexp
func IsRegexMatch(str, regex string) bool {
	reg := regexp.MustCompile(regex)
	return reg.MatchString(str)
}

// IsStrongPassword check if the string is strong password, if len(password) is less than the length param, return false
// Strong password: alpha(lower+upper) + number + special chars(!@#$%^&*()?><)
func IsStrongPassword(password string, length int) bool {
	if len(password) < length {
		return false
	}
	var num, lower, upper, special bool
	for _, r := range password {
		switch {
		case unicode.IsDigit(r):
			num = true
		case unicode.IsUpper(r):
			upper = true
		case unicode.IsLower(r):
			lower = true
		case unicode.IsSymbol(r), unicode.IsPunct(r):
			special = true
		}
	}

	return num && lower && upper && special

	// go doesn't support regexp (?=re)
	//pattern := `^(?=.*[0-9])(?=.*[a-zA-Z])(?=.*[@#$%^&+=])(?=\S+$).$`
	//reg := regexp.MustCompile(pattern)
	//return reg.MatchString(password)
}

// IsWeakPassword check if the string is weak password
// Weak password: only letter or only number or letter + number
func IsWeakPassword(password string) bool {
	var num, letter, special bool
	for _, r := range password {
		switch {
		case unicode.IsDigit(r):
			num = true
		case unicode.IsLetter(r):
			letter = true
		case unicode.IsSymbol(r), unicode.IsPunct(r):
			special = true
		}
	}

	return (num || letter) && !special
}

// IsZeroValue checks if value is a zero value
func IsZeroValue(value interface{}) bool {
	if value == nil {
		return true
	}

	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return true
	}

	switch rv.Kind() {
	case reflect.String:
		return rv.Len() == 0
	case reflect.Bool:
		return !rv.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() == 0
	case reflect.Ptr, reflect.Chan, reflect.Func, reflect.Interface, reflect.Slice, reflect.Map:
		return rv.IsNil()
	}

	return reflect.DeepEqual(rv.Interface(), reflect.Zero(rv.Type()).Interface())
}

// IsGBK check if data encoding is gbk
// Note: this function is implemented by whether double bytes fall within the encoding range of gbk,
// while each byte of utf-8 encoding format falls within the encoding range of gbk.
// Therefore, utf8.valid() should be called first to check whether it is not utf-8 encoding,
// and then call IsGBK() to check gbk encoding. like below
/**
	data := []byte("你好")
	if utf8.Valid(data) {
		fmt.Println("data encoding is utf-8")
	}else if(IsGBK(data)) {
		fmt.Println("data encoding is GBK")
	}
	fmt.Println("data encoding is unknown")
**/
func IsGBK(data []byte) bool {
	i := 0
	for i < len(data) {
		if data[i] <= 0xff {
			i++
			continue
		} else {
			if data[i] >= 0x81 &&
				data[i] <= 0xfe &&
				data[i+1] >= 0x40 &&
				data[i+1] <= 0xfe &&
				data[i+1] != 0xf7 {
				i += 2
				continue
			} else {
				return false
			}
		}
	}

	return true
}

// IsNumberStr check if the value is number(integer, float) or not.
func IsNumber(v interface{}) bool {
	return IsInt(v) || IsFloat(v)
}

// IsFloat check if the value is float(float32, float34) or not.
func IsFloat(v interface{}) bool {
	switch v.(type) {
	case float32, float64:
		return true
	}
	return false
}

// IsInt check if the value is integer(int, unit) or not.
func IsInt(v interface{}) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
		return true
	}
	return false
}

// IsASCII checks if string is all ASCII char.
func IsASCII(str string) bool {
	for i := 0; i < len(str); i++ {
		if str[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// IsPrintable checks if string is all printable chars.
func IsPrintable(str string) bool {
	for _, r := range str {
		if !unicode.IsPrint(r) {
			if r == '\n' || r == '\r' || r == '\t' || r == '`' {
				continue
			}
			return false
		}
	}
	return true
}

// IsBin check if a give string is a valid binary value or not.
func IsBin(v string) bool {
	return binMatcher.MatchString(v)
}

// IsHex check if a give string is a valid hexadecimal value or not.
func IsHex(v string) bool {
	return hexMatcher.MatchString(v)
}

// IsBase64URL check if a give string is a valid URL-safe Base64 encoded string.
func IsBase64URL(v string) bool {
	return base64URLMatcher.MatchString(v)
}

// IsJWT check if a give string is a valid JSON Web Token (JWT).
func IsJWT(v string) bool {
	strings := strings.Split(v, ".")
	if len(strings) != 3 {
		return false
	}

	for _, s := range strings {
		if !IsBase64URL(s) {
			return false
		}
	}

	return true
}

// IsVisa check if a give string is a valid visa card nubmer or not.
func IsVisa(v string) bool {
	return visaMatcher.MatchString(v)
}

// IsMasterCard check if a give string is a valid master card nubmer or not.
func IsMasterCard(v string) bool {
	return masterCardMatcher.MatchString(v)
}

// IsAmericanExpress check if a give string is a valid american expression card nubmer or not.
func IsAmericanExpress(v string) bool {
	return americanExpressMatcher.MatchString(v)
}

// IsUnionPay check if a give string is a valid union pay nubmer or not.
func IsUnionPay(v string) bool {
	return unionPay.MatchString(v)
}

// IsChinaUnionPay check if a give string is a valid china union pay nubmer or not.
func IsChinaUnionPay(v string) bool {
	return chinaUnionPay.MatchString(v)
}
