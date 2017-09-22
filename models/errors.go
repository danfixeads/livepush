package models

import "errors"

// -----------------------
// ERRORS -------------
// -----------------------

// ErrMissingClientID variable for Error handling
var ErrMissingClientID = errors.New("ClientID cannot be empty")

// ErrMissingVitalFields variable for Error Handling
var ErrMissingVitalFields = errors.New("One of the following fields cannot be empty: PEMFile, P12File or FCMToken")

// ErrMissingPassPhrase variable for Error Handling
var ErrMissingPassPhrase = errors.New("PassPhrase cannot be empty")
