package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/freeloginname/otusGoBasicProject/internal/repository/notes"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetUser(ctx context.Context, dbc *pgxpool.Pool, name string) (notes.User, error) {

	tx, err := dbc.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		fmt.Printf("failed to start transaction: %v", err)
		return notes.User{}, err
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			fmt.Printf("failed to rollback transaction: %v", err)
		}
	}()

	requestor := notes.New(tx)
	requestor.WithTx(tx)
	user, err := requestor.GetUserByName(ctx, name)
	if err != nil {
		fmt.Printf("failed to find user: %v", err)
		return notes.User{}, err
	}
	return *user, nil
}

func CreateUser(ctx context.Context, dbc *pgxpool.Pool, name string, password string) (string, error) {

	tx, err := dbc.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		fmt.Printf("failed to start transaction: %v", err)
		return "", err
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			fmt.Printf("failed to rollback transaction: %v", err)
		}
	}()

	requestor := notes.New(tx)
	requestor.WithTx(tx)
	users, err := requestor.GetAllUsers(ctx)
	if err != nil {
		fmt.Printf("failed to get all users: %v", err)
		return "", err
	}
	for _, user := range users {
		if user.Name == name {
			errorString := fmt.Sprintf("user with name %v already existst with id: %v", user.Name, user.ID)
			fmt.Println(errorString)
			err = errors.New(errorString)
			return "", err
		}
	}
	userID, err := requestor.CreateUser(ctx, notes.CreateUserParams{
		Name:     name,
		Password: password,
	})
	if err != nil {
		fmt.Printf("failed to create user: %v", err)
		return "", err
	}
	err = tx.Commit(ctx)
	if err != nil {
		fmt.Printf("failed to commit transaction: %v", err)
		return "", err
	}
	return userID.String(), nil
}

func CreateNote(ctx context.Context, dbc *pgxpool.Pool, userName string, name string, text string) (string, error) {
	tx, err := dbc.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		fmt.Printf("failed to start transaction: %v", err)
		return "", err
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			fmt.Printf("failed to rollback transaction: %v", err)
		}
	}()

	requestor := notes.New(tx)
	requestor.WithTx(tx)
	users, err := requestor.GetAllUsers(ctx)
	if err != nil {
		fmt.Printf("failed to get all users: %v", err)
		return "", err
	}
	userID := pgtype.UUID{
		Bytes: [16]byte{},
		Valid: false,
	}
	for _, user := range users {
		if user.Name == userName {
			userID = user.ID
			break
		}
	}
	if !userID.Valid {
		errorString := fmt.Sprintf("user with name %v not found", userName)
		fmt.Println(errorString)
		err = errors.New(errorString)
		return "", err
	}
	allNotes, err := requestor.GetUserNotes(ctx, userID)
	if err != nil {
		fmt.Printf("failed to get user notes: %v", err)
		return "", err
	}
	for _, note := range allNotes {
		if note.Name == name {
			errorString := fmt.Sprintf("note with name %v already existst with id: %v", note.Name, note.ID)
			fmt.Println(errorString)
			err = errors.New(errorString)
			return "", err
		}
	}

	noteID, err := requestor.CreateNote(ctx, notes.CreateNoteParams{
		UserID: userID,
		Name:   name,
		Text:   text,
	})
	if err != nil {
		fmt.Printf("failed to create note: %v", err)
		return "", err
	}
	err = tx.Commit(ctx)
	if err != nil {
		fmt.Printf("failed to commit transaction: %v", err)
		return "", err
	}
	return noteID.String(), nil
}

func GetAllUserNotes(ctx context.Context, dbc *pgxpool.Pool, userName string) ([]*notes.Note, error) {
	tx, err := dbc.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		fmt.Printf("failed to start transaction: %v", err)
		return []*notes.Note{}, err
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			fmt.Printf("failed to rollback transaction: %v", err)
		}
	}()

	requestor := notes.New(tx)
	requestor.WithTx(tx)
	users, err := requestor.GetAllUsers(ctx)
	if err != nil {
		fmt.Printf("failed to get all users: %v", err)
		return []*notes.Note{}, err
	}
	userID := pgtype.UUID{
		Bytes: [16]byte{},
		Valid: false,
	}
	for _, user := range users {
		if user.Name == userName {
			userID = user.ID
			break
		}
	}
	if !userID.Valid {
		errorString := fmt.Sprintf("user with name %v not found", userName)
		fmt.Println(errorString)
		err = errors.New(errorString)
		return []*notes.Note{}, err
	}
	allNotes, err := requestor.GetUserNotes(ctx, userID)
	if err != nil {
		fmt.Printf("failed to get user notes: %v", err)
		return []*notes.Note{}, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		fmt.Printf("failed to commit transaction: %v", err)
		return []*notes.Note{}, err
	}
	return allNotes, nil
}

