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

	"github.com/donutmonger/game_engine/systems"
	"github.com/donutmonger/game_engine/world"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"math/rand"
	"time"
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

	// Configure global settings
	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.CULL_FACE)
	gl.ClearColor(0.1, 0.1, 0.1, 0.0)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	w := world.NewWorld(10000)

	// Load the texture
	spaceshipTexture, err := texture.NewTextureFromFile("res/textures/spaceship.png")
	if err != nil {
		log.Fatalln(err)
	}
	lazorTexture, err := texture.NewTextureFromFile("res/textures/red.png")
	if err != nil {
		log.Fatalln(err)
	}
	w.CreatePlayerSpaceship(spaceshipTexture, lazorTexture)

	enemyTexture, err := texture.NewTextureFromFile("res/textures/enemy_fighter.png")
	if err != nil {
		log.Fatalln(err)
	}
	w.CreateEnemyFighter(enemyTexture)

	allSystems := make([]systems.System, 0)
	allSystems = append(allSystems, systems.NewRenderSystem(*shaderProgram))
	allSystems = append(allSystems, systems.NewMovementSystem())
	allSystems = append(allSystems, systems.NewPlayerInputSystem(window.GlfwWindow))
	allSystems = append(allSystems, systems.NewParticleCleanupSystem())
	allSystems = append(allSystems, systems.NewParticleEmitterSystem())
	allSystems = append(allSystems, systems.NewCollisionSystem())

	for !window.GlfwWindow.ShouldClose() {

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		for _, system := range allSystems {
			system.Update(&w)
		}

		// Maintenance
		window.GlfwWindow.SwapBuffers()
		glfw.PollEvents()
	}
}
