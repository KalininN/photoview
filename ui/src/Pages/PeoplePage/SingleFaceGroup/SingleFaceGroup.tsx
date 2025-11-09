import { gql, useQuery } from '@apollo/client'
import { useIsAdmin } from '../../../components/routes/AuthorizedRoute'
import React, { useEffect, useReducer } from 'react'
import { useTranslation } from 'react-i18next'
import useURLParameters from '../../../hooks/useURLParameters'
import PaginateLoader from '../../../components/PaginateLoader'
import MediaGallery from '../../../components/photoGallery/MediaGallery'
import { mediaGalleryReducer } from '../../../components/photoGallery/mediaGalleryReducer'
import useScrollPagination from '../../../hooks/useScrollPagination'
import FaceGroupTitle from './FaceGroupTitle'
import {
  singleFaceGroup,
  singleFaceGroupVariables,
  singleFaceGroup_faceGroup_imageFaces,
} from './__generated__/singleFaceGroup'
import { set } from 'husky'
import { FaceGroupsWrapper } from '../PeoplePage'
import { FaceCircleImage } from '../FaceCircleImage'


export const SINGLE_FACE_GROUP = gql`
  query singleFaceGroup($id: ID!, $minRating: Int, $limit: Int!, $offset: Int!) {
    faceGroup(id: $id) {
      id
      displayLabel
      imageFaces(minRating: $minRating, paginate: { limit: $limit, offset: $offset }) {
        id
        rectangle {
          minX
          maxX
          minY
          maxY
        }
        media {
          id
          type
          title
          blurhash
          thumbnail {
            url
            width
            height
          }
          highRes {
            url
          }
          favorite
        }
      }
    }
  }
`

type SingleFaceGroupProps = {
  faceGroupID: string
}

type FaceDetectedProps = {
  face: singleFaceGroup_faceGroup_imageFaces
}

export const FaceDetected = ({ face }: FaceDetectedProps) => {
  const previewFace = face

  return (
    <div className="m-3">
      <FaceCircleImage imageFace={previewFace} selectable />
    </div>
  )
}

const SingleFaceGroup = ({ faceGroupID }: SingleFaceGroupProps) => {
  const { t } = useTranslation()

  const { getParam, setParam } = useURLParameters()

  const sanitizeRating = (rating: number) => {
    return Math.max(Math.min(rating, 5), 3)
  }

  const minRatingParam = getParam('minRating')
  const minRating = sanitizeRating(minRatingParam !== null ? parseInt(minRatingParam, 10) : 5)
  const setMinRating = (rating: number) => {
    setParam('minRating', sanitizeRating(rating).toString())
  }

  const { data, error, loading, refetch, fetchMore } = useQuery<
    singleFaceGroup,
    singleFaceGroupVariables
  >(SINGLE_FACE_GROUP, {
    variables: {
      limit: 200,
      offset: 0,
      id: faceGroupID,
      minRating: minRating,
    },
  })

  const [mediaState, dispatchMedia] = useReducer(mediaGalleryReducer, {
    presenting: false,
    activeIndex: -1,
    media: [],
  })

  const { containerElem, finished: finishedLoadingMore } =
    useScrollPagination<singleFaceGroup>({
      loading,
      fetchMore,
      data,
      getItems: data => data.faceGroup.imageFaces,
    })

  useEffect(() => {
    const media = data?.faceGroup.imageFaces.map(x => x.media) || []
    dispatchMedia({ type: 'replaceMedia', media })
  }, [data])

  useEffect(() => {
    refetch({
      minRating: minRating,
    })
  }, [minRating])

  const faceGroup = data?.faceGroup

  const isAdmin = useIsAdmin()
  let faces: JSX.Element[] | null = null
  if (data) {
    faces = faceGroup?.imageFaces.map(imageFace => (
      <FaceDetected face={imageFace} />
    ))
  }

  if (error) {
    return <div>{error.message}</div>
  }

  return (
    <div ref={containerElem}>
      <FaceGroupTitle
        faceGroup={faceGroup}
        minRating={minRating}
        setMinRating={setMinRating}
      />
      <div>
        <MediaGallery
          loading={loading}
          dispatchMedia={dispatchMedia}
          mediaState={mediaState}
        />
        {isAdmin && (
          <FaceGroupsWrapper ref={containerElem}>{faces}</FaceGroupsWrapper>
        )}
        <PaginateLoader
          active={!finishedLoadingMore && !loading}
          text={t('general.loading.paginate.media', 'Loading more media')}
        />
      </div>
    </div>
  )
}

export default SingleFaceGroup
