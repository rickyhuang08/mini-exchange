# Mini Exchange

## A. Cara Menjalankan Project
Prasyarat
- Go versi 1.24 atau lebih baru
- OpenSSL (untuk generate RSA key pair)
- Git (untuk cloning repository)
#
### Langkah-langkah
1. Clone Repository
    - `git clone <repository-url>`
    - `cd mini-exchange`
2. Konfigurasi Environment
    - Salin file `.env.development.example` menjadi `.env.development`:
        - `cp .env.development.example .env.development`
    - Sesuaikan isi file `.env.development` jika diperlukan (misal: port, API key Finnhub, dll).
3. Konfigurasi YAML
    - Salin file `config/config.development.yaml.example` menjadi `config/config.development.yaml`:
        - `cp config/config.development.yaml.example config/config.development.yaml`
    - Edit file tersebut sesuai kebutuhan (misal: pengaturan logger, database, dll). Catatan: Project ini menggunakan inmemory storage, jadi tidak perlu konfigurasi database.
4. Generate RSA Key Pair (untuk JWT)
    - `cd config/keys`
    - `openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048`
    - `openssl rsa -pubout -in private.pem -out public.pem`
    - `cd ../..`
5. Install Dependencies
    - Jalankan script install_deps.sh untuk mengunduh semua dependensi Go:
        - `./install_deps.sh`
    - Script ini akan menjalankan go mod download dan `go mod tidy`.
6. Jalankan Aplikasi
    - `go run cmd/server/main.go`
    - Server akan berjalan di http://localhost:8080 (default).
    - Untuk mengubah port, atur variabel PORT di file `config/config.development.yaml`
### Kredensial JWT (untuk testing)
1. Gunakan salah satu akun berikut untuk mendapatkan token JWT melalui endpoint login (jika diimplementasikan):
    - Email: ybtech@example.com, Password: ybtech1234
    - Email: alice@example.com, Password: alice1234
    - Email: bob@example.com, Password: bob1234
    - Email: charlie@example.com, Password: charlie1234
    - Catatan: Fitur autentikasi JWT bersifat opsional (bonus). Jika diaktifkan, semua endpoint REST (kecuali WebSocket) memerlukan header Authorization: Bearer (token).
#
## B. Design Arsitektur (Singkat)
1. Aplikasi ini mengadopsi clean architecture dengan pemisahan layer sebagai berikut:
    - Domain Layer (internal/domain) Berisi entity bisnis (Order, Trade, OrderBook) dan value object.
    - Entity Layer (internal/domain) Berisi entity delivery to https (Order, Trade, OrderBook) dan value object.
    - Repository Layer (internal/repository) Implementasi penyimpanan inmemory dengan concurrency safe map (sync.RWMutex). Menyimpan data order, trade, dan order book.
    - Usecase Layer (internal/usecase) Berisi logic bisnis: matching engine, order management, market data, dan price simulator.
    - Delivery Layer
        - REST Handler (internal/delivery/http): Menyediakan endpoint HTTP (Gin).
        - WebSocket Hub & Client (internal/delivery/websocket): Menangani koneksi WebSocket realtime, subscription channel, dan broadcast.
    - Adapter Layer (internal/adapter/external) Untuk integrasi dengan sumber data eksternal (misal: Finnhub WebSocket) – opsional.
    - Pkg Layer (pkg/) Berisi komponen reusable yang bersifat modular dan bisa digunakan lintas layer maupun lintas project. Layer ini biasanya berisi abstraksi atau wrapper yang cukup “besar” dan memiliki tanggung jawab teknis tertentu.
        - Contoh isi:
            - JWT Service
            - Logger
    - Helpers Layer (helpers/) Berisi global constant dan generic function sederhana yang tidak termasuk utility kompleks atau service modular. Layer ini lebih ringan dibanding pkg/ dan biasanya berisi fungsi stateless atau nilai konstan yang sering dipakai di banyak tempat.
        - Contoh isi:
            - Global constant (status enum, role enum, dll)
            - Function sederhana seperti:
                - String formatter
                - Mapping kecil
                - Helper konversi tipe data
    - Middleware Layer (middleware/) Berisi komponen middleware yang digunakan pada layer delivery (HTTP / WebSocket) untuk menangani cross cutting concerns sebelum request diproses oleh handler.
        - Contoh tanggung jawab middleware:
            - JWT Authentication & Authorization
            - Logging request & response
            - CORS handling
    - Config Layer (config/) Berisi seluruh konfigurasi aplikasi yang dibutuhkan saat startup. Layer ini bertanggung jawab untuk membaca, memparsing, dan menyediakan konfigurasi dalam bentuk struct yang dapat diinject ke komponen lain.
        - Konfigurasi biasanya bersumber dari:
            - File .env
            - File YAML / JSON
            - Environment variable sistem
            - File key (misal RSA private/public key)

