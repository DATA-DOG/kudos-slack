package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func loadDatabase() {
	conn, err := sql.Open("sqlite3", config.Database)
	checkErr(err)

	db = conn

	createDatabase()

	fmt.Println("Loading database...")

	loadKudos()
}

func createDatabase() {
	fmt.Println("Creating database...")

	sqlStmt := `CREATE TABLE IF NOT EXISTS kudos (id integer primary key, user text, userFrom text, kudo text, likes int);`
	_, err := db.Exec(sqlStmt)
	checkErr(err)

	sqlStmt = `CREATE TABLE IF NOT EXISTS likes (user text, kudo int);`
	_, err = db.Exec(sqlStmt)
	checkErr(err)

	sqlStmt = `ALTER TABLE kudos ADD COLUMN value INTEGER DEFAULT 1`
	_, err = db.Exec(sqlStmt)
}

func dbSaveKudo(kudo *Kudo) {
	if kudo.ID != 0 {
		return
	}

	stmt, err := db.Prepare("INSERT INTO kudos (user, userFrom, kudo, likes, value) VALUES (?, ?, ?, ?, ?)")
	checkErr(err)

	res, err := stmt.Exec(kudo.MemberTo.ID, kudo.MemberFrom.ID, kudo.Kudo, kudo.LikeCount, kudo.Value)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	kudo.ID = id
}

func dbUpdateKudoLikes(kudo Kudo) {
	stmt, err := db.Prepare("UPDATE kudos SET likes = ? WHERE id = ?")
	checkErr(err)

	_, err = stmt.Exec(kudo.LikeCount, kudo.ID)
	checkErr(err)
}

func dbKudoLiked(kudo Kudo, member Member) bool {
	stmt, err := db.Prepare("SELECT COUNT(kudo) FROM likes WHERE kudo = ? AND user = ?")
	checkErr(err)

	row := stmt.QueryRow(kudo.ID, member.ID)

	var amount int
	row.Scan(&amount)

	return amount > 0
}

func addMemberToLike(kudo Kudo, member *Member) {
	stmt, err := db.Prepare("INSERT INTO	 likes (user, kudo) VALUES (?, ?)")
	checkErr(err)

	_, err = stmt.Exec(member.ID, kudo.ID)
	checkErr(err)
}

func loadKudos() {
	rows, err := db.Query("SELECT id, user, userFrom, kudo, likes, value FROM kudos")
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var kudo Kudo
		var memberFrom, memberTo string
		err = rows.Scan(&kudo.ID, &memberTo, &memberFrom, &kudo.Kudo, &kudo.LikeCount, &kudo.Value)
		checkErr(err)

		kudo.MemberFrom, _ = findMemberByID(memberFrom)
		kudo.MemberTo, _ = findMemberByID(memberTo)

		kudos = append(kudos, kudo)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
