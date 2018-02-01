package storage

import "github.com/aspcartman/eros/env"

var revisions = [...]func(tx *Tx){
	0: func(tx *Tx) {
		tx.Exec(`CREATE TABLE "public"."app" (
		    "key" TEXT,
		    "value" INT NOT NULL,
		    PRIMARY KEY ("key")
		)`)

		tx.Exec(`INSERT INTO app VALUES ('revision','0')`)
	},
	1: func(tx *Tx) {
		tx.Exec(`CREATE TABLE "public"."users" (
		    "id" TEXT NOT NULL,
		    "name" TEXT NOT NULL,
		  	"age" INTEGER NOT NULL,
		  	"school" TEXT NOT NULL,
		  	"job" TEXT NOT NULL,
		  	"info" TEXT NOT NULL,
			"instagram" TEXT NOT NULL,
			"deleted" BOOL NOT NULL,
			"matched" BOOL NOT NULL,

		    PRIMARY KEY ("id")
		)`)
	},
	2: func(tx *Tx) {
		tx.Exec(`CREATE TABLE "public"."photos" (
			    "id" serial,
			    "user_id" text NOT NULL,
			    "source" text NOT NULL,
			    "url" text NOT NULL,
			    "original" bytea,
			    "small" bytea,
			    "placeholder" bytea,
			    PRIMARY KEY ("id"),
			    FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE ON UPDATE CASCADE
		)`)
	},
	3: func(tx *Tx) {
		ch, count := GetOldUsers()
		env.Log.WithField("count", count).Info("Will migrate users")
		for u := range ch {
			tx.SaveUser(&User{
				ID:        u.SID,
				Name:      u.Name,
				Age:       ptoi(u.Age),
				School:    ptos(u.School),
				Job:       ptos(u.Job),
				Info:      ptos(u.Info),
				Deleted:   u.Deleted,
				Matched:   u.Matched,
				Instagram: ptos(u.InstagramNick),
			})
			for _, p := range u.Photos {
				tx.SavePhoto(&Photo{
					UserID: u.SID,
					Source: "tinder",
					URL:    p,
				})
			}
		}
	},
}

func ptos(s *string) string {
	if s == nil {
		return ""
	} else {
		return *s
	}
}

func ptoi(i *int) int {
	if i == nil {
		return 0
	} else {
		return *i
	}
}
