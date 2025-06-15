package database

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type UUID[T any] uuid.UUID

func (u *UUID[T]) GormDataType() string {
	return "binary(16)"
}

func (u *UUID[T]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "binary"
}

func (u *UUID[T]) Scan(value any) (err error) {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal UUID value:", value))
	}
	parseBytes, err := uuid.FromBytes(bytes)
	*u = UUID[T](parseBytes)
	return
}

func (u UUID[T]) Value() (bytes driver.Value, err error) {
	bytes, err = uuid.UUID(u).MarshalBinary()
	return
}

func (u UUID[T]) String() string {
	return uuid.UUID(u).String()
}