func GetAllNotes(ctx context.Context, dbc *pgxpool.Pool) ([]*notes.Note, error) {
	tx, err := dbc.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		fmt.Printf("failed to start transaction: %v", err)
		return []*notes.Note{}, err
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			fmt.Printf("failed to rollback transaction: %v", err)
		}
	}()

	requestor := notes.New(tx)
	requestor.WithTx(tx)
	allNotes, err := requestor.GetAllNotes(ctx)
	if err != nil {
		fmt.Printf("failed to get all notes: %v", err)
		return []*notes.Note{}, err
	}
	return allNotes, nil
}

func GetNotes(ctx context.Context, dbc *pgxpool.Pool, userName string) ([]*notes.Note, error) {
	tx, err := dbc.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		fmt.Printf("failed to start transaction: %v", err)
		return []*notes.Note{}, err
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			fmt.Printf("failed to rollback transaction: %v", err)
		}
	}()

	requestor := notes.New(tx)
	requestor.WithTx(tx)
	user, err := requestor.GetUserByName(ctx, userName)
	if err != nil {
		fmt.Printf("failed to find user: %v", err)
		return []*notes.Note{}, err
	}
	allNotes, err := requestor.GetUserNotes(ctx, user.ID)
	if err != nil {
		fmt.Printf("failed to get user notes: %v", err)
		return []*notes.Note{}, err
	}
	return allNotes, nil

}

func GetNoteByID(ctx context.Context, dbc *pgxpool.Pool, noteID string) (*notes.Note, error) {
	tx, err := dbc.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		fmt.Printf("failed to start transaction: %v", err)
		return &notes.Note{}, err
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			fmt.Printf("failed to rollback transaction: %v", err)
		}
	}()

	requestor := notes.New(tx)
	requestor.WithTx(tx)
	var noteIDPg pgtype.UUID
	noteIDPg.Scan(noteID)
	note, err := requestor.GetNote(ctx, noteIDPg)
	if err != nil {
		fmt.Printf("failed to get user note: %v", err)
		return &notes.Note{}, err
	}
	return note, nil
}

// func GetNote(ctx context.Context, dbc *pgxpool.Pool, userName string, noteName string) (*notes.Note, error) {
// 	tx, err := dbc.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
// 	if err != nil {
// 		fmt.Printf("failed to start transaction: %v", err)
// 		return &notes.Note{}, err
// 	}
// 	defer func() {
// 		err = tx.Rollback(ctx)
// 		if err != nil {
// 			fmt.Printf("failed to rollback transaction: %v", err)
// 		}
// 	}()

// 	requestor := notes.New(tx)
// 	requestor.WithTx(tx)
// 	user, err := requestor.GetUserByName(ctx, userName)
// 	if err != nil {
// 		fmt.Printf("failed to find user: %v", err)
// 		return &notes.Note{}, err
// 	}
// 	note, err := requestor.GetUserNoteByName(ctx, notes.GetUserNoteByNameParams{
// 		UserID: user.ID,
// 		Name:   noteName,
// 	})
// 	if err != nil {
// 		fmt.Printf("failed to get user note: %v", err)
// 		return &notes.Note{}, err
// 	}
// 	return note, nil
// }

func GetUserNoteByName(ctx context.Context, dbc *pgxpool.Pool, userName string, noteName string) (*notes.Note, error) {
	tx, err := dbc.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		fmt.Printf("failed to start transaction: %v", err)
		return &notes.Note{}, err
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			fmt.Printf("failed to rollback transaction: %v", err)
		}
	}()

	requestor := notes.New(tx)
	requestor.WithTx(tx)
	user, err := requestor.GetUserByName(ctx, userName)
	if err != nil {
		fmt.Printf("failed to find user: %v", err)
		return &notes.Note{}, err
	}
	note, err := requestor.GetUserNoteByName(ctx, notes.GetUserNoteByNameParams{
		UserID: user.ID,
		Name:   noteName,
	})
	if err != nil {
		fmt.Printf("failed to get user note: %v", err)
		return &notes.Note{}, err
	}
	return note, nil
}

