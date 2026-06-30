
# EnnexOS Kiosk for Displaying Dashboards

This application launches a headless Chromium web browser, automatically logs the user in, and displays the designated EnnexOS dashboard.

---



## 🚀 Setup & Installation

> For compiling the .go file, you need the go compiler installed. For more informations, reefer to[go.dev](https://go.dev/)

### 1. Configuration
1. Copy the `.env.example` file into the same directory as your executable.
2. Rename it to `.env`.
3. Open the `.env` file and fill in your credentials.
```SUNNY_USER="example@example.com"
SUNNY_PASS="yourpassword"
#Replace XXXXXXXX with your actual dashboard ID in the TARGET_URL below.
TARGET_URL="https://ennexos.sunnyportal.com/xxxxxxxx/dashboard" 
``` 

### 2. Building the Executable
Run the appropriate command in your terminal to build the executable from the Go source file:

#### 🖥️ Windows
``` go build -o sunny-kiosk.exe ```

#### Linux 🍓(Raspberry Pi)
``` 
$env:GOOS="linux"
$env_GOARCH="arm64"
go build -o sunny-kiosk
```



## Autostart setup
#### 🖥️ Windows
For autostart setup, move your `sunny-kiosk.exe` and your `.env` in the autostart folder of windows.
Tip: You can open this folder quickly by pressing Win + R, typing shell:startup, and hitting Enter.

#### Linux 🍓(Raspberry Pi)
1. Move your `sunny-kiosk` and `.env` in your home directory. 
2. Change rights of the executable: `chmod +x  sunny-kiosk`
3. Create a Cronjob.
``` crontab -e
# Add the following line to run on boot:
@reboot /home/your-username/sunny-kiosk
```


### Usage
Simply double click the `sunny-kiosk.exe` or run it via shell on linux (`./sunny-kiosk`)


### Problems

If you encounter any problems or have questions, please open a new issue on the [EnnexOS_Kioks Issues page](https://github.com/kinba3000/EnnexOS_kiosk/issues). 


### Contribution

Contributions are highly welcome!

1. Fork the repository.
2. Create your feature branch.
3. Open a Pull Request with a detailed description of your changes.