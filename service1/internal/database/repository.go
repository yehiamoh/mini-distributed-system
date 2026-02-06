package database

import (
	"context"
	"service1/internal/models"

	"github.com/jackc/pgx/v5"
)

type MessageRepository struct {
	conn *pgx.Conn
}

func NewMessageRepository(conn *pgx.Conn) *MessageRepository{
	return &MessageRepository{conn: conn}
}

func(r *MessageRepository)CreateMessageWithOutbox(ctx context.Context,messagePayload models.MessagePayLoad,payload []byte)error{
	tx,err:=r.conn.Begin(ctx)
	if err!=nil{
		return err
	}

	_,err=tx.Exec(ctx,`INSERT INTO sum_results(num_1,num_2,result)Values($1,$2,$3)`,messagePayload.Num_1,messagePayload.Nums_2,messagePayload.Result)

	if err!=nil{
		tx.Rollback(ctx)
		return err
	}

	_,err=tx.Exec(ctx,`INSERT INTO outbox_events(event_type,payload)Values($1,$2)`,"message.created",payload)

	if err!=nil{
		tx.Rollback(ctx)
		return err
	}

	err=tx.Commit(ctx)

	if err!=nil{
		return err
	}

	return nil

}

func(r *MessageRepository)GetAllUnProcessedMessages(ctx context.Context,limit int)([]models.OutboxEvent,error){
	tx,err:=r.conn.Begin(ctx)
	if err!=nil{
		return nil,err
	}

	rows,err:=tx.Query(ctx,`
		SELECT id,event_type,payload,processed,created_at
		FROM outbox_events
		WHERE processed=false
		ORDER BY created_at
		LIMIT $1
	`,limit)

	if err!=nil{
		tx.Rollback(ctx)
		return nil,err
		
	}
	defer rows.Close()

	events:=make([]models.OutboxEvent,0)

	for rows.Next(){
		var e models.OutboxEvent

		err:=rows.Scan(
			&e.ID,
			&e.EventType,
			&e.Payload,
			&e.Processed,
			&e.CreatedAt,
		)
		if err!=nil{
			tx.Rollback(ctx)
			return nil,err
		}

		events = append(events, e)

		if rows.Err()!=nil{
			return nil,rows.Err()
		}

	}
		return events,nil
}

func(r *MessageRepository)UpdateMessageProcessingStatus(ctx context.Context,id int)error{
	tx,err:=r.conn.Begin(ctx)

	_,err=tx.Exec(ctx,`
	UPDATE outbox_events
	SET processed=true
	WHERE id=$1
	`,id)

	if err!=nil{
		tx.Rollback(ctx)
		return err
	}

	return nil

}