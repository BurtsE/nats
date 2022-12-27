package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"main/service"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type config struct {
	Host     string `json:host`
	Port     int    `json:port`
	User     string `json:user`
	Password string `json:password`
	Dbname   string `json:dbname`
}

var database *sql.DB

func GetDatabase() *sql.DB {
	return database
}
func ConnectTODB() {
	log.Println("connecting to db...")
	var conf config
	file, err := os.Open("../static/config.json")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	cdata, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(cdata, &conf)
	if err != nil {
		log.Fatal(err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, conf.Dbname)

	db, err := sql.Open("pgx", psqlInfo)

	if err != nil {
		log.Fatal("unable to connect to db", err)
	}
	log.Println("connected to db")
	database = db
}

func AddToDB(m service.Message) error {
	log.Println("adding to database...")
	if orderInDB(m.Order_uid) {
		return errors.New("already exists")
	}
	log.Println("1...")
	insert_order(m)

	log.Println("2...")
	insert_items(m)

	log.Println("3...")
	insert_delivery(m)

	log.Println("4...")
	insert_payment(m)
	log.Println("finished")
	return nil
}

func RecoverFromDB() error {
	log.Println("recovering from database")

	var data []service.Message
	rows, err := database.Query("SELECT * FROM orders;")
	defer rows.Close()
	if err != nil {
		return err
	}

	var m service.Message
	for rows.Next() {
		err = rows.Scan(&m.Order_uid, &m.Track_number, &m.Entry, &m.Locale, &m.Internal_signature, &m.Customer_id, &m.Delivery_service, &m.Shardkey,
			&m.Sm_id, &m.Date_created, &m.Oof_shard)
		if err != nil {
			return err
		}
		payment, err := database.Query("SELECT transaction, request_id, currency, provider, amount, payment_dt,"+
			"bank, delivery_cost, goods_total, custom_fee FROM payments WHERE order_id=$1;", m.Order_uid)
		defer payment.Close()
		if err != nil {
			return err
		}
		for payment.Next() {
			err = payment.Scan(&m.Payment.Transaction, &m.Payment.Request_id, &m.Payment.Currency, &m.Payment.Provider, &m.Payment.Amount,
				&m.Payment.Payment_dt, &m.Payment.Bank, &m.Payment.Delivery_cost, &m.Payment.Goods_total, &m.Payment.Custom_fee)
			if err != nil {
				return err
			}
		}

		delivery, err := database.Query("SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_id=$1;", m.Order_uid)
		defer delivery.Close()
		if err != nil {
			log.Println(err)
		}
		for delivery.Next() {
			err = delivery.Scan(&m.Delivery.Name, &m.Delivery.Phone, &m.Delivery.Zip, &m.Delivery.City, &m.Delivery.Address, &m.Delivery.Region,
				&m.Delivery.Email)
			if err != nil {
				return err
			}
		}
		sqlfile, err := os.Open("../sql/get_items.sql")
		defer sqlfile.Close()
		if err != nil {
			return err
		}
		sqldata, err := ioutil.ReadAll(sqlfile)
		if err != nil {
			return err
		}
		var sql string = string(sqldata)
		items, err := database.Query(sql, m.Order_uid)
		defer items.Close()
		if err != nil {
			return err
		}
		for items.Next() {
			var it service.Item
			err = items.Scan(&it.Chrt_id, &it.Track_number, &it.Price, &it.Rid, &it.Name, &it.Sale, &it.Size, &it.Total_price,
				&it.Nm_id, &it.Brand, &it.Status)
			if err != nil {
				return err
			}
			m.Items = append(m.Items, it)
		}
		data = append(data, m)
	}
	for _, m := range data {
		memory.Set(m.Order_uid, m)
	}
	log.Println("data recovered")
	return nil
}

func orderInDB(id string) bool {
	rows, err := database.Query(fmt.Sprintf("select order_uid from orders where order_uid='%s'", id))
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return true
	}
	if rows.Next() {
		log.Println("message already exists in database")
		return true
	}
	return false
}

func itemInDB(id int) bool {
	rows, err := database.Query(fmt.Sprintf("select chrt_id from items where chrt_id='%d'", id))
	defer rows.Close()
	if err != nil {
		log.Println(err)
	}
	if rows.Next() {
		return true
	}
	return false
}

func insert_items(m service.Message) {
	log.Println("insert items")
	sqlfile1, err := os.Open("../sql/insert_item.sql")
	if err != nil {
		log.Println(err)
	}
	defer sqlfile1.Close()
	sqlfile2, err := os.Open("../sql/insert_item2.sql")
	if err != nil {
		log.Println(err)
	}
	defer sqlfile2.Close()

	sqldata1, _ := ioutil.ReadAll(sqlfile1)
	var sql1 string = string(sqldata1)

	sqldata2, _ := ioutil.ReadAll(sqlfile2)
	var sql2 string = string(sqldata2)

	for _, item := range m.Items {
		if !itemInDB(item.Chrt_id) {
			_, err = database.Exec(sql1, item.Chrt_id, item.Track_number, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.Total_price,
				item.Nm_id, item.Brand, item.Status)
			if err != nil {
				log.Println(err)
			}
			_, err = database.Exec(sql2, m.Order_uid, item.Chrt_id)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func insert_order(m service.Message) {
	log.Println("insert order")
	sqlfile, err := os.Open("../sql/insert_order.sql")
	defer sqlfile.Close()
	log.Println(err)
	sqldata, _ := ioutil.ReadAll(sqlfile)
	var sql string = string(sqldata)
	_, err = database.Exec(sql, m.Order_uid, m.Track_number, m.Entry, m.Locale, m.Internal_signature, m.Customer_id, m.Delivery_service, m.Shardkey,
		m.Sm_id, m.Date_created, m.Oof_shard)
	if err != nil {
		log.Println(err)
	}
}

func insert_delivery(m service.Message) {
	log.Println("insert delivery")
	sqlfile, err := os.Open("../sql/insert_delivery.sql")
	defer sqlfile.Close()
	log.Println(err)
	sqldata, _ := ioutil.ReadAll(sqlfile)
	var sql string = string(sqldata)
	_, err = database.Exec(sql, m.Delivery.Name, m.Delivery.Phone, m.Delivery.Zip, m.Delivery.City, m.Delivery.Address, m.Delivery.Region,
		m.Delivery.Email, m.Order_uid)
	if err != nil {
		log.Println(err)
	}
}

func insert_payment(m service.Message) {
	log.Println("insert payment")
	sqlfile, err := os.Open("s../ql/insert_payment.sql")
	defer sqlfile.Close()
	log.Println(err)
	sqldata, _ := ioutil.ReadAll(sqlfile)
	var sql string = string(sqldata)
	_, err = database.Exec(sql, m.Payment.Transaction, m.Payment.Request_id, m.Payment.Currency, m.Payment.Provider, m.Payment.Amount,
		m.Payment.Payment_dt, m.Payment.Bank, m.Payment.Delivery_cost, m.Payment.Goods_total, m.Payment.Custom_fee, m.Order_uid)
	if err != nil {
		log.Println(err)
	}
}
