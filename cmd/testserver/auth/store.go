package auth

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
)

// User is one authenticated principal.
type User struct {
	Username string `yaml:"username"`
	// PasswordHash is a bcrypt hash. Prefer this in real deployments.
	PasswordHash string `yaml:"password_hash"`
	// Password is a lab-only plaintext password; hashed into memory at load time.
	Password            string   `yaml:"password"`
	Role                string   `yaml:"role"`
	HypervisorAllowlist []string `yaml:"hypervisor_allowlist"`
}

// Config holds the multi-user auth store.
type Config struct {
	Users []User `yaml:"users"`
}

// Claims are JWT claims issued by the testserver.
type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// Store validates credentials and issues/verifies JWTs.
type Store struct {
	users      map[string]User
	appName    string
	signingKey []byte
	lifetime   time.Duration
}

// LoadUsersFile reads a YAML user store from path.
func LoadUsersFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	if len(cfg.Users) == 0 {
		return Config{}, errors.New("users file contains no users")
	}
	return cfg, nil
}

// NewStore builds an auth store from config.
func NewStore(cfg Config, appName string, signingKey []byte, lifetime time.Duration) (*Store, error) {
	if appName == "" {
		return nil, errors.New("application name is required")
	}
	if len(signingKey) == 0 {
		return nil, errors.New("jwt signing key is required")
	}
	users := make(map[string]User, len(cfg.Users))
	for _, u := range cfg.Users {
		if u.Username == "" {
			return nil, fmt.Errorf("user entry missing username")
		}
		if u.PasswordHash == "" {
			if u.Password == "" {
				return nil, fmt.Errorf("user %q missing password_hash (or lab password)", u.Username)
			}
			hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
			if err != nil {
				return nil, fmt.Errorf("hash password for %q: %w", u.Username, err)
			}
			u.PasswordHash = string(hash)
			u.Password = ""
		}
		users[u.Username] = u
	}
	return &Store{
		users:      users,
		appName:    appName,
		signingKey: signingKey,
		lifetime:   lifetime,
	}, nil
}

// AuthenticateBasic validates HTTP basic credentials and returns a JWT.
func (s *Store) AuthenticateBasic(username, password string) (string, error) {
	user, ok := s.users[username]
	if !ok {
		return "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	now := time.Now()
	claims := Claims{
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Username,
			Issuer:    s.appName,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.lifetime)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(s.signingKey)
}

// ParseBearer validates a Bearer token and returns claims.
func (s *Store) ParseBearer(authorizationHeader string) (*Claims, error) {
	const prefix = "Bearer "
	if !strings.HasPrefix(authorizationHeader, prefix) {
		return nil, errors.New("missing bearer token")
	}
	raw := strings.TrimSpace(strings.TrimPrefix(authorizationHeader, prefix))
	token, err := jwt.ParseWithClaims(raw, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS512 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.Issuer != s.appName {
		return nil, errors.New("invalid token issuer")
	}
	return claims, nil
}

// AllowHypervisor reports whether the user may use the given hypervisor URI.
func (s *Store) AllowHypervisor(username, uri string) bool {
	user, ok := s.users[username]
	if !ok {
		return false
	}
	if len(user.HypervisorAllowlist) == 0 {
		return true
	}
	for _, allowed := range user.HypervisorAllowlist {
		if allowed == uri {
			return true
		}
	}
	return false
}
