package db

import (
	"backend/internal/model"
)

func GetAllPaymentEvents() ([]model.PaymentEvent, error) {
	query := `SELECT user_id, payment_id, deposit_amount FROM payment_event`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.PaymentEvent
	for rows.Next() {
		var e model.PaymentEvent
		if err := rows.Scan(&e.UserID, &e.PaymentID, &e.DepositAmount); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func InsertPaymentEvent(event model.PaymentEvent) error {
	query := `
	INSERT INTO payment_event (user_id, payment_id, deposit_amount)
	VALUES (?, ?, ?)
	`
	_, err := DB.Exec(query, event.UserID, event.PaymentID, event.DepositAmount)
	if err != nil {
		return err
	}
	return nil
}

func InsertSkippedMessage(event model.PaymentEvent) error {
	// id auto
	query := `
	INSERT INTO skipped_message (user_id, payment_id, deposit_amount)
	VALUES (?, ?, ?)
	`
	_, err := DB.Exec(query, event.UserID, event.PaymentID, event.DepositAmount)
	if err != nil {
		return err
	}
	return nil
}
