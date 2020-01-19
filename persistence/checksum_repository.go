package persistence

import (
	"github.com/astaxie/beego/orm"
	"github.com/cloudsonic/sonic-server/model"
)

type checkSumRepository struct {
	ormer orm.Ormer
}

const checkSumId = "1"

type checksum struct {
	ID  string `orm:"pk;column(id)"`
	Sum string
}

func NewCheckSumRepository(o orm.Ormer) model.ChecksumRepository {
	r := &checkSumRepository{ormer: o}
	return r
}

func (r *checkSumRepository) GetData() (model.ChecksumMap, error) {
	loadedData := make(map[string]string)

	var all []checksum
	_, err := r.ormer.QueryTable(&checksum{}).Limit(-1).All(&all)
	if err != nil {
		return nil, err
	}

	for _, cks := range all {
		loadedData[cks.ID] = cks.Sum
	}

	return loadedData, nil
}

func (r *checkSumRepository) SetData(newSums model.ChecksumMap) error {
	_, err := r.ormer.Raw("delete from checksum").Exec()
	if err != nil {
		return err
	}

	var checksums []checksum
	for k, v := range newSums {
		cks := checksum{ID: k, Sum: v}
		checksums = append(checksums, cks)
	}
	_, err = r.ormer.InsertMulti(batchSize, &checksums)
	if err != nil {
		return err
	}
	return nil
}

var _ model.ChecksumRepository = (*checkSumRepository)(nil)
