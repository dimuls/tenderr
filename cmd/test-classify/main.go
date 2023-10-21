package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"

	"tednerr/entity"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Fatal(err.Error())
	}

	f, err := excelize.OpenFile("logs.xlsx")
	if err != nil {
		logger.Fatal("open logs xlsx", zap.Error(err))
		return
	}
	defer f.Close()

	rows, err := f.GetRows("Логи")
	if err != nil {
		logger.Fatal("get rows of logs sheet", zap.Error(err))
	}

	for i, row := range rows {
		if i == 0 {
			// skip headers
			continue
		}

		log := entity.Log{
			Time:    time.Now(), // use now for testing
			ID:      row[0],
			Message: row[2],
		}

		//log.Time, err = time.Parse(time.DateTime, row[1])
		//if err != nil {
		//	logger.Error("parse log time", zap.Error(err))
		//	continue
		//}

		logJSON, err := json.Marshal(log)
		if err != nil {
			logger.Error("json marshal log", zap.Error(err))
			continue
		}

		res, err := http.Post("http://localhost:8080/api/logs", "application/json", bytes.NewReader(logJSON))
		if err != nil {
			logger.Error("post log", zap.Error(err))
			continue
		}

		res.Body.Close()

		if res.StatusCode != http.StatusOK {
			logger.Error("post log: not ok status code", zap.Int("statusCode", res.StatusCode))
			continue
		}
	}
}
