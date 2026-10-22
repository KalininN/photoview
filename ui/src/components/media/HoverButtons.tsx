import styled from 'styled-components'

export const HoverIcon = styled.button`
  font-size: 1.3em;
  margin: 10px 10px 10px 10px;
  color: white;
  text-shadow: 0 0 4px black;
  opacity: 0;
  position: relative;

  border-radius: 17px;
  width: auto;
  height: 34px;

  padding: 0px 7px 0px 7px;

  &:hover,
  &:focus {
    background-color: rgba(0, 0, 0, 0.4);
  }

  transition: opacity 100ms, background-color 100ms;
`

type FavoriteIconProps = {
  favorite: boolean
}

export const FavoriteIcon = ({ favorite }: FavoriteIconProps) => {
  return (
    <svg
    className="text-white m-auto mt-1"
    width="19px"
    height="17px"
    viewBox="0 0 19 17"
    version="1.1"
    >
    <path
        d="M13.999086,1 C15.0573371,1 16.0710089,1.43342987 16.8190212,2.20112483 C17.5765039,2.97781012 18,4.03198704 18,5.13009709 C18,6.22820714 17.5765039,7.28238406 16.8188574,8.05923734 L16.8188574,8.05923734 L15.8553647,9.04761889 L9.49975689,15.5674041 L3.14414912,9.04761889 L2.18065643,8.05923735 C1.39216493,7.2503776 0.999999992,6.18971057 1,5.13009711 C1.00000001,4.07048366 1.39216496,3.00981663 2.18065647,2.20095689 C2.95931483,1.40218431 3.97927681,1.00049878 5.00042783,1.00049878 C6.02157882,1.00049878 7.04154078,1.4021843 7.82019912,2.20095684 L7.82019912,2.20095684 L9.4997569,3.92390079 L11.1794784,2.20078881 C11.9271631,1.43342987 12.9408349,1 13.999086,1 L13.999086,1 Z"
        fill={favorite ? 'currentColor' : 'none'}
        stroke="currentColor"
        strokeWidth={favorite ? '0' : '2'}
    ></path>
    </svg>
  )
}

type LikesIconProps = {
  liked: boolean
  likes: number
}

export const LikesIcon = ({ liked, likes }: LikesIconProps) => {
  return (
    <span style={{ verticalAlign: 'middle', marginBottom: '3px', display: 'inline-flex', alignItems: 'start', gap: '6px', height: '17px' }}>
      <svg
        className="text-white m-auto"
        width="19px"
        height="17px"
        viewBox="0 0 19 17"
        version="1.1"
      >
        <path
          d="M13.999086,1 C15.0573371,1 16.0710089,1.43342987 16.8190212,2.20112483 C17.5765039,2.97781012 18,4.03198704 18,5.13009709 C18,6.22820714 17.5765039,7.28238406 16.8188574,8.05923734 L16.8188574,8.05923734 L15.8553647,9.04761889 L9.49975689,15.5674041 L3.14414912,9.04761889 L2.18065643,8.05923735 C1.39216493,7.2503776 0.999999992,6.18971057 1,5.13009711 C1.00000001,4.07048366 1.39216496,3.00981663 2.18065647,2.20095689 C2.95931483,1.40218431 3.97927681,1.00049878 5.00042783,1.00049878 C6.02157882,1.00049878 7.04154078,1.4021843 7.82019912,2.20095684 L7.82019912,2.20095684 L9.4997569,3.92390079 L11.1794784,2.20078881 C11.9271631,1.43342987 12.9408349,1 13.999086,1 L13.999086,1 Z"
          fill={liked ? 'currentColor' : 'none'}
          stroke="currentColor"
          strokeWidth={liked ? '0' : '2'}
        ></path>
      </svg>
      {likes > 0 ?
      <span className="text-white m-auto" style={{lineHeight: '17px', height: '17px'}}>{likes}</span> : null}
    </span>
  )
}

export const SidebarIcon = () => (
    <svg
        width="20px"
        height="20px"
        viewBox="0 0 20 20"
        version="1.1"
        className="m-auto"
    >
        <path
        d="M10,0 C15.5228475,0 20,4.4771525 20,10 C20,15.5228475 15.5228475,20 10,20 C4.4771525,20 0,15.5228475 0,10 C0,4.4771525 4.4771525,0 10,0 Z M10,9 C9.44771525,9 9,9.44771525 9,10 L9,10 L9,14 L9.00672773,14.1166211 C9.06449284,14.6139598 9.48716416,15 10,15 C10.5522847,15 11,14.5522847 11,14 L11,14 L11,10 L10.9932723,9.88337887 C10.9355072,9.38604019 10.5128358,9 10,9 Z M10.01,5 L9.88337887,5.00672773 C9.38604019,5.06449284 9,5.48716416 9,6 C9,6.55228475 9.44771525,7 10,7 L10,7 L10.1266211,6.99327227 C10.6239598,6.93550716 11.01,6.51283584 11.01,6 C11.01,5.44771525 10.5622847,5 10.01,5 L10.01,5 Z"
        fill="#FFFFFF"
        ></path>
    </svg>
)
