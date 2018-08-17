package gorester

import (
	"strconv"
	"time"
	"errors"
)

type BaseFields struct {
	PkField
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
type PkField struct {
	ID int `gorm:"primary_key" json:"id"`
}

func (bf PkField) ListFields() []string {
	return []string{"*"}
}
func (bf PkField) InfoFields() []string {
	return bf.ListFields()
}
func (bf PkField) ForbidUpdateFields() []string {
	return SetForbidUpdateFields()
}
func SetForbidUpdateFields(fs ...string) []string {
	res := []string{"id", "created_at", "updated_at", "deleted_at"}
	for _, value := range fs {
		res = append(res, value)
	}
	return res
}

type ResourceInterface interface {
	ListFields() []string
	InfoFields() []string
	ForbidUpdateFields() []string
}

type ForbidUpdateResource struct{}

func (bf ForbidUpdateResource) ForbidUpdate() bool {
	return true
}

func FindOneByID(res ResourceInterface, id int) ResourceInterface {

	if err := Db.Where("id = ?", id).Last(res).Error; err == nil {
		return res
	} else {
		panic(errors.New("ByID:(" + strconv.Itoa(id) + ") data not found "))
	}
}

func FindOneByMap(res ResourceInterface, where map[string]interface{}) ResourceInterface {
	if err := Db.Where(where).First(res).Error; err == nil {
		return res
	} else {
		panic(QueryDaoError(err))
	}
}

func FindListByMap(res interface{}, where map[string]interface{}, order string, page int, perPage int) {
	offset := perPage * (page - 1)
	if err := Db.Where(where).Order(order).Offset(offset).Limit(perPage).Find(res).Error; err != nil {
		panic(QueryDaoError(err))
	}
}

func UpdateByID(id int, res ResourceInterface) ResourceInterface {
	if err := Db.Model(res).Where("id = ?", id).Updates(res).Error; err == nil {
		return FindOneByID(res, id)
	} else {
		panic(QueryDaoError(err))
	}
}

func UpdateWhere(where map[string]interface{}, res ResourceInterface) ResourceInterface {
	if err := Db.Model(res).Where(where).Updates(res).Error; err == nil {
		if val, ok := where["id"]; ok {
			return FindOneByID(res, val.(int))
		}
		return res
	} else {
		panic(QueryDaoError(err))
	}
}

func DeleteByID(res ResourceInterface, id int) ResourceInterface {
	if err := Db.Where("id = ?", id).Delete(res).Error; err == nil {
		return res
	} else {
		panic(QueryDaoError(err))
	}
}

func Create(res ResourceInterface) ResourceInterface {
	if err := Db.Create(res).Error; err == nil {
		return res
	} else {
		panic(QueryDaoError(err))
	}
}

func ExsitAndFirst(res ResourceInterface) {
	if err := Db.Where(res).First(res).Error; err != nil {
		res = nil
	}
}
