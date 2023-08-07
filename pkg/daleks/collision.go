package daleks

// ArePlayersColliding - Are the two sprites colliding?
func ArePlayersColliding(Player1, Player2 *Player) bool {

	return Player1.xPos < Player2.xPos+float64(Player2.GetPlayerImageWidth()) &&
		Player1.xPos+float64(Player1.GetPlayerImageWidth()) > Player2.xPos &&
		Player1.yPos < Player2.yPos+float64(Player2.GetPlayerImageHeight()) &&
		Player1.yPos+float64(Player1.GetPlayerImageHeight()) > Player2.yPos
}

// CheckHeroCollision - Checks if the hero has collided with a robot
func CheckHeroCollision(Player *Player, Robots []*Player) {
	// Check for collisions among sprites
	for index := range Robots {
		if ArePlayersColliding(Player, Robots[index]) {
			Player.isAlive = false
			//fmt.Print("Hero is dead - collision")

		}
	}
}

// CheckRobotsCollision - Check if the Robots are colliding with each other
func CheckRobotsCollision(RobotPlayer []*Player) {
	for i := 0; i < len(RobotPlayer); i++ {
		for j := i + 1; j < len(RobotPlayer); j++ {
			// e.g., perform actions like removing sprites, triggering events, etc.
			if ArePlayersColliding(RobotPlayer[i], RobotPlayer[j]) {
				// Handle collision between sprites[i] and sprites[j]
				RobotPlayer[i].isAlive = false
				RobotPlayer[j].isAlive = false
			}
		}
	}

}

// CheckAllRobotsAlive - Checks if all robots in the slice are alive
func CheckAllRobotsAlive(RobotPlayer []*Player) bool {

	//Get the robots in the slice
	robotAliveVal := len(RobotPlayer)

	// Loop over the robot players
	for index := range RobotPlayer {
		if RobotPlayer[index].isAlive {
			return true
		} else if !RobotPlayer[index].isAlive {
			robotAliveVal -= 1
			if robotAliveVal == 0 {
				return false
			}
		}
	}

	return true
}