func UpdateNote(ctx context.Context, dbc *pgxpool.Pool, userName string, noteName string, text string) error {
	tx, err := dbc.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		fmt.Printf("failed to start transaction: %v", err)
		return err
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			fmt.Printf("failed to rollback transaction: %v", err)
		}
	}()

	requestor := notes.New(tx)
	requestor.WithTx(tx)
	user, err := requestor.GetUserByName(ctx, userName)
	if err != nil {
		fmt.Printf("failed to find user: %v", err)
		return err
	}
	_, err = requestor.GetUserNoteByName(ctx, notes.GetUserNoteByNameParams{
		UserID: user.ID,
		Name:   noteName,
	})
	if err != nil {
		fmt.Printf("failed to get user note: %v", err)
		return err
	}

	err = requestor.UpdateUserNoteByName(ctx, notes.UpdateUserNoteByNameParams{
		Name:   noteName,
		UserID: user.ID,
		Text:   text,
	})
	if err != nil {
		fmt.Printf("failed to update user note %v: %v", noteName, err)
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		fmt.Printf("failed to commit transaction: %v", err)
		return err
	}
	return nil
}

func DeleteUserNoteByName(ctx context.Context, dbc *pgxpool.Pool, userName string, noteName string) error {
	tx, err := dbc.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		fmt.Printf("failed to start transaction: %v", err)
		return err
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			fmt.Printf("failed to rollback transaction: %v", err)
		}
	}()

	requestor := notes.New(tx)
	requestor.WithTx(tx)
	user, err := requestor.GetUserByName(ctx, userName)
	if err != nil {
		fmt.Printf("failed to find user: %v", err)
		return err
	}
	note, err := requestor.GetUserNoteByName(ctx, notes.GetUserNoteByNameParams{
		UserID: user.ID,
		Name:   noteName,
	})
	if err != nil {
		fmt.Printf("failed to get user note: %v", err)
		return err
	}
	err = requestor.DeleteNoteById(ctx, note.ID)
	if err != nil {
		fmt.Printf("failed to delete user note %v: %v", noteName, err)
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		fmt.Printf("failed to commit transaction: %v", err)
		return err
	}
	return nil
}

func DeleteNoteByID(ctx context.Context, dbc *pgxpool.Pool, noteID string) error {
	tx, err := dbc.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		fmt.Printf("failed to start transaction: %v", err)
		return err
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			fmt.Printf("failed to rollback transaction: %v", err)
		}
	}()

	requestor := notes.New(tx)
	requestor.WithTx(tx)
	var noteIDPg pgtype.UUID
	noteIDPg.Scan(noteID)
	err = requestor.DeleteNoteById(ctx, noteIDPg)
	if err != nil {
		fmt.Printf("failed to delete user note %v: %v", noteID, err)
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		fmt.Printf("failed to commit transaction: %v", err)
		return err
	}
	return nil
}

func DeleteNote(ctx context.Context, dbc *pgxpool.Pool, userName string, noteName string) error {
	tx, err := dbc.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		fmt.Printf("failed to start transaction: %v", err)
		return err
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			fmt.Printf("failed to rollback transaction: %v", err)
		}
	}()

	requestor := notes.New(tx)
	requestor.WithTx(tx)
	user, err := requestor.GetUserByName(ctx, userName)
	if err != nil {
		fmt.Printf("failed to find user: %v", err)
		return err
	}
	note, err := requestor.GetUserNoteByName(ctx, notes.GetUserNoteByNameParams{
		UserID: user.ID,
		Name:   noteName,
	})
	if err != nil {
		fmt.Printf("failed to get user note: %v", err)
		return err
	}
	err = requestor.DeleteNoteById(ctx, note.ID)
	if err != nil {
		fmt.Printf("failed to delete user note: %v", err)
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		fmt.Printf("failed to commit transaction: %v", err)
		return err
	}
	return nil
}

// func GetOrdersByUser(ctx context.Context, dsn string, userName string) ([]*product.Order, error) {
// 	var userID uuid.UUID
// 	dbc, err := pgdb.New(ctx, dsn, 1)
// 	if err != nil {
// 		fmt.Printf("failed to connect to DB: %v", err)
// 		return nil, err
// 	}
// 	defer dbc.Close()

// 	tx, err := dbc.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
// 	if err != nil {
// 		fmt.Printf("failed to start transaction: %v", err)
// 		return []*product.Order{}, err
// 	}
// 	defer func() {
// 		err = tx.Rollback(ctx)
// 		if err != nil {
// 			fmt.Printf("failed to rollback transaction: %v", err)
// 		}
// 	}()

// 	requestor := product.New(tx)
// 	requestor.WithTx(tx)
// 	users, err := requestor.GetAllUsers(ctx)
// 	if err != nil {
// 		fmt.Printf("failed to get information about users: %v", err)
// 		return []*product.Order{}, err
// 	}

// 	for _, user := range users {
// 		if user.Name == userName {
// 			userID = user.ID
// 			break
// 		}
// 	}
// 	if userID == [16]byte{} {
// 		errorString := fmt.Sprintf("user with name %v not found", userName)
// 		fmt.Println(errorString)
// 		err = errors.New(errorString)
// 		return []*product.Order{}, err
// 	}

