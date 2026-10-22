/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: setGroupLabel
// ====================================================

export interface setGroupLabel_setFaceGroupLabel {
  __typename: "FaceGroup";
  id: string;
  /**
   * The display name
   */
  displayLabel: string | null;
}

export interface setGroupLabel {
  /**
   * Assign a label to a face group, set label to null to remove the current one
   */
  setFaceGroupLabel: setGroupLabel_setFaceGroupLabel;
}

export interface setGroupLabelVariables {
  groupID: string;
  displayLabel?: string | null;
}
