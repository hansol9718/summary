package sessions

import (
	"errors"
	"net/http"
	"fmt"
	"strings"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	//TODO:
	//- create a new SessionID
	//- save the sessionState to the store
	//- add a header to the ResponseWriter that looks like this:
	//    "Authorization: Bearer <sessionID>"
	//  where "<sessionID>" is replaced with the newly-created SessionID
	//  (note the constants declared for you above, which will help you avoid typos)
	sessionID, err := NewSessionID(signingKey)
	if err != nil {
		return InvalidSessionID, fmt.Errorf("Error when creating new sessionID: %v", err)
	}
	err = store.Save(sessionID, sessionState)
	if err != nil {
		return InvalidSessionID, fmt.Errorf("Error when saving session state: %v", err)
	}

	w.Header().Add(headerAuthorization, schemeBearer+sessionID.String())

	return sessionID, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	//TODO: get the value of the Authorization header,
	//or the "auth" query string parameter if no Authorization header is present,
	//and validate it. If it's valid, return the SessionID. If not
	//return the validation error.
	headerAuth := r.Header.Get(headerAuthorization)

	if len(headerAuth) == 0 {
		headerAuth = r.URL.Query().Get(paramAuthorization)
	}
	if !strings.HasPrefix(headerAuth, schemeBearer) {
		return InvalidSessionID, errors.New("no scheme prefix")
	}
	headerID := strings.TrimPrefix(headerAuth, schemeBearer)
	sessionID, err := ValidateID(headerID, signingKey)

	if err != nil {
		return InvalidSessionID, fmt.Errorf("Error when validating sessionID: %v", err)
	}

	return sessionID, nil
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	//TODO: get the SessionID from the request, and get the data
	//associated with that SessionID from the store.
	sessionID, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, fmt.Errorf("Error when getting sessionID: %v", err)
	}
	err = store.Get(sessionID, sessionState)
	if err != nil {
		return InvalidSessionID, err
	}
	return sessionID, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {

	sessionID, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, fmt.Errorf("Error when getting sessionID: %v", err)
	}
	err = store.Delete(sessionID)
	
	if err != nil {
		return InvalidSessionID, fmt.Errorf("Error deleting sessionID: %v", err)
	}
	return sessionID, nil
}
