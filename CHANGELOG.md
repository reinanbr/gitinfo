# Changelog

All notable changes to this project will be documented in this file.

## v0.4.0

- Add `GetProfile`, combining identity, repo/follower/star counts, total commits, and streaks into a single concurrent call.
- Add `GetTotalStars`, summing `stargazerCount` across a user's public repositories.
- Extend `UserInfo` with `Location`, `Company`, `WebsiteUrl`, `TwitterUsername`, and `IsHireable`.
- Extend `RepoNode` with `StargazerCount`.
- Note: GitHub's API doesn't expose a paid-plan ("Pro") flag for arbitrary users, only for the authenticated viewer — `ProfileInfo` has no such field.

## v0.3.3

- Fix `GetStreaks` reporting `current_streak: 0` when today has no contribution yet. The current streak now stays valid through a one-day grace period (last contribution today or yesterday) instead of breaking before the day is even over.
- Add unit tests for `GetContributionStreaks` covering the grace period, a genuinely broken streak, and a same-day contribution.

## v0.3.1

- Add `GetUserInfo` usage to README quick start and API reference.
- Add `GetUserInfo` example and user CLI example.

## v0.3.0

- Release v0.3.0.
- Refactor streaks to use typed `StreakResponse` and update README.
- Remove tracked binaries.

## v0.2.0

- Refactor tests to use typed `StreakResponse`.

## v0.1.3

- Improve README.

## v0.1.2

- Add `example_test.go` for pkg.go.dev.

## v0.1.1

- Add LICENSE.

## v0.1.0

- Initial release.
