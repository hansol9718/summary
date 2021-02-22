package handlers
import (
	"time"
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-hansol9718/servers/gateway/models/users"
)

type SessionState struct {
	BeginTime time.Time   `json:"beginTime"`
	User      *users.User `json:"user"`
}