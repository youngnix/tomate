package main

import (
	"container/list"
	"fmt"
	"image/color"
	"math"
	"os"
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
)

const (
	TIMER_STATE_PAUSED = iota
	TIMER_STATE_RUNNING
)

type Cycle struct {
	length float32
	title  string
}

func (c Cycle) Countdown() float32 {
	return c.length * 60.0
}

type Sound struct {
	streamer beep.StreamSeekCloser
	format   beep.Format
}

func NewSound(path string) (*Sound, error) {
	fp, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	streamer, format, err := wav.Decode(fp)

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
	countdown          float32 = 0
	countdown_string   binding.String
	current_cycle      *list.Element
	cycles             *list.List
	state              uint = TIMER_STATE_PAUSED
	myApp              fyne.App
	focus_notification *fyne.Notification
	break_notification *fyne.Notification
	focus_sound        *Sound
	rest_sound         *Sound
	start_icon         fyne.Resource
	pause_icon         fyne.Resource
	stop_icon          fyne.Resource
	skip_icon          fyne.Resource
)

func init() {
	var err error

	focus_sound, err = NewSound("res/chime_sound.wav")
	if err != nil {
		panic(err)
	}

	rest_sound, err = NewSound("res/chime_sound_slow.wav")
	if err != nil {
		panic(err)
	}

	start_icon, err = fyne.LoadResourceFromPath("res/start_icon.png")

	if err != nil {
		return
	}

	pause_icon, err = fyne.LoadResourceFromPath("res/pause_icon.png")

	if err != nil {
		return
	}

	stop_icon, err = fyne.LoadResourceFromPath("res/stop_icon.png")

	if err != nil {
		return
	}

	skip_icon, err = fyne.LoadResourceFromPath("res/skip_icon.png")

	if err != nil {
		return
	}

	focus_notification = fyne.NewNotification("Time to Focus", "Focus on your tasks.")
	break_notification = fyne.NewNotification("Break Time", "Take a break. Relax and hydrate.")
}

func SetCountdown(c float32) {
	countdown = c
	countdown_string.Set(fmt.Sprintf("%02f:%02f", countdown/60, math.Mod(float64(countdown), 60)))
}

func NextCycle() {
	if current_cycle.Next() != nil {
		current_cycle = current_cycle.Next()
	} else {
		current_cycle = cycles.Front()
	}

	SetCountdown(current_cycle.Value.(Cycle).Countdown())
}

func Countdown() {
	for range time.Tick(time.Second) {
		if state != TIMER_STATE_RUNNING {
			return
		}

		if countdown == 0 {
			NextCycle()
			switch current_cycle.Value.(Cycle).title {
			case "break":
				myApp.SendNotification(break_notification)
				rest_sound.Play()
				break
			case "focus":
				myApp.SendNotification(focus_notification)
				focus_sound.Play()
				break
			}
		}

		SetCountdown(countdown - 1)
	}
}

func main() {
	defer focus_sound.Close()
	defer rest_sound.Close()

	myApp = app.New()
	appWindow := myApp.NewWindow("Tomate")

	cycles = list.New()

	cycles.PushBack(Cycle{title: "focus", length: 60})
	cycles.PushBack(Cycle{title: "break", length: 10})

	current_cycle = cycles.Front()

	speaker.Init(focus_sound.format.SampleRate, focus_sound.format.SampleRate.N(time.Second/10))

	countdown_string = binding.NewString()

	SetCountdown(current_cycle.Value.(Cycle).Countdown())

	var start_pause_button, skip_button, stop_button *widget.Button
	countdown_label := canvas.NewText("", color.White)

	countdown_label.TextSize = 48
	countdown_label.Alignment = fyne.TextAlignCenter

	countdown_string.AddListener(binding.NewDataListener(func() {
		newValue, _ := countdown_string.Get()
		countdown_label.Text = newValue
		countdown_label.Refresh()
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

	ResetCycles := func() {
		current_cycle = cycles.Front()

		SetCountdown(current_cycle.Value.(Cycle).Countdown())

		state = TIMER_STATE_PAUSED

		start_pause_button.SetIcon(start_icon)
	}

	stop_button = widget.NewButtonWithIcon("", stop_icon, ResetCycles)

	appWindow.SetContent(
		container.NewCenter(
			container.NewVBox(
				countdown_label,
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
