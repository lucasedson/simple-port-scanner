build:
	go build -ldflags "-H=windowsgui -s -w" .\src\main.go

build-android:
	cd .\src\ && fyne package -os android -appID com.example.myapp -icon ..\mobile-icon.png