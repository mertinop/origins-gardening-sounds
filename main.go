package main

import (
	"flag"
	"log"

	g "github.com/thauanvargas/goearth"
	"github.com/thauanvargas/goearth/shockwave/in"
	"github.com/thauanvargas/goearth/shockwave/room"
)

var p = flag.String("p", "", "Port")

var ext = g.NewExt(g.ExtInfo{
	Title:       "origins-gardening-sounds",
	Description: "Plays a sound when plants can be watered or harvested",
	Author:      "mertinop",
	Version:     "1.0",
})

func main() {
	ext.Initialized(onInitialized)
	ext.Connected(onConnected)
	ext.Disconnected(onDisconnected)

	roomMgr := room.NewManager(ext)

	ext.Intercept(in.PLANTDATAUPDATE).With(func(e *g.Intercept) {
		target := e.Packet.ReadString()
		value := e.Packet.ReadString()
		canBeWatered := e.Packet.ReadBool()
		waterSecondsLeft := e.Packet.ReadInt()
		remainingHarvests := e.Packet.ReadInt()
		animation := e.Packet.ReadInt()
		log.Printf("target=%s value=%s canBeWatered=%v waterSecondsLeft=%d remainingHarvests=%d animation=%d", target, value, canBeWatered, waterSecondsLeft, remainingHarvests, animation)
		canBeHarvested := value == "6" && remainingHarvests > 0
		canCompost := remainingHarvests == 0 && roomMgr.IsOwner()
		if canBeHarvested || canCompost || canBeWatered {
			go playSound("stopped_fishing.wav")
		}

	})

	ext.Run()

}

func onInitialized(e g.InitArgs) {
	log.Println("Extension initialized")
}

func onConnected(e g.ConnectArgs) {
	log.Printf("Game connected (%s)\n", e.Host)
}

func onDisconnected() {
	log.Println("Game disconnected")
}
