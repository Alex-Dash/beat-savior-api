package main

import (
	"fmt"
	"log"
	"os"

	_ "embed"

	"bsvapi/bsddb"
	"bsvapi/bsdfilereader"
	"bsvapi/bsdweb"
	"bsvapi/types"

	"fyne.io/systray"
	"github.com/gen2brain/beeep"
	"golang.design/x/clipboard"
)

var (
	bsd_reader *bsdfilereader.BSD_Reader = &bsdfilereader.BSD_Reader{}
	db         *bsddb.DBwrap             = &bsddb.DBwrap{}
)

var (
	//go:embed assets/ico/red_a.ico
	ICON_RED_A []byte
	//go:embed assets/ico/red_d.ico
	ICON_RED_D []byte
	//go:embed assets/ico/blue_a.ico
	ICON_BLUE_A []byte
	//go:embed assets/ico/blue_d.ico
	ICON_BLUE_D []byte
	//go:embed assets/ico/green_a.ico
	ICON_GREEN_A []byte
	//go:embed assets/ico/green_d.ico
	ICON_GREEN_D []byte
	//go:embed assets/ico/orange_a.ico
	ICON_ORANGE_A []byte
	//go:embed assets/ico/orange_d.ico
	ICON_ORANGE_D []byte
)

func startSetup() error {

	sess_ch := make(chan *types.BSD_Session)
	song_ch := make(chan *types.BSD_Song)
	web_ch := types.WEB_Settings{
		OnNewSession:  &sess_ch,
		OnNewSongData: &song_ch,
	}

	systray.SetTooltip("Creating database")
	err := db.Init()
	if err != nil {
		return err
	}

	systray.SetTooltip("Indexing song data")
	err = bsd_reader.Init(db, &web_ch)
	if err != nil {
		return err
	}

	if bsd_reader.Headers != nil {
		log.Printf("Header count: %d\n", len(*bsd_reader.Headers))
	}

	if bsd_reader.Songs != nil {
		log.Printf("Song data count: %d\n", len(*bsd_reader.Songs))
	}

	err = bsdweb.Init(&web_ch)
	if err != nil {
		return err
	}

	return nil
}

func onURLcopy(c chan struct{}) {
	<-c

	err := clipboard.Init()
	if err != nil {
		log.Println(err)
	} else {
		clipboard.Write(clipboard.FmtText, []byte("http://127.0.0.1:1337"))
	}
	log.Println("Copied URL")

	go onURLcopy(c)
}

func onQuit(c chan struct{}) {
	<-c
	systray.Quit()
}

func onToggleLan(c chan struct{}, mi *systray.MenuItem) {
	<-c
	if mi.Checked() {
		systray.SetIcon(ICON_GREEN_A)
		mi.SetTitle("Open to LAN")
		mi.Uncheck()
	} else {
		systray.SetIcon(ICON_BLUE_A)
		mi.SetTitle("Switch to LOCAL")
		mi.Check()
	}
	go onToggleLan(c, mi)
}

func onReady() {
	systray.SetIcon(ICON_ORANGE_D)
	systray.SetTitle("Beat Savior API")
	systray.SetTooltip("Starting Beat Savior API")
	mURLcopy := systray.AddMenuItem("Copy API URL", "Copy API web address into clipboard")
	go onURLcopy(mURLcopy.ClickedCh)

	mRebuild := systray.AddMenuItem("Rebuild Database", "Completely rebuild the database using new data")
	mRebuild.SetIcon(ICON_RED_A)

	mLan := systray.AddMenuItem("Open to LAN", "Open server to external connections")
	mLan.SetIcon(ICON_BLUE_A)
	go onToggleLan(mLan.ClickedCh, mLan)

	mQuit := systray.AddMenuItem("Quit", "Close API server")
	mQuit.SetIcon(ICON_RED_D)
	go onQuit(mQuit.ClickedCh)

	err := startSetup()
	if err != nil {
		log.Println(err)
		systray.SetIcon(ICON_RED_D)
		systray.SetTooltip(fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	err = beeep.Notify("Beat Savior API", "Running in the background\nClick the tray icon for more settings", "assets/png/green_a.png")
	if err != nil {
		panic(err)
	}
	systray.SetTooltip("Beat Savior API is running")
	systray.SetIcon(ICON_GREEN_A)
}

func onExit() {
	db.Close()
	os.Exit(0)
}

func main() {
	systray.Run(onReady, onExit)
}
