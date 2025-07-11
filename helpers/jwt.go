package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type SignedDetails struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Uid       string `json:"uid"`
	UserType  string `json:"user_type"`
	jwt.RegisteredClaims
}

var SECRET_KEY = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email, firstname, lastname, usertype, uid string) (string, string, error) {
	if SECRET_KEY == "" {
		return "", "", errors.New("SECRET_KEY environment variable not set")
	}

	claims := SignedDetails{
		Email:     email,
		FirstName: firstname,
		LastName:  lastname,
		UserType:  usertype,
		Uid:       uid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Puneet",
			Subject:   "access_token",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshClaims := SignedDetails{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Puneet",
			Subject:   "refresh_token",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create tokens with your custom jwt.NewWithClaims (pass your signing method)
	// Assuming SigningMethodHS256 is defined in your jwt package:
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	// Sign tokens
	tokenString, err := accessToken.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

func ValidateToken(signedToken string) (*SignedDetails, error) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(t *jwt.Token) (interface{}, error) {
		//Using a function instead of a static key provides flexibility:
		//  1.🔑 Dynamic keys	: You can look up the signing key from a DB or cache (e.g. based on kid header).
		//  2.🔀 Support for multiple algorithms: You can check token.Method.Alg() to reject unexpected algs.
		//  3.🔍 Safety	: Forces you to validate token algorithm inside the function (best practice).

		// Ensure it's HMAC and not something else (e.g. RSA, none)
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, errors.New("token expired")
}

// RefreshHandler validates refresh token and issues new tokens
func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("Authorization") // e.g. Bearer <token>
	if refreshToken == "" {
		http.Error(w, "No refresh token provided", http.StatusUnauthorized)
		return
	}

	if len(refreshToken) > 7 && refreshToken[:7] == "Bearer " {
		refreshToken = refreshToken[7:]
	}

	claims, err := ValidateToken(refreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Generate new tokens using the user details
	newAccessToken, newRefreshToken, err := GenerateAllTokens(
		claims.Email,
		claims.FirstName,
		claims.LastName,
		claims.UserType,
		claims.Uid,
	)
	if err != nil {
		http.Error(w, "Failed to generate new tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}
