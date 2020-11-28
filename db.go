package mydis

type DB struct {
	Conn
}

func OpenDB(c Conn, err error) (*DB, error) {
	if err != nil {
		return nil, err
	}
	return &DB{Conn: c}, nil
}
