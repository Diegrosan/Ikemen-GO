//go:build windows

package main

import (
	"github.com/ikemen-engine/Ikemen-GO/packages/glfont"
	glfw "github.com/go-gl/glfw/v3.3/glfw"
)

func (s *System) initRenderer() {
	if s.cfg.Video.RenderMode == "OpenGL 2.1" {
		gfx = &Renderer_GL21{}
		gfxFont = &glfont.FontRenderer_GL21{}
	} else {
		// Default to OpenGL 3.2 for Windows
		gfx = &Renderer_GL32{}
		gfxFont = &glfont.FontRenderer_GL32{}
	}
}

func (s *System) initWindowHint() {
	if sys.cfg.Video.RenderMode == "OpenGL 3.2" {
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 2)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	} else {
		glfw.WindowHint(glfw.ContextVersionMajor, 2)
		glfw.WindowHint(glfw.ContextVersionMinor, 1)
	}
}