2. Komunikasi antar layer menggunakan interface untuk menjaga ketergantungan terbalik (dependency inversion). Semua komponen berjalan di atas goroutine dan channel untuk mencapai konkurensi tinggi.
#
## C. Flow System
1. Order Flow (REST)
    - Client mengirim POST /orders dengan JSON {stock_code, side, price, quantity}.
    - Handler memvalidasi input, membuat objek Order, lalu memanggil OrderUsecase.PlaceOrder.
    - Usecase menyimpan order ke repository, lalu mengirim order ke matching engine melalui channel per stock (SubmitOrder).
    - Matching engine (goroutine per stock) memproses order secara serial, melakukan matching dengan order yang ada di order book.
    - Jika terjadi trade, matching engine membuat Trade, menyimpannya, dan mengirim ke channel TradeBroadcast.
    - Order book diperbarui, order yang terisi penuh dihapus dari antrian.
2. Market Data Flow (Price Simulator)
    - PriceSimulator berjalan sebagai goroutine terpisah, memperbarui harga setiap 2 detik untuk saham tertentu (AAPL, GOOG, MSFT).
    - Setiap update, harga disimpan di TickerRepository (opsional) dan di‑broadcast ke subscriber channel ticker via WebSocket hub.
3. WebSocket Flow
    - Client melakukan koneksi ke ws://localhost:8080/api/realtime/ws.
    - Server melakukan upgrade, membuat objek Client, dan mendaftarkannya ke Hub.
    - Client dapat mengirim pesan JSON untuk subscribe/unsubscribe ke channel tertentu (ticker, trade, orderbook) untuk suatu saham.
    - Hub menyimpan mapping client per channel & saham.
    - Ketika ada event (trade, ticker update, orderbook change), hub mengirim pesan ke semua client yang subscribe.
4. External Data Integration (Opsional)
    - FinnhubAdapter (atau adapter lain) terkoneksi ke WebSocket publik, menerima trade realtime.
    - Trade tersebut dikirim ke channel externalTradeChan, lalu diproses (disimpan, broadcast via hub) seperti trade internal.
#
## D. Assumption yang Digunakan
1. Penyimpanan in memory Data tidak persisten, akan hilang saat server restart.
2. Matching sederhana Hanya mendukung limit order, tidak ada market order, stop order, dll.
3. FIFO (First In First Out) Untuk order dengan harga sama, prioritas diberikan berdasarkan urutan masuk (karena order disimpan dalam slice dan tidak diurutkan ulang berdasarkan waktu).
4. Harga matching Untuk trade, harga yang digunakan adalah harga order sell (sesuai implementasi). Alternatif bisa menggunakan harga buy atau midpoint.
5. Partial fill didukung Order dapat terisi sebagian, status menjadi PARTIAL.
6. Tidak ada fee atau komisi.
7. Satu stock diproses oleh satu goroutine – Menjamin tidak ada race condition untuk stock yang sama.
8. WebSocket channel Hanya tiga channel: ticker, trade, orderbook. Format pesan subscribe: {"action":"subscribe","channel":"ticker","stock":"AAPL"}.
9. Simulasi harga internal Digunakan jika tidak ada sumber eksternal; harga berubah secara random walk.
#
## E. Penjelasan Potensi Race Condition dan Pencegahannya
1. Potensi Race Condition
    - Update order book secara concurrent – Jika dua goroutine mengakses order book yang sama untuk stock yang sama.
    - Update status order Saat order diubah (misal dikurangi quantity) oleh lebih dari satu proses.
    - Akses ke repository Map di repository diakses oleh banyak goroutine (misal dari stock processor yang berbeda).
