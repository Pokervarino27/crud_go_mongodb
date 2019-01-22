package main

//Movie estructura de peliculas
type Movie struct {
	Name     string `json:"name"`
	Year     int    `json:"year"`
	Director string `json:"director"`
}

//Movies arreglo de struct Movie
type Movies []Movie

//Message struct del tipo mensaje de API
type Message struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
