package user

import (
    "database/sql"
    "log"
)

var DB *sql.DB

func GetAllUsers() []User {
    rows, err := DB.Query("SELECT id, name, email FROM users")
    if err != nil {
        log.Println(err)
        return []User{}
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var u User
        if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
            log.Println(err)
            continue
        }
        users = append(users, u)
    }

    return users
}