2. Strategi Pencegahan
    - Per‑stock goroutine
    - Setiap stock memiliki satu goroutine (processOrders) yang menangani semua order untuk stock tersebut. Semua operasi order book dan matching dilakukan dalam goroutine ini secara serial, sehingga tidak ada race condition untuk stock yang sama.
    - Mutex di Repository Setiap repository (order, trade, order book) menggunakan sync.RWMutex untuk melindungi map internal. Operasi baca/tulis aman untuk concurrent akses dari goroutine stock yang berbeda.
    - Channel untuk Komunikasi Pengiriman order ke matching engine menggunakan channel buffered, sehingga pengirim tidak perlu menunggu pemrosesan. Ini menghindari blocking dan race.
    - Atomic Operation Untuk operasi sederhana seperti increment counter, tidak diperlukan karena kita sudah serial per stock.
    - WebSocket Hub Map subscriber dilindungi dengan sync.RWMutex untuk menghindari race saat subscribe/unsubscribe dan broadcast.
3. Dengan pendekatan ini, sistem aman terhadap concurrent order submission dan koneksi WebSocket simultan.
#
## F. Strategi Broadcast Non‑Blocking
1. Broadcast dilakukan oleh WebSocket hub ke banyak client. Untuk menghindari satu client lambat menghambat client lain, digunakan strategi berikut:
    - Setiap client memiliki buffered channel Send (ukuran 256).
    - Saat broadcast, hub melakukan range atas subscriber, lalu mencoba mengirim pesan ke channel Send client dengan select:
    - ```
        go
        select {
        case client.Send <- data:
            // pesan terkirim ke buffer client
        default:
            // buffer penuh, pesan di-drop dan client dianggap lambat
            log.Println("Client send buffer full, dropping message")
        }
    - Jika buffer penuh, pesan di‑drop (tidak memblokir). Client yang terlalu lambat mungkin akan kehilangan beberapa pesan, tetapi tidak mempengaruhi client lain.
    - Goroutine WritePump setiap client bertugas membaca dari channel Send dan menulis ke koneksi WebSocket secara serial. Jika terjadi error saat menulis, client akan di‑unregister.
2. Pendekatan ini memastikan broadcast tetap lancar meskipun ada client dengan koneksi lambat atau terputus.
#
## G. Tiga Bottleneck Utama dan Cara Mengatasinya
1. Single‑stock Processor
    - Bottleneck: Untuk stock dengan volume order sangat tinggi, goroutine tunggal yang memproses semua order bisa menjadi bottleneck.
    - Solusi:
        - Sharding Mempartisi order berdasarkan kriteria lain (misal: range harga) atau menggunakan multiple goroutine per stock dengan mekanisme locking yang lebih canggih (contoh: stripe locking).
        - Optimasi kode matching Gunakan struktur data yang efisien (heap, tree) daripada slice dan sorting setiap kali.
        - Batch processing Kumpulkan beberapa order sebelum diproses.
2. Repository Locks
    - Bottleneck: Jika banyak stock berbeda mengakses repository yang sama (misal OrderRepository), mutex global dapat menyebabkan kontensi.
    - Solusi:
        - Sharded mutex Gunakan map of mutex per stock atau per key range.
        - Per‑stock repository Pisahkan penyimpanan per stock (contoh: map of maps) sehingga kunci per stock.
        - Optimasi baca Gunakan RLock untuk operasi baca agar tidak saling blokir.

3. WebSocket Broadcast
    - Bottleneck: Broadcast ke ribuan client dapat membebani CPU dan network.
    - Solusi:
        - Batch broadcast Jika banyak client subscribe ke channel yang sama, kirim pesan sekali ke goroutine worker yang kemudian mendistribusikan.
        - Gunakan connection pool Batasi jumlah koneksi per IP.
        - Kompresi pesan Gunakan websocket.Compress untuk mengurangi ukuran data.
        - Scale out Jalankan multiple instance server dengan load balancer, dan gunakan pub/sub seperti Redis untuk menyinkronkan broadcast antar instance.
#
## Catatan Tambahan
- Untuk menjalankan integrasi dengan Finnhub, pastikan Anda memiliki API key dan mengaturnya di .env.development sebagai FINNHUB_API_KEY.
- Jika ingin menonaktifkan price simulator internal, komentari pemanggilan priceSimulator.Start() di main.go.
- Dokumentasi API lengkap (endpoint REST) dan WebSocket (cara connect, subscribe, format pesan) tersedia di file terpisah (API.md dan WEBSOCKET.md).