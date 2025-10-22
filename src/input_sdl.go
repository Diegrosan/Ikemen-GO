//go:build sdl
// +build sdl

package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	sdl "github.com/veandco/go-sdl2/sdl"
)

const MAX_JOYSTICK_COUNT = 4

type Input struct {
	joysticks        [MAX_JOYSTICK_COUNT]*sdl.Joystick
	gameControllers  [MAX_JOYSTICK_COUNT]*sdl.GameController
	useGameController [MAX_JOYSTICK_COUNT]bool
	deviceIndex      [MAX_JOYSTICK_COUNT]int32
}

type Key = sdl.Keycode
type ModifierKey = sdl.Keymod

const (
	KeyUnknown = sdl.K_UNKNOWN
	KeyEscape  = sdl.K_ESCAPE
	KeyEnter   = sdl.K_RETURN
	KeyInsert  = sdl.K_INSERT
	KeyF12     = sdl.K_F12
)

var KeyToStringLUT = map[sdl.Keycode]string{
	sdl.K_RETURN:       "RETURN",
	sdl.K_ESCAPE:       "ESCAPE",
	sdl.K_BACKSPACE:    "BACKSPACE",
	sdl.K_TAB:          "TAB",
	sdl.K_SPACE:        "SPACE",
	sdl.K_QUOTE:        "QUOTE",
	sdl.K_COMMA:        "COMMA",
	sdl.K_MINUS:        "MINUS",
	sdl.K_PERIOD:       "PERIOD",
	sdl.K_SLASH:        "SLASH",
	sdl.K_0:            "0",
	sdl.K_1:            "1",
	sdl.K_2:            "2",
	sdl.K_3:            "3",
	sdl.K_4:            "4",
	sdl.K_5:            "5",
	sdl.K_6:            "6",
	sdl.K_7:            "7",
	sdl.K_8:            "8",
	sdl.K_9:            "9",
	sdl.K_SEMICOLON:    "SEMICOLON",
	sdl.K_EQUALS:       "EQUALS",
	sdl.K_LEFTBRACKET:  "LBRACKET",
	sdl.K_BACKSLASH:    "BACKSLASH",
	sdl.K_RIGHTBRACKET: "RBRACKET",
	sdl.K_BACKQUOTE:    "BACKQUOTE",
	sdl.K_a:            "a",
	sdl.K_b:            "b",
	sdl.K_c:            "c",
	sdl.K_d:            "d",
	sdl.K_e:            "e",
	sdl.K_f:            "f",
	sdl.K_g:            "g",
	sdl.K_h:            "h",
	sdl.K_i:            "i",
	sdl.K_j:            "j",
	sdl.K_k:            "k",
	sdl.K_l:            "l",
	sdl.K_m:            "m",
	sdl.K_n:            "n",
	sdl.K_o:            "o",
	sdl.K_p:            "p",
	sdl.K_q:            "q",
	sdl.K_r:            "r",
	sdl.K_s:            "s",
	sdl.K_t:            "t",
	sdl.K_u:            "u",
	sdl.K_v:            "v",
	sdl.K_w:            "w",
	sdl.K_x:            "x",
	sdl.K_y:            "y",
	sdl.K_z:            "z",
	sdl.K_CAPSLOCK:     "CAPSLOCK",
	sdl.K_F1:           "F1",
	sdl.K_F2:           "F2",
	sdl.K_F3:           "F3",
	sdl.K_F4:           "F4",
	sdl.K_F5:           "F5",
	sdl.K_F6:           "F6",
	sdl.K_F7:           "F7",
	sdl.K_F8:           "F8",
	sdl.K_F9:           "F9",
	sdl.K_F10:          "F10",
	sdl.K_F11:          "F11",
	sdl.K_F12:          "F12",
	sdl.K_PRINTSCREEN:  "PRINTSCREEN",
	sdl.K_SCROLLLOCK:   "SCROLLLOCK",
	sdl.K_PAUSE:        "PAUSE",
	sdl.K_INSERT:       "INSERT",
	sdl.K_HOME:         "HOME",
	sdl.K_PAGEUP:       "PAGEUP",
	sdl.K_DELETE:       "DELETE",
	sdl.K_END:          "END",
	sdl.K_PAGEDOWN:     "PAGEDOWN",
	sdl.K_RIGHT:        "RIGHT",
	sdl.K_LEFT:         "LEFT",
	sdl.K_DOWN:         "DOWN",
	sdl.K_UP:           "UP",
	sdl.K_NUMLOCKCLEAR: "NUMLOCKCLEAR",
	sdl.K_KP_DIVIDE:    "KP_DIVIDE",
	sdl.K_KP_MULTIPLY:  "KP_MULTIPLY",
	sdl.K_KP_MINUS:     "KP_MINUS",
	sdl.K_KP_PLUS:      "KP_PLUS",
	sdl.K_KP_ENTER:     "KP_ENTER",
	sdl.K_KP_1:         "KP_1",
	sdl.K_KP_2:         "KP_2",
	sdl.K_KP_3:         "KP_3",
	sdl.K_KP_4:         "KP_4",
	sdl.K_KP_5:         "KP_5",
	sdl.K_KP_6:         "KP_6",
	sdl.K_KP_7:         "KP_7",
	sdl.K_KP_8:         "KP_8",
	sdl.K_KP_9:         "KP_9",
	sdl.K_KP_0:         "KP_0",
	sdl.K_KP_PERIOD:    "KP_PERIOD",
	sdl.K_KP_EQUALS:    "KP_EQUALS",
	sdl.K_F13:          "F13",
	sdl.K_F14:          "F14",
	sdl.K_F15:          "F15",
	sdl.K_F16:          "F16",
	sdl.K_F17:          "F17",
	sdl.K_F18:          "F18",
	sdl.K_F19:          "F19",
	sdl.K_F20:          "F20",
	sdl.K_F21:          "F21",
	sdl.K_F22:          "F22",
	sdl.K_F23:          "F23",
	sdl.K_F24:          "F24",
	sdl.K_MENU:         "MENU",
	sdl.K_LCTRL:        "LCTRL",
	sdl.K_LSHIFT:       "LSHIFT",
	sdl.K_LALT:         "LALT",
	sdl.K_LGUI:         "LGUI",
	sdl.K_RCTRL:        "RCTRL",
	sdl.K_RSHIFT:       "RSHIFT",
	sdl.K_RALT:         "RALT",
	sdl.K_RGUI:         "RGUI",
}

