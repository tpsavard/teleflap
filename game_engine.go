package main

import (
	"math/rand"
)

const MIN_HEIGHT = 0
const MAX_HEIGHT = 1000
const GATE_SIZE = 150
const GATE_DISTANCE_MS = 4 * 1000

// Init Methods

func get_first_gate() (int, int) {
	const PAD = 50
	var lower = rand.Intn(MAX_HEIGHT-PAD-GATE_SIZE-PAD) + PAD
	return (lower + GATE_SIZE), lower
}

func get_starting_player_position() int {
	return rand.Intn(MAX_HEIGHT)
}

// Game Tick Methods

func get_next_gate(prev_upper int, prev_lower int, elapsed_time_ms int64) (int, int) {
	if is_player_at_gate(elapsed_time_ms) {
		// TODO: Restrict this based on the previous height
		return get_first_gate()
	} else {
		return prev_upper, prev_lower
	}
}

func get_next_player_position(prev_position int, time_since_last_jump_ms int64) int {
	var basic_current_position = get_basic_position(time_since_last_jump_ms)

	var player_position = 0
	if time_since_last_jump_ms == 0 {
		player_position = prev_position + basic_current_position
	} else {
		var offset = prev_position - get_basic_position(time_since_last_jump_ms-1)
		player_position = basic_current_position + offset
	}

	if player_position > MAX_HEIGHT {
		player_position = MAX_HEIGHT
	} else if player_position < MIN_HEIGHT {
		player_position = MIN_HEIGHT
	}

	return player_position
}

func get_basic_position(time_since_last_jump_ms int64) int {
	const JUMP_HEIGHT = 200
	return int(JUMP_HEIGHT - (time_since_last_jump_ms / 10))
}

func is_player_at_gate(elapsed_time_ms int64) bool {
	if elapsed_time_ms == 0 {
		return false
	} else {
		return elapsed_time_ms%GATE_DISTANCE_MS == 0
	}
}

// Game Status Methods

func get_distance_to_next_gate(elapsed_time_ms int64) int {
	return int(GATE_DISTANCE_MS - (elapsed_time_ms % GATE_DISTANCE_MS))
}

func is_player_alive(gate_upper int, gate_lower int, player_position int, elapsed_time_ms int64) bool {
	if is_player_at_gate(elapsed_time_ms) {
		return player_position > gate_lower && player_position < gate_upper
	} else {
		return true
	}
}
