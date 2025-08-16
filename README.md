# Gator RSS Feed Aggregator

A command-line RSS feed aggregator built in Go that allows you to subscribe to RSS feeds, fetch posts, and manage your feed subscriptions.

## Prerequisites

Before running Gator, make sure you have the following installed:

- **Go** (version 1.19 or higher) - [Install Go](https://golang.org/doc/install)
- **PostgreSQL** - [Install PostgreSQL](https://www.postgresql.org/download/)

## Installation

Install the `gator` CLI tool using Go:

```bash
go install github.com/sleklere/gator@latest
```

Make sure your `$GOPATH/bin` is in your system's `PATH` so you can run the `gator` command from anywhere.

## Setup

### 1. Database Setup

First, create a PostgreSQL database for Gator:

```sql
CREATE DATABASE gator;
```

### 2. Configuration File

Create a configuration file at `~/.gatorconfig.json` with your database connection details:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable"
}
```

Replace `username` and `password` with your PostgreSQL credentials.

### 3. Run Database Migrations

The application will automatically run database migrations when you first use it, but you can also run them manually using goose if needed.

## Usage

Here are the main commands you can use with Gator:

### User Management

```bash
# Register a new user
gator register <username>

# Login as a user
gator login <username>

# View current user
gator users

# Reset database (delete all users and data)
gator reset
```

### Feed Management

```bash
# Add a new RSS feed
gator addfeed <name> <url>

# View all available feeds
gator feeds

# Follow a feed
gator follow <feed_url>

# Unfollow a feed
gator unfollow <feed_url>

# View feeds you're following
gator following
```

### Reading Posts

```bash
# Browse recent posts from your followed feeds
gator browse [limit]

# Manually fetch latest posts from all feeds
gator agg
```

## Examples

1. **Register and login:**
   ```bash
   gator register john
   gator login john
   ```

2. **Add and follow some feeds:**
   ```bash
   gator addfeed "TechCrunch" "https://techcrunch.com/feed/"
   gator follow "https://techcrunch.com/feed/"
   ```

3. **Fetch and read posts:**
   ```bash
   gator agg
   gator browse 5
   ```

## Project Structure

- **Database migrations** - Located in the `sql/` directory
- **CLI commands** - Each command is implemented as a separate handler
- **RSS parsing** - Automatically fetches and parses RSS feeds
- **User system** - Multi-user support with individual feed subscriptions

