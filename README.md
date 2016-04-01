# go-mandelbrot
A version of the mandelbrot set in Go

###Installation
```
sudo apt-get install libsdl2-dev
go get github.com/veandco/go-sdl2
go get github.com/tomcraven/go-mandelbrot
```

###Run
```
go run *.go [optional: theme-name]
e.g.
go run *.go fire
```
Where ```theme-name``` is one of the themes in colour.go, currently:
* ```full-spectrum```
* ```fire```
* ```leaf```
* ```water```
* ```beach``` (default)
* ```random```

###Usage
Arrow keys to move a round, 'z' and 'x' to zoom, 'a' and 's' in increase/decrease the level of detail

###Pictures
![](https://cloud.githubusercontent.com/assets/5236109/14201409/a84d9cea-f7e9-11e5-9c40-975420b1d3c8.png)
![](https://cloud.githubusercontent.com/assets/5236109/14201325/10dbec90-f7e9-11e5-9f21-a71066e66029.png)
![](https://cloud.githubusercontent.com/assets/5236109/14075860/af6d4a74-f4d4-11e5-9769-7c32183542b2.png)
![](https://cloud.githubusercontent.com/assets/5236109/14201323/0d784116-f7e9-11e5-93e6-48f640f2e479.png)

###Known issues
When the zoom or detail level is high, there is a fair bit of screen tearing. The benefit of this is that overall on most machines, the image takes less time to render. I need to spend some time fiddling with using SDL2 properly and double buffering the screen.
