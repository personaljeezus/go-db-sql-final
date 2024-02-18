package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {

	res, err := s.db.Exec("INSERT INTO parcel (Client, Status, Address, Created_At) VALUES (:Client, :Status, :Address, :Created_At)",
		sql.Named("Client", p.Client),
		sql.Named("Status", p.Status),
		sql.Named("Address", p.Address),
		sql.Named("Created_At", p.CreatedAt))

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	ind, _ := res.LastInsertId()
	return int(ind), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	var p Parcel
	row := s.db.QueryRow("SELECT Client, Status, Address, Created_At FROM parcel WHERE Number = :Number",
		sql.Named("Number", number))
	err := row.Scan(&p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		fmt.Println(err)
		return p, err
	}
	return p, err
}
func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	var res []Parcel
	rows, err := s.db.Query("SELECT Client, Status, Address, Created_At FROM parcel WHERE Client = :Client", sql.Named("Client", client))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		par := Parcel{}
		err := rows.Scan(&par.Client, &par.Status, &par.Address, &par.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, par)
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET Status = :Status WHERE Number = :Number",
		sql.Named("Status", status),
		sql.Named("Number", number))

	if err != nil {

		fmt.Println(err)
		return err
	}
	return nil
}
func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := s.db.Exec("UPDATE parcel SET Address = :Address WHERE Number = :Number AND Status = :Status",
		sql.Named("Address", address),
		sql.Named("Number", number),
		sql.Named("Status", ParcelStatusRegistered))

	if err != nil {

		fmt.Println(err)
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	_, err := s.db.Exec("DELETE FROM parcel WHERE Number = :Number AND Status = :Status",
		sql.Named("Number", number),
		sql.Named("Status", ParcelStatusRegistered))

	if err != nil {

		fmt.Println(err)
		return err
	}
	return nil
}
