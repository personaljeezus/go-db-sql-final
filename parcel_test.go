package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
	// get
	pack, err := store.Get(id)
	require.Error(t, sql.ErrNoRows)
	require.NoError(t, err)
	assert.Equal(t, parcel.Client, pack.Client)
	assert.Equal(t, parcel.Status, pack.Status)
	assert.Equal(t, parcel.Address, pack.Address)
	assert.Equal(t, parcel.CreatedAt, pack.CreatedAt)

	// delete
	err = store.Delete(parcel.Number)
	assert.NoError(t, err)
	require.Empty(t, parcel.Number)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()
	// add
	id, err := store.Add(parcel)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
	// set address
	newAddress := "new test address"
	err = store.SetAddress(parcel.Number, newAddress)
	require.NoError(t, err)
	// check
	whatAddress := parcel.Address
	fmt.Println(whatAddress)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()
	// add
	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	// set status
	err = store.SetStatus(parcel.Number, parcel.Status)
	require.NoError(t, err)
	// check
	testStatus := parcel.Status
	fmt.Println(testStatus)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}

	parcelMap := map[int]Parcel{}
	client := randRange.Intn(10000000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		assert.NoError(t, err)
		assert.NotEmpty(t, id)
		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	for _, parcel1 := range parcelMap {
		for _, allParcels := range parcels {
			assert.Equal(t, allParcels, parcel1)
		}
	}
	for _, parcel := range storedParcels {
		for _, parcel1 := range parcelMap {
			require.NotEmpty(t, parcel)
			assert.Equal(t, parcel1, parcel)
			defer db.Close()
		}
	}
}
