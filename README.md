# go-mandelbrot
A version of the mandelbrot set in Go

###Installation
```
sudo apt-get install libsdl2-dev
go get https://github.com/veandco/go-sdl2
```

###Run
```
go run *.go [theme-name]
e.g.
go run *.go fire
```
Where ```theme-name``` is one of the themes in colour.go, currently:
* ```full-spectrum```
* ```fire```
* ```leaf```

###Usage
Arrow keys to move a round, 'z' and 'x' to zoom, 'a' and 's' in increase/decrease the level of detail

###Pictures
![](https://cloud.githubusercontent.com/assets/5236109/14068923/209fd186-f489-11e5-8f6c-b43d536f6210.png)
![](https://cloud.githubusercontent.com/assets/5236109/14068931/48d78090-f489-11e5-8e93-668d48d43fd2.png)
![](https://cloud.githubusercontent.com/assets/5236109/14075860/af6d4a74-f4d4-11e5-9769-7c32183542b2.png)
![](https://cloud.githubusercontent.com/assets/5236109/14068940/931a2e5a-f489-11e5-8e1e-78e2c80e1f43.png)

###Known issues
When the zoom or detail level is high, there is a fair bit of screen tearing. The benefit of this is that overall on most machines, the image takes less time to render. I need to spend some time fiddling with using SDL2 properly and double buffering the screen.
