import React from 'react'
import PresentNavigationOverlay from './PresentNavigationOverlay'
import { fireEvent, render, screen, act } from '@testing-library/react'
import { MockedProvider } from '@apollo/client/testing'
import { MediaType } from '../../../__generated__/globalTypes'
import { MediaGalleryFields } from '../__generated__/MediaGalleryFields'

// Fork-specific: the overlay shows a like button for the media
const media: MediaGalleryFields = {
  __typename: 'Media',
  id: '1',
  type: MediaType.Photo,
  blurhash: null,
  thumbnail: null,
  highRes: null,
  videoWeb: null,
  favorite: false,
  likes: 0,
}

const renderOverlay = (dispatchMedia: () => void) =>
  render(
    <MockedProvider>
      <PresentNavigationOverlay dispatchMedia={dispatchMedia} media={media} />
    </MockedProvider>
  )

vi.useFakeTimers()

describe('PresentNavigationOverlay component', () => {
  test('simple render', () => {
    const dispatchMedia = vi.fn()
    renderOverlay(dispatchMedia)

    expect(screen.getByLabelText('Previous image')).toBeInTheDocument()
    expect(screen.getByLabelText('Next image')).toBeInTheDocument()
    expect(screen.getByLabelText('Exit presentation mode')).toBeInTheDocument()
  })

  test('click buttons', () => {
    const dispatchMedia = vi.fn()
    renderOverlay(dispatchMedia)

    expect(dispatchMedia).not.toHaveBeenCalled()

    fireEvent.click(screen.getByLabelText('Next image'))
    expect(dispatchMedia).lastCalledWith({ type: 'nextImage' })

    fireEvent.click(screen.getByLabelText('Previous image'))
    expect(dispatchMedia).lastCalledWith({ type: 'previousImage' })
  })

  test('mouse move, show and hide', () => {
    const dispatchMedia = vi.fn()
    const { container } = renderOverlay(dispatchMedia)

    expect(screen.getByLabelText('Next image')).toHaveClass('hide')

    fireEvent.mouseMove(container.firstChild!)
    expect(screen.getByLabelText('Next image')).not.toHaveClass('hide')

    act(() => {
      vi.advanceTimersByTime(3000)
    })

    expect(screen.getByLabelText('Next image')).toHaveClass('hide')
  })
})
