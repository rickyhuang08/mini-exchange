# API Documentation

## Base URL
- Semua endpoint berada di http://localhost:8080/api.
- Untuk environment production, sesuaikan dengan domain yang digunakan.

## Autentikasi
- Sebagian besar endpoint bersifat protected dan memerlukan token JWT yang dikirim melalui header:
    - Authorization: Bearer (token)
- Token dapat diperoleh dengan melakukan login ke endpoint publik /public/login.
- Jika ingin menguji tanpa autentikasi (misal untuk development), Anda dapat mengganti prefix /protected/ menjadi /public/ pada URL.
- Catatan: Endpoint public hanya tersedia untuk mode development dan tidak aman untuk production.

## Daftar Endpoint
1. Login (Publik)
    - Mendapatkan token JWT untuk akses endpoint protected.
    - Endpoint:
        - `POST /public/login`
    - Request Body (JSON)
        - ```
            {
                "email": "string",
                "password": "string"
            }
    - Contoh Request
        - ```
            curl --location 'http://localhost:8080/api/public/login' \
            --header 'Content-Type: application/json' \
            --data-raw '{
                "email": "ybtech@example.com",
                "password": "ybtech1234"
            }'
    - Response Sukses (200 Ok)
        ```
        {
            "status": "success",
            "message": "Login successful",
            "data": {
                "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InlidGVjaEBleGFtcGxlLmNvbSIsImV4cCI6MTc3MjcxNDM5NCwiaWF0IjoxNzcyNjI3OTk0LCJyb2xlIjoxLCJ1c2VyX2lkIjoxfQ.RL8Iau5SvggY_DID1TNoYhYVM5GaN4HBfT8PNavHpM5CnpO6wY-CwuvOvHWLVl39zZkqUKJdvMEpWpc4bl6UCteH3sSmNA5XVJGgrQdCGixXJBM1lbx194lv9b22s-b3lWgY_AxNBmmmEAtld375dvlppdb6lGAejJscdfDuM-eUVaJjEMWEogvoO7oom16ffw4VKYg38RVEZpNwS10mGJ8NIvH9A4ItAko936iPwafCA9oz52TxMizye_W1y88ODFMhh9o63xzxEWjzvgot6t36xqQ5Nm_yQpP90ej6PxJXBtS0CvUZwtEU5acTxJvNNfvTcToy9OxGab4sE93qxw",
                "user": {
                    "id": 1,
                    "name": "YBTech",
                    "email": "ybtech@example.com",
                    "role": 1
                }
            }
        }
        ```
    - Kredensial yang tersedia:
        | Email | Password |
        |-------|----------|
        | ybtech@example.com | ybtech1234 |
        | alice@example.com | alice1234 |
        | bob@example.com | bob1234 |
        | charlie@example.com | charlie1234 |

2. Create Order (Protected)
    - Membuat order baru (BUY atau SELL).
    - Endpoint: 
        - `POST /protected/orders`
        - Request Headers:
            - `Authorization: Bearer <token>`
            - `Content-Type: application/json`
        - Request Body (JSON)
            | Field       | Tipe    | Deskripsi                              |
            |------------|---------|------------------------------------------|
            | stock_code | string  | Kode saham (contoh: "AAPL")             |
            | side       | string  | "BUY" atau "SELL"                       |
            | price      | number  | Harga per unit                          |
            | quantity   | integer | Jumlah lot (unit)                       |
        - Contoh Request
            ```
            curl --location 'http://localhost:8080/api/protected/orders' \
            --header 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...' \
            --header 'Content-Type: application/json' \
            --data '{
                "stock_code": "AAPL",
                "Side": "SELL",
                "Price": 263,
                "quantity": 10
            }'
            ```
        - Response Sukses (200 Ok)
            ```
            {
                "status": "success",
                "message": "Order placed successfully",
                "data": {
                    "order_id": "ef07c7a5-bbb3-4b90-81df-5a1a78b9622e"
                }
            }
            ```

