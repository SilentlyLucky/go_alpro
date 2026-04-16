# Materi Backend Camin ALPRO

---

## Daftar Isi

- [Environment Setup](#environment-setup)
    -  [Instalasi Go](#instalasi-go)
    -  [Setup Project Go](#setup-project-go)
- [Go Crash Course](#go-crash-course)
    - [Variabel & Tipe Data](#variabel--tipe-data)
    - [Array & Slices](#array--slices)
    - [Fungsi](#fungsi)
    - [Struct](#struct)
    - [Error Handling](#error-handling)
- [Backend di Go](#backend-di-go)
    - [Gin Web Framework](#gin-web-framework)
    - [Gorm Library](#gorm-library)
    - [Contoh Boilerplate](#contoh-boilerplate)
- [Latihan Implementasi](#latihan-implementasi)
- [Teori-teori Backend Lanjutan](#teori-teori-backend-lanjutan)

## Environment Setup
### Instalasi Go

Kamu dapat mendownload Go di website resmi mereka [disini](https://go.dev/doc/install).

<img width="1862" height="926" alt="image" src="https://github.com/user-attachments/assets/2aa526f3-1bf1-4770-82c5-a07a520b0c1a" />

Setelah penginstallan selesai, kamu dapat mengecek jika berhasil di terminal, dan menjalankan command berikut:

```bash
go version
# Output yang diharapkan: go version go1.21.x linux/amd64 (atau sejenisnya)
```

### Setup Project Go

Setelah Go berhasil diinstall, kamu dapat membuat project dengan menjalankan command di bawah:

```bash
# Buat folder project
mkdir go-workshop && cd go-workshop

# Inisialisasi Go module
go mod init {username kamu}/go-workshop
```

`go mod init` bertanggung jawab untuk membuat file `go.mod`, yang dapat kamu anggap seperti pelekat project kamu dengan dependency di dalam ataupun di luar project.  

Isi di dalam `go.mod` kamu harusnya ini.
```
module github.com/username/go-workshop

go 1.26.2
```

Untuk mulai menulis kode kamu, kamu dapat membuat file baru `main.go` sebagai file mulainya program kamu.  
Sebagai pemanasan, kita coba print Hello world sebagai tradisi mempelajari bahasa baru!

```go
// main.go
package main

import "fmt"

func main() {
    fmt.Println("Hello world!")
}
```

```bash
go run main.go
# Output: Hello world!
```

## Go Crash Course

### Variabel & Tipe Data

Go memiliki dua cara mendeklarasikan variabel:

```go
// Cara 1: deklarasi eksplisit dengan 'var'
var nama string = "Budi"
var umur int = 25

// Cara 2: short declaration dengan ':='
nama := "Budi"
umur := 25
isActive := true
```

Kamu juga dapat mendeklarasikan 2 variabel seperti di bahasa Python

```go
result, err := getName()
```

### Array & Slices

Kamu dapat mendeklarasikan array seperti berikut:

```go
angka = [7]int{1, 2, 3, 4, 5, 6, 7}

// Kamu dapat mendapatkan nilai-nilai dalam array seperti berikut
fmt.Println(angka[3])

// Mendapatkan dari index 1 ke 3.
fmt.Println(angka[1:4])

// Mendapatkan dari index 0 ke 1.
fmt.Println(angka[:2])
```

Array dalam Go bersifat statis dan besarnya tidak bisa diganti.

Kamu juga dapat menyimpan cuplikan dari array dalam bentuk slice. Slice berupa reference ke array sebenarnya, dan seluruh perubahan pada slice akan berlaku di arraynya, dan vice versa.

```go
a = angka[0:2]
b = angka[1:4]
fmt.Println(a, b)

a[1] = 10
fmt.Println(a, b)
```

### Fungsi

```go
// Fungsi tanpa return value
func sapa(nama string) {
    fmt.Println("Halo,", nama)
}

// Fungsi dengan return value
func tambah(a int, b int) int {
    return a + b
}

// Fungsi dengan multiple return (idiom khas Go)
func bagi(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("tidak bisa dibagi nol")
    }
    return a / b, nil
}
```

### Struct

Go bukan bahasa yang mendukung object oriented programming. Sebagai gantinya, gunakan `struct` untuk mengelompokkan data.

```go
// Mendefinisikan struct
type User struct {
    ID    uint
    Name  string
    Email string
}

// Membuat instance dari struct
user := User{
    ID:    1,
    Name:  "Budi",
    Email: "budi@email.com",
}

fmt.Println(user.Name) // Output: Budi
```

### Error Handling

Go tidak menggunakan `try-catch`. Error dikembalikan sebagai nilai return biasa.  

Go memiliki aturan yang sangat ketat yang ditetapkan, salah satunya adalah **tidak bolehnya ada variabel yang tidak dipakai**. Sebagai efeknya, error yang dikembalikan wajib diurus oleh programmer, jika tidak, program Go tidak akan jalan!

```go
result, err := bagi(10, 0)
if err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println("Hasil:", result)
```

Untuk menambah pengetahuan mengenai Go, dapat mengikuti [tour ini](https://go.dev/tour/list).


## Backend Di Go

Go adalah salah satu bahasa yang paling populer dalam backend development, dengan alasan-alasan berikut:
- Go cukup gampang untuk dimengerti, sehingga kode-kode yang dibuat lebih readable dan gampang dimaintain
- Go adalah bahasa yang sangat _opinionated_ sampai aturan formattingnya ditetapkan secara universal, sehingga bahkan seorang BE developer yang tidak pernah menyentuh codebase, dapat mengerti control flow nya dengan cepat.
- Fitur-fitur dasar Go sangat cocok untuk backend development, berupa `net/http` dan `goroutines`, dan sudah cukup robust untuk production, dan tidak terlalu perlu memanggil dependency dari library eksternal.
- Go sangat gampang untuk dideploy, karena semua hampir dependency dapat diinstall dengan Go sebagai "package manager" nya.

Walau `net/http` cukup untuk membangun projek Backend yang solid, kita akan memakai framework web `gin` untuk menghindari _reinventing the wheel_.

### Gin Web Framework

Pada singkatnya, Gin adalah framework web yang membuat routing, parsing request, dan penulisan response JSON menjadi jauh lebih ringkas.

Untuk menginstall gin, jalankan command ini di direktori projek Go kalian.
```bash
# Install Gin
go get github.com/gin-gonic/gin
```

Kamu dapat merubah `main.go` kamu menjadi dibawah untuk mencoba apakah Gin berhasil diinstall.
```go
package main

import (
  "net/http"

  "github.com/gin-gonic/gin"
)

func main() {
  // Create a Gin router with default middleware (logger and recovery)
  r := gin.Default()

  // Define a simple GET endpoint
  r.GET("/ping", func(c *gin.Context) {
    // Return JSON response
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })

  // Start server on port 8080 (default)
  // Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
  r.Run()
}
```

Jalankan dengan `go run main.go`, lalu buka browser ke `GET http://localhost:8080/ping`. Hasilnya:

```json
{
  "message": "pong"
}
```

---

### Gorm Library

ORM (Object-Relational Mapper) adalah alat yang menjembatani kode Go dengan database. Dengan GORM, kamu tidak perlu menulis SQL secara manual untuk operasi CRUD dasar.

Untuk menginstall gorm, jalankan command ini di direktori projek Go kalian. Workshop ini akan memakai `postgres` sebagai sistem databasenya.
```bash
# Install GORM dan driver PostgreSQL
go get gorm.io/gorm
go get gorm.io/driver/postgres
```

GORM memetakan `struct` Go menjadi tabel di database secara otomatis.:

```go
type User struct {
    ID       uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
    Name     string    `gorm:"not null" json:"name"`
    Email    string    `gorm:"unique;not null" json:"email"`
    Password string    `gorm:"not null" json:"-"` // tidak muncul di JSON response
}
```

Dalam struct GORM, kamu dapat memberikan atribut ke baris di entity dengan menambahkan field tags.  
Untuk field tags yang lengkap, dapat dilihat di [dokumentasi GORM](https://gorm.io/docs/models.html#Fields-Tags)

```go
package main

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name     string    `gorm:"not null" json:"name"`
	Email    string    `gorm:"unique;not null" json:"email"`
	Password string    `gorm:"not null" json:"-"` // tidak muncul di JSON response
}

func main() {
	// Harusnya menggunakan .env but for demonstration only
	dsn := "host=localhost user=username password=password dbname=go_workshop port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Karena user menggunakan uuid
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	db.AutoMigrate(&User{})

	newUser := User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashed_password_here",
	}
	db.Create(&newUser)
	fmt.Println("Created User ID:", newUser.ID)

	var user User
	db.First(&user, "id = ?", newUser.ID)
	fmt.Println("User:", user)

	db.Delete(&user, "id = ?", user.ID)
	fmt.Println("Deleted user!")
}
```

Perlu diingat bahwa kamu perlu mensetup Postgres kamu terlebih dahulu:
- Buat user di postgres `(CREATE USER your_username WITH PASSWORD 'your_password';)`
- Buat database di postgres `(CREATE DATABASE your_dbname;)`
- Berikan user akses ke database tersebut `(ALTER DATABASE your_dbname OWNER TO your_username;)`

---

Jika kita menggabungkan kedua library ini, bentuk `main.go` akan seperti ini:

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name     string    `gorm:"not null" json:"name"`
	Email    string    `gorm:"unique;not null" json:"email"`
	Password string    `gorm:"not null" json:"-"` // Does not appear in JSON response
}

func main() {
	// Harusnya menggunakan .env but for demonstration only
	dsn := "host=localhost user=postgres password=root dbname=go_workshop port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Karena user menggunakan uuid
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	db.AutoMigrate(&User{})
	fmt.Println("Database connected and migrated successfully!")

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.POST("/users", func(c *gin.Context) {
		var input struct {
			Name     string `json:"name" binding:"required"`
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newUser := User{
			Name:     input.Name,
			Email:    input.Email,
			Password: input.Password, // Di dalam praktek, jangan lupa untuk menghash password terlebih dahulu!
		}

		if err := db.Create(&newUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "data": newUser})
	})

	r.GET("/users", func(c *gin.Context) {
		var users []User

		if err := db.Find(&users).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": users})
	})

	r.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		if err := db.Delete(&User{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	})

	fmt.Println("Server is running on http://localhost:8080")
	r.Run(":8080")
}
```

Melihat kode `main.go`, hasilnya cukup panjang dan ribet bukan? Dan ini baru satu API CRUD entity.  
Oleh karena itu, walaupun tidak diwajibkan oleh bahasa, kita perlu struktur project yang jelas dan terstruktur.

### Contoh Boilerplate

Boilerplate singkat kata adalah template yang sering dipakai berulang-ulang dan biasanya sedikit / tidak ada variasi, yang bertujuan untuk memberikan struktur ke project kamu.  

Terdapat beberapa boilerplate backend Go yang tersedia, tetapi di workshop ini kita akan menggunakan [boilerplate ini](https://github.com/Caknoooo/go-gin-clean-starter/tree/main).

Berikut adalah struktur folder boilerplate yang kita gunakan, beserta penjelasan tanggung jawab masing-masing:

```
.
├── cmd/
│   └── main.go                        # Entry point — hanya panggil providers & jalankan server
│
├── config/
│   ├── database.go                    # Setup koneksi GORM ke PostgreSQL
│   ├── email.go                       # Konfigurasi SMTP untuk kirim email
│   ├── logger.go                      # Setup query logger ke file
│   └── logs/query_log                 # Output log query SQL
│
├── database/
│   ├── entities/                      # Definisi struct yang di-mapping ke tabel DB
│   │   ├── common.go                  # Base struct (ID, CreatedAt, UpdatedAt, dll)
│   │   ├── user_entity.go             # Tabel `users`
│   │   └── refresh_token_entity.go    # Tabel `refresh_tokens`
│   ├── migrations/                    # Versi skema database (seperti git untuk DB)
│   │   ├── 20240101000000_create_users_table.go
│   │   └── 20240101000001_create_refresh_tokens_table.go
│   ├── seeders/                       # Data awal untuk development/testing
│   │   ├── json/users.json            # Data seed dalam format JSON
│   │   └── seeds/user_seed.go         # Logic untuk insert seed data
│   ├── manager.go                     # Orchestrator: jalankan migrate & seed
│   ├── migration.go                   # Runner untuk file-file migration
│   └── seeder.go                      # Runner untuk file-file seeder
│
├── middlewares/
│   ├── authentication.go              # Validasi JWT token di setiap protected route
│   └── cors.go                        # Izinkan/blokir request lintas domain
│
├── modules/                           # ← FOKUS UTAMA workshop
│   ├── auth/                          # Semua yang berkaitan dengan login/logout/token
│   │   ├── controller/auth_controller.go
│   │   ├── dto/auth_dto.go
│   │   ├── repository/refresh_token_repository.go
│   │   ├── service/auth_service.go
│   │   ├── service/jwt_service.go
│   │   ├── validation/auth_validation.go
│   │   ├── tests/auth_validation_test.go
│   │   └── routes.go
│   └── user/                          # Semua yang berkaitan dengan data user
│       ├── controller/user_controller.go
│       ├── dto/user_dto.go
│       ├── query/user_query.go
│       ├── repository/user_repository.go
│       ├── service/user_service.go
│       ├── validation/user_validation.go
│       ├── tests/user_validation_test.go
│       └── routes.go
│
├── pkg/
│   ├── constants/common.go            # Konstanta global (pesan error, status, dll)
│   ├── helpers/password.go            # Helper bcrypt: hash & compare password
│   └── utils/
│       ├── aes.go                     # Enkripsi/dekripsi data sensitif
│       ├── email.go                   # Fungsi kirim email via SMTP
│       ├── file.go                    # Helper upload & manajemen file
│       └── response.go                # Standarisasi format JSON response
│
├── providers/
│   └── core.go                        # Dependency injection: wiring semua layer
│
├── script/
│   ├── command.go                     # Definisi command CLI (migrate, seed, dll)
│   └── script.go                      # Runner untuk perintah dari terminal
│
├── docker/
│   ├── Dockerfile                     # Build image untuk production
│   ├── nginx/default.conf             # Konfigurasi reverse proxy
│   └── postgresql/                    # Konfigurasi PostgreSQL container
│
├── docker-compose.yml                 # Jalankan seluruh stack (App + DB + Nginx)
├── Makefile                           # Shortcut command (make migrate, make run, dll)
└── create_module.sh                   # Script otomatis buat module baru
```

#### Alur Request — Dari HTTP ke Database

Setiap request yang masuk melewati lapisan-lapisan ini secara berurutan:

```
HTTP Request
    │
    ▼
[Middleware]           → Auth check (JWT valid?), CORS header
    │
    ▼
[Controller]           → Terima request Gin, panggil Validation, kirim response
    │
    ▼
[Validation]           → Validasi input (field required, format email, dll)
    │
    ▼
[Service]              → Business logic (hash password, generate token, dll)
    │
    ▼
[Repository]           → Query ke database via GORM
    │
    ▼
[Database / Entity]    → PostgreSQL
```

#### Pola Per-Module

Setiap fitur baru dibuat dalam satu folder `modules/<nama_fitur>/` yang memiliki struktur seragam:

| File | Tanggung Jawab |
|---|---|
| `controller/` | Terima `*gin.Context`, parsing input, kirim JSON response |
| `dto/` | Struct untuk request body & response (Data Transfer Object) |
| `validation/` | Aturan validasi input sebelum masuk ke service |
| `service/` | Business logic — tidak boleh tahu soal HTTP atau database secara langsung |
| `repository/` | Semua query GORM — satu-satunya layer yang boleh sentuh DB |
| `query/` | Raw query atau filter kompleks yang dipakai oleh repository |
| `routes.go` | Daftarkan semua endpoint milik module ini |
| `tests/` | Unit test untuk validation & service |

> **Pola ini membuat kode mudah ditemukan.** Kalau ada bug di response format, cari di `controller`. Kalau ada bug di kalkulasi bisnis, cari di `service`. Kalau ada bug di query lambat, cari di `repository`.

#### Cara Jalankan Project

```bash
# Clone dan install dependency
go mod tidy

# Jalankan seluruh stack dengan Docker
docker-compose up -d

# Jalankan migrasi database
make migrate
# atau: go run script/script.go migrate

# Jalankan seeder (isi data awal)
make seed
# atau: go run script/script.go seed

# Jalankan development server
go run cmd/main.go
```

---

## Phase 4 — Hands-on Implementation
### `01:00 - 01:30`

Teori sudah cukup. Sekarang waktunya tangan kotor dengan kode sungguhan di dalam boilerplate. Semua file yang dibuat mengikuti pola yang sudah ada di module `auth` dan `user`.

### 01:00 - 01:10 | Demo Mentor: `POST /users` (Create User)

Mentor mendemonstrasikan alur pembuatan satu endpoint utuh, dari mendaftarkan route sampai data tersimpan ke database.

**`database/entities/user_entity.go`** — Skema tabel
```go
package entities

type User struct {
    Common                    // embed ID, CreatedAt, UpdatedAt, DeletedAt
    Name     string `gorm:"not null"`
    Email    string `gorm:"unique;not null"`
    Password string `gorm:"not null"`
    Role     string `gorm:"default:'user'"`
}
```

**`modules/user/dto/user_dto.go`** — Shape data masuk & keluar
```go
package dto

type CreateUserRequest struct {
    Name     string `json:"name"     binding:"required"`
    Email    string `json:"email"    binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

type UserResponse struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Role  string `json:"role"`
}
```

**`modules/user/validation/user_validation.go`** — Validasi input
```go
package validation

import (
    "github.com/gin-gonic/gin"
    "github.com/username/boilerplate/modules/user/dto"
)

func ValidateCreateUser(c *gin.Context) (*dto.CreateUserRequest, error) {
    var req dto.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        return nil, err
    }
    return &req, nil
}
```

**`modules/user/repository/user_repository.go`** — Query database
```go
package repository

import (
    "github.com/username/boilerplate/database/entities"
    "gorm.io/gorm"
)

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *entities.User) error {
    return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*entities.User, error) {
    var user entities.User
    err := r.db.First(&user, id).Error
    return &user, err
}

func (r *UserRepository) FindAll() ([]entities.User, error) {
    var users []entities.User
    err := r.db.Find(&users).Error
    return users, err
}
```

**`modules/user/service/user_service.go`** — Business logic
```go
package service

import (
    "github.com/username/boilerplate/database/entities"
    "github.com/username/boilerplate/modules/user/dto"
    "github.com/username/boilerplate/modules/user/repository"
    "github.com/username/boilerplate/pkg/helpers"
)

type UserService struct {
    repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) CreateUser(req *dto.CreateUserRequest) (*entities.User, error) {
    // Business logic: hash password sebelum disimpan
    hashedPassword, err := helpers.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }

    user := &entities.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: hashedPassword,
    }

    err = s.repo.Create(user)
    return user, err
}
```

**`modules/user/controller/user_controller.go`** — Handle HTTP
```go
package controller

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/username/boilerplate/modules/user/service"
    "github.com/username/boilerplate/modules/user/validation"
    "github.com/username/boilerplate/pkg/utils"
)

type UserController struct {
    service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
    return &UserController{service: service}
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
    req, err := validation.ValidateCreateUser(c)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    user, err := ctrl.service.CreateUser(req)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat user")
        return
    }

    utils.SuccessResponse(c, http.StatusCreated, "User berhasil dibuat", user)
}
```

**`modules/user/routes.go`** — Daftarkan endpoint
```go
package user

import (
    "github.com/gin-gonic/gin"
    "github.com/username/boilerplate/middlewares"
    "github.com/username/boilerplate/modules/user/controller"
)

func RegisterUserRoutes(r *gin.RouterGroup, ctrl *controller.UserController) {
    users := r.Group("/users")
    {
        users.POST("", ctrl.CreateUser)                                          // POST /api/users
        users.GET("", middlewares.Authentication(), ctrl.GetAllUsers)            // GET  /api/users (protected)
        users.GET("/:id", middlewares.Authentication(), ctrl.GetUserByID)        // GET  /api/users/:id (protected)
    }
}
```

> **Perhatikan** `middlewares.Authentication()` — endpoint yang membutuhkan login dibungkus middleware ini. Middleware akan memvalidasi JWT token sebelum request diteruskan ke controller.

---

### 01:10 - 01:30 | Challenge Peserta

Setelah demo selesai, peserta mengerjakan tantangan berikut **secara mandiri** mengikuti pola yang sudah dicontohkan.

**Challenge A — `GET /users/:id`**

> Ambil satu user berdasarkan ID. Kembalikan `404` jika tidak ditemukan.

Tambahkan di `repository`:
```go
func (r *UserRepository) FindByID(id uint) (*entities.User, error) {
    var user entities.User
    err := r.db.First(&user, id).Error
    return &user, err
}
```

Tambahkan di `service`:
```go
func (s *UserService) GetUserByID(id uint) (*entities.User, error) {
    return s.repo.FindByID(id)
}
```

Tambahkan di `controller`:
```go
func (ctrl *UserController) GetUserByID(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
        return
    }

    user, err := ctrl.service.GetUserByID(uint(id))
    if err != nil {
        utils.ErrorResponse(c, http.StatusNotFound, "User tidak ditemukan")
        return
    }

    utils.SuccessResponse(c, http.StatusOK, "OK", user)
}
```

**Challenge B — `GET /users`**

> Ambil semua user. Kembalikan array JSON.

```go
// repository
func (r *UserRepository) FindAll() ([]entities.User, error) {
    var users []entities.User
    err := r.db.Find(&users).Error
    return users, err
}

// service
func (s *UserService) GetAllUsers() ([]entities.User, error) {
    return s.repo.FindAll()
}

// controller
func (ctrl *UserController) GetAllUsers(c *gin.Context) {
    users, err := ctrl.service.GetAllUsers()
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data")
        return
    }
    utils.SuccessResponse(c, http.StatusOK, "OK", users)
}
```

> **Tips debugging:** Error paling umum di Go adalah lupa menangani `if err != nil`. Kalau ada `panic`, cari baris yang tidak menangani error return-nya. Gunakan `log.Println(err)` untuk print error ke terminal.

---

## Phase 5 — Backend 101: Fundamental Concepts
### `01:30 - 02:00`

*Cooling down* dari sesi coding. Phase ini membahas konsep arsitektur yang wajib dipahami oleh setiap backend engineer, terlepas dari bahasa pemrograman yang digunakan.

### 01:30 - 01:45 | Jenis-Jenis Database

Tidak ada satu jenis database yang cocok untuk semua kebutuhan. Pilihan database harus disesuaikan dengan karakteristik data dan kebutuhan sistem.

#### Relational Database (SQL)

**Contoh:** PostgreSQL, MySQL, SQLite

Menyimpan data dalam bentuk tabel dengan baris dan kolom. Hubungan antar tabel didefinisikan secara eksplisit.

| Kapan digunakan | Contoh kasus |
|---|---|
| Data terstruktur dengan skema yang jelas | Data transaksi keuangan |
| Butuh ACID (Atomicity, Consistency, Isolation, Durability) | Sistem pemesanan tiket |
| Ada relasi kompleks antar entitas | Sistem manajemen pengguna |

```sql
-- Contoh query SQL
SELECT users.name, orders.total
FROM users
JOIN orders ON users.id = orders.user_id
WHERE users.id = 1;
```

#### NoSQL Database

**Contoh:** MongoDB (document), Redis (key-value), Cassandra (wide-column)

Menyimpan data dalam format yang lebih fleksibel — bisa dokumen JSON, pasangan key-value, graph, dll.

| Kapan digunakan | Contoh kasus |
|---|---|
| Skema data sering berubah | Katalog produk e-commerce |
| Volume data sangat besar dan perlu scale horizontal | Social media feed |
| Kebutuhan read/write sangat cepat (caching) | Session storage, rate limiting |

> **Aturan praktis:** Mulai dengan SQL. Pindah ke NoSQL hanya jika ada alasan teknis yang jelas, bukan karena terasa "lebih modern".

---

### 01:45 - 01:55 | Authentication: Session vs JWT

Authentication adalah proses verifikasi identitas: *"Apakah kamu benar-benar siapa yang kamu klaim?"*

#### Session-Based Authentication (Stateful)

Server menyimpan informasi sesi di memorinya (atau di database/Redis). Setiap request dari client membawa session ID, dan server mencarinya di penyimpanan.

```
[Client]  →  Login dengan username/password
[Server]  →  Buat session, simpan di memory/Redis
[Server]  →  Kirim session_id ke client (biasanya via cookie)
[Client]  →  Setiap request berikutnya bawa cookie session_id
[Server]  →  Lookup session_id, verifikasi, proses request
```

| ✅ Kelebihan | ❌ Kekurangan |
|---|---|
| Mudah di-revoke (hapus session dari server) | Stateful: server harus simpan data session |
| Cocok untuk aplikasi monolitik | Sulit di-scale horizontal (session harus di-share antar server) |

#### JWT — JSON Web Token (Stateless)

Server tidak menyimpan apapun. Semua informasi yang dibutuhkan ada di dalam token yang dikirim oleh client.

```
[Client]  →  Login dengan username/password
[Server]  →  Generate JWT, tanda-tangani dengan secret key
[Server]  →  Kirim JWT ke client
[Client]  →  Simpan JWT (localStorage atau httpOnly cookie)
[Client]  →  Setiap request kirim JWT di header: Authorization: Bearer <token>
[Server]  →  Verifikasi signature JWT dengan secret key — tidak perlu database lookup
```

Struktur JWT terdiri dari tiga bagian yang dipisahkan titik:
```
header.payload.signature

eyJhbGciOiJIUzI1NiJ9.eyJ1c2VySWQiOjF9.SflKxwRJSMeKKF2QT4fwpMeJf36P
```

| ✅ Kelebihan | ❌ Kekurangan |
|---|---|
| Stateless: mudah di-scale | Sulit di-revoke sebelum expired |
| Tidak butuh database lookup untuk verifikasi | Payload ter-encode (bukan terenkripsi) — jangan simpan data sensitif |
| Bisa digunakan lintas domain/service | Token besar → overhead di setiap request |

> **Rekomendasi:** JWT cocok untuk microservices dan API publik. Session cocok untuk web app tradisional di mana kontrol penuh atas session lebih penting.

---

### 01:55 - 02:00 | Q&A & Wrap Up

Pertanyaan yang sering muncul di akhir sesi ini:

**Q: Kapan saya harus pakai Gin vs framework lain seperti Echo atau Fiber?**
A: Gin adalah pilihan yang aman untuk pemula — dokumentasinya banyak dan komunitasnya besar. Echo dan Fiber juga baik; perbedaannya minor untuk project skala kecil-menengah.

**Q: Apakah Go wajib untuk backend?**
A: Tidak. Go unggul di performa, concurrency, dan binary deployment. Tapi Python, Node.js, atau Java juga valid tergantung kebutuhan tim dan project.

**Q: Apa langkah selanjutnya setelah workshop ini?**
A: Lihat bagian referensi di bawah.

---

<div align="center">

*Workshop ini adalah awal, bukan akhir. Yang paling penting adalah terus membangun sesuatu.*

**Selamat ngoding!** 🎉

</div>
