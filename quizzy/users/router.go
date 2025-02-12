package users

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
)

func ConfigureRoutes(rt *gin.RouterGroup) {
	rt.POST("/users", createUser)
}

func getGooglePublicKeys() (map[string]*rsa.PublicKey, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v3/certs")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var keys struct {
		Keys []struct {
			Kid string `json:"kid"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}

	if err := json.Unmarshal(body, &keys); err != nil {
		return nil, err
	}

	// Convertir les clés en format RSA
	publicKeys := make(map[string]*rsa.PublicKey)
	for _, key := range keys.Keys {
		nBytes, _ := base64.RawURLEncoding.DecodeString(key.N)
		eBytes, _ := base64.RawURLEncoding.DecodeString(key.E)

		e := 0
		for _, b := range eBytes {
			e = e<<8 + int(b)
		}

		pubKey := &rsa.PublicKey{
			N: new(big.Int).SetBytes(nBytes),
			E: e,
		}
		publicKeys[key.Kid] = pubKey
	}

	return publicKeys, nil
}

// Fonction pour extraire le kid du JWT
func getJWTHeaderKid(tokenString string) (string, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) < 2 {
		return "", errors.New("invalid token format")
	}

	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", err
	}

	var header struct {
		Kid string `json:"kid"`
	}
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		return "", err
	}

	return header.Kid, nil
}

func createUser(c *gin.Context) {
	if _, exists := c.Get("firebase-services"); !exists {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"error": "Missing or invalid token"})
			return
		}

		// Extraire le token JWT
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Extraire le kid du token
		kid, err := getJWTHeaderKid(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid JWT format"})
			return
		}

		// Récupérer les clés publiques de Google
		publicKeys, err := getGooglePublicKeys()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch Google public keys"})
			return
		}

		// Vérifier que la clé publique correspondante existe
		publicKey, exists := publicKeys[kid]
		if !exists {
			c.JSON(401, gin.H{"error": "No matching public key found"})
			return
		}

		// Vérifier le token avec la clé publique
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return publicKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			return
		}

		// Afficher les claims
		fmt.Println("Token Claims:", claims)
		c.JSON(200, gin.H{"message": "Token valid", "claims": claims})
	}
}
