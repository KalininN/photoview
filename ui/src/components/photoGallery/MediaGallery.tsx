import React, { useContext } from 'react'
import styled from 'styled-components'
import { MediaThumbnail, MediaPlaceholder } from './MediaThumbnail'
import PresentView from './presentView/PresentView'
import {
  openPresentModeAction,
  PhotoGalleryAction,
  MediaGalleryState,
} from './mediaGalleryReducer'
import { MediaGalleryFields } from './__generated__/MediaGalleryFields'
import {
  toggleFavoriteAction,
  useMarkFavoriteMutation,
} from './photoGalleryMutations'
import MediaSidebar from '../sidebar/MediaSidebar/MediaSidebar'
import { SidebarContext } from '../sidebar/Sidebar'
import { useMutation, gql } from '@apollo/client'
import { likeMedia, likeMediaVariables } from './__generated__/likeMedia'

const Gallery = styled.div`
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  min-height: 200px;
  position: relative;
  margin: -4px;

  @media (max-width: 1000px) {
    /* Compensate for tab bar on mobile */
    margin-bottom: 76px;
  }
`

export const PhotoFiller = styled.div`
  height: 200px;
  flex-grow: 999999;
`

export const MEDIA_GALLERY_FRAGMENT = gql`
  fragment MediaGalleryFields on Media {
    id
    type
    blurhash
    thumbnail {
      url
      width
      height
    }
    highRes {
      url
    }
    videoWeb {
      url
    }
    favorite
    likes
  }
`

const LIKE_MEDIA_MUTATION = gql`
  mutation likeMedia($mediaId: ID!) {
    likeMedia(mediaId: $mediaId) {
      id
      likes
    }
  }
`

export const getLocalStorageLikedKey = (mediaId: string) => {
  return `liked_media_${mediaId}`
}

type MediaGalleryProps = {
  loading: boolean
  mediaState: MediaGalleryState
  dispatchMedia: React.Dispatch<PhotoGalleryAction>
}

const MediaGallery = ({ mediaState, dispatchMedia }: MediaGalleryProps) => {
  const [markFavorite] = useMarkFavoriteMutation()

  const { media, activeIndex, presenting } = mediaState

  const { updateSidebar } = useContext(SidebarContext)
  const [likeMediaMutation] = useMutation<likeMedia, likeMediaVariables>(
    LIKE_MEDIA_MUTATION
  )

  const likePressed = ({ media }: { media: MediaGalleryFields }) => {
    const wasLiked = localStorage.getItem(getLocalStorageLikedKey(media.id)) === 'true'
    if (wasLiked) {
      return
    }
    localStorage.setItem(getLocalStorageLikedKey(media.id), 'true')
    // Optimistic UI update
    const newLikes = (media.likes || 0) + 1
    likeMediaMutation({
      variables: {
        mediaId: media.id,
      },
      optimisticResponse: {
        likeMedia: {
          __typename: 'Media',
          id: media.id,
          likes: newLikes,
        },
      },
    })
  }

  let mediaElements = []
  if (media) {
    mediaElements = media.map((media, index) => {
      const active = activeIndex == index

      return (
        <MediaThumbnail
          key={media.id}
          media={media}
          active={active}
          selectImage={() => {
            dispatchMedia({
              type: 'selectImage',
              index,
            })
            updateSidebar(<MediaSidebar media={mediaState.media[index]} />)
          }}
          clickFavorite={() => {
            toggleFavoriteAction({
              media,
              markFavorite,
            })
          }}
          clickLike={() => {
            likePressed({
              media
            })
          }}
          clickPresent={() => {
            openPresentModeAction({ dispatchMedia, activeIndex: index })
          }}
        />
      )
    })
  } else {
    for (let i = 0; i < 6; i++) {
      mediaElements.push(<MediaPlaceholder key={i} />)
    }
  }

  return (
    <>
      <Gallery data-testid="photo-gallery-wrapper">
        {mediaElements}
        <PhotoFiller />
      </Gallery>
      {presenting && (
        <PresentView
          activeMedia={mediaState.media[mediaState.activeIndex]}
          dispatchMedia={dispatchMedia}
        />
      )}
    </>
  )
}

export default MediaGallery
