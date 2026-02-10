package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func Open()(*pgx.Conn,error){
	conn,err:=pgx.Connect(context.Background(),"postgres://postgres:postgres@postgres:5432/app_db")
	
	if err!=nil{
		return nil,err
	}

	if err:=createTables(conn);err!=nil{
		conn.Close(context.Background())
		return nil,err
	}

	if err:=conn.Ping(context.Background());err!=nil{
		conn.Close(context.Background())
		return nil,err

	}
	fmt.Println("Connected to the database")
	return conn,nil
}

func createTables(conn *pgx.Conn)error{
	query:=`
	CREATE TABLE IF NOT EXISTS  sum_results(
		id SERIAL PRIMARY KEY,
		num_1 INT NOT NULL,
		num_2 INT NOT NULL,
		result INT NOT NULL,
		created_at TIMESTAMP DEFAULT now()
	);

	CREATE TABLE IF NOT EXISTS outbox_events(
		id SERIAL PRIMARY KEY,
		event_type TEXT NOT NULL,
		payload JSONB NOT NULL,
		processed BOOLEAN DEFAULT false,
		created_at TIMESTAMP DEFAULT now()
	);
	
	`
	_,err:=conn.Exec(context.Background(),query)

	return err
}