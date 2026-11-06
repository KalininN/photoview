/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL query operation: guesserEnabled
// ====================================================

export interface guesserEnabled_siteInfo {
  __typename: 'SiteInfo'
  /**
   * Whether the 4D-guesser is enabled, false when PHOTOVIEW_DISABLE_4D_GUESSER is set
   */
  guesserEnabled: boolean
}

export interface guesserEnabled {
  siteInfo: guesserEnabled_siteInfo
}
