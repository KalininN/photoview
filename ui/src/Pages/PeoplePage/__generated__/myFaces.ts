/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL query operation: myFaces
// ====================================================

export interface myFaces_myFaceGroups_thumbnail_rectangle {
  __typename: "FaceRectangle";
  minX: number;
  maxX: number;
  minY: number;
  maxY: number;
}

export interface myFaces_myFaceGroups_thumbnail_media_thumbnail {
  __typename: "MediaURL";
  /**
   * URL for previewing the image
   */
  url: string;
  /**
   * Width of the image in pixels
   */
  width: number;
  /**
   * Height of the image in pixels
   */
  height: number;
}

export interface myFaces_myFaceGroups_thumbnail_media {
  __typename: "Media";
  id: string;
  title: string;
  /**
   * URL to display the media in a smaller resolution
   */
  thumbnail: myFaces_myFaceGroups_thumbnail_media_thumbnail | null;
}

export interface myFaces_myFaceGroups_thumbnail {
  __typename: "ImageFace";
  id: string;
  /**
   * A bounding box of where on the image the face is present
   */
  rectangle: myFaces_myFaceGroups_thumbnail_rectangle;
  /**
   * A reference to the image the face appears on
   */
  media: myFaces_myFaceGroups_thumbnail_media;
}

export interface myFaces_myFaceGroups {
  __typename: "FaceGroup";
  id: string;
  /**
   * The name of the person
   */
  label: string | null;
  /**
   * The display name
   */
  displayLabel: string | null;
  /**
   * The id of the cover photo
   */
  thumbnail: myFaces_myFaceGroups_thumbnail;
  /**
   * The total number of images in this collection
   */
  imageFaceCount: number;
}

export interface myFaces {
  /**
   * Get a list of `FaceGroup`s for the logged in user
   */
  myFaceGroups: myFaces_myFaceGroups[];
}

export interface myFacesVariables {
  limit?: number | null;
  offset?: number | null;
}
