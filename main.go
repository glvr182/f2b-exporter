package main

import (
	"github.com/alicebob/sqlittle"
	log "github.com/sirupsen/logrus"
)

type prisoner struct {
	jail      string
	ip        string
	timeofban int
	bantime   int
}

func main() {
	log.Println("Starting exporter")
	db, err := sqlittle.Open("/var/lib/fail2ban/fail2ban.sqlite3")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	prisoners, err := jailed(db)
	if err != nil {
		panic(err)
	}

	log.Printf("%v", prisoners)
}

func jailed(db *sqlittle.DB) ([]prisoner, error) {
	var (
		prisoners []prisoner
		p         prisoner
	)

	db.Select("bans", func(r sqlittle.Row) {
		err := r.Scan(&p.jail, &p.ip, &p.timeofban, &p.bantime)
		if err != nil {
			panic(err)
		}
		prisoners = append(prisoners, p)

	}, "jail", "ip", "timeofban", "bantime")

	return prisoners, nil
}
