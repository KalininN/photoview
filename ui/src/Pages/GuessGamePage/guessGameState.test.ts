import {
  formatDistance,
  formatMonthYear,
  gameView,
  guessedCount,
  isGameOpen,
  nextRound,
  shortMonthNames,
  unwrapLongitude,
  wrapLongitude,
} from './guessGameState'
import { testGame } from './testFixtures'

describe('gameView', () => {
  test('a new game starts at the first round', () => {
    const view = gameView(testGame(0), null)
    expect(view.kind).toBe('guess')
    expect(view.kind === 'guess' && view.round.index).toBe(0)
  })

  test('a reload mid-game resumes at the first unguessed round', () => {
    const view = gameView(testGame(2), null)
    expect(view.kind === 'guess' && view.round.index).toBe(2)
  })

  test('the round just guessed stays revealed until the player moves on', () => {
    const view = gameView(testGame(3), 'round-2')
    expect(view).toMatchObject({
      kind: 'reveal',
      round: { index: 2 },
      hasNext: true,
    })
  })

  test('the last reveal leads to the summary', () => {
    const game = testGame(5)
    expect(gameView(game, 'round-4')).toMatchObject({
      kind: 'reveal',
      hasNext: false,
    })
    expect(gameView(game, null)).toEqual({ kind: 'summary' })
  })

  test('an unguessed or unknown reveal id is ignored', () => {
    expect(gameView(testGame(1), 'round-3').kind).toBe('guess')
    expect(gameView(testGame(1), 'nope').kind).toBe('guess')
  })

  test('rounds are ordered by index, not by array position', () => {
    const game = testGame(1)
    game.rounds.reverse()
    expect(nextRound(game)?.index).toBe(1)
  })
})

test('isGameOpen', () => {
  expect(isGameOpen(testGame(0))).toBe(true)
  expect(isGameOpen(testGame(4))).toBe(true)
  expect(isGameOpen(testGame(5))).toBe(false)
  expect(isGameOpen(testGame(2, 5, { abandoned: true }))).toBe(false)
  expect(guessedCount(testGame(3))).toBe(3)
})

test('longitudes', () => {
  expect(wrapLongitude(190)).toBeCloseTo(-170)
  expect(wrapLongitude(-190)).toBeCloseTo(170)
  expect(wrapLongitude(540)).toBe(180)
  expect(wrapLongitude(12.5)).toBe(12.5)
  // The line from 179 to -179 goes over the antimeridian, not around the world
  expect(unwrapLongitude(179, -179)).toBe(181)
  expect(unwrapLongitude(-179, 179)).toBe(-181)
  expect(unwrapLongitude(10, 20)).toBe(20)
})

test('formatting', () => {
  expect(formatDistance(0.0254, 'en')).toBe('25 m')
  expect(formatDistance(3.456, 'en')).toBe('3.5 km')
  expect(formatDistance(1234.4, 'en')).toBe('1,234 km')
  expect(formatMonthYear(2019, 3, 'en')).toBe('March 2019')
  expect(formatMonthYear(2019, 12, 'not a locale!')).toContain('2019')
  expect(shortMonthNames('en')).toHaveLength(12)
  expect(shortMonthNames('en')[0]).toBe('Jan')
})
