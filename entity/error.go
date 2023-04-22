package entity

import "errors"

var ErrUserNotFound = errors.New("user not found")

var ErrUserAlreadyExists = errors.New("user already exists")

var ErrAuthenticationAlreadyExists = errors.New("authentication already exists")
