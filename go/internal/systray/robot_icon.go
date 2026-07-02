//go:build windows

package systray

import (
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

// GDI & User32 API declarations for custom primitive drawing
var (
	gdi32                   = windows.NewLazySystemDLL("gdi32.dll")
	pCreateBitmap           = gdi32.NewProc("CreateBitmap")
	pCreateCompatibleDC     = gdi32.NewProc("CreateCompatibleDC")
	pCreateCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
	pSelectObject           = gdi32.NewProc("SelectObject")
	pDeleteDC               = gdi32.NewProc("DeleteDC")
	pDeleteObject           = gdi32.NewProc("DeleteObject")
	pCreateSolidBrush       = gdi32.NewProc("CreateSolidBrush")
	pCreatePen              = gdi32.NewProc("CreatePen")
	pRoundRect              = gdi32.NewProc("RoundRect")
	pEllipse                = gdi32.NewProc("Ellipse")
	pPatBlt                 = gdi32.NewProc("PatBlt")

	pCreateIconIndirect = user32.NewProc("CreateIconIndirect")
	pGetDC              = user32.NewProc("GetDC")
	pReleaseDC          = user32.NewProc("ReleaseDC")
)

type ICONINFO struct {
	FIcon    int32
	XHotspot uint32
	YHotspot uint32
	HbmMask  windows.Handle
	HbmColor windows.Handle
}

const (
	ICON_SIZE = 16 // Standard tray icon size is 16x16 pixels
)

var (
	robotNormalIcon windows.Handle
	robotAlertIcon  windows.Handle
	robotIconOnce   sync.Once
)

// initRobotIcons generates custom GDI-rendered robot icons.
func initRobotIcons() {
	robotIconOnce.Do(func() {
		robotNormalIcon = drawCustomIcon("normal")
		robotAlertIcon = drawCustomIcon("alert")
	})
}

// drawCustomIcon creates a 16x16 tray icon programmatically using GDI primitives
func drawCustomIcon(style string) windows.Handle {
	hDC, _, _ := pGetDC.Call(0)
	if hDC == 0 {
		return 0
	}
	defer pReleaseDC.Call(0, hDC)

	hMemDC, _, _ := pCreateCompatibleDC.Call(hDC)
	if hMemDC == 0 {
		return 0
	}
	defer pDeleteDC.Call(hMemDC)

	// Create color bitmap matching display compatibility
	hbmColor, _, _ := pCreateCompatibleBitmap.Call(hDC, ICON_SIZE, ICON_SIZE)
	if hbmColor == 0 {
		return 0
	}
	defer pDeleteObject.Call(hbmColor)

	hOldObj, _, _ := pSelectObject.Call(hMemDC, hbmColor)
	defer pSelectObject.Call(hMemDC, hOldObj)

	// Create monochrome mask bitmap
	maskBits := make([]byte, (ICON_SIZE*ICON_SIZE)/8)
	for i := range maskBits {
		maskBits[i] = 0xFF // Default to transparent
	}
	hbmMask, _, _ := pCreateBitmap.Call(ICON_SIZE, ICON_SIZE, 1, 1, uintptr(unsafe.Pointer(&maskBits[0])))
	if hbmMask == 0 {
		return 0
	}
	defer pDeleteObject.Call(hbmMask)

	// 1. Draw head background block (Vibrant Dark Blue/Grey: RGB 30, 41, 59)
	hBrushBg, _, _ := pCreateSolidBrush.Call(0x003B291E)
	hOldBrush, _, _ := pSelectObject.Call(hMemDC, hBrushBg)
	defer func() {
		pSelectObject.Call(hMemDC, hOldBrush)
		pDeleteObject.Call(hBrushBg)
	}()

	hPen, _, _ := pCreatePen.Call(0, 1, 0x00475569) // border
	hOldPen, _, _ := pSelectObject.Call(hMemDC, hPen)
	defer func() {
		pSelectObject.Call(hMemDC, hOldPen)
		pDeleteObject.Call(hPen)
	}()

	// Draw head block
	pRoundRect.Call(hMemDC, 1, 2, ICON_SIZE-1, ICON_SIZE-1, 4, 4)

	// 2. Draw glowing status indicators (eyes & antenna bulb)
	var glowColor uintptr
	if style == "alert" {
		glowColor = 0x000000FF // Crimson Red (BGR: 00 00 FF)
	} else {
		glowColor = 0x00FFFF00 // Neon Cyan/Aqua (BGR: FF FF 00)
	}

	hGlowBrush, _, _ := pCreateSolidBrush.Call(glowColor)
	pSelectObject.Call(hMemDC, hGlowBrush)
	defer pDeleteObject.Call(hGlowBrush)

	hNoPen, _, _ := pCreatePen.Call(5, 1, 0) // PS_NULL
	pSelectObject.Call(hMemDC, hNoPen)
	defer pDeleteObject.Call(hNoPen)

	// Left Eye: (4, 6) to (7, 9)
	pEllipse.Call(hMemDC, 3, 6, 6, 9)
	// Right Eye: (9, 6) to (12, 9)
	pEllipse.Call(hMemDC, 10, 6, 13, 9)

	// Mouth bar: (5, 11) to (11, 13)
	pRoundRect.Call(hMemDC, 5, 11, 11, 13, 1, 1)

	// Antenna Bulb: (7, 0) to (9, 2)
	pEllipse.Call(hMemDC, 7, 0, 9, 2)

	// 3. Build monochrome mask DC to configure transparency bounds
	hMaskDC, _, _ := pCreateCompatibleDC.Call(0)
	if hMaskDC != 0 {
		hOldMask, _, _ := pSelectObject.Call(hMaskDC, hbmMask)

		// PatBlt with WHITENESS (all bits set to 1 = transparent background)
		pPatBlt.Call(hMaskDC, 0, 0, ICON_SIZE, ICON_SIZE, 0x00F00021) // WHITENESS

		// Black Brush/Pen (0 = opaque region for the head shape)
		hBlackBrush, _, _ := pCreateSolidBrush.Call(0x00000000)
		pSelectObject.Call(hMaskDC, hBlackBrush)
		defer pDeleteObject.Call(hBlackBrush)

		hBlackPen, _, _ := pCreatePen.Call(0, 1, 0x00000000)
		pSelectObject.Call(hMaskDC, hBlackPen)
		defer pDeleteObject.Call(hBlackPen)

		// Opaque head region matching the head drawing bounds
		pRoundRect.Call(hMaskDC, 1, 2, ICON_SIZE-1, ICON_SIZE-1, 4, 4)

		pSelectObject.Call(hMaskDC, hOldMask)
		pDeleteDC.Call(hMaskDC)
	}

	var ii ICONINFO
	ii.FIcon = 1
	ii.HbmMask = windows.Handle(hbmMask)
	ii.HbmColor = windows.Handle(hbmColor)

	hIcon, _, _ := pCreateIconIndirect.Call(uintptr(unsafe.Pointer(&ii)))
	return windows.Handle(hIcon)
}

func getNormalIcon() windows.Handle {
	initRobotIcons()
	return robotNormalIcon
}

func getAlertIcon() windows.Handle {
	initRobotIcons()
	return robotAlertIcon
}
