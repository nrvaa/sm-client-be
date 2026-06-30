# Dokumentasi Teknis — SM Client Portal

> Stanley Marthin Restoration · Portal Progres Restorasi Mobil Klasik

---

## Daftar Isi

1. [Gambaran Umum Proyek](#1-gambaran-umum-proyek)
2. [Arsitektur Sistem](#2-arsitektur-sistem)
3. [Struktur Proyek](#3-struktur-proyek)
4. [Sistem Autentikasi](#4-sistem-autentikasi)
   - 4.1 [URL Unik per Client](#41-url-unik-per-client)
   - 4.2 [JWT (JSON Web Token)](#42-jwt-json-web-token)
   - 4.3 [Redis — Session & Blacklist](#43-redis--session--blacklist)
   - 4.4 [Cookies](#44-cookies)
5. [Alur Lengkap Sistem](#5-alur-lengkap-sistem)
6. [Implementasi Backend (Golang)](#6-implementasi-backend-golang)
   - 6.1 [Struktur Database](#61-struktur-database)
   - 6.2 [Generate URL Token](#62-generate-url-token)
   - 6.3 [Endpoint Login](#63-endpoint-login)
   - 6.4 [Middleware JWT + Redis](#64-middleware-jwt--redis)
   - 6.5 [Endpoint Vehicles](#65-endpoint-vehicles)
   - 6.6 [Endpoint Logout](#66-endpoint-logout)
7. [Implementasi Frontend (Svelte)](#7-implementasi-frontend-svelte)
   - 7.1 [Komponen & Routing](#71-komponen--routing)
   - 7.2 [Portal.svelte — Halaman Login](#72-portalsvelte--halaman-login)
   - 7.3 [App.svelte — State & Fetch](#73-appsvelte--state--fetch)
   - 7.4 [Dashboard.svelte](#74-dashboardsvelte)
   - 7.5 [ManagerPanel.svelte](#75-managerpanelsvelte)
8. [Keamanan — Skenario Serangan](#8-keamanan--skenario-serangan)
9. [Konfigurasi & Environment Variables](#9-konfigurasi--environment-variables)
10. [Panduan Setup & Menjalankan Proyek](#10-panduan-setup--menjalankan-proyek)

---

## 1. Gambaran Umum Proyek

SM Client Portal adalah aplikasi web **eksklusif** yang memungkinkan client Stanley Marthin Restoration memantau perkembangan restorasi mobil klasik mereka secara real-time. Setiap client mendapat akses pribadi ke dashboard yang menampilkan:

- Status & persentase progres restorasi
- Log timeline pengerjaan (jurnal restorasi)
- Galeri foto & video progres (bento grid)
- List part yang diganti

Aplikasi terdiri dari dua bagian:

| Bagian            | Teknologi                      | Fungsi                              |
| ----------------- | ------------------------------ | ----------------------------------- |
| **Frontend**      | Svelte 5 + Vite + Tailwind CSS | UI client & manager                 |
| **Backend**       | Golang                         | API, autentikasi, data              |
| **Cache/Session** | Redis                          | Session management, blacklist token |
| **Database**      | PostgreSQL / SQLite            | Data client & kendaraan             |

---

## 2. Arsitektur Sistem

```
┌─────────────────────────────────────────────────────────┐
│                      CLIENT BROWSER                      │
│                                                          │
│  stanleymarthin.com/portal/a3f9x2kp                     │
│  ┌─────────────┐  ┌──────────────┐  ┌────────────────┐ │
│  │ Portal.svelte│  │ App.svelte   │  │Dashboard.svelte│ │
│  │  (login)    │  │ (state mgmt) │  │  (UI client)   │ │
│  └──────┬──────┘  └──────┬───────┘  └────────────────┘ │
└─────────┼────────────────┼────────────────────────────-─┘
          │ fetch()        │ fetch() + Bearer Token
          ▼                ▼
┌─────────────────────────────────────────────────────────┐
│                    BACKEND (Golang)                       │
│                                                          │
│  POST /api/login          GET /api/vehicles/:slug        │
│  ┌──────────────────────────────────────────────────┐   │
│  │  1. Validasi url_token + access_code             │   │
│  │  2. Generate JWT                                 │   │
│  │  3. Simpan session di Redis                      │   │
│  │  4. Set Cookie HttpOnly                          │   │
│  └──────────────────────────────────────────────────┘   │
└───────────────┬─────────────────────┬───────────────────┘
                │                     │
                ▼                     ▼
┌──────────────────────┐   ┌─────────────────────────────┐
│       REDIS           │   │        DATABASE              │
│                      │   │                              │
│  session:{slug}      │   │  clients                     │
│  blacklist:{token}   │   │  vehicles                    │
│  TTL: 7 hari         │   │  timeline_events             │
└──────────────────────┘   │  progress_photos             │
                           └─────────────────────────────┘
```

---

## 3. Struktur Proyek

```
sm-client-vite/                  # Frontend Svelte
├── src/
│   ├── App.svelte               # Root component, state management
│   ├── main.ts                  # Entry point
│   ├── app.css                  # Global styles
│   └── lib/
│       ├── Portal.svelte        # Halaman login
│       ├── Dashboard.svelte     # Dashboard utama client
│       ├── VehicleSelection.svelte  # Pilih kendaraan (multi-kendaraan)
│       ├── ManagerPanel.svelte  # Panel admin/manager
│       ├── TransitionOverlay.svelte # Animasi transisi antar halaman
│       └── mockData.ts          # Type definitions & data dummy
└── package.json

sm-client-backend/               # Backend Golang
├── main.go                      # Entry point server
├── auth/
│   └── jwt.go                  # Login, logout handler
├── handlers/
│   └── handlers.go               # Login, logout handler
├── middleware/
│   └── auth.go                   # JWT verification middleware
├── models/
│   ├── client.go                # Client model
│   └── vehicle.go               # Vehicle model
├── data/
│   └── mock.go                  # Mock data (sementara)
└── payload.json                 # Contoh payload login
```

---

## 4. Sistem Autentikasi

Proyek ini menggunakan sistem autentikasi berlapis yang terdiri dari empat komponen yang saling melengkapi.

### 4.1 URL Unik per Client

Setiap client mendapat URL acak yang di-generate **sekali saja** saat admin membuat akun client baru. URL ini tidak berubah selama akun masih aktif.

```
stanleymarthin.com/portal/a3f9x2kp   ← Budi
stanleymarthin.com/portal/m7k2np4q   ← Dedi
stanleymarthin.com/portal/x9b3wt6r   ← Anto
```

**Mengapa acak?** Kalau URL-nya predictable seperti `/portal/sm-budi`, orang yang tahu nama client lain bisa menebak URL-nya. Dengan token acak 8 karakter (hex), ada 4 miliar kemungkinan kombinasi — praktis tidak bisa ditebak.

**Analogi:** URL = username, access code = password. Keduanya harus cocok untuk bisa masuk.

**Generate token di Golang:**

```go
import (
    "crypto/rand"
    "encoding/hex"
)

func generateURLToken() string {
    bytes := make([]byte, 4) // 4 bytes = 8 karakter hex
    rand.Read(bytes)
    return hex.EncodeToString(bytes) // contoh: "a3f9x2kp"
}
```

Token ini disimpan permanen di kolom `url_token` tabel `clients` di database.

---

### 4.2 JWT (JSON Web Token)

JWT adalah token digital yang di-generate backend setelah login berhasil. Token ini berisi informasi client yang sudah di-sign dengan secret key, sehingga tidak bisa dipalsukan.

**Struktur JWT payload:**

```json
{
  "slug": "sm-budi",
  "exp": 1234567890,
  "iat": 1234481490
}
```

**Di backend Golang:**

```go
import "github.com/golang-jwt/jwt/v5"

func generateJWT(slug string) (string, error) {
    claims := jwt.MapClaims{
        "slug": slug,
        "exp":  time.Now().Add(7 * 24 * time.Hour).Unix(),
        "iat":  time.Now().Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
```

**Di frontend (App.svelte)**, JWT dikirim di setiap request sebagai Bearer token:

```js
const res = await fetch(`http://localhost:3000/api/vehicles/${slug}`, {
  headers: {
    Authorization: `Bearer ${token}`,
  },
});
```

**Kenapa tidak cukup hanya JWT?** JWT tidak bisa "dicabut" sebelum expired. Kalau user logout, token lama tetap valid sampai waktu expiry-nya habis — kecuali kita punya mekanisme tambahan. Di sinilah Redis berperan.

---

### 4.3 Redis — Session & Blacklist

Redis adalah in-memory database yang sangat cepat. Di proyek ini digunakan untuk dua hal:

#### Session Storage

Saat login berhasil, backend simpan sesi aktif di Redis:

```
Key:   session:sm-budi
Value: { "login_at": "2024-01-15T10:30:00Z", "ip": "192.168.1.1" }
TTL:   7 hari (604800 detik)
```

Setiap request yang masuk, backend **wajib cek Redis** apakah sesi masih ada. Kalau key sudah tidak ada (TTL habis atau sudah dihapus), request ditolak walaupun JWT-nya masih valid.

```go
// Cek sesi di Redis
val, err := redisClient.Get(ctx, "session:"+slug).Result()
if err == redis.Nil {
    return c.JSON(401, "Sesi tidak ditemukan, silakan login ulang")
}
```

#### Token Blacklist

Saat logout, backend masukkan token ke blacklist Redis:

```
Key:   blacklist:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Value: "1"
TTL:   Sama dengan sisa expiry JWT
```

Setiap request, backend cek blacklist sebelum proses apapun:

```go
// Cek blacklist
exists, _ := redisClient.Exists(ctx, "blacklist:"+tokenString).Result()
if exists > 0 {
    return c.JSON(401, "Token sudah tidak berlaku")
}
```

**Ringkasan peran Redis:**

| Key Pattern         | Kapan dibuat        | Kapan dihapus           | TTL             |
| ------------------- | ------------------- | ----------------------- | --------------- |
| `session:{slug}`    | Saat login berhasil | Saat logout / TTL habis | 7 hari          |
| `blacklist:{token}` | Saat logout         | Otomatis saat TTL habis | Sisa expiry JWT |

---

### 4.4 Cookies

Cookies dengan flag `HttpOnly` adalah cara paling aman menyimpan JWT di browser karena JavaScript tidak bisa mengaksesnya — melindungi dari serangan XSS (Cross-Site Scripting).

**Backend set cookie saat login:**

```go
c.SetCookie(&http.Cookie{
    Name:     "sm_jwt",
    Value:    tokenString,
    MaxAge:   7 * 24 * 60 * 60, // 7 hari dalam detik
    Path:     "/",
    HttpOnly: true,              // tidak bisa diakses JS
    Secure:   true,              // hanya HTTPS
    SameSite: http.SameSiteStrictMode,
})
```

**Perbandingan localStorage vs Cookie HttpOnly:**

|                      | localStorage (sekarang) | Cookie HttpOnly (rekomendasi) |
| -------------------- | ----------------------- | ----------------------------- |
| Bisa diakses JS      | ✅ Ya                   | ❌ Tidak (lebih aman)         |
| Rentan XSS           | ✅ Ya                   | ❌ Tidak                      |
| Auto kirim ke server | ❌ Manual               | ✅ Otomatis                   |
| Bisa set expiry      | ❌ Tidak                | ✅ Ya (`Max-Age`)             |
| Mudah diimplementasi | ✅ Lebih mudah          | Perlu setup CORS              |

**Sinkronisasi durasi — ketiganya harus sama:**

```
Cookie Max-Age  = 7 hari  ← browser pegang token selama ini
JWT exp         = 7 hari  ← token valid selama ini
Redis TTL       = 7 hari  ← sesi aktif selama ini
```

Yang paling pendek yang menang. Kalau Redis TTL 1 hari tapi JWT 7 hari, user tetap logout setelah 1 hari.

---

## 5. Alur Lengkap Sistem

### Login

```
1. Client buka link: stanleymarthin.com/portal/a3f9x2kp
2. Frontend baca url_token dari URL path: "a3f9x2kp"
3. Client input access code: "SM-BUDI"
4. Frontend kirim POST /api/login:
   {
     "url_token": "a3f9x2kp",
     "access_code": "SM-BUDI"
   }
5. Backend lookup DB: cari client dengan url_token = "a3f9x2kp"
6. Backend validasi: client.access_code == "SM-BUDI"? ✅
7. Backend generate JWT dengan payload { slug: "sm-budi", exp: ... }
8. Backend simpan session di Redis: SET session:sm-budi ... EX 604800
9. Backend kirim response: { token: "eyJ...", slug: "sm-budi" }
   + Set-Cookie: sm_jwt=eyJ...; HttpOnly; Secure; Max-Age=604800
10. Frontend simpan token di localStorage
11. Frontend redirect ke dashboard dengan slug
```

### Fetch Data Kendaraan

```
1. Frontend fetch: GET /api/vehicles/sm-budi
   Header: Authorization: Bearer eyJ...
2. Backend extract token dari header
3. Backend cek blacklist Redis: ada? → 401 Unauthorized
4. Backend decode & verify JWT signature
5. Backend ambil slug dari JWT: "sm-budi"
6. Backend cek: slug di JWT == slug di URL path? ✅
7. Backend cek session Redis: session:sm-budi ada? ✅
8. Backend query DB: ambil vehicles milik slug "sm-budi"
9. Backend return data kendaraan
10. Frontend render dashboard
```

### Logout

```
1. User klik tombol Keluar
2. Frontend kirim POST /api/logout
   Header: Authorization: Bearer eyJ...
3. Backend masukkan token ke Redis blacklist:
   SET blacklist:eyJ... 1 EX {sisa_expiry}
4. Backend hapus session Redis: DEL session:sm-budi
5. Backend clear cookie: Set-Cookie: sm_jwt=; Max-Age=0
6. Frontend hapus localStorage: removeItem("sm_token")
7. Frontend redirect ke halaman login
```

---

## 6. Implementasi Backend (Golang)

### 6.1 Struktur Database

```sql
-- Tabel clients
CREATE TABLE clients (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL,
    access_code VARCHAR(20)  NOT NULL UNIQUE,  -- "SM-BUDI"
    url_token   VARCHAR(16)  NOT NULL UNIQUE,  -- "a3f9x2kp"
    slug        VARCHAR(50)  NOT NULL UNIQUE,  -- "sm-budi"
    created_at  TIMESTAMP DEFAULT NOW()
);

-- Tabel vehicles
CREATE TABLE vehicles (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_slug           VARCHAR(50) REFERENCES clients(slug),
    brand                 VARCHAR(50),
    model                 VARCHAR(50),
    year                  INT,
    license_plate         VARCHAR(20),
    vin                   VARCHAR(50),
    status                TEXT,
    completion_percentage INT DEFAULT 0,
    estimated_completion  DATE,
    banner_image          TEXT,
    restoration_type      VARCHAR(50),
    created_at            TIMESTAMP DEFAULT NOW()
);

-- Tabel timeline_events
CREATE TABLE timeline_events (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vehicle_id  UUID REFERENCES vehicles(id),
    status      VARCHAR(100),
    date        DATE,
    description TEXT,
    image       TEXT,
    completed   BOOLEAN DEFAULT FALSE
);

-- Tabel progress_photos
CREATE TABLE progress_photos (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vehicle_id UUID REFERENCES vehicles(id),
    url        TEXT,
    title      VARCHAR(100),
    date       DATE,
    stage      VARCHAR(100),
    type       VARCHAR(10) DEFAULT 'photo'
);
```

---

### 6.2 Generate URL Token

```go
// models/client.go
package models

import (
    "crypto/rand"
    "encoding/hex"
    "strings"
)

type Client struct {
    ID         string `db:"id"`
    Name       string `db:"name"`
    AccessCode string `db:"access_code"`
    URLToken   string `db:"url_token"`
    Slug       string `db:"slug"`
}

// Dipanggil SEKALI saat admin buat client baru
func GenerateURLToken() string {
    bytes := make([]byte, 4)
    rand.Read(bytes)
    return hex.EncodeToString(bytes) // "a3f9x2kp"
}

func AccessCodeToSlug(accessCode string) string {
    // "SM-BUDI" → "sm-budi"
    return strings.ToLower(accessCode)
}
```

---

### 6.3 Endpoint Login

```go
// handlers/handlers.go
package handlers

import (
	"fmt"
	"sm-client-backend/auth"
	"sm-client-backend/data"

	"github.com/gofiber/fiber/v3"
)

type LoginRequest struct {
	Username string `json:"username"`
}

// Login handles user authentication and JWT generation
func Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		fmt.Println("Bind Error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	fmt.Printf("Received Login Request - Username: '%s'", req.Username)

	// Verify against mock users
	for _, user := range data.MockUsers {
		if user.Username == req.Username {
			// Valid credentials, generate token
			token, err := auth.GenerateToken(user.ID, user.Slug)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to generate token",
				})
			}
			return c.JSON(fiber.Map{
				"token": token,
				"slug":  user.Slug,
			})
		}
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Invalid credentials",
	})
}
```

---

### 6.4 Middleware JWT

```go
// middleware/auth.go
package middleware

import (
	"strings"

	"sm-client-backend/auth"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

// Protected verifies the JWT and validates the layer 2 security (matching slug)
func Protected() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing Authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Authorization header format",
			})
		}

		tokenString := parts[1]

		// Layer 1: Verify Token
		token, err := jwt.ParseWithClaims(tokenString, &auth.JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
			return auth.SecretKey, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		claims, ok := token.Claims.(*auth.JWTClaim)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		// Layer 2: Match Slug
		requestSlug := c.Params("slug")
		if requestSlug != "" && requestSlug != claims.Slug {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied",
			})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("slug", claims.Slug)

		return c.Next()
	}
}
```

---

### 6.5 Endpoint Vehicles

```go
// handlers/handlers.go (lanjutan)
func GetVehicle(c fiber.Ctx) error {
	slug := c.Params("slug")

	var results []models.Vehicle
	for _, v := range data.InitialVehicles {
		if v.ClientCode == slug {
			results = append(results, v)
		}
	}

	if len(results) > 0 {
		return c.JSON(results)
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "Vehicle not found",
	})
}
```

---

## 7. Implementasi Frontend (Svelte)

### 7.1 Komponen & Routing

Proyek menggunakan Svelte 5 dengan reactive state (`$state`, `$derived`, `$effect`). Routing dilakukan secara manual di `App.svelte` berdasarkan state `loggedInClientCode`.

```
App.svelte (root)
│
├── [!loggedInClientCode]    → Portal.svelte (halaman login)
├── [loggedInClientCode && !selectedVehicleId]
│                            → VehicleSelection.svelte (pilih mobil)
└── [loggedInClientCode && selectedVehicleId]
                             → Dashboard.svelte (dashboard utama)
                             └── ManagerPanel.svelte (panel admin, overlay)
```

---

### 7.2 Portal.svelte — Halaman Login

Halaman login dengan glassmorphism design. Fitur utama:

- Input access code (dengan toggle show/hide)
- Rate limiting di frontend: maksimal 5 percobaan, lockout 30 detik
- Memanggil `POST /api/login` ke backend Golang

**Yang perlu diupdate** — tambahkan `url_token` dari URL path ke request body:

```js
// Ambil url_token dari URL path saat ini
// contoh: stanleymarthin.com/portal/a3f9x2kp → "a3f9x2kp"
const urlToken = window.location.pathname.split("/").pop();

const response = await fetch("http://localhost:3000/api/login", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    url_token: urlToken, // ← tambahkan ini
    access_code: cleanCode, // "SM-BUDI"
  }),
});
```

**Rate limiting (sudah ada di kode):**

```js
let loginAttempts = $state(0);
let lockoutUntil = $state(0);

// Cek lockout sebelum proses login
if (Date.now() < lockoutUntil) {
  const remaining = Math.ceil((lockoutUntil - Date.now()) / 1000);
  errorMessage = `Terlalu banyak percobaan. Coba lagi dalam ${remaining} detik.`;
  return;
}

// Setelah 5 gagal → lockout 30 detik
if (loginAttempts >= 5) {
  lockoutUntil = Date.now() + 30000;
}
```

---

### 7.3 App.svelte — State & Fetch

Root component yang mengelola seluruh state aplikasi.

**State utama:**

```js
let vehicles = $state<Vehicle[]>([]);
let loggedInClientCode = $state<string | null>(null);
let selectedVehicleId = $state<string>("");
let showManager = $state(false);
```

**handleLogin** — dipanggil setelah login berhasil di Portal.svelte:

```js
async function handleLogin(code: string) {
    const token = localStorage.getItem("sm_token");

    if (token) {
        const res = await fetch(`http://localhost:3000/api/vehicles/${code}`, {
            headers: { "Authorization": `Bearer ${token}` }
        });

        if (res.ok) {
            vehicles = await res.json();
        } else {
            handleLogout(); // Token invalid → paksa logout
            return;
        }
    }

    loggedInClientCode = code;
}
```

**handleLogout** — bersihkan semua state:

```js
function handleLogout() {
  localStorage.removeItem("sm_token");
  vehicles = [];
  loggedInClientCode = null;
  selectedVehicleId = "";
  showManager = false;
  // Idealnya juga panggil POST /api/logout ke backend
}
```

---

### 7.4 Dashboard.svelte

Halaman utama yang dilihat client. Terdiri dari tiga chapter:

**Chapter I — Progres Saat Ini**

- Status unit (teks deskriptif)
- Progress bar persentase
- Estimasi selesai

**Chapter II — Jurnal Restorasi**

- List part yang diganti (tabel: Part, Dari, Menjadi)
- Timeline log progres (urutan terbaru di atas)
- Setiap entry timeline bisa punya foto

**Chapter III — Galeri Visual**

- Bento grid foto & video
- Lightbox saat diklik
- Tombol download media

**Data types yang digunakan (`mockData.ts`):**

```typescript
interface Vehicle {
  id: string;
  clientCode: string;
  clientName: string;
  brand: string;
  model: string;
  year: number;
  licensePlate: string;
  vin: string;
  status: string;
  completionPercentage: number;
  estimatedCompletion: string;
  bannerImage: string;
  restorationType: string;
  timeline: TimelineEvent[];
  gallery: ProgressPhoto[];
  replacedParts?: ReplacedPart[];
}
```

---

### 7.5 ManagerPanel.svelte

Panel overlay untuk admin/mekanik mengupdate progres. Fitur:

- Pilih kendaraan yang ingin diupdate
- Edit status, persentase, estimasi selesai, tipe restorasi
- Tambah entry timeline baru (dengan opsional foto)
- Upload foto (mock URL atau file upload)

Panel ini seharusnya hanya bisa diakses admin — perlu ditambahkan proteksi agar client biasa tidak bisa membuka panel ini.

---

## 8. Keamanan — Skenario Serangan

| Skenario                               | Mekanisme Perlindungan                                                              | Hasil                |
| -------------------------------------- | ----------------------------------------------------------------------------------- | -------------------- |
| Nebak URL orang lain                   | URL token 8 karakter acak (4 miliar kombinasi)                                      | ❌ Praktis mustahil  |
| Tau URL tapi tidak tau access code     | Validasi `url_token + access_code` harus cocok di backend                           | ❌ Ditolak           |
| Tau access code tapi tidak tau URL     | URL acak tidak bisa ditebak dari access code                                        | ❌ Ditolak           |
| Buka URL Budi dengan access code Dedi  | Backend: `url_token` milik Budi, `access_code` Dedi → tidak cocok                   | ❌ Ditolak           |
| Gunakan JWT Budi untuk akses data Dedi | Backend cek: slug di JWT vs slug di URL path harus sama                             | ❌ Ditolak           |
| Brute force access code                | Rate limiting 5 percobaan + lockout 30 detik di frontend, tambahkan juga di backend | ❌ Dilambatkan       |
| Curi JWT via XSS                       | Cookie HttpOnly: JavaScript tidak bisa akses cookie                                 | ❌ Tidak bisa dicuri |
| Gunakan JWT lama setelah logout        | Redis blacklist: token di-blacklist saat logout                                     | ❌ Ditolak           |
| Session expired tapi JWT masih valid   | Redis TTL: sesi dihapus Redis setelah 7 hari                                        | ❌ Ditolak           |

---

## 9. Konfigurasi & Environment Variables

**Backend `.env`:**

```env
# Server
PORT=3000

# Database
DATABASE_URL=postgres://user:password@localhost:5432/sm_client

# Redis
REDIS_URL=redis://localhost:6379

# JWT
JWT_SECRET=your-super-secret-key-minimal-32-karakter

# CORS
ALLOWED_ORIGIN=https://stanleymarthin.com

# Environment
APP_ENV=production
```

**Frontend `.env`:**

```env
VITE_API_URL=http://localhost:3000
```

---

## 10. Panduan Setup & Menjalankan Proyek

### Prerequisites

```bash
# Backend
go 1.21+
redis-server
postgresql (atau sqlite untuk development)

# Frontend
node 18+
npm atau pnpm
```

### Frontend

```bash
cd sm-client-vite
npm install
npm run dev        # development: http://localhost:5173
npm run build      # production build
```

### Backend

```bash
cd sm-client-backend

# Install dependencies
go mod tidy

# Jalankan Redis (pastikan sudah terinstall)
redis-server

# Copy dan isi environment variables
cp .env.example .env

# Jalankan backend
go run main.go     # development: http://localhost:3000
```

### Generate URL Token untuk Client Baru

Saat admin menambah client baru, panggil endpoint admin (perlu dibuat):

```bash
curl -X POST http://localhost:3000/admin/clients \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {ADMIN_TOKEN}" \
  -d '{
    "name": "Budi Santoso",
    "access_code": "SM-BUDI"
  }'

# Response:
# {
#   "slug": "sm-budi",
#   "url_token": "a3f9x2kp",
#   "portal_url": "https://stanleymarthin.com/portal/a3f9x2kp"
# }
```

URL yang dikembalikan inilah yang dikirim ke client via WhatsApp/email bersama access code-nya.

---

_Dokumentasi ini dibuat berdasarkan source code sm-client-vite dan sm-client-backend per Juni 2026._
