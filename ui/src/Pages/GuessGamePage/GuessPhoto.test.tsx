import React from 'react'
import { fireEvent, render, screen } from '@testing-library/react'
import GuessPhoto from './GuessPhoto'
import { testRound } from './testFixtures'

test('the photo has a fixed alt text and no title', () => {
  const round = testRound(0, false)
  const { container } = render(
    <GuessPhoto photo={round.photo} preload={testRound(1, false).photo} />
  )

  const image = screen.getByAltText('Photo to guess')
  expect(image.getAttribute('src')).toContain('/photo/highres_0.jpg')
  expect(container.querySelector('[title]')).toBeNull()

  // The next photo is loaded in the background, hidden
  expect(screen.getByTestId('guess-photo-preload')).toHaveClass('hidden')
})

test('clicking the photo opens it full screen', () => {
  render(<GuessPhoto photo={testRound(0, false).photo} />)
  fireEvent.click(
    screen.getByRole('button', { name: 'Show the photo full screen' })
  )
  expect(screen.getByRole('dialog')).toBeInTheDocument()
  fireEvent.keyDown(document, { key: 'Escape' })
  expect(screen.queryByRole('dialog')).toBeNull()
})

test('a deleted photo still lets the player guess', () => {
  render(<GuessPhoto photo={null} />)
  expect(screen.getByText(/This photo was deleted/)).toBeInTheDocument()
})
