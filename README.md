# Backend Go 

---

Readme ini menjelaskan struktur code project, cara kerjanya, lalu membahas Challenge A dan B secara terpisah.

> Built with ❤️ for **Tugas Day 2 Workshop Camin Alpro 2026**
---
## 1. Struktur Code

Project ini adalah REST API sederhana dengan Go, Gin, dan Gorm, dengan fitur utamanya adalah:

1. Registrasi user.
2. Login user dan pembuatan JWT.
3. Endpoint membaca data user (dilindungi JWT).
4. Setiap fitur dibuat dalam pola per-module.

Struktur foldernya:

- `cmd/main.go` untuk bootstrap aplikasi.
- `config/` untuk koneksi database.
- `database/entities/` untuk model yang dipetakan ke tabel database.
- `middlewares/` untuk middleware (authentication JWT).
- `modules/auth/` untuk login dan JWT.
- `modules/user/` untuk fitur user.
- `pkg/helpers/` untuk helper umum seperti password hashing.
- `pkg/utils/` untuk helper response HTTP.

### Alur request secara umum

Saat request datang, alurnya seperti ini:

```
HTTP Request
    │
    v
[Middleware]     -> Cek JWT token (khusus route protected)
    │
    v
[Controller]     -> Baca request, panggil Validation, kirim response
    │
    v
[Validation]     -> Validasi dan parse input (body JSON / path param)
    │
    v
[Service]        -> Business logic (hash password, mapping data, dll)
    │
    v
[Repository]     -> Query ke database via Gorm
    │
    v
[Database]       -> PostgreSQL
```

### `cmd/main.go`

File ini adalah titik awal aplikasi. Fungsinya, yaitu:

1. Load file `.env`.
2. Koneksi ke database.
3. Auto migrate tabel `User`.
4. Membuat instance Gin.
5. Menyusun repository, service, dan controller.
6. Mendaftarkan route auth dan user.
7. Menjalankan server.

Jadi `main.go` ibarat 'otak' atau ruang kontrol yang merangkai semua komponen sebelum server berjalan.

### Entity user

Entity user ada di `database/entities/user_entity.go`.

```go
type User struct {
  Common
  Name     string `gorm:"not null" json:"name"`
  Email    string `gorm:"unique;not null" json:"email"`
  Password string `gorm:"not null" json:"-"`
  Role     string `gorm:"default:'user'" json:"role"`
}
```

Penjelasan variable dan constraintnya:

- `Common` membawa field umum seperti `ID`, `CreatedAt`, `UpdatedAt`, dan `DeletedAt` (mirip kyk composite).
- `Name` wajib diisi.
- `Email` wajib diisi dan harus unik.
- `Password` disimpan di database, tetapi disembunyikan dari JSON response (`json:"-"`).
- `Role` punya nilai default `user`.

### DTO (Data Transfer Object)

DTO ada di `modules/user/dto/user_dto.go`.

Ada dua fungsi penting, yaitu:

- `CreateUserRequest` untuk input saat membuat user.
- `UserResponse` untuk menentukan data yang akan dikirim ke client.

### Validasi input

Validasi ada di `modules/user/validation/user_validation.go`.

```go
func ValidateCreateUser(c *gin.Context) (*dto.CreateUserRequest, error) {
  var req dto.CreateUserRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    return nil, err
  }
  return &req, nil
}

func ValidateGetUserByID(c *gin.Context) (uint, error) {
  id, err := strconv.ParseUint(c.Param("id"), 10, 64)
  if err != nil {
    return 0, err
  }
  return uint(id), nil
}
```

- `ValidateCreateUser` membaca JSON dari request body dan mengisinya ke struct `CreateUserRequest`. Kalau field wajib kosong atau format email salah, langsung mengembalikan error.
- `ValidateGetUserByID` mem-parse path parameter `:id` dari URL menjadi angka `uint`. Dengan cara ini, controller tidak perlu berurusan langsung dengan string parsing.

### Middleware authentication

Middleware ada di `middlewares/authentication.go`.

```go
func Authentication(jwtService *service.JWTService) gin.HandlerFunc {
  return func(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
      utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized: Token tidak valid")
      c.Abort()
      return
    }

    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    claims, err := jwtService.ValidateToken(tokenString)
    if err != nil {
      utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
      c.Abort()
      return
    }

    c.Set("user_id", claims.UserID)
    c.Set("email", claims.Email)
    c.Set("role", claims.Role)
    c.Next()
  }
}
```

