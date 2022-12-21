Yet another implementation of CHIP-8 in golang.

### Dependencies
SDL2 for golang
```
go get -u github.com/veandco/go-sdl2/sdl
```

Since the package above is a wrapper over SDL2, you also need the original SDL2 and SDL2_TTF libraries. Installation depends on your system. For example in archlinux:
```
sudo pacman -S sdl2 sdl2_ttf
```
### Build
```
make build
```
Keep in mind that this may take some time because of the cgo. 

### Usage
```
bin/miya --fname Pong.ch8
```
Key bindings:
```
|1 2 3 4|	  |1 2 3 C|
|Q W E R|     |4 5 6 D|
|A S D F|     |7 8 9 E|
|Z X C V|     |A 0 B F|
```

#### Additional options:
```
bin/miya --fname Pong.ch8 --delay 2
```
Delay description

```
bin/miya --fname Pong.ch8 --background-color 0x000000FF --pixel-color 0xFFFFFFFF
```
Colors description

```
bin/miya --fname Pong.ch8  --delay 10 --debug-mode
```
Debug mode description

