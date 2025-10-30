<div align="center">

# GoDaleks v1.2.0 

A modern Go/Ebiten (and faithful) recreation of the classic Apple Macintosh game **Daleks**, itself inspired by Johan Strandberg’s 1984 _Daleks_ and the older BSD UNIX game _Robots_.  
This version keeps the spirit of the original while adding smooth animations, mouse support, and modern gameplay tweaks including sounds.

[![Build Status](https://github.com/AaronSaikovski/godaleks/workflows/build/badge.svg)](https://github.com/AaronSaikovski/godaleks/actions)
![version](https://img.shields.io/badge/version-1.2.0-blue)
[![Licence](https://img.shields.io/github/license/AaronSaikovski/godaleks)](LICENSE)

</div>

## 📜 Background

The version you see here is a modern port designed to run cross-platform on today’s systems, preserving the tension and strategy that made the original game so addictive.
All efforts have been made to make this game as faithful to the original as possible.

This code is written in Go and can be played locally or on the web via a Web assembly compiled version and is fully self contained.

See [here](https://www.macintoshrepository.org/3913-daleks)

---

## 🎮 Gameplay

[GoDaleks- Online Playable version](https://www.godaleks.com)

In this game, you attempt to survive by avoiding steadily converging robots If you are overrun by the robots, or move into their immediate zone of control, you are disintegrated.
By guiding the robots with your actions, you can get them to destroy themselves as they collide with each other.
You can escape by teleporting out of range, or you can destroy adjacent robots once each round with a sonic screwdriver.

![Dalek main screen](./images/daleks1.jpg)

![Dalek game play screen](./images/daleks2.jpg)

Daleks move one step per turn toward you. Survive by making them crash into each other, creating scrap heaps, or by destroying them with your **Sonic Screwdriver**.

**You win a level** when all Daleks are destroyed.  
**You lose** if a Dalek catches you.

---

## 🕹️ Controls

### **Keyboard**

| Key                | Action                                          |
| ------------------ | ----------------------------------------------- |
| Arrow Keys / Mouse | Move up, down, left, right and down             |
| Q / E / Z / C      | Diagonal movement                               |
| `SPACE` or `.`     | Wait in place                                   |
| `T`                | Teleport randomly                               |
| `N`                | Start a New game                                |
| `R`                | Safe teleport (avoid near Daleks)               |
| `S`                | Use Sonic Screwdriver (destroy adjacent Daleks) |
| `L`                | Last Stand (Daleks rush continuously)           |
| `G`                | Toggle grid on/off                              |
| `D`                | Debug info (speed, daleks left, etc.)           |

### **Mouse**

- **Click adjacent cell**: Move there
- **Click on your position**: Wait in place

---

## 🛠 Features

### Core Gameplay
- Smooth Dalek movement with easing animations
- Cool sound effects and audio feedback
- Mouse and keyboard control support
- Scrap heaps from Dalek-to-Dalek collisions
- Optional grid overlay for strategic positioning
- Level progression with increasing difficulty

### Power-ups & Tools
- **Teleports** (normal random & safe guaranteed)
- **Sonic Screwdriver** (destroy adjacent Daleks)
- **Last Stands** (bonus mode - continuous Dalek rush)

### Visual Effects
- Teleportation effects with visual feedback
- Sonic Screwdriver visual effects
- Collision explosion animations
- Smooth interpolated movement between grid cells

### Boss Mechanics (v1.2.0+)
- **Dalek Emperor**: Invulnerable boss character
  - Spawns on level 2+ with 60% chance
  - Crushes regular Daleks on collision
  - Provides strategic challenge as two-phase battle
  - Becomes vulnerable only when all normal Daleks defeated
  - **50 point bonus** for defeating the Emperor

---

## 📈 Scoring

| Action | Points |
|--------|--------|
| Dalek destroyed by collision | +2 |
| Dalek destroyed by Sonic Screwdriver | +5 |
| Regular Dalek crushed by Emperor | +2 |
| **Dalek Emperor defeated** | **+50** |
| Level completion bonus | +10 × level number |
| Last Stand completion bonus | +50 |

---

## 🚀 Building & Running

The toolchain is driven by using [Taskfile](https://taskfile.dev/) and all commands are managed via the file `Taskfile.yml`

The list of commands is as follows:

```bash
* build:            Compiles the code.
* clean:            Cleans the project.
* deps:             Updates/installs and dependencies.
* goreleaser:       Builds using Goreleaser.
* lint:             Lints and tidies up the project.
* release:          Builds a release version (smaller binary) of the project.
* run:              Executes the project.
* seccheck:         Checks for security vulnerabilities in the project.
* staticcheck:      Runs a static check of the project.
* test:             Executes and tests for the project.
* generate:         Updates the project build version.
* vet:              Vet examines Go source code and reports suspicious constructs.
```

## 📋 Game Mechanics

### Level Progression
- **Starting Daleks**: 5 + current level count
- **Difficulty**: Increases with each level
- **Resources per level**:
  - +2 Teleports per level
  - +2 Sonic Screwdrivers per level
  - +1 Last Stand every 5 levels
- **Level completion**: All Daleks must be destroyed

### Dalek Behavior
- **Normal Movement**: Move one grid cell per turn toward the player
- **Animation**: 0.85 second smooth easing animation per move
- **Collision Detection**: Grid-based (50×33 grid, 16 pixels per cell)
- **Collision Rules**:
  - Dalek hits scrap: Dies (+2 points)
  - Dalek hits Dalek: Both die (+2 points each), scrap created
  - Dalek hits Emperor: Dies (+2 points), Emperor survives
  - Emperor hits nothing: Survives until all normal Daleks defeated

### Emperor Mechanics (v1.2.0+)
- **Spawn Conditions**: Level 2+, 60% chance per level
- **Spawn Distance**: At least 3 cells from player
- **Visual**: Unique sprite (red emperor Dalek)
- **Audio**: "dalek_emperor" alert sound on spawn + visual warning
- **Invulnerability**: Cannot be harmed while normal Daleks exist
- **Combat**: Crushes any normal Dalek on collision
- **Defeat**: Dies when all normal Daleks eliminated (+50 bonus points)
- **Scrap Immunity**: Cannot be destroyed by scrap heaps

### Last Stand Mode
- **Activation**: Press `L` key (requires available Last Stand)
- **Behavior**: Daleks accelerate continuously toward player
- **Initial Speed**: 2.0 cells/second
- **Acceleration**: 1.1× per second
- **Maximum Speed**: 8.0 cells/second
- **End Condition**: All Daleks destroyed or player caught
- **Bonus**: +50 points on successful completion

---

## 📝 Changelog

### v1.2.0 - Emperor Update
- **Added**: Dalek Emperor boss character (invulnerable until normal Daleks defeated)
- **Added**: Emperor collision mechanics - crushes regular Daleks
- **Added**: Emperor visual sprite and audio alert
- **Improved**: Two-phase final boss battle system
- **Added**: 50-point bonus for defeating Emperor
- **Fixed**: Collision logic to allow Emperor to actively engage in combat

### v1.1.0 - Performance Update
- **Improved**: Big performance fixes and optimizations
- **Improved**: Better animations and smooth movement
- **Added**: Enhanced visual effects

### v1.0.0 - Initial Release
- Core game mechanics from classic Daleks
- Smooth animations and modern graphics
- Mouse and keyboard controls
- Sound effects and audio

---

## 🔧 Architecture

### Core Files

| File | Purpose |
|------|---------|
| `cmd/types.go` | Entity definitions (Dalek, Game, CollisionEffect) |
| `cmd/collision.go` | Collision detection and resolution system |
| `cmd/game.go` | Game state management and level progression |
| `cmd/movement.go` | Dalek movement and animation logic |
| `cmd/draw.go` | Rendering and visual display |
| `cmd/effect.go` | Player actions (teleport, screwdriver, etc.) |
| `cmd/menu.go` | Menu and UI rendering |
| `cmd/loadsounds.go` | Audio loading and playback |
| `cmd/root.go` | Game loop and update cycle |

### Key Systems

**Collision Detection** (`collision.go`):
- Position-based grid collision detection
- Emperor invulnerability tracking
- Scrap hazard collision
- Player safety checks
- Multi-dalek collision handling

**Movement System** (`movement.go`):
- Grid-based normal movement (0.85s per move)
- Smooth easing interpolation for animation
- Last Stand continuous acceleration mode
- Distance-based collision for continuous movement

**Game State** (`game.go`):
- Level management and progression
- Entity spawning and lifecycle
- Position occupation validation
- Resource tracking (teleports, screwdrivers, Last Stands)

---

## 🐛 Reporting an Issue

Please feel free to lodge an [issue or pull request on GitHub](https://github.com/AaronSaikovski/godaleks/issues).