Middleware dijalankan **sebelum** controller. Tugasnya membaca header `Authorization`, memvalidasi JWT, dan kalau valid, menyimpan data user ke context Gin agar bisa dipakai di handler selanjutnya. Kalau token tidak valid atau tidak ada, request langsung dihentikan dengan `c.Abort()`.

### Helper response

Format response standar ada di `pkg/utils/response.go`.

- `ErrorResponse` dipakai untuk response error.
- `SuccessResponse` dipakai untuk response sukses.

Pola response yang konsisten memudahkan client membaca hasil API.

## 2. Cara Code Bekerja

### 2.1 Wiring saat aplikasi start

Di `cmd/main.go`, urutannya seperti ini:

1. Load environment dari `.env`.
2. Koneksi database lewat `config.SetupDatabase()`.
3. Auto migrate entity `User`.
4. Buat Gin engine.
5. Buat dependency berantai:
  - JWT service,
  - user repository,
  - user service,
  - auth service,
  - user controller,
  - auth controller.
6. Register route auth dan user di group `/api`.
7. Jalankan server dengan `r.Run()`.

### 2.2 Alur dari HTTP Request ke Database

Pola layer yang dipakai adalah:

```
HTTP Request -> Middleware -> Controller -> Validation -> Service -> Repository -> Database
```

#### Contoh alur: `POST /api/users` (create user, tanpa login)

1. Route di `modules/user/routes.go` menerima request `POST /api/users`.
2. Tidak ada middleware (endpoint bersifat public).
3. Controller `CreateUser` di `modules/user/controller/user_controller.go` memanggil validasi.
4. Validasi `ValidateCreateUser` di `modules/user/validation/user_validation.go` melakukan bind JSON ke `CreateUserRequest`.
5. Kalau valid, controller memanggil service `CreateUser`.
6. Service melakukan hash password, membentuk object `entities.User`, lalu memanggil repository `Create(user)`.
7. Repository menjalankan query insert ke database melalui Gorm.
8. Service memetakan hasil ke `UserResponse` dan mengembalikannya ke controller.
9. Controller mengirim response `201 Created`.

#### Contoh alur: `GET /api/users/:id` (protected)

1. Route di `modules/user/routes.go` menerima request `GET /api/users/:id`.
2. Middleware `Authentication(jwtSvc)` dijalankan lebih dulu untuk mengecek header `Authorization`.
3. Kalau token tidak valid, request dihentikan dengan `401 Unauthorized`.
4. Kalau token valid, request diteruskan ke controller `GetUserByID`.
5. Controller memanggil `ValidateGetUserByID(c)` untuk mem-parse `:id` dari URL.
6. Kalau ID bukan angka, response `400 Bad Request` dikirim.
7. Controller memanggil service `GetUserByID(id)`.
8. Service memanggil repository `FindByID(id)`.
9. Repository mencari user di database.
10. Jika tidak ada, service mengembalikan error, controller merespon `404 Not Found`.
11. Jika ada, service memetakan ke `UserResponse` dan controller mengirim `200 OK`.

#### Contoh alur: `POST /api/auth/login`

1. Route di `modules/auth/routes.go` menerima request `POST /api/auth/login`.
2. Tidak ada middleware (endpoint bersifat public).
3. Controller auth memvalidasi payload login.
4. Service mencari user berdasarkan email lewat user repository.
5. Password dari request dicek dengan hash di database memakai `helpers.CheckPasswordHash`.
6. Jika cocok, JWT dibuat dengan `JWTService` dan dikembalikan ke client.

### 2.3 Alur dari Database ke HTTP Response

Arah balik datanya seperti ini:

```
Database -> Repository -> Service -> Controller -> HTTP Response
```

Penjelasan per layer:

1. Repository mengembalikan hasil query atau error dari Gorm.
2. Service menerapkan aturan bisnis (mapping data, filter field sensitif).
3. Controller menerjemahkan hasil ke status code HTTP.
4. Controller mengirim JSON response standar lewat helper di `pkg/utils/response.go`.

Contoh mapping hasil ke response:

- Validasi input gagal -> `400 Bad Request`.
- Token JWT tidak ada atau tidak valid -> `401 Unauthorized`.
- User tidak ditemukan -> `404 Not Found`.
- Error database -> `500 Internal Server Error`.
- Create user berhasil -> `201 Created`.
- Login berhasil -> `200 OK` + token.

### 2.4 Peran tiap layer

| Layer | Tanggung Jawab |
|---|---|
| Middleware | Cek JWT sebelum request masuk ke controller |
| Route | Menentukan endpoint, method HTTP, dan middleware yang dipakai |
| Controller | Menerima request, panggil validation dan service, tentukan status code, kirim response |
| Validation | Mem-parse dan memvalidasi input (body JSON / path parameter) |
| Service | Business logic |
| Repository | Satu-satunya layer yang boleh menyentuh database |
| Entity | Definisi tabel database |
| DTO | Format data yang aman untuk dikirim/diterima dari client |

