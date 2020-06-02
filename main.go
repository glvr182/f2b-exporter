package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/alicebob/sqlittle"
	"github.com/glvr182/f2b-exporter/provider"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

var (
	geocount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "f2b_banned_ip",
			Help: "Number of banned IPs per country / region",
		},
		[]string{"country", "geohash", "jail", "currently_banned"},
	)
)

func init() {
	prometheus.MustRegister(geocount)
}

func main() {
	log.Println("Starting exporter")
	go func() {
		for {
			time.Sleep(time.Minute)
			if err := update(); err != nil {
				log.Fatal(err)
			}
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func update() error {
	db, err := sqlittle.Open("/var/lib/fail2ban/fail2ban.sqlite3")
	if err != nil {
		return err
	}
	defer db.Close()

	provider, err := provider.New("freeGeoIP")
	if err != nil {
		return err
	}

	prisoners, err := jailed(db, provider)
	if err != nil {
		return err
	}

	geocount.Reset()
	for _, prisoner := range prisoners {
		geocount.With(prometheus.Labels{"country": prisoner.country, "geohash": prisoner.geohash, "jail": prisoner.jail, "currently_banned": strconv.FormatBool(prisoner.currentlyBanned)}).Inc()
	}

	return nil
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
