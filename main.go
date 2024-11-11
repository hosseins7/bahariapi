package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    _ "github.com/lib/pq"
)

var db *sql.DB

type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

func main() {
    connStr := "user=hossein dbname=godb sslmode=disable password=123456Hh&"
    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Error connecting to the database: ", err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal("Cannot ping the database: ", err)
    }
    fmt.Println("Successfully connected to PostgreSQL database!")

    router := gin.Default()

    router.GET("/users", getUsers)
    router.POST("/users", createUser)
    router.GET("/users/:id", getUser)
    router.PUT("/users/:id", updateUser)
    router.DELETE("/users/:id", deleteUser)

    router.Run("0.0.0.0:8080")
}

func createUser(c *gin.Context) {
    var user User

    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
    err := db.QueryRow(query, user.Name, user.Email, user.Password).Scan(&user.ID)
    if err != nil {
        log.Printf("Error inserting user: %v", err) // Log the error for debugging
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"id": user.ID, "name": user.Name, "email": user.Email})
}

func getUsers(c *gin.Context) {
    rows, err := db.Query("SELECT id, name, email FROM users")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch users"})
        return
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to scan user"})
            return
        }
        users = append(users, user)
    }

    c.JSON(http.StatusOK, users)
}

func getUser(c *gin.Context) {
    id := c.Param("id")
    var user User

    err := db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.JSON(http.StatusOK, user)
}

func updateUser(c *gin.Context) {
    id := c.Param("id")
    var user User

    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    query := `UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4`
    _, err := db.Exec(query, user.Name, user.Email, user.Password, id)
    if err != nil {
        log.Printf("Error updating user: %v", err) // Log the error for debugging
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func deleteUser(c *gin.Context) {
    id := c.Param("id")

    query := `DELETE FROM users WHERE id = $1`
    _, err := db.Exec(query, id)
    if err != nil {
        log.Printf("Error deleting user: %v", err) // Log the error for debugging
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

