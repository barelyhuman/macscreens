# macscreens

A tiny implementation to save and restore multi-monitor layout configs for Mac.


## MVP 
- [x] Save the current display configuration
- [x] Restore a display configuration
- [x] List all available configurations 


## Installation 
Due to the dependency on the native CoreGraphic libraries on macOS, you'll need a few prerequisites on your system 

### Prerequisites
- `go>=1.16`, you'll need a minimum of golang version 1.16 
- A macOS compatible system (Book Air, Book Pro, Studio Pro, Mini, etc)

### Build 

- Clone this repository 
```js
$ git clone https://github.com/barelyhuman/macscreens.git macscreens
```

- change the active directory to `macscreens` 
```js
$ cd macscreens
```

- build and install it using `go`
```js
$ go mod tidy; go build; go install
```

`go mod tidy` is optional since the binary has no dependencies right now in the MVP stage, we will have a few deps in the future for an aesthetic CLI. 


## Usage 

```sh
Usage of macscreens:
  -apply NAME
    	apply the configuration with the name NAME
  -list
    	list all saved configurations
  -save NAME
    	Save the current configuration as NAME
```

**--save**
1. Go into your display settings from _System Preferences_ > _Displays_ > _Arrange_ and use the UI to define a layout for your Monitors
2. In a terminal, you can now save this layout using `macscreens`

```
$ macscreens --save bottom-main-display
```

3. This will now be saved in `~/.config/macscreens/bottom-main-display.json` with the X and Y coordinates of each monitor.

**--apply**
1. You can now restore any saved configuration directly from the terminal using 
```
$ macscreens --apply bottom-main-display
``` 


## License 
[MIT](/LICENSE)