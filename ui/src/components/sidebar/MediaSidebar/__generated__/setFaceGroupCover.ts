/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: setFaceGroupCover
// ====================================================

export interface setFaceGroupCover_setFaceGroupCover {
  __typename: "FaceGroup";
  id: string;
}

export interface setFaceGroupCover {
  /**
   * Assign a cover photo to a face group, set id to null to remove the current one
   */
  setFaceGroupCover: setFaceGroupCover_setFaceGroupCover;
}

export interface setFaceGroupCoverVariables {
  faceGroupID: string;
  faceID: string;
}
