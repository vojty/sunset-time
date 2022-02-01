package main

import (
	"fmt"
	"time"

	"github.com/kermieisinthehouse/systray"
	"github.com/sixdouglas/suncalc"
	"github.com/skratchdot/open-golang/open"
)

var latitude = 50.073658
var longitude = 14.41854

var icons = map[suncalc.DayTimeName][]byte{
	suncalc.Sunrise: sunriseIcon,
	suncalc.Sunset:  sunsetIcon,
}

func getTimes(timestamp time.Time, latitude float64, longitude float64) map[suncalc.DayTimeName]suncalc.DayTime {
	times := suncalc.GetTimes(timestamp, latitude, longitude)

	return times
}

type event struct {
	time time.Time
	name suncalc.DayTimeName
}

func getNextEvent() event {
	now := time.Now()
	times := getTimes(now, latitude, longitude)

	// today's sunset has already happened -> recompute for tomorrow
	if times[suncalc.Sunset].Time.Sub(now) < 0 {
		times := getTimes(now.AddDate(0, 0, 1), latitude, longitude)
		return event{time: times[suncalc.Sunrise].Time, name: suncalc.Sunrise}
	}

	return event{time: times[suncalc.Sunset].Time, name: suncalc.Sunset}
}

func formatTime(t time.Time) string {
	return fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
}

func updateTray() {
	event := getNextEvent()
	systray.SetTemplateIcon(icons[event.name], icons[event.name])
	systray.SetTitle(formatTime(event.time))
	systray.SetTooltip(fmt.Sprintf("Next %s @ %s", event.name, formatTime(event.time)))
}

func updateMenuDayItem(dateItem *systray.MenuItem, timesItem *systray.MenuItem, time time.Time) {
	dateItem.SetTitle(fmt.Sprintf("%02d.%02d.%d", time.Day(), time.Month(), time.Year()))
	dateItem.Disable()

	times := getTimes(time, latitude, longitude)
	timesItem.SetTitle(fmt.Sprintf("↗ %s     ↘ %s", formatTime(times[suncalc.Sunrise].Time), formatTime(times[suncalc.Sunset].Time)))
	timesItem.Disable()
}

func updateMenu(mTodayDate *systray.MenuItem, mTodayTimes *systray.MenuItem, mTomorrowDate *systray.MenuItem, mTomorrowTimes *systray.MenuItem) {
	updateMenuDayItem(mTodayDate, mTodayTimes, time.Now())
	updateMenuDayItem(mTomorrowDate, mTomorrowTimes, time.Now().AddDate(0, 0, 1))
}

func onReady() {
	updateTray()

	systray.AddMenuItem("Today", "Today").Disable()
	mTodayDate := systray.AddMenuItem("", "")
	mTodayTimes := systray.AddMenuItem("", "")
	systray.AddSeparator()

	systray.AddMenuItem("Tomorrow", "Tomorrow").Disable()
	mTomorrowDate := systray.AddMenuItem("", "")
	mTomorrowTimes := systray.AddMenuItem("", "")
	systray.AddSeparator()

	updateMenu(mTodayDate, mTodayTimes, mTomorrowDate, mTomorrowTimes)

	mAbout := systray.AddMenuItem("About", "About")
	mQuit := systray.AddMenuItem("Quit", "Quit")

	// Menu handlers
	go func() {
		for {
			select {
			case <-mAbout.ClickedCh:
				open.Run("https://github.com/vojty/sunset-time")

			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}

	}()

	// Update tray every minute
	// https://stackoverflow.com/a/16466581
	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				updateTray()
				updateMenu(mTodayDate, mTodayTimes, mTomorrowDate, mTomorrowTimes)
			}
		}
	}()
}

func main() {
	systray.Run(onReady, func() {})
}
