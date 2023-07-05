package mssql

// import (
//   "directpay.online/mno/ecocash-zimbabwe-refund/models"
//   "database/sql"
//   "fmt"
//   "github.com/DATA-DOG/go-sqlmock"
//   "github.com/stretchr/testify/assert"
//   "log"
//   "testing"
// )

// var u = &models.UserEntity{
//   ID:    1,
//   Name:  "Momo",
//   Email: "momo@mail.com",
//   Phone: "08123456789",
// }

// func NewMock() (*sql.DB, sqlmock.Sqlmock) {
//   db, mock, err := sqlmock.New()
//   if err != nil {
//     log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//   }

//   return db, mock
// }

// func TestNewRepositoryBadDriver(t *testing.T) {
//   dbConfig := &models.DBConfig{
//     Dialect: "001sqlmock",
//     Host: "host",
//     Port: 1234,
//     Database: "mock",
//     User: "mockUser",
//     Password: "mockPassword",
//   }
//   connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s",
//      dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Database)
//   // first we create a mock db, under a connString
//   _, _, err := sqlmock.NewWithDSN(connString)
//   if err != nil {
//     panic("Got an unexpected error.")
//   }
//   _, err = NewRepository(dbConfig)
//   if err == nil {
//     t.Error("NewRepository bad driver test failed")
//   }
// }

// func TestNewRepositoryWrongDB(t *testing.T) {
//   dbConfig := &models.DBConfig{
//     Dialect: "sqlmock",
//     Host: "host",
//     Port: 1234,
//     Database: "mock",
//     User: "mockUser",
//     Password: "mockPassword",
//   }
//   connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=wrong",
//     dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Port)
//   // first we create a mock db, under a connString
//   _, _, err := sqlmock.NewWithDSN(connString)
//   if err != nil {
//     panic("Got an unexpected error.")
//   }
//   _, err = NewRepository(dbConfig)
//   if err == nil {
//     t.Error("NewRepository wrong db test failed")
//   }
// }

// func TestNewRepository(t *testing.T) {
//   dbConfig := &models.DBConfig{
//     Dialect: "sqlmock",
//     Host: "host",
//     Port: 1234,
//     Database: "mock",
//     User: "mockUser",
//     Password: "mockPassword",
//   }
//   connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
//     dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Port, dbConfig.Database)
//   // first we create a mock db, under a connString
//   _, _, err := sqlmock.NewWithDSN(connString)
//   if err != nil {
//     fmt.Print(err.Error())
//     panic("Got an unexpected error.")
//   }
//   repo, _ := NewRepository(dbConfig)
//   if repo == nil {
//     t.Error("NewRepository test failed")
//   }
//   repo.Close()
// }

// func TestFindUserByID(t *testing.T) {
//   db, mock := NewMock()
//   repo := &repository{db}
//   defer func() {
//     repo.Close()
//   }()

//   query := "SELECT id, name, email, phone FROM users WHERE id = \\?"

//   rows := sqlmock.NewRows([]string{"id", "name", "email", "phone"}).
//     AddRow(1, u.Name, u.Email, u.Phone)

//   mock.ExpectQuery(query).WithArgs(u.ID).WillReturnRows(rows)

//   user, err := repo.FindUserByID(u.ID)
//   assert.NotNil(t, user)
//   assert.NoError(t, err)
// }

// func TestFindUserByIDError(t *testing.T) {
//   db, mock := NewMock()
//   repo := &repository{db}
//   defer func() {
//     repo.Close()
//   }()

//   query := "SELECT id, name, email, phone FROM user WHERE id = \\?"

//   rows := sqlmock.NewRows([]string{"id", "name", "email", "phone"})

//   mock.ExpectQuery(query).WithArgs(u.ID).WillReturnRows(rows)

//   user, err := repo.FindUserByID(u.ID)
//   assert.Empty(t, user)
//   assert.Error(t, err)
// }

// func TestCreateUser(t *testing.T) {
//   db, mock := NewMock()
//   repo := &repository{db}
//   defer func() {
//     repo.Close()
//   }()

//   query := "INSERT INTO users \\(id, name, email, phone\\) VALUES \\(\\?, \\?, \\?, \\?\\)"

//   prep := mock.ExpectPrepare(query)
//   prep.ExpectExec().WithArgs(u.ID, u.Name, u.Email, u.Phone).WillReturnResult(sqlmock.NewResult(0, 1))

//   err := repo.CreateUser(u)
//   assert.NoError(t, err)
// }

// func TestCreateError(t *testing.T) {
//   db, mock := NewMock()
//   repo := &repository{db}
//   defer func() {
//     repo.Close()
//   }()

//   query := "INSERT INTO user \\(id, name, email, phone\\) VALUES \\(\\?, \\?, \\?, \\?\\)"

//   prep := mock.ExpectPrepare(query)
//   prep.ExpectExec().WithArgs(u.ID, u.Name, u.Email, u.Phone).WillReturnResult(sqlmock.NewResult(0, 0))

//   err := repo.CreateUser(u)
//   assert.Error(t, err)
// }

// func TestUpdateUser(t *testing.T) {
//   db, mock := NewMock()
//   repo := &repository{db}
//   defer func() {
//     repo.Close()
//   }()

//   query := "UPDATE users SET name = \\?, email = \\?, phone = \\? WHERE id = \\?"

//   prep := mock.ExpectPrepare(query)
//   prep.ExpectExec().WithArgs(u.Name, u.Email, u.Phone, u.ID).WillReturnResult(sqlmock.NewResult(0, 1))

//   err := repo.UpdateUser(u)
//   assert.NoError(t, err)
// }

// func TestUpdateUserErr(t *testing.T) {
//   db, mock := NewMock()
//   repo := &repository{db}
//   defer func() {
//     repo.Close()
//   }()

//   query := "UPDATE user SET name = \\?, email = \\?, phone = \\? WHERE id = \\?"

//   prep := mock.ExpectPrepare(query)
//   prep.ExpectExec().WithArgs(u.Name, u.Email, u.Phone, u.ID).WillReturnResult(sqlmock.NewResult(0, 0))

//   err := repo.UpdateUser(u)
//   assert.Error(t, err)
// }
