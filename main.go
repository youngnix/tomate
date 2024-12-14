package main

import (
	"bytes"
	"container/list"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
	"youngnix.com/tomate/internal/config"
)

const (
	TIMER_STATE_PAUSED = iota
	TIMER_STATE_RUNNING
)

type Sound struct {
	streamer beep.StreamSeekCloser
	format   beep.Format
}

func NewSound(wave []byte) (*Sound, error) {
	buffer := bytes.NewBuffer(wave)

	streamer, format, err := wav.Decode(buffer)

	if err != nil {
		return nil, err
	}

	return &Sound{
		streamer,
		format,
	}, nil
}

func (s *Sound) Play() {
	go func() {
		speaker.PlayAndWait(s.streamer)
	}()
}

func (s *Sound) Close() {
	s.streamer.Close()
}

var (
	countdown          float64 = 0
	countdown_string, title_string   binding.String
	current_cycle      *list.Element
	cycles             *list.List
	state              uint = TIMER_STATE_PAUSED
	myApp              fyne.App
	focus_notification, break_notification *fyne.Notification
	focus_sound, rest_sound *Sound
	start_icon, pause_icon, stop_icon, skip_icon fyne.Resource
	start_pause_button, skip_button, stop_button *widget.Button
)

func init() {
	var err error

	focus_sound, err = NewSound(chime_snd)
	if err != nil {
		panic(err)
	}

	rest_sound, err = NewSound(chime_slow_snd)
	if err != nil {
		panic(err)
	}

	start_icon = fyne.NewStaticResource("start_icon", start_icon_img)
	pause_icon = fyne.NewStaticResource("pause_icon", pause_icon_img)
	stop_icon = fyne.NewStaticResource("stop_icon", stop_icon_img)
	skip_icon = fyne.NewStaticResource("skip_icon", skip_icon_img)

	countdown_string = binding.NewString()
	title_string = binding.NewString()
}

func SetCountdown(c float64) {
	countdown = c

	duration := time.Duration(c * float64(time.Second))

	countdown_string.Set(duration.String())
}

func NextCycle() {
	if current_cycle.Next() != nil {
		current_cycle = current_cycle.Next()
	} else {
		current_cycle = cycles.Front()
	}

	SetCountdown(current_cycle.Value.(*config.Cycle).Duration.Seconds())
}

func Countdown() {
	for range time.Tick(time.Second) {
		if state != TIMER_STATE_RUNNING {
			return
		}

		if countdown == 0 {
			NextCycle()

			if (config.Config.Bell) {
				focus_sound.Play()
			}

			if (config.Config.Notify) {
				myApp.SendNotification(current_cycle.Value.(*config.Cycle).Notification)
			}
		}

		title_string.Set(current_cycle.Value.(*config.Cycle).Title)

		SetCountdown(countdown - 1)
	}
}

func ResetCycles() {
	current_cycle = cycles.Front()

	SetCountdown(current_cycle.Value.(*config.Cycle).Duration.Seconds())

	state = TIMER_STATE_PAUSED

	start_pause_button.SetIcon(start_icon)
}

func main() {
	defer focus_sound.Close()
	defer rest_sound.Close()

	myApp = app.New()
	appWindow := myApp.NewWindow("Tomate")

	cycles = list.New()

	for _, c := range config.Config.Cycles {
		cycles.PushBack(c)
	}

	current_cycle = cycles.Front()

	speaker.Init(focus_sound.format.SampleRate, focus_sound.format.SampleRate.N(time.Second/10))

	SetCountdown(current_cycle.Value.(*config.Cycle).Duration.Seconds())

	countdown_label := canvas.NewText("", color.White)

	countdown_label.TextSize = 48
	countdown_label.Alignment = fyne.TextAlignCenter

	countdown_string.AddListener(binding.NewDataListener(func() {
		newValue, _ := countdown_string.Get()
		countdown_label.Text = newValue
		countdown_label.Refresh()
	}))

	title_label := canvas.NewText("", color.White)
	title_label.Alignment = fyne.TextAlignCenter

	title_string.AddListener(binding.NewDataListener(func() {
		newValue, _ := title_string.Get()
		title_label.Text = newValue
		title_label.Refresh()
	}))

	start_pause_button = widget.NewButtonWithIcon("", start_icon, func() {
		switch state {
		case TIMER_STATE_PAUSED:
			state = TIMER_STATE_RUNNING
			start_pause_button.SetIcon(pause_icon)
			go Countdown()
			break
		case TIMER_STATE_RUNNING:
			state = TIMER_STATE_PAUSED
			start_pause_button.SetIcon(start_icon)
			break
		}
	})

	skip_button = widget.NewButtonWithIcon("", skip_icon, NextCycle)
	stop_button = widget.NewButtonWithIcon("", stop_icon, ResetCycles)

	appWindow.SetContent(
		container.NewCenter(
			container.NewVBox(
				countdown_label,
				title_label,
				container.NewCenter(
					container.NewHBox(
						start_pause_button,
						skip_button,
						stop_button,
					),
				),
			),
		),
	)

	appWindow.ShowAndRun()
}
