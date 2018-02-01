package auth

type MapStore map[string]string

func (m MapStore) Set(key, value string) {
	m[key] = value
}

func (m MapStore) Get(key string, clb func(value string)) error {
	if v, ok := m[key]; ok {
		clb(v)
		return nil
	}
	return ErrNotFound
}
