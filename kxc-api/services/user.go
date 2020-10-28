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
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dgrijalva/jwt-go"
	"github.com/didil/kubexcloud/kxc-api/requests"
	"github.com/didil/kubexcloud/kxc-api/responses"
	cloudv1alpha1 "github.com/didil/kubexcloud/kxc-operator/api/v1alpha1"
)

// UserSvc interface
type UserSvc interface {
	Login(ctx context.Context, userName, password string) (string, error)
	Create(ctx context.Context, reqData *requests.CreateUser) error
	HasRole(ctx context.Context, userName, role string) (bool, error)
	List(ctx context.Context) (*responses.ListUser, error)
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
	user, err := svc.find(ctx, userName)
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

const (
	UserRoleRegular string = "regular"
	UserRoleAdmin   string = "admin"
)

func (svc *UserService) Create(ctx context.Context, reqData *requests.CreateUser) error {
	client := svc.k8sSvc.Client()

	// check if the user exists
	user, err := svc.find(ctx, reqData.Name)
	if err != nil {
		return err
	}
	if user != nil {
		return fmt.Errorf("user already exists: %s", reqData.Name)
	}

	if len(reqData.Password) < minPasswordLength {
		return fmt.Errorf("password should be at least already %v chars long", minPasswordLength)
	}

	if reqData.Role != UserRoleRegular && reqData.Role != UserRoleAdmin {
		return fmt.Errorf("unknown role: %s", reqData.Role)
	}

	// hash password
	passwordHash, err := hashAndSalt([]byte(reqData.Password))
	if err != nil {
		return err
	}

	user = &cloudv1alpha1.UserAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: reqData.Name,
		},
		Spec: cloudv1alpha1.UserAccountSpec{
			Password: passwordHash,
			Role:     reqData.Role,
		},
	}

	err = client.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("create user: %v", err)
	}

	return nil
}

func (svc *UserService) find(ctx context.Context, userName string) (*cloudv1alpha1.UserAccount, error) {
	client := svc.k8sSvc.Client()

	user := &cloudv1alpha1.UserAccount{}
	err := client.Get(ctx, types.NamespacedName{Name: userName}, user)
	if errors.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find user: %v", err)
	}

	return user, nil
}

func (svc *UserService) HasRole(ctx context.Context, userName, role string) (bool, error) {
	user, err := svc.find(ctx, userName)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, fmt.Errorf("user doesn't exist: %s", userName)
	}

	return user.Spec.Role == role, nil
}

func (svc *UserService) List(ctx context.Context) (*responses.ListUser, error) {
	cl := svc.k8sSvc.Client()

	userList := &cloudv1alpha1.UserAccountList{}
	listOpts := []client.ListOption{}
	if err := cl.List(ctx, userList, listOpts...); err != nil {
		return nil, fmt.Errorf("failed to list users: %v", err)
	}

	respData := &responses.ListUser{
		Users: []responses.ListUserEntry{},
	}

	for _, user := range userList.Items {
		respData.Users = append(respData.Users, responses.ListUserEntry{
			Name: user.Name,
			Role: user.Spec.Role,
		})
	}

	return respData, nil
}

// auth helpers

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
