package sessions

import (
	"time"
	"github.com/go-redis/redis"
	"encoding/json"
	"fmt"
)

//RedisStore represents a session.Store backed by redis.
type RedisStore struct {
	//Redis client used to talk to redis server.
	Client *redis.Client
	//Used for key expiry time on redis.
	SessionDuration time.Duration
}

//NewRedisStore constructs a new RedisStore
func NewRedisStore(client *redis.Client, sessionDuration time.Duration) *RedisStore {
	//initialize and return a new RedisStore struct
	return &RedisStore{Client: client, SessionDuration: sessionDuration,}
}

//Store implementation

//Save saves the provided `sessionState` and associated SessionID to the store.
//The `sessionState` parameter is typically a pointer to a struct containing
//all the data you want to associated with the given SessionID.
func (rs *RedisStore) Save(sid SessionID, sessionState interface{}) error {
	//TODO: marshal the `sessionState` to JSON and save it in the redis database,
	//using `sid.getRedisKey()` for the key.
	//return any errors that occur along the way.
	j, err := json.Marshal(sessionState)
	if nil != err {
		return err
	}
	rs.Client.Set(sid.getRedisKey(), j, rs.SessionDuration)
	
	return nil
}

//Get populates `sessionState` with the data previously saved
//for the given SessionID
func (rs *RedisStore) Get(sid SessionID, sessionState interface{}) error {
	//TODO: get the previously-saved session state data from redis,
	//unmarshal it back into the `sessionState` parameter
	//and reset the expiry time, so that it doesn't get deleted until
	//the SessionDuration has elapsed.

	j := rs.Client.Get(sid.getRedisKey())
	if j.Err() != nil {
		return ErrStateNotFound
	}
	result, err := j.Result()
	if err != nil  {
		return fmt.Errorf("Error converting to bytes")
	}
	err = json.Unmarshal([]byte(result), sessionState)
	if err != nil {
		return fmt.Errorf("error unmarshalling: %v", err)
	}
	
	exp := rs.Client.Expire(sid.getRedisKey(), rs.SessionDuration)
	if exp.Err() != nil {
		return fmt.Errorf("error with expiration: %v", err)
	}

	return nil
}

//Delete deletes all state data associated with the SessionID from the store.
func (rs *RedisStore) Delete(sid SessionID) error {
	err := rs.Client.Del(sid.getRedisKey())
	if err.Err() != nil {
		return fmt.Errorf("Error deleting state data: %v", err.Err())
	}
	return nil
}

//getRedisKey() returns the redis key to use for the SessionID
func (sid SessionID) getRedisKey() string {
	//convert the SessionID to a string and add the prefix "sid:" to keep
	//SessionID keys separate from other keys that might end up in this
	//redis instance
	return "sid:" + sid.String()
}
