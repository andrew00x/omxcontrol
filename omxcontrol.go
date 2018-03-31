package omxcontrol

import (
	"errors"
	"fmt"
	"github.com/godbus/dbus"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
	"time"
)

const (
	playerInterface = "org.mpris.MediaPlayer2.Player"
	propertyGetter  = "org.freedesktop.DBus.Properties.Get"
	propertySetter  = "org.freedesktop.DBus.Properties.Set"
)

type OmxCtrl struct {
	conn      *dbus.Conn
	omxPlayer dbus.BusObject
}

type Stream struct {
	Index    int    `json:"index"`
	Language string `json:"lang"`
	Name     string `json:"name"`
	Codec    string `json:"codec"`
	Active   bool   `json:"active"`
}

func Create() (*OmxCtrl, error) {
	user := os.Getenv("USER")
	address, err := ioutil.ReadFile(fmt.Sprintf("/tmp/omxplayerdbus.%s", user))
	if err != nil {
		return nil, err
	}
	pid, err := ioutil.ReadFile(fmt.Sprintf("/tmp/omxplayerdbus.%s.pid", user))
	if err != nil {
		return nil, err
	}
	os.Setenv("DBUS_SESSION_BUS_ADDRESS", string(address))
	os.Setenv("DBUS_SESSION_BUS_PID", string(pid))
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}

	omxPlayer := conn.Object("org.mpris.MediaPlayer2.omxplayer", dbus.ObjectPath("/org/mpris/MediaPlayer2"))
	return &OmxCtrl{conn: conn, omxPlayer: omxPlayer}, nil
}

func (ctrl *OmxCtrl) Action(action KeyboardAction) error {
	return ctrl.omxPlayer.Call(methodFullName("Action"), 0, action).Err
}

func (ctrl *OmxCtrl) AudioTracks() (audios []Stream, err error) {
	var raw []string
	err = ctrl.omxPlayer.Call(methodFullName("ListAudio"), 0).Store(&raw)
	if err == nil {
		audios = make([]Stream, 0, len(raw))
		for _, s := range raw {
			audios = append(audios, parseStreamInfo(s))
		}
	}
	return
}

func (ctrl *OmxCtrl) CanControl() (res bool, err error) {
	err = ctrl.omxPlayer.Call(propertyGetter, 0, playerInterface, "CanControl").Store(&res)
	return
}

func (ctrl *OmxCtrl) Close() error {
	return ctrl.conn.Close()
}

func (ctrl *OmxCtrl) Duration() (duration time.Duration, err error) {
	var v int64
	err = ctrl.omxPlayer.Call(propertyGetter, 0, playerInterface, "Duration").Store(&v)
	if err == nil {
		duration = time.Duration(v) * time.Microsecond
	}
	return
}

func (ctrl *OmxCtrl) HideSubtitles() error {
	return ctrl.omxPlayer.Call(methodFullName("HideSubtitles"), 0).Err
}

func (ctrl *OmxCtrl) Mute() error {
	return ctrl.omxPlayer.Call(methodFullName("Mute"), 0).Err
}

func (ctrl *OmxCtrl) PlaybackStatus() (status Status, err error) {
	var r string
	err = ctrl.omxPlayer.Call(propertyGetter, 0, playerInterface, "PlaybackStatus").Store(&r)
	if err == nil {
		switch r {
		case "Playing":
			status = Playing
		case "Paused":
			status = Paused
		default:
			status = Unknown
		}
	}
	return
}

func (ctrl *OmxCtrl) Playing() (playing string, err error) {
	err = ctrl.omxPlayer.Call(methodFullName("GetSource"), 0).Store(&playing)
	return
}

func (ctrl *OmxCtrl) Pause() error {
	return ctrl.omxPlayer.Call(methodFullName("Pause"), 0).Err
}

func (ctrl *OmxCtrl) Play() error {
	return ctrl.omxPlayer.Call(methodFullName("Play"), 0).Err
}

func (ctrl *OmxCtrl) PlayPause() error {
	return ctrl.omxPlayer.Call(methodFullName("PlayPause"), 0).Err
}

func (ctrl *OmxCtrl) Position() (pos time.Duration, err error) {
	var v int64
	err = ctrl.omxPlayer.Call(propertyGetter, 0, playerInterface, "Position").Store(&v)
	if err == nil {
		pos = time.Duration(v) * time.Microsecond
	}
	return
}

func (ctrl *OmxCtrl) Seek(offset time.Duration) (err error) {
	var res int64
	err = ctrl.omxPlayer.Call(methodFullName("Seek"), 0, int64(offset/time.Microsecond)).Store(&res)
	if err == nil {
		if res == 0 {
			err = errors.New(fmt.Sprintf("invalid seek offset: %d", offset))
		}
	}
	return
}

func (ctrl *OmxCtrl) SelectAudio(index int) (res bool, err error) {
	err = ctrl.omxPlayer.Call(methodFullName("SelectAudio"), 0, index).Store(&res)
	return
}

func (ctrl *OmxCtrl) SelectSubtitle(index int) (res bool, err error) {
	err = ctrl.omxPlayer.Call(methodFullName("SelectSubtitle"), 0, index).Store(&res)
	return
}

func (ctrl *OmxCtrl) SetPosition(position time.Duration) (err error) {
	var res int64
	err = ctrl.omxPlayer.Call(methodFullName("SetPosition"), 0, dbus.ObjectPath("/"), int64(position/time.Microsecond)).Store(&res)
	if err == nil {
		if position != 0 && res == 0 {
			err = errors.New(fmt.Sprintf("invalid possition: %d", position))
		}
	}
	return
}

func (ctrl *OmxCtrl) ShowSubtitles() error {
	return ctrl.omxPlayer.Call(methodFullName("ShowSubtitles"), 0).Err
}

func (ctrl *OmxCtrl) SetVolume(vol float64) (res float64, err error) {
	err = ctrl.omxPlayer.Call(propertySetter, 0, playerInterface, "Volume", vol).Store(&res)
	return
}

func (ctrl *OmxCtrl) Stop() error {
	return ctrl.omxPlayer.Call(methodFullName("Stop"), 0).Err
}

func (ctrl *OmxCtrl) Subtitles() (subtitles []Stream, err error) {
	var raw []string
	err = ctrl.omxPlayer.Call(methodFullName("ListSubtitles"), 0).Store(&raw)
	if err == nil {
		subtitles = make([]Stream, 0, len(raw))
		for _, s := range raw {
			subtitles = append(subtitles, parseStreamInfo(s))
		}
	}
	return
}

func (ctrl *OmxCtrl) Unmute() error {
	return ctrl.omxPlayer.Call(methodFullName("Unmute"), 0).Err
}

func (ctrl *OmxCtrl) Volume() (vol float64, err error) {
	err = ctrl.omxPlayer.Call(propertyGetter, 0, playerInterface, "Volume").Store(&vol)
	return
}

func methodFullName(shortName string) string {
	return fmt.Sprintf("%s.%s", playerInterface, shortName)
}

func parseStreamInfo(s string) Stream {
	tokens := strings.Split(s, ":")
	index, _ := strconv.Atoi(tokens[0])
	language, name, codec, active := tokens[1], tokens[2], tokens[3], tokens[4] == "active"
	return Stream{Index: index, Language: language, Name: name, Codec: codec, Active: active}
}
