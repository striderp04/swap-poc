package models

import (
	"fmt"
	"errors"
	"swapnil-ex/models/db"
	"time"
	"gorm.io/gorm"
	"gopkg.in/validator.v2"
	"strings"
)

type Student struct {
	ID            						uint    `json:"id"`
	Inil     									string `json:"inil"`
	FirstName     						string `json:"first_name" validate:"nonzero"`
	MiddleName     						string `json:"middle_name" validate:"nonzero"`
	LastName      						string `json:"last_name" validate:"nonzero"`
	BirthDate  								time.Time `json:"birth_date"`
	ParentName								string `json:"parent_name" validate:"nonzero"`
	ParentOccupation					string `json:"parent_occupation" validate:"nonzero"`
	ContactNumber 						string  `json:"contact_number" gorm:"contact_number" validate:"nonzero,min=10,max=12"`
	WhNumber									string  `json:"wh_number" validate:"nonzero,min=10,max=12"`
	Status 										string `json:"status"`
	Town 											string `json:"town" validate:"nonzero"`
	HasHostel									bool `json:"has_hostel" gorm:"default:false"`
	Balance 									float64 `json:"balance" gorm:"default:0.0"`
	BatchStandardStudents     []BatchStandardStudent 
	CreatedAt 								time.Time
	UpdatedAt 								time.Time
  DeletedAt 								gorm.DeletedAt `gorm:"index"`
}

