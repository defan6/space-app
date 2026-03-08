package service

import "errors"

var ErrNotFound = errors.New("order not found")

var ErrDoesNotEnoughPart = errors.New("does not enough part")
