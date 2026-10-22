import { MediaGalleryFields } from './__generated__/MediaGalleryFields'
import { gql, MutationFunction, useMutation } from '@apollo/client'
import {
  markMediaFavorite,
  markMediaFavoriteVariables,
} from './__generated__/markMediaFavorite'
import { likeMedia, likeMediaVariables } from './__generated__/likeMedia'

const markFavoriteMutation = gql`
  mutation markMediaFavorite($mediaId: ID!, $favorite: Boolean!) {
    favoriteMedia(mediaId: $mediaId, favorite: $favorite) {
      id
      favorite
    }
  }
`

export const useMarkFavoriteMutation = () => {
  return useMutation<markMediaFavorite, markMediaFavoriteVariables>(
    markFavoriteMutation
  )
}

export const toggleFavoriteAction = ({
  media,
  markFavorite,
}: {
  media: MediaGalleryFields
  markFavorite: MutationFunction<markMediaFavorite, markMediaFavoriteVariables>
}) => {
  return markFavorite({
    variables: {
      mediaId: media.id,
      favorite: !media.favorite,
    },
    optimisticResponse: {
      favoriteMedia: {
        id: media.id,
        favorite: !media.favorite,
        __typename: 'Media',
      },
    },
  })
}

const LIKE_MEDIA_MUTATION = gql`
  mutation likeMedia($mediaId: ID!) {
    likeMedia(mediaId: $mediaId) {
      id
      likes
    }
  }
`

const getLocalStorageLikedKey = (mediaId: string) => {
  return `liked_media_${mediaId}`
}

export const isMediaLiked = (mediaId: string): boolean => {
  return localStorage.getItem(getLocalStorageLikedKey(mediaId)) === 'true'
}

export const useLikeMediaMutation = () => {
  return useMutation<likeMedia, likeMediaVariables>(LIKE_MEDIA_MUTATION)
}

export const likePressed = ({ media, likeMedia }:
  { media: MediaGalleryFields,
    likeMedia: MutationFunction<likeMedia, likeMediaVariables>
  }) => {
  const wasLiked = isMediaLiked(media.id)
  if (wasLiked) {
    return
  }
  localStorage.setItem(getLocalStorageLikedKey(media.id), 'true')
  return likeMedia({
    variables: {
      mediaId: media.id,
    },
    optimisticResponse: {
      likeMedia: {
        __typename: 'Media',
        id: media.id,
        likes: (media.likes || 0) + 1,
      },
    },
  })
}
