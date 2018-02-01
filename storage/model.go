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