3. Get Order List (Protected)
    - Mengambil daftar order yang telah dibuat, dengan filter opsional.
    - Endpoint:
        - `GET /protected/orders`
        - Query Parameters (opsional)
            | Parameter  | Tipe   | Deskripsi                                                           |
            |------------|--------|---------------------------------------------------------------------|
            | stock_code | string | Filter berdasarkan kode saham                                      |
            | status     | string | Filter berdasarkan status (NEW, PARTIAL, FILLED, CANCELLED)        |
        - Request Headers
            - `Authorization: Bearer <token>`
            - Contoh Request
                ```
                curl --location 'http://localhost:8080/api/protected/orders?stock_code=AAPL' \
                --header 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...'
                ```
            - Response Sukses (200 OK)
                ```
                {
                    "status": "success",
                    "message": "Order list retrieved successfully",
                    "data": [
                        {
                        "id": "520ae97b-b3ef-48e7-9136-91514656b2e5",
                        "stock_code": "AAPL",
                        "side": "BUY",
                        "price": 263,
                        "quantity": 0,
                        "status": "FILLED",
                        "created_at": "2026-03-04T19:34:04.067788912+07:00",
                        "updated_at": "2026-03-04T19:34:13.929397038+07:00"
                        }
                    ]
                }
                ```
4. Get Trade History (Protected)
    - Menampilkan riwayat trade yang telah terjadi.
    - Endpoint:
        - `GET /protected/trades`
        - Query Parameters (opsional)
            | Parameter  | Tipe   | Deskripsi                                                           |
            |------------|--------|---------------------------------------------------------------------|
            | stock_code | string | Filter berdasarkan kode saham                                      |
        - Request Headers
            - `Authorization: Bearer <token>`
        - Contoh Request
            ```
            curl --location 'http://localhost:8080/api/protected/trades?stock_code=AAPL' \
            --header 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...'
            ```
        - Response Sukses (200 OK)
            ```
            {
                "status": "success",
                "message": "Trade history retrieved successfully",
                "data": [
                    {
                    "id": "trade-uuid",
                    "stock_code": "AAPL",
                    "buy_order_id": "buy-order-uuid",
                    "sell_order_id": "sell-order-uuid",
                    "price": 263,
                    "quantity": 10,
                    "traded_at": "2026-03-04T19:34:13.929182438+07:00"
                    }
                ]
            }
            ```
5. Get Market Snapshot (Protected)
    - Mengambil snapshot pasar terkini untuk suatu saham: harga terakhir, perubahan, volume, order book (bid/ask), dan trade terbaru.
    - Endpoint: 
        - `GET /protected/market/{stock}/snapshot`
        - Path Parameter
            | Parameter  | Tipe   | Deskripsi                                                           |
            |------------|--------|---------------------------------------------------------------------|
            | stock | string | Kode saham (contoh: AAPL)                                      |
        - Request Headers
            - `Authorization: Bearer <token>`
        - Contoh Request
            ```
            curl --location 'http://localhost:8080/api/protected/market/AAPL/snapshot' \
            --header 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...'
            ```
        - Response Sukses (200 OK)
            ```
            {
                "status": "success",
                "message": "Trade history retrieved successfully",
                "data": {
                    "stock_code": "AAPL",
                    "last_price": 263.5,
                    "change": 0.5,
                    "volume": 1000,
                    "order_book": {
                        "bids": [
                        [262.5, 100],
                        [262.0, 200]
                        ],
                        "asks": [
                        [263.5, 150],
                        [264.0, 300]
                        ]
                    },
                    "recent_trades": [
                        {
                        "price": 263.5,
                        "quantity": 10,
                        "time": 1678901234
                        }
                    ]
                }
            }
            ```
6. Catatan Penting
- Semua endpoint protected memerlukan token JWT yang valid. Token dikirim melalui header `Authorization`.
- Untuk pengujian tanpa autentikasi, ganti prefix `/protected/` dengan `/public/` pada URL. Contoh: `POST /public/orders` akan sama dengan versi protected tetapi tanpa validasi token.
- Data bersifat in-memory, sehingga akan hilang jika server di-restart.
- Untuk WebSocket documentation, lihat file terpisah `WEBSOCKET.md`.