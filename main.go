package main

import (
	// #cgo LDFLAGS: -framework CoreGraphics
	// #include <CoreGraphics/CoreGraphics.h>
	"C"
	"encoding/json"
	"fmt"
	"unsafe"
)

import (
	"flag"
	"github.com/s0rg/compflag"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	PointX    int
	PointY    int
	DisplayID C.CGDirectDisplayID
}

func bail(err error) {
	if err != nil {
		panic(err)
	}
}

var baseConfigPath string

func main() {
	var err error
	var displayIds C.uint
	var displayCount C.uint
	var configs []Config

	saveConfigFlag := flag.String("save", "", "Save the current configuration as `NAME`")
	applyConfigFlag := flag.String("apply", "", "apply the configuration with the name `NAME`")

	listConfigFlag := flag.Bool("list", false, "list all saved configurations")

	compflag.Complete()

	flag.Parse()

	baseConfigPath, err = os.UserHomeDir()
	bail(err)

	baseConfigPath = filepath.Join(baseConfigPath, ".config", "macscreens")
	os.MkdirAll(baseConfigPath, os.ModePerm)

	cErr := C.CGGetActiveDisplayList(10, nil, &displayCount)
	if cErr != 0 {
		panic(fmt.Errorf("Failed to operate"))
	}

	cErr = C.CGGetActiveDisplayList(displayCount, &displayIds, &displayCount)
	if cErr != 0 {
		panic(fmt.Errorf("Failed to operate"))
	}

	displays := getDisplays(&displayIds, displayCount)

	for _, displayId := range displays {
		var box C.CGRect
		box = C.CGDisplayBounds(displayId)
		x := int(box.origin.x)
		y := int(box.origin.y)

		a := Config{
			PointX:    x,
			PointY:    y,
			DisplayID: displayId,
		}

		configs = append(configs, a)

	}

	if *listConfigFlag {
		listConfigs()
		return
	}

	saveConfigName := *saveConfigFlag
	applyConfigName := *applyConfigFlag

	if len(saveConfigName) > 0 {
		saveConfig(configs, saveConfigName)
		fmt.Println(">> Saved:", filepath.Join(baseConfigPath, saveConfigName))
		return
	}

	if len(applyConfigName) > 0 {
		applyConfig(applyConfigName)
		fmt.Println(">> Applied:", filepath.Join(baseConfigPath, applyConfigName))
		return
	}

}

func saveConfig(configs []Config, name string) {
	jsonData, err := json.Marshal(configs)
	bail(err)

	toFile := filepath.Join(baseConfigPath, name+".json")
	os.WriteFile(toFile, jsonData, os.ModePerm)
}

func applyConfig(name string) {
	normalizedName := name

	if !strings.HasSuffix(name, ".json") {
		normalizedName = normalizedName + ".json"
	}

	fromFile := filepath.Join(baseConfigPath, normalizedName)
	fileD, err := os.ReadFile(fromFile)
	bail(err)

	configFromFile := []Config{}
	json.Unmarshal(fileD, &configFromFile)

	configRef := getDisplayRef()

	for _, savedConfig := range configFromFile {
		err := C.CGConfigureDisplayOrigin(*configRef, savedConfig.DisplayID, C.int(savedConfig.PointX), C.int(savedConfig.PointY))
		if err != 0 {
			panic(fmt.Errorf("Failed to operate"))
		}
	}

	// Confirm the configuration
	cErr := C.CGCompleteDisplayConfiguration(*configRef, C.kCGConfigurePermanently)
	if cErr != 0 {
		panic(fmt.Errorf("Failed to Apply Configuration"))
	}

}

func getDisplayRef() *C.CGDisplayConfigRef {
	var refConfig C.CGDisplayConfigRef
	err := C.CGBeginDisplayConfiguration(&refConfig)
	if err != 0 {
		panic(fmt.Errorf("Failed to operate"))
	}
	return &refConfig
}

func getDisplays(arr *C.uint, size C.uint) []C.CGDirectDisplayID {
	var slice []C.CGDirectDisplayID
	gSlice := (*[1 << 10]C.uint)(unsafe.Pointer(arr))[:size:size]

	for i := 0; i < int(size); i++ {
		slice = append(slice, gSlice[i])
	}

	return slice
}

func listConfigs() {
	dirEntries, err := os.ReadDir(filepath.Join(baseConfigPath))
	bail(err)
	for _, entry := range dirEntries {
		fname := entry.Name()
		if !strings.HasSuffix(fname, ".json") {
			continue
		}
		fmt.Println("- ", strings.TrimSuffix(fname, ".json"))
	}
}
