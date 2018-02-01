package storage

import (
	"github.com/aspcartman/exceptions"
	"github.com/aspcartman/eros/env"
)

type oldUser struct {
	ID              int
	SID             string
	Name            string
	Age             *int
	School          *string
	Job             *string
	Info            *string
	Photos          []string
	Deleted         bool
	Matched         bool
	InstagramNick   *string
	InstagramPhotos []string
}

func GetOldUsers() (<-chan oldUser, int) {
	db := NewDB("eros.ru", 7770, "postgres", "", "himeros")
	tx := db.tx()

	count := 0
	err := tx.QueryRow(`select count(1) from "user"`).Scan(&count)
	if err != nil {
		e.Throw("failed getting user count from the old db", err)
	}

	env.Log.WithField("count", count).Info("Found some old users to fetch")

	rows, err := tx.Query(`select id, sid, name, age, school, job, info, photos, deleted, matched, instagram_nick, instagram_photos from "user" order by id`)
	if err != nil {
		e.Throw("failed getting users from the old db", err)
	}

	ch := make(chan oldUser, 100)
	go func() {
		defer tx.Rollback()
		defer close(ch)

		env.Log.Info("Reading old users")
		for rows.Next() {
			usr := oldUser{}
			err := rows.Scan(&usr.ID, &usr.SID, &usr.Name, &usr.Age, &usr.School, &usr.Job, &usr.Info, &usr.Photos, &usr.Deleted, &usr.Matched, &usr.InstagramNick, &usr.InstagramPhotos)
			if err != nil {
				e.Throw("failed scanning an old user", err)
			}
			ch <- usr
		}
		env.Log.Info("All old users had been read and sent to the channel")
	}()

	return ch, count
}
