package imdb

import (
	"github.com/BurntSushi/csql"
)

type Inserter struct {
	db *DB
	tx *Tx
	*csql.Inserter
}

func (db *DB) NewInserter(
	tx *Tx,
	size int,
	table string,
	columns ...string,
) (*Inserter, error) {
	ins := &Inserter{db, tx, nil}
	err := csql.Safe(func() {
		var err error
		if ins.tx == nil {
			ins.tx, err = db.Begin()
			csql.SQLPanic(err)
		}
		ins.Inserter, err = csql.NewInserter(ins.tx.Tx, db.Driver,
			size, table, columns...)
		csql.SQLPanic(err)
		db.inserters = append(db.inserters, ins)
	})
	return ins, err
}

func (ins *Inserter) Close() error {
	if ins.tx.closed {
		return nil
	}
	return ins.tx.Commit()
}
