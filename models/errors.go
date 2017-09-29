package models

import "errors"

// -----------------------
// ERRORS -------------
// -----------------------

// ErrInvalidPayload variable for Error handling
var ErrInvalidPayload = errors.New("Invalid payload")

// ErrRecordNotFound variable for Error handling
var ErrRecordNotFound = errors.New("Record was not found")

// ErrMissingClientID variable for Error handling
var ErrMissingClientID = errors.New("ClientID cannot be empty")

// ErrMissingVitalFields variable for Error Handling
var ErrMissingVitalFields = errors.New("One of the following fields cannot be empty: PEMFile, P12File or FCMAuthKey")

// ErrMissingVitalIOSFields variable for Error Handling
var ErrMissingVitalIOSFields = errors.New("The following fields cannot be empty: PEMFile (or P12File), Passphrase and BundleIdentifier")

// ErrMissingPassPhrase variable for Error Handling
var ErrMissingPassPhrase = errors.New("PassPhrase cannot be empty")

// ErrMissingBundleIdentifier variable for Error Handling
var ErrMissingBundleIdentifier = errors.New("BundleIdentifier cannot be empty")

// ErrFailedToLoadPEMFile variable for Error handling
var ErrFailedToLoadPEMFile = errors.New("Failed to loaded either PEM or P12 file")

// ErrFailedToSendPush variable for Error handling
var ErrFailedToSendPush = errors.New("Failed to send push (or all pushes)")
