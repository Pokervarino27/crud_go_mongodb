package main

import(
  "encoding/json"
  "fmt"
  "log"
  "context"
  "net/http"
  "github.com/mongodb/mongo-go-driver/mongo"
  "github.com/mongodb/mongo-go-driver/bson"
  "github.com/mongodb/mongo-go-driver/mongo/options"
  "github.com/mongodb/mongo-go-driver/bson/primitive"
  "github.com/gorilla/mux"
)

func getClient() *mongo.Client{
  client, err := mongo.Connect(context.TODO(), "mongodb://localhost:27017")

  if err != nil{
    log.Fatal(err)
  }
  //Chequea la conexión
  err = client.Ping(context.TODO(), nil)

  if err != nil{
    log.Fatal(err)
  }

  fmt.Println("Connected to Mongo")
  return client
}

func response(w http.ResponseWriter, status int, results Movie){
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(200)
  json.NewEncoder(w).Encode(results)
}

func Index(w http.ResponseWriter, r *http.Request){
  fmt.Fprintf(w, "Hola soy CpeBot")
}

var collection = getClient().Database("testDb").Collection("movies")

func MovieAdd(w http.ResponseWriter, r *http.Request){
  decoder := json.NewDecoder(r.Body)

  var movie_data Movie
  err := decoder.Decode(&movie_data)
  if (err != nil){
    panic(err)
  }
  defer r.Body.Close()

  insertResult, err := collection.InsertOne(context.TODO(), movie_data)
  if err != nil{
    log.Fatal(err)
    w.WriteHeader(500)
    return
  }

  fmt.Println("Movie Insert: ", insertResult.InsertedID)
  response(w, 200, movie_data)
}

func MovieUpdate(w http.ResponseWriter, r *http.Request){

  params := mux.Vars(r)

  movie_id, err := primitive.ObjectIDFromHex(params["id"])
  if err != nil {
    w.WriteHeader(404)
    return
  }

  decoder := json.NewDecoder(r.Body)

  var movie_data Movie

  err = decoder.Decode(&movie_data)

  if err != nil{
    panic(err)
    w.WriteHeader(500)
    return
  }

  defer r.Body.Close()

  filter := bson.D{{"_id", movie_id}}

  update := bson.D{{"$set", movie_data}}

  updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
  if err != nil{
    log.Fatal(err)
  }
  response(w, 200, movie_data)

  fmt.Printf("Matched %v documents and updated %v documents. \n", updateResult.MatchedCount, updateResult.ModifiedCount)
}

func MovieList(w http.ResponseWriter, r *http.Request){

  var results []*Movie

  filter := bson.D{{}}
  options := options.Find()


  cur, err := collection.Find(context.TODO(), filter, options)
  if err != nil{
    log.Fatal(err)
  }

  for cur.Next(context.TODO()){
    var elem Movie
    err := cur.Decode(&elem)
    if err != nil{
      log.Fatal(err)
    }
    fmt.Println(elem)
    results = append(results, &elem)
  }

  if err := cur.Err(); err != nil{
    log.Fatal(err)
  }

  cur.Close(context.TODO())

  fmt.Printf("Se encontraron documentos %+v\n", results)

}

func MovieShow(w http.ResponseWriter, r *http.Request){
  //Lee parametros desde url y los carga en el arreglo mux.vars
  params := mux.Vars(r)
  var result Movie

  movie_id := params["id"]

  id, err := primitive.ObjectIDFromHex(movie_id)
  if err != nil{
      log.Fatal(err)
  }

  filter:= bson.D{{"_id", id}}

  err = collection.FindOne(context.TODO(), filter).Decode(&result)
  if err != nil{
    log.Fatal(err)
  }
  response(w, 200, result)
  fmt.Printf("Found a single movie: %v\n", result)
}

func MovieDelete(w http.ResponseWriter, r *http.Request){
  params := mux.Vars(r)

  movie_id, err := primitive.ObjectIDFromHex(params["id"])
  if err != nil{
      w.WriteHeader(404)
      return
  }

  filter := bson.D{{"_id", movie_id}}

  deleteResult, err := collection.DeleteOne(context.TODO(), filter)
  if err != nil{
    log.Fatal(err)
  }

  fmt.Println("Documento eliminado %v \n", deleteResult.DeletedCount)
  results := Message{"success", "El documento ha sido eliminado"}
  w.Header().Set("Content-Type","application/json")
  w.WriteHeader(200)
  json.NewEncoder(w).Encode(results)
}
