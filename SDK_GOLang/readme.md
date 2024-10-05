<!-- Here is the formatted content for `README.md`:

```markdown
# MindsDB Golang SDK

This project demonstrates how to build a client for MindsDB using Go and MySQL. The client connects to MindsDB, creates a table (if not already created), inserts sample data, and queries the `predictors` table. The solution emphasizes robustness, clarity, and ease of use.

## Features
1. **Database Connection**: Connect to MindsDB using a Data Source Name (DSN) and the MySQL driver.
2. **Table Management**: Create a `predictors` table if it doesn't exist.
3. **Sample Data Insertion**: Insert sample records into the table.
4. **Querying Data**: Retrieve a limited number of records from the `predictors` table.
5. **Error Handling**: Comprehensive error handling for connection failures and SQL execution issues.

## Setup Instructions

### Pre-requisites
- **Go** (version 1.16+)
- **MySQL Driver for Go** (`github.com/go-sql-driver/mysql`)
- **MindsDB** running locally or accessible via a TCP/IP connection

### Install Dependencies
Run the following command to install the MySQL driver:

```bash
go get -u github.com/go-sql-driver/mysql
```

### MindsDB Setup
Ensure that MindsDB is up and running. If using Docker, start the MindsDB container:

```bash
docker run -p 47334:47334 -d mindsdb/mindsdb
```

Confirm that MindsDB is accessible:

```bash
mysql -h localhost -P 47334 -u root
```

### Project Structure
- **`main.go`**: Contains the core logic for interacting with MindsDB.

---

## Code Walkthrough

### Connecting to MindsDB
The `NewMindsDBClient` function establishes a connection with MindsDB using the MySQL driver. It also tests the connection with a 10-second timeout.

```go
func NewMindsDBClient(dsn string) (*MindsDBClient, error) {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := db.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    return &MindsDBClient{db: db}, nil
}
```

---

### Creating the Predictors Table
The `CreatePredictorsTable` function checks if the `predictors` table exists. If not, it creates the table with fields for `id`, `name`, and timestamps.

```go
func (client *MindsDBClient) CreatePredictorsTable() error {
    query := `
    CREATE TABLE IF NOT EXISTS predictors (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    );`

    _, err := client.db.Exec(query)
    return err
}
```

---

### Inserting Sample Data
Once the table is created, the program inserts sample records into the `predictors` table.

```go
sampleData := []string{"Predictor 1", "Predictor 2", "Predictor 3"}
for _, name := range sampleData {
    _, err := client.db.Exec("INSERT INTO predictors (name) VALUES (?);", name)
    if err != nil {
        log.Fatalf("Error inserting sample data: %v", err)
    }
}
```

---

### Querying the Predictors
The `QueryPredictors` function retrieves a limited number of records from the `predictors` table.

```go
func (client *MindsDBClient) QueryPredictors(limit int) ([]string, error) {
    query := fmt.Sprintf("SELECT name FROM mindsdb.predictors LIMIT %d;", limit)
    rows, err := client.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("error executing query: %w", err)
    }
    defer rows.Close()

    var predictors []string
    for rows.Next() {
        var name string
        if err := rows.Scan(&name); err != nil {
            return nil, fmt.Errorf("error scanning row: %w", err)
        }
        predictors = append(predictors, name)
    }

    return predictors, nil
}
```

---

## Running the Project

### Build the Project
Compile the Go program using the following command:

```bash
go build -o mindsdb-client
```

### Run the Program
Execute the binary:

```bash
./mindsdb-client
```

### Expected Output
You should see output similar to:

```bash
Successfully connected to MindsDB!
Predictors:
Predictor 1
Predictor 2
Predictor 3
```

---

## Error Handling
The project uses structured error handling for various steps, including database connection, table creation, data insertion, and query execution. The program logs any errors and exits if a critical issue is encountered.

---

## Testing

1. **Connection Timeout**: The connection uses a 10-second timeout to ensure the program doesn't hang indefinitely if MindsDB is unreachable.
2. **SQL Query Validation**: The program checks for SQL execution errors and logs any issues encountered during data manipulation.
3. **End-to-End Query**: After inserting data into the `predictors` table, the query ensures that the data is correctly stored and retrieved.

---

## Next Steps

- **Enhancements**: Consider adding more advanced queries or predictive capabilities from MindsDB.
- **Scaling**: Implement retry logic and connection pooling to handle larger datasets or more complex queries.
- **Testing**: Add unit tests to cover the critical paths, including table creation and data querying.

This project offers a foundational understanding of how to interact with MindsDB using Go and sets the stage for building more advanced machine learning applications on top of MindsDB. -->



# MindsDB GO SDK

This project implements a simple REST API for managing **Predictors** using Go and MongoDB Atlas. The API allows users to create and retrieve predictors from a MongoDB collection.

## Prerequisites

To run this project, ensure that you have the following tools installed on your machine:

- Go (version 1.16 or above)
- MongoDB Atlas account and cluster (or a MongoDB URI for connection)
- Git (to clone the repository)

## Features

- **Create a Predictor**: Add a new predictor to the MongoDB collection.
- **Retrieve Predictors**: Get a list of all predictors stored in the MongoDB collection.

## Endpoints

### 1. **Create a Predictor**

- **Endpoint**: `POST /predictors`
- **Description**: Add a new predictor to the MongoDB collection.
- **Request Body** (JSON format):
  ```json
  {
    "name": "Predictor Name"
  }
  ```
- **Response**:
  - `201 Created` on success with the newly created predictor in the response body.

- **Example cURL Command**:
  ```bash
  curl -X POST http://localhost:8080/predictors \
  -H "Content-Type: application/json" \
  -d '{"name": "Predictor 1"}'
  ```

### 2. **Retrieve Predictors**

- **Endpoint**: `GET /predictors`
- **Description**: Retrieve all predictors from the MongoDB collection.
- **Response** (JSON format):
  ```json
  [
    {
      "id": "some_id",
      "name": "Predictor 1"
    },
    {
      "id": "some_other_id",
      "name": "Predictor 2"
    }
  ]
  ```

- **Example cURL Command**:
  ```bash
  curl http://localhost:8080/predictors
  ```

## Project Setup and Installation

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/mindsdb-predictors-api.git
cd mindsdb-predictors-api
```

