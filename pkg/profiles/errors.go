package profiles

import "errors"

var (
	ErrDeletingDefaultProfile   = errors.New("cannot delete the default profile")
	ErrProfileIsEmpty           = errors.New("error profile is empty")
	ErrProfileIncorrectType     = errors.New("erros profile is not of type ProfileConfig")
	ErrCreatingPlatform         = errors.New("error when creating platform")
	ErrCreatingNewProfile       = errors.New("error creating profile")
	ErrUnknownProfileDriverType = errors.New("error unknown profile driver type")
)
