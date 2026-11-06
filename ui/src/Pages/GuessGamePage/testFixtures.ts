// Test data for the 4D-guesser tests. Fork-specific: not part of upstream Photoview.
import { GuessDifficulty } from '../../__generated__/globalTypes'
import {
  GuessGameFields,
  GuessGameFields_rounds,
} from './__generated__/GuessGameFields'

export function testRound(
  index: number,
  guessed: boolean,
  overrides: Partial<GuessGameFields_rounds> = {}
): GuessGameFields_rounds {
  return {
    __typename: 'GuessRound',
    id: `round-${index}`,
    index,
    photo: {
      __typename: 'Media',
      id: `${100 + index}`,
      blurhash: null,
      album: { __typename: 'Album', id: '7' },
      thumbnail: {
        __typename: 'MediaURL',
        url: `/photo/thumbnail_${index}.jpg`,
        width: 1024,
        height: 768,
      },
      highRes: {
        __typename: 'MediaURL',
        url: `/photo/highres_${index}.jpg`,
        width: 4000,
        height: 3000,
      },
    },
    result: guessed
      ? {
          __typename: 'GuessRoundResult',
          guess: {
            __typename: 'Coordinates',
            latitude: 48.85,
            longitude: 2.35,
          },
          guessYear: 2019,
          guessMonth: 3,
          answer: {
            __typename: 'Coordinates',
            latitude: 51.5,
            longitude: -0.12,
          },
          answerYear: 2020,
          answerMonth: 1,
          distanceKm: 343.6,
          monthsOff: 10,
          locationScore: 1000,
          timeScore: 2000,
          score: 3000,
        }
      : null,
    ...overrides,
  }
}

export function testGame(
  guessed: number,
  total = 5,
  overrides: Partial<GuessGameFields> = {}
): GuessGameFields {
  return {
    __typename: 'GuessGame',
    id: 'game-1',
    createdAt: '2026-11-06T12:00:00Z',
    finishedAt: guessed === total ? '2026-11-06T12:30:00Z' : null,
    abandoned: false,
    difficulty: GuessDifficulty.AUTO,
    totalScore: guessed * 3000,
    maxScore: total * 10000,
    rounds: Array.from({ length: total }, (_, i) => testRound(i, i < guessed)),
    ...overrides,
  }
}
