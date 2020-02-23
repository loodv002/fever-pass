package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Record struct {
	ID     uint32 `gorm:"primary_key"`
	UserID string
	Pass   bool
	Time   time.Time

	RecordedBy Account `json:"-"`
	AccountID  uint32  `json:"-"`
}

// insert a new record into database
func (h Handler) newRecord(w http.ResponseWriter, r *http.Request) {
	var record Record
	var err error
	record.UserID = r.FormValue("user_id")

	record.Pass, err = parseBool(r.FormValue("pass")) // default false
	if err != nil {
		http.Error(w, err.Error(), 415)
		return
	}

	record.Time = time.Now()

	if acct, ok := r.Context().Value(KeyAccount).(Account); ok {
		record.AccountID = acct.ID
	} else {
		http.Error(w, "cannot read account from session", 500)
		return
	}

	err = h.db.Create(&record).Error
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(w)
	enc.Encode(record)
}

func parseBool(str string) (bool, error) {
	switch str {
	case "on", "true":
		return true, nil
	case "off", "false", "":
		return false, nil
	default:
		return false, fmt.Errorf("'%s' cannot parse to bool", str)
	}
}

func (h Handler) findRecord(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	var record Record
	err := h.db.Where("user_id = ? and time > ?", userID, today()).Order("time desc").First(&record).Error
	if gorm.IsRecordNotFoundError(err) {
		http.Error(w, "record not found", 404)
		return
	} else if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	enc := json.NewEncoder(w)
	if err = enc.Encode(&record); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
}

func (h Handler) listRecord(w http.ResponseWriter, r *http.Request) {
	var records []Record
	err := h.db.Where("time > ?", today()).Order("time desc").Find(&records).Error

	enc := json.NewEncoder(w)
	if err = enc.Encode(&records); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h Handler) deleteRecord(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), 415)
		return
	}

	err = h.db.Delete(&Record{}, "id = ?", id).Error
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h Handler) editPage(w http.ResponseWriter, r *http.Request) {
	var records []Record
	if acct, ok := r.Context().Value(KeyAccount).(Account); ok {
		err := h.db.Where("account_id = ?", acct.ID).Order("time desc").Limit(20).Find(&records).Error
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		http.Error(w, "cannot read account from session", 500)
		return
	}
	page := struct {
		Records []Record
	}{records}
	h.tpl.ExecuteTemplate(w, "new.htm", page)
}
