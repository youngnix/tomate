package main

import (
	_ "embed"
)

var (
	//go:embed res/start_icon.png
	start_icon_img []byte
	//go:embed res/pause_icon.png
	pause_icon_img []byte
	//go:embed res/skip_icon.png
	skip_icon_img []byte
	//go:embed res/stop_icon.png
	stop_icon_img []byte
	//go:embed res/chime_sound.wav
	chime_snd []byte
	//go:embed res/chime_sound_slow.wav
	chime_slow_snd []byte
)
