### Dependencies
SDL2 and SDL2_ttf for golang
```
go get -u -v github.com/veandco/go-sdl2/sdl
go get -u -v github.com/veandco/go-sdl2/ttf
```

Since the package above is a wrapper over SDL2, you also need the original SDL2 and SDL2_ttf libraries. Installation depends on your system. For example in archlinux:
```
sudo pacman -S sdl2 sdl2_ttf
```
### Build
```
make
```
Keep in mind that this may take some time because of Cgo

### Usage
```
bin/miya --fname Pong.ch8
```
Key bindings:
```
|1 2 3 4|     |1 2 3 C|
|Q W E R|     |4 5 6 D|
|A S D F|     |7 8 9 E|
|Z X C V|     |A 0 B F|
```

#### Additional options:
```
bin/miya --fname Pong.ch8 --delay 2
```
Delay in milliseconds for virtual machine and full rendering cycle

```
bin/miya --fname Pong.ch8 --background-color 0x000000FF --pixel-color 0xFFFFFFFF
```
Colors for the background and the pixels on the *screen*

```
bin/miya --fname Pong.ch8  --delay 10 --debug-mode
```
Run in debug mode. Additional window with registers, stack, etc.

