package main

import (
	"bufio"
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"golang.org/x/net/proxy"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

func main() {
	fmt.Println("Tor Scraper Başlatılıyor...")

	os.MkdirAll("logs", 0755)
	os.MkdirAll("screenshots", 0755)
	os.MkdirAll("scraper/html", 0755)

	logFile, _ := os.Create("logs/scan_report.log")
	defer logFile.Close()
	logger := bufio.NewWriter(logFile)
	
	writeLog := func(msg string) { 
		zaman := time.Now().Format("2006-01-02 15:04:05")
		mesaj := fmt.Sprintf("[%s] %s", zaman, msg)
		fmt.Println(mesaj) 
		logger.WriteString(msg + "\n")
		logger.Flush()
	}

	writeLog("--- TARAMA BAŞLANGICI: " + time.Now().Format(time.RFC850) + " ---")
	//  Tor Proxy Ayarları (SOCKS5)
	torProxy := "127.0.0.1:9150" 
	writeLog("INFO:Tor Proxy bağlantısı kontrol ediliyor... ")
	//  Proxy Bağlantısını Oluşturma
	dialer, err := proxy.SOCKS5("tcp", torProxy, nil, proxy.Direct)
	if err != nil {
		writeLog("HATA: Tor bağlantısı hatası: " + err.Error())
		return
	}

	tr := &http.Transport{Dial: dialer.Dial}
	client := &http.Client{Transport: tr,Timeout:   time.Second * 90,}

	writeLog("Hedef dosyası okunuyor...")

	//  Dosyayı Aç
	file, err := os.Open("targets.yaml") 
	if err != nil {
		writeLog("Dosya açılamadı:" + err.Error())
		return
	}
	defer file.Close()

	u := launcher.New().
		Proxy("socks5://127.0.0.1:9150").
		Headless(true).
		Leakless(false).
		MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	//  Dosyayı Satır Satır Oku
	scanner := bufio.NewScanner(file)	
	sayac := 0
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		url = strings.TrimPrefix(url, "-") 
		url = strings.TrimSpace(url)

		if url == "" || !strings.HasPrefix(url, "http") {
			continue
		}
		sayac++
		writeLog(fmt.Sprintf("[%d] Hedef Bulundu: %s\n", sayac, url))
		resp, err := client.Get(url)
		if err != nil {
			writeLog(fmt.Sprintf("[HATA] %s adresine gidilemedi: %v\n", url, err))
		} else {
			writeLog(fmt.Sprintf("[BAŞARILI] %s -> Durum Kodu: %d\n", url, resp.StatusCode))

			dosyaAdi := fmt.Sprintf("site_%d_%d", sayac, time.Now().Unix())

			htmlData, _ := ioutil.ReadAll(resp.Body)
			ioutil.WriteFile("scraper/html/"+dosyaAdi+".html", htmlData, 0644)
			writeLog(" -> HTML Kaydedildi.")
			resp.Body.Close()

			page, err := browser.Page(proto.TargetCreateTarget{URL: url})
			if err != nil {
				writeLog(" -> Sayfa açılırken hata: " + err.Error())
				continue
			}

			// Çözünürlüğü ayarla (Geniş ekran)
			page.MustSetViewport(1920, 1080, 1.0, false)

			// Yüklenmesini bekle (Timeout koyuyoruz ki sonsuza kadar beklemesin)
			err = page.Timeout(90 * time.Second).WaitLoad()
			if err != nil {
				writeLog(" -> Sayfa yükleme zaman aşımı (Screenshot yine de denenecek).")
			}
			
			// Tor yavaş olduğu için ekstra bekleme (Senin isteğin üzerine)
			time.Sleep(5 * time.Second)
			

			img, err := page.Screenshot(true, &proto.PageCaptureScreenshot{
                Format:  proto.PageCaptureScreenshotFormatPng, 
            })

			if err == nil {
				ioutil.WriteFile("screenshots/"+dosyaAdi+".png", img, 0644)
				writeLog(" -> Ekran Görüntüsü Kaydedildi.")
			} else {
				writeLog(" -> Screenshot alınamadı: " + err.Error())
			}
		}
	}

	fmt.Println("Tüm işlemler tamamlandı.")
}