var StringToKeyLUT = map[string]sdl.Keycode{}

func init() {
	sdl.JoystickEventState(sdl.ENABLE)
	for k, v := range KeyToStringLUT {
		StringToKeyLUT[v] = k
	}
	
	if err := sdl.InitSubSystem(sdl.INIT_GAMECONTROLLER); err != nil {
		fmt.Printf("[init] Failed to initialize GameController: %v\n", err)
	}
	
	AddCommonMappings()
	LoadGameControllerDB()
}

func AddCommonMappings() {
	fmt.Printf("[AddCommonMappings] Loading built-in controller mappings...\n")
	
	commonMappings := []string{
		"030000005e0400008e02000014010000,Xbox 360 Controller,a:b0,b:b1,back:b6,dpdown:h0.4,dpleft:h0.8,dpright:h0.2,dpup:h0.1,guide:b8,leftshoulder:b4,leftstick:b9,lefttrigger:a2,leftx:a0,lefty:a1,rightshoulder:b5,rightstick:b10,righttrigger:a5,rightx:a3,righty:a4,start:b7,x:b2,y:b3,platform:Linux,",
		"030000005e040000ea02000001030000,Xbox One Controller,a:b0,b:b1,back:b6,dpdown:h0.4,dpleft:h0.8,dpright:h0.2,dpup:h0.1,guide:b8,leftshoulder:b4,leftstick:b9,lefttrigger:a2,leftx:a0,lefty:a1,rightshoulder:b5,rightstick:b10,righttrigger:a5,rightx:a3,righty:a4,start:b7,x:b2,y:b3,platform:Linux,",
		"030000004c050000c405000000010000,PS4 Controller,a:b1,b:b2,back:b8,dpdown:h0.4,dpleft:h0.8,dpright:h0.2,dpup:h0.1,guide:b12,leftshoulder:b4,leftstick:b10,lefttrigger:a3,leftx:a0,lefty:a1,rightshoulder:b5,rightstick:b11,righttrigger:a4,rightx:a2,righty:a5,start:b9,x:b0,y:b3,platform:Linux,",
		"030000004c050000e60c000000010000,PS5 Controller,a:b1,b:b2,back:b8,dpdown:h0.4,dpleft:h0.8,dpright:h0.2,dpup:h0.1,guide:b12,leftshoulder:b4,leftstick:b10,lefttrigger:a3,leftx:a0,lefty:a1,rightshoulder:b5,rightstick:b11,righttrigger:a4,rightx:a2,righty:a5,start:b9,x:b0,y:b3,platform:Linux,",
		"030000007e0500000920000001000000,Nintendo Switch Pro Controller,a:b0,b:b1,back:b8,dpdown:h0.4,dpleft:h0.8,dpright:h0.2,dpup:h0.1,guide:b12,leftshoulder:b4,leftstick:b10,lefttrigger:b6,leftx:a0,lefty:a1,rightshoulder:b5,rightstick:b11,righttrigger:b7,rightx:a2,righty:a3,start:b9,x:b2,y:b3,platform:Linux,",
		"03000000790000001100000000000000,Generic USB Controller,a:b2,b:b1,back:b8,dpdown:h0.4,dpleft:h0.8,dpright:h0.2,dpup:h0.1,leftshoulder:b4,leftstick:b10,lefttrigger:b6,leftx:a0,lefty:a1,rightshoulder:b5,rightstick:b11,righttrigger:b7,rightx:a3,righty:a2,start:b9,x:b3,y:b0,platform:Linux,",
	}
	
	count := 0
	for _, mapping := range commonMappings {
		if sdl.GameControllerAddMapping(mapping) >= 0 {
			count++
		}
	}
	
	fmt.Printf("[AddCommonMappings] Loaded %d built-in mappings (will be overridden by gamecontrollerdb.txt if present)\n", count)
}