## 3. Challenge A -- `GET /users/:id`

### Tugasnya

Challenge A meminta kita untuk mengambil satu user berdasarkan ID dan mengembalikan `404` kalau user tidak ditemukan.

### Cara mengerjakannya

Supaya endpoint ini jalan, ada beberapa bagian yang harus ditambah:

1. Route baru untuk `GET /users/:id` dengan middleware authentication.
2. Fungsi validasi untuk mem-parse `:id`.
3. Method repository untuk cari user berdasarkan ID.
4. Method service dan controller untuk mengubah hasil query menjadi response yang aman.

### Kode route

```go
func RegisterUserRoutes(r *gin.RouterGroup, ctrl *controller.UserController, jwtSvc *authService.JWTService) {
  users := r.Group("/users")
  {
    users.POST("", ctrl.CreateUser)
    users.GET("", middlewares.Authentication(jwtSvc), ctrl.GetUsers)
    users.GET("/:id", middlewares.Authentication(jwtSvc), ctrl.GetUserByID)
  }
}
```

### Penjelasan route

- `r.Group("/users")` membuat prefix `/users` untuk semua route user.
- `users.GET("/:id", middlewares.Authentication(jwtSvc), ctrl.GetUserByID)`, tanda `:id` adalah path parameter yang nilainya dibaca dari URL. Middleware `Authentication(jwtSvc)` diletakkan sebelum controller sehingga request hanya bisa masuk kalau JWT valid.
- `POST /api/users` tidak diberi middleware karena endpoint registrasi memang harus bisa diakses tanpa login.

### Kode validasi

```go
func ValidateGetUserByID(c *gin.Context) (uint, error) {
  id, err := strconv.ParseUint(c.Param("id"), 10, 64)
  if err != nil {
    return 0, err
  }
  return uint(id), nil
}
```

### Penjelasan validasi

- `c.Param("id")` membaca string ID dari URL.
- `strconv.ParseUint` mengubah string menjadi angka positif.
- Kalau `:id` bukan angka valid, fungsi mengembalikan error, controller hanya tinggal cek `if err != nil` sehingga parsing input dilakukan di layer validation, bukan di controller.

### Kode repository

```go
func (r *UserRepository) FindByID(id uint) (*entities.User, error) {
  var user entities.User
  err := r.db.First(&user, id).Error
  return &user, err
}
```

### Penjelasan repository

- `r.db.First(&user, id)` mencari record pertama dengan primary key `id`.
- Kalau user ada, data akan masuk ke struct `user`.
- Kalau tidak ada, Gorm akan mengembalikan error `record not found`.

### Kode service

```go
func (s *UserService) GetUserByID(id uint) (*dto.UserResponse, error) {
  user, err := s.repo.FindByID(id)
  if err != nil {
    return nil, err
  }

  return mapUserResponse(user), nil
}

func mapUserResponse(user *entities.User) *dto.UserResponse {
  return &dto.UserResponse{
    ID:    user.ID,
    Name:  user.Name,
    Email: user.Email,
    Role:  user.Role,
  }
}
```

### Penjelasan service

- Service memanggil repository untuk mengambil user.
- Kalau repository error, error itu diteruskan ke controller.
- Kalau sukses, entity user diubah ke `UserResponse` lewat `mapUserResponse`.
- Fungsi `mapUserResponse` memastikan field `Password` tidak ikut dikirim ke client.

### Kode controller

```go
func (ctrl *UserController) GetUserByID(c *gin.Context) {
  id, err := validation.ValidateGetUserByID(c)
  if err != nil {
    utils.ErrorResponse(c, http.StatusBadRequest, "ID user tidak valid")
    return
  }

  user, err := ctrl.service.GetUserByID(id)
  if err != nil {
    utils.ErrorResponse(c, http.StatusNotFound, "User tidak ditemukan")
    return
  }

  utils.SuccessResponse(c, http.StatusOK, "User berhasil ditemukan", user)
}
```

### Penjelasan controller

1. Controller memanggil `ValidateGetUserByID(c)` (untuk parsing `:id` terjadi di layer validation).
2. Kalau ID tidak valid (bukan angka), controller membalas `400 Bad Request`.
3. Controller memanggil service `GetUserByID(id)`.
4. Kalau user tidak ditemukan, controller membalas `404 Not Found`.
5. Kalau sukses, data user dikirim sebagai JSON response `200 OK`.

