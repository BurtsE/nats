package storage

import (
	"main/service"
	"testing"
)

// Не выполняются параллельно

func TestCreationDB(t *testing.T) {
	ConnectTODB()
	db := GetDatabase()
	defer db.Close()
	if db == nil {
		t.Error("No connection to database")
	}
	_, err := db.Exec("select;")
	if err != nil {
		t.Errorf("unable to execute request: %s", err)
	}
}

func TestInsertOrder(t *testing.T) {
	ConnectTODB()
	db := GetDatabase()
	defer db.Close()
	m := service.Message{
		Order_uid:    "0000000",
		Date_created: "1000-10-01",
	}

	insert_order(m)
	var ans string
	rows, err := db.Query("select order_uid from orders where order_uid='0000000' and date_created='1000-10-01';")
	if err != nil {
		t.Error(err)
	}
	for rows.Next() {
		rows.Scan(&ans)
	}
	if ans != "0000000" {
		t.Error("data was changed or lost")
	}
	db.Exec("delete from orders where order_uid='0000000' and date_created='1000-10-01' and locale='';")
}

func TestOrderInDB(t *testing.T) {
	ConnectTODB()
	db := GetDatabase()
	defer db.Close()

	_, err := db.Exec("insert into orders(order_uid) values('000000000');")
	if err != nil {
		t.Error("unable to execute request")
	}
	if !orderInDB("000000000") {
		t.Error("unable to find order")
	}
	db.Exec("delete from orders where order_uid='000000000' and date_created is NULL")
}

func TestItemInDB(t *testing.T) {
	ConnectTODB()
	db := GetDatabase()
	defer db.Close()

	_, err := db.Exec("insert into items(chrt_id) values(123);")
	if err != nil {
		t.Error("unable to execute request")
	}
	if !itemInDB(123) {
		t.Error("unable to find item")
	}
	db.Exec("delete from items where chrt_id=123 and track_number is NULL")

}

func TestInsertItems(t *testing.T) {
	ConnectTODB()
	db := GetDatabase()
	defer db.Close()
	var item = service.Item{
		Chrt_id: 0,
	}
	var m = service.Message{
		Order_uid:    "0",
		Items:        []service.Item{item},
		Date_created: "0001-01-01",
	}
	_, err := db.Exec("insert into orders(order_uid) values('0');")
	if err != nil {
		t.Error("unable to execute request")
	}
	insert_items(m)
	rows, err := db.Query("select chrt_id from items inner join order_item on(chrt_id=item_id) where chrt_id='0' and status is NULL")
	if err != nil {
		t.Error(err)
	}
	var ans int
	for rows.Next() {
		rows.Scan(&ans)
	}
	if ans != 0 {
		t.Error("data was changed or lost")
	}
	//db.Exec("delete from order_item where order_id='0' and item_id=0;")
	db.Exec("delete from orders where order_uid='0' and locale is NULL;")
	db.Exec("delete from items where chrt_id=0 and track_number='';")
}

func TestInsertDelivery(t *testing.T) {
	ConnectTODB()
	db := GetDatabase()
	defer db.Close()
	m := service.Message{
		Order_uid:    "0",
		Date_created: "1000-10-01",
	}
	m.Delivery.Phone = "-012"
	_, err := db.Exec("insert into orders(order_uid) values('0');")
	if err != nil {
		t.Error("unable to execute request")
	}
	insert_delivery(m)
	var ans string
	rows, err := db.Query("select order_id from delivery where order_id='0' and phone='-012'  and name='';")
	if err != nil {
		t.Error(err)
	}
	for rows.Next() {
		rows.Scan(&ans)
	}
	if ans != "0" {
		t.Errorf("data was changed or lost: %s", ans)
	}
	db.Exec("delete from orders where order_uid='0' and locale is  NULL;")
	db.Exec("delete from delivery where order_id='0' and phone='-012' and and name='';")
}

func TestInsertPayment(t *testing.T) {
	ConnectTODB()
	db := GetDatabase()
	defer db.Close()
	m := service.Message{
		Order_uid:    "0",
		Date_created: "1000-10-01",
	}
	m.Payment.Bank = "-012"
	_, err := db.Exec("insert into orders(order_uid) values('0');")
	if err != nil {
		t.Error("unable to execute request")
	}
	insert_payment(m)
	var ans string
	rows, err := db.Query("select order_id from payments where order_id='0' and bank='-012';")
	if err != nil {
		t.Error(err)
	}
	for rows.Next() {
		rows.Scan(&ans)
	}
	if ans != "0" {
		t.Errorf("data was changed or lost: %s", ans)
	}
	db.Exec("delete from orders where order_uid='0' and locale is  NULL;")
	db.Exec("delete from payments where order_id='0' and and bank='-012';")
}
func TestUsage(t *testing.T) {
	ConnectTODB()
	db := GetDatabase()
	defer db.Close()
	c := GetCache()
	m := service.Message{
		Order_uid:    "0",
		Date_created: "0001-01-01",
	}
	AddToDB(m)
	err := RecoverFromDB()
	if err != nil {
		t.Error("unable to recover data from database")
	}
	val, ok := c.storage["0"]
	if !ok || val.Order_uid != "0" {
		t.Error("data was changed or erased")
	}
	db.Exec("delete from payments where order_id='0' and bank='';")
	db.Exec("delete from delivery where order_id='0' and city='';")
	db.Exec("delete from orders where order_uid='0' and date_created='0001-01-01';")
}
