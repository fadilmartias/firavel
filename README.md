# Firavel 🌀

**Firavel** adalah starter kit backend untuk Go, menggunakan [Fiber v2](https://gofiber.io/), dengan struktur dan filosofi yang mengadopsi gaya Laravel — tapi tetap ringan, modular, dan siap untuk produksi.

## 🔧 Fitur Utama

* **Struktur folder ala Laravel**
  Susunan kode familiar seperti `app/`, `routes/`, `database/`, `app/controllers/`, `app/models/`, `app/middleware/`, dan lainnya.

* **Fiber v2**
  Framework HTTP super cepat, ringan, dan mudah digunakan.

* **Middleware support**
  Dukungan middleware customizable dengan modular approach.

* **Routing terorganisir**
  Route handler ditata dalam file terpisah sesuai domain/fungsi.

* **Hot-reload saat development**
  Dilengkapi *auto-reload* lewat Air untuk iterasi cepat.

* **Environment config \*\*\*\*`.env`**
  Konfigurasi sederhana melalui `.env`, ideal untuk berbagai environment.

* **Validator built-in**
  Validasi request dengan pendekatan langsung dan efisien.

* **Modular controller & model**
  Clean architecture memudahkan pengembangan dan testing.

* **Database via GORM**
  ORM powerful + dukungan migration & seeder.

* **Redis cache/session**
  Integrasi Redis siap pakai.

* **Log rotasi harian**
  Logging otomatis dengan rotasi harian untuk file log.

* **Seeder & migration tools**
  Bangun, reset, dan isi database dengan mudah.

---

## 📦 Instalasi Cepat

1. Clone repo:

   ```bash
   git clone https://github.com/fadilmartias/firavel.git
   cd firavel
   ```
2. Setup environment:

   ```bash
   cp .env.example .env
   ```
3. Jalankan untuk development:

   ```bash
   go mod tidy
   air
   ```
4. Atau build untuk production:

   ```bash
   go build -o firavel cmd/firavel/main.go
   ./firavel
   ```

## ❤️ Contribute & Support

Terima kasih sudah menggunakan Firavel!
Kalau kamu menemukan bug, fitur yang bisa ditambahkan, atau ingin berkontribusi — silakan buka issue atau pull request di repo GitHub.
