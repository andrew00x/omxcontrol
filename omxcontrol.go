package omxcontrol

import (
	"errors"
	"fmt"
	"github.com/godbus/dbus"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
)

type OmxCtrl struct {
	conn      *dbus.Conn
	omxPlayer dbus.BusObject
}

type Stream struct {
	Index    int
	Language string
	Name     string
	Codec    string
	Active   bool
}

type Status int

const (
	Unknown Status = iota
	Playing
	Paused
)

var statusValues = [...]string{"Unknown", "Playing", "Paused"}

func (s Status) String() string {
	return statusValues[s]
}

type KeyboardAction int

const (
	ActionDecreaseSpeed         KeyboardAction = iota + 1
	ActionIncreaseSpeed
	ActionRewind
	ActionFastForward
	ActionShowInfo
	ActionPreviousAudio
	ActionNextAudio
	ActionPreviousChapter
	ActionNextChapter
	ActionPreviousSubtitle
	ActionNextSubtitle
	ActionToggleSubtitle
	ActionDecreaseSubtitleDelay
	ActionIncreaseSubtitleDelay
	ActionExit
	ActionPlayPause
	ActionDecreaseVolume
	ActionIncreaseVolume
	ActionSeekBackSmall
	ActionSeekForwardSmall
	ActionSeekBackLarge
	ActionSeekForwardLarge
	ActionStep
	ActionBlank
	ActionSeekRelative
	ActionSeekAbsolute
	ActionMoveVideo
	ActionHideVideo
	ActionUnhideVideo
	ActionHideSubtitles
	ActionShowSubtitles
	ActionSetAlpha
	ActionSetAspectMode
	ActionCropVideo
	ActionPause
	ActionPlay
)

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
	return ctrl.omxPlayer.Call("org.mpris.MediaPlayer2.Player.Action", 0, action).Err
}

func (ctrl *OmxCtrl) AudioTracks() (audios []Stream, err error) {
	var raw []string
	err = ctrl.omxPlayer.Call("org.mpris.MediaPlayer2.Player.ListAudio", 0).Store(&raw)
	if err == nil {
		audios = make([]Stream, 0, len(raw))
		for _, s := range raw {
			audios = append(audios, parseStreamInfo(s))
		}
	}
	return
}

func (ctrl *OmxCtrl) Close() error {
	return ctrl.conn.Close()
}

func (ctrl *OmxCtrl) Duration() (duration int64, err error) {
	err = ctrl.omxPlayer.Call("org.freedesktop.DBus.Properties.Get", 0, "org.mpris.MediaPlayer2.Player", "Duration").Store(&duration)
	return
}

func (ctrl *OmxCtrl) HideSubtitles() error {
	return ctrl.omxPlayer.Call("org.mpris.MediaPlayer2.Player.HideSubtitles", 0).Err
}

func (ctrl *OmxCtrl) Mute() error {
	return ctrl.omxPlayer.Call("org.mpris.MediaPlayer2.Player.Mute", 0).Err
}

func (ctrl *OmxCtrl) PlaybackStatus() (status Status, err error) {
	var r string
	err = ctrl.omxPlayer.Call("org.freedesktop.DBus.Properties.Get", 0, "org.mpris.MediaPlayer2.Player", "PlaybackStatus").Store(&r)
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
	err = ctrl.omxPlayer.Call("org.mpris.MediaPlayer2.Player.GetSource", 0).Store(&playing)
	return
}

func (ctrl *OmxCtrl) PlayPause() error {
	return ctrl.omxPlayer.Call("org.mpris.MediaPlayer2.Player.PlayPause", 0).Err
}

func (ctrl *OmxCtrl) Position() (pos int64, err error) {
	err = ctrl.omxPlayer.Call("org.freedesktop.DBus.Properties.Get", 0, "org.mpris.MediaPlayer2.Player", "Position").Store(&pos)
	return
}

func (ctrl *OmxCtrl) Seek(offset int64) (err error) {
	var res int64
	err = ctrl.omxPlayer.Call("org.mpris.MediaPlayer2.Player.Seek", 0, offset).Store(&res)
	if err == nil {
		if res == 0 {
			err = errors.New(fmt.Sprintf("invalid seek offset: %d", offset))
		}
	}
	return
}

func (ctrl *OmxCtrl) SelectAudio(index int) (res bool, err error) {
	err = ctrl.omxPlayer.Call("org.mpris.MediaPlayer2.Player.SelectAudio", 0, index).Store(&res)
	return
}

func (ctrl *OmxCtrl) SelectSubtitle(index int) (res bool, err error) {
	err = ctrl.omxPlayer.Call("org.mpris.MediaPlayer2.Player.SelectSubtitle", 0, index).Store(&res)
	return
}

func (ctrl *OmxCtrl) SetPosition(offset int64) (err error) {
	var res int64
	err = ctrl.omxPlayer.Call("org.mpris.MediaPlayer2.Player.SetPosition", 0, dbus.ObjectPath("/"), offset).Store(&res)
	if err == nil {
		if res == 0 {
			err = errors.New(fmt.Sprintf("invalid possition: %d", offset))
		}
	}
	return
}

func (ctrl *OmxCtrl) ShowSubtitles() error {
	return ctrl.omxPlayer.Call("org.mpris.MediaPlayer2.Player.ShowSubtitles", 0).Err
}

func (ctrl *OmxCtrl) SetVolume(vol float64) (res float64, err error) {
	err = ctrl.omxPlayer.Call("org.freedesktop.DBus.Properties.Set", 0, "org.mpris.MediaPlayer2.Player", "Volume", vol).Store(&res)
	return
}

func (ctrl *OmxCtrl) Stop() error {
	return ctrl.omxPlayer.Call("org.mpris.MediaPlayer2.Player.Stop", 0).Err
}

func (ctrl *OmxCtrl) Subtitles() (subtitles []Stream, err error) {
	var raw []string
	err = ctrl.omxPlayer.Call("org.mpris.MediaPlayer2.Player.ListSubtitles", 0).Store(&raw)
	if err == nil {
		subtitles = make([]Stream, 0, len(raw))
		for _, s := range raw {
			subtitles = append(subtitles, parseStreamInfo(s))
		}
	}
	return
}

func (ctrl *OmxCtrl) Unmute() error {
	return ctrl.omxPlayer.Call("org.mpris.MediaPlayer2.Player.Unmute", 0).Err
}

func (ctrl *OmxCtrl) Volume() (vol float64, err error) {
	err = ctrl.omxPlayer.Call("org.freedesktop.DBus.Properties.Get", 0, "org.mpris.MediaPlayer2.Player", "Volume").Store(&vol)
	return
}

func parseStreamInfo(s string) Stream {
	tokens := strings.Split(s, ":")
	index, _ := strconv.Atoi(tokens[0])
	language, name, codec, active := tokens[1], tokens[2], tokens[3], tokens[4] == "active"
	return Stream{Index: index, Language: language, Name: name, Codec: codec, Active: active}
}