func migrateStudent() {
	fmt.Println("migrating student..")
	err := db.Driver.AutoMigrate(&Student{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func NewStudent(studentData map[string]interface{}) *Student {
	student := &Student{}
	student.Assign(studentData)
	return student
}

func (s *Student) Validate() error {
	if errs := validator.Validate(s); errs != nil {
		return errs
	} else {
		return nil
	}
}

func (s *Student) AssignClass() error {	
	if s.Status == "Admission" {
		return nil
	} else {
		return errors.New("Already assigned Class")
	}
}

func (s *Student) Assign(studentData map[string]interface{}) {
	fmt.Printf("%+v\n", studentData)
	if firstName, ok := studentData["first_name"]; ok {
		s.FirstName = firstName.(string)
	}

	if middleName, ok := studentData["middle_name"]; ok {
		s.MiddleName = middleName.(string)
	}

	if lastName, ok := studentData["last_name"]; ok {
		s.LastName = lastName.(string)
	}

	if birthDate, ok := studentData["birth_date"]; ok {
		s.BirthDate, _ = time.Parse("2006-01-02T15:04:05.999999999Z", birthDate.(string))
	}	

	if parentName, ok := studentData["parent_name"]; ok {
		s.ParentName = parentName.(string)
	}
	if parentOccupation, ok := studentData["parent_occupation"]; ok {
		s.ParentOccupation = parentOccupation.(string)
	}
	if contactNumber, ok := studentData["contact_number"]; ok {
		s.ContactNumber = contactNumber.(string)
	}

	if whNumber, ok := studentData["wh_number"]; ok {
		s.WhNumber = whNumber.(string)
	}

	if town, ok := studentData["town"]; ok {
		s.Town = town.(string)
	}
}

func (s *Student) All(page int, search string) ([]Student, error) {
	var students []Student
	query := db.Driver.Limit(10).Offset((page - 1) * 10)
	search = strings.Trim(search, " ")
	if len([]rune(search)) > 0 {
		search = "%" + search + "%"
		query = query.Where("first_name like ? or middle_name like ? or last_name like ?", search, search, search)
	}
	err := query.Find(&students).Error
	return students, err
}

func (s *Student) Count(search string) (int64, error) {
	var count int64
	query := db.Driver.Model(&Student{})
	search = strings.Trim(search, " ")
	if len([]rune(search)) > 0 {
		search = "%" + search + "%"
		query = query.Where("first_name like ? or middle_name like ? or last_name like ?", search, search, search)
	}
	err := query.Count(&count).Error
	return count, err
}

func (s *Student) Find() error {
	err := db.Driver.First(s, "ID = ?", s.ID).Error
	return err
}

func (s *Student) Create() error {
	err := db.Driver.Create(s).Error
	return err
}

func (s *Student) Update() error {
	err := db.Driver.Save(s).Error
	return err
}

func (s *Student) Delete() error {
	err := db.Driver.Delete(s).Error
	return err
}

func (s *Student) AdmissionStatus() bool {
	return s.Status == "Admission"
}

func (s *Student) ConfirmedStatus() bool {
	return s.Status == "Confirmed"
}

func (s *Student) GetBatchStandardStudents() []BatchStandardStudent {
	var batchStandardStudents []BatchStandardStudent
	_ = db.Driver.Where("student_id = ?", s.ID).Find(&batchStandardStudents).Error
	return batchStandardStudents
}

func (s *Student) RemoveBatchStandard(batchStandard *BatchStandard) error {
	totalDebits, totalCredits := s.GetBalance()
	balance := totalCredits - totalDebits
	if balance > 0.0 {
		return errors.New("Please Clear Balance first")
	}

	batchStandardStudent := &BatchStandardStudent{StudentId: s.ID, BatchStandardId: batchStandard.ID}
	err := db.Driver.Find(batchStandardStudent).Error
	if err != nil {
		return err
	}
	return batchStandardStudent.Delete()
}

func (s *Student) AssignBatchStandard(batchStandard *BatchStandard) error {
	//check student already assign to batch standard
	//if assign remove from current batch standard before this check transaction 
	// after that assign new batch standard
	batchStandardStudents := s.GetBatchStandardStudents()
	if len(batchStandardStudents) > 0 {
		return errors.New("Already assigned to Class")
	} else {
		batchStandardStudent := &BatchStandardStudent{}
		
		batchStandardStudent.BatchId = batchStandard.BatchId
		batchStandardStudent.StandardId = batchStandard.StandardId
		batchStandardStudent.StudentId = s.ID
		batchStandardStudent.BatchStandardId = batchStandard.ID
		batchStandardStudent.Fee = batchStandard.Fee
		
		return batchStandardStudent.Create()
	}
	
} 

func (s *Student) AssignHostel(h *Hostel, hr *HostelRoom) error {
	var hostelStudent = HostelStudent{StudentId: s.ID, HostelId: h.ID, HostelRoomId: hr.ID}
	err := db.Driver.Where("student_id = ?",s.ID).First(&hostelStudent).Error
	
	if err != nil {
		err = hostelStudent.Create()
		if err == nil {
			s.HasHostel = true
			s.Update()
			s.SaveBalance()
		}
	}
	return err
}

func (s *Student) ChangeHostel(h *Hostel, hr *HostelRoom) error {
	var hostelStudent = HostelStudent{StudentId: s.ID}
	err := db.Driver.Where("student_id = ?",s.ID).First(&hostelStudent).Error
	if err == nil {
		hostelStudent.HostelId = h.ID
		hostelStudent.HostelRoomId = hr.ID
		s.HasHostel = true
		s.Update()
		return hostelStudent.Update()
	}
	return err
}

func (s *Student) GetStudentHostel() (HostelStudent, error) {
	var hostelStudent = HostelStudent{}
	err := db.Driver.Where("student_id = ?", s.ID).Preload("Hostel").Preload("HostelRoom").First(&hostelStudent).Error
	
	return hostelStudent, err
}

func (s *Student) GetTransactions() ([]Transaction, error) {
	transactions := []Transaction{}
	err := db.Driver.Where("student_id = ?", s.ID).Find(&transactions).Error
	return transactions, err
}

func (s *Student) TotalDebits() float64 {
	transactions, err := s.GetTransactions()
	var total = 0.0
	if err == nil {
		for _, transaction := range transactions {
			if transaction.TransactionType == "debit" {
				total = total + transaction.Amount
			}
		}
	}
	return total
}

func (s *Student) TotalCridits() float64 {
	transactions, err := s.GetTransactions()
	var total = 0.0
	if err == nil {
		for _, transaction := range transactions {
			if transaction.TransactionType == "credit" {
				total = total + transaction.Amount
			}
		}
	}
	return total
}

func (s *Student) SaveBalance() error{
	debits, credits := s.GetBalance()
	s.Balance =  credits - debits
	return s.Update()
}

func (s *Student) GetBalance() (float64, float64) {
	transactions, err := s.GetTransactions()
	var totalDebits = 0.0
	var totalCredits = 0.0
	if err == nil {
		for _, transaction := range transactions {
			if transaction.TransactionType == "debit" {
				totalDebits = totalDebits + transaction.Amount
			} else {
				totalCredits = totalCredits + transaction.Amount
			}
		}
	}
	return totalDebits, totalCredits
}
