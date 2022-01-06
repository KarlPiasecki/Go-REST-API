package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "root"
	password = ""
	hostname = "localhost:3306"
	dbname   = "games"
)

type Game struct {
	category    string
	title       string
	description string
}

func DSN(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)
}

func connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", DSN(""))
	if err != nil {
		log.Printf("Error %s", err)
		return nil, err
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		log.Printf("Error %s", err)
		return nil, err
	}
	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s", err)
		return nil, err
	}
	log.Printf("rows affected %d\n", no)

	db.Close()
	db, err = sql.Open("mysql", DSN(dbname))
	if err != nil {
		log.Printf("Error %s", err)
		return nil, err
	}

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Error %s", err)
		return nil, err
	}
	log.Printf("Connected to DB %s successfully\n", dbname)
	return db, nil
}

func createProductTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS games(id int PRIMARY KEY AUTO_INCREMENT NOT NULL, category varchar(255) NOT NULL, 
        title varchar(255) NOT NULL, description text NOT NULL)`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s", err)
		return err
	}
	log.Printf("Rows affected when creating table: %d", rows)
	return nil
}

func insert(db *sql.DB, g Game) error {
	query := "INSERT INTO games(category, title, description) VALUES (?, ?, ?)"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, g.category, g.title, g.description)
	if err != nil {
		log.Printf("Error %s", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s", err)
		return err
	}
	log.Printf("%d game created ", rows)
	return nil
}

func main() {
	db, err := connect()
	if err != nil {
		log.Printf("Error %s", err)
		return
	}
	defer db.Close()
	log.Printf("Successfully connected to database")
	err = createProductTable(db)
	if err != nil {
		log.Printf("Error %s", err)
		return

	}

	g := Game{
		category:    "Retro",
		title:       "Sonic the Hedgehog 3 & Knuckles",
		description: "lorem ipsum something",
	}
	err = insert(db, g)
	if err != nil {
		log.Printf("Error %s", err)
		return
	}
}
