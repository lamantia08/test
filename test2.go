package main

import (
    "encoding/json"
    "net/http"
     "context"
        "go.mongodb.org/mongo-driver/bson"
        "go.mongodb.org/mongo-driver/bson/primitive"
        "go.mongodb.org/mongo-driver/mongo"
        "go.mongodb.org/mongo-driver/mongo/options"
        "log"
        "time"

)

type Student struct {
        ID    primitive.ObjectID `bson:"_id"`
        Nome  string             `bson:"nome"`
        Mamma string             `bson:"mamma"`
        Padre string             `bson:"padre"`
        Anni  int                `bson:"anni"`
}




func main() {
    // Definisci una funzione per gestire le richieste HTTP sulla route "/"
    http.HandleFunc("/classe", func(w http.ResponseWriter, r *http.Request) {
        // Crea un oggetto Message con un messaggio di risposta

            // Imposta l'URI di connessione
        uri := "mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.0.0"

        // Crea un context con timeout per gestire la connessione
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        // Connessione al database
        client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
        if err != nil {
                log.Fatal(err)
        }
        defer client.Disconnect(ctx)

        // Specifica il nome del database e della collezione
        databaseName := "test"
        collectionName := "classe"

        // Ottieni un handler per la collezione
        collection := client.Database(databaseName).Collection(collectionName)

        // Esegui una query per recuperare i documenti
        filter := bson.M{"anni": bson.M{"$gt": 5}}

        cursor, err := collection.Find(ctx, filter)
        if err != nil {
                log.Fatalf("Errore nella creazione del cursore: %s", err)
        }
        defer cursor.Close(ctx)

        // Itera sui documenti e decodificali
        var students []Student
        for cursor.Next(ctx) {
                var student Student
                if err := cursor.Decode(&student); err != nil {
                        log.Fatal(err, "errore decodifica")
                }
                students = append(students, student)
        }

    









	// Converti l'oggetto Message in formato JSON
        jsonResponse, err := json.Marshal(students)
        if err != nil {
            http.Error(w, "Errore nella creazione della risposta JSON", http.StatusInternalServerError)
            return
        }

        // Imposta l'intestazione della risposta per indicare che è JSON
        w.Header().Set("Content-Type", "application/json")

        // Scrivi la risposta JSON al client
        _, err = w.Write(jsonResponse)
        if err != nil {
            http.Error(w, "Errore nell'invio della risposta JSON", http.StatusInternalServerError)
            return
        }
    })

    // Avvia il server HTTP sulla porta 8080
    if err := http.ListenAndServe(":8080", nil); err != nil {
        panic(err)
    }
}

