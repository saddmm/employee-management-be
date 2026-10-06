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
Tabel database akan di-migrasi otomatis dan akun demo akan di-seed.

---

## 👤 Akun Bawaan (Default Seed)

| Role   | Email                | Password   |
|--------|----------------------|------------|
| Admin  | `admin@example.com`  | `admin123` |
| Viewer | `viewer@example.com` | `viewer123`|

---

## 📖 Ringkasan API Endpoints

### Public / Health
- `GET /health` — Cek status server

### Auth
- `POST /api/auth/login` — Login & dapatkan token JWT
- `GET /api/auth/me` — Info akun yang sedang login (Protected)

### Departments
- `GET /api/departments` — List departemen (Protected)
- `GET /api/departments/:id` — Detail departemen (Protected)
- `POST /api/departments` — Tambah departemen (Admin only)
- `PUT /api/departments/:id` — Edit departemen (Admin only)
- `DELETE /api/departments/:id` — Hapus departemen (Admin only)

### Employees
- `GET /api/employees` — List karyawan dengan filter, sort & pagination (Protected)
  - Query Params: `?search=`, `?department_id=`, `?status=`, `?sort_by=`, `?sort_order=`, `?page=`, `?limit=`
- `GET /api/employees/:id` — Detail karyawan (Protected)
- `POST /api/employees` — Tambah karyawan (Admin only)
- `PUT /api/employees/:id` — Edit karyawan (Admin only)
- `DELETE /api/employees/:id` — Hapus karyawan (Admin only)
- `GET /api/employees/export/csv` — Download CSV list karyawan (Protected)

### Audit Logs
- `GET /api/audit-logs` — Riwayat aktivitas sistem (Admin only)
