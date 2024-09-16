package service

import "errors"

var (
	NotExpectedErrorFromStorage = errors.New("not expected error from storage")
	TenderNotFound              = errors.New("tender not found")
	OrganizationNotFound        = errors.New("organization not found")
	UserNotFound                = errors.New("user not found")
	TenderAlreadyExists         = errors.New("tender already exists")
	NotCorrectParams            = errors.New("not correct params")
	BidNotFound                 = errors.New("bid not found")
)
