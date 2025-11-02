/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: setAlbumTitle
// ====================================================

export interface setAlbumTitle_setAlbumTitle {
  __typename: "Album";
  id: string;
  title: string;
}

export interface setAlbumTitle {
  /**
   * Assign a custom album title
   */
  setAlbumTitle: setAlbumTitle_setAlbumTitle;
}

export interface setAlbumTitleVariables {
  albumID: string;
  title: string;
}
