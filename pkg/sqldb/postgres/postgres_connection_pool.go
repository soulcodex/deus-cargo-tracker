package postgres

import (
	"database/sql"

	"github.com/soulcodex/deus-cargoes-tracker/pkg/sqldb"
)

type ConnectionPool struct {
	writer *sql.DB
	reader *sql.DB
}

func NewConnectionPoolFromCredentials(
	writerCredentials,
	readerCredentials Credentials,
) (*ConnectionPool, error) {
	writer, err := NewWriter(writerCredentials)
	if err != nil {
		return nil, err
	}

	reader, err := NewReader(readerCredentials)
	if err != nil {
		return nil, err
	}

	return NewConnectionPool(writer, reader)
}

func NewConnectionPool(writer, reader *sql.DB) (*ConnectionPool, error) {
	if writer == nil || reader == nil {
		return nil, sqldb.NewPoolConfigProvidedError(driver)
	}

	return &ConnectionPool{
		writer: writer,
		reader: reader,
	}, nil
}

func WithWriterOnly(writer *sql.DB) (*ConnectionPool, error) {
	if writer == nil {
		return nil, sqldb.NewPoolConfigProvidedError(driver)
	}

	return &ConnectionPool{
		writer: writer,
		reader: nil,
	}, nil
}

func (p *ConnectionPool) Writer() *sql.DB {
	return p.writer
}

func (p *ConnectionPool) Reader() *sql.DB {
	if nil == p.reader {
		return p.writer
	}

	return p.reader
}
