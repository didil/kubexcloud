package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/dgrijalva/jwt-go"
	cloudv1alpha1 "github.com/didil/kubexcloud/kxc-operator/api/v1alpha1"
)

// UserSvc interface
type UserSvc interface {
	Login(ctx context.Context, userName, password string) (string, error)
	Create(ctx context.Context, userName, password string) error
}

type UserService struct {
	k8sSvc K8sSvc
}

// NewUserService builds a new user service
func NewUserService(k8sSvc K8sSvc) *UserService {
	return &UserService{
		k8sSvc: k8sSvc,
	}
}

func (svc *UserService) Login(ctx context.Context, userName, password string) (string, error) {
	// check if the user exists
	user, err := svc.Get(ctx, userName)
	if err != nil {
		return "", err

	}
	if user == nil {
		return "", fmt.Errorf("user not found: %s", userName)
	}

	ok, err := comparePasswords([]byte(user.Spec.Password), []byte(password))
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("password invalid")
	}

	token, err := SignJWT(userName)
	if err != nil {
		return "", err
	}

	return token, nil
}

const minPasswordLength = 6

func (svc *UserService) Create(ctx context.Context, userName, password string) error {
	client := svc.k8sSvc.Client()

	// check if the user exists
	user, err := svc.Get(ctx, userName)
	if err != nil {
		return err

	}
	if user != nil {
		return fmt.Errorf("user already exists: %s", userName)
	}

	if len(password) < minPasswordLength {
		return fmt.Errorf("password should be at least already %v chars long", minPasswordLength)
	}

	// hash password
	passwordHash, err := hashAndSalt([]byte(password))
	if err != nil {
		return err
	}

	user = &cloudv1alpha1.UserAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: userName,
		},
		Spec: cloudv1alpha1.UserAccountSpec{
			Password: passwordHash,
		},
	}

	err = client.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("create user: %v", err)
	}

	return nil
}

func (svc *UserService) Get(ctx context.Context, userName string) (*cloudv1alpha1.UserAccount, error) {
	client := svc.k8sSvc.Client()

	user := &cloudv1alpha1.UserAccount{}
	err := client.Get(ctx, types.NamespacedName{Name: userName}, user)
	if errors.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user: %v", err)
	}

	return user, nil
}

// hashAndSalt hashes and salts a password
func hashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// comparePasswords compares a hashed password with a plain password
func comparePasswords(passwordHash []byte, plainPwd []byte) (bool, error) {
	if len(passwordHash) == 0 && len(plainPwd) == 0 {
		return true, nil
	}

	if len(passwordHash) == 0 && len(plainPwd) > 0 {
		return false, nil
	}

	err := bcrypt.CompareHashAndPassword(passwordHash, plainPwd)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// jwt custom claims
type customClaims struct {
	UserName string `json:"user"`
	jwt.StandardClaims
}

func getJWTSecret() ([]byte, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) == 0 {
		return nil, fmt.Errorf("JWT signing Secret is empty")
	}

	return secret, nil
}

// SignJWT returns a jwt signed token
func SignJWT(userName string) (string, error) {
	secret, err := getJWTSecret()
	if err != nil {
		return "", err
	}

	claims := customClaims{
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(60 * 24 * time.Hour).Unix(),
			Issuer:    "kxc-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// ParseJWT parses a jwt token
func ParseJWT(tokenStr string) (string, error) {
	secret, err := getJWTSecret()
	if err != nil {
		return "", err
	}

	token, err := jwt.ParseWithClaims(tokenStr, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*customClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("jwt token invalid")
	}

	return claims.UserName, nil
}
