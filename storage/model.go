package storage

type User struct {
	ID        string
	Name      string
	Age       int
	School    string
	Job       string
	Info      string
	Instagram string
	Deleted   bool
	Matched   bool
}

func (tx *Tx) SaveUser(u *User) {
	tx.Exec(`
		INSERT INTO users (id, name, age, school, job, info, instagram, deleted, matched)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (id) DO UPDATE SET (name, age, school, job, info, instagram, deleted, matched) = ($2,$3,$4,$5,$6,$7,$8,$9)
	`, u.ID, u.Name, u.Age, u.School, u.Job, u.Info, u.Instagram, u.Deleted, u.Matched)
}

func (tx *Tx) SavePhoto(p *Photo) {
	if p.ID == 0 {
		tx.Exec(`INSERT INTO photos (user_id, source, url, original, small, placeholder)
				 VALUES ($1,$2,$3,$4,$5,$6)`,
			p.UserID, p.Source, p.URL, p.Original, p.Small, p.Placeholder)
	} else {
		tx.Exec(`UPDATE photos
				 SET	(user_id, source, url, original, small, placeholder) = ($2,$3,$4,$5,$6)
				 WHERE 	id = $1`,
			p.ID, p.UserID, p.Source, p.URL, p.Original, p.Small, p.Placeholder)
	}

}

type Photo struct {
	ID          uint
	UserID      string
	Source      string
	URL         string
	Original    []byte
	Small       []byte
	Placeholder []byte
}
