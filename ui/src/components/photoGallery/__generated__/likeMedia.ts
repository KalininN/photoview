/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: likeMedia
// ====================================================

export interface likeMedia_likeMedia {
  __typename: "Media";
  id: string;
  /**
   * Number of likes
   */
  likes: number | null;
}

export interface likeMedia {
  /**
   * Add 1 to likes count for a media
   */
  likeMedia: likeMedia_likeMedia;
}

export interface likeMediaVariables {
  mediaId: string;
}
