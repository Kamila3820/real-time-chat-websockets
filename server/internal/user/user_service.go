package user

import (
	"context"
	"fmt"
	"server/util"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	secretKey = "secret"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repository Repository) Service {
	return &service{
		repository,
		time.Duration(2) * time.Second,
	}
}

func (s *service) CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// TODO: hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	newUser := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	user, err := s.Repository.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	res := &CreateUserRes{
		ID:       strconv.Itoa(int(user.ID)),
		Username: user.Username,
		Email:    user.Email,
	}

	return res, nil
}

type MyJWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *service) Login(c context.Context, req *LoginReq) (*LoginRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	user, err := s.Repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return &LoginRes{}, err
	}

	fmt.Println("Email: ", req.Email)
	fmt.Println("Password: ", req.Password)

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		return &LoginRes{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		ID:       strconv.Itoa(int(user.ID)),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(user.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return &LoginRes{}, err
	}

	res := &LoginRes{
		accessToken: ss,
		ID:          user.ID,
		Username:    user.Username,
	}

	return res, nil
}
