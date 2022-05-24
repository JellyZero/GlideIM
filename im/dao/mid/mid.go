package mid

import (
	"github.com/glide-im/glideim/pkg/db"
)

func GetMid() (int64, error) {
	k := "im:msg:mid:incr"
	seq, err := db.Redis.Incr(k).Result()
	if err != nil {
		return 0, err
	}
	return seq, nil
}
