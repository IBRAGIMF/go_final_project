package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	DateFormat = "20060102" //valid formtat
	MaxDigit   = 400
	MsgErr     = "Error detected"
	MsgOk      = "OK"
	WebDir     = "web"
	TxtErr     = "Ошибка при добавлении задания в БД"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	t, _ := time.Parse(DateFormat, date)
	rep := strings.ReplaceAll(repeat, " ", "")

	if len(rep) > 0 && rep[0] == 'y' {
		origMonth := t.Month()
		origDay := t.Day()

		for {
			year := t.Year() + 1
			month := origMonth
			day := origDay

			// обработка для невисокосному году
			if origMonth == time.February && origDay == 29 {
				if year%4 != 0 || (year%100 == 0 && year%400 != 0) {
					month = time.March
					day = 1
				}
			}

			candidate := time.Date(year, month, day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
			if candidate.Before(now) {
				t = candidate
				continue
			}
			return candidate.Format(DateFormat), nil
		}
	} else {
		num, err := strconv.Atoi(rep[1:])
		if err != nil {
			return MsgErr, err
		}

		next := t.AddDate(0, 0, num)
		for next.Before(now) {
			next = next.AddDate(0, 0, num)
		}
		return next.Format(DateFormat), nil
	}
}
func NextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowVal := r.FormValue("now")
	dateVal := r.FormValue("date")
	repeatVal := r.FormValue("repeat")

	if CheckNextDate(nowVal, dateVal, repeatVal) != MsgOk {
		http.Error(w, MsgErr, http.StatusBadRequest)
		return
	}
	timeNow, err := time.Parse(DateFormat, nowVal)
	if err != nil {
		http.Error(w, MsgErr, http.StatusBadRequest)
		return
	}

	num, err := NextDate(timeNow, dateVal, repeatVal)
	if err != nil {
		http.Error(w, MsgErr, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(num))
}

func CheckNextDate(now string, date string, repeat string) string {
	repeat = strings.ReplaceAll(repeat, " ", "")

	if len(date) == 0 || len(repeat) == 0 {
		return MsgErr
	}

	if _, err := time.Parse(DateFormat, now); err != nil {
		return MsgErr
	}
	_, err := time.Parse(DateFormat, date)
	if err != nil {
		return MsgErr
	}

	if repeat == "y" && len(repeat) == 1 {
		return MsgOk
	} else if repeat[0:1] == "d" && len(repeat) == 1 {
		return MsgErr
	} else if repeat[0:1] == "y" && len(repeat) > 1 {
		return MsgErr
	}
	num, err := strconv.Atoi(repeat[1:])
	if err != nil {
		return MsgErr
	}

	if repeat[0:1] == "d" && num <= MaxDigit {
		return MsgOk
	} else {
		return MsgErr
	}
}
