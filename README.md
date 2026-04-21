# rss-aggregator

A command-line RSS feed aggregator written in Go. Register users, subscribe to feeds, run a background aggregator to fetch posts, and browse the latest content.

## Prerequisites

- [Go](https://go.dev/dl/) 1.22 or later
- [PostgreSQL](https://www.postgresql.org/download/) (any recent version)

## Installation

```bash
go install github.com/jasonsoprovich/rss-aggregator@latest
```

This installs the `rss-aggregator` binary to your `$GOPATH/bin`. Make sure that directory is on your `$PATH`.

## Configuration

Create a config file at `~/.gatorconfig.json`:

```json
{
  "db_url": "postgres://username:password@localhost:5432/rss-aggregator?sslmode=disable",
  "current_user_name": ""
}
```

Replace `username`, `password`, and the database name as appropriate for your Postgres setup. The `current_user_name` field is managed automatically by the `login` and `register` commands.

### Database setup

Create the database and run the migrations in order:

```bash
createdb gator
psql gator -f sql/schema/001_users.sql
psql gator -f sql/schema/002_feeds.sql
psql gator -f sql/schema/003_feed_follows.sql
psql gator -f sql/schema/004_feed_lastfetched.sql
psql gator -f sql/schema/005_posts.sql
```

## Usage

### User management

```bash
# Create a new user and log in as them
rss-aggregator register alice

# Switch to an existing user
rss-aggregator login alice

# List all registered users
rss-aggregator users
```

### Managing feeds

```bash
# Add a feed and automatically follow it (requires login)
rss-aggregator addfeed "Hacker News" https://news.ycombinator.com/rss

# List all feeds in the system
rss-aggregator feeds

# Follow an existing feed by URL
rss-aggregator follow https://news.ycombinator.com/rss

# List feeds you're following
rss-aggregator following

# Unfollow a feed
rss-aggregator unfollow https://news.ycombinator.com/rss
```

### Aggregating posts

Run the aggregator in a separate terminal. It fetches feeds on a continuous loop at the interval you specify:

```bash
# Fetch feeds every 30 seconds
rss-aggregator agg 30s

# Fetch feeds every 5 minutes
rss-aggregator agg 5m
```

### Browsing posts

```bash
# Show your 2 most recent posts (default)
rss-aggregator browse

# Show your 10 most recent posts
rss-aggregator browse 10
```

### Reset

```bash
# Delete all users (useful for development/testing)
rss-aggregator reset
```
