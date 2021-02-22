package sessions

import (
	"crypto/sha256"
	"errors"
	"crypto/hmac"
	"crypto/rand"
	"fmt"
	"encoding/base64"
)

//InvalidSessionID represents an empty, invalid session ID
const InvalidSessionID SessionID = ""

//idLength is the length of the ID portion
const idLength = 32

//signedLength is the full length of the signed session ID
//(ID portion plus signature)
const signedLength = idLength + sha256.Size

//SessionID represents a valid, digitally-signed session ID.
//This is a base64 URL encoded string created from a byte slice
//where the first `idLength` bytes are crytographically random
//bytes representing the unique session ID, and the remaining bytes
//are an HMAC hash of those ID bytes (i.e., a digital signature).
//The byte slice layout is like so:
//+-----------------------------------------------------+
//|...32 crypto random bytes...|HMAC hash of those bytes|
//+-----------------------------------------------------+
type SessionID string

//ErrInvalidID is returned when an invalid session id is passed to ValidateID()
var ErrInvalidID = errors.New("Invalid Session ID")

//NewSessionID creates and returns a new digitally-signed session ID,
//using `signingKey` as the HMAC signing key. An error is returned only
//if there was an error generating random bytes for the session ID
func NewSessionID(signingKey string) (SessionID, error) {
	//TODO: if `signingKey` is zero-length, return InvalidSessionID
	//and an error indicating that it may not be empty
	if len(signingKey) == 0 {
		return InvalidSessionID, fmt.Errorf("Error with session ID")
	}

	slice := make([]byte, idLength)
    _, err := rand.Read(slice)
    if err != nil {
        return InvalidSessionID, fmt.Errorf("error generating salt: %v", err)
    }
	h := hmac.New(sha256.New, []byte(signingKey))
	h.Write(slice)
	signature := h.Sum(nil)
	ss := append(slice, signature...)
	encodedSlice := SessionID(base64.URLEncoding.EncodeToString(ss))
	
	
	return encodedSlice, nil
}

//ValidateID validates the string in the `id` parameter
//using the `signingKey` as the HMAC signing key
//and returns an error if invalid, or a SessionID if valid
func ValidateID(id string, signingKey string) (SessionID, error) {

	decode, err := base64.URLEncoding.DecodeString(id)
	if err != nil {
		return InvalidSessionID, fmt.Errorf("Error validating ID: %v", err)
	}
	idPortion := decode[:idLength]
	previousSig := decode[idLength:]

	h := hmac.New(sha256.New, []byte(signingKey))
	h.Write([]byte(idPortion))
	signature := h.Sum(nil)
	if hmac.Equal(previousSig, signature) {
		return SessionID(id), nil
	}
	return InvalidSessionID, ErrInvalidID
}

//String returns a string representation of the sessionID
func (sid SessionID) String() string {
	return string(sid)
}
