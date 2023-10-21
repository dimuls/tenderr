package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
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

	f, err := excelize.OpenFile("docs/logs.xlsx")
	if err != nil {
		logger.Fatal("open logs xlsx", zap.Error(err))
		return
	}
	defer f.Close()

	rows, err := f.GetRows("Логи")
	if err != nil {
		logger.Fatal("get rows of logs sheet", zap.Error(err))
	}

	rows = rows[1:] // skip header

	f.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	clientErrorRegexpStrings := []string{
		"(?i)регистрац",
		"Unable to upload files",
		"Unable to get data",
		"Execution Timeout Expired.  The timeout period elapsed prior to completion of the operation or the server is not responding.",
		"Невозможно скачать файл по ссылке",
	}

	var elementIDs []string
	for i := 0; i < 10; i++ {
		var t time.Time

		determStr := t.Add(time.Duration(i) * time.Hour).String()
		md5Sum := md5.Sum([]byte(determStr))

		elementIDs = append(elementIDs, hex.EncodeToString(md5Sum[:]))
	}

	var clientErrorRegexpes []*regexp.Regexp
	for _, rx := range clientErrorRegexpStrings {
		clientErrorRegexpes = append(clientErrorRegexpes, regexp.MustCompile(rx))
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		time.Sleep(time.Duration(rand.Intn(100)+10) * time.Millisecond)

		i := rand.Intn(len(rows))
		row := rows[i]

		if i == 0 {
			// skip headers
			continue
		}

		log := entity.Log{
			Time:    time.Now(), // use now for testing
			ID:      row[0],
			Message: row[2],
		}

		var client bool
		for _, rx := range clientErrorRegexpes {
			if rx.MatchString(row[2]) {
				client = true
				break
			}
		}

		if client && rand.Float64() > 0.5 {
			i := rand.Intn(len(elementIDs))
			log.ElementID = elementIDs[i]
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