// 	orders, err := requestor.GetOrdersByUser(ctx, userID)
// 	if err != nil {
// 		fmt.Printf("failed to get information about orders: %v", err)
// 		return []*product.Order{}, err
// 	}
// 	err = tx.Commit(ctx)
// 	if err != nil {
// 		fmt.Printf("failed to commit transaction: %v", err)
// 		return []*product.Order{}, err
// 	}
// 	return orders, nil
// }

// func CreateOrder(ctx context.Context, dsn string, userName string, totalAmount string) error {
// 	var numeric pgtype.Numeric
// 	numeric.Scan(totalAmount)
// 	dbc, err := pgdb.New(ctx, dsn, 1)
// 	if err != nil {
// 		fmt.Printf("failed to connect to DB: %v", err)
// 		return err
// 	}
// 	defer dbc.Close()
// 	tx, err := dbc.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
// 	if err != nil {
// 		fmt.Printf("failed to start transaction: %v", err)
// 		return err
// 	}
// 	defer func() {
// 		err = tx.Rollback(ctx)
// 		if err != nil {
// 			fmt.Printf("failed to rollback transaction: %v", err)
// 		}
// 	}()

// 	requestor := product.New(tx)
// 	requestor.WithTx(tx)

// 	users, err := requestor.GetAllUsers(ctx)
// 	if err != nil {
// 		fmt.Printf("failed to get all users: %v", err)
// 		return err
// 	}
// 	for _, user := range users {
// 		if user.Name == userName {
// 			_, err = requestor.CreateOrderWithCurrentDate(ctx, product.CreateOrderWithCurrentDateParams{
// 				UserID:      user.ID,
// 				TotalAmount: numeric,
// 			})
// 			if err != nil {
// 				fmt.Printf("failed to create order: %v", err)
// 				return err
// 			}
// 			err = tx.Commit(ctx)
// 			if err != nil {
// 				fmt.Printf("failed to commit transaction: %v", err)
// 				return err
// 			}
// 			return nil
// 		}
// 	}
// 	errorString := fmt.Sprintf("user with name %v not found", userName)
// 	fmt.Println(errorString)
// 	err = errors.New(errorString)
// 	return err
// }

// func GetProductByName(ctx context.Context, dsn string, name string) (product.Product, error) {
// 	dbc, err := pgdb.New(ctx, dsn, 1)
// 	if err != nil {
// 		fmt.Printf("failed to connect to DB: %v", err)
// 		return product.Product{}, err
// 	}
// 	defer dbc.Close()
// 	requestor := product.New(dbc)
// 	product, err := requestor.GetProductByName(ctx, name)
// 	if err != nil {
// 		fmt.Printf("failed to get product by name: %v", err)
// 		return *product, err
// 	}
// 	return *product, nil
// }

// func CreateProduct(ctx context.Context, dsn string, name string, price string) (string, error) {
// 	var numeric pgtype.Numeric
// 	numeric.Scan(price)
// 	dbc, err := pgdb.New(ctx, dsn, 1)
// 	if err != nil {
// 		fmt.Printf("failed to connect to DB: %v", err)
// 		return "", err
// 	}
// 	defer dbc.Close()

// 	tx, err := dbc.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
// 	if err != nil {
// 		fmt.Printf("failed to start transaction: %v", err)
// 		return "", err
// 	}
// 	defer func() {
// 		err = tx.Rollback(ctx)
// 		if err != nil {
// 			fmt.Printf("failed to rollback transaction: %v", err)
// 		}
// 	}()

// 	requestor := product.New(tx)
// 	requestor.WithTx(tx)
// 	// TODO добавить проверку наличия данного товара перед созданием
// 	products, err := requestor.GetAllProducts(ctx)
// 	if err != nil {
// 		fmt.Printf("failed to get all products: %v", err)
// 		return "", err
// 	}
// 	for _, product := range products {
// 		if product.Name == name {
// 			errorString := fmt.Sprintf("product with name %v already existst with id: %v", product.Name, product.ID)
// 			fmt.Println(errorString)
// 			err = errors.New(errorString)
// 			return "", err
// 		}
// 	}
// 	productID, err := requestor.CreateProduct(ctx, product.CreateProductParams{
// 		Name:  name,
// 		Price: numeric,
// 	})
// 	if err != nil {
// 		fmt.Printf("failed to create product: %v", err)
// 		return "", err
// 	}
// 	err = tx.Commit(ctx)
// 	if err != nil {
// 		fmt.Printf("failed to commit transaction: %v", err)
// 		return "", err
// 	}
// 	return productID.String(), nil
// }
