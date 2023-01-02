package murult

import (
	"database/sql"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	db *sql.DB
}

func NewDb(path string) *Db {
	db, err := sql.Open("sqlite3", path)

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
	_, err := db.db.Exec("CREATE TABLE IF NOT EXISTS posts (dc TEXT NOT NULL, duty TEXT NOT NULL, tags TEXT NOT NULL, description TEXT NOT NULL, creator TEXT NOT NULL, expTime INTEGER NOT NULL, updTime INTEGER NOT NULL, party TEXT NOT NULL, channelId TEXT NOT NULL, messageId TEXT NOT NULL, PRIMARY KEY(channelId, messageId))")

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

	Logger.Printf("Inserted (guildId='%s', channelId='%s') into channels table\n", guildId, channelId)
	return true
}

func (db *Db) RemoveChannel(guildId, channelId string) bool {
	_, err := db.db.Exec("DELETE FROM posts WHERE channelId=?", channelId)

	if err != nil {
		Logger.Printf("Unable to clear posts with channelId='%s' because '%s'\n", channelId, err)
		return false
	}
	Logger.Printf("Removed all posts related to channelId='%s'\n", channelId)

	_, err = db.db.Exec("DELETE FROM duties WHERE channelId=?", channelId)

	if err != nil {
		Logger.Printf("Unable to clear duties with channelId='%s' because '%s'\n", channelId, err)
		return false
	}
	Logger.Printf("Removed all duties related to channelId='%s'\n", channelId)

	_, err = db.db.Exec("DELETE FROM regions WHERE channelId=?", channelId)

	if err != nil {
		Logger.Printf("Unable to clear regions with channelId='%s' because '%s'\n", channelId, err)
		return false
	}
	Logger.Printf("Removed all regions related to channelId='%s'\n", channelId)

	_, err = db.db.Exec("DELETE FROM channels WHERE guildId=? AND channelId=?", guildId, channelId)

	if err != nil {
		Logger.Printf("Unable to remove from channels table because '%s'\n", err)
		return false
	}
	Logger.Printf("Removed all channels related to channelId='%s'\n", channelId)

	return true
}

func (db *Db) InsertRegion(channelId string, region Region) bool {
	_, err := db.db.Exec("INSERT INTO regions (channelId, region) VALUES (?, ?)", channelId, region)

	if err != nil {
		Logger.Printf("Unable to insert into regions table because '%s'\n", err)
		return false
	}

	Logger.Printf("Inserted (channelId='%s', region='%s') into regions table\n", channelId, region)
	return true
}

func (db *Db) RemoveRegion(channelId, region Region) bool {
	_, err := db.db.Exec("DELETE FROM regions WHERE channelId=? AND region=?", channelId, region)

	if err != nil {
		Logger.Printf("Unable to remove from regions table because '%s'\n", err)
		return false
	}

	Logger.Printf("Removed (channelId='%s', region='%s') from regions table\n", channelId, region)
	return true
}

func (db *Db) InsertDuty(channelId, name string) bool {
	_, err := db.db.Exec("INSERT INTO duties (channelId, name) VALUES (?, ?)", channelId, name)

	if err != nil {
		Logger.Printf("Unable to insert into duties table because '%s'\n", err)
		return false
	}

	Logger.Printf("Inserted (channelId='%s', name='%s') into duties table\n", channelId, name)
	return true
}

func (db *Db) RemoveDuty(channelId, name string) bool {
	_, err := db.db.Exec("DELETE FROM duties WHERE channelId=? AND name=?", channelId, name)

	if err != nil {
		Logger.Printf("Unable to remove from duties table because '%s'\n", err)
		return false
	}

	Logger.Printf("Removed (channelId='%s', name='%s') from duties table\n", channelId, name)
	return true
}

func (db *Db) InsertPost(p *Post) bool {
	_, err := db.db.Exec(
		`INSERT INTO posts (
			dc, duty, tags, description, creator, expTime, updTime, party, channelId, messageId
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		) ON CONFLICT (channelId, messageId) DO 
		UPDATE SET 
			dc = ?, duty = ?, tags = ?, description = ?, creator = ?, expTime = ?, updTime = ?, party = ?`,
		p.DataCentre, p.Duty, p.Tags, p.Description, p.Creator, p.Expires.Unix(), p.Updated.Unix(), serializeParty(p.Slots), p.ChannelId, p.MessageId,
		p.DataCentre, p.Duty, p.Tags, p.Description, p.Creator, p.Expires.Unix(), p.Updated.Unix(), serializeParty(p.Slots))

	if err != nil {
		Logger.Printf("Unable to insert into posts table because '%s'\n", err)
		return false
	}

	Logger.Printf("Inserted '%#v' into the posts table\n", p)
	return true
}

