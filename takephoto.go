package main

import "os/exec"

func takePhoto() (filename string, err error) {
	filename = "/tmp/pesik_photo.jpg"
	cmd := exec.Command("fswebcam", "-r 640x480", "--no-banner", "--no-timestamp", filename)
	err = cmd.Run()
	return
}