### Cara kerja endpointnya

Contoh request (perlu menyertakan JWT di header):

```http
GET /api/users/2
Authorization: Bearer <token>
```

Alurnya:

1. Middleware mengecek JWT. Kalau tidak valid, berhenti di sini dengan `401`.
2. Controller memanggil validasi untuk parse ID `2`.
3. Service memanggil repository.
4. Repository mencari user dengan ID 2.
5. Jika ada, response `200 OK` berisi data user.
6. Jika tidak ada, response `404 Not Found`.

## 4. Challenge B -- `GET /users`

### Tugasnya

Challenge B meminta kita untuk mengambil semua user dan mengembalikannya dalam bentuk array JSON.

### Cara mengerjakannya

Sama seperti Challenge A, kita butuh beberapa bagian:

1. Route `GET /users` dengan middleware authentication.
2. Method repository untuk mengambil semua user.
3. Method service dan controller untuk mengirim array response yang aman.

### Kode route

```go
func RegisterUserRoutes(r *gin.RouterGroup, ctrl *controller.UserController, jwtSvc *authService.JWTService) {
  users := r.Group("/users")
  {
    users.POST("", ctrl.CreateUser)
    users.GET("", middlewares.Authentication(jwtSvc), ctrl.GetUsers)
    users.GET("/:id", middlewares.Authentication(jwtSvc), ctrl.GetUserByID)
  }
}
```

### Penjelasan route

- `users.GET("", middlewares.Authentication(jwtSvc), ctrl.GetUsers)`, request ke `/api/users` akan masuk middleware dulu, baru diteruskan ke handler `GetUsers`.

### Kode repository

```go
func (r *UserRepository) FindAll() ([]entities.User, error) {
  var users []entities.User
  err := r.db.Find(&users).Error
  return users, err
}
```

### Penjelasan repository

- `r.db.Find(&users)` mengambil semua record dari tabel user.
- Hasilnya disimpan ke slice `users`.
- Kalau query sukses, repository mengembalikan semua data user.

### Kode service

```go
func (s *UserService) GetAllUsers() ([]dto.UserResponse, error) {
  users, err := s.repo.FindAll()
  if err != nil {
    return nil, err
  }

  responses := make([]dto.UserResponse, 0, len(users))
  for i := range users {
    responses = append(responses, *mapUserResponse(&users[i]))
  }

  return responses, nil
}
```

### Penjelasan service

- Service memanggil repository untuk mengambil semua user.
- Setiap entity diubah menjadi `UserResponse` lewat `mapUserResponse`.
- Hasil akhir berupa slice response yang aman untuk client (tanpa password).

### Kode controller

```go
func (ctrl *UserController) GetUsers(c *gin.Context) {
  users, err := ctrl.service.GetAllUsers()
  if err != nil {
    utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar user")
    return
  }

  utils.SuccessResponse(c, http.StatusOK, "Daftar user berhasil diambil", users)
}
```

### Penjelasan controller

1. Controller memanggil service untuk mengambil semua user.
2. Kalau ada error database, response `500 Internal Server Error` dikirim.
3. Kalau sukses, response `200 OK` berisi array JSON.

### Cara kerja endpointnya

Contoh request (perlu menyertakan JWT di header):

```http
GET /api/users
Authorization: Bearer <token>
```

Alurnya:

1. Middleware mengecek JWT. Kalau tidak valid, berhenti di sini dengan `401`.
2. Controller memanggil service.
3. Service mengambil seluruh user dari repository.
4. Setiap entity dikonversi menjadi `UserResponse`.
5. Client menerima array JSON dengan response `200 OK`.

## 5. Ringkasan Perubahan Yang Dibuat

Untuk memenuhi Challenge A dan B, perubahan yang ditambahkan adalah:

1. Route `GET /users` -- dilindungi `middlewares.Authentication(jwtSvc)`.
2. Route `GET /users/:id` -- dilindungi `middlewares.Authentication(jwtSvc)`.
3. Validasi `ValidateGetUserByID(c)` -- parsing `:id` dari URL dilakukan di layer validation.
4. Repository `FindByID(id)` -- cari satu user berdasarkan primary key.
5. Repository `FindAll()` -- ambil semua user.
6. Service `GetUserByID(id)` -- return `*dto.UserResponse`.
7. Service `GetAllUsers()` -- return `[]dto.UserResponse`.
8. Controller `GetUserByID` -- memakai `ValidateGetUserByID`, response `400`/`404`/`200`.
9. Controller `GetUsers` -- response array JSON `200`/`500`.
10. Helper `mapUserResponse` -- memetakan entity ke DTO agar password tidak bocor.

