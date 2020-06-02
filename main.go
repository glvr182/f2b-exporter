package main

import (
	"log"
	"time"

	"github.com/alicebob/sqlittle"
	"github.com/glvr182/f2b-exporter/provider"
)

// prisoner is an ip that has been (temp) banned by f2b
type prisoner struct {
	// jail represents the jail f2b assigned
	jail string
	// ip is the prisoners ip
	ip string
	// timeofban indicates moment he was banned
	timeofban int
	// bantime indicates how long he will be banned (-1 = infinity)
	bantime int
	// country is the general location of the prisoner
	country string
	// geohash is a more accurate location of the prisoner
	geohash string
	// currentlyBanned indicates if the prisoner is currently banned
	currentlyBanned bool
}

func main() {
	log.Println("Starting exporter")
	db, err := sqlittle.Open("/var/lib/fail2ban/fail2ban.sqlite3")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	provider, err := provider.New("freeGeoIP")
	if err != nil {
		panic(err)
	}

	prisoners, err := jailed(db, provider)
	if err != nil {
		panic(err)
	}

	log.Printf("%v", prisoners)
}

func jailed(db *sqlittle.DB, provider provider.Provider) ([]prisoner, error) {
	var (
		prisoners []prisoner
		p         prisoner
		err       error
	)

	db.SelectDone("bans", func(r sqlittle.Row) bool {
		err = r.Scan(&p.jail, &p.ip, &p.timeofban, &p.bantime)
		if err != nil {
			return true
		}

		payload, err := provider.Lookup(p.ip)
		if err != nil {
			return true
		}
		p.country = payload.CountryCode
		p.geohash = payload.GeoHash

		if int64(p.timeofban+p.bantime) > time.Now().Unix() || p.bantime < 0 {
			p.currentlyBanned = true
		}

		prisoners = append(prisoners, p)
		return false

	}, "jail", "ip", "timeofban", "bantime")

	if err != nil {
		return nil, err
	}

	return prisoners, nil
}
