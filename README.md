# 🏢 Employee Management System - Backend API

RESTful API untuk Employee Management System menggunakan **Go (Golang)**, **Fiber v2**, **GORM**, dan **MySQL**.

---

## 🚀 Fitur Utama

- **Authentication & Authorization**: JWT token auth dengan 2 roles (`admin` dan `viewer`).
- **Department Management**: CRUD data departemen (hanya `admin` yang bisa create/update/delete).
- **Employee Management**: CRUD karyawan dengan relasi ke departemen, validasi server-side lengkap.
- **Search, Filter & Sort**: Sorting & filtering dieksekusi langsung pada query database MySQL.
- **Audit Logging**: Mencatat aktor, entity, tindakan (create/update/delete), dan payload snapshot sebelum & sesudah.
- **CSV Export**: Generate file `.csv` data karyawan langsung dari server.
- **Rate Limiting**: Melindungi API dari brute force/spam request.
- **Health Check**: Endpoint `/health` untuk status monitoring.
- **Swagger Documentation**: Dokumentasi interaktif OpenAPI / Swagger UI di `/swagger/`.

---

## 📋 Prasyarat

- Go 1.25+
- Docker & Docker Compose (untuk database MySQL lokal)

---

## 🛠️ Cara Menjalankan

### 1. Salin Environment Variables
```bash
cp .env.example .env
```

### 2. Jalankan Database MySQL (Docker Compose)
```bash
docker compose up -d
```
Container `ems_mysql` akan berjalan di port `3306`.

### 3. Jalankan Aplikasi Backend
```bash
go run cmd/server/main.go
```
Server akan berjalan di `http://localhost:8080`.
Dokumentasi Swagger dapat diakses di `http://localhost:8080/swagger/`.

---

## 👤 Akun Bawaan (Default Seed)

| Role   | Email                | Password   |
|--------|----------------------|------------|
| Admin  | `admin@example.com`  | `admin123` |
| Viewer | `viewer@example.com` | `viewer123`|