func LoadGameControllerDB() {
	envPath := os.Getenv("SDL_GAMECONTROLLERDB")
	
	if envPath == "" {
		fmt.Printf("[LoadGameControllerDB] SDL_GAMECONTROLLERDB environment variable not set\n")
		fmt.Printf("[LoadGameControllerDB] Using built-in mappings as fallback\n")
		return
	}
	
	fmt.Printf("[LoadGameControllerDB] Searching for gamecontrollerdb.txt in:\n")
	fmt.Printf("[LoadGameControllerDB]   - %s (SDL_GAMECONTROLLERDB env variable)\n", envPath)
	
	if _, err := os.Stat(envPath); err != nil {
		fmt.Printf("[LoadGameControllerDB] File not found: %s\n", envPath)
		fmt.Printf("[LoadGameControllerDB] Using built-in mappings as fallback\n")
		return
	}
	
	data, err := os.ReadFile(envPath)
	if err != nil {
		fmt.Printf("[LoadGameControllerDB] Failed to read file: %v\n", err)
		fmt.Printf("[LoadGameControllerDB] Using built-in mappings as fallback\n")
		return
	}
	
	lines := strings.Split(string(data), "\n")
	count := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		if sdl.GameControllerAddMapping(line) >= 0 {
			count++
		}
	}
	
	fmt.Printf("[LoadGameControllerDB] Loaded %d mappings from: %s\n", count, envPath)
	fmt.Printf("[LoadGameControllerDB] External database has PRIORITY over built-in mappings\n")
}

func StringToKey(s string) sdl.Keycode {
	if key, ok := StringToKeyLUT[s]; ok {
		return key
	}
	return sdl.K_UNKNOWN
}

func KeyToString(k sdl.Keycode) string {
	if s, ok := KeyToStringLUT[k]; ok {
		return s
	}
	return ""
}

func NewModifierKey(ctrl, alt, shift bool) (mod ModifierKey) {
	if ctrl {
		mod |= sdl.KMOD_CTRL
	}
	if alt {
		mod |= sdl.KMOD_ALT
	}
	if shift {
		mod |= sdl.KMOD_SHIFT
	}
	return mod
}

var input Input

func (input *Input) GetMaxJoystickCount() int {
	return len(input.joysticks)
}

func (input *Input) IsJoystickPresent(joy int) bool {
	if joy < 0 || joy >= len(input.joysticks) {
		return false
	}
	if input.gameControllers[joy] != nil {
		return true
	}
	if input.joysticks[joy] != nil {
		return input.joysticks[joy].Attached()
	}
	return false
}

