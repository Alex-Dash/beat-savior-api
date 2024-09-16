module bsvapi

go 1.21.5

require fyne.io/systray v1.11.0

require (
	github.com/gen2brain/beeep v0.0.0-20240516210008-9c006672e7f4 // indirect
	github.com/go-toast/toast v0.0.0-20190211030409-01e6764cf0a4 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.23 // indirect
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d // indirect
	github.com/tadvi/systray v0.0.0-20190226123456-11a2b8fa57af // indirect
	golang.design/x/clipboard v0.7.0 // indirect
	golang.org/x/exp v0.0.0-20190731235908-ec7cb31e5a56 // indirect
	golang.org/x/image v0.6.0 // indirect
	golang.org/x/mobile v0.0.0-20230301163155-e0f57694e12c // indirect
	golang.org/x/sys v0.25.0 // indirect
)

require bsvapi/bsdfilereader v0.0.0

replace bsvapi/bsdfilereader => ./modules/bsdfilereader

require bsvapi/types v0.0.0

replace bsvapi/types => ./modules/types

require bsvapi/bsddb v0.0.0

replace bsvapi/bsddb => ./modules/bsddb
