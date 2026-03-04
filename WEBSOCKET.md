# WebSocket Documentation
## Base URL WebSocket
- `ws://localhost:8080/api/realtime/ws`
- Untuk koneksi aman (WSS) jika menggunakan production dengan SSL, gunakan:
- `wss://your-domain.com/api/realtime/ws`

1. Cara Connect
- Koneksi WebSocket dapat dilakukan menggunakan berbagai client. Pastikan Anda telah memiliki token JWT? Pada implementasi ini, WebSocket bersifat publik (tidak memerlukan autentikasi) untuk memudahkan pengujian.
    - Contoh connect menggunakan JavaScript (browser):
        ```
        const ws = new WebSocket('ws://localhost:8080/api/realtime/ws');

        ws.onopen = () => {
            console.log('Connected to WebSocket');
        };

        ws.onmessage = (event) => {
            console.log('Received:', event.data);
        };

        ws.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        ws.onclose = () => {
            console.log('Disconnected');
        };
        ```
    - Atau bisa membuka folder test-websocket/ (terdapat html -> dengan memeriksa devtools console)
    - Contoh connect menggunakan websocat (command line):
        - `websocat ws://localhost:8080/api/realtime/ws`
    - Contoh connect menggunakan wscat:
        - `wscat -c ws://localhost:8080/api/realtime/ws`
2. Cara Subscribe
- Setelah koneksi terbuka, client harus mengirim pesan JSON untuk subscribe ke channel tertentu. Format pesannya adalah sebagai berikut:
    - Format Pesan Subscribe:
        ```
        {
            "action": "subscribe",
            "channel": "ticker",
            "stock": "AAPL"
        }
        ```
        - action: wajib, bernilai "subscribe" atau "unsubscribe".
        - channel: jenis channel yang diinginkan. Nilai yang tersedia:
            - "ticker" : update harga terakhir (perubahan harga).
            - "trade" : trade yang terjadi.
            - "orderbook" : depth order book (bids/asks) – opsional.
        - stock: kode saham (misal "AAPL", "GOOG", "MSFT"). Client dapat subscribe ke beberapa saham dengan mengirim pesan terpisah untuk setiap saham.
        - Contoh subscribe ke channel ticker untuk AAPL:
            - `{"action":"subscribe","channel":"ticker","stock":"AAPL"}`
        - Contoh subscribe ke channel trade untuk AAPL:
            - `{"action":"subscribe","channel":"trade","stock":"AAPL"}`
        - Contoh subscribe ke channel orderbook untuk AAPL:
            - `{"action":"subscribe","channel":"orderbook","stock":"AAPL"}`
    - Unsubscribe
        - Untuk berhenti menerima update dari suatu channel, kirim pesan dengan action "unsubscribe":
            - `Untuk berhenti menerima update dari suatu channel, kirim pesan dengan action "unsubscribe":`

3. Format Message (dari Server)
    ```
    {
        "channel": "ticker",
        "stock": "AAPL",
        "data": {
            // konten tergantung channel
        }
    }
    ```
    - Channel `ticker`
        - Data berisi harga terbaru dan timestamp (dalam detik Unix):
        ```
        {
            "channel": "ticker",
            "stock": "AAPL",
            "data": {
                "price": 145.67,
                "time": 1678901234
            }
        }
        ```
    - Channel trade
        - Data berisi detail trade:
        ```
        {
            "channel": "trade",
            "stock": "AAPL",
            "data": {
                "id": "550e8400-e29b-41d4-a716-446655440000",
                "price": 145.67,
                "quantity": 100,
                "traded_at": "2026-03-04T12:34:56Z"
            }
        }
        ```
    - Channel orderbook
        - Data berisi daftar bid dan ask (masing-masing array [price, quantity]):
        ```
        {
            "channel": "orderbook",
            "stock": "AAPL",
            "data": {
                "bids": [[145.0, 500], [144.5, 300]],
                "asks": [[145.5, 200], [146.0, 400]]
            }
        }
        ```
4. Rekomendasi Tools dan Cara Testing
- Berikut beberapa tools yang dapat digunakan untuk menguji WebSocket:
    1. Browser Console (JavaScript)
    - Cara paling mudah. Buka Developer Tools (F12) di browser, lalu paste kode berikut:
    ```
    const ws = new WebSocket('ws://localhost:8080/api/realtime/ws');

    ws.onopen = () => {
        console.log('Connected');
        ws.send(JSON.stringify({action: "subscribe", channel: "ticker", stock: "AAPL"}));
    };

    ws.onmessage = (e) => {
        console.log('Received:', JSON.parse(e.data));
    };
    ```

    2. websocat (command line)
    - Install: `brew install websocat` (macOS) atau download dari github.com/vi/websocat.
    - Contoh penggunaan:
        - `websocat ws://localhost:8080/api/realtime/ws`
        - Setelah koneksi, ketikkan JSON subscribe lalu Enter:
        - `{"action":"subscribe","channel":"ticker","stock":"AAPL"}`
        - Pesan masuk akan ditampilkan di terminal.
    3. wscat (Node.js)
    - Install: `npm install -g wscat`
    - Conoh:
        - `wscat -c ws://localhost:8080/api/realtime/ws`
        - Kemudian kirim pesan subscribe seperti di atas.
    4. Postman
    - Postman mendukung WebSocket. Buat request baru dengan tipe WebSocket, masukkan URL, lalu connect. Setelah connect, kirim pesan JSON di bagian "Message".
    5. Online WebSocket Tester
    - PieSocket WebSocket Tester
    - LivePerson WebSocket Tester
    - Masukkan URL ws://localhost:8080/api/realtime/ws, connect, lalu kirim pesan subscribe.
    6. Atau dengan membuka forlder `test-websocket/`
        - dimana didalamnya terdapat file `test-websocket.html` untuk melihat hasilnya dapat membuka Devtools (F12) -> Console
5. Catatan Penting
- Pastikan server sedang berjalan sebelum mencoba koneksi.
- Jika menggunakan URL localhost, pastikan tidak ada proxy yang menghalangi.
- Untuk production, gunakan wss:// jika server menggunakan SSL.
- Jika tidak menerima pesan setelah subscribe, periksa apakah ada data yang dikirim oleh server (misal, price simulator berjalan atau ada trade).
- Untuk unsubscribe, kirim pesan dengan action "unsubscribe".