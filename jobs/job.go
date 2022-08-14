package jobs

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Boozoorg/GreatProjeck/client"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func StartJob(jobParam int64, dsn string) {
	timer := time.NewTicker(time.Hour * time.Duration(jobParam))
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			usefulCode(jobParam, dsn)
		}
	}
}

func usefulCode(jobParam int64, dsn string) {
	file, err := os.Open("jobs/text/log.txt")
	if err != nil {
		log.Print(err)
		return
	}
	defer func() {
		cerr := file.Close()
		if cerr != nil {
			log.Print(cerr)
		}
	}()

	file, err = os.Create("jobs/text/log.txt")
	if err != nil {
		log.Print(err)
		return
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Print(err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Print(err)
		}
	}()
	rows, err := db.Query(`SELECT * FROM messanger`)
	if err != nil {
		log.Print(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item = &client.Chat{}
		err = rows.Scan(&item.SendlerID, &item.ReceiverID, &item.Message, &item.Time)
		if err != nil {
			log.Print(err)
			return
		}
		id1 := strconv.FormatUint(item.SendlerID, 10)
		id2 := strconv.FormatUint(item.ReceiverID, 10)
		i := id1 + " " + id2 + " " + "[" + item.Message + "]" + " " + string(item.Time.String()) + "\n"
		file.Write([]byte(i))
	}
}
