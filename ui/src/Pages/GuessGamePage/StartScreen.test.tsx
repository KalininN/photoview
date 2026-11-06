import React from 'react'
import { MockedProvider } from '@apollo/client/testing'
import { fireEvent, render, screen } from '@testing-library/react'
import StartScreen from './StartScreen'
import { GUESS_GAME_INFO_QUERY } from './guessGameQueries'
import { testGame } from './testFixtures'
import { GuessDifficulty } from '../../__generated__/globalTypes'

const infoMock = (eligibleCount: number) => ({
  request: {
    query: GUESS_GAME_INFO_QUERY,
    variables: { filter: { minRating: 4 } },
  },
  result: {
    data: {
      guessGameInfo: {
        __typename: 'GuessGameInfo',
        eligibleCount,
        minYear: 2005,
        maxYear: 2025,
        southWest: null,
        northEast: null,
      },
    },
  },
})

const renderStart = (
  lastGame: ReturnType<typeof testGame> | null,
  eligibleCount = 120
) => {
  const onContinue = vi.fn()
  const onStart = vi.fn()
  render(
    <MockedProvider mocks={[infoMock(eligibleCount)]}>
      <StartScreen
        lastGame={lastGame}
        onContinue={onContinue}
        onStart={onStart}
        starting={false}
      />
    </MockedProvider>
  )
  return { onContinue, onStart }
}

test('an unfinished last game can be continued', async () => {
  const { onContinue } = renderStart(testGame(2))
  fireEvent.click(
    screen.getByRole('button', { name: 'Continue (round 3 of 5)' })
  )
  expect(onContinue).toHaveBeenCalled()
  expect(
    await screen.findByText('120 photos with a place and a date')
  ).toBeInTheDocument()
  expect(
    screen.getByRole('button', { name: 'Start a new game instead' })
  ).toBeEnabled()
})

test('a finished or replaced game cannot be continued', () => {
  renderStart(testGame(5))
  expect(screen.queryByRole('button', { name: /Continue/ })).toBeNull()
  expect(
    screen.getByText('Your last game here: 15,000 of 50,000 points.')
  ).toBeInTheDocument()
})

test('a replaced game is not shown at all', () => {
  renderStart(testGame(2, 5, { abandoned: true }))
  expect(screen.queryByRole('button', { name: /Continue/ })).toBeNull()
  expect(screen.queryByText(/Your last game/)).toBeNull()
})

test('a new game starts with the defaults: 5 rounds, auto, rated 4+', async () => {
  const { onStart } = renderStart(null)
  const start = screen.getByRole('button', { name: 'Start' })
  await screen.findByText('120 photos with a place and a date')
  fireEvent.click(start)
  expect(onStart).toHaveBeenCalledWith({
    rounds: 5,
    difficulty: GuessDifficulty.AUTO,
    minRating: 4,
  })
})

test('a small pool shortens the game, an empty one cannot start', async () => {
  const { onStart } = renderStart(null, 3)
  await screen.findByText(/Only 3 photos can be used/)
  fireEvent.click(screen.getByRole('button', { name: 'Start' }))
  expect(onStart).toHaveBeenCalledWith(expect.objectContaining({ rounds: 3 }))
})

test('an empty pool explains why', async () => {
  renderStart(null, 0)
  await screen.findByText(/No photos can be used/)
  expect(screen.getByRole('button', { name: 'Start' })).toBeDisabled()
})
