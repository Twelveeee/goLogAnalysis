package entity

import (
	"fmt"
	"strings"
	"time"
)

type OssLog struct {
	ID          int64     `gorm:"type:int(64) UNSIGNED;primary_key;auto_increment:true;not null"`
	RequestTime time.Time `gorm:"type:timestamp;not null;default:'1970-01-01 00:00:01';index:idx_requestTime"`
	IP          uint32    `gorm:"type:int(64) UNSIGNED;not null;default:0;index:idx_ip_referer,priority:1"`
	UA          string    `gorm:"type:varchar(1024);not null;default:''"`
	Path        string    `gorm:"type:varchar(1024);not null;default:''"`
	Referer     string    `gorm:"type:varchar(1024);not null;default:'';index:idx_ip_referer,priority:2,length:20"`
	Host        string    `gorm:"type:varchar(64);not null;default:''"`
	Bucket      string    `gorm:"type:varchar(64);not null;default:'';"`
}

// TableName returns the entity table name.
func (m *OssLog) TableName() string {
	if m.RequestTime.IsZero() {
		return "osslog"
	}
	return "osslog_" + m.RequestTime.Format("200601")
}

// Save updates the record in the database or inserts a new record if it does not already exist.
func (m *OssLog) Save() error {
	return Db().Table(m.TableName()).Save(m).Error
}

// Create inserts the entity to the database.
func (m *OssLog) Create() error {
	if err := m.CreateTable(); err != nil {
		return err
	}

	return Db().Table(m.TableName()).Create(m).Error
}

func (m *OssLog) CreateTable() error {
	if HasTableCache(m.TableName()) {
		return nil
	}

	CreateTableLock()
	defer CreateTableUnlock()

	if !Db().Table(m.TableName()).Migrator().HasTable(m) {
		if err := Db().Table(m.TableName()).Migrator().CreateTable(m); err != nil {
			return err
		}
	}
	SetHasTableCache(m.TableName())
	return nil
}

func (m *OssLog) CreateInBatches(mList []OssLog, batchSize int) (int64, error) {
	mMap := make(map[string][]OssLog)
	for _, v := range mList {
		// if v.IsFilter() {
		// 	continue
		// }

		mMap[v.TableName()] = append(mMap[v.TableName()], v)
	}

	var RowsAffected int64 = 0

	for _, ossList := range mMap {
		if len(ossList) == 0 {
			continue
		}

		if err := ossList[0].CreateTable(); err != nil {
			log.Fatal().Msg(err.Error())
			continue
		}

		result := Db().Table(ossList[0].TableName()).CreateInBatches(ossList, batchSize)
		log.Debug().Msg(fmt.Sprintf("CreateInBatches: %s, RowsAffected: %d", ossList[0].TableName(), result.RowsAffected))
		RowsAffected += result.RowsAffected
	}

	return RowsAffected, nil
}

func (m *OssLog) IsFilter() (string, bool) {
	if m.IP == 0 {
		return "empty ip", true
	}

	if m.Referer == "-" || len(m.Referer) == 0 {
		return "empty referer", true
	}

	if m.RequestTime.IsZero() {
		return "empty requestTime", true
	}

	if strings.Contains(m.UA, "aliyun-sdk") {
		return "aliyun-sdk", true
	}

	return "", false
}