func (input *Input) GetJoystickName(joy int) string {
	if joy < 0 || joy >= len(input.joysticks) {
		return ""
	}
	if input.gameControllers[joy] != nil {
		return input.gameControllers[joy].Name()
	}
	if input.joysticks[joy] != nil {
		return input.joysticks[joy].Name()
	}
	return ""
}

func (input *Input) GetJoystickAxis(joy int, axis int) int16 {
	if joy < 0 || joy >= len(input.joysticks) {
		return 0
	}
	if input.joysticks[joy] != nil {
		return input.joysticks[joy].Axis(axis)
	}
	return 0
}

func (input *Input) GetJoystickAxes(joy int) []int16 {
	if joy < 0 || joy >= len(input.joysticks) {
		return []int16{}
	}
	if input.joysticks[joy] == nil {
		return []int16{}
	}
	numAxes := input.joysticks[joy].NumAxes()
	axes := make([]int16, numAxes)
	for i := 0; i < numAxes; i++ {
		axes[i] = input.joysticks[joy].Axis(i)
	}
	return axes
}

func (input *Input) GetJoystickButtons(joy int) []byte {
	if joy < 0 || joy >= len(input.joysticks) {
		return []byte{}
	}
	if input.joysticks[joy] == nil {
		return []byte{}
	}
	numButtons := input.joysticks[joy].NumButtons()
	buttons := make([]byte, numButtons)
	for i := 0; i < numButtons; i++ {
		buttons[i] = input.joysticks[joy].Button(i)
	}
	return buttons
}

func (input *Input) GetJoystickHats(joy int) []byte {
	if joy < 0 || joy >= len(input.joysticks) {
		return []byte{}
	}
	if input.joysticks[joy] == nil {
		return []byte{}
	}
	numHats := input.joysticks[joy].NumHats()
	hats := make([]byte, numHats)
	for i := 0; i < numHats; i++ {
		hats[i] = input.joysticks[joy].Hat(i)
	}
	return hats
}

func getAxisValue(value int16, sensitivity int16) int16 {
	if value > sensitivity {
		return value
	} else if value < -sensitivity {
		return value
	}
	return 0
}

func JoystickState(joy, button int) bool {
	if joy < 0 {
		return sys.keyState[Key(button)]
	}
	if joy >= input.GetMaxJoystickCount() {
		return false
	}
	
	js := input.joysticks[joy]
	if js == nil {
		return false
	}
	
	isDpadButton := button == sys.joystickConfig[joy].dU ||
		button == sys.joystickConfig[joy].dR ||
		button == sys.joystickConfig[joy].dD ||
		button == sys.joystickConfig[joy].dL
	
	if isDpadButton {
		if js.NumHats() > 0 {
			hatValue := js.Hat(0)
			switch button {
			case sys.joystickConfig[joy].dU:
				if hatValue&sdl.HAT_UP != 0 {
					return true
				}
			case sys.joystickConfig[joy].dR:
				if hatValue&sdl.HAT_RIGHT != 0 {
					return true
				}
			case sys.joystickConfig[joy].dD:
				if hatValue&sdl.HAT_DOWN != 0 {
					return true
				}
			case sys.joystickConfig[joy].dL:
				if hatValue&sdl.HAT_LEFT != 0 {
					return true
				}
			}
		}
		
		if js.NumAxes() >= 2 {
			axis0 := getAxisValue(js.Axis(0), sys.controllerStickSensitivitySDL)
			axis1 := getAxisValue(js.Axis(1), sys.controllerStickSensitivitySDL)
			
			switch button {
			case sys.joystickConfig[joy].dU:
				if axis1 < 0 {
					return true
				}
			case sys.joystickConfig[joy].dR:
				if axis0 > 0 {
					return true
				}
			case sys.joystickConfig[joy].dD:
				if axis1 > 0 {
					return true
				}
			case sys.joystickConfig[joy].dL:
				if axis0 < 0 {
					return true
				}
			}
		}
		
		if button >= 0 && button < js.NumButtons() {
			return js.Button(button) != 0
		}
		
		return false
	}
	
	if button >= 0 && button < js.NumButtons() {
		return js.Button(button) != 0
	}
	
	if button < 0 {
		var axis int
		isPositive := (button & 1) == 0
		
		if isPositive {
			axis = (-button - 1) / 2
		} else {
			axis = -button / 2
		}
		
		if axis < js.NumAxes() {
			value := getAxisValue(js.Axis(axis), sys.controllerStickSensitivitySDL)
			if isPositive {
				return value > 0
			} else {
				return value < 0
			}
		}
	}
	
	return false
}

