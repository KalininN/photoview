# Fork notes

This repository is a fork of [photoview/photoview](https://github.com/photoview/photoview).
Core updates come from upstream; this file lists what the fork adds and how to sync.

## Branches

- `origin-master`: an exact mirror of upstream `master`. Only fast-forwarded, never committed to.
- `master`: `origin-master` plus the fork's changes listed below.
- Fixes that are useful upstream are branched from `origin-master` and sent upstream as PRs
  (for example [#1260](https://github.com/photoview/photoview/pull/1260)).

## Fork changes

### Star rating

Reads the EXIF/XMP `Rating` tag, stores it as `media.rating` (default 3), shows it in the sidebar,
and adds a minimum-rating filter to the album, shared-album, face and timeline pages
("best only" = 5, "story" = 4, "all" = 3). The filter defaults to 5 on the timeline and 4 elsewhere,
is clamped to 3–5, and is not applied when the API gets no `minRating`. The map shows media rated 4+.
When a file changes on disk, its EXIF data and rating are read again on the next scan.

Fork-only files:

- `api/scanner/externaltools/exiftool/rating.go`
- `api/scanner/scanner_tasks/rating_task.go`, `rating_task_test.go`
- `api/graphql/models/rating.go`
- `api/graphql/models/actions/timeline_rating_test.go`
- `api/scanner/externaltools/exif/rating_test.go`, `test_data/rated*.jpg`
- `ui/src/components/album/RatingDropdown.tsx`

Hooks in upstream files:

- `api/graphql/models/media.go`, `media_exif.go`: `Rating` fields
- `api/graphql/resolvers/album.graphql`, `media.graphql`, `timeline.graphql`, `faces.graphql`: `minRating` arguments, `rating` fields
- `api/graphql/resolvers/album.go`, `timeline.go`, `faces.go`, `media_geo_json.go`, `models/actions/timeline_actions.go`:
  `FilterMinRating` calls
- `api/scanner/externaltools/exif/exif.go`: reads the `exiftool.Rating` group
- `api/scanner/scanner_tasks/scanner_tasks.go`: registers `RatingTask` after `ExifTask`
- UI: `AlbumFilter.tsx`, `AlbumGallery.tsx`, `AlbumPage.tsx`, `AlbumSharePage.tsx`, `SharePage.tsx`,
  `TimelineFilters.tsx`, `TimelineGallery.tsx`, `MediaSidebar.tsx`, `MediaSidebarExif.tsx` and two test fixtures

### Faces from metadata

Built with the `faces_exif` tag (the Dockerfile and `dev-compose.yaml` set it), Photoview reads named
face regions (XMP MWG `RegionName`/`RegionArea*`) with exiftool instead of detecting faces. Each name
becomes a face group, with a display label (the first name by default) and a cover face. Regions
without a name are skipped. Faces are read when a photo is found or changes on disk, by `RatingTask`.
Without the tag, upstream's face detector works as before.

Fork-only files:

- `api/scanner/face_detection/face_detector_exif.go`, `face_detector_not_exif.go`, `face_regions.go` (and test)
- `api/scanner/externaltools/exif/face_regions.go` (and test, `test_data/faces.jpg`): reads the region tags
  with the scanner's exiftool
- `api/scanner/scanner_tasks/rating_task_faces_test.go`: only runs with `go test -tags faces_exif`
- `api/database/migrations/face_group_display_label.go` (and test)
- `ui/src/primitives/LabelEditable.tsx`

Hooks in upstream files:

- `api/scanner/face_detection/face_detector_impl.go`: `!faces_exif` build tag
- `api/scanner/scanner_tasks/face_detection_task.go`: skipped when faces come from metadata
- `api/graphql/models/face_detection.go`: `DisplayLabel`, `CoverID` fields
- `api/graphql/resolvers/faces.graphql`, `faces.go`: `displayLabel`, `thumbnail`, `setFaceGroupCover`;
  `setFaceGroupLabel` sets the display label
- `api/database/database.go`: runs the display label migration
- UI: the people pages, `FaceCircleImage.tsx`, `MediaSidebarPeople.tsx` and their tests

### Likes

A `likes` counter on media, the `likeMedia` mutation, and like buttons on thumbnails, the timeline and
the present view. The browser remembers what it liked, so each browser likes a photo once.

- Fork-only: `api/graphql/models/actions/like_actions.go` (and test), `ui/src/components/media/HoverButtons.tsx`
- Hooks: `models/media.go` (`Likes`), `media.graphql`, `resolvers/media.go`, `photoGalleryMutations.ts`,
  `MediaGallery.tsx`, `MediaThumbnail.tsx`, `TimelineGroupAlbum.tsx`, `PresentNavigationOverlay.tsx`,
  `PresentView.tsx`, gallery reducers and test fixtures

### Album titles and covers

Admins can rename an album (`setAlbumTitle`, `resetAlbumTitle`). An album without a cover uses its
sub-albums' cover. The albums page sorts by title in reverse order by default.

- Fork-only: `api/graphql/models/actions/album_title_actions.go` (and test)
- Hooks: `album.graphql`, `resolvers/album.go`, `models/album.go` (`Thumbnail`), `AlbumTitle.tsx`,
  `AlbumsPage.tsx`, `useOrderingParams.ts` (default direction), `album_actions_test.go`

### Admin-only management

Setting album covers and titles, editing face groups, and deleting, protecting, labelling or setting
the expiry of share tokens need an admin (`@isAdmin` in `album.graphql`, `faces.graphql`,
`share_token.graphql`). The UI hides the album and face controls from other users; the share token
controls in the sharing sidebar are still shown to them, and fail.

### Smaller UI changes

- The favorites filter is hidden (`AlbumFilter.tsx`, `TimelineFilters.tsx`).
- The sidebar shows the thumbnail, and closes on page changes unless pinned (`MediaSidebar.tsx`, `Sidebar.tsx`).

### GPS longitude up to ±180 in the GPS migration

`api/database/migrations/exif_invalid_gps.go` (and its test) only clears longitudes beyond ±180.
Upstream still clears everything beyond ±90 on every startup, which deletes valid coordinates;
this is worth sending upstream. (The scanner side of this fix is upstream as
[#1260](https://github.com/photoview/photoview/pull/1260).)

### Build script

`build_and_push.sh` regenerates the GraphQL code, builds the Docker image and pushes it as
`kalininn/photoview:latest`.

## Syncing with upstream

```bash
git remote add upstream https://github.com/photoview/photoview.git   # once
git config rerere.enabled true                                        # once
git fetch upstream
git push origin upstream/master:origin-master
git switch -c sync/upstream-$(date +%F) master
git merge origin-master
# Generated files: take either side, then regenerate
(cd api && go tool gqlgen generate)
(cd ui && npm run genSchemaTypes)   # keep only the changes to the fork's queries
# Build, run the tests, open a PR into master
```

Keep new parameters at the end of signatures and argument lists, and keep fork logic in the
fork-only files above, so upstream changes merge without conflicts.
