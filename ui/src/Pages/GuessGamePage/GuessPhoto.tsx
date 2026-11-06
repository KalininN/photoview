// The photo of a 4D-guesser round. It never shows the photo's title or file name,
// not even as alt text. Fork-specific: not part of upstream Photoview.
import React, { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { BlurhashCanvas } from 'react-blurhash'
import { ProtectedImage } from '../../components/photoGallery/ProtectedMedia'
import { GuessGameFields_rounds_photo } from './__generated__/GuessGameFields'

type GuessPhotoProps = {
  photo: GuessGameFields_rounds_photo | null
  /** The next round's photo, loaded in the background */
  preload?: GuessGameFields_rounds_photo | null
}

const photoUrl = (photo: GuessGameFields_rounds_photo) =>
  photo.highRes?.url ?? photo.thumbnail?.url

const GuessPhoto = ({ photo, preload }: GuessPhotoProps) => {
  const { t } = useTranslation()
  const [loaded, setLoaded] = useState(false)
  const [zoomed, setZoomed] = useState(false)

  useEffect(() => {
    setLoaded(false)
    setZoomed(false)
  }, [photo?.id])

  useEffect(() => {
    if (!zoomed) return
    const onKey = (event: KeyboardEvent) => {
      if (event.key === 'Escape') setZoomed(false)
    }
    document.addEventListener('keydown', onKey)
    return () => document.removeEventListener('keydown', onKey)
  }, [zoomed])

  const alt = t('guess_game.photo_alt', 'Photo to guess')
  const preloadUrl = preload ? photoUrl(preload) : undefined

  if (photo === null) {
    return (
      <div className="w-full h-full flex items-center justify-center text-center p-6 bg-gray-100 dark:bg-dark-bg2 rounded">
        {t(
          'guess_game.photo_deleted',
          'This photo was deleted. The round still counts: guess anyway.'
        )}
      </div>
    )
  }

  return (
    <div className="relative w-full h-full bg-black rounded overflow-hidden">
      {photo.blurhash && !loaded && (
        <BlurhashCanvas
          className="absolute w-full h-full top-0 left-0"
          hash={photo.blurhash}
        />
      )}
      <button
        type="button"
        className="absolute inset-0 w-full h-full cursor-zoom-in"
        onClick={() => setZoomed(true)}
        aria-label={t('guess_game.zoom_photo', 'Show the photo full screen')}
      >
        <ProtectedImage
          key={photo.id}
          src={photoUrl(photo)}
          alt={alt}
          className="w-full h-full object-contain"
          onLoad={() => setLoaded(true)}
        />
      </button>
      {zoomed && (
        <div
          className="fixed inset-0 z-50 bg-black flex items-center justify-center cursor-zoom-out"
          onClick={() => setZoomed(false)}
          role="dialog"
          aria-label={alt}
        >
          <ProtectedImage
            src={photoUrl(photo)}
            alt={alt}
            className="max-w-full max-h-full object-contain"
          />
        </div>
      )}
      {preloadUrl && (
        <ProtectedImage
          src={preloadUrl}
          alt=""
          aria-hidden
          className="hidden"
          data-testid="guess-photo-preload"
        />
      )}
    </div>
  )
}

export default GuessPhoto
