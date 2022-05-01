package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/hypebeast/go-osc/osc"
)

var (
	avatarPath          = "/avatar/parameters/"
	characterParameters = []string{
		avatarPath + "OSC_DC_0",
		avatarPath + "OSC_DC_1",
		avatarPath + "OSC_DC_2",
		avatarPath + "OSC_DC_3",
		avatarPath + "OSC_DC_4",
		avatarPath + "OSC_DC_5",
		avatarPath + "OSC_DC_6",
		avatarPath + "OSC_DC_7",
		avatarPath + "OSC_DC_8",
		avatarPath + "OSC_DC_9",
		avatarPath + "OSC_DC_10",
		avatarPath + "OSC_DC_11",
	}
	characterLimit                 = len(characterParameters)
	isShowingNotificationParameter = avatarPath + "OSC_DC_IS_SHOWING"
	letterMap                      = map[string]float32{
		" ": .0,
		"a": .33,
		"b": .34,
		"c": .35,
		"d": .36,
		"e": .37,
		"f": .38,
		"g": .39,
		"h": .40,
		"i": .41,
		"j": .42,
		"k": .43,
		"l": .44,
		"m": .45,
		"n": .46,
		"o": .47,
		"p": .48,
		"q": .49,
		"r": .50,
		"s": .51,
		"t": .52,
		"u": .53,
		"v": .54,
		"w": .55,
		"x": .56,
		"y": .57,
		"z": .58,
	}
)

type PermissionGrantedMessage struct {
	WasPermissionGranted bool
}

type DiscordNotificationMessage struct {
	Username string
}

func main() {
	cmd := exec.Command("./app/VRCDiscordNotifications.exe")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to read discord notifications... %v", err)
	}
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	if scanner.Scan() {
		b := scanner.Bytes()
		var wasPermissionGrantedMessage PermissionGrantedMessage
		err = json.Unmarshal(b, &wasPermissionGrantedMessage)
		if err != nil {
			log.Fatalf("Failed to get permissions to read notifications... %v", err)
		}

		if !wasPermissionGrantedMessage.WasPermissionGranted {
			log.Fatal("Failed to get permissions to read notifications.")
		}
	}
	sender := osc.NewClient("localhost", 9000)

	hideNotifTimer := 10
	go func() {
		for {
			if hideNotifTimer > 0 {
				hideNotifTimer--
				time.Sleep(time.Second * 1)
				continue
			}
			log.Println("Hiding notification")
			msg := osc.NewMessage(isShowingNotificationParameter)
			msg.Append(false)
			sender.Send(msg)
			hideNotifTimer = 10
			time.Sleep(time.Second * 1)
		}
	}()

	for scanner.Scan() {
		b := scanner.Bytes()
		var dnm DiscordNotificationMessage
		err = json.Unmarshal(b, &dnm)
		if err != nil {
			log.Fatalf("Failed to read discord notification %v", err)
		}
		hideNotifTimer = 10
		log.Printf("Message Recieved from: %v", dnm.Username)
		payload := SerializeToVRCFloatArr(dnm.Username)
		msg := osc.NewMessage(isShowingNotificationParameter)
		msg.Append(true)
		sender.Send(msg)
		for i, n := range payload {
			msg := osc.NewMessage(characterParameters[i])
			msg.Append(n)
			sender.Send(msg)
		}
	}
	cmd.Wait()
}

func SerializeToVRCFloatArr(str string) []float32 {
	l := make([]float32, characterLimit)
	str = strings.Trim(str, " ")
	str = strings.ToLower(str)
	strArr := strings.Split(str, "")
	for i := 0; i < characterLimit; i++ {
		if i >= len(strArr) {
			l[i] = .0
			continue
		}
		if f, ok := letterMap[strArr[i]]; ok {
			l[i] = f
		} else {
			l[i] = .0
		}
	}
	return l
}
