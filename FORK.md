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
and adds a minimum-rating filter (default 5) to the album, timeline and shared-album pages.
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
- `api/graphql/resolvers/album.graphql`, `media.graphql`, `timeline.graphql`: `minRating` arguments, `rating` fields
- `api/graphql/resolvers/album.go`, `timeline.go`, `models/actions/timeline_actions.go`: `FilterMinRating` calls
- `api/scanner/externaltools/exif/exif.go`: reads the `exiftool.Rating` group
- `api/scanner/scanner_tasks/scanner_tasks.go`: registers `RatingTask` after `ExifTask`
- UI: `AlbumFilter.tsx`, `AlbumGallery.tsx`, `AlbumPage.tsx`, `AlbumSharePage.tsx`, `SharePage.tsx`,
  `TimelineFilters.tsx`, `TimelineGallery.tsx`, `MediaSidebar.tsx`, `MediaSidebarExif.tsx` and two test fixtures

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
