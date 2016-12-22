// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3 and OpenGL 4.1 core forward-compatible profile.
package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/donutmonger/game_engine/shader"
	"github.com/donutmonger/game_engine/texture"
	"github.com/donutmonger/game_engine/window"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"time"
	"github.com/donutmonger/game_engine/world"
	"math/rand"
)

const windowWidth = 800
const windowHeight = 600

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	window := window.NewWindow(windowWidth, windowHeight)
	window.GlfwWindow.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure the vertex and fragment shaders
	vertexShader, err := shader.NewVertexShaderFromFile("res/shaders/basic_sprite.vs")
	if err != nil {
		panic(err)
	}
	fragmentShader, err := shader.NewFragmentShaderFromFile("res/shaders/basic_sprite.fs")
	if err != nil {
		panic(err)
	}
	shaderProgram, err := shader.NewShaderProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}
	shaderProgram.Use()

	gl.BindFragDataLocation(shaderProgram.GLid, 0, gl.Str("outputColor\x00"))

	// Load the texture
	enemyTexture, err := texture.NewTextureFromFile("enemy1.png")
	if err != nil {
		log.Fatalln(err)
	}

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(flatMeshVertices)*4, gl.Ptr(flatMeshVertices), gl.STATIC_DRAW)

	posAttrib := uint32(shaderProgram.GetAttribLocation("position"))
	gl.EnableVertexAttribArray(posAttrib)
	gl.VertexAttribPointer(posAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0))

	texcoordAttrib := uint32(shaderProgram.GetAttribLocation("texcoord"))
	gl.EnableVertexAttribArray(texcoordAttrib)
	gl.VertexAttribPointer(texcoordAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))

	// Configure global settings
	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.CULL_FACE)
	gl.ClearColor(0.2, 0.2, 0.4, 0.0)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	w := world.NewWorld(100)
	for i := 0; i < 99; i++ {
		w.CreateEnemy(enemyTexture)
	}

	for !window.GlfwWindow.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Render
		shaderProgram.Use()
		gl.BindVertexArray(vao)

		w.Draw(shaderProgram)

		// Maintenance
		window.GlfwWindow.SwapBuffers()
		glfw.PollEvents()
	}
}

var flatMeshVertices = []float32{
	-1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, 1.0, 0.0,

	1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, 0.0, 1.0,
}
