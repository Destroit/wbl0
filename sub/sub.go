package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	stan "github.com/nats-io/stan.go"
	"html/template"
	"log"
	"net/http"
)

var tmpl *template.Template
var conn *pgx.Conn
var orderCache = make(map[string]JsonForm)
var validate *validator.Validate

// Page with search form
func httpRootHandler(w http.ResponseWriter, req *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

// Returns JSON if successful. If it isn't, error message
func httpSearchHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	orderId := req.FormValue("orderid")
	var orderByte []byte
	resForm, ok := orderCache[orderId]
	if !ok {
		fmt.Fprintln(w, "Id not found")
		return
	}
	orderByte, err := json.Marshal(resForm)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, "Id not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(orderByte))
}

func doAck(msg *stan.Msg) {
	if err := msg.Ack(); err != nil {
		log.Printf("Ack error: %v\n", err)
	}
}

// On new message unmarshal, validate, insert into cache and DB
func onStanMsg(m *stan.Msg) {
	var orderForm JsonForm
	err := json.Unmarshal(m.Data, &orderForm)
	if err != nil {
		log.Printf("JSON unmarshal error: %v\n", err)
		doAck(m)
		return
	}
	orderUid := orderForm.Order_uid

	err = validate.Struct(orderForm)
	if err != nil {
		log.Printf("Validation error:\n%v\n", err)
		doAck(m)
		return
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO orders (id, data) VALUES ($1, $2) ON CONFLICT DO NOTHING", orderUid, m.Data)
	if err != nil {
		log.Printf("postrgreSQL insert error: %v\n", err)
		return
	}
	_, exists := orderCache[orderUid]
	if !exists {
		orderCache[orderUid] = orderForm
		log.Printf("added order with uid:%s\n", orderUid)
	} else {
		log.Printf("order with uid:%s already exists\n", orderUid)
	}
	doAck(m)
}

func pgxInit(tmpAddr string) *pgx.Conn {
	tmpConn, err := pgx.Connect(context.Background(), tmpAddr)
	if err != nil {
		log.Fatalf("postgreSQL connection error: %v\n", err)
	}
	return tmpConn
}

func stanInit(tmpCluster string, tmpClient string) stan.Conn {
	tmpSc, err := stan.Connect("test-cluster", "test-client")
	if err != nil {
		log.Fatalf("stan connection error: %v\n", err)
	}
	return tmpSc
}

func main() {
	//postgres init
	addr := "postgres://l0user:l0jka@localhost:5432/l0db"
	conn = pgxInit(addr)
	defer conn.Close(context.Background())
	log.Println("postgreSQL init complete")

	//stan init
	cluster := "test-cluster"
	client := "test-client"
	sc := stanInit(cluster, client)
	opt := stan.SetManualAckMode()
	_, err := sc.Subscribe("orders", onStanMsg, stan.DurableName("l0-durable"), opt)
	if err != nil {
		log.Fatalf("stan subscription error: %v\n", err)
	}
	defer sc.Close()
	log.Println("stan init complete")

	//validator for JSON struct
	validate = validator.New(validator.WithRequiredStructEnabled())

	//copy from db into cache(map)
	rows, err := conn.Query(context.Background(), "SELECT * FROM orders")
	if err != nil {
		log.Fatalf("error copying data from db to cache: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var form JsonForm
		var id string
		err := rows.Scan(&id, &form)
		if err != nil {
			log.Fatalf("db rows parsing error: %v\n", err)
		} else {
			orderCache[id] = form
		}
	}
	log.Println("DB to cache copy complete")

	//Preparing HTTP server
	tmpl, err = template.ParseFiles("./index.html")
	if err != nil {
		log.Fatalf("template parsing error: %v\n", err)
	}

	http.HandleFunc("/", httpRootHandler)
	http.HandleFunc("/search", httpSearchHandler)
	port := ":8080"
	log.Println("HTTP server init complete")
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("HTTP server error: %v\n", err)
	}
}