func (ir *InputReader) LocalInput(in int) (bool, bool, bool, bool, bool, bool, bool, bool, bool, bool, bool, bool, bool, bool) {
	var U, D, L, R, a, b, c, x, y, z, s, d, w, m bool
	
	if in < len(sys.keyConfig) {
		joy := sys.keyConfig[in].Joy
		if joy == -1 {
			U = sys.keyConfig[in].U()
			D = sys.keyConfig[in].D()
			L = sys.keyConfig[in].L()
			R = sys.keyConfig[in].R()
			a = sys.keyConfig[in].a()
			b = sys.keyConfig[in].b()
			c = sys.keyConfig[in].c()
			x = sys.keyConfig[in].x()
			y = sys.keyConfig[in].y()
			z = sys.keyConfig[in].z()
			s = sys.keyConfig[in].s()
			d = sys.keyConfig[in].d()
			w = sys.keyConfig[in].w()
			m = sys.keyConfig[in].m()
		}
	}
	
	if in < len(sys.joystickConfig) {
		sdl.JoystickUpdate()
		joyS := sys.joystickConfig[in].Joy
		if joyS >= 0 {
			U = U || sys.joystickConfig[in].U()
			D = D || sys.joystickConfig[in].D()
			L = L || sys.joystickConfig[in].L()
			R = R || sys.joystickConfig[in].R()
			a = a || sys.joystickConfig[in].a()
			b = b || sys.joystickConfig[in].b()
			c = c || sys.joystickConfig[in].c()
			x = x || sys.joystickConfig[in].x()
			y = y || sys.joystickConfig[in].y()
			z = z || sys.joystickConfig[in].z()
			s = s || sys.joystickConfig[in].s()
			d = d || sys.joystickConfig[in].d()
			w = w || sys.joystickConfig[in].w()
			m = m || sys.joystickConfig[in].m()
		}
	}
	
	if sys.inputButtonAssist {
		a, b, c, x, y, z, s, d, w = ir.ButtonAssistCheck(a, b, c, x, y, z, s, d, w)
	}
	
	return U, D, L, R, a, b, c, x, y, z, s, d, w, m
}

func checkAxisForDpad(joy int, axes *[]int16, base int) string {
	if len(*axes) < 2 {
		return ""
	}
	
	axis0 := getAxisValue((*axes)[0], sys.controllerStickSensitivitySDL)
	axis1 := getAxisValue((*axes)[1], sys.controllerStickSensitivitySDL)
	
	if axis0 > 0 {
		return strconv.Itoa(2 + base)
	} else if axis0 < 0 {
		return strconv.Itoa(1 + base)
	}
	
	if axis1 > 0 {
		return strconv.Itoa(3 + base)
	} else if axis1 < 0 {
		return strconv.Itoa(base)
	}
	
	return ""
}

func checkAxisForTrigger(joy int, axes *[]int16) string {
	name := input.GetJoystickName(joy) + "." + runtime.GOOS + "." + runtime.GOARCH + ".sdl"
	
	restingTriggers := map[string][]int{
		"XInput Gamepad (GLFW).windows.amd64.sdl":     {4, 5},
		"PS4 Controller.windows.amd64.sdl":            {4, 5},
		"Steam Virtual Gamepad.linux.amd64.glfw":      {2, 5},
		"Steam Deck Controller.linux.amd64.sdl":       {2, 5},
		"PS3 Controller.linux.amd64.sdl":              {2, 5},
		"Logitech Dual Action.linux.amd64.sdl":        {2, 5},
		"Gamepad.linux.arm64.sdl":                     {4, 5},
	}
	
	exceptions := restingTriggers[name]
	
	for i := range *axes {
		value := getAxisValue((*axes)[i], sys.controllerStickSensitivitySDL)
		
		isException := false
		for _, ex := range exceptions {
			if i == ex && value < 0 {
				isException = true
				break
			}
		}
		
		if isException {
			continue
		}
		
		if value < 0 {
			return strconv.Itoa(-i*2 - 1)
		} else if value > 0 {
			return strconv.Itoa(-i*2 - 2)
		}
	}
	
	return ""
}
