/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

import {
  GuessDifficulty,
  GuessGameFilter,
} from './../../../__generated__/globalTypes'

// ====================================================
// GraphQL mutation operation: startGuessGame
// ====================================================

export interface startGuessGame_startGuessGame_rounds_photo_album {
  __typename: 'Album'
  id: string
}

export interface startGuessGame_startGuessGame_rounds_photo_thumbnail {
  __typename: 'MediaURL'
  /**
   * URL for previewing the image
   */
  url: string
  /**
   * Width of the image in pixels
   */
  width: number
  /**
   * Height of the image in pixels
   */
  height: number
}

export interface startGuessGame_startGuessGame_rounds_photo_highRes {
  __typename: 'MediaURL'
  /**
   * URL for previewing the image
   */
  url: string
  /**
   * Width of the image in pixels
   */
  width: number
  /**
   * Height of the image in pixels
   */
  height: number
}

export interface startGuessGame_startGuessGame_rounds_photo {
  __typename: 'Media'
  id: string
  /**
   * A short string that can be used to generate a blured version of the media, to show while the original is loading
   */
  blurhash: string | null
  /**
   * The album that holds the media
   */
  album: startGuessGame_startGuessGame_rounds_photo_album
  /**
   * URL to display the media in a smaller resolution
   */
  thumbnail: startGuessGame_startGuessGame_rounds_photo_thumbnail | null
  /**
   * URL to display the photo in full resolution, will be null for videos
   */
  highRes: startGuessGame_startGuessGame_rounds_photo_highRes | null
}

export interface startGuessGame_startGuessGame_rounds_result_guess {
  __typename: 'Coordinates'
  /**
   * GPS latitude in degrees
   */
  latitude: number
  /**
   * GPS longitude in degrees
   */
  longitude: number
}

export interface startGuessGame_startGuessGame_rounds_result_answer {
  __typename: 'Coordinates'
  /**
   * GPS latitude in degrees
   */
  latitude: number
  /**
   * GPS longitude in degrees
   */
  longitude: number
}

export interface startGuessGame_startGuessGame_rounds_result {
  __typename: 'GuessRoundResult'
  guess: startGuessGame_startGuessGame_rounds_result_guess
  guessYear: number
  guessMonth: number
  answer: startGuessGame_startGuessGame_rounds_result_answer
  answerYear: number
  answerMonth: number
  distanceKm: number
  monthsOff: number
  locationScore: number
  timeScore: number
  score: number
}

export interface startGuessGame_startGuessGame_rounds {
  __typename: 'GuessRound'
  /**
   * Random id of the round
   */
  id: string
  /**
   * Position of the round in its game, starting at 0
   */
  index: number
  /**
   * The photo to guess, null if it was deleted. Do not show its title or path
   */
  photo: startGuessGame_startGuessGame_rounds_photo | null
  /**
   * The guess and its score, null until the round is guessed
   */
  result: startGuessGame_startGuessGame_rounds_result | null
}

export interface startGuessGame_startGuessGame {
  __typename: 'GuessGame'
  /**
   * Random id of the game; games are only reachable by it
   */
  id: string
  createdAt: Time
  finishedAt: Time | null
  /**
   * True when the game was replaced by a newer game on its device, or expired
   */
  abandoned: boolean
  difficulty: GuessDifficulty
  totalScore: number
  maxScore: number
  rounds: startGuessGame_startGuessGame_rounds[]
}

export interface startGuessGame {
  /**
   * Start a game. deviceKey is a random id the browser keeps, used to avoid repeating photos.
   * replaceGameId names the device's previous unfinished game, which is then abandoned.
   */
  startGuessGame: startGuessGame_startGuessGame
}

export interface startGuessGameVariables {
  deviceKey: string
  rounds?: number | null
  difficulty?: GuessDifficulty | null
  filter?: GuessGameFilter | null
  replaceGameId?: string | null
}