func (db *Db) RemovePost(p *Post) bool {
	_, err := db.db.Exec("DELETE FROM posts WHERE channelId=? AND messageId=?", p.ChannelId, p.MessageId)

	if err != nil {
		Logger.Printf("Unable to remove from posts table because '%s'\n", err)
		return false
	}
	Logger.Printf("Removed '%#v' from the posts table\n", p)
	return true
}

func (db *Db) SelectAllChannels() (map[string]*Channel, bool) {
	channels := make(map[string]*Channel, 0)
	dbChannels, err := db.db.Query("SELECT channelId, guildId FROM channels")

	if err != nil {
		Logger.Printf("Unable to query for channels: '%s'\n", err)
		return map[string]*Channel{}, false
	}

	for dbChannels.Next() {
		channelId := ""
		guildId := ""
		err = dbChannels.Scan(&channelId, &guildId)

		if err != nil {
			Logger.Printf("Unable to read from a row: '%s'\n", err)
			return map[string]*Channel{}, false
		}

		regions, ok1 := db.SelectAllRegionsForChannel(channelId)
		duties, ok2 := db.SelectAllDutiesForChannel(channelId)
		posts, ok3 := db.SelectAllPostsFromChannel(channelId)

		if !ok1 || !ok2 || !ok3 {
			return map[string]*Channel{}, false
		}

		channels[channelId] = NewChannel(guildId, channelId, regions, duties, posts)
	}

	return channels, true
}

func (db *Db) SelectAllPostsFromChannel(channelId string) (map[string]*Post, bool) {
	posts := make(map[string]*Post, 0)

	dbPosts, err := db.db.Query("SELECT dc, duty, tags, description, creator, expTime, updTime, party, messageId FROM posts WHERE channelId=?", channelId)

	if err != nil {
		Logger.Printf("Unable to query for posts: '%s'\n", err)
		return map[string]*Post{}, false
	}

	for dbPosts.Next() {
		dc := ""
		duty := ""
		tags := ""
		description := ""
		creator := ""
		var expTimeInt int64
		var updTimeInt int64
		party := ""
		messageId := ""
		err = dbPosts.Scan(
			&dc,
			&duty,
			&tags,
			&description,
			&creator,
			&expTimeInt,
			&updTimeInt,
			&party,
			&messageId,
		)

		if err != nil {
			Logger.Printf("Unable to read from a row: '%s'\n", err)
			return map[string]*Post{}, false
		}

		post := &Post{
			ChannelId:   channelId,
			MessageId:   messageId,
			DataCentre:  dc,
			Duty:        duty,
			Tags:        tags,
			Description: description,
			Creator:     creator,
			Expires:     time.Unix(expTimeInt, 0),
			Updated:     time.Unix(updTimeInt, 0),
			Slots:       deserializeParty(party),
		}

		posts[creator] = post
	}

	return posts, true
}

func (db *Db) SelectAllDutiesForChannel(channelId string) (map[string]struct{}, bool) {
	duties := make(map[string]struct{}, 0)
	dbDuties, err := db.db.Query("SELECT name FROM duties WHERE channelId=?", channelId)

	if err != nil {
		Logger.Printf("Unable to query for duties: '%s'\n", err)
		return map[string]struct{}{}, false
	}

	for dbDuties.Next() {
		duty := ""
		err = dbDuties.Scan(&duty)

		if err != nil {
			Logger.Printf("Unable to read from a row: '%s'\n", err)
			return map[string]struct{}{}, false
		}

		duties[duty] = struct{}{}
	}

	return duties, true
}

func (db *Db) SelectAllRegionsForChannel(channelId string) (map[string]struct{}, bool) {
	regions := make(map[string]struct{}, 0)
	dbRegions, err := db.db.Query("SELECT region FROM regions WHERE channelId=?", channelId)

	if err != nil {
		Logger.Printf("Unable to query for regions: '%s'\n", err)
		return map[string]struct{}{}, false
	}

	for dbRegions.Next() {
		region := ""
		err = dbRegions.Scan(&region)

		if err != nil {
			Logger.Printf("Unable to read from a row: '%s'\n", err)
			return map[string]struct{}{}, false
		}

		regions[region] = struct{}{}
	}

	return regions, true
}

func (db *Db) Close() error {
	return db.db.Close()
}

func deserializeParty(party string) []string {
	return strings.Split(party, ",")
}

func serializeParty(party []string) string {
	return strings.Join(party, ",")
}
