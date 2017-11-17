package noaweb

import (
	"errors"
	"net/http"

	"github.com/gorilla/csrf"
)

// HelperFunctions struct used to structure package
type HelperFunctions struct{}

// Helper variable used to structure package
var Helper HelperFunctions

// GetUserConfig returns the given value as string
func (HelperFunctions) GetUserConfig(cfg string) (string, error) {
	if val, ok := noawebinst.UserConfig[cfg]; ok {
		return val, nil
	}
	return "", errors.New("Could not find UserConfig " + cfg)
}

// GetCSRF uses gorilla/csrf package to get a valid csrf token.
func (HelperFunctions) GetCSRF(r *http.Request) string {
	return csrf.Token(r)
}

// GetFaviconString returns favicon string to be included in header.
// All favicons must be available at /favicons/favicon-size.png.
// The following sizes must be present.
//  - favicon-32.png
//  - favicon-96.png
//  - favicon-120.png
//  - favicon-180.png
//  - favicon-152.png
//  - favicon-167.png
//  - favicon-128.png
//  - favicon-196.png
//  - favicon-228.png
func (HelperFunctions) GetFaviconString() string {
	str := `
<!-- General favicons -->
	<link rel="icon" type="image/png" href="/favicons/favicon-32.png" sizes="32x32">
	<link rel="icon" type="image/png" href="/favicons/favicon-96.png" sizes="96x96">

	<!-- Apple devices -->
	<link rel="apple-touch-icon" href="/favicons/favicon-120.png">
	<link rel="apple-touch-icon" sizes="180x180" href="/favicons/favicon-180.png">
	<link rel="apple-touch-icon" sizes="152x152" href="/favicons/favicon-152.png">
	<link rel="apple-touch-icon" sizes="167x167" href="/favicons/favicon-167.png">

	<!-- Other -->
	<link rel="icon" type="image/png" href="/favicons/favicon-128.png" sizes="128x128">
	<link rel="icon" type="image/png" href="/favicons/favicon-196.png" sizes="196x196">
	<link rel="icon" type="image/png" href="/favicons/favicon-228.png" sizes="228x228">`

	return str
}
