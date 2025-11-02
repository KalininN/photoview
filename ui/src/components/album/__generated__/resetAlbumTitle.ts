/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: resetAlbumTitle
// ====================================================

export interface resetAlbumTitle_resetAlbumTitle {
  __typename: "Album";
  id: string;
}

export interface resetAlbumTitle {
  /**
   * Reset a curtom album title
   */
  resetAlbumTitle: resetAlbumTitle_resetAlbumTitle;
}

export interface resetAlbumTitleVariables {
  albumID: string;
}
