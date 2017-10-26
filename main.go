package main

import (  
  "net/http"
  "Projects/VinMVC/views"
  "fmt"  
  "database/sql"
  _"github.com/mattn/go-sqlite3"
  //"strings"
)

var index *views.View  
var contact *views.View
var results []string
var db *sql.DB

type chatMsg struct {
  id int
  email string
  msg string
}
 
 
type TestItem struct {
  id int
  email string
  msg string
}
   
func InitDB(filepath string) {
  var err error  
  db, err = sql.Open("sqlite3", filepath)
  if err != nil {panic(err)}
  if db == nil {panic("db nil")}
//  return db
}


func CreateTable(){
  // create if not exist
  sqltable := "CREATE TABLE IF NOT EXISTS chat(id INTEGER PRIMARY KEY, Email TEXT, Msg TEXT, Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP);"

    _,err := db.Exec(sqltable)
    checkErr(err)
}
  
func StoreItem( item TestItem){
  sqladdItem := "INSERT OR REPLACE INTO chat( Email, Msg) values(?,?)"
  stmt, err := db.Prepare(sqladdItem)
  if err != nil {panic(err)}
  defer stmt.Close()

  
  _, err2 := stmt.Exec(item.email, item.msg)
  if err2 != nil{panic(err2)}
  

}
 
func ReadItem() []TestItem {
  sqlreadall := "SELECT Id, Email, Msg FROM chat ORDER BY datetime(TimeStamp) DESC"
 
  rows, err := db.Query(sqlreadall)
  if err != nil {panic(err)}
  defer rows.Close()
 
  var result []TestItem
  for rows.Next() {
    item := TestItem{}
    err2 := rows.Scan(&item.id, &item.email, &item.msg )
    if err2 != nil {panic(err2)}
    result = append(result, item)
  }
  return result
}

 

func main() {  
  index = views.NewView("bootstrap", "views/index.gohtml")
  contact = views.NewView("bootstrap", "views/contact.gohtml")
  
  const dbPath = "chatter.db"
  InitDB(dbPath)
  //defer db.Close()
  CreateTable()
  fmt.Print("start this baby")
  http.HandleFunc("/", indexHandler)
  http.HandleFunc("/contact", contactHandler)
  
  http.ListenAndServe(":3000", nil)
  
   
}
   
func indexHandler(w http.ResponseWriter, r *http.Request) {  
  info := ReadItem()
  //for _, itemA := range info{
    //fmt.Println(itemA.email + " ddd " + itemA.msg)
  //}
  index.Render(w, info)
}
   
func contactHandler(w http.ResponseWriter, r *http.Request) {  
  if r.Method == "GET"{
  contact.Render(w, nil)
  } else{
    r.ParseForm()
    //fmt.Printf("username: %s",r.FormValue("email"))
    //fmt.Println("message",r.FormValue("message"))
   // results = append(results,r.Form["email"])//, r.Form["message"])
    //http.RedirectHandler("/", 500)
    var item TestItem 
    item.msg = r.FormValue("message")
    item.email = r.FormValue("email")
    if item.msg != "" {
      if item.email != "" {
      StoreItem(item)
    }}
    //readItem := ReadItem()
 //   for _, itemA := range readItem{
   ///   fmt.Println(itemA.email + " : " + itemA.msg)
    //}/
      http.Redirect(w, r, "/", http.StatusSeeOther)
  } 
} 
//$key, $value :=  
 
func checkErr(err error) {
  if err != nil {
      panic(err)
  }
}
