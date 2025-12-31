#  Thor's Scraper

Go (Golang) ve Rod kütüphanesi kullanılarak geliştirilmiş, Tor ağı üzerindeki .onion sitelerini tarayan, durumlarını kontrol eden ve kanıt olarak ekran görüntüsü alan siber istihbarat aracıdır.

## 🚀 Özellikler
- **Tam Gizlilik:** Tüm trafik Tor Proxy (SOCKS5) üzerinden geçer.
- **Otomasyon:** `targets.yaml` listesindeki siteleri sırayla tarar.
- **Kanıt Toplama:** Sitelerin HTML kaynak kodlarını ve PNG formatında ekran görüntülerini kaydeder.

## 🛠️ Kurulum

1. Bilgisayarınızda **Tor Browser**'ın açık olduğundan emin olun (Port: 9150).
2. Repoyu klonlayın:
   ```bash
   git clone [https://github.com/Sinemsena/Thor-un_Scraperi.git](https://github.com/Sinemsena/Thor-un_Scraperi.git)


3.Gerekli kütüphaneleri indirin:
    ```bash

      go mod tidy

4.Çalıştırın:
     ```bash

       go run main.go

📂 Çıktılar
Program çalıştığında aşağıdaki klasörleri otomatik oluşturur:

/logs: Tarama raporu ve zaman damgaları.

/output/screenshots: Sitelerin ekran görüntüleri.

/scraper/html: Sitelerin kaynak kodları.
