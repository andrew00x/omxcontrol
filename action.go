package omxcontrol

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
