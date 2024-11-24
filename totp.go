package main

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"strconv"

	"bufio"
	"bytes"
	"encoding/base32"
	"fmt"
	"image/png"
	"log"
	"os"
	"time"
)

func display(key *otp.Key, data []byte) {
	fmt.Printf("Issuer:       %s\n", key.Issuer())
	fmt.Printf("Account Name: %s\n", key.AccountName())
	fmt.Printf("Secret:       %s\n", key.Secret())
	fmt.Println("Writing PNG to qr-code.png....")
	os.WriteFile("qr-code.png", data, 0644)
	fmt.Println("")
	fmt.Println("Please add your TOTP to your OTP Application now!")
	fmt.Println("")
}

func promptForPasscode() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Passcode: ")
	text, _ := reader.ReadString('\n')
	return text
}

// Demo function, not used in main
// Generates Passcode using a UTF-8 (not base32) secret and custom parameters
func GeneratePassCode(utf8string string) string {
	secret := base32.StdEncoding.EncodeToString([]byte(utf8string))
	passcode, err := totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA512,
	})
	if err != nil {
		panic(err)
	}
	return passcode
}

func environmentVariable(Name string) string {
	Value, OK := os.LookupEnv(Name)
	if !OK {
		log.Fatalf("Missing environment variable %s", Name)
	}
	return Value
}
func main() {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      environmentVariable("TOTP_ISSUER"),
		AccountName: environmentVariable("TOTP_ACCOUNT_NAME"),
		Period: func() uint {
			P := environmentVariable("TOTP_PERIOD")
			Pu, err := strconv.ParseUint(P, 10, 32)
			if err != nil {
				log.Fatalf("Bad environment variable TOTP_PERIOD, %s", err)
			}
			return uint(Pu)
		}(),
	})
	if err != nil {
		panic(err)
	}
	// Convert TOTP key into a PNG
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		panic(err)
	}
	png.Encode(&buf, img)

	// display the QR code to the user.
	display(key, buf.Bytes())

	// Now Validate that the user's successfully added the passcode.
	fmt.Println("Validating TOTP...")
	passcode := promptForPasscode()
	valid := totp.Validate(passcode, key.Secret())
	if valid {
		println("Valid passcode!")
		os.Exit(0)
	} else {
		println("Invalid passcode!")
		os.Exit(1)
	}
}
