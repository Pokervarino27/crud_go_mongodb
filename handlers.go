package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

func getClient() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), "mongodb://localhost:27017")

	if err != nil {
		log.Fatal(err)
	}
	//Chequea la conexión
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Mongo")
	return client
}

var collection = getClient().Database("testDb").Collection("movies")

func response(w http.ResponseWriter, status int, results Movie) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(results)
}

//Index funcion de saludo de inicio
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hola, este es el inicio")
}

//MovieAdd recibe un json que añade una pelicula
func MovieAdd(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var movieData Movie
	err := decoder.Decode(&movieData)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	insertResult, err := collection.InsertOne(context.TODO(), movieData)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(500)
		return
	}

	fmt.Println("Movie Insert: ", insertResult.InsertedID)
	response(w, 200, movieData)
}

//MovieUpdate Recibe una actualizacion de una pelicula y la actualiza
func MovieUpdate(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	movieID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		w.WriteHeader(404)
		return
	}

	decoder := json.NewDecoder(r.Body)

	var movieData Movie

	err = decoder.Decode(&movieData)

	if err != nil {
		panic(err)
		w.WriteHeader(500)
		return
	}

	defer r.Body.Close()

	filter := bson.D{{"_id", movieID}}

	update := bson.D{{"$set", movieData}}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	response(w, 200, movieData)

	fmt.Printf("Matched %v documents and updated %v documents. \n", updateResult.MatchedCount, updateResult.ModifiedCount)
}

//MovieList Lista todas las peliculas disponibles.
func MovieList(w http.ResponseWriter, r *http.Request) {

	var results []*Movie

	filter := bson.D{{}}
	options := options.Find()

	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem Movie
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(elem)
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	fmt.Printf("Se encontraron documentos %+v\n", results)

}

//MovieShow muestra una película, seleccionada por ID
func MovieShow(w http.ResponseWriter, r *http.Request) {
	//Lee parametros desde url y los carga en el arreglo mux.vars
	params := mux.Vars(r)
	var result Movie

	movieID := params["id"]

	id, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.D{{"_id", id}}

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	response(w, 200, result)
	fmt.Printf("Found a single movie: %v\n", result)
}

//MovieDelete Elimina una pelicula a través del ID
func MovieDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	movieID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		w.WriteHeader(404)
		return
	}

	filter := bson.D{{"_id", movieID}}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Documento eliminado %v \n", deleteResult.DeletedCount)
	results := Message{"success", "El documento ha sido eliminado"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(results)
}
