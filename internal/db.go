package murult

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	db *sql.DB
}

func NewDb() *Db {
	db, err := sql.Open("sqlite3", "murult.db")

	if err != nil {
		Logger.Printf("Unable to open sqlite3 DB because '%s'\n", err)
		return nil
	}

	return &Db{
		db: db,
	}
}

func (db *Db) CreateChannelsTable() {
	_, err := db.db.Exec("CREATE TABLE IF NOT EXISTS channels (channelId TEXT NOT NULL, guildId TEXT NOT NULL, PRIMARY KEY(channelId, guildId))")

	if err != nil {
		Logger.Printf("Unable to create channels table because '%s'\n", err)
	}
}

func (db *Db) CreateRegionsTable() {
	_, err := db.db.Exec("CREATE TABLE IF NOT EXISTS regions (channelId TEXT NOT NULL, region TEXT NOT NULL, PRIMARY KEY(channelId, region))")

	if err != nil {
		Logger.Printf("Unable to create regions table because '%s'\n", err)
	}
}

func (db *Db) CreateDutiesTable() {
	_, err := db.db.Exec("CREATE TABLE IF NOT EXISTS duties (channelId TEXT NOT NULL, name TEXT NOT NULL, PRIMARY KEY(channelId, name))")

	if err != nil {
		Logger.Printf("Unable to create duties table because '%s'\n", err)
	}
}

func (db *Db) CreatePostsTable() {
	_, err := db.db.Exec("CREATE TABLE IF NOT EXISTS posts (channelId TEXT NOT NULL, messageId TEXT NOT NULL, creator TEXT NOT NULL, PRIMARY KEY(channelId, messageId, creator))")

	if err != nil {
		Logger.Printf("Unable to create posts table because '%s'\n", err)
	}
}

func (db *Db) InsertChannel(guildId, channelId string) bool {
	_, err := db.db.Exec("INSERT INTO channels (guildId, channelId) VALUES (?, ?)", guildId, channelId)

	if err != nil {
		Logger.Printf("Unable to insert into channels table because '%s'\n", err)
		return false
	}

	return true
}

func (db *Db) RemoveChannel(guildId, channelId string) bool {
	_, err := db.db.Exec("DELETE FROM posts WHERE channelId=?", channelId)

	if err != nil {
		Logger.Printf("Unable to clear posts with channelId='%s' because '%s'\n", channelId, err)
		return false
	}

	_, err = db.db.Exec("DELETE FROM duties WHERE channelId=?", channelId)

	if err != nil {
		Logger.Printf("Unable to clear duties with channelId='%s' because '%s'\n", channelId, err)
		return false
	}

	_, err = db.db.Exec("DELETE FROM channels WHERE guildId=? AND channelId=?", guildId, channelId)

	if err != nil {
		Logger.Printf("Unable to remove from channels table because '%s'\n", err)
		return false
	}

	return true
}

func (db *Db) InsertRegion(channelId string, region Region) bool {
	_, err := db.db.Exec("INSERT INTO regions (channelId, region) VALUES (?, ?)", channelId, region)

	if err != nil {
		Logger.Printf("Unable to insert into regions table because '%s'\n", err)
		return false
	}

	return true
}

func (db *Db) RemoveRegion(channelId, region Region) bool {
	_, err := db.db.Exec("DELETE FROM regions WHERE channelId=? AND region=?", channelId, region)

	if err != nil {
		Logger.Printf("Unable to remove from regions table because '%s'\n", err)
		return false
	}

	return true
}

func (db *Db) InsertDuty(channelId, name string) bool {
	_, err := db.db.Exec("INSERT INTO duties (channelId, name) VALUES (?, ?)", channelId, name)

	if err != nil {
		Logger.Printf("Unable to insert into duties table because '%s'\n", err)
		return false
	}

	return true
}

func (db *Db) RemoveDuty(channelId, name string) bool {
	_, err := db.db.Exec("DELETE FROM duties WHERE channelId=? AND name=?", channelId, name)

	if err != nil {
		Logger.Printf("Unable to remove from duties table because '%s'\n", err)
		return false
	}

	return true
}

func (db *Db) InsertPost(channelId, messageId, creator string) bool {
	_, err := db.db.Exec("INSERT INTO posts (channelId, messageId, creator) VALUES (?, ?, ?)", channelId, messageId, creator)

	if err != nil {
		Logger.Printf("Unable to insert into posts table because '%s'\n", err)
		return false
	}

	return true
}

func (db *Db) RemovePost(channelId, messageId, creator string) bool {
	_, err := db.db.Exec("DELETE FROM posts WHERE channelId=? AND messageId=? AND creator=?", channelId, messageId, creator)

	if err != nil {
		Logger.Printf("Unable to remove from posts table because '%s'\n", err)
		return false
	}

	return true
}

func (db *Db) SelectAllChannels() (map[string]*Channel, bool) {
	channels := make(map[string]*Channel, 0)
	dbChannels, err := db.db.Query("SELECT channelId, guildId FROM channels")

	if err != nil {
		Logger.Printf("Unable to query for channels: '%s'\n", err)
		return channels, false
	}

	for dbChannels.Next() {
		channelId := ""
		guildId := ""
		regions := make(map[string]struct{}, 0)
		duties := make(map[string]struct{}, 0)
		posts := make(map[string]*Post, 0)
		err = dbChannels.Scan(&channelId, &guildId)

		if err != nil {
			Logger.Printf("Unable to read from a row: '%s'\n", err)
			return channels, false
		}

		// Regions
		{
			dbRegions, err := db.db.Query("SELECT region FROM regions WHERE channelId=?", channelId)

			if err != nil {
				Logger.Printf("Unable to query for regions: '%s'\n", err)
				return channels, false
			}

			for dbRegions.Next() {
				region := ""
				err = dbRegions.Scan(&region)

				if err != nil {
					Logger.Printf("Unable to read from a row: '%s'\n", err)
					return channels, false
				}

				regions[region] = struct{}{}
			}
		}

		// Duties
		{
			dbDuties, err := db.db.Query("SELECT name FROM duties WHERE channelId=?", channelId)

			if err != nil {
				Logger.Printf("Unable to query for duties: '%s'\n", err)
				return channels, false
			}

			for dbDuties.Next() {
				duty := ""
				err = dbDuties.Scan(&duty)

				if err != nil {
					Logger.Printf("Unable to read from a row: '%s'\n", err)
					return channels, false
				}

				duties[duty] = struct{}{}
			}
		}

		// Posts
		{
			dbPosts, err := db.db.Query("SELECT messageId, creator FROM posts WHERE channelId=?", channelId)

			if err != nil {
				Logger.Printf("Unable to query for posts: '%s'\n", err)
				return channels, false
			}

			for dbPosts.Next() {
				messageId := ""
				creator := ""
				err = dbPosts.Scan(&messageId, &creator)

				if err != nil {
					Logger.Printf("Unable to read from a row: '%s'\n", err)
					return channels, false
				}

				post := NewPost()
				post.MessageId = messageId
				post.Creator = creator
				posts[creator] = post
			}
		}

		channels[channelId] = NewChannel(guildId, channelId, regions, duties, posts)
	}

	return channels, true
}

func (db *Db) Close() error {
	return db.db.Close()
}
