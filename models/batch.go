package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
)

type Batch struct {
	ID            	uint    `json:"id"`
	Name     				string `json:"name"`
	Year      			int `json:"year"`
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
  DeletedAt 			gorm.DeletedAt `gorm:"index"`
}

func migrateBatch() {
	fmt.Println("migrating student..")
	err := db.Driver.AutoMigrate(&Batch{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewBatch(batchData map[string]interface{}) *Batch {
	batch := &Batch{}
	batch.Assign(batchData)
	return batch
}

func (b *Batch) Validate() error {
	return nil
}

func (b *Batch) Assign(batchData map[string]interface{}) {
	fmt.Printf("%+v\n", batchData)
	if name, ok := batchData["name"]; ok {
		b.Name = name.(string)
	}

	if year, ok := batchData["year"]; ok {
		b.Year = int(year.(int64))
	}
}

func (b *Batch) All() ([]Batch, error) {
	var batchs []Batch
	err := db.Driver.Find(&batchs).Error
	return batchs, err
}

func (b *Batch) Find() error {
	err := db.Driver.First(b, "ID = ?", b.ID).Error
	return err
}

func (b *Batch) Create() error {
	//transaction block
	err := db.Driver.Create(b).Error
	s := &Standard{}
	stds, stdErr := s.All()
	if stdErr == nil {
		for _, std := range stds {
			bs := &BatchStandard{}
			bs.BatchId = b.ID
			bs.StandardId = std.ID
			bsErr := bs.Create()
			if bsErr != nil {
				// uncommit 
				break;
			}
		}
	} else {
		// uncommit 
		// return block 
	}
	db.Commit()
	return err
}

func (b *Batch) Update() error {
	err := db.Driver.Save(b).Error
	db.Commit()
	return err
}

func (b *Batch) Delete() error {
	err := db.Driver.Delete(b).Error
	_ := db.Commit()
	return err
}
