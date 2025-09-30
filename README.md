<div align="center">

# GoDaleks v1.0

A modern Go/Ebiten (and faithful) recreation of the classic Apple Macintosh game **Daleks**, itself inspired by Johan Strandberg’s 1984 _Daleks_ and the older BSD UNIX game _Robots_.  
This version keeps the spirit of the original while adding smooth animations, mouse support, and modern gameplay tweaks including sounds.

[![Build Status](https://github.com/AaronSaikovski/godaleks/workflows/build/badge.svg)](https://github.com/AaronSaikovski/godaleks/actions)
![version](https://img.shields.io/badge/version-1.0.0-blue)
[![Licence](https://img.shields.io/github/license/AaronSaikovski/godaleks)](LICENSE)

</div>

## 📜 Background

The version you see here is a modern port designed to run cross-platform on today’s systems, preserving the tension and strategy that made the original game so addictive.
All efforts have been made to make this game as faithful to the original as possible.

This code is written in Go and can be played locally or on the web via a Web assembly compiled version and is fully self contained.

See [here](https://www.macintoshrepository.org/3913-daleks)

---

## 🎮 Gameplay

[GoDaleks- Online Playable version](https://AaronSaikovski.github.io/godaleks)

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

- Smooth Dalek movement
- Cool sounds
- Mouse and keyboard control
- Teleportation effects & Sonic Screwdriver visual effects
- Scrap heaps from Dalek collisions
- **Last Stand mode**: Continuous rush of Daleks for bonus points
- Safe teleport option to avoid instant death
- Optional grid overlay
- Level progression with score bonuses
- Power-ups:
  - **Teleports** (normal & safe)
  - **Screwdrivers**
  - **Last Stands**

---

## 📈 Scoring

- Dalek destroyed by collision: **+2 points**
- Dalek destroyed by screwdriver: **+5 points**
- Level completion: **+10 × level number**
- Surviving a Last Stand: **+50 bonus**

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

## Reporting an issue

Please feel free to lodge an [issue or pull request on GitHub](https://github.com/AaronSaikovski/godaleks/issues).
