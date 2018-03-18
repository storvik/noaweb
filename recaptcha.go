package noaweb

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

const recaptchaURL = "https://www.google.com/recaptcha/api/siteverify"

var recaptchaSecret string

// RecaptchaFunctions struct
type RecaptchaFunctions struct{}

// Recaptcha variable to be used when calling recaptcha functions
var Recaptcha RecaptchaFunctions

type recaptchaResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes,omitempty"`
}

// Init is used to set recaptcha secret key
func (RecaptchaFunctions) Init (secret string) error {
	recaptchaSecret = secret
	return nil
}

// Verify takes recaptcha response as an argument and uses the recaptcha api to
// validate the response. Returns success bool and error. 
func (RecaptchaFunctions) Verify (recaptResponse string) (bool, error) {
	var err error

	formValue := []byte(`secret=` + recaptchaSecret + `&response=` + recaptResponse)
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := netClient.Post(recaptchaURL, "application/x-www-form-urlencoded; charset=utf-8", bytes.NewBuffer(formValue))
	if err != nil {
		return false, errors.New("Noaweb.Recaptcha.Verify: Could not reach recaptcha endpoint.")
	}
	defer response.Body.Close()

	resultBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, errors.New("Noaweb.Recaptcha.Verify: Could not read recaptcha response body.")
	}

	var result recaptchaResponse
	err = json.Unmarshal(resultBody, &result)
	if err != nil {
		return false, errors.New("Noaweb.Recaptcha.Verify: Invalid recaptcha response body.")
	}
	return result.Success, nil
}
