<div align="center">

# GoDaleks

A modern Go/Ebiten (and faithful) recreation of the classic Mac game **Daleks**, itself inspired by Johan Strandberg’s 1984 _Daleks_ and the older BSD UNIX game _Robots_.  
This version keeps the spirit of the original while adding smooth animations, mouse support, and modern gameplay tweaks.
</div>

## 📜 Background

_"Daleks Forever"_ was originally written by Mike Gleason in 2010 as a tribute to the Mac classic, featuring smooth Dalek movement, advanced gameplay tweaks, and faithful sound/graphics.  
The Go/Ebiten edition you see here is a modern port designed to run cross-platform on today’s systems, preserving the tension and strategy that made the game addictive.

For the history of the original and its Macintosh versions, see [About Daleks Forever](About%20Daleks%20Forever.txt).

---

## 🎮 Gameplay

The player is a lone human on a grid, hunted by deadly Daleks.  
Daleks move one step per turn toward you. Survive by making them crash into each other, creating scrap heaps, or by destroying them with your **Sonic Screwdriver**.

**You win a level** when all Daleks are destroyed.  
**You lose** if a Dalek catches you.

---

## 🕹️ Controls

### **Keyboard**

| Key                | Action                                          |
| ------------------ | ----------------------------------------------- |
| Arrow Keys / H J K | Move up, left, down                             |
| Y / U / B / N      | Diagonal movement                               |
| `SPACE` or `.`     | Wait in place                                   |
| `T`                | Teleport randomly                               |
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

- Smooth Dalek movement animations
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
* docker-build:     Builds a docker image based on the docker file.
* docker-run:       Runs the docker container.
```


## Reporting an issue

Please feel free to lodge an [issue or pull request on GitHub](https://github.com/AaronSaikovski/godaleks/issues).