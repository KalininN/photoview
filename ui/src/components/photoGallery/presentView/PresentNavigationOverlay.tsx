import React, { useState, useRef, useEffect } from 'react'
import styled from 'styled-components'
import { debounce, DebouncedFn } from '../../../helpers/utils'
import { closePresentModeAction, GalleryAction } from '../mediaGalleryReducer'

import { useSwipeable } from 'react-swipeable'

import ExitIcon from './icons/Exit'
import InfoIcon from './icons/Info'
import NextIcon from './icons/Next'
import PrevIcon from './icons/Previous'
import { MediaGalleryFields } from '../__generated__/MediaGalleryFields'
import { isMediaLiked, useLikeMediaMutation, likePressed } from '../photoGalleryMutations'

const StyledOverlayContainer = styled.div`
  width: 100%;
  height: 100%;
  position: relative;
`

const OverlayButton = styled.button`
  width: 64px;
  height: 64px;
  background: none;
  border: none;
  outline: none;
  cursor: pointer;
  position: absolute;
  font-size: 32px;

  & svg {
    width: 32px;
    height: 32px;
    overflow: visible !important;
  }

  & svg path {
    stroke: rgba(255, 255, 255, 0.5);
    transition-property: stroke, filter;
    transition-duration: 140ms;
  }

  &:hover svg path {
    stroke: rgba(255, 255, 255, 1);
    filter: drop-shadow(0px 0px 2px rgba(0, 0, 0, 0.6));
  }

  /* outline: none above hides the ring for pointer users; keyboard focus
     still needs to be visible, all the more so now that a key press reveals
     the controls and the first stop is a button on top of a photo. The thin
     dark edge keeps the white ring legible over a bright one. */
  &:focus-visible {
    outline: 3px solid rgba(255, 255, 255, 0.95);
    outline-offset: -6px;
    box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.7);
  }

  &:focus-visible svg path {
    stroke: rgba(255, 255, 255, 1);
    filter: drop-shadow(0px 0px 2px rgba(0, 0, 0, 0.6));
  }

  &:hover svg path.disabled {
    stroke: rgba(255, 255, 255, 0.5);
    filter: none;
  }

  &.hide svg path {
    stroke: rgba(255, 255, 255, 0);
    transition: stroke 300ms;
  }

  /* An invisible button must not stay usable: on touch devices the first tap
     is meant to reveal the controls via the container's click handler, not to
     fire whatever button happens to sit under the finger - and a keyboard
     user must not be able to tab to a control they cannot see, which
     pointer-events alone would still allow. visibility also takes them out
     of the accessibility tree. The delay lets the stroke finish fading out
     first; revealing them again is instant, since no transition applies in
     that direction. */
  &.hide {
    pointer-events: none;
    visibility: hidden;
    transition: visibility 0s 300ms;
  }

  & svg path.filled {
    fill: rgba(255, 255, 255, 0.5);
    transition-property: stroke, filter, fill;
  }

  &:hover svg path.filled {
    fill: rgba(255, 255, 255, 1);
  }

  &:hover svg path.filled.disabled {
    fill: rgba(255, 255, 255, 0.5);
  }

  &.hide svg path.filled {
    fill: rgba(255, 255, 255, 0);
  }

  & span {
    color: rgba(255, 255, 255, 0.5);
    transition-property: color;
    transition-duration: 140ms;
  }

  &:hover span {
    color: rgba(255, 255, 255, 1);
  }

  &:hover span.disabled {
    color: rgba(255, 255, 255, 0.5);
  }

  &.hide span {
    color: rgba(255, 255, 255, 0);
    transition: color 300ms;
  }
`

const ExitButton = styled(OverlayButton)`
  left: 28px;
  top: 28px;
`

const InfoButton = styled(OverlayButton)`
  right: 28px;
  top: 28px;
`

const NavigationButton = styled(OverlayButton)<{ align: 'left' | 'right' }>`
  height: 80%;
  width: 20%;
  top: 10%;

  ${({ align: float }) => (float == 'left' ? 'left: 0;' : null)}
  ${({ align: float }) => (float == 'right' ? 'right: 0;' : null)}

  & svg {
    margin: auto;
    width: 48px;
    height: 64px;
  }
`

type PresentNavigationOverlayProps = {
  children?: React.ReactChild
  dispatchMedia: React.Dispatch<GalleryAction>
  disableSaveCloseInHistory?: boolean
  onInfoClick?: () => void
  media: MediaGalleryFields
}

const LikesButtonWrapper = styled(OverlayButton)`
  bottom: 28px;
  left: 28px;
  width: auto;
`

type LikesIconProps = {
  liked: boolean
  likes: number
}

