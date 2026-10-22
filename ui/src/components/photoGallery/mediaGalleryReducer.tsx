import React, { useEffect } from 'react'
import { MediaGalleryFields } from './__generated__/MediaGalleryFields'
import { act } from 'react-dom/test-utils'

export interface MediaGalleryState {
  presenting: boolean
  activeIndex: number
  media: MediaGalleryFields[]
}

export type GalleryAction =
  | { type: 'nextImage' }
  | { type: 'previousImage' }
  | { type: 'closePresentMode' }

export type PhotoGalleryAction =
  | GalleryAction
  | { type: 'openPresentMode'; activeIndex: number }
  | { type: 'selectImage'; index: number }
  | { type: 'replaceMedia'; media: MediaGalleryFields[] }

const arraysEqual = (a: string[], b: string[]) => {
  if (a.length !== b.length) return false
  for (let i = 0; i < a.length; i++) {
    if (a[i] !== b[i]) return false
  }
  return true
}

export function mediaGalleryReducer(
  state: MediaGalleryState,
  action: PhotoGalleryAction
): MediaGalleryState {
  switch (action.type) {
    case 'nextImage':
      return {
        ...state,
        activeIndex: (state.activeIndex + 1) % state.media.length,
      }
    case 'previousImage':
      if (state.activeIndex <= 0) {
        return {
          ...state,
          activeIndex: state.media.length - 1,
        }
      } else {
        return {
          ...state,
          activeIndex: state.activeIndex - 1,
        }
      }
    case 'openPresentMode':
      return {
        ...state,
        presenting: true,
        activeIndex: action.activeIndex,
      }
    case 'closePresentMode':
      return {
        ...state,
        presenting: false,
      }
    case 'selectImage':
      return {
        ...state,
        activeIndex: Math.max(
          0,
          Math.min(state.media.length - 1, action.index)
        ),
      }
    case 'replaceMedia':
      const existingMediaList = state.media.map(m => m.id)
      const updatedMediaList = action.media.map(m => m.id)
      const mediaChanged = !arraysEqual(existingMediaList, updatedMediaList);
      const updatedActiveIndex = (mediaChanged ? -1 : state.activeIndex);
      const updatedPresenting = (mediaChanged ? false : state.presenting);
      return {
        ...state,
        media: action.media,
        activeIndex: updatedActiveIndex,
        presenting: updatedPresenting,
      }
  }
}

export interface MediaGalleryPopStateEvent extends PopStateEvent {
  state: MediaGalleryState
}

export const urlPresentModeSetupHook = ({
  dispatchMedia,
  openPresentMode,
}: {
  dispatchMedia: React.Dispatch<GalleryAction>
  openPresentMode: (event: MediaGalleryPopStateEvent) => void
}) => {
  useEffect(() => {
    const urlChangeListener = (event: MediaGalleryPopStateEvent) => {
      if (event.state.presenting === true) {
        openPresentMode(event)
      } else {
        dispatchMedia({ type: 'closePresentMode' })
      }
    }

    window.addEventListener('popstate', urlChangeListener)

    history.replaceState({ presenting: false }, '')

    return () => {
      window.removeEventListener('popstate', urlChangeListener)
    }
  }, [])
}

export const openPresentModeAction = ({
  dispatchMedia,
  activeIndex,
}: {
  dispatchMedia: React.Dispatch<PhotoGalleryAction>
  activeIndex: number
}) => {
  dispatchMedia({
    type: 'openPresentMode',
    activeIndex: activeIndex,
  })

  history.pushState({ presenting: true, activeIndex }, '')
}

export const closePresentModeAction = ({
  dispatchMedia,
}: {
  dispatchMedia: React.Dispatch<GalleryAction>
}) => {
  dispatchMedia({
    type: 'closePresentMode',
  })

  history.back()
}
