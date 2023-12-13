package utils

import (
	"database/sql"
	"encoding/json"
	"log"
	"nats-service/app/database"
	"nats-service/app/model"
)

func CreateTableIfNotExist(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS orders_json (" +
		"id serial primary key," +
		"data jsonb" +
		");")

	if err != nil {
		log.Println(err)
	}
}

func InsertToTable(db *sql.DB, data []byte) {
	_, err := db.Exec("INSERT INTO orders_json (data) VALUES ($1)", data)
	if err != nil {
		log.Printf("Error inserting data into PostgreSQL: %v", err)
	}
}

func RetrieveOrdersFromDB(user string, password string, dbname string, host string, port int) (map[int]model.Order, error) {
	orders := make(map[int]model.Order)

	// Открываем соединение с базой данных
	// database := &database.Database{}
	database := database.New(user, password, dbname, host, port)
	db, err := database.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// testStruct := struct {
	// 	id    int
	// 	order Order
	// }{}

	// Выполняем запрос к базе данных
	rows, err := db.Query("SELECT id, data FROM orders_json")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Итерируем по результатам запроса
	for rows.Next() {
		var id int
		var data []byte
		if err := rows.Scan(&id, &data); err != nil {
			return nil, err
		}

		// Декодируем данные в структуру Order
		var order model.Order
		if err := json.Unmarshal(data, &order); err != nil {
			return nil, err
		}

		// Добавляем заказ в карту
		orders[id] = order
	}

	// Проверяем наличие ошибок после завершения итерации
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