export const LikesIcon = ({ liked, likes }: LikesIconProps) => {
  return (
    <span style={{ verticalAlign: 'middle', display: 'inline-flex', alignItems: 'start', gap: '10px', height: '32px' }}>
      <svg
        width="1em"
        height="1em"
        viewBox="0 0 19 17"
        version="1.1"
      >
        <path
          d="M13.999086,1 C15.0573371,1 16.0710089,1.43342987 16.8190212,2.20112483 C17.5765039,2.97781012 18,4.03198704 18,5.13009709 C18,6.22820714 17.5765039,7.28238406 16.8188574,8.05923734 L16.8188574,8.05923734 L15.8553647,9.04761889 L9.49975689,15.5674041 L3.14414912,9.04761889 L2.18065643,8.05923735 C1.39216493,7.2503776 0.999999992,6.18971057 1,5.13009711 C1.00000001,4.07048366 1.39216496,3.00981663 2.18065647,2.20095689 C2.95931483,1.40218431 3.97927681,1.00049878 5.00042783,1.00049878 C6.02157882,1.00049878 7.04154078,1.4021843 7.82019912,2.20095684 L7.82019912,2.20095684 L9.4997569,3.92390079 L11.1794784,2.20078881 C11.9271631,1.43342987 12.9408349,1 13.999086,1 L13.999086,1 Z"
          fill={liked ? 'undefined' : 'none'}
          stroke="#000"
          strokeWidth={liked ? '0' : '2'}
          className={(liked ? 'filled' : '') + ' ' + (liked ? 'disabled' : '')}
        ></path>
      </svg>
      {likes > 0 ?
      <span
        style={{lineHeight: '32px', height: '32px'}}
        className={liked ? 'disabled' : ''}
      >{likes}
      </span> : null}
    </span>
  )
}

type LikesButtonProps = {
  liked: boolean
  likes: number
  hide: boolean
  onClick(e: React.MouseEvent<HTMLButtonElement, MouseEvent>): void
}

const LikesButton = ({ liked, likes, hide, onClick }: LikesButtonProps) => {
  return (
    <LikesButtonWrapper
      onClick={onClick}
      className={hide ? 'hide' : undefined}
    >
      <LikesIcon liked={liked} likes={likes} />
    </LikesButtonWrapper>
  )
}


const PresentNavigationOverlay = ({
  children,
  dispatchMedia,
  disableSaveCloseInHistory,
  onInfoClick,
  media,
}: PresentNavigationOverlayProps) => {
  const [hide, setHide] = useState(true)
  const onMouseMove = useRef<null | DebouncedFn<() => void>>(null)
  const [likeMedia] = useLikeMediaMutation()

  useEffect(() => {
    onMouseMove.current = debounce(
      () => {
        setHide(hide => !hide)
      },
      2000,
      true
    )

    return () => {
      onMouseMove.current?.cancel()
    }
  }, [])

  useEffect(() => {
    // Hidden controls are taken out of the tab order entirely, so a keyboard
    // user has no way to reach them - and no pointer to move. Any key press
    // reveals them, the same as a mouse movement or a tap, which then puts
    // them back within reach of Tab.
    const revealOnKey = () => {
      onMouseMove.current && onMouseMove.current()
    }

    document.addEventListener('keydown', revealOnKey)

    return () => {
      document.removeEventListener('keydown', revealOnKey)
    }
  }, [])

  const handlers = useSwipeable({
    onSwipedLeft: () => dispatchMedia({ type: 'nextImage' }),
    onSwipedRight: () => dispatchMedia({ type: 'previousImage' }),
    preventScrollOnSwipe: true,
    trackMouse: false,
  })

  const liked = isMediaLiked(media.id)
  const likes = media.likes ? media.likes : 0
  
  return (
    <StyledOverlayContainer
      data-testid="present-overlay"
      onMouseMove={() => {
        onMouseMove.current && onMouseMove.current()
      }}
      onClick={() => {
        // Touch devices never fire mousemove, so the controls would
        // otherwise stay hidden forever. A tap toggles them the same way a
        // mouse movement does, and reuses the same auto-hide timer.
        onMouseMove.current && onMouseMove.current()
      }}
    >
      <div {...handlers}>
        {children}
        <NavigationButton
          aria-label="Previous image"
          className={hide ? 'hide' : undefined}
          align="left"
          onClick={() => dispatchMedia({ type: 'previousImage' })}
        >
          <PrevIcon />
        </NavigationButton>
        <NavigationButton
          aria-label="Next image"
          className={hide ? 'hide' : undefined}
          align="right"
          onClick={() => dispatchMedia({ type: 'nextImage' })}
        >
          <NextIcon />
        </NavigationButton>
        <ExitButton
          aria-label="Exit presentation mode"
          className={hide ? 'hide' : undefined}
          onClick={() => {
            if (disableSaveCloseInHistory === true) {
              dispatchMedia({ type: 'closePresentMode' })
            } else {
              closePresentModeAction({ dispatchMedia })
            }
          }}
        >
          <ExitIcon />
        </ExitButton>
        {onInfoClick && (
          <InfoButton
            aria-label="Show media info"
            className={hide ? 'hide' : undefined}
            onClick={onInfoClick}
          >
            <InfoIcon />
          </InfoButton>
        )}
        <LikesButton
          liked={liked}
          likes={likes}
          hide={hide}
          onClick={e => {
            e.stopPropagation()
            likePressed({
              media,
              likeMedia,
            })
          }}
        ></LikesButton>
      </div>
    </StyledOverlayContainer>
  )
}

export default PresentNavigationOverlay
