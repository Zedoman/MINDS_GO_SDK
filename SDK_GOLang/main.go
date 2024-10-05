// package main

// import (
//     "database/sql"
//     "context"
//     "fmt"
//     "log"
//     "time"
//     _ "github.com/go-sql-driver/mysql"
// )

// // MindsDBClient represents a client for MindsDB.
// type MindsDBClient struct {
//     db *sql.DB
// }

// // NewMindsDBClient initializes a new MindsDB client.
// func NewMindsDBClient(dsn string) (*MindsDBClient, error) {
//     db, err := sql.Open("mysql", dsn)
//     if err != nil {
//         return nil, fmt.Errorf("failed to open database: %w", err)
//     }

//     // Test the connection with a timeout
//     ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//     defer cancel()

//     if err := db.PingContext(ctx); err != nil {
//         return nil, fmt.Errorf("failed to ping database: %w", err)
//     }

//     return &MindsDBClient{db: db}, nil
// }

// // Close closes the database connection.
// func (client *MindsDBClient) Close() error {
//     return client.db.Close()
// }

// // QueryPredictors retrieves a limited number of predictors from the database.
// func (client *MindsDBClient) QueryPredictors(limit int) ([]string, error) {
//     query := fmt.Sprintf("SELECT name FROM mindsdb.predictors LIMIT %d;", limit)
//     rows, err := client.db.Query(query)
//     if err != nil {
//         return nil, fmt.Errorf("error executing query: %w", err)
//     }
//     defer rows.Close()

//     var predictors []string
//     for rows.Next() {
//         var name string
//         if err := rows.Scan(&name); err != nil {
//             return nil, fmt.Errorf("error scanning row: %w", err)
//         }
//         predictors = append(predictors, name)
//     }

//     return predictors, nil
// }

// // CreatePredictorsTable creates the predictors table if it doesn't exist.
// func (client *MindsDBClient) CreatePredictorsTable() error {
//     query := `
//     CREATE TABLE IF NOT EXISTS predictors (
//         id INT AUTO_INCREMENT PRIMARY KEY,
//         name VARCHAR(255) NOT NULL,
//         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
//     );`

//     _, err := client.db.Exec(query)
//     return err
// }

// func main() {
//     // Set the Data Source Name (DSN)
//     dsn := "root@tcp(localhost:47334)/mindsdb?timeout=10s"
//     log.Println("Attempting to connect to MindsDB...")

//     client, err := NewMindsDBClient(dsn)
//     if err != nil {
//         log.Fatalf("Error creating MindsDB client: %v", err)
//     }
//     defer client.Close()

//     fmt.Println("Successfully connected to MindsDB!")

//     // Create the predictors table if it doesn't exist
//     if err := client.CreatePredictorsTable(); err != nil {
//         log.Fatalf("Error creating predictors table: %v", err)
//     }

//     // Optional: Insert sample data into the predictors table
//     sampleData := []string{"Predictor 1", "Predictor 2", "Predictor 3"}
//     for _, name := range sampleData {
//         _, err := client.db.Exec("INSERT INTO predictors (name) VALUES (?);", name)
//         if err != nil {
//             log.Fatalf("Error inserting sample data: %v", err)
//         }
//     }

//     // Query data from MindsDB
//     predictors, err := client.QueryPredictors(10)
//     if err != nil {
//         log.Fatalf("Error querying predictors: %v", err)
//     }

//     // Print the query results
//     fmt.Println("Predictors:")
//     for _, predictor := range predictors {
//         fmt.Println(predictor)
//     }
// }


package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
    "github.com/gorilla/mux"
)

// MindsDBClient represents a client for MongoDB.
type MindsDBClient struct {
    collection *mongo.Collection
}

// Predictor represents the structure for predictor.
type Predictor struct {
    ID   string `json:"id" bson:"_id,omitempty"`
    Name string `json:"name" bson:"name"`
}

// NewMindsDBClient initializes a new MongoDB client for MindsDB using MongoDB Atlas.
func NewMindsDBClient(uri string, dbName string, collectionName string) (*MindsDBClient, error) {
    clientOptions := options.Client().ApplyURI(uri)
    client, err := mongo.NewClient(clientOptions)
    if err != nil {
        return nil, fmt.Errorf("failed to create MongoDB client: %v", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err = client.Connect(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
    }

    collection := client.Database(dbName).Collection(collectionName)

    return &MindsDBClient{collection: collection}, nil
}

// CreatePredictor creates a new predictor in the MongoDB collection.
func (client *MindsDBClient) CreatePredictor(predictor Predictor) error {
    _, err := client.collection.InsertOne(context.TODO(), predictor)
    return err
}

// GetPredictors retrieves all predictors from the collection.
func (client *MindsDBClient) GetPredictors() ([]Predictor, error) {
    var predictors []Predictor
    cursor, err := client.collection.Find(context.TODO(), bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.TODO())

    for cursor.Next(context.TODO()) {
        var predictor Predictor
        err := cursor.Decode(&predictor)
        if err != nil {
            return nil, err
        }
        predictors = append(predictors, predictor)
    }
    return predictors, nil
}

// CreatePredictorHandler handles the creation of a predictor via POST request.
func CreatePredictorHandler(client *MindsDBClient, w http.ResponseWriter, r *http.Request) {
    var predictor Predictor
    err := json.NewDecoder(r.Body).Decode(&predictor)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    err = client.CreatePredictor(predictor)
    if err != nil {
        http.Error(w, "Failed to create predictor", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(predictor)
}

// GetPredictorsHandler handles retrieving the list of predictors via GET request.
func GetPredictorsHandler(client *MindsDBClient, w http.ResponseWriter, r *http.Request) {
    predictors, err := client.GetPredictors()
    if err != nil {
        http.Error(w, "Failed to retrieve predictors", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(predictors)
}

func main() {
    // MongoDB Atlas connection string
    uri := "mongodb+srv://<username>:<password>@cluster0.kpxtb.mongodb.net/<dbname>?retryWrites=true&w=majority"

    // Replace with your own credentials and database name
    dbName := "mindsdb"
    collectionName := "predictors"

    client, err := NewMindsDBClient(uri, dbName, collectionName)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    // Set up router
    r := mux.NewRouter()
    r.HandleFunc("/predictors", func(w http.ResponseWriter, r *http.Request) {
        GetPredictorsHandler(client, w, r)
    }).Methods("GET")
    r.HandleFunc("/predictors", func(w http.ResponseWriter, r *http.Request) {
        CreatePredictorHandler(client, w, r)
    }).Methods("POST")

    log.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
