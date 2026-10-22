import React from 'react'
import styled from 'styled-components'
import { ProtectedImage } from './ProtectedMedia'
import { MediaType } from '../../__generated__/globalTypes'
import { ReactComponent as VideoThumbnailIconSVG } from './icons/videoThumbnailIcon.svg'
import { MediaGalleryFields } from './__generated__/MediaGalleryFields'
import { isMediaLiked } from './photoGalleryMutations'
import { HoverIcon, FavoriteIcon, LikesIcon, SidebarIcon } from '../media/HoverButtons'

const MediaContainer = styled.div`
  flex-grow: 1;
  flex-basis: 0;
  height: 200px;
  margin: 4px 2px;
  background-color: #eee;
  position: relative;
  overflow: hidden;
`

const StyledPhoto = styled(ProtectedImage)`
  height: 200px;
  min-width: 100%;
  position: relative;
  object-fit: cover;

  transition: opacity 300ms;
`

type LazyPhotoProps = {
  src?: string
  blurhash: string | null
}

const LazyPhoto = (photoProps: LazyPhotoProps) => {
  return <StyledPhoto {...photoProps} lazyLoading />
}

const PhotoOverlay = styled.div<{ active: boolean }>`
  width: 100%;
  height: 100%;
  position: absolute;
  top: 0;
  left: 0;

  ${({ active }) =>
    active &&
    `
      outline: 4px solid rgba(65, 131, 196, 0.6);
      outline-offset: -4px;
    `}
`

const LikesButtonWrapper = styled(HoverIcon)`
  margin-top: 160px !important;
  ${MediaContainer}:hover &, ${MediaContainer}:focus-within & {
    opacity: 1 !important;
  }
`

type LikesButtonProps = {
  liked: boolean
  likes: number
  onClick(e: React.MouseEvent<HTMLButtonElement, MouseEvent>): void
}

const LikesButton = ({ liked, likes, onClick }: LikesButtonProps) => {
  return (
    <LikesButtonWrapper
      onClick={onClick}
      style={{ opacity: liked ? '0.75' : undefined }}
    >
      <LikesIcon liked={liked} likes={likes} />
    </LikesButtonWrapper>
  )
}

type FavoriteButtonProps = {
  favorite: boolean
  onClick(e: React.MouseEvent<HTMLButtonElement, MouseEvent>): void
}

const FavoriteButton = ({ favorite, onClick }: FavoriteButtonProps) => {
  return (
    <LikesButtonWrapper
      onClick={onClick}
      style={{ opacity: favorite ? '0.75' : undefined }}
    >
      <FavoriteIcon favorite={favorite} />
    </LikesButtonWrapper>
  )
}

const SidebarIconWrapper = styled(HoverIcon)`
  margin: 10px !important;
  position: absolute;
  top: 0;
  right: 0;
  ${MediaContainer}:hover &, ${MediaContainer}:focus-within & {
    opacity: 1 !important;
  }
`

type SidebarIconProps = {
  onClick(e: React.MouseEvent<HTMLButtonElement, MouseEvent>): void
}

const SidebarButton = ({ onClick }: SidebarIconProps) => (
  <SidebarIconWrapper onClick={onClick}>
    <SidebarIcon/>
  </SidebarIconWrapper>
)

const VideoThumbnailIcon = styled(VideoThumbnailIconSVG)`
  color: rgba(255, 255, 255, 0.8);
  position: absolute;
  left: calc(50% - 17.5px);
  top: calc(50% - 22px);
`

type MediaThumbnailProps = {
  media: MediaGalleryFields
  active: boolean
  selectImage(): void
  clickPresent(): void
  clickFavorite(): void
  clickLike(): void
}

export const MediaThumbnail = ({
  media,
  active,
  selectImage,
  clickPresent,
  clickFavorite,
  clickLike
}: MediaThumbnailProps) => {
  let heartIcon = null
  if (media.favorite !== undefined) {
    heartIcon = (
      <FavoriteButton
        favorite={media.favorite}
        onClick={e => {
          e.stopPropagation()
          clickFavorite()
        }}
      />
    )
  }
  if (media.likes !== undefined) {
    heartIcon = (
      <LikesButton
        liked={isMediaLiked(media.id)}
        likes={media.likes? media.likes : 0}
        onClick={e => {
          e.stopPropagation()
          clickLike()
        }}
      />
    )
  }

  let videoIcon = null
  if (media.type == MediaType.Video) {
    videoIcon = <VideoThumbnailIcon />
  }

  let minWidth = 100
  if (media.thumbnail) {
    minWidth = Math.floor(
      (media.thumbnail.width / media.thumbnail.height) * 200
    )
  }

  return (
    <MediaContainer
      key={media.id}
      style={{
        cursor: 'pointer',
        minWidth: `clamp(124px, ${minWidth}px, 100% - 8px)`,
      }}
      onClick={() => {
        clickPresent()
      }}
    >
      <div
        style={{
          // minWidth: `min(${minWidth}px, 100%)`,
          minWidth: `${minWidth}px`,
          height: `200px`,
        }}
      >
        <LazyPhoto src={media.thumbnail?.url} blurhash={media.blurhash} />
      </div>
      <PhotoOverlay active={active}>
        {videoIcon}
        <SidebarButton
          onClick={e => {
            e.stopPropagation()
            selectImage()
          }}
        />
        {heartIcon}
      </PhotoOverlay>
    </MediaContainer>
  )
}

export const MediaPlaceholder = styled.div`
  flex-grow: 1;
  height: 200px;
  width: 300px;
  margin: 4px 2px;
  background-color: #eee;
  position: relative;
`
