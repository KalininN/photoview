/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

import { GuessInput } from './../../../__generated__/globalTypes'

// ====================================================
// GraphQL mutation operation: submitGuess
// ====================================================

export interface submitGuess_submitGuess {
  __typename: 'GuessRoundResult'
  score: number
}

export interface submitGuess {
  /**
   * Guess a round. Guessing an already guessed round returns its stored result
   */
  submitGuess: submitGuess_submitGuess
}

export interface submitGuessVariables {
  roundId: string
  guess: GuessInput
}
