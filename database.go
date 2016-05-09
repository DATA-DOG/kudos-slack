package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func loadDatabase() {
	conn, err := sql.Open("sqlite3", config.Database)
	checkErr(err)

	db = conn

	createDatabase()

	log.Println("Loading database...")

	loadKudos()
}

func createDatabase() {
	log.Println("Creating database...")

	sqlStmt := `CREATE TABLE IF NOT EXISTS kudos (id integer primary key, userFrom text, kudo text, value INTEGER DEFAULT 1, date text);`
	_, err := db.Exec(sqlStmt)
	checkErr(err)
}

func dbSaveKudo(kudo *Kudo) {
	if kudo.ID != 0 {
		return
	}

	stmt, err := db.Prepare("INSERT INTO kudos (userFrom, kudo, value, date) VALUES (?, ?, ?, ?)")
	checkErr(err)

	res, err := stmt.Exec(kudo.MemberFrom.ID, kudo.Original, kudo.Value, kudo.Date.Format(time.RFC3339))
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	kudo.ID = id
}

func loadKudos() {
	rows, err := db.Query("SELECT id, userFrom, kudo, value, date FROM kudos ORDER BY id DESC LIMIT 9")
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var kudo Kudo
		var memberFrom, message, date string
		err = rows.Scan(&kudo.ID, &memberFrom, &message, &kudo.Value, &date)
		checkErr(err)

		parsed := parseKudoCommand(message)
		kudo.MemberFrom, _ = findMemberByID(memberFrom)
		kudo.Text = parsed.Text
		kudo.Recipients = parsed.Members
		kudo.Color = randomColor()
		kudo.Date, _ = time.Parse(time.RFC3339, date)

		kudos = append(kudos, kudo)
	}
}

func loadKudosGaveList() []KudosStats {

	rows, err := db.Query(`
			SELECT userFrom, count(id) as pts FROM kudos
			GROUP by userFrom ORDER BY pts DESC
			LIMIT 5`)
	checkErr(err)
	defer rows.Close()

	var statsList []KudosStats
	i := 1
	var max float32 = 0
	for rows.Next() {
		var kudoStats KudosStats
		var memberFrom string
		var pts int
		err = rows.Scan(&memberFrom, &pts)
		checkErr(err)

		if i == 1 {
			max = float32(pts)
		}
		kudoStats.Member, _ = findMemberByID(memberFrom)
		kudoStats.Pts = pts
		kudoStats.Position = i
		kudoStats.HasCrown = i == 1
		kudoStats.Prc = float32(pts * 85) / max

		i = i + 1
		statsList = append(statsList, kudoStats)
	}
	return statsList
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
