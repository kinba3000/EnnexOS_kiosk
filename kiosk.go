package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/mxschmitt/playwright-go"
)

func main() {
	// 1. Installiert die benötigten Browser im Hintergrund, falls sie fehlen
	log.Println("Prüfe/Installiere Playwright-Browser...")
	err := playwright.Install(&playwright.RunOptions{
		Browsers: []string{"chromium"}, // oder "firefox", "webkit"
	})
	if err != nil {
		log.Fatalf("Fehler beim Installieren der Browser: %v", err)
	}

	// 1. .env Datei laden
	errr := godotenv.Load()
	if errr != nil {
		log.Println("Hinweis: Keine .env Datei gefunden, verwende Systemeinstellungen.")
	}

	user := os.Getenv("SUNNY_USER")
	pass := os.Getenv("SUNNY_PASS")
	targetURL := os.Getenv("TARGET_URL")
	baseURL := "https://ennexos.sunnyportal.com"

	// 2. Playwright starten
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("Fehler beim Starten von Playwright: %v", err)
	}
	defer pw.Stop()

	// 3. Chromium im Kiosk-Modus starten
	fmt.Println("Starte Chromium im Kiosk-Modus...")
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
		Args:     []string{"--kiosk"},
	})
	if err != nil {
		log.Fatalf("Fehler beim Starten von Chromium: %v", err)
	}
	defer browser.Close()

	// 4. Kontext und Seite erstellen (Auflösung festlegen)
	// Default viewport
	width := 1920
	height := 1080
	if os.Getenv("KIOSK_WIDTH") != "" && os.Getenv("KIOSK_HEIGHT") != "" {
		if w, err := strconv.Atoi(os.Getenv("KIOSK_WIDTH")); err == nil {
			width = w
		}
		if h, err := strconv.Atoi(os.Getenv("KIOSK_HEIGHT")); err == nil {
			height = h
		}
	}
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		Viewport: &playwright.Size{Width: width, Height: height},
	})
	if err != nil {
		log.Fatalf("Fehler beim Erstellen des Browser-Kontexts: %v", err)
	}
	page, err := context.NewPage()
	if err != nil {
		log.Fatalf("Fehler beim Erstellen der Seite: %v", err)
	}

	// 5. Hauptseite aufrufen
	fmt.Println("Öffne Ennexos Sunny Portal...")
	_, err = page.Goto(baseURL)
	if err != nil {
		log.Fatalf("Fehler beim Aufrufen der Basis-URL: %v", err)
	}
	page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{State: playwright.LoadStateNetworkidle})

	// 6. Ersten Login-Button suchen und klicken
	loginStartButton := "a:has-text('Login'), button:has-text('Login'), a:has-text('Anmelden'), button:has-text('Anmelden')"
	buttonLocator := page.Locator(loginStartButton)

	visible, _ := buttonLocator.IsVisible()
	if visible {
		fmt.Println("Klicke auf den initialen Login-Button...")
		err = buttonLocator.Click()
		if err != nil {
			log.Printf("Warnung beim Klicken des Login-Buttons: %v", err)
		}
		// Kurz warten, bis die Formularfelder geladen sind
		time.Sleep(2 * time.Second)
	}

	// 7. Login-Formular ausfüllen
	usernameInput := "input[type='email'], input[name='username'], #txtUsername"
	passwordInput := "input[type='password'], #txtPassword"

	inputLocator := page.Locator(usernameInput)
	fieldsVisible, _ := inputLocator.IsVisible()

	if fieldsVisible {
		fmt.Println("Loginfelder gefunden. Setze Credentials ein...")
		page.Fill(usernameInput, user)
		page.Fill(passwordInput, pass)

		// Absenden (Submit-Button klicken)
		err = page.Click("button[type='submit'], #btnLogin")
		if err != nil {
			log.Fatalf("Fehler beim Klicken des Submit-Buttons: %v", err)
		}

		page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{State: playwright.LoadStateNetworkidle})
		fmt.Println("Login erfolgreich übermittelt!")
	} else {
		fmt.Println("Bereits eingeloggt oder Login-Felder nicht gefunden. Navigiere direkt weiter...")
	}

	// 8. Zur Zielseite im Kiosk navigieren
	fmt.Printf("Navigiere zur Zielseite: %s\n", targetURL)
	_, err = page.Goto(targetURL)
	if err != nil {
		log.Fatalf("Fehler beim Navigieren zur Zielseite: %v", err)
	}

	// 9. Accept cookies if the button is present
	cookieButton := "a:has-text('Accept all'), a:has-text('Alle akzeptieren')"
	cookieLocator := page.Locator(cookieButton)
	cookieVisible, _ := cookieLocator.IsVisible()

	if cookieVisible {
		fmt.Println("Cookie-Banner gefunden. Akzeptiere Cookies...")
		err = cookieLocator.Click()
		if err != nil {
			log.Printf("Warnung beim Klicken des Cookie-Buttons: %v", err)
		}
	}

	fmt.Println("Kiosk-Modus aktiv. Drücke STRG+C im Terminal zum Beenden.")

	// 9. Signal-Handler, damit das Programm offen bleibt, bis es beendet wird
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	fmt.Println("\nSkript wird beendet...")
}
