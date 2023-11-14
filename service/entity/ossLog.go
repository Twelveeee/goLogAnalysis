package entity

import "strings"

type OssLog struct {
	ID          int64  `gorm:"type:int(64) UNSIGNED;primary_key;auto_increment:true;not null"`
	RequestTime int64  `gorm:"type:int(64) UNSIGNED;not null;default:0;index:idx_requesttime"`
	IP          uint32 `gorm:"type:int(64) UNSIGNED;not null;default:0;index:idx_ip_referer,priority:1"`
	UA          string `gorm:"type:varchar(1024);not null;default:''"`
	Path        string `gorm:"type:varchar(1024);not null;default:''"`
	Referer     string `gorm:"type:varchar(1024);not null;default:'';index:idx_ip_referer,priority:2,length:20"`
	Host        string `gorm:"type:varchar(64);not null;default:''"`
	Bucket      string `gorm:"type:varchar(64);not null;default:'';"`
}

// TableName returns the entity table name.
func (OssLog) TableName() string {
	return "osslog"
}

// Save updates the record in the database or inserts a new record if it does not already exist.
func (m *OssLog) Save() error {
	return Db().Save(m).Error
}

// Create inserts the entity to the database.
func (m *OssLog) Create() error {
	return Db().Create(m).Error
}

func (m *OssLog) CreateTable() error {
	if !Db().Migrator().HasTable(m) {
		return Db().Migrator().CreateTable(m)
	}
	return nil
}

func (m *OssLog) CreateInBatches(mList []OssLog, batchSize int) (int64, error) {
	result := Db().CreateInBatches(mList, batchSize)
	return result.RowsAffected, result.Error
}

func (m *OssLog) IsFilter() bool {
	if m.IP == 0 {
		return true
	}

	if m.Referer == "-" || len(m.Referer) == 0 {
		return true
	}

	if m.RequestTime == 0 {
		return true
	}

	if strings.Contains(m.UA, "aliyun-sdk") {
		return true
	}

	return false
}
