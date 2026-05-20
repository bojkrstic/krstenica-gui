# Session Summary

Session ID: `019e4480-4367-73e3-a7c7-09c8e1892886`

Date: `2026-05-20`

Branch: `develop`

Last commit: `ce658e8` (`Add in-app user guide`)

Project root: `/home/bojan/develop/horisen/Krstenica/Krstenica-gui/krstenica`

## Completed Work

- Created `USER_GUIDE.md` with a short Cyrillic user guide for the GUI.
- Created `VIZUELNO_UPUTSTVO.html` and exported `VIZUELNO_UPUTSTVO.pdf` as a Cyrillic visual guide with images.
- Added an in-app Cyrillic guide page at `/ui/uputstvo` backed by `web/templates/uputstvo/index.html`.
- Copied the visual guide images into `web/static/guide/` and linked them from inside the app.
- Tightened the in-app guide copy and added mobile-friendly layout rules so the page reads better on phones.
- Added a PDF download link for the visual guide at `/static/guide/VIZUELNO_UPUTSTVO.pdf`.
- Created `AGENTS.md` with project rules, coding style, documentation requirements, and the rule to always maintain `SESSION.md`.
- Committed the in-app guide, visual guide assets, and documentation updates into `ce658e8`.
- Updated `build-and-push.sh` to point at `version1.1.3`.
- Applied a small UI polish pass:
  - larger page titles and stronger card hierarchy
  - more prominent primary buttons
  - denser table rows for faster scanning
  - sticky modal footer styling for form actions
- Made the shared config server-safe again and moved local-only overrides to `config/config.local.yaml` with `.gitignore` protection.
- Verified that the local/server config split works as intended for development and deployment.
- Added local/development API auth shortcut logic in `internal/handler/auth.go`.
- `requireAPIAuth()` now skips API authentication when `config.ENV` is `local`, `dev`, or `development`.
- `requireRole()` now skips role checks in the same local/dev environments.
- Added a helper to centralize the local/development environment check.
- New krstenica form now defaults date fields when the modal opens:
  - `birth_date` uses the current local date and time.
  - `baptism` uses today's local date.
  - `certificate` uses today's local date.
- Date picker trigger was moved from the left side of date inputs to the right side globally.
- Date picker trigger was restyled as a small calendar button instead of a plain/blue square.
- Additional CSS was added to try to vertically center the calendar SVG inside the trigger button.
- Updated `SESSION.md` to the required handoff structure.

## Changed Files

- `USER_GUIDE.md`
- `VIZUELNO_UPUTSTVO.html`
- `VIZUELNO_UPUTSTVO.pdf`
- `AGENTS.md`
- `SESSION.md`
- `.gitignore`
- `config/config.yaml`
- `config/config.local.yaml`
- `internal/handler/auth.go`
- `internal/handler/gui.go`
- `web/templates/krstenice/new.html`
- `web/templates/layouts/base.html`
- `web/templates/users/index.html`
- `web/templates/uputstvo/index.html`
- `web/static/guide/01-prijava.png`
- `web/static/guide/02-pocetni-ekran.png`
- `web/static/guide/03-pregled-krstenica.png`
- `web/static/guide/04-nova-krstenica.png`
- `web/static/guide/05-pregled-krstenica.png`
- `web/static/guide/06-eparhije-izmena.png`
- `web/static/guide/07-eparhija-nova.png`
- `web/static/guide/08-hram-izmena.png`
- `web/static/guide/09-hram-novi.png`
- `web/static/guide/10-svestenik-izmena.png`
- `web/static/guide/11-svestenik-novi.png`
- `web/static/guide/12-osobe-izmena.png`
- `web/static/guide/13-osoba-nova.png`
- `web/static/guide/VIZUELNO_UPUTSTVO.pdf`
- `build-and-push.sh`

## Pending Tasks

- Visually re-check the date picker icon after hard refresh (`Ctrl+F5`), because the user still saw the icon aligned too low after CSS changes.
- If the icon is still low, inspect computed browser styles for `.date-input-icon` and `.date-input-icon svg`; likely a global button style or browser rendering detail is still affecting alignment.
- Confirm the new hierarchy and denser table styling feel balanced in the running GUI.
- Review the `/ui/uputstvo` page in the browser and confirm the sequence, text, and image scaling are acceptable.
- Decide whether the local/development API auth shortcut should remain committed as-is.
- Confirm whether `config/config.local.yaml` should stay local-only or be replaced with an environment-variable-based override.
- Add tests for `requireAPIAuth()`, `requireRole()`, and the local/development environment helper.
- Verify the new krstenica defaults in the running GUI by opening `/ui/krstenice/new` multiple times and confirming the time updates each time.

## Blockers

- `go test ./...` still fails unless local Postgres is running on `127.0.0.1:5560`, because `cmd/krstenica` starts auto migration.
- The sandbox blocks normal Go build cache writes; targeted tests need elevated execution or a writable cache.
- Full visual verification was not possible from terminal-only execution; browser screenshot/user feedback is still needed for final icon positioning.
- Untracked files currently present and not part of these code changes:
  - `2026-05-20_13-05.png`
  - `internal/dto/krstenica-api`
  - `VIZUELNO_UPUTSTVO.html`
- PDF export on the server is blocked unless `doc/template_files/krstenica-template-empty.xlsx` and `doc/template_files/krstenica-template.xlsx` are copied into the runtime image.

## Verification

- `go test ./internal/handler ./internal/dto ./internal/service` passed after allowing normal Go build cache access.
- `go test ./internal/handler` passed after the date picker CSS changes.
- `go test ./internal/handler` passed after the UI hierarchy/table polish pass.
- `go test ./internal/handler` passed after adding the in-app guide page.
- `go test ./internal/handler` passed after the mobile-friendly guide layout update.
- `go test ./internal/handler` passed after adding the PDF download link.
- `go test ./internal/config` passed after the config loader changes.
- `VIZUELNO_UPUTSTVO.pdf` was generated successfully with headless Chrome.
- The visual guide PDF is also available as a direct app download.
- No new tests were run for this documentation update.

## Next Steps

1. Hard refresh the GUI and confirm whether the date picker icon is vertically centered.
2. If still misaligned, replace the current button-based trigger with a non-button wrapper plus explicit click handler, or use a CSS pseudo-element icon to avoid inherited button styles completely.
3. Review the `/ui/uputstvo` page in the browser on desktop and mobile and adjust copy or image order if needed.
4. Start local Postgres on `127.0.0.1:5560` and run `go test ./...`.
5. Decide what to do with untracked files `2026-05-20_13-05.png`, the source `Sl*.png` images, `internal/dto/krstenica-api`, and `VIZUELNO_UPUTSTVO.html`.
6. Review whether the local override should stay file-based or move to an environment-variable override.
