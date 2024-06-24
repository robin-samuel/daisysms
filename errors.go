package daisysms

import "errors"

var (
	ErrMaxPriceExceeded     = errors.New("max price exceeded")
	ErrNoNumbers            = errors.New("no numbers left")
	ErrTooManyActiveRentals = errors.New("too many active rentals")
	ErrNoMoney              = errors.New("no money")

	ErrWrongID        = errors.New("wrong id")
	ErrRentalCanceled = errors.New("rental canceled")
)