### 2. Install Dependencies

Run the following command to install the Go dependencies required for the project:

```bash
go mod tidy
```

### 3. Update MongoDB URI

In the `main.go` file, update the `uri` with your MongoDB Atlas connection string. For example:

```go
uri := "mongodb+srv://<username>:<password>@cluster0.kpxtb.mongodb.net/<dbname>?retryWrites=true&w=majority"
```

### 4. Run the Application

Start the server by running:

```bash
go run main.go
```

You should see a message like:

```
Server is running on port 8080
```

### 5. Testing the API

Use tools like **Postman**, **Insomnia**, or **cURL** to test the API.

- To **create a new predictor**, send a `POST` request to `http://localhost:8080/predictors`.
- To **retrieve the list of predictors**, send a `GET` request to `http://localhost:8080/predictors`.

## Code Overview

### `main.go`

The main Go file that defines the API and MongoDB client.

- **MindsDBClient**: Represents a MongoDB client connected to the specified collection.
- **Predictor**: A struct that defines the schema for predictors, containing an ID and a Name.
- **NewMindsDBClient**: Function to initialize and connect to MongoDB Atlas using a provided URI.
- **CreatePredictorHandler**: HTTP handler for adding a new predictor via `POST` request.
- **GetPredictorsHandler**: HTTP handler for retrieving predictors via `GET` request.

### Dependencies

- `go.mongodb.org/mongo-driver/mongo`: MongoDB driver for Go.
- `github.com/gorilla/mux`: A powerful router for handling HTTP requests.

## MongoDB Setup

1. Create a free MongoDB Atlas account at [MongoDB Atlas](https://www.mongodb.com/cloud/atlas).
2. Create a new cluster and database.
3. Obtain your MongoDB connection string (URI) from the MongoDB Atlas dashboard.
4. Ensure that your cluster is properly configured with IP whitelisting and credentials.



## License

This project is licensed under the MIT License.
