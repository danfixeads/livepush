package models

import "errors"

// -----------------------
// ERRORS -------------
// -----------------------

// ErrRecordNotFound variable for Error handling
var ErrRecordNotFound = errors.New("Record was not found")

// ErrMissingClientID variable for Error handling
var ErrMissingClientID = errors.New("ClientID cannot be empty")

// ErrMissingVitalFields variable for Error Handling
var ErrMissingVitalFields = errors.New("One of the following fields cannot be empty: PEMFile, P12File or FCMAuthKey")

// ErrMissingPassPhrase variable for Error Handling
var ErrMissingPassPhrase = errors.New("PassPhrase cannot be empty")

// ErrMissingBundleIdentifier variable for Error Handling
var ErrMissingBundleIdentifier = errors.New("BundleIdentifier cannot be empty")
