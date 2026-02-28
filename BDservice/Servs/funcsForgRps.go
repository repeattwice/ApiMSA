package servs

import "context"

func CreateAccount(user_name string, last_name string, email string, ctx context.Context) error {
	conn, err := GetBDConnection(ctx)
	if err != nil {
		return err
	}
	sqlQuery := `
	INSERT INTO users (user_name, last_name, email)
	SELECT $1, $2, $3
	`
	_, err = conn.Exec(ctx, sqlQuery, user_name, last_name, email)
	if err != nil {
		return err
	}
	return nil
}
