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
}

func createDatabase() {
	log.Println("Creating database...")

	sqlStmt1 := `CREATE TABLE IF NOT EXISTS kudos (id integer primary key, userFrom text, kudo text, message text, value INTEGER DEFAULT 1, date text);`
	_, err1 := db.Exec(sqlStmt1)
	checkErr(err1)

	sqlStmt2 := `CREATE TABLE IF NOT EXISTS kudos_receiver (id integer primary key, userTo text, kudos_id INTEGER);`
	_, err2 := db.Exec(sqlStmt2)
	checkErr(err2)
}

// Remove data migration in the future.
func migrateKudos() {
	stmt1 := `ALTER TABLE kudos ADD COLUMN message text`
	db.Exec(stmt1)

	rows, err := db.Query("SELECT id, kudo FROM kudos ORDER BY id ASC")
	checkErr(err)
	defer rows.Close()

	stmt2, err2 := db.Prepare("INSERT INTO kudos_receiver (userTo, kudos_id) VALUES (?, ?)")
	checkErr(err2)
	stmt3, err3 := db.Prepare("UPDATE kudos SET message = ? WHERE id = ?")
	checkErr(err3)

  var kudosList []Kudo
	for rows.Next() {
		var kudo Kudo
		err = rows.Scan(&kudo.ID, &kudo.Text)
		checkErr(err)
		kudosList = append(kudosList, kudo)
	}
	rows.Close()

	for _, element := range kudosList {
		parsed := parseKudoCommand(element.Text)
		log.Printf("Kudo #%d\n", element.ID)
		stmt3.Exec(parsed.Text, element.ID)
		for _, member := range parsed.Members {
			log.Printf("Receiver #%s\n", member.ID)
			stmt2.Exec(member.ID, element.ID)
		}
		log.Println("________")
	}
}

func dbSaveKudo(kudo *Kudo) {
	if kudo.ID != 0 {
		return
	}

	stmt, err := db.Prepare("INSERT INTO kudos (userFrom, kudo, message, value, date) VALUES (?, ?, ?, ?, ?)")
	checkErr(err)

	res, err := stmt.Exec(kudo.MemberFrom.ID, kudo.Original, kudo.Text, kudo.Value, kudo.Date.Format(time.RFC3339))
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	stmt2, err := db.Prepare("INSERT INTO kudos_receiver (userTo, kudos_id) VALUES (?, ?)")
	checkErr(err)

	for _, member := range kudo.Recipients {
		stmt2.Exec(member.ID, id)
	}

	kudo.ID = id
}

func loadKudos() {
	rows, err := db.Query("SELECT id, userFrom, message, value, date FROM kudos ORDER BY id DESC LIMIT 9")
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var kudo Kudo
		var memberFrom, date string
		err = rows.Scan(&kudo.ID, &memberFrom, &kudo.Text, &kudo.Value, &date)
		checkErr(err)

		kudo.MemberFrom, _ = findMemberByID(memberFrom)
		kudo.Recipients = loadKudosRecipients(kudo.ID)
		kudo.Color = randomColor()
		kudo.Date, _ = time.Parse(time.RFC3339, date)
		kudos = append(kudos, kudo)
	}
}

func loadKudosRecipients(kudoId int64) []Member {
		stmt, err := db.Prepare("SELECT userTo FROM kudos_receiver WHERE kudos_id = ?")
		checkErr(err)
		memberRows, err := stmt.Query(kudoId)
		checkErr(err)
		defer memberRows.Close()

		var members []Member
		for memberRows.Next() {
			var memberToId string
			err = memberRows.Scan(&memberToId)
			checkErr(err)

			memberTo, err := findMemberByID(memberToId)
			if err != nil {
				log.Printf("Invalid user provided: #%s\n", memberToId)
			} else {
				members = append(members, memberTo)
			}
		}
		return members
}

func loadKudosGaveList() []KudosStats {

	rows, err := db.Query(`
			SELECT userFrom, count(id) as pts FROM kudos
			WHERE kudos.date > date('now','-1 year')
			GROUP by userFrom HAVING count(id) > 0  ORDER BY pts DESC
			`)
	checkErr(err)
	defer rows.Close()

		var statsList []KudosStats
		i := 1
		var max float32 = 0
		for rows.Next() {
			var memberFrom string
			var pts int
			err = rows.Scan(&memberFrom, &pts)
			checkErr(err)
			if i == 1 {
				max = float32(pts)
			}
			statsList = append(statsList, generateStatsResult(memberFrom, pts, i, max))
			i = i + 1
		}
		return statsList
}

func loadKudosReceivedList() []KudosStats {

	rows, err := db.Query(`
			SELECT userTo, SUM(value) as pts FROM kudos ku
			INNER JOIN kudos_receiver kr ON kr.kudos_id = ku.id
			WHERE ku.date > date('now','-1 year')
			GROUP by userTo HAVING SUM(value) > 0 ORDER BY pts DESC
			`)
	checkErr(err)
	defer rows.Close()

	var statsList []KudosStats
	i := 1
	var max float32 = 0
	for rows.Next() {
		var memberFrom string
		var pts int
		err = rows.Scan(&memberFrom, &pts)
		checkErr(err)
		if i == 1 {
			max = float32(pts)
		}
		statsList = append(statsList, generateStatsResult(memberFrom, pts, i, max))
		i = i + 1
	}
	return statsList
}

func loadKudosReceivedByUser(userId string) []Kudo {
	stmt, err := db.Prepare(`SELECT ku.id, ku.userFrom, ku.message, ku.value, ku.date
			FROM kudos ku
			JOIN kudos_receiver ON ku.id = kudos_receiver.kudos_id
			WHERE userTo = ?
			ORDER BY ku.date DESC
	`)

	checkErr(err)
	rows, err := stmt.Query(userId)
	defer rows.Close()

	var kudos []Kudo

	for rows.Next() {
		var kudo Kudo
		var memberFrom, date string
		err = rows.Scan(&kudo.ID, &memberFrom, &kudo.Text, &kudo.Value, &date)
		kudo.MemberFrom, _ = findMemberByID(memberFrom)
		kudo.Date, _ = time.Parse(time.RFC3339, date)
		kudos = append(kudos, kudo)
	}

	return kudos
}

func generateStatsResult(memberFrom string, pts int, i int, max float32) KudosStats {

	prc := float32(pts * 85) / max
	if prc < 0 {
		prc = 0
	}
	var kudoStats KudosStats
	kudoStats.Member, _ = findMemberByID(memberFrom)
	kudoStats.Pts = pts
	kudoStats.Position = i
	kudoStats.HasCrown = i == 1
	kudoStats.Prc = prc

	return kudoStats
